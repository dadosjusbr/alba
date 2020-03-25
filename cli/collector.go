package main

import (
	"time"

	"github.com/dadosjusbr/alba/storage"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

var flagsAddCollector = []cli.Flag{
	&cli.StringFlag{Name: "id", Usage: "Initials entity like 'trt13'"},
	&cli.StringFlag{Name: "entity", Usage: "Entity from which the collector extracts data like 'Tribunal Regional do Trabalho 13° Região'"},
	&cli.StringFlag{Name: "city", Usage: "City of the entity from which the collector extracts data"},
	&cli.StringFlag{Name: "fu", Usage: "Federation unit of the entity from which the collector extracts data"},
	&cli.StringFlag{Name: "path", Usage: "Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'"},
	&cli.IntFlag{Name: "frequency", Usage: "Frequency of the collector execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'"},
	&cli.IntFlag{Name: "startDay", Usage: "Day of the month for the collector execution. Values must be between 1 and 30"},
	&cli.IntFlag{Name: "limitMonthBackward", Usage: "The limit month to which the collector must be executed in its historical execution"},
	&cli.IntFlag{Name: "limitYearBackward", Usage: "The limit year until which the collector must be executed in its historical execution"},
}

var flagsAddCollectorFromFile = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{Name: "id", Usage: "Initials entity like 'trt13'"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "entity", Usage: "Entity from which the collector extracts data like 'Tribunal Regional do Trabalho 13° Região'"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "city", Usage: "City of the entity from which the collector extracts data"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "fu", Usage: "Federation unit of the entity from which the collector extracts data"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "path", Usage: "Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "frequency", Usage: "Frequency of the collector execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "startDay", Usage: "Day of the month for the collector execution. Values must be between 1 and 30"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "limitMonthBackward", Usage: "The limit month to which the collector must be executed in its historical execution"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "limitYearBackward", Usage: "The limit year until which the collector must be executed in its historical execution"}),
	&cli.StringFlag{Name: "file", Usage: "File with collector data", Required: true},
}

// GetCommandAddCollector return the stuct for add-collector command.
func GetCommandAddCollector() []*cli.Command {
	command := []*cli.Command{
		{
			Name:   "add-collector",
			Usage:  "Register a collector from parameters",
			Action: addCollector,
			Flags:  flagsAddCollector,
			Subcommands: []*cli.Command{
				{
					Name:   "from-file",
					Usage:  "Register a collector from a JSON file",
					Action: addCollector,
					Before: altsrc.InitInputSourceWithContext(flagsAddCollectorFromFile, altsrc.NewJSONSourceFromFlagFunc("file")),
					Flags:  flagsAddCollectorFromFile,
				},
			},
		},
	}

	return command

}

func addCollector(c *cli.Context) error {
	newCollector, err := parseCollectorFromContext(c)
	if err != nil {
		return err
	}
	if err := storage.InsertCollector(newCollector); err != nil {
		return err
	}
	return nil
}

func parseCollectorFromContext(c *cli.Context) (storage.Collector, error) {
	id := c.String("id")
	entity := c.String("entity")
	city := c.String("city")
	fu := c.String("fu")
	path := c.String("path")
	frequency := c.Int("frequency")
	startDay := c.Int("startDay")
	limitMonthBackward := c.Int("limitMonthBackward")
	limitYearBackward := c.Int("limitYearBackward")
	updateDate := time.Now()

	if id == "" {
		return storage.Collector{}, cli.Exit("--id were not provided completely. Please provide all parameters to continue", 1)
	}
	if entity == "" {
		return storage.Collector{}, cli.Exit("--entity were not provided completely. Please provide all parameters to continue", 1)

	}
	if city == "" {
		return storage.Collector{}, cli.Exit("--city were not provided completely. Please provide all parameters to continue", 1)

	}
	if fu == "" {
		return storage.Collector{}, cli.Exit("--fu were not provided completely. Please provide all parameters to continue", 1)

	}
	if path == "" {
		return storage.Collector{}, cli.Exit("--path were not provided completely. Please provide all parameters to continue", 1)

	}
	if frequency == 0 {
		return storage.Collector{}, cli.Exit("--frequency were not provided completely. Please provide all parameters to continue", 1)

	}
	if startDay == 0 {
		return storage.Collector{}, cli.Exit("--startDay were not provided completely. Please provide all parameters to continue", 1)

	}
	if limitMonthBackward == 0 {
		return storage.Collector{}, cli.Exit("--limitMonthBackward were not provided completely. Please provide all parameters to continue", 1)

	}
	if limitYearBackward == 0 {
		return storage.Collector{}, cli.Exit("--limitYearBackward were not provided completely. Please provide all parameters to continue", 1)
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

	return newCollector, nil
}
