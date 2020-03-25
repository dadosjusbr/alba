package main

import (
	"os"
	"testing"
)

func TestAddCollector(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "add-collector")
	args = append(args, "-id=mppb")
	args = append(args, "-entity=\"Ministério Público da Paraíba\"")
	args = append(args, "-city=\"João Pessoa\"")
	args = append(args, "-fu=PB")
	args = append(args, "-path=\"github.com/dadosjusbr/coletores/mppb\"")
	args = append(args, "-frequency=30")
	args = append(args, "-startDay=5")
	args = append(args, "-limitMonthBackward=1")
	args = append(args, "-limitYearBackward=2018")

	run(args)
}

func TestAddCollectorFromFile(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "add-collector")
	args = append(args, "from-file")
	args = append(args, "-file=input.json")

	run(args)
}
