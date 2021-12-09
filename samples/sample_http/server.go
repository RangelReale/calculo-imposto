package main

import (
	"fmt"
	"net/http"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
	"github.com/RangelReale/calculo-imposto/calculo_pf"
	"github.com/RangelReale/calculo-imposto/calculo_pj"
)

func main() {
	http.HandleFunc("/calculo", calculo)
	http.ListenAndServe(":8091", nil)
}

func calculo(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	//faturamentoMensal := 45_000.00
	exportacao := true

	t := NewTemplate(w)

	faturamentoMensalBase := 5000.00
	for i := 1; i < 11; i++ {
		faturamentoMensal := faturamentoMensalBase * float64(i)

		t.Tfmtln(`Faturamento Mensal: %.2f<br/>`, faturamentoMensal)

		t.Tfmtln(`<table border="1">`)
		t.Tfmtln(`<tr><th>Simples Anexo 5</th><th>Simples Anexo 5 Fator R</th><th>Lucro Presumido</th></tr>`)
		t.Tfmtln(`<tr>`)

		t.Tfmtln(`<td valign="top">`)
		// Simples Anexo 5
		err := calculo_SimplesAnexo5(t, faturamentoMensal, exportacao)
		if err != nil {
			t.CheckError(err)
			fmt.Printf("Erro em SimplesAnexo5: %s\n", err.Error())
		}
		t.Tfmtln(`</td>`)

		t.Tfmtln(`<td valign="top">`)
		// Simples Anexo 5 fator R
		err = calculo_SimplesAnexo5FatorR(t, faturamentoMensal, exportacao)
		if err != nil {
			t.CheckError(err)
			fmt.Printf("Erro em calculo_SimplesAnexo5FatorR: %s\n", err.Error())
		}
		t.Tfmtln(`</td>`)

		t.Tfmtln(`<td valign="top">`)
		// Lucro Presumido
		err = calculo_LucroPresumido(t, faturamentoMensal, exportacao)
		if err != nil {
			t.CheckError(err)
			fmt.Printf("Erro em calculo_LucroPresumido: %s\n", err.Error())
		}
		t.Tfmtln(`</td>`)

		t.Tfmtln(`</tr>`)
		t.Tfmtln(`</table>`)
	}
}

func outputCalc(t *Template, res_pj *calculo_imposto.CalculoResultado, res_pf *calculo_imposto.CalculoResultado) {
	t.Tfmtln(`<table border="1">`)
	t.Tfmtln(`<tr><th>PJ</th><th>PF</th><th>%%</th><th>Total</th></tr>`)
	t.Tfmtln(`<tr>`)
	t.Tfmtln(`<td>%.2f</td>`, res_pj.Total().ValorImposto/12.0)
	t.Tfmtln(`<td>%.2f</td>`, res_pf.Total().ValorImposto/12.0)
	t.Tfmtln(`<td>%.2f%%</td>`, (res_pj.Total().ValorImposto+res_pf.Total().ValorImposto)/res_pj.Total().ValorOriginal*100.0)
	t.Tfmtln(`<td>%.2f</td>`, (res_pj.Total().ValorImposto+res_pf.Total().ValorImposto)/12.0)
	t.Tfmtln(`</tr>`)

	for _, timpres := range []map[calculo_imposto.TipoImposto]*calculo_imposto.CalculoResultadoImposto{
		res_pj.Total().Impostos, res_pf.Total().Impostos} {
		// usa lista fixa para manter ordem
		for _, timpid := range calculo_imposto.TipoImpostoLista {
			if timp, tok := timpres[timpid]; tok {
				t.Tfmtln(`<tr>`)
				t.Tfmtln(`<th colspan="2">%s</th>`, timpid.String())
				t.Tfmtln(`<td align="right">%.2f%%</td>`, timp.Aliquota)
				t.Tfmtln(`<td align="right">%.2f</td>`, timp.ValorImposto/12.0)
				t.Tfmtln(`</tr>`)
			}
		}
	}

	t.Tfmtln(`</table>`)
}

