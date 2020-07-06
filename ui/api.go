package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type execution struct {
	Date   string
	Status string
	Result string
}

type executionDetails struct {
	Entity     string
	Executions []execution
}

func getCollectors(f finder, c echo.Context) error {
	results, err := f.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if len(results) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, results)
}

func getExecutions(f finder, c echo.Context) error {
	// Mockup.
	data := executionDetails{
		Entity: "Nome do órgão",
		Executions: []execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.JSON(http.StatusOK, data)
}
