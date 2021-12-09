package calculo_pj

import (
	"strings"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
)

// https://portaldacontabilidade.clmcontroller.com.br/como-calcular-o-lucro-presumido/

const CalculoLucroPresumido_SOURCE = "pj-lucro-presumido"

type CalculoLucroPresumido struct {
	consts             calculo_imposto.Consts
	irpjLucroPresumido float64
	csllLucroPresumido float64
	iss                float64
	icms               float64
	impostoAplicado    calculo_imposto.ImpostoAplicado
	explain            calculo_imposto.ExplainIntf
}

func NewCalculoLucroPresumido(irpjLucroPresumido float64, csllLucroPresumido float64,
	opt ...CalculoLucroPresumidoOpt) *CalculoLucroPresumido {
	ret := &CalculoLucroPresumido{
		consts:             calculo_imposto.Consts_Atual,
		irpjLucroPresumido: irpjLucroPresumido,
		csllLucroPresumido: csllLucroPresumido,
		impostoAplicado:    &calculo_imposto.ImpostoAplicado_Normal{},
		explain:            &calculo_imposto.ExplainEmpty{},
	}
	for _, o := range opt {
		o(ret)
	}
	return ret
}

// TODO: fazer apuração trimestral
func (c *CalculoLucroPresumido) Calculo(faturamento Faturamento) (*calculo_imposto.CalculoResultado, error) {
	ret := &calculo_imposto.CalculoResultado{}

	for mes := 0; mes < faturamento.Meses(); mes++ {
		rmes, err := c.calculoMes(faturamento, mes)
		if err != nil {
			return nil, err
		}
		ret.Items = append(ret.Items, rmes)
	}

	return ret, nil
}

