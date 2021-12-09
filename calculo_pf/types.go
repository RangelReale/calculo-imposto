package calculo_pf

import calculo_imposto "github.com/RangelReale/calculo-imposto"

const (
	CalculoResultadoExtra_TaxaINSS calculo_imposto.CalculoResultadoExtra = "taxa-inss"
)

type IRPF_Tipo int

const (
	IRPF_Tipo_Mensal IRPF_Tipo = 0
	IRPF_Tipo_Anual            = 1
)

type IRPF struct {
	Ano    int
	Tipo   IRPF_Tipo
	Faixas []IRPFFaixa
}

type Prolabore interface {
	Meses() int
	ValorMes(mes int) float64
}

type TabelaINSS interface {
	Taxa(valorMensal float64) (taxa float64, valorMax float64)
}

type TabelaINSS_Faixa struct {
	LimiteMaximo float64 `json:"limite-maximo"`
	Aliquota     float64 `json:"aliquota"`
}

type Calculo interface {
	Calculo(prolabore Prolabore) (*calculo_imposto.CalculoResultado, error)
}
