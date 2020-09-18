package main

import (
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

	if len(results) == 0 {
		return echo.ErrNotFound
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

	if len(results) == 0 {
		return echo.ErrNotFound
	}

	data := struct {
		Executions []storage.Execution
		Entity     string
	}{
		results,
		results[0].Entity,
	}

	return c.Render(http.StatusOK, "executionsDetails.html", data)
}