func (c *CalculoLucroPresumido) calculoMes(faturamento Faturamento, mes int) (*calculo_imposto.CalculoResultadoItem, error) {
	c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		"Processando mês %<mes>d", map[string]interface{}{
			"mes": mes,
		})

	ret := &calculo_imposto.CalculoResultadoItem{
		Periodo:       mes,
		TipoPeriodo:   calculo_imposto.TipoPeriodo_Mes,
		ValorOriginal: faturamento.ValorMes(mes),
		Impostos:      map[calculo_imposto.TipoImposto]*calculo_imposto.CalculoResultadoImposto{},
	}

	if faturamento.FolhaDePagamentoMes(mes) != nil {
		ret.Extra = map[calculo_imposto.CalculoResultadoExtra]interface{}{
			calculo_imposto.CalculoResultadoExtra_ValorFolhaDePagamento: *faturamento.FolhaDePagamentoMes(mes),
		}

		c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Usando valor mês %<valorMes>.2f, folha de pagamento mês %<folhadepagamentoMes>.2f", map[string]interface{}{
				"valorMes":            faturamento.ValorMes(mes),
				"folhadepagamentoMes": *faturamento.FolhaDePagamentoMes(mes),
			})
	} else {
		c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Usando valor mês %<valorMes>.2f, sem folha de pagamento", map[string]interface{}{
				"valorMes": faturamento.ValorMes(mes),
			})
	}

	var impostosExcluidos []string

	// CPP (Contribuição Previdenciária Patronal)
	if faturamento.FolhaDePagamentoMes(mes) != nil {
		if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CPP) {
			impostoCPP := &calculo_imposto.CalculoResultadoImposto{
				Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_CPP),
				ValorImposto: *faturamento.FolhaDePagamentoMes(mes) *
					(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_CPP) / 100.0),
			}
			ret.Impostos[calculo_imposto.TipoImposto_CPP] = impostoCPP
			ret.ValorImposto += impostoCPP.ValorImposto
		} else {
			impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_CPP.String())
		}
	}

	// IRPJ
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_IRPJ) {
		irpjValorLucroPresumido := faturamento.ValorMes(mes) * (c.irpjLucroPresumido / 100.0)

		impostoIRPJ := &calculo_imposto.CalculoResultadoImposto{
			Aliquota:     c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ),
			ValorImposto: irpjValorLucroPresumido * (c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ) / 100.0),
		}
		ret.Impostos[calculo_imposto.TipoImposto_IRPJ] = impostoIRPJ
		ret.ValorImposto += impostoIRPJ.ValorImposto

		c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Calculo IRPJ com lucro presumido de %<lucroPresumido>.2f (%<lucroPresumidoPct>.2f%%), "+
				"aliquota %<aliquota>.2f%%", map[string]interface{}{
				"lucroPresumido":    irpjValorLucroPresumido,
				"lucroPresumidoPct": c.irpjLucroPresumido,
				"aliquota":          c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ),
			})

		// IRPJ Adicional
		if irpjValorLucroPresumido > c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES) {
			impostoIRPJAdicional := &calculo_imposto.CalculoResultadoImposto{
				Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ_ADICIONAL),
				ValorImposto: (irpjValorLucroPresumido - c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES)) *
					(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ_ADICIONAL) / 100.0),
			}
			ret.Impostos[calculo_imposto.TipoImposto_IRPJ_Adicional] = impostoIRPJAdicional
			ret.ValorImposto += impostoIRPJAdicional.ValorImposto

			c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
				"Calculo IRPJ adicional por lucro presumido mês maior que %<irpjAdicionalValorMes>.2f, sobre "+
					"adicional de %<valorBaseIRPJAdicional>.2f, alíquota %<aliquota>.2f%%", map[string]interface{}{
					"irpjAdicionalValorMes":  c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES),
					"valorBaseIRPJAdicional": irpjValorLucroPresumido - c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES),
					"aliquota":               c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ_ADICIONAL),
				})
		}
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_IRPJ.String())
	}

	// CSLL
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CSLL) {
		impostoCSLL := &calculo_imposto.CalculoResultadoImposto{
			Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_CSLL),
			ValorImposto: faturamento.ValorMes(mes) * (c.csllLucroPresumido / 100.0) *
				(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_CSLL) / 100.0),
		}
		ret.Impostos[calculo_imposto.TipoImposto_CSLL] = impostoCSLL
		ret.ValorImposto += impostoCSLL.ValorImposto
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_CSLL.String())
	}

	// PIS
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_PIS) {
		impostoPIS := &calculo_imposto.CalculoResultadoImposto{
			Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_PIS_CUMULATIVO),
			ValorImposto: faturamento.ValorMes(mes) *
				(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_PIS_CUMULATIVO) / 100.0),
		}
		ret.Impostos[calculo_imposto.TipoImposto_PIS] = impostoPIS
		ret.ValorImposto += impostoPIS.ValorImposto
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_PIS.String())
	}

	// Cofins
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_Cofins) {
		impostoCofins := &calculo_imposto.CalculoResultadoImposto{
			Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_COFINS_CUMULATIVO),
			ValorImposto: faturamento.ValorMes(mes) *
				(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_COFINS_CUMULATIVO) / 100.0),
		}
		ret.Impostos[calculo_imposto.TipoImposto_Cofins] = impostoCofins
		ret.ValorImposto += impostoCofins.ValorImposto
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_Cofins.String())
	}

	// ISS
	if c.iss > 0 {
		if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_ISS) {
			impostoISS := &calculo_imposto.CalculoResultadoImposto{
				Aliquota:     c.iss,
				ValorImposto: faturamento.ValorMes(mes) * (c.iss / 100.0),
			}
			ret.Impostos[calculo_imposto.TipoImposto_ISS] = impostoISS
			ret.ValorImposto += impostoISS.ValorImposto
		} else {
			impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_ISS.String())
		}
	}

	// ICMS
	if c.icms > 0 {
		if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_ICMS) {
			impostoICMS := &calculo_imposto.CalculoResultadoImposto{
				Aliquota:     c.icms,
				ValorImposto: faturamento.ValorMes(mes) * (c.icms / 100.0),
			}
			ret.Impostos[calculo_imposto.TipoImposto_ICMS] = impostoICMS
			ret.ValorImposto += impostoICMS.ValorImposto
		} else {
			impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_ICMS.String())
		}
	}

	if len(impostosExcluidos) > 0 {
		c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Excluindo impostos '%<impostosExcluidos>s'", map[string]interface{}{
				"impostosExcluidos": strings.Join(impostosExcluidos, ", "),
			})
	}

	impostosMsg := "Impostos: "
	impostosParams := map[string]interface{}{}
	first := true
	// Usa lista para manter ordem
	for _, impidx := range calculo_imposto.TipoImpostoLista {
		if imp, impok := ret.Impostos[impidx]; impok {
			if !first {
				impostosMsg += ", "
			}
			impostoId := "imposto_" + strings.Replace(string(impidx), "-", "_", -1)

			impostosMsg += "%<" + impostoId + ">s: %<" + impostoId + "_aliquota>.2f%%"
			impostosParams[impostoId] = impidx.String()
			impostosParams[impostoId+"_aliquota"] = imp.Aliquota
			first = false
		}
	}
	c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		impostosMsg, impostosParams)

	c.explain.Add(CalculoLucroPresumido_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		"Cálculo com alíquota %<aliquota>.2f%%, valor do imposto %<valorImposto>.2f", map[string]interface{}{
			"aliquota":     ret.AliquotaImposto(),
			"valorImposto": ret.ValorImposto,
		})

	return ret, nil
}

type CalculoLucroPresumidoOpt func(*CalculoLucroPresumido)

func WithCLP_Consts(consts calculo_imposto.Consts) CalculoLucroPresumidoOpt {
	return func(c *CalculoLucroPresumido) {
		c.consts = consts
	}
}

func WithCLP_ISS(iss float64) CalculoLucroPresumidoOpt {
	return func(c *CalculoLucroPresumido) {
		c.iss = iss
	}
}

func WithCLP_ICMS(icms float64) CalculoLucroPresumidoOpt {
	return func(c *CalculoLucroPresumido) {
		c.icms = icms
	}
}

func WithCLP_ImpostoAplicado(impostoAplicado calculo_imposto.ImpostoAplicado) CalculoLucroPresumidoOpt {
	return func(c *CalculoLucroPresumido) {
		c.impostoAplicado = impostoAplicado
	}
}

func WithCLP_Explain(explain calculo_imposto.ExplainIntf) CalculoLucroPresumidoOpt {
	return func(c *CalculoLucroPresumido) {
		c.explain = explain
	}
}
