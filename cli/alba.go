package main

import (
	"log"
	"os"

	"github.com/dadosjusbr/alba/cli/collector"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."
	app.Commands = []*cli.Command{collector.AddCommand}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
