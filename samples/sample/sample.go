package main

import (
	"fmt"
	"strings"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
	"github.com/RangelReale/calculo-imposto/calculo_pf"
	"github.com/RangelReale/calculo-imposto/calculo_pj"
)

func main() {
	var err error

	err = t_simples1()
	//err = t_lucropresumido1()
	//err = t_lucroreal1()
	//err = t_irpf1()

	if err != nil {
		fmt.Printf("Erro no cálculo: %s\n", err.Error())
		return
	}
}

func t_explain(explain calculo_imposto.ExplainIntf) {
	fmt.Printf("%s Explain %s\n", strings.Repeat("====", 10), strings.Repeat("====", 10))
	for _, e := range explain.Items() {
		fmt.Printf("+[P%d] {%s} %s\n", e.Periodo, e.Source, e.FormatMessage())
	}
}

func t_simples1() error {
	explain := calculo_imposto.NewExplain()

	//fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorAnual(540_000),
	//	calculo_pj.WithFS_FolhadePagamentoAnual(540_000*0.28))
	fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(20_000.0),
		calculo_pj.WithFS_FolhadePagamentoMensal(20_000.0*0.28))

	//calc := calculo_pj.NewCalculoSimples(fat, calculo_pj.Simples2018AnexoIII, nil, calculo_pj.SimplesImpostoAplicado_Exportacao{})
	//calc := calculo_pj.NewCalculoSimples(fat, calculo_pj.Simples2018AnexoV, nil, calculo_pj.SimplesImpostoAplicado_Exportacao{})
	calc := calculo_pj.NewCalculoSimples(fat, calculo_pj.Simples2018AnexoV, calculo_pj.WithCS_AnexoFatorR(calculo_pj.Simples2018AnexoIII),
		//calculo_pj.WithCS_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}),
		calculo_pj.WithCS_Explain(explain))

	ret, err := calc.Calculo(fat)
	if err != nil {
		return err
	}

	for _, cm := range ret.Items {
		fmt.Printf("Mes: %d\n", cm.Periodo)
		fmt.Printf("Faturamento: %f\n", cm.ValorOriginal)
		fmt.Printf("Imposto: %f\n", cm.ValorImposto)
		fmt.Printf("Fator R: %f\n", cm.Extra[calculo_pj.CalculoResultadoExtra_Simples_FatorR].(float64))
		for _, ti := range calculo_imposto.TipoImpostoLista {
			if imp, ok := cm.Impostos[ti]; ok {
				if imp.ValorImposto > 0 {
					fmt.Printf("ValorImposto %s: %f (%.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
				}
			}
		}
		fmt.Printf("%s\n", strings.Repeat("-", 10))
	}

	t_explain(explain)

	return nil
}

func t_lucropresumido1() error {
	explain := calculo_imposto.NewExplain()

	calc := calculo_pj.NewCalculoLucroPresumido(
		calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_IRPJ_PRESTACAO_DE_SERVICOS),
		calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_CSLL_PRESTACAO_DE_SERVICOS),
		calculo_pj.WithCLP_ISS(calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_ALIQUOTA_ISS_SOFTWARE)),
		calculo_pj.WithCLP_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}),
		calculo_pj.WithCLP_Explain(explain))

	//fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorAnual(540_000),
	//	calculo_pj.WithFS_FolhadePagamentoAnual(1100.0*12))
	fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorAnual(840_000),
		calculo_pj.WithFS_FolhadePagamentoAnual(1100.0*12))

	ret, err := calc.Calculo(fat)
	if err != nil {
		return err
	}

	for _, cm := range ret.Items {
		fmt.Printf("Mes: %d\n", cm.Periodo)
		fmt.Printf("Faturamento: %f\n", cm.ValorOriginal)
		fmt.Printf("Imposto: %f\n", cm.ValorImposto)
		if fp, ok := cm.Extra[calculo_imposto.CalculoResultadoExtra_ValorFolhaDePagamento]; ok {
			fmt.Printf("Folha de pagamento: %f\n", fp.(float64))
		}
		for _, ti := range calculo_imposto.TipoImpostoLista {
			if imp, ok := cm.Impostos[ti]; ok {
				if imp.ValorImposto > 0 {
					fmt.Printf("ValorImposto %s: %f\n", ti.String(), imp.ValorImposto)
				}
			}
		}
		fmt.Printf("%s\n", strings.Repeat("-", 10))
	}

	t_explain(explain)

	return nil
}

func t_lucroreal1() error {
	explain := calculo_imposto.NewExplain()

	calc := calculo_pj.NewCalculoLucroReal(
		calculo_pj.WithCLR_ISS(calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_ALIQUOTA_ISS_SOFTWARE)),
		calculo_pj.WithCLR_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}),
		calculo_pj.WithCLR_Explain(explain))

	fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorAnual(540_000),
		calculo_pj.WithFS_LucroAnual(400_000))

	ret, err := calc.Calculo(fat)
	if err != nil {
		return err
	}

	for _, cm := range ret.Items {
		fmt.Printf("Mes: %d\n", cm.Periodo)
		fmt.Printf("Faturamento: %f\n", cm.ValorOriginal)
		fmt.Printf("Imposto: %f\n", cm.ValorImposto)
		for _, ti := range calculo_imposto.TipoImpostoLista {
			if imp, ok := cm.Impostos[ti]; ok {
				if imp.ValorImposto > 0 {
					fmt.Printf("ValorImposto %s: %f\n", ti.String(), imp.ValorImposto)
				}
			}
		}
		fmt.Printf("%s\n", strings.Repeat("-", 10))
	}

	t_explain(explain)

	return nil
}

func t_irpf1() error {
	calc := calculo_pf.NewCalculoPF(calculo_pf.IRPF2021_Mensal, calculo_pf.WithCPF_TabelaINSS(calculo_pf.NewTabelaINSS_PJ_2021()))

	prolabore := calculo_pf.NewProlabore_Static(12, calculo_pf.WithPS_ValorMensal(1100.0))

	ret, err := calc.Calculo(prolabore)
	if err != nil {
		return err
	}

	for _, cm := range ret.Items {
		fmt.Printf("Mes: %d\n", cm.Periodo)
		fmt.Printf("Prolabore: %f\n", cm.ValorOriginal)
		fmt.Printf("Imposto: %f\n", cm.ValorImposto)
		for _, ti := range calculo_imposto.TipoImpostoLista {
			if imp, ok := cm.Impostos[ti]; ok {
				if imp.ValorImposto > 0 {
					fmt.Printf("ValorImposto %s: %f\n", ti.String(), imp.ValorImposto)
				}
			}
		}
		fmt.Printf("%s\n", strings.Repeat("-", 10))
	}

	return nil
}
