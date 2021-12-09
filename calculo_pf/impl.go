package calculo_pf

import (
	"github.com/RangelReale/calculo-imposto/calculo_pj"
)

//**************
// Prolabore
//*************

//
// Prolabore_Static
//

type Prolabore_Static struct {
	meses       int
	valorMensal float64
}

type Prolabore_Static_Option func(*Prolabore_Static)

func NewProlabore_Static(meses int, opt ...Prolabore_Static_Option) *Prolabore_Static {
	ret := &Prolabore_Static{
		meses: meses,
	}
	for _, o := range opt {
		o(ret)
	}
	return ret
}

func WithPS_ValorMensal(valorMensal float64) Prolabore_Static_Option {
	return func(f *Prolabore_Static) {
		f.valorMensal = valorMensal
	}
}

func WithPS_ValorAnual(valorAnual float64) Prolabore_Static_Option {
	return func(f *Prolabore_Static) {
		f.valorMensal = valorAnual / float64(f.meses)
	}
}

func (f *Prolabore_Static) Meses() int {
	return f.meses
}

func (f *Prolabore_Static) ValorMes(mes int) float64 {
	return f.valorMensal
}

//
// Prolabore_Faturamento_Adapter
//

type Prolabore_Faturamento_Adapter struct {
	faturamento calculo_pj.Faturamento
}

func NewProlabore_Faturamento_Adapter(faturamento calculo_pj.Faturamento) *Prolabore_Faturamento_Adapter {
	return &Prolabore_Faturamento_Adapter{faturamento}
}

func (f *Prolabore_Faturamento_Adapter) Meses() int {
	return f.faturamento.Meses()
}

func (f *Prolabore_Faturamento_Adapter) ValorMes(mes int) float64 {
	ret := f.faturamento.FolhaDePagamentoMes(mes)
	if ret == nil {
		return 0
	}
	return *ret
}

//**************
// TabelaINSS
//*************

//
// TabelaINSS_RegimeGeral
//

type TabelaINSS_RegimeGeral struct {
	faixas []TabelaINSS_Faixa
}

// Fonte: https://ingracio.adv.br/contribuicoes-inss-2021/
func NewTabelaINSS_RegimeGeral(faixas []TabelaINSS_Faixa) *TabelaINSS_RegimeGeral {
	return &TabelaINSS_RegimeGeral{
		faixas: faixas,
	}
}

func (t *TabelaINSS_RegimeGeral) Taxa(valorMensal float64) (taxa float64, valorMax float64) {
	var faixaAnterior *TabelaINSS_Faixa
	for _, faixa := range t.faixas {
		if valorMensal <= faixa.LimiteMaximo {
			return faixa.Aliquota, faixa.LimiteMaximo
		}
		faixaAnterior = &faixa
	}
	// se o valor mensal for maior que o máximo, usa a última aliquota
	return faixaAnterior.Aliquota, faixaAnterior.LimiteMaximo
}

// Fonte: https://ingracio.adv.br/contribuicoes-inss-2021/
func NewTabelaINSS_RegimeGeral_2021() *TabelaINSS_RegimeGeral {
	return NewTabelaINSS_RegimeGeral([]TabelaINSS_Faixa{
		{
			LimiteMaximo: 1100.0,
			Aliquota:     7.5,
		},
		{
			LimiteMaximo: 2203.48,
			Aliquota:     9.0,
		},
		{
			LimiteMaximo: 3305.22,
			Aliquota:     12.0,
		},
		{
			LimiteMaximo: 6433.57,
			Aliquota:     14.0,
		},
	})
}

//
// TabelaINSS_Autonomo
//

// fonte: https://www.contabilizei.com.br/contabilidade-online/inss-autonomo/

type TabelaINSS_Autonomo struct {
	valorTeto float64
}

func NewTabelaINSS_Autonomo(valorTeto float64) *TabelaINSS_Autonomo {
	return &TabelaINSS_Autonomo{valorTeto}
}

func (t *TabelaINSS_Autonomo) Taxa(valorMensal float64) (taxa float64, valorMax float64) {
	return 20.0, t.valorTeto
}

func NewTabelaINSS_Autonomo_2021() *TabelaINSS_Autonomo {
	return NewTabelaINSS_Autonomo(6433.57)
}

//
// TabelaINSS_Autonomo_Simplificado
//

// fonte: https://www.contabilizei.com.br/contabilidade-online/inss-autonomo/

type TabelaINSS_Autonomo_Simplificado struct {
	valorMinimo float64
	aliquota    float64
}

func NewTabelaINSS_Autonomo_Simplificado(valorMinimo float64, aliquota float64) *TabelaINSS_Autonomo_Simplificado {
	return &TabelaINSS_Autonomo_Simplificado{valorMinimo, aliquota}
}

func (t *TabelaINSS_Autonomo_Simplificado) Taxa(valorMensal float64) (taxa float64, valorMax float64) {
	return t.aliquota, t.valorMinimo
}

func NewTabelaINSS_Autonomo_Simplificado_2021() *TabelaINSS_Autonomo_Simplificado {
	return &TabelaINSS_Autonomo_Simplificado{6433.57, 11.0}
}

func NewTabelaINSS_Autonomo_Simplificado_Minimo_2021() *TabelaINSS_Autonomo_Simplificado {
	return &TabelaINSS_Autonomo_Simplificado{1100.0, 11.0}
}

// fonte: https://www.contabilizei.com.br/contabilidade-online/como-funciona-o-recolhimento-do-inss-para-pessoa-juridica/
func NewTabelaINSS_PJ_2021() *TabelaINSS_Autonomo_Simplificado {
	return &TabelaINSS_Autonomo_Simplificado{6433.57, 11.0}
}
