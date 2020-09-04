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
        "build-env": "GIT_COMMIT=",
        "run-env": "--mes=,--ano="
    },
    {
        "name": "Empacotamento",
        "dir": "packager",
        "run-env": "OUTPUT_FOLDER=/output"
    },
    {
        "name": "Armazenamento",
        "dir": "store",
        "run-env": "OUTPUT_FOLDER=/output,MONGODB_URI=,MONGODB_DBNAME=,MONGODB_MICOL=,MONGODB_AGCOL=,SWIFT_USERNAME=,SWIFT_APIKEY=,SWIFT_AUTHURL=,SWIFT_DOMAIN=,SWIFT_CONTAINER=" 
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
Os parâmetros `GIT_COMMIT`, `--mes` e `--ano` são padrões do contexto do [DadosJusBR](https://github.com/dadosjusbr/coletores/blob/master/TUTORIAL.md). Os valores desses parâmetros podem variar ao longo do tempo e a cada execução, por isso, na definição de um pipeline DadosJusBr deixamos seus valores vazios:

``` json
"default-build-env": "GIT_COMMIT=",
"default-build-env": "--mes=,--ano="
``` 

E consideramos as seguintes regras de negócio:
- Se o `GIT_COMMIT` não estiver preenchido o pacote cli faz o download a última versão do código (a partir do endereço em `repo`) e carrega a informação com o `git rev-list -1 HEAD`.
- No caso de `--mes` e `--ano`:
    - Quando a execução for iniciada via cli, os valores devem ser passados por parâmetro [no comando]().
    - Quando a execução for iniciada pelo [worker](), ele irá avaliar quais são os valores a partir de execuções anteriores.

---

## Execução da CLI

| Para que a CLI funcione corretamente é preciso que as instruções para a [configuração do ambiente](https://github.com/dadosjusbr/alba/blob/master/README.md) tenham sido concluídas. |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

Fazer o build do projeto criando um executável de nome **alba**:

`go build -o alba`

### Visualizar os comandos da CLI através do comando:

`./alba`

### Cadastrar um Pipeline
Passando como paraâmetro o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/master/cli/collector/.pipeline.json). 
Nesse arquivo descrevemos os pipelines para os coletores do [Tribunal Regional do Trabalho - 13ª região](https://github.com/dadosjusbr/coletores/tree/master/trt13) e [Ministério Público da Paraíba](https://github.com/dadosjusbr/coletores/tree/master/mppb).

`./alba add --from-file pipeline/pipeline-example.json`

### Executando um Pipeline
