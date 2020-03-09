# Alba
## Sistema de Orquestração de Execuções DadosJusBR

Sistema para orquestração e escalonamento de execuções, visando a automatização de processos do [DadosJusBR](https://dadosjusbr.com/). A função principal da ferramenta é gerenciar o processo de libertação contínua de dados de remuneração do sistema de justiça brasileiro, que inclui as etapas de:

- **Coleta:** Etapa responsável por encontrar, fazer o download dos arquivos e consolidar/traduzir as informações para um formato único do DadosJusBr. Cada coletor é responsável por um determinado órgão e recebe como parâmetro um mês/ano e o identificador do órgão
- **Validação:** Responsável por fazer validações nos dados de acordo a cada contexto;
- **Empacotamento:** Responsável por padronizar os dados no formato de datapackages;
- **Armazenamento:** Responsável por armazenar os dados extraídos, além de versionar também os artefatos baixados e gerados durante a coleta; 

Esse projeto é financiado na modalidade de Flash Grants pela [Shuttleworth Foundation](https://www.shuttleworthfoundation.org/), a quem agredecemos muito pelo suporte e incentivo.

## Configuração de ambiente
Após instalar as ferramentas docker e docker-composer é possível levantar as instâncias do projeto utilizando:

 `docker-compose up -d`

 Esse comando vai levantar um container para o banco de dados mongodb e fazer o build de outro container com as configurações necessárias para execução dos pacotes em go.

Além de utilizar a linha de comando, é possível se conectar ao banco de dados utilizando a ferramenta [Mongo Compass Community](https://www.mongodb.com/download-center/compass?jmp=docs). Ao abrir a ferramenta deve-se selecionar o modo de autenticação Username / Password, onde `Username` é `root` e `Password` é `example`.

## Execução da CLI

É possível visualizar os comandos da CLI através do comando:

`docker-compose run golang go run cli/alba.go`

Para o cadastro de um coletor deve-se configurar um arquivo em formato JSON com as informações necessárias. Conforme o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/comando-cli/cli/input.json).

Para executar o cadastro:

`docker-compose run golang go run cli/alba.go addCollector --file="cli/input.json"`