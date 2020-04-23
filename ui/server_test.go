package main

import (
	"testing"
	"time"

	"github.com/dadosjusbr/alba/storage"
)

var resultDB storage.Collector = storage.Collector{
	ID:                 "trt13",
	Entity:             "Tribunal Regional do Trabalho 13ª Região",
	City:               "João Pessoa",
	FU:                 "PB",
	UpdateDate:         time.Now(),
	Path:               "github.com/dadosjusbr/coletores/trt13",
	Frequency:          30,
	StartDay:           5,
	LimitMonthBackward: 2,
	LimitYearBackward:  2018,
}

func TestViewAllCollectors(t *testing.T) {

}

func testViewCollectorByID(t *testing.T) {

}

func testViewCollectorByPath(t *testing.T) {

}
