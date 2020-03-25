# CLI

Interface de linha de comando para interagir com as funcionalidades de persisitência e gerenciamento de coletores.

## Execução da CLI

É possível visualizar os comandos da CLI através do comando:

`docker-compose run golang go run cli/alba.go`

Para o cadastro de um coletor pode-se:

1) Fazer o cadastro via parâmetros:
 ```
docker-compose run golang go run cli/alba.go add --id=mppb --entity="Ministério Público da Paraíba" --city="João Pessoa" --
fu=PB --path="github.com/dadosjusbr/coletores/mppb" --frequency=30 --startDay=5 --limitMonthBackward=1 --limitYearBackward=2018
```

2) Configurar um arquivo em formato JSON com as informações necessárias, conforme o [arquivo de exemplo](https://github.com/dadosjusbr/alba/blob/comando-cli/cli/input.json).

`docker-compose run golang go run cli/alba.go add fromFile --file="cli/input.json"`

### Testes

Para execução dos testes de cada comando da CLI:

`docker-compose run golang bash -c "cd cli/ && go test -v"`