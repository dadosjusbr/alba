package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

//Run execute the cli
func Run(args []string) {
	fmt.Println(args)

	app := cli.NewApp()
	app.Name = "Alba"
	app.Usage = "A tool for manage the process of continuous data release through steps such as: collection, validation, packaging and storage."
	app.Commands = GetCommandAddCollector()

	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run(os.Args)
}
