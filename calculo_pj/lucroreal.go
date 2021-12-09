package calculo_pj

import (
	"fmt"
	"strings"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
)

// https://portaldacontabilidade.clmcontroller.com.br/como-calcular-o-lucro-real/

const CalculoLucroReal_SOURCE = "pj-lucro-real"

type CalculoLucroReal struct {
	consts          calculo_imposto.Consts
	iss             float64
	icms            float64
	impostoAplicado calculo_imposto.ImpostoAplicado
	explain         calculo_imposto.ExplainIntf
}

func NewCalculoLucroReal(opt ...CalculoLucroRealOpt) *CalculoLucroReal {
	ret := &CalculoLucroReal{
		consts:          calculo_imposto.Consts_Atual,
		impostoAplicado: &calculo_imposto.ImpostoAplicado_Normal{},
		explain:         &calculo_imposto.ExplainEmpty{},
	}
	for _, o := range opt {
		o(ret)
	}
	return ret
}

func (c *CalculoLucroReal) Calculo(faturamento Faturamento) (*calculo_imposto.CalculoResultado, error) {
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

func (c *CalculoLucroReal) calculoMes(faturamento Faturamento, mes int) (*calculo_imposto.CalculoResultadoItem, error) {
	if faturamento.LucroMes(mes) == nil {
		return nil, fmt.Errorf("LucroMes não informado para o mes %d", mes)
	}

	c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		"Processando mês %<mes>d", map[string]interface{}{
			"mes": mes,
		})

	lucroMes := *faturamento.LucroMes(mes)

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

		c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Usando valor mês %<valorMes>.2f, lucro mês %<lucroMes>.2f, folha de pagamento mês %<folhadepagamentoMes>.2f", map[string]interface{}{
				"valorMes":            faturamento.ValorMes(mes),
				"lucroMes":            lucroMes,
				"folhadepagamentoMes": *faturamento.FolhaDePagamentoMes(mes),
			})
	} else {
		c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Usando valor mês %<valorMes>.2f, lucro mês %<lucroMes>.2f, sem folha de pagamento", map[string]interface{}{
				"valorMes": faturamento.ValorMes(mes),
				"lucroMes": lucroMes,
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
		irpjValorLucroReal := lucroMes

		impostoIRPJ := &calculo_imposto.CalculoResultadoImposto{
			Aliquota:     c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ),
			ValorImposto: irpjValorLucroReal * (c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ) / 100.0),
		}
		ret.Impostos[calculo_imposto.TipoImposto_IRPJ] = impostoIRPJ
		ret.ValorImposto += impostoIRPJ.ValorImposto

		c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Calculo IRPJ com lucro de %<lucroReal>.2f (%<lucroRealPct>.2f%%), "+
				"aliquota %<aliquota>.2f%%", map[string]interface{}{
				"lucroReal":    irpjValorLucroReal,
				"lucroRealPct": irpjValorLucroReal / faturamento.ValorMes(mes) * 100.0,
				"aliquota":     c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ),
			})

		// IRPJ Adicional
		if irpjValorLucroReal > c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES) {
			impostoIRPJAdicional := &calculo_imposto.CalculoResultadoImposto{
				Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ_ADICIONAL),
				ValorImposto: (irpjValorLucroReal - c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES)) *
					(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ_ADICIONAL) / 100.0),
			}
			ret.Impostos[calculo_imposto.TipoImposto_IRPJ_Adicional] = impostoIRPJAdicional
			ret.ValorImposto += impostoIRPJAdicional.ValorImposto

			c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
				"Calculo IRPJ adicional por lucro mês maior que %<irpjAdicionalValorMes>.2f, sobre "+
					"adicional de %<valorBaseIRPJAdicional>.2f, alíquota %<aliquota>.2f%%", map[string]interface{}{
					"irpjAdicionalValorMes":  c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES),
					"valorBaseIRPJAdicional": irpjValorLucroReal - c.consts.Value(calculo_imposto.ConstItem_IRPJ_ADICIONAL_VALOR_MES),
					"aliquota":               c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_IRPJ_ADICIONAL),
				})
		}
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_IRPJ.String())
	}

	// CSLL
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CSLL) {
		impostoCSLL := &calculo_imposto.CalculoResultadoImposto{
			Aliquota:     c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_CSLL),
			ValorImposto: lucroMes * (c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_CSLL) / 100.0),
		}
		ret.Impostos[calculo_imposto.TipoImposto_CSLL] = impostoCSLL
		ret.ValorImposto += impostoCSLL.ValorImposto
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_CSLL.String())
	}

	// PIS
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_PIS) {
		impostoPIS := &calculo_imposto.CalculoResultadoImposto{
			Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_PIS_NAO_CUMULATIVO),
			ValorImposto: faturamento.ValorMes(mes) *
				(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_PIS_NAO_CUMULATIVO) / 100.0),
		}
		// TODO: deduzir custos
		ret.Impostos[calculo_imposto.TipoImposto_PIS] = impostoPIS
		ret.ValorImposto += impostoPIS.ValorImposto
	} else {
		impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_PIS.String())
	}

	// Cofins
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_Cofins) {
		impostoCofins := &calculo_imposto.CalculoResultadoImposto{
			Aliquota: c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_COFINS_NAO_CUMULATIVO),
			ValorImposto: faturamento.ValorMes(mes) *
				(c.consts.Value(calculo_imposto.ConstItem_ALIQUOTA_COFINS_NAO_CUMULATIVO) / 100.0),
		}
		// TODO: deduzir custos
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
		c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
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
	c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		impostosMsg, impostosParams)

	c.explain.Add(CalculoLucroReal_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		"Cálculo com alíquota %<aliquota>.2f%%, valor do imposto %<valorImposto>.2f", map[string]interface{}{
			"aliquota":     ret.AliquotaImposto(),
			"valorImposto": ret.ValorImposto,
		})

	return ret, nil
}

type CalculoLucroRealOpt func(*CalculoLucroReal)

func WithCLR_Consts(consts calculo_imposto.Consts) CalculoLucroRealOpt {
	return func(c *CalculoLucroReal) {
		c.consts = consts
	}
}

func WithCLR_ISS(iss float64) CalculoLucroRealOpt {
	return func(c *CalculoLucroReal) {
		c.iss = iss
	}
}

func WithCLR_ICMS(icms float64) CalculoLucroRealOpt {
	return func(c *CalculoLucroReal) {
		c.icms = icms
	}
}

func WithCLR_ImpostoAplicado(impostoAplicado calculo_imposto.ImpostoAplicado) CalculoLucroRealOpt {
	return func(c *CalculoLucroReal) {
		c.impostoAplicado = impostoAplicado
	}
}

func WithCLR_Explain(explain calculo_imposto.ExplainIntf) CalculoLucroRealOpt {
	return func(c *CalculoLucroReal) {
		c.explain = explain
	}
}
