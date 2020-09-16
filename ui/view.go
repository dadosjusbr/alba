package main

import (
	"fmt"
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

func index(f finder, c echo.Context) error {
	results, err := f.GetPipelines()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	// TODO: Retornar página hmtl indicando que não existem resultados.
	if len(results) == 0 {
		c.Render(http.StatusOK, "home.html", "")
	}

	data := struct {
		Pipelines []storage.Pipeline
	}{
		results,
	}

	return c.Render(http.StatusOK, "home.html", data)
}

func viewExecutions(f finder, c echo.Context) error {
	id := c.Param("id")
	results, err := f.GetExecutionsByID(id)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	// TODO: Retornar página html indicando que não existem resultados.
	if len(results) == 0 {
		c.Render(http.StatusOK, "home.html", "")
	}

	data := struct {
		Executions []storage.Execution
	}{
		results,
	}

	fmt.Println(data)

	return c.Render(http.StatusOK, "executionsDetails.html", data)
}
