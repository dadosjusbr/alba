# CLI

Interface de linha de comando para interagir com as funcionalidades de persisitência e gerenciamento de [Pipelines](https://github.com/dadosjusbr/executor).

## Pipeline DadosJusBR
O Pipeline DadosJusBR tem os seguintes estágios:

### Coleta 
Etapa responsável por encontrar os dados no site oficial do órgão, fazer o download dos arquivos e traduzir as informações para um formato único do DadosJusBr [(*crawling result*)](https://github.com/dadosjusbr/storage/blob/master/agency.go#L27). 

- [DadosJusBR Collectors](https://github.com/dadosjusbr/coletores)

### Validação
Responsável por fazer validações nos dados de acordo a cada contexto.

- [*Work in progress*](https://github.com/dadosjusbr/coletores)

### Empacotamento
Responsável por padronizar os dados no formato de datapackages.

- [DadosJusBR Packager](https://github.com/dadosjusbr/coletores/tree/master/packager)

### Armazenamento
Responsável por armazenar os dados extraídos, além de versionar também os artefatos baixados e gerados durante a coleta.

- [DadosJusBR Store](https://github.com/dadosjusbr/coletores/tree/master/store)

### Estágio de tratamento de erros
- [DadosJusBR Store Error](https://github.com/dadosjusbr/coletores/tree/master/store-error)

### Definindo um Pipeline 

Para que a Alba seja capaz de gerenciar e automatizar a execução periódica de Pipeline é preciso que a definição para cada órgão siga o seguinte formato:

``` json
{
"name": "",
"default-base-dir": "",
"default-build-env": "chave=valor,chave=valor,...",
"default-run-env": "chave=valor,chave=valor,...",
"stages":[
    {
        "name": "",
        "dir": "",
        "base-dir": "", 
        "build-env": "chave=valor,chave=valor,...",
        "run-env": "chave=valor,chave=valor,..."
    },
    {
        "name": "",
        "dir": "",
        "base-dir": "", 
        "build-env": "chave=valor,chave=valor,...",
        "run-env": "chave=valor,chave=valor,..."
    },
],
"error-handler": {
    "name": "",
    "dir": "",
    "base-dir": "",
    "build-env": "chave=valor,chave=valor,...", 
    "run-env": "chave=valor,chave=valor,..."
},
"entity": "",
"city": "",
"fu": "",
"repo": "",
"frequency": "",
"start-day": "",
"limit-month-backward": "",
"limit-year-backward": ""
}
```

**Exemplo para o coletor do [Tribunal Regional do Trabalho - 13ª região](https://github.com/dadosjusbr/coletores/tree/master/trt13)**
``` json
{
"name": "Tribunal Regional do Trabalho 13ª Região",
"default-base-dir": "github.com/dadosjusbr/coletores/",
"default-build-env": "",
"default-run-env": "",
"stages":[
    {
        "name": "Coleta",
        "dir": "trt13",
        "base-dir": "",
        "build-env": "",
        "run-env": ""
    },
    {
        "name": "Empacotamento",
        "dir": "packager",
        "base-dir": "",
        "build-env": "",
        "run-env": ""
    },
    {
        "name": "Armazenamento",
        "dir": "store",
        "base-dir": "",
        "build-env": "",
        "run-env": ""
    }
],
"error-handler": {
    "name": "Armazenamento de Erros",
    "dir": "store-error",
    "base-dir": "",
    "build-env": "",
    "run-env": ""
},
"entity": "Tribunal Regional do Trabalho 13ª Região",
"city": "João Pessoa",
"fu": "PB",
"repo": "github.com/dadosjusbr/coletores",
"frequency": 30,
"start-day": 5,
"limit-month-backward": 2,
"limit-year-backward": 2018
```
**Todo**: Adicionar tag para GIT_COMMIT - Porque pode variar com o tempo \\
**Todo**: Adicionar tag para MES e ANO - Porque varia a cada execução

Um exemplo preenchido para cadastro pode ser visto nesse [arquivo](https://github.com/dadosjusbr/alba/blob/master/cli/collector/input.json).

---

## Execução da CLI

| Para que a CLI funcione corretamente é preciso que as instruções para a [configuração do ambiente](https://github.com/dadosjusbr/alba/blob/master/README.md) tenham sido concluídas. |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

Fazer o build do projeto criando um executável de nome **alba**:

`go build -o alba`

### Visualizar os comandos da CLI através do comando:**

`./alba`

### Cadastrar um Pipeline
Conforme o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/master/cli/collector/input.json).

*Exemplo para o coletor do [Tribunal Regional do Trabalho - 13ª região](https://github.com/dadosjusbr/coletores/tree/master/trt13)*

`./alba add-collector --from-file=collector/input.json`

### Executando um Pipeline