package api

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

type getterCollector interface {
	GetCollectors() ([]storage.Collector, error)
}

type prodGetCollector struct {
}

func (prod prodGetCollector) GetCollectors() ([]storage.Collector, error) {
	return storage.GetCollectors()
}

//AddGetCollector set prodGetCollector
func AddGetCollector(c echo.Context) error {
	return getCollector(c, prodGetCollector{})
}

//NotFound information for the user when their search has no results
type NotFound struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"docmentation_url"`
}

const msgNotFound = "Not found"
const docmentationURL = "https://github.com/dadosjusbr/alba/wiki/API"

//Execution represents a execution
type Execution struct {
	Date   string
	Status string
	Result string
}

//ExecutionDetails represents the details of the execution from collector
type ExecutionDetails struct {
	Entity     string
	Executions []Execution
}

func getCollector(c echo.Context, getter getterCollector) error {
	results, err := getter.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if len(results) == 0 {
		return c.JSON(http.StatusOK, NotFound{msgNotFound, docmentationURL})
	}

	return c.JSON(http.StatusOK, results)
}

//ExecutionByID returns all executions by collector ID. e.GET("/alba/api/coletores/execucoes/:id", ExecutionsByID)
func ExecutionByID(c echo.Context) error {
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
