package calculo_pf

// https://www.contabilizei.com.br/contabilidade-online/irpf-2021-o-que-e-descontado-na-fonte-e-como-funciona-o-ajuste-anual/

import (
	"fmt"
	"math"

	calculo_imposto "github.com/RangelReale/calculo-imposto"
)

type CalculoPF struct {
	irpf            *IRPF // precisa ser tabela MENSAL
	tabelaINSS      TabelaINSS
	impostoAplicado calculo_imposto.ImpostoAplicado
}

func NewCalculoPF(irpf *IRPF, opt ...CalculoPFOpt) *CalculoPF {
	ret := &CalculoPF{
		irpf:            irpf,
		impostoAplicado: &calculo_imposto.ImpostoAplicado_Normal{},
	}
	for _, o := range opt {
		o(ret)
	}
	return ret
}

func (c *CalculoPF) Calculo(prolabore Prolabore) (*calculo_imposto.CalculoResultado, error) {
	ret := &calculo_imposto.CalculoResultado{}

	for mes := 0; mes < prolabore.Meses(); mes++ {
		rmes, err := c.calculoMes(prolabore, mes)
		if err != nil {
			return nil, err
		}
		ret.Items = append(ret.Items, rmes)
	}

	return ret, nil
}

func (c *CalculoPF) calculoMes(prolabore Prolabore, mes int) (*calculo_imposto.CalculoResultadoItem, error) {
	if c.irpf.Tipo != IRPF_Tipo_Mensal {
		return nil, fmt.Errorf("Tabela de IRPF deve ser mensal para este cálculo")
	}

	valorMes := prolabore.ValorMes(mes)
	limiteAnterior := 0.0
	for _, faixa := range c.irpf.Faixas {
		if valorMes >= limiteAnterior && valorMes <= faixa.LimiteMaximo {
			ret := &calculo_imposto.CalculoResultadoItem{
				Periodo:       mes,
				TipoPeriodo:   calculo_imposto.TipoPeriodo_Mes,
				ValorOriginal: valorMes,
				Impostos:      map[calculo_imposto.TipoImposto]*calculo_imposto.CalculoResultadoImposto{},
			}

			// INSS
			var impostoINSS *calculo_imposto.CalculoResultadoImposto

			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_INSS) {
				if c.tabelaINSS == nil {
					return nil, fmt.Errorf("Tabela de INSS deve ser informada")
				}

				taxaINSS, maxINSS := c.tabelaINSS.Taxa(valorMes)
				impostoINSS = &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     taxaINSS,
					ValorImposto: math.Min(maxINSS, valorMes) * taxaINSS / 100.0,
				}
				if impostoINSS.ValorImposto > 0 {
					ret.Impostos[calculo_imposto.TipoImposto_INSS] = impostoINSS
					ret.ValorImposto += impostoINSS.ValorImposto
				}
			} else {
				impostoINSS = &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     0,
					ValorImposto: 0.0,
				}
			}

			// IRPF
			if c.impostoAplicado.ImpostoAplicavel(calculo_imposto.TipoImposto_IRPF) {
				// fonte: https://www.pontotel.com.br/calcular-irrf/
				impostoIRPJ := &calculo_imposto.CalculoResultadoImposto{
					Aliquota:     faixa.Aliquota,
					ValorImposto: ((valorMes - impostoINSS.ValorImposto) * faixa.Aliquota / 100.0) - faixa.ValorDeducao,
				}
				if impostoIRPJ.ValorImposto > 0 {
					ret.Impostos[calculo_imposto.TipoImposto_IRPF] = impostoIRPJ
					ret.ValorImposto += impostoIRPJ.ValorImposto
				}
			}

			return ret, nil
		}
		limiteAnterior = faixa.LimiteMaximo + 0.01
	}
	return nil, fmt.Errorf("Valor %f está fora das faixas do IRPF", valorMes)
}

type CalculoPFOpt func(*CalculoPF)

func WithCPF_TabelaINSS(tabelaINSS TabelaINSS) CalculoPFOpt {
	return func(c *CalculoPF) {
		c.tabelaINSS = tabelaINSS
	}
}

func WithCPF_ImpostoAplicado(impostoAplicado calculo_imposto.ImpostoAplicado) CalculoPFOpt {
	return func(c *CalculoPF) {
		c.impostoAplicado = impostoAplicado
	}
}
