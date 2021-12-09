package calculo_imposto

type ConstItem string

const (
	ConstItem_ALIQUOTA_CPP                          ConstItem = "aliquota-cpp"
	ConstItem_ALIQUOTA_PIS_CUMULATIVO               ConstItem = "aliquota-pis-cumulativo"
	ConstItem_ALIQUOTA_PIS_NAO_CUMULATIVO           ConstItem = "aliquota-pis-nao-cumulativo"
	ConstItem_ALIQUOTA_COFINS_CUMULATIVO            ConstItem = "aliquota-cofins-cumulativo"
	ConstItem_ALIQUOTA_COFINS_NAO_CUMULATIVO        ConstItem = "aliquota-cofins-nao-cumulativo"
	ConstItem_ALIQUOTA_CSLL                         ConstItem = "aliquota-csll"
	ConstItem_ALIQUOTA_IRPJ                         ConstItem = "aliquota-irpj"
	ConstItem_ALIQUOTA_IRPJ_ADICIONAL               ConstItem = "aliquota-irpj-adicional"
	ConstItem_IRPJ_ADICIONAL_VALOR_MES              ConstItem = "irpj-adicional-valor-mes"
	ConstItem_ALIQUOTA_INSS_INDIVIDUAL              ConstItem = "aliquota-inss-individual"
	ConstItem_ALIQUOTA_INSS_INDIVIDUAL_SIMPLIFICADO ConstItem = "aliquota-inss-individual-simplificado"
	ConstItem_ALIQUOTA_INSS_PJ                      ConstItem = "aliquota-inss-pj"
	ConstItem_PERCENTUAL_FATOR_R                    ConstItem = "percentual-fator-r"
	ConstItem_ALIQUOTA_ISS_SOFTWARE                 ConstItem = "aliquota-iss-software"
	ConstItem_ALIQUOTA_ICMS_SOFTWARE                ConstItem = "aliquota-icms-software"

	ConstItem_LUCRO_PRESUMIDO_IRPJ_REVENDA_COMBUSTIVEIS   ConstItem = "lucro-presumido-irpj-revenda-combustiveis"
	ConstItem_LUCRO_PRESUMIDO_IRPJ_REGRA_GERAL            ConstItem = "lucro-presumido-irpj-regra-geral"
	ConstItem_LUCRO_PRESUMIDO_IRPJ_SERVICOS_DE_TRANSPORTE ConstItem = "lucro-presumido-irpj-servicos-de-transporte"
	ConstItem_LUCRO_PRESUMIDO_IRPJ_PRESTACAO_DE_SERVICOS  ConstItem = "lucro-presumido-irpj-prestacao-de-servicos"

	ConstItem_LUCRO_PRESUMIDO_CSLL_REGRA_GERAL           ConstItem = "lucro-presumido-csll-regra-geral"
	ConstItem_LUCRO_PRESUMIDO_CSLL_PRESTACAO_DE_SERVICOS ConstItem = "lucro-presumido-csll-prestacao-de-servicos"
)

type Consts interface {
	Exists(item ConstItem) bool
	Value(item ConstItem) float64
}
