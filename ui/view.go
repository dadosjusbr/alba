package main

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

func index(f finder, c echo.Context) error {
	results, err := f.GetCollectors()
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

func viewExecutions(f finder, c echo.Context) error {
	//mockup
	data := executionDetails{
		Entity: "Nome do órgão",
		Executions: []execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.Render(http.StatusOK, "executionsDetails.html", data)
}
