package main

import (
	"log"
	"os"
	"time"

	"github.com/dadosjusbr/alba/storage"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

func main() {
	app := cli.NewApp()
	app.Name = "Alba Executions Orchestration System"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."

	flagsAddCollector := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{Name: "id", Usage: "Initials entity like 'trt13'"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "entity", Usage: "Entity from which the collector extracts data like 'Tribunal Regional do Trabalho 13° Região'"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "city", Usage: "City of the entity from which the collector extracts data"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "fu", Usage: "Federation unit of the entity from which the collector extracts data"}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "path", Usage: "Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "frequency", Usage: "Frequency of the collector execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "startDay", Usage: "Day of the month for the collector execution. Values must be between 1 and 30"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "limitMonthBackward", Usage: "The limit month to which the collector must be executed in its historical execution"}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "limitYearBackward", Usage: "The limit year until which the collector must be executed in its historical execution"}),
		&cli.StringFlag{Name: "file"},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "addCollector",
			Usage: "Register a collector",
			Action: func(c *cli.Context) error {
				newCollector := parseCollectorFromContext(c)

				err := storage.InsertCollector(newCollector)
				if err != nil {
					return err
				}
				return nil
			},
			Before: altsrc.InitInputSourceWithContext(flagsAddCollector, altsrc.NewJSONSourceFromFlagFunc("file")),
			Flags:  flagsAddCollector,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func parseCollectorFromContext(c *cli.Context) storage.Collector {
	id := c.String("id")
	entity := c.String("entity")
	city := c.String("city")
	fu := c.String("fu")
	path := c.String("path")
	frequency := c.Int("frequency")
	startDay := c.Int("startDay")
	limitMonthBackward := c.Int("limitMonthBackward")
	limitYearBackward := c.Int("limitYearBackward")
	updateDate := getUpdateDate()

	if id == "" || entity == "" || city == "" || fu == "" || path == "" || frequency == 0 ||
		startDay == 0 || limitMonthBackward == 0 || limitYearBackward == 0 {
		log.Fatal("Parameters were not provided completely. Please provide those to continue")
		os.Exit(1)
	}

	newCollector := storage.Collector{
		ID:                 id,
		Entity:             entity,
		City:               city,
		FU:                 fu,
		UpdateDate:         updateDate,
		Path:               path,
		Frequency:          frequency,
		StartDay:           startDay,
		LimitMonthBackward: limitMonthBackward,
		LimitYearBackward:  limitYearBackward}

	return newCollector
}

func getUpdateDate() time.Time {
	return time.Now()
}
