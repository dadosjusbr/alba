package api

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

//CollectorsURL represents the collectors url.
const CollectorsURL = "/alba/api/collectors"

//RunsURL represents the executions url.
const RunsURL = "/alba/api/runs/:id"

type collectorsGetter interface {
	GetCollectors() ([]storage.Collector, error)
}

type prodCollectorsGetter struct {
}

func (prod prodCollectorsGetter) GetCollectors() ([]storage.Collector, error) {
	return storage.GetCollectors()
}

//GetCollectorsHandler set prodGetCollector.
func GetCollectorsHandler(c echo.Context) error {
	return getCollectors(c, prodCollectorsGetter{})
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

func getCollectors(c echo.Context, getter collectorsGetter) error {
	results, err := getter.GetCollectors()
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
func ExecutionsByID(c echo.Context) error {
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
