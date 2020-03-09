package main

import (
	"testing"

	"github.com/dadosjusbr/alba/storage"
)

var expectedCollector = storage.Collector{
	ID:                 "trt13",
	Entity:             "Tribunal Regional do Trabalho 13ª Região",
	City:               "João Pessoa",
	FU:                 "PB",
	Path:               "github.com/dadosjusbr/coletores/trt13",
	Frequency:          30,
	StartDay:           5,
	LimitMonthBackward: 2,
	LimitYearBackward:  2018,
}

// Casos de teste:
// - Verificar se um input completo foi persisitido no banco
// - Verificar se um input incompleto foi persistido no banco
// - Verificar o comando add fromFile
// - Verificar o comando add
func TestAddCollector(t *testing.T) {

}

func TestAddCollectorFromFile(t *testing.T) {

}
