package main

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

type collectorsGetter interface {
	GetCollectors() ([]storage.Collector, error)
}

//API represents the functions for build and return html pages.
type api struct {
	client *storage.DBClient
	getter collectorsGetter
}

//Execution represents a execution
type Execution struct {
	Date   string
	Status string
	Result string
}

//ExecutionDetails represents the details of the execution from collector.
type ExecutionDetails struct {
	Entity     string
	Executions []Execution
}

func (a *api) getCollectors(c echo.Context) error {
	results, err := a.getter.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if len(results) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, results)
}

//ExecutionsByID returns all executions by collector ID.
func (a *api) executionsByID(c echo.Context) error {
	//Mockup
	data := ExecutionDetails{
		Entity: "Nome do órgão",
		Executions: []Execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.JSON(http.StatusOK, data)
}
