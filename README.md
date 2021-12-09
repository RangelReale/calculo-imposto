# Rotinas de cálculo de impostos do Brasil Pessoas Físicas e Jurídicas

Biblioteca para realização de cálculos de impostos do Brasil:

* Pessoa física: IRPF e INSS
* Pessoa jurídica: Simples, Lucro Presumido, Lucro Real

## Pessoa Física

```go
// Simula um prolabore de 1100.00 mensais por 12 meses
prolabore := calculo_pf.NewProlabore_Static(12, calculo_pf.WithPS_ValorMensal(4000.00))

// Cria um cálculo para pessoa física usando a tabela mensal de IRPF 2021, e a
// tabela INSS para PJ de 2021
calc := calculo_pf.NewCalculoPF(calculo_pf.IRPF2021_Mensal, 
    calculo_pf.WithCPF_TabelaINSS(calculo_pf.NewTabelaINSS_PJ_2021()))

// Executa o cálculo
ret, err := calc.Calculo(prolabore)
if err != nil {
    return err
}

// Imprime cada um dos meses
for _, cm := range ret.Items {
    fmt.Printf("Mes: %d\n", cm.Periodo)
    fmt.Printf("Prolabore: %.2f\n", cm.ValorOriginal)
    // Imprime o total de impostos pago no mes
    fmt.Printf("Imposto: %.2f (%.2f%%)\n", cm.ValorImposto, cm.AliquotaImposto())
    // Imprime valor de cada imposto usando uma ordem pré-definida
    for _, ti := range calculo_imposto.TipoImpostoLista {
        if imp, ok := cm.Impostos[ti]; ok {
            if imp.ValorImposto > 0 {
                fmt.Printf("ValorImposto %s: %.2f (aliq: %.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
            }
        }
    }
    fmt.Printf("%s\n", strings.Repeat("-", 10))
}

// Imprime o total de todos os meses
fmt.Printf("\n%s TOTAL %s\n", strings.Repeat("-", 5), strings.Repeat("-", 5))
rettotal := ret.Total()
fmt.Printf("Imposto: %.2f (%.2f%%)\n", rettotal.ValorImposto, rettotal.AliquotaImposto())
for _, ti := range calculo_imposto.TipoImpostoLista {
    if imp, ok := rettotal.Impostos[ti]; ok {
        if imp.ValorImposto > 0 {
            fmt.Printf("ValorImposto %s: %.2f (aliq: %.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
        }
    }
}
```

Resultado:

```
Mes: 0
Prolabore: 4000.00                      
Imposto: 604.87 (15.12%)                
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------                              
Mes: 1                                  
Prolabore: 4000.00                      
Imposto: 604.87 (15.12%)                
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------                              
Mes: 2                                  
Prolabore: 4000.00                      
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 3
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 4
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 5
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 6
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 7
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 8
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 9
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 10
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------
Mes: 11
Prolabore: 4000.00
Imposto: 604.87 (15.12%)
ValorImposto IRPF: 164.87 (aliq: 22.50%)
ValorImposto INSS: 440.00 (aliq: 11.00%)
----------

----- TOTAL -----
Imposto: 7258.44 (15.12%)
ValorImposto IRPF: 1978.44 (aliq: 22.50%)
ValorImposto INSS: 5280.00 (aliq: 11.00%)
```

## Pessoa Jurídica - Simples

```go
// Simula um faturamento de 20000.00 mensais por 12 meses, com uma folha de pagamento mensal de 28% de 20000.00
fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorMensal(20_000.0),
    calculo_pj.WithFS_FolhadePagamentoMensal(20_000.0*0.28))

// Cria um cálculo de Simples usando o Anexo V de 2018, ou o Anexo III de 2018 caso seja válido o Fator R,
// usando o imposto de exportação (sem PIS, Cofins, nem ISS)
calc := calculo_pj.NewCalculoSimples(fat, calculo_pj.Simples2018AnexoV, 
	calculo_pj.WithCS_AnexoFatorR(calculo_pj.Simples2018AnexoIII),
    calculo_pj.WithCS_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}))

// Executa o cálculo
ret, err := calc.Calculo(fat)
if err != nil {
    return err
}

// Imprime cada um dos meses
for _, cm := range ret.Items {
    fmt.Printf("Mes: %d\n", cm.Periodo)
    fmt.Printf("Faturamento: %.2f\n", cm.ValorOriginal)
    // Imprime o total de impostos pago no mes
    fmt.Printf("Imposto: %.2f (%.2f%%)\n", cm.ValorImposto, cm.AliquotaImposto())
    // Imprime o fator R do mes
    fmt.Printf("Fator R: %.2f\n", cm.Extra[calculo_pj.CalculoResultadoExtra_Simples_FatorR].(float64))
    // Imprime valor de cada imposto usando uma ordem pré-definida
    for _, ti := range calculo_imposto.TipoImpostoLista {
        if imp, ok := cm.Impostos[ti]; ok {
            if imp.ValorImposto > 0 {
                fmt.Printf("ValorImposto %s: %.2f (aliq: %.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
            }
        }
    }
    fmt.Printf("%s\n", strings.Repeat("-", 10))
}

// Imprime o total de todos os meses
fmt.Printf("\n%s TOTAL %s\n", strings.Repeat("-", 5), strings.Repeat("-", 5))
rettotal := ret.Total()
fmt.Printf("Imposto: %.2f (%.2f%%)\n", rettotal.ValorImposto, rettotal.AliquotaImposto())
for _, ti := range calculo_imposto.TipoImpostoLista {
    if imp, ok := rettotal.Impostos[ti]; ok {
        if imp.ValorImposto > 0 {
            fmt.Printf("ValorImposto %s: %.2f (aliq: %.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
        }
    }
}
```

