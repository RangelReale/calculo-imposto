package calculo_pj

//**************
// Faturamento
//*************

//
// Faturamento_Static
//

type Faturamento_Static struct {
	meses                  int
	valorMensal            float64
	lucroMensal            *float64
	folhadepagamentoMensal *float64
}

type Faturamento_Static_Option func(*Faturamento_Static)

func NewFaturamento_Static(meses int, opt ...Faturamento_Static_Option) *Faturamento_Static {
	ret := &Faturamento_Static{
		meses: meses,
	}
	for _, o := range opt {
		o(ret)
	}
	return ret
}

func WithFS_ValorMensal(valorMensal float64) Faturamento_Static_Option {
	return func(f *Faturamento_Static) {
		f.valorMensal = valorMensal
	}
}

func WithFS_ValorAnual(valorAnual float64) Faturamento_Static_Option {
	return func(f *Faturamento_Static) {
		f.valorMensal = valorAnual / float64(f.meses)
	}
}

func WithFS_LucroMensal(lucroMensal float64) Faturamento_Static_Option {
	return func(f *Faturamento_Static) {
		f.lucroMensal = &lucroMensal
	}
}

func WithFS_LucroAnual(lucroAnual float64) Faturamento_Static_Option {
	return func(f *Faturamento_Static) {
		f.lucroMensal = new(float64)
		*f.lucroMensal = lucroAnual / float64(f.meses)
	}
}

func WithFS_FolhadePagamentoMensal(folhadepagamentoMensal float64) Faturamento_Static_Option {
	return func(f *Faturamento_Static) {
		f.folhadepagamentoMensal = &folhadepagamentoMensal
	}
}

func WithFS_FolhadePagamentoAnual(folhadepagamentoAnual float64) Faturamento_Static_Option {
	return func(f *Faturamento_Static) {
		f.folhadepagamentoMensal = new(float64)
		*f.folhadepagamentoMensal = folhadepagamentoAnual / float64(f.meses)
	}
}

func (f *Faturamento_Static) Meses() int {
	return f.meses
}

func (f *Faturamento_Static) ValorMes(mes int) float64 {
	return f.valorMensal
}

func (f *Faturamento_Static) LucroMes(mes int) *float64 {
	return f.lucroMensal
}

func (f *Faturamento_Static) FolhaDePagamentoMes(mes int) *float64 {
	return f.folhadepagamentoMensal
}
