package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dadosjusbr/alba/server/model"
)

func readCollectorData(name, entity, city, uf string) model.Collector {

	updateDate := time.Now()

	fmt.Print("1/6 - Enter collector path: ")
	var path string
	fmt.Scanf("%s", &path)

	fmt.Print("2/6 - Enter collector id_version: ")
	var idVersion string
	fmt.Scanf("%s", &idVersion)

	fmt.Print("3/6 - Enter collector frequency: ")
	var frequency int
	fmt.Scanf("%d", &frequency)

	fmt.Print("4/6 - Enter start day of collector : ")
	var startDay int
	fmt.Scanf("%d", &startDay)

	fmt.Print("5/6 - Enter collector limit month backward: ")
	var limitMonthBackward int
	fmt.Scanf("%d", &limitMonthBackward)

	fmt.Print("6/6 - Enter collector limit year backward: ")
	var limitYearBackward int
	fmt.Scanf("%d", &limitYearBackward)

	return model.Collector{name, entity, city, uf, updateDate,
		path, idVersion, frequency, startDay, limitMonthBackward,
		limitYearBackward}
}

func main() {

	name := flag.String("name", "", "Name of collector")
	entity := flag.String("entity", "", "Name of the entity from which the collector extracts data")
	city := flag.String("city", "", "Name of the entity city")
	uf := flag.String("uf", "", "UF of entity city")

	flag.Parse()

	if *name == "" || *entity == "" || *city == "" || *uf == "" {
		log.Fatal("Name, Entity, City or UF not provided. Please provide those to continue. --name={} --entity={} --city={} --uf={}\n")
		os.Exit(1)
	}

	newCollector := readCollectorData(*name, *entity, *city, *uf)

	fmt.Println("New Collector: ", newCollector)
	model.InsertCollector(newCollector)
}