Resultado:

```
Mes: 0
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)     
ValorImposto CPP: 633.64 (aliq: 43.40%)    
ValorImposto CSLL: 51.10 (aliq: 3.50%)     
ValorImposto PIS: 44.53 (aliq: 3.05%)      
ValorImposto Cofins: 205.13 (aliq: 14.05%) 
ValorImposto ISS: 467.20 (aliq: 32.00%)    
----------
Mes: 1
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 2
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 3
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 4
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 5
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 6
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 7
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 8
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 9
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 10
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------
Mes: 11
Faturamento: 20000.00
Imposto: 1460.00 (7.30%)
Fator R: 28.00
ValorImposto Simples: 1460.00 (aliq: 7.30%)
ValorImposto IRPJ: 58.40 (aliq: 4.00%)
ValorImposto CPP: 633.64 (aliq: 43.40%)
ValorImposto CSLL: 51.10 (aliq: 3.50%)
ValorImposto PIS: 44.53 (aliq: 3.05%)
ValorImposto Cofins: 205.13 (aliq: 14.05%)
ValorImposto ISS: 467.20 (aliq: 32.00%)
----------

----- TOTAL -----
Imposto: 17520.00 (7.30%)
ValorImposto Simples: 17520.00 (aliq: 7.30%)
ValorImposto IRPJ: 700.80 (aliq: 4.00%)
ValorImposto CPP: 7603.68 (aliq: 43.40%)
ValorImposto CSLL: 613.20 (aliq: 3.50%)
ValorImposto PIS: 534.36 (aliq: 3.05%)
ValorImposto Cofins: 2461.56 (aliq: 14.05%)
ValorImposto ISS: 5606.40 (aliq: 32.00%)
```

## Pessoa Jurídica: Lucro Presumido

```go
// Simula um faturamento de 840000.00 anuais por 12 meses, com uma folha de pagamento anual de 1100.00 * 12
fat := calculo_pj.NewFaturamento_Static(12, calculo_pj.WithFS_ValorAnual(840_000.00),
    calculo_pj.WithFS_FolhadePagamentoAnual(1100.00 * 12))

// Cria um cálculo de Lucro Presumido, com alíquota de IRPJ presumido de prestação de serviços,
// com alíquota de CSLL presumido de prestação de serviços, com aliquota de ISS de desenvolvimento de software, 
// usando o imposto de exportação (sem PIS, Cofins, nem ISS)
calc := calculo_pj.NewCalculoLucroPresumido(
    calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_IRPJ_PRESTACAO_DE_SERVICOS),
    calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_LUCRO_PRESUMIDO_CSLL_PRESTACAO_DE_SERVICOS),
    calculo_pj.WithCLP_ISS(calculo_imposto.Consts_Atual.Value(calculo_imposto.ConstItem_ALIQUOTA_ISS_SOFTWARE)),
    calculo_pj.WithCLP_ImpostoAplicado(calculo_imposto.ImpostoAplicado_Exportacao{}))

// Executa o cálculo
ret, err := calc.Calculo(fat)
if err != nil {
    return err
}

// Imprime cada um dos meses
for _, cm := range ret.Items {
    fmt.Printf("Mes: %d\n", cm.Periodo)
    fmt.Printf("Faturamento: %.2f\n", cm.ValorOriginal)
    // Imprime o total de impostos pago no mes
    fmt.Printf("Imposto: %.2f (%.2f%%)\n", cm.ValorImposto, cm.AliquotaImposto())
    if fp, ok := cm.Extra[calculo_imposto.CalculoResultadoExtra_ValorFolhaDePagamento]; ok {
        fmt.Printf("Folha de pagamento: %.2f\n", fp.(float64))
    }
    // Imprime valor de cada imposto usando uma ordem pré-definida
    for _, ti := range calculo_imposto.TipoImpostoLista {
        if imp, ok := cm.Impostos[ti]; ok {
            if imp.ValorImposto > 0 {
                fmt.Printf("ValorImposto %s: %.2f (aliq: %.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
            }
        }
    }
    fmt.Printf("%s\n", strings.Repeat("-", 10))
}

// Imprime o total de todos os meses
fmt.Printf("\n%s TOTAL %s\n", strings.Repeat("-", 5), strings.Repeat("-", 5))
rettotal := ret.Total()
fmt.Printf("Imposto: %.2f (%.2f%%)\n", rettotal.ValorImposto, rettotal.AliquotaImposto())
for _, ti := range calculo_imposto.TipoImpostoLista {
    if imp, ok := rettotal.Impostos[ti]; ok {
        if imp.ValorImposto > 0 {
            fmt.Printf("ValorImposto %s: %.2f (aliq: %.2f%%)\n", ti.String(), imp.ValorImposto, imp.Aliquota)
        }
    }
}
```

Resultado:

```
Mes: 0
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 1
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 2
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 3
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 4
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 5
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 6
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 7
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 8
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 9
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 10
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------
Mes: 11
Faturamento: 70000.00
Imposto: 5836.00 (8.34%)
Folha de pagamento: 1100.00
ValorImposto IRPJ: 3360.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 240.00 (aliq: 10.00%)
ValorImposto CPP: 220.00 (aliq: 20.00%)
ValorImposto CSLL: 2016.00 (aliq: 9.00%)
----------

----- TOTAL -----
Imposto: 70032.00 (8.34%)
ValorImposto IRPJ: 40320.00 (aliq: 15.00%)
ValorImposto IRPJ Adicional: 2880.00 (aliq: 10.00%)
ValorImposto CPP: 2640.00 (aliq: 20.00%)
ValorImposto CSLL: 24192.00 (aliq: 9.00%)
```

### Autor

Rangel Reale <rangelreale@gmail.com>
