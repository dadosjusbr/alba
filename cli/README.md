# CLI

Interface de linha de comando para interagir com as funcionalidades de persisitência e gerenciamento de [Pipelines](https://github.com/dadosjusbr/executor).

## Definição de um Pipeline 

Para gerenciar e automatizar a execução periódica de um Pipeline é preciso que a definição siga o seguinte formato:

``` json
{
"name": "",
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
"id": "",
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
O dicionário que descreve essa estrutura está disponível em [dadosjusbr/alba/storage/pipeline.go]().
O campo do [executor.Pipeline.DefaultBaseDir](https://github.com/dadosjusbr/executor/blob/3f2bad506ad914557b101fd3f0d78b5c897d8ec3/pipeline.go#L35) não é passado na definição aqui porque ele é gerado a partir do download do repositório que é gerenciado pela Alba.

### Exemplo de Pipeline DadosJusBR

**Exemplo para o coletor do [Tribunal Regional do Trabalho - 13ª região](https://github.com/dadosjusbr/coletores/tree/master/trt13)**
``` json
{
"name": "Tribunal Regional do Trabalho 13ª Região",
"stages":[
    {
        "name": "Coleta",
        "dir": "trt13",
    },
    {
        "name": "Empacotamento",
        "dir": "packager",
    },
    {
        "name": "Armazenamento",
        "dir": "store",
        "run-env": "MONGODB_URI=,MONGODB_DBNAME=,MONGODB_MICOL=,MONGODB_AGCOL=,SWIFT_USERNAME=,SWIFT_APIKEY=,SWIFT_AUTHURL=,SWIFT_DOMAIN=,SWIFT_CONTAINER=" 
    }
],
"error-handler": {
    "name": "Armazenamento de Erros",
    "dir": "store-error",
    "run-env": "MONGODB_URI=,MONGODB_DBNAME=,MONGODB_MICOL=,MONGODB_AGCOL=,SWIFT_USERNAME=,SWIFT_APIKEY=,SWIFT_AUTHURL=,SWIFT_DOMAIN=,SWIFT_CONTAINER=" 

},
"id": "trt13",
"entity": "Tribunal Regional do Trabalho 13ª Região",
"city": "João Pessoa",
"fu": "PB",
"repo": "github.com/dadosjusbr/coletores",
"frequency": 30,
"start-day": 5,
"limit-month-backward": 2,
"limit-year-backward": 2018
}
```

Por padrão passamos para **todos** os estágios as seguintes variáveis de ambiente:
- `GIT_COMMIT` para o `build-env`
-  `YEAR`, `OUTPUT_FOLDER`, `MONTH` para o `run-env`

Elas são sobrescritas caso sejam preenchidas na definição do pipeline.
- O GIT_COMMIT é carregado a partir da última versão do código, que é baixada a cada execução
- O valor padrão do `OUTPUT_FOLDER` é `"/output"`.
- `YEAR` e `MONTH` são carregadas a partir de variáveis de ambiente.

---

## Execução da CLI

| Para que a CLI funcione corretamente é preciso que as instruções para a [configuração do ambiente](https://github.com/dadosjusbr/alba/blob/master/README.md) tenham sido concluídas. |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

Fazendo o build do projeto criando um executável de nome **alba**:

`go build -o alba`

### Visualizar os comandos da CLI através do comando:

`./alba`

### Cadastrar um Pipeline
Passando como paraâmetro o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/master/cli/pipeline/pipeline-example.json). 
Nesse arquivo descrevemos os pipelines para os coletores do [Tribunal Regional do Trabalho - 13ª região](https://github.com/dadosjusbr/coletores/tree/master/trt13) e [Ministério Público da Paraíba](https://github.com/dadosjusbr/coletores/tree/master/mppb).

`./alba add --from-file pipeline/pipeline-example.json`

### Executando um Pipeline 

`./alba run --id trt13`

Sendo `trt13` o id do pipeline para o Tribunal Regional do Trabalho - 13ª região.
