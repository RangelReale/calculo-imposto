package calculo_pj

import (
	"testing"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
	"github.com/stretchr/testify/assert"
)

func TestLucroPresumido2021(t *testing.T) {
	testCases := []struct {
		inputFaturamentoMensal     float64
		inputFolhaMensal           float64
		expectedDescontoPercentual float64
		exportacao                 bool
	}{
		{10_000.00, 10_000.00 * 0.28, 19.830, false},
		{10_000.00, 1_100.00, 16.430, false},
		{20_000.00, 20_000.00 * 0.28, 19.830, false},
		{20_000.00, 1_100.00, 15.330, false},
		{10_000.00, 10_000.00 * 0.28, 13.280, true},
		{10_000.00, 1_100.00, 9.880, true},
		{20_000.00, 20_000.00 * 0.28, 13.280, true},
		{20_000.00, 1_100.00, 8.780, true},
	}

	consts := calculo_imposto.NewConsts_2021()

	for _, tc := range testCases {
		fat := NewFaturamento_Static(1, WithFS_ValorMensal(tc.inputFaturamentoMensal),
			WithFS_FolhadePagamentoMensal(tc.inputFolhaMensal))

		calcopt := []CalculoLucroPresumidoOpt{
			WithCLP_ISS(calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_ALIQUOTA_ISS_SOFTWARE)),
			WithCLP_Consts(consts),
		}
		if tc.exportacao {
			calcopt = append(calcopt, WithCLP_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}))
		}

		calc := NewCalculoLucroPresumido(consts.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_IRPJ_PRESTACAO_DE_SERVICOS),
			consts.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_CSLL_PRESTACAO_DE_SERVICOS), calcopt...)

		cres, err := calc.Calculo(fat)
		if err != nil {
			t.Fatal(err)
		}

		crestotal := cres.Total()

		assert.InDeltaf(t, tc.expectedDescontoPercentual, crestotal.AliquotaImposto(),
			0.01, "Percentual de imposto lucro presumido incorreto, esperado %.4f, calculado %.4f",
			tc.expectedDescontoPercentual, crestotal.AliquotaImposto())
	}
}
