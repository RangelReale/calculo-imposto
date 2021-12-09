package calculo_imposto

type TipoImposto string

const (
	TipoImposto_IRPJ           TipoImposto = "irpj"
	TipoImposto_IRPJ_Adicional TipoImposto = "irpj-adicional"
	TipoImposto_CPP            TipoImposto = "cpp"
	TipoImposto_CSLL           TipoImposto = "csll"
	TipoImposto_PIS            TipoImposto = "PIS"
	TipoImposto_Cofins         TipoImposto = "cofins"
	TipoImposto_ISS            TipoImposto = "iss"
	TipoImposto_ICMS           TipoImposto = "icms"
	TipoImposto_IRPF           TipoImposto = "irpf"
	TipoImposto_INSS           TipoImposto = "inss"
	TipoImposto_SIMPLES        TipoImposto = "simples"
)

var TipoImpostoLista = []TipoImposto{TipoImposto_SIMPLES, TipoImposto_IRPJ, TipoImposto_IRPJ_Adicional, TipoImposto_CPP,
	TipoImposto_CSLL, TipoImposto_PIS, TipoImposto_Cofins, TipoImposto_ISS, TipoImposto_ICMS, TipoImposto_IRPF,
	TipoImposto_INSS}

func (ti TipoImposto) String() string {
	switch ti {
	case TipoImposto_IRPJ:
		return "IRPJ"
	case TipoImposto_IRPJ_Adicional:
		return "IRPJ Adicional"
	case TipoImposto_CPP:
		return "CPP"
	case TipoImposto_CSLL:
		return "CSLL"
	case TipoImposto_PIS:
		return "PIS"
	case TipoImposto_Cofins:
		return "Cofins"
	case TipoImposto_ISS:
		return "ISS"
	case TipoImposto_ICMS:
		return "ICMS"
	case TipoImposto_IRPF:
		return "IRPF"
	case TipoImposto_INSS:
		return "INSS"
	case TipoImposto_SIMPLES:
		return "Simples"
	}
	return "Unknown"
}

// ImpostoAplicado determina quais impostos são aplicáveis
type ImpostoAplicado interface {
	ImpostoAplicavel(imposto TipoImposto) bool
}

type TipoPeriodo int

const (
	TipoPeriodo_Mes       TipoPeriodo = 0
	TipoPeriodo_Trimestre TipoPeriodo = 1
)

type CalculoResultadoExtra string

const (
	CalculoResultadoExtra_ValorFolhaDePagamento CalculoResultadoExtra = "valor-folha-de-pagamento"
	CalculoResultadoExtra_Faixa                                       = "faixa"
)

type CalculoResultado struct {
	Items []*CalculoResultadoItem `json:"items"`
}

func (cr *CalculoResultado) Total() *CalculoResultadoItem {
	ret := &CalculoResultadoItem{
		Impostos: map[TipoImposto]*CalculoResultadoImposto{},
	}
	type impostoAliqotaData struct {
		Soma       float64
		Quantidade int
	}
	impostoAliquota := map[TipoImposto]*impostoAliqotaData{}

	for _, item := range cr.Items {
		ret.ValorOriginal += item.ValorOriginal
		ret.ValorImposto += item.ValorImposto
		for impid, imp := range item.Impostos {
			if _, fok := ret.Impostos[impid]; !fok {
				ret.Impostos[impid] = &CalculoResultadoImposto{}
			}
			ret.Impostos[impid].ValorImposto += imp.ValorImposto

			if _, fok := impostoAliquota[impid]; !fok {
				impostoAliquota[impid] = &impostoAliqotaData{}
			}
			impostoAliquota[impid].Soma += imp.Aliquota
			impostoAliquota[impid].Quantidade += 1
		}
	}

	for impid, ia := range impostoAliquota {
		ret.Impostos[impid].Aliquota = ia.Soma / float64(ia.Quantidade)
	}

	return ret
}

type CalculoResultadoImposto struct {
	Aliquota     float64 `json:"aliquota"`
	ValorImposto float64 `json:"valor-imposto"`
}

type CalculoResultadoItem struct {
	Periodo       int                                      `json:"periodo"`
	TipoPeriodo   TipoPeriodo                              `json:"tipo-periodo"`
	ValorOriginal float64                                  `json:"valor-original"`
	ValorImposto  float64                                  `json:"valor-imposto"`
	Impostos      map[TipoImposto]*CalculoResultadoImposto `json:"impostos"`
	Extra         map[CalculoResultadoExtra]interface{}    `json:"extra"`
}

func (c *CalculoResultadoItem) AliquotaImposto() float64 {
	if c.ValorOriginal == 0 {
		return 0
	}
	return c.ValorImposto / c.ValorOriginal * 100.0
}
