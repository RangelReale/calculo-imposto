package calculo_pf

// fonte: https://www.debit.com.br/tabelas/tabelas-irrf.php
var IRPF2021_Mensal = &IRPF{
	Ano:  2021,
	Tipo: IRPF_Tipo_Mensal,
	Faixas: []IRPFFaixa{
		{LimiteMaximo: 1_903.98, Aliquota: 0, ValorDeducao: 0},
		{LimiteMaximo: 2_826.65, Aliquota: 7.5, ValorDeducao: 142.80},
		{LimiteMaximo: 3_751.05, Aliquota: 15, ValorDeducao: 354.80},
		{LimiteMaximo: 4_664.68, Aliquota: 22.5, ValorDeducao: 636.13},
		{LimiteMaximo: 999999999999999999.99, Aliquota: 27.5, ValorDeducao: 869.36},
	},
}

var IRPF_Mensal_Atual = IRPF2021_Mensal

var IRPF2021_Anual = &IRPF{
	Ano:  2021,
	Tipo: IRPF_Tipo_Anual,
	Faixas: []IRPFFaixa{
		{LimiteMaximo: 22_847.76, Aliquota: 0, ValorDeducao: 0},
		{LimiteMaximo: 33_919.80, Aliquota: 7.5, ValorDeducao: 1_713.58},
		{LimiteMaximo: 45_012.60, Aliquota: 15, ValorDeducao: 4_257.57},
		{LimiteMaximo: 55_976.16, Aliquota: 22.5, ValorDeducao: 7_633.51},
		{LimiteMaximo: 999999999999999999.99, Aliquota: 27.5, ValorDeducao: 10_432.32},
	},
}

var IRPF_Anual_Atual = IRPF2021_Anual
