package calculo_pj

import calculo_imposto "github.com/RangelReale/calculo-imposto"

const (
	CalculoResultadoExtra_Simples_FatorR     calculo_imposto.CalculoResultadoExtra = "simples-fator-r"
	CalculoResultadoExtra_Simples_ValorRBT12                                       = "simples-valor-rbt12"
)

type SimplesFaixa struct {
	LimiteMaximo       float64                   `json:"limite-maximo"`
	Aliquota           float64                   `json:"aliquota"`
	ValorDeducao       float64                   `json:"valor-deducao"`
	ReparticaoImpostos SimplesReparticaoImpostos `json:"reparticao-impostos"`
}

// Fonte: https://blog.contabilizei.com.br/contabilidade-online/anexo-3-simples-nacional/
type SimplesReparticaoImpostos struct {
	CPP    float64 `json:"cpp"`
	ISS    float64 `json:"iss"`
	CSLL   float64 `json:"csll"`
	IRPJ   float64 `json:"irpj"`
	Cofins float64 `json:"cofins"`
	PIS    float64 `json:"pis"`
}

type SimplesAnexo struct {
	Titulo string         `json:"titulo"`
	Numero string         `json:"numero"`
	Faixas []SimplesFaixa `json:"faixas"`
}

type Simples struct {
	AnoVigencia int                      `json:"ano-vigencia"`
	MesVigencia int                      `json:"mes-vigencia"`
	Anexos      map[string]*SimplesAnexo `json:"anexos"`
}
