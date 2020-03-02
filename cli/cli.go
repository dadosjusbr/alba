package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

func main() {
	app := cli.NewApp()
	app.Name = "Alba Executions Orchestration System"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."

	flags := []cli.Flag{
		&cli.StringFlag{
			Name: "id",
		},
		&cli.StringFlag{
			Name: "entity",
		},
		&cli.StringFlag{
			Name: "city",
		},
		&cli.StringFlag{
			Name: "fu",
		},
		&cli.StringFlag{
			Name: "path",
		},
		&cli.IntFlag{
			Name: "frequency",
		},
		&cli.IntFlag{
			Name: "limit_month_backward",
		},
		&cli.IntFlag{
			Name: "limit_year_backward",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "addCollector",
			Usage: "Register a collector",
			Action: func(c *cli.Context) error {
				//newCollector := storage.Collector{ID, entity, city, FU, updateDate,
				//	path, frequency, startDay, limitMonthBackward,
				//	limitYearBackward}

				//return storage.InsertCollector(newCollector)
				fmt.Println(flags)
				return nil
			},
			Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewJSONSource("{}")),
			Flags:  flags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getUpdateDate() time.Time {
	return time.Now()
}