func calculo_SimplesAnexo5(t *Template, faturamentoMensal float64, exportacao bool) error {
	prolaboreMensal := 1_100.00

	// PJ
	fat12Mes := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(faturamentoMensal),
		calculo_pj.WithFS_FolhadePagamentoMensal(prolaboreMensal))

	var impostoAplicado calculo_imposto.ImpostoAplicado
	if exportacao {
		impostoAplicado = calculo_imposto.ImpostoAplicado_Exportacao{}
	} else {
		impostoAplicado = calculo_imposto.ImpostoAplicado_Normal{}
	}

	calc := calculo_pj.NewCalculoSimples(fat12Mes, calculo_pj.Simples2018AnexoV,
		calculo_pj.WithCS_ImpostoAplicado(impostoAplicado))

	fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(faturamentoMensal),
		calculo_pj.WithFS_FolhadePagamentoMensal(prolaboreMensal))

	res_pj, err := calc.Calculo(fat)
	if err != nil {
		return err
	}

	// PF
	calc_pf := calculo_pf.NewCalculoPF(calculo_pf.IRPF2021_Mensal,
		calculo_pf.WithCPF_TabelaINSS(calculo_pf.NewTabelaINSS_PJ_2021()))

	prolabore := calculo_pf.NewProlabore_Faturamento_Adapter(fat)

	res_pf, err := calc_pf.Calculo(prolabore)
	if err != nil {
		return err
	}

	outputCalc(t, res_pj, res_pf)

	return nil
}

func calculo_SimplesAnexo5FatorR(t *Template, faturamentoMensal float64, exportacao bool) error {
	prolaboreMensal := faturamentoMensal * 0.28

	// PJ
	fat12Mes := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(faturamentoMensal),
		calculo_pj.WithFS_FolhadePagamentoMensal(prolaboreMensal))

	var impostoAplicado calculo_imposto.ImpostoAplicado
	if exportacao {
		impostoAplicado = calculo_imposto.ImpostoAplicado_Exportacao{}
	} else {
		impostoAplicado = calculo_imposto.ImpostoAplicado_Normal{}
	}

	calc := calculo_pj.NewCalculoSimples(fat12Mes, calculo_pj.Simples2018AnexoV,
		calculo_pj.WithCS_AnexoFatorR(calculo_pj.Simples2018AnexoIII),
		calculo_pj.WithCS_ImpostoAplicado(impostoAplicado))

	fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(faturamentoMensal),
		calculo_pj.WithFS_FolhadePagamentoMensal(prolaboreMensal))

	res_pj, err := calc.Calculo(fat)
	if err != nil {
		return err
	}

	// PF
	calc_pf := calculo_pf.NewCalculoPF(calculo_pf.IRPF2021_Mensal,
		calculo_pf.WithCPF_TabelaINSS(calculo_pf.NewTabelaINSS_PJ_2021()))

	prolabore := calculo_pf.NewProlabore_Faturamento_Adapter(fat)

	res_pf, err := calc_pf.Calculo(prolabore)
	if err != nil {
		return err
	}

	outputCalc(t, res_pj, res_pf)

	return nil
}

func calculo_LucroPresumido(t *Template, faturamentoMensal float64, exportacao bool) error {
	prolaboreMensal := 1_100.00

	// PJ
	var impostoAplicado calculo_imposto.ImpostoAplicado
	if exportacao {
		impostoAplicado = calculo_imposto.ImpostoAplicado_Exportacao{}
	} else {
		impostoAplicado = calculo_imposto.ImpostoAplicado_Normal{}
	}

	calc := calculo_pj.NewCalculoLucroPresumido(
		calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_IRPJ_PRESTACAO_DE_SERVICOS),
		calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_CSLL_PRESTACAO_DE_SERVICOS),
		calculo_pj.WithCLP_ISS(calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_ALIQUOTA_ISS_SOFTWARE)),
		calculo_pj.WithCLP_ImpostoAplicado(impostoAplicado))

	fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(faturamentoMensal),
		calculo_pj.WithFS_FolhadePagamentoMensal(prolaboreMensal))

	res_pj, err := calc.Calculo(fat)
	if err != nil {
		return err
	}

	// PF
	calc_pf := calculo_pf.NewCalculoPF(calculo_pf.IRPF2021_Mensal,
		calculo_pf.WithCPF_TabelaINSS(calculo_pf.NewTabelaINSS_PJ_2021()))

	prolabore := calculo_pf.NewProlabore_Faturamento_Adapter(fat)

	res_pf, err := calc_pf.Calculo(prolabore)
	if err != nil {
		return err
	}

	outputCalc(t, res_pj, res_pf)

	return nil
}
