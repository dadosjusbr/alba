[![Build Status](https://travis-ci.org/dadosjusbr/alba.svg?branch=master)](https://travis-ci.org/dadosjusbr/alba) [![codecov.io](http://codecov.io/github/dadosjusbr/alba/coverage.svg?branch=master)](http://codecov.io/github/dadosjusbr/alba?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/dadosjusbr/alba)](https://goreportcard.com/report/github.com/dadosjusbr/alba) [![GoDoc](https://godoc.org/github.com/dadosjusbr/alba?status.svg)](https://godoc.org/github.com/dadosjusbr/alba)

# Alba

## Sistema de Orquestração de Execuções DadosJusBR

Sistema para orquestração e escalonamento de execuções, visando a automatização de processos do [DadosJusBR](https://dadosjusbr.org/). A função principal da ferramenta é gerenciar o processo de libertação contínua de dados de remuneração do sistema de justiça brasileiro, que inclui as etapas de:

- **Coleta:** Etapa responsável por encontrar, fazer o download dos arquivos e consolidar/traduzir as informações para um formato único do DadosJusBr. Cada coletor é responsável por um determinado órgão e recebe como parâmetro um mês/ano e o identificador do órgão
- **Validação:** Responsável por fazer validações nos dados de acordo a cada contexto;
- **Empacotamento:** Responsável por padronizar os dados no formato de datapackages;
- **Armazenamento:** Responsável por armazenar os dados extraídos, além de versionar também os artefatos baixados e gerados durante a coleta; 

Esse projeto é financiado na modalidade de Flash Grants pela [Shuttleworth Foundation](https://www.shuttleworthfoundation.org/), a quem agradecemos muito pelo suporte e incentivo.


## Configuração do ambiente

1) Após realizar o git clone do projeto é necessário exportar a variável de ambiente para o servidor Mongo:

`export MONGODB=mongodb://<usuario>:<senha>@<ip-do-servidor>:<porta>`

Se for executar para o servidor Mongo configurado no [docker-composer.yml]() é só usar `export MONGODB mongodb://root:example@localhost:28017`

2) Levantar o container do banco de dados executando:
`docker-compose up -d`