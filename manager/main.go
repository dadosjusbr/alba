package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dadosjusbr/alba/manager/model"
)

func readCollectorData(ID, entity, city, FU string) model.Collector {

	updateDate := time.Now()

	fmt.Print("1/6 - Enter collector path: ")
	var path string
	fmt.Scanf("%s", &path)

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

	return model.Collector{ID, entity, city, FU, updateDate,
		path, frequency, startDay, limitMonthBackward,
		limitYearBackward}
}

func main() {

	id := flag.String("id", "", "Initials entity")
	entity := flag.String("entity", "", "Name of the entity from which the collector extracts data")
	city := flag.String("city", "", "Name of the entity city")
	fu := flag.String("fu", "", "FU of entity city")

	flag.Parse()

	if *id == "" || *entity == "" || *city == "" || *fu == "" {
		log.Fatal("ID, Entity, City or FU not provided. Please provide those to continue. --id={} --entity={} --city={} --fu={}\n")
		os.Exit(1)
	}

	newCollector := readCollectorData(*id, *entity, *city, *fu)

	fmt.Println("New Collector: ", newCollector)
	model.InsertCollector(newCollector)
}
