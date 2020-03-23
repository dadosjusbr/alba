package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

func main() {
	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."

	flagsAddCollector := []cli.Flag{
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

	flagsAddCollectorFromFile := []cli.Flag{
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

	app.Commands = []*cli.Command{
		{
			Name:   "add-collector",
			Usage:  "Register a collector from parameters",
			Action: AddCollector,
			Flags:  flagsAddCollector,
			Subcommands: []*cli.Command{
				{
					Name:   "from-file",
					Usage:  "Register a collector from a JSON file",
					Action: AddCollector,
					Before: altsrc.InitInputSourceWithContext(flagsAddCollectorFromFile, altsrc.NewJSONSourceFromFlagFunc("file")),
					Flags:  flagsAddCollectorFromFile,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
