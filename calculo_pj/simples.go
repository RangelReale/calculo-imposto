package calculo_pj

// https://www.contabilizei.com.br/contabilidade-online/calculo-simples-nacional/

import (
	"fmt"
	"strings"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
)

const CalculoSimples_SOURCE = "pj-simples"

type CalculoSimples struct {
	consts             calculo_imposto.Consts
	faturamento12Meses Faturamento
	anexo              *SimplesAnexo
	anexoFatorR        *SimplesAnexo
	impostoAplicado    calculo_imposto.ImpostoAplicado
	explain            calculo_imposto.ExplainIntf
}

func NewCalculoSimples(faturamento12Meses Faturamento, anexo *SimplesAnexo, opt ...CalculoSimplesOpt) *CalculoSimples {
	ret := &CalculoSimples{
		consts:             calculo_imposto.Consts_Atual,
		faturamento12Meses: faturamento12Meses,
		anexo:              anexo,
		impostoAplicado:    &calculo_imposto.ImpostoAplicado_Normal{},
		explain:            &calculo_imposto.ExplainEmpty{},
	}
	for _, o := range opt {
		o(ret)
	}
	return ret
}

func (c *CalculoSimples) Calculo(faturamento Faturamento) (*calculo_imposto.CalculoResultado, error) {
	if c.faturamento12Meses.Meses() != 12 {
		return nil, fmt.Errorf("O faturamento de 12 meses deve conter exatamente 12 meses")
	}

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

// rbt12 retorna a receita bruta da sua empresa dos últimos 12 meses a partir do mes passado
func (c *CalculoSimples) rbt12(faturamento Faturamento, mes int) float64 {
	var ret float64
	// pega N meses do faturamento passado
	for cm := mes; cm <= 11; cm++ {
		ret += c.faturamento12Meses.ValorMes(cm)
	}
	// pega N meses do faturamento atual
	for cm := 0; cm < mes; cm++ {
		ret += faturamento.ValorMes(cm)
	}
	return ret
}

// folhadepagamento12 retorna a folha de pagamento da sua empresa dos últimos 12 meses a partir do mes passado
func (c *CalculoSimples) folhadepagamento12(faturamento Faturamento, mes int) float64 {
	var ret float64
	// pega N meses do faturamento passado
	for cm := mes; cm <= 11; cm++ {
		if c.faturamento12Meses.FolhaDePagamentoMes(cm) != nil {
			ret += *c.faturamento12Meses.FolhaDePagamentoMes(cm)
		}
	}
	// pega N meses do faturamento atual
	for cm := 0; cm < mes; cm++ {
		if faturamento.FolhaDePagamentoMes(cm) != nil {
			ret += *faturamento.FolhaDePagamentoMes(cm)
		}
	}
	return ret
}

func (c *CalculoSimples) calcImpostoPercentual(reparticaoImpostos SimplesReparticaoImpostos) float64 {
	ret := 0.0
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CPP) {
		ret += reparticaoImpostos.CPP
	}
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_ISS) {
		ret += reparticaoImpostos.ISS
	}
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CSLL) {
		ret += reparticaoImpostos.CSLL
	}
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_IRPJ) {
		ret += reparticaoImpostos.IRPJ
	}
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_Cofins) {
		ret += reparticaoImpostos.Cofins
	}
	if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_PIS) {
		ret += reparticaoImpostos.PIS
	}

	return ret
}

