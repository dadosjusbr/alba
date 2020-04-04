# CLI

Interface de linha de comando para interagir com as funcionalidades de persisitência e gerenciamento de coletores.

Um coletor é uma entidade capaz de extrair informações referentes às remunerações do sistema de justiça brasileiro. Ele é responsável por duas tarefas: baixar os dados do site oficial do órgão e convertê-los para o formato padronizado de [resultado de coleta](https://github.com/dadosjusbr/storage/blob/master/agency.go#L27) (_crawling result_). 

Para que a Alba seja capaz de gerenciar e automatizar a execução periódica dos coletores de cada órgão é preciso que as informações do cadastro de cada coletor siga o seguinte formato:

```
"id": "trt13", \\ Iniciais da entidade, como trt13.
"entity": "Tribunal Regional do Trabalho 13ª Região", \\ Entidade da qual o coletor extrai dados como 'Tribunal Regional do Trabalho 13 ° Região'.
"city": "João Pessoa", \\ Cidade da entidade da qual o coletor extrai dados.
"fu": "PB", \\ Unidade de federação da entidade da qual o coletor extrai dados.
"path": "github.com/dadosjusbr/coletores/trt13", \\ Caminho do repositório do coletor. Usando o padrão de importação do golang como 'github.com/dadosjusbr/coletores/trt13'.
"frequency": 30, \\ Frequência de execução do coletor em dias. Os valores devem estar entre 1 e 30. Para ser executado mensalmente, deve ser preenchido com '30'.
"startDay": 5, \\ Dia do mês para a execução do coletor. Os valores devem estar entre 1 e 31.
"limitMonthBackward": 2, \\ O mês limite para o qual o coletor deve ser executado de forma histórica. Exemplo: "Gostaria que o coletor que acabei de criar em Abril de 2020 fosse executado de forma histórica até Janeiro de 2019 visto que a estrutura dos dados passados é a mesma."
"limitYearBackward": 2018 \\ O ano limite até o qual o coletor deve ser executado de forma histórica, semelhante ao campo acima.
```

***

## Execução da CLI

| Para que a CLI funcione corretamente é preciso que as instruções para a [configuração do ambiente](https://github.com/dadosjusbr/alba/blob/master/README.md) tenham sido concluídas. |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

Fazer o build do projeto criando um executável de nome **alba**:
`go build -o alba`

Visualizar os comandos da CLI através do comando:

`./alba`

O cadastro de um coletor na base de dados da Alba, conforme a estrutura explicada acima, pode ser realizado de duas formas:

1) Fazer o cadastro via parâmetros:
*Exemplo para o coletor do [Ministério Público da Paraíba](https://github.com/dadosjusbr/coletores/tree/master/mppb)*
 ```
./alba add-collector --id=mppb --entity="Ministério Público da Paraíba" --city="João Pessoa" --
fu=PB --path="github.com/dadosjusbr/coletores/mppb" --frequency=30 --start-day=5 --limit-month-backward=1 --limit-year-backward=2018
```

2) Configurar um arquivo em formato JSON com as informações necessárias, conforme o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/comando-cli/cli/input.json).
*Exemplo para o coletor do [Tribunal Regional do Trabalho - 13ª região](https://github.com/dadosjusbr/coletores/tree/master/trt13)*

`./alba add-collector --from-file=collector/input.json`
