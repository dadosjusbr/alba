[![Build Status](https://travis-ci.org/dadosjusbr/alba.svg?branch=master)](https://travis-ci.org/dadosjusbr/alba) [![codecov.io](http://codecov.io/github/dadosjusbr/alba/coverage.svg?branch=master)](http://codecov.io/github/dadosjusbr/alba?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/dadosjusbr/alba)](https://goreportcard.com/report/github.com/dadosjusbr/alba) [![GoDoc](https://godoc.org/github.com/dadosjusbr/alba?status.svg)](https://godoc.org/github.com/dadosjusbr/alba)

# Alba

## Sistema de Orquestração de Execuções DadosJusBR

Sistema para orquestração e escalonamento de execuções, visando a automatização de processos do [DadosJusBR](https://dadosjusbr.org/). A função principal da ferramenta é gerenciar o processo de libertação contínua de dados de remuneração do sistema de justiça brasileiro, que inclui as etapas de:

- **Coleta:** Etapa responsável por encontrar, fazer o download dos arquivos e consolidar/traduzir as informações para um formato único do DadosJusBr. Cada coletor é responsável por um determinado órgão e recebe como parâmetro um mês/ano e o identificador do órgão
- **Validação:** Responsável por fazer validações nos dados de acordo a cada contexto;
- **Empacotamento:** Responsável por padronizar os dados no formato de datapackages;
- **Armazenamento:** Responsável por armazenar os dados extraídos, além de versionar também os artefatos baixados e gerados durante a coleta; 

Esse projeto é financiado na modalidade de Flash Grants pela [Shuttleworth Foundation](https://www.shuttleworthfoundation.org/), a quem agradecemos muito pelo suporte e incentivo.

***

## Configuração do ambiente

### 1 - Variável de ambiente para o MongoDB

Após realizar o git clone do projeto é necessário exportar a variável de ambiente para o servidor de banco de dados MongoDB:

`export MONGODB=mongodb://<usuario>:<senha>@<ip-do-servidor>:<porta>`

É interessante fazer isso de forma permanente para que a informação não seja perdida toda vez que precisar reiniciar o computador. Se o seu sistema é Ubuntu, uma forma de fazer isso é editando o arquivo `~/.profile` e depois reiniciando o computador ou executando `source ~/.profile`.

- No caso de querer utilizar a versão do Mongo instalada na sua máquina **o passo 2 não é necessário** e a variável de ambiente deve ser montada de acordo com o usuário, senha, e porta configurados na hora da instalação, passando `localhost` como `<ip-do-servidor>`.

- Se preferir utilizar o servidor Mongo configurado no arquivo [docker-composer.yml](https://github.com/dadosjusbr/alba/blob/master/docker-compose.yml) é só usar `export MONGODB=mongodb://root:example@localhost:28017`

### 2 - Levantar o container do banco de dados executando:

Para levantar o container do banco de dados execute:

`docker-compose up -d`

> É possível visualizar as informações persisitidas no banco de dados através do terminal ou utilizando a ferramenta [Mongo Compass Community](https://www.mongodb.com/download-center/compass?jmp=docs). Uma vez utilizando servidor Mongo configurado no [docker-composer.yml](https://github.com/dadosjusbr/alba/blob/master/docker-compose.yml), ao abrir a ferramenta deve-se utilizar como porta `28017` e selecionar o modo de autenticação Username / Password, onde Username é `root` e Password é `example`.
