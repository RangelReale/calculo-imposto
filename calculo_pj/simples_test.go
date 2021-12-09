package calculo_pj

import (
	"testing"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
	"github.com/stretchr/testify/assert"
)

func TestSimples2021(t *testing.T) {
	// https://contjet.com.br/calculadora-simples-nacional/
	testCases := []struct {
		inputFaturamentoMensal     float64
		expectedDescontoPercentual float64
		exportacao                 bool
	}{
		{10_000.00, 6.000, false},
		{20_000.00, 7.300, false},
		{23_333.33, 7.857, false},
		{33_000.00, 9.045, false},
		{60_000.00, 11.050, false},
		{10_000.00, 3.054, true},
		{20_000.00, 3.715, true},
		{23_333.33, 3.999, true},
		{33_000.00, 4.604, true},
		{60_000.00, 5.624, true},
	}

	for _, tc := range testCases {
		fat12m := NewFaturamento_Static(12, WithFS_ValorMensal(tc.inputFaturamentoMensal))
		fat := NewFaturamento_Static(1, WithFS_ValorMensal(tc.inputFaturamentoMensal))

		calcopt := []CalculoSimplesOpt{WithCS_Consts(calculo_imposto.NewConsts_2021())}
		if tc.exportacao {
			calcopt = append(calcopt, WithCS_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}))
		}
		calc := NewCalculoSimples(fat12m, Simples2018AnexoIII, calcopt...)

		cres, err := calc.Calculo(fat)
		if err != nil {
			t.Fatal(err)
		}

		crestotal := cres.Total()

		assert.InDeltaf(t, tc.expectedDescontoPercentual, crestotal.AliquotaImposto(),
			0.01, "Percentual de imposto simples incorreto, esperado %.4f, calculado %.4f",
			tc.expectedDescontoPercentual, crestotal.AliquotaImposto())
	}
}

func TestSimplesFatorR2021(t *testing.T) {
	testCases := []struct {
		inputFaturamentoMensal     float64
		inputFolhaMensal           float64
		expectedDescontoPercentual float64
		exportacao                 bool
	}{
		{10_000.00, 10_000.00 * 0.28, 6.000, false},
		{10_000.00, 1_100.00, 15.500, false},
		{20_000.00, 20_000.00 * 0.28, 7.300, false},
		{20_000.00, 1_100.00, 16.125, false},
		{10_000.00, 10_000.00 * 0.28, 3.054, true},
		{10_000.00, 1_100.00, 10.671, true},
		{20_000.00, 20_000.00 * 0.28, 3.715, true},
		{20_000.00, 1_100.00, 10.618, true},
	}

	for _, tc := range testCases {
		fat12m := NewFaturamento_Static(12, WithFS_ValorMensal(tc.inputFaturamentoMensal),
			WithFS_FolhadePagamentoMensal(tc.inputFolhaMensal))
		fat := NewFaturamento_Static(1, WithFS_ValorMensal(tc.inputFaturamentoMensal),
			WithFS_FolhadePagamentoMensal(tc.inputFolhaMensal))

		calcopt := []CalculoSimplesOpt{
			WithCS_Consts(calculo_imposto.NewConsts_2021()),
			WithCS_AnexoFatorR(Simples2018AnexoIII),
		}
		if tc.exportacao {
			calcopt = append(calcopt, WithCS_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}))
		}
		calc := NewCalculoSimples(fat12m, Simples2018AnexoV, calcopt...)

		cres, err := calc.Calculo(fat)
		if err != nil {
			t.Fatal(err)
		}

		crestotal := cres.Total()

		assert.InDeltaf(t, tc.expectedDescontoPercentual, crestotal.AliquotaImposto(),
			0.01, "Percentual de imposto simples incorreto, esperado %.4f, calculado %.4f",
			tc.expectedDescontoPercentual, crestotal.AliquotaImposto())
	}
}
