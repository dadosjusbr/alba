package main

import (
	"log"
	"os"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

type application struct {
	app *echo.Echo
}

type finder interface {
	GetPipelines() ([]storage.Pipeline, error)
	GetExecutions() ([]storage.Execution, error)
	GetPipeline(id string) (storage.Pipeline, error)
	GetExecutionsByID(id string) ([]storage.Execution, error)
}

// URLs definition.
const (
	apiPipelines    = "/alba/api/v1/pipelines"
	apiPipelineByID = "/alba/api/v1/pipeline/:id"
	apiRuns         = "/alba/api/v1/runs"
	apiRunsByid     = "/alba/api/v1/runs/:id"
)

func newApp(dbClient *storage.DBClient) *application {
	app := echo.New()

	app.GET(apiPipelines, func(c echo.Context) error {
		return getPipelines(dbClient, c)
	})
	app.GET(apiPipelineByID, func(c echo.Context) error {
		return getPipelineByID(dbClient, c)
	})
	app.GET(apiRuns, func(c echo.Context) error {
		return getExecutions(dbClient, c)
	})
	app.GET(apiRunsByid, func(c echo.Context) error {
		return getExecutionsByID(dbClient, c)
	})

	return &application{
		app: app,
	}
}

func (a *application) start() {
	a.app.Logger.Fatal(a.app.Start(":8080"))
}

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

	newApp(client).start()
}
