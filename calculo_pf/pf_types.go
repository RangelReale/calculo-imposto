package calculo_pf

type IRPFFaixa struct {
	LimiteMaximo float64 `json:"limite-maximo"`
	Aliquota     float64 `json:"aliquota"`
	ValorDeducao float64 `json:"valor-deducao"`
}
