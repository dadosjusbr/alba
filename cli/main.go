package main

import (
	"log"
	"os"

	"github.com/dadosjusbr/alba/cli/pipeline"
	"github.com/dadosjusbr/alba/storage"

	"github.com/urfave/cli/v2"
)

func main() {
	var client *storage.DBClient

	uri := os.Getenv("MONGODB")
	if uri == "" {
		log.Fatal("error trying get environment variable: $MONGODB is empty")
	}

	client, err := storage.NewDBClient(uri)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(); err != nil {
		log.Fatal(err)

	}
	defer client.Disconnect()

	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through configurable pipelines runs."
	app.Commands = []*cli.Command{pipeline.NewAddCommand(client), pipeline.NewExecCommand(client)}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
