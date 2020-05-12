package main

import (
	"log"
	"os"

	"github.com/dadosjusbr/alba/cli/collector"
	"github.com/dadosjusbr/alba/storage"

	"github.com/urfave/cli/v2"
)

func main() {
	var client *storage.DBClient

	uri := os.Getenv("MONGODB")
	if uri == "" {
		log.Fatal("error trying get environment variable: $MONGODB is empty")
	}

	client, err := storage.NewClientDB(uri)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(); err != nil {
		log.Fatal(err)

	}
	defer client.Disconnect()

	add := collector.Add{Inserter: client}

	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."
	app.Commands = []*cli.Command{add.AddCommand()}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
