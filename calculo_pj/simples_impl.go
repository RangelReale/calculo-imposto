package calculo_pj

// Fonte: http://www.planalto.gov.br/CCivil_03/leis/LCP/Lcp123.htm

//
// SIMPLES 2018
//

// Fonte: https://contabilizei.com.br/contabilidade-online/anexo-3-simples-nacional
var Simples2018AnexoIII = &SimplesAnexo{
	Titulo: "Anexo III",
	Numero: "III",
	Faixas: []SimplesFaixa{
		{
			LimiteMaximo: 180_000.0,
			Aliquota:     6,
			ValorDeducao: 0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    43.4,
				ISS:    33.5,
				CSLL:   3.5,
				IRPJ:   4,
				Cofins: 12.82,
				PIS:    2.78,
			},
		},
		{
			LimiteMaximo: 360_000.0,
			Aliquota:     11.2,
			ValorDeducao: 9_360.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    43.4,
				ISS:    32,
				CSLL:   3.5,
				IRPJ:   4,
				Cofins: 14.05,
				PIS:    3.05,
			},
		},
		{
			LimiteMaximo: 720_000.0,
			Aliquota:     13.5,
			ValorDeducao: 17_640.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    43.4,
				ISS:    32.5,
				CSLL:   3.5,
				IRPJ:   4,
				Cofins: 13.64,
				PIS:    2.96,
			},
		},
		{
			LimiteMaximo: 1_800_000.0,
			Aliquota:     16,
			ValorDeducao: 35_640.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    43.4,
				ISS:    32.5,
				CSLL:   3.5,
				IRPJ:   4,
				Cofins: 13.64,
				PIS:    2.96,
			},
		},
		{
			LimiteMaximo: 3_600_000.0,
			Aliquota:     21,
			ValorDeducao: 125_640.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    43.4,
				ISS:    33.5, // TODO: O percentual efetivo máximo devido ao ISS será de 5%.
				CSLL:   3.5,
				IRPJ:   4,
				Cofins: 12.82,
				PIS:    2.78,
			},
		},
		{
			LimiteMaximo: 4_800_000.0,
			Aliquota:     33,
			ValorDeducao: 648_000.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    30.5,
				ISS:    0,
				CSLL:   15,
				IRPJ:   35,
				Cofins: 16.03,
				PIS:    3.47,
			},
		},
	},
}

// Fonte: https://www.contabilizei.com.br/contabilidade-online/anexo-5-simples-nacional/
var Simples2018AnexoV = &SimplesAnexo{
	Titulo: "Anexo V",
	Numero: "V",
	Faixas: []SimplesFaixa{
		{
			LimiteMaximo: 180_000.0,
			Aliquota:     15.5,
			ValorDeducao: 0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    28.85,
				ISS:    14,
				CSLL:   15,
				IRPJ:   25,
				Cofins: 14.1,
				PIS:    3.05,
			},
		},
		{
			LimiteMaximo: 360_000.0,
			Aliquota:     18,
			ValorDeducao: 4_500.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    27.85,
				ISS:    17,
				CSLL:   15,
				IRPJ:   23,
				Cofins: 14.10,
				PIS:    3.05,
			},
		},
		{
			LimiteMaximo: 720_000.0,
			Aliquota:     19.5,
			ValorDeducao: 9_900.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    23.85,
				ISS:    19,
				CSLL:   15,
				IRPJ:   24,
				Cofins: 14.92,
				PIS:    3.23,
			},
		},
		{
			LimiteMaximo: 1_800_000.0,
			Aliquota:     20.5,
			ValorDeducao: 17_100,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    23.85,
				ISS:    21,
				CSLL:   15,
				IRPJ:   21,
				Cofins: 15.74,
				PIS:    3.41,
			},
		},
		{
			LimiteMaximo: 3_600_000.0,
			Aliquota:     23,
			ValorDeducao: 62_100.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    23.85,
				ISS:    23.5,
				CSLL:   12.5,
				IRPJ:   23,
				Cofins: 14.1,
				PIS:    3.05,
			},
		},
		{
			LimiteMaximo: 4_800_000.0,
			Aliquota:     30.5,
			ValorDeducao: 540_000.0,
			ReparticaoImpostos: SimplesReparticaoImpostos{
				CPP:    29.5,
				ISS:    0,
				CSLL:   15.5,
				IRPJ:   35,
				Cofins: 16.44,
				PIS:    3.56,
			},
		},
	},
}

var Simples2018 = &Simples{
	AnoVigencia: 2018,
	MesVigencia: 1,
	Anexos: map[string]*SimplesAnexo{
		"III": Simples2018AnexoIII,
		"V":   Simples2018AnexoV,
	},
}

var Simples_Atual = Simples2018
