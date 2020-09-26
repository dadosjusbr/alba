package main

import (
	"log"
	"os"

	"alba/pipeline"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through configurable pipelines runs."
	app.Commands = []*cli.Command{pipeline.NewAddCommand(), pipeline.NewRunCommand()}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
