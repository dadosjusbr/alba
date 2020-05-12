package main

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

//View represents the functions for build and return html pages.
type view struct {
	client *storage.DBClient
}

//e.GET("/alba", index)
func (v view) index(c echo.Context) error {
	results, err := v.client.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	//TODO: Retornar página hmtl indicando que não existem resultados
	if len(results) == 0 {
		c.Render(http.StatusOK, "home.html", "")
	}

	data := struct {
		Collectors []storage.Collector
	}{
		results,
	}

	return c.Render(http.StatusOK, "home.html", data)
}

//e.GET("/alba/:id", viewExecutionsByID)
func (v view) executionsByID(c echo.Context) error {
	//mockup
	data := ExecutionDetails{
		Entity: "Nome do órgão",
		Executions: []Execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.Render(http.StatusOK, "executionsDetails.html", data)
}
