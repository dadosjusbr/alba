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

func getPipelines(f finder, c echo.Context) error {
	results, err := f.GetPipelines()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if len(results) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, results)
}

func getPipelineByID(f finder, c echo.Context) error {
	id := c.Param("id")
	result, err := f.GetPipeline(id)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if result.ID == "" {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, result)
}

func getExecutions(f finder, c echo.Context) error {
	results, err := f.GetExecutions()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	if len(results) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, results)
}

func getExecutionsByID(f finder, c echo.Context) error {
	id := c.Param("id")
	results, err := f.GetExecutionsByID(id)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	if len(results) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, results)
}
