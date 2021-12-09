package calculo_imposto

//
// Consts_Impl
//

type Consts_Impl struct {
	values map[ConstItem]float64
}

func NewConsts_Impl(values map[ConstItem]float64) *Consts_Impl {
	return &Consts_Impl{values}
}

func (c *Consts_Impl) Exists(item ConstItem) bool {
	_, ex := c.values[item]
	return ex
}

func (c *Consts_Impl) Value(item ConstItem) float64 {
	return c.values[item]
}

//
// Consts_Override
//

type Consts_Override struct {
	base   Consts
	values map[ConstItem]float64
}

func NewConsts_Override(base Consts, values map[ConstItem]float64) *Consts_Override {
	return &Consts_Override{base, values}
}

func (c *Consts_Override) Exists(item ConstItem) bool {
	_, ex := c.values[item]
	if ex {
		return true
	}
	return c.base.Exists(item)
}

func (c *Consts_Override) Value(item ConstItem) float64 {
	v, ok := c.values[item]
	if ok {
		return v
	}
	return c.base.Value(item)
}

//
// Consts_2021
//

func NewConsts_2021() *Consts_Impl {
	return NewConsts_Impl(map[ConstItem]float64{
		ConstItem_ALIQUOTA_CPP: 20.0,

		ConstItem_ALIQUOTA_PIS_CUMULATIVO:     0.65,
		ConstItem_ALIQUOTA_PIS_NAO_CUMULATIVO: 1.65,

		ConstItem_ALIQUOTA_COFINS_CUMULATIVO:     3.0,
		ConstItem_ALIQUOTA_COFINS_NAO_CUMULATIVO: 7.6,

		ConstItem_ALIQUOTA_CSLL: 9.0,

		ConstItem_ALIQUOTA_IRPJ:            15.0,
		ConstItem_ALIQUOTA_IRPJ_ADICIONAL:  10.0,
		ConstItem_IRPJ_ADICIONAL_VALOR_MES: 20_000.0,

		ConstItem_ALIQUOTA_INSS_INDIVIDUAL:              20.0,
		ConstItem_ALIQUOTA_INSS_INDIVIDUAL_SIMPLIFICADO: 11.0,
		ConstItem_ALIQUOTA_INSS_PJ:                      11.0,

		// Fonte: http://normas.receita.fazenda.gov.br/sijut2consulta/link.action?idAto=92278
		ConstItem_PERCENTUAL_FATOR_R: 28.0,

		// Fonte: https://portaldacontabilidade.clmcontroller.com.br/impostos-sobre-software/
		ConstItem_ALIQUOTA_ISS_SOFTWARE:  2.9,
		ConstItem_ALIQUOTA_ICMS_SOFTWARE: 5.0,

		// Fonte: https://portaldacontabilidade.clmcontroller.com.br/como-calcular-o-lucro-presumido/
		ConstItem_LUCRO_PRESUMIDO_IRPJ_REVENDA_COMBUSTIVEIS:   1.6,
		ConstItem_LUCRO_PRESUMIDO_IRPJ_REGRA_GERAL:            8.0,
		ConstItem_LUCRO_PRESUMIDO_IRPJ_SERVICOS_DE_TRANSPORTE: 16.0,
		ConstItem_LUCRO_PRESUMIDO_IRPJ_PRESTACAO_DE_SERVICOS:  32.0,

		ConstItem_LUCRO_PRESUMIDO_CSLL_REGRA_GERAL:           12.0,
		ConstItem_LUCRO_PRESUMIDO_CSLL_PRESTACAO_DE_SERVICOS: 32.0,
	})
}

var Consts_Atual = NewConsts_2021()

//
// ImpostoAplicado_Normal
//

type ImpostoAplicado_Normal struct {
}

func (i ImpostoAplicado_Normal) ImpostoAplicavel(imposto TipoImposto) bool {
	return true
}

//
// ImpostoAplicado_Manual
//

type ImpostoAplicado_Manual struct {
	Impostos []TipoImposto
}

func (i ImpostoAplicado_Manual) ImpostoAplicavel(imposto TipoImposto) bool {
	for _, i := range i.Impostos {
		if i == imposto {
			return true
		}
	}
	return false
}

//
// ImpostoAplicado_Exportacao
// Exportação é isenta de PIS/Cofins e possivelmente ISS (dependendo do estado)
//

type ImpostoAplicado_Exportacao struct {
	IncluirISS bool
}

func (i ImpostoAplicado_Exportacao) ImpostoAplicavel(imposto TipoImposto) bool {
	switch imposto {
	case TipoImposto_PIS, TipoImposto_Cofins:
		return false
	case TipoImposto_ISS:
		return i.IncluirISS
	}
	return true
}
