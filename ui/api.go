package main

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

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

//APIGetCollectors returns  all collectors. e.GET("/alba/api/coletores", APIGetCollectors)
func APIGetCollectors(c echo.Context) error {
	results, err := storage.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if len(results) == 0 {
		return c.JSON(http.StatusOK, NotFound{msgNotFound, docmentationURL})
	}

	return c.JSON(http.StatusOK, results)
}

//APIExecutionsByID returns all executions by collector ID. e.GET("/alba/api/coletores/execucoes/:id", viewAPIExecutionsByID)
func APIExecutionsByID(c echo.Context) error {
	id := c.Param("id")

	//Mockup
	data := ExecutionDetails{
		Entity: id,
		Executions: []Execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.JSON(http.StatusOK, data)
}
