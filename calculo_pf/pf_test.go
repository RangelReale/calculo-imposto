package calculo_pf

import (
	"testing"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
	"github.com/stretchr/testify/assert"
)

func TestIRPF2021(t *testing.T) {
	// http://www.receita.fazenda.gov.br/Aplicacoes/ATRJO/Simulador/simulador.asp
	testCases := []struct {
		prolaboreMensal float64
		expectedIRPF    float64
	}{
		{1_100.00, 0.00},
		{2_400.00, 37.20},
		{3_200.00, 125.20},
		{5_400.00, 615.64},
		{10_000.00, 1_880.64},
		{15_000.00, 3_255.64},
	}

	impostoAplicado := &calculo_imposto.ImpostoAplicado_Manual{
		Impostos: []calculo_imposto.TipoImposto{calculo_imposto.TipoImposto_IRPF},
	}

	for _, tc := range testCases {
		calc := NewCalculoPF(IRPF2021_Mensal, WithCPF_TabelaINSS(NewTabelaINSS_PJ_2021()),
			WithCPF_ImpostoAplicado(impostoAplicado))
		prolabore := NewProlabore_Static(1, WithPS_ValorMensal(tc.prolaboreMensal))

		cres, err := calc.Calculo(prolabore)
		if err != nil {
			t.Fatal(err)
		}
		crestotal := cres.Total()

		assert.InDeltaf(t, tc.expectedIRPF, crestotal.ValorImposto, 0.01, "Valor de imposto incorreto, esperado %.2f, calculado %.2f",
			tc.expectedIRPF, crestotal.ValorImposto)
	}
}
