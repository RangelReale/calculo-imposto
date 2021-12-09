package calculo_pj

import calculo_imposto "github.com/RangelReale/calculo-imposto"

// Faturamento abstrai uma lista de meses de faturamento
type Faturamento interface {
	Meses() int
	ValorMes(mes int) float64
	LucroMes(mes int) *float64
	FolhaDePagamentoMes(mes int) *float64
}

type Calculo interface {
	Calculo(faturamento Faturamento) (*calculo_imposto.CalculoResultado, error)
}