func (c *CalculoSimples) calculoMes(faturamento Faturamento, mes int) (*calculo_imposto.CalculoResultadoItem, error) {
	c.explain.Add(CalculoSimples_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		"Processando mês %<mes>d", map[string]interface{}{
			"mes": mes,
		})

	faturamentoAnual := c.rbt12(faturamento, mes)
	folhadepagamentoAnual := c.folhadepagamento12(faturamento, mes)

	var fatorR float64
	if folhadepagamentoAnual > 0 && faturamentoAnual == 0 {
		fatorR = 28.0
	} else if folhadepagamentoAnual == 0 && faturamentoAnual > 0 {
		fatorR = 1.0
	} else if folhadepagamentoAnual > 0 && faturamentoAnual > 0 {
		fatorR = folhadepagamentoAnual / faturamentoAnual * 100.0
	}

	anexo := c.anexo
	if c.anexoFatorR != nil && fatorR >= c.consts.Value(calculo_imposto.ConstItem_PERCENTUAL_FATOR_R) {
		c.explain.Add(CalculoSimples_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
			"Aplicando fator R %<fatorR>.2f (>= %<fatorRPercentual>.2f%%), selecionado "+
				"anexo '%<anexoFatorR>s' ao invés de '%<anexo>s'", map[string]interface{}{
				"fatorR":           fatorR,
				"fatorRPercentual": c.consts.Value(calculo_imposto.ConstItem_PERCENTUAL_FATOR_R),
				"anexoFatorR":      c.anexoFatorR.Titulo,
				"anexo":            anexo.Titulo,
			})

		anexo = c.anexoFatorR
	}

	c.explain.Add(CalculoSimples_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
		"Usando valor mês %<valorMes>.2f, rbt12 %<rbt12>.2f, folha de pagamento anual %<folhadepagamentoAnual>.2f, "+
			"fator R %<fatorR>.2f%%, anexo '%<anexo>s'", map[string]interface{}{
			"valorMes":              faturamento.ValorMes(mes),
			"rbt12":                 faturamentoAnual,
			"folhadepagamentoAnual": folhadepagamentoAnual,
			"fatorR":                fatorR,
			"anexo":                 anexo.Titulo,
		})

	// https://www.jornalcontabil.com.br/simples-nacional-saiba-como-funciona-a-exportacao/

	limiteAnterior := 0.0
	for faixaIdx, faixa := range anexo.Faixas {
		if faturamentoAnual >= limiteAnterior && faturamentoAnual <= faixa.LimiteMaximo {
			// aliquota efetiva leva em conta o RBT12 e a dedução da faixa
			aliquotaEfetiva := (((faturamentoAnual * (faixa.Aliquota / 100.0)) - faixa.ValorDeducao) / faturamentoAnual) * 100.0
			valorImpostoBase := faturamento.ValorMes(mes) * aliquotaEfetiva / 100.0

			ret := &calculo_imposto.CalculoResultadoItem{
				Periodo:       mes,
				TipoPeriodo:   calculo_imposto.TipoPeriodo_Mes,
				ValorOriginal: faturamento.ValorMes(mes),
				Extra: map[calculo_imposto.CalculoResultadoExtra]interface{}{
					calculo_imposto.CalculoResultadoExtra_Faixa: faixaIdx + 1,
					CalculoResultadoExtra_Simples_FatorR:        fatorR,
					CalculoResultadoExtra_Simples_ValorRBT12:    faturamentoAnual,
				},
				Impostos: map[calculo_imposto.TipoImposto]*calculo_imposto.CalculoResultadoImposto{},
			}

			// aplica cada um dos impostos aplicáveis
			var impostosExcluidos []string
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_IRPJ) {
				impostoIRPJ := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.ReparticaoImpostos.IRPJ,
					ValorImposto: valorImpostoBase * (faixa.ReparticaoImpostos.IRPJ / 100.0),
				}
				ret.Impostos[calculo_imposto.TipoImposto_IRPJ] = impostoIRPJ
				ret.ValorImposto += impostoIRPJ.ValorImposto
			} else {
				impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_IRPJ.String())
			}
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CSLL) {
				impostoCSLL := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.ReparticaoImpostos.CSLL,
					ValorImposto: valorImpostoBase * (faixa.ReparticaoImpostos.CSLL / 100.0),
				}
				ret.Impostos[calculo_imposto.TipoImposto_CSLL] = impostoCSLL
				ret.ValorImposto += impostoCSLL.ValorImposto
			} else {
				impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_CSLL.String())
			}
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_Cofins) {
				impostoCofins := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.ReparticaoImpostos.Cofins,
					ValorImposto: valorImpostoBase * (faixa.ReparticaoImpostos.Cofins / 100.0),
				}
				ret.Impostos[calculo_imposto.TipoImposto_Cofins] = impostoCofins
				ret.ValorImposto += impostoCofins.ValorImposto
			} else {
				impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_Cofins.String())
			}
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_PIS) {
				impostoPIS := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.ReparticaoImpostos.PIS,
					ValorImposto: valorImpostoBase * (faixa.ReparticaoImpostos.PIS / 100.0),
				}
				ret.Impostos[calculo_imposto.TipoImposto_PIS] = impostoPIS
				ret.ValorImposto += impostoPIS.ValorImposto
			} else {
				impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_PIS.String())
			}
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_CPP) {
				impostoCPP := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.ReparticaoImpostos.CPP,
					ValorImposto: valorImpostoBase * (faixa.ReparticaoImpostos.CPP / 100.0),
				}
				ret.Impostos[calculo_imposto.TipoImposto_CPP] = impostoCPP
				ret.ValorImposto += impostoCPP.ValorImposto
			} else {
				impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_CPP.String())
			}
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_ISS) {
				impostoISS := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.ReparticaoImpostos.ISS,
					ValorImposto: valorImpostoBase * (faixa.ReparticaoImpostos.ISS / 100.0),
				}
				ret.Impostos[calculo_imposto.TipoImposto_ISS] = impostoISS
				ret.ValorImposto += impostoISS.ValorImposto
			} else {
				impostosExcluidos = append(impostosExcluidos, calculo_imposto.TipoImposto_ISS.String())
			}

			aliquotaSimples := ret.ValorImposto / ret.ValorOriginal * 100.0

			if len(impostosExcluidos) == 0 {
				c.explain.Add(CalculoSimples_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
					"Cálculo com faixa %<faixa>d do anexo '%<anexo>s', alíquota %<aliquota>.2f, "+
						"alíquota efetiva %<aliquotaEfetiva>.2f, valor do imposto %<valorImposto>.2f", map[string]interface{}{
						"faixa":           faixaIdx + 1,
						"anexo":           anexo.Titulo,
						"aliquota":        faixa.Aliquota,
						"aliquotaEfetiva": aliquotaEfetiva,
						"valorImposto":    valorImpostoBase,
					})

			} else {
				c.explain.Add(CalculoSimples_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
					"Cálculo com faixa %<faixa>d do anexo '%<anexo>s', alíquota base %<aliquotaBase>.2f, "+
						"alíquota efetiva base %<aliquotaEfetivaBase>.2f, valor base do imposto %<valorImpostoBase>.2f", map[string]interface{}{
						"faixa":               faixaIdx + 1,
						"anexo":               anexo.Titulo,
						"aliquotaBase":        faixa.Aliquota,
						"aliquotaEfetivaBase": aliquotaEfetiva,
						"valorImpostoBase":    valorImpostoBase,
					})

				c.explain.Add(CalculoSimples_SOURCE, mes, calculo_imposto.TipoPeriodo_Mes,
					"Excluindo impostos '%<impostosExcluidos>s', aliquota efetiva %<aliquotaEfetiva>.2f, "+
						"valor imposto %<valorImposto>.2f", map[string]interface{}{
						"impostosExcluidos": strings.Join(impostosExcluidos, ", "),
						"aliquotaEfetiva":   aliquotaSimples,
						"valorImposto":      ret.ValorImposto,
					})
			}

			ret.Impostos[calculo_imposto.TipoImposto_SIMPLES] = &calculo_imposto.CalculoResultadoImposto{
				Aliquota:     aliquotaSimples,
				ValorImposto: ret.ValorImposto,
			}

			return ret, nil
		}
		limiteAnterior = faixa.LimiteMaximo + 0.01
	}
	return nil, fmt.Errorf("Valor anual %f está fora das faixas do simples", faturamentoAnual)
}

type CalculoSimplesOpt func(*CalculoSimples)

func WithCS_Consts(consts calculo_imposto.Consts) CalculoSimplesOpt {
	return func(c *CalculoSimples) {
		c.consts = consts
	}
}

func WithCS_AnexoFatorR(anexoFatorR *SimplesAnexo) CalculoSimplesOpt {
	return func(c *CalculoSimples) {
		c.anexoFatorR = anexoFatorR
	}
}

func WithCS_ImpostoAplicado(impostoAplicado calculo_imposto.ImpostoAplicado) CalculoSimplesOpt {
	return func(c *CalculoSimples) {
		c.impostoAplicado = impostoAplicado
	}
}

func WithCS_Explain(explain calculo_imposto.ExplainIntf) CalculoSimplesOpt {
	return func(c *CalculoSimples) {
		c.explain = explain
	}
}
