# CLI

Interface de linha de comando para interagir com as funcionalidades de persisitência e gerenciamento de coletores.

## Execução da CLI

Fazer o build do projeto:
`go build -o alba`

Visualizar os comandos da CLI através do comando:
`./alba`

Para o cadastro de um coletor pode-se:

1) Fazer o cadastro via parâmetros:

 ```
./alba add-collector --id=mppb --entity="Ministério Público da Paraíba" --city="João Pessoa" --
fu=PB --path="github.com/dadosjusbr/coletores/mppb" --frequency=30 --startDay=5 --limitMonthBackward=1 --limitYearBackward=2018
```

2) Configurar um arquivo em formato JSON com as informações necessárias, conforme o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/comando-cli/cli/input.json).

`./alba add-collector --from-file=collector/input.json`
