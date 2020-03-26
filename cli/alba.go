package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."
	app.Commands = GetCommandAddCollector()

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
