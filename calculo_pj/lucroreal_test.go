package calculo_pj

import (
	"testing"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
	"github.com/stretchr/testify/assert"
)

func TestLucroReal2021(t *testing.T) {
	testCases := []struct {
		inputFaturamentoMensal     float64
		inputFolhaMensal           float64
		inputLucroMensal           float64
		expectedDescontoPercentual float64
		exportacao                 bool
	}{
		{10_000.00, 10_000.00 * 0.28, 5_000.00, 29.750, false},
		{10_000.00, 1_100.00, 5_000.00, 26.350, false},
		{20_000.00, 20_000.00 * 0.28, 10_000.00, 29.750, false},
		{20_000.00, 1_100.00, 10_000.00, 25.250, false},
		{10_000.00, 10_000.00 * 0.28, 5_000.00, 17.600, true},
		{10_000.00, 1_100.00, 5_000.00, 14.200, true},
		{20_000.00, 20_000.00 * 0.28, 10_000.00, 17.600, true},
		{20_000.00, 1_100.00, 10_000.00, 13.100, true},
	}

	consts := calculo_imposto.NewConsts_2021()

	for _, tc := range testCases {
		fat := NewFaturamento_Static(1, WithFS_ValorMensal(tc.inputFaturamentoMensal),
			WithFS_FolhadePagamentoMensal(tc.inputFolhaMensal), WithFS_LucroMensal(tc.inputLucroMensal))

		calcopt := []CalculoLucroRealOpt{
			WithCLR_ISS(calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_ALIQUOTA_ISS_SOFTWARE)),
			WithCLR_Consts(consts),
		}
		if tc.exportacao {
			calcopt = append(calcopt, WithCLR_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}))
		}

		calc := NewCalculoLucroReal(calcopt...)

		cres, err := calc.Calculo(fat)
		if err != nil {
			t.Fatal(err)
		}

		crestotal := cres.Total()

		assert.InDeltaf(t, tc.expectedDescontoPercentual, crestotal.AliquotaImposto(),
			0.01, "Percentual de imposto lucro real incorreto, esperado %.4f, calculado %.4f",
			tc.expectedDescontoPercentual, crestotal.AliquotaImposto())
	}
}
