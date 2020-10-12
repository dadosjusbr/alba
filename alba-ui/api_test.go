package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dadosjusbr/alba/storage"
	"github.com/dadosjusbr/executor"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var pipelinesDB = []storage.Pipeline{
	{
		Pipeline: executor.Pipeline{
			Name: "Tribunal Regional do Trabalho 13ª Região",
		},
		ID:                 "trt13",
		Entity:             "Tribunal Regional do Trabalho 13ª Região",
		City:               "João Pessoa",
		FU:                 "PB",
		Repo:               "github.com/dadosjusbr/coletores",
		Frequency:          30,
		StartDay:           5,
		LimitMonthBackward: 2,
		LimitYearBackward:  2018,
	},
}

var executionsDB = []storage.Execution{
	{
		PipelineResult: executor.PipelineResult{
			Name:   "Tribunal Regional do Trabalho 13ª Região",
			Status: "OK",
		},
		ID:     "trt13",
		Entity: "Tribunal Regional do Trabalho 13ª Região",
	},
	{
		PipelineResult: executor.PipelineResult{
			Name:   "Ministério Público da Paraíba",
			Status: "OK",
		},
		ID:     "mppb",
		Entity: "Ministério Público da Paraíba",
	},
}

type fakeFinder struct {
	pipelines  []storage.Pipeline
	executions []storage.Execution
	err        error
}

func (fake fakeFinder) GetPipelines() ([]storage.Pipeline, error) {
	return fake.pipelines, fake.err
}
func (fake fakeFinder) GetExecutions() ([]storage.Execution, error) {
	return fake.executions, fake.err
}
func (fake fakeFinder) GetPipeline(id string) (storage.Pipeline, error) {
	return fake.pipelines[0], fake.err
}
func (fake fakeFinder) GetExecutionsByID(id string) ([]storage.Execution, error) {
	return fake.executions, fake.err
}

func TestGetPipelines_NotFoundURL(t *testing.T) {
	app := echo.New()
	app.GET(apiPipelines, func(c echo.Context) error {
		return getPipelines(fakeFinder{pipelines: pipelinesDB, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get("/alba/api/v1/non-existent").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetPipelines_Sucess(t *testing.T) {
	app := echo.New()
	app.GET(apiPipelines, func(c echo.Context) error {
		return getPipelines(fakeFinder{pipelines: pipelinesDB, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiPipelines).
		Expect(t).
		Assert(jsonpath.Contains(`$[0].pipeline.Name`, "Tribunal Regional do Trabalho 13ª Região")).
		Status(http.StatusOK).
		End()
}

func TestGetPipelines_EmptyDB(t *testing.T) {
	app := echo.New()
	app.GET(apiPipelines, func(c echo.Context) error {
		return getPipelines(fakeFinder{pipelines: nil, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiPipelines).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetPipelines_InternalServerError(t *testing.T) {
	app := echo.New()
	app.GET(apiPipelines, func(c echo.Context) error {
		return getPipelines(fakeFinder{pipelines: pipelinesDB, err: errors.New("get collector: internal server error")}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiPipelines).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}

func TestGetPipelineByID_Sucess(t *testing.T) {
	app := echo.New()
	app.GET(apiPipelineByID, func(c echo.Context) error {
		return getPipelineByID(fakeFinder{pipelines: pipelinesDB, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiPipelineByID).
		Query("idi", "trt13").
		Expect(t).
		Assert(jsonpath.Contains(`$.pipeline.Name`, "Tribunal Regional do Trabalho 13ª Região")).
		Status(http.StatusOK).
		End()
}

func TestGetPipelineByID_InternalServerError(t *testing.T) {
	app := echo.New()
	app.GET(apiPipelineByID, func(c echo.Context) error {
		return getPipelineByID(fakeFinder{pipelines: pipelinesDB, err: errors.New("get collector: internal server error")}, c)
	})

	apitest.New().
		Handler(app).
		Get("/alba/api/v1/pipeline/trt13").
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}

func TestGetExecutions_Sucess(t *testing.T) {
	app := echo.New()
	app.GET(apiRuns, func(c echo.Context) error {
		return getExecutions(fakeFinder{executions: executionsDB, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiRuns).
		Expect(t).
		Assert(jsonpath.Contains(`$[0].entity`, "Tribunal Regional do Trabalho 13ª Região")).
		Status(http.StatusOK).
		End()
}

func TestGetExecutions_EmptyDB(t *testing.T) {
	app := echo.New()
	app.GET(apiRuns, func(c echo.Context) error {
		return getExecutions(fakeFinder{executions: nil, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiRuns).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetExecutions_InternalServerError(t *testing.T) {
	app := echo.New()
	app.GET(apiRuns, func(c echo.Context) error {
		return getPipelines(fakeFinder{executions: executionsDB, err: errors.New("get collector: internal server error")}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiRuns).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}

func TestGetExecutionsByID_Sucess(t *testing.T) {
	app := echo.New()
	app.GET(apiRunsByID, func(c echo.Context) error {
		return getExecutionsByID(fakeFinder{executions: executionsDB, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiRunsByID).
		Expect(t).
		Assert(jsonpath.Contains(`$[1].entity`, "Ministério Público da Paraíba")).
		Status(http.StatusOK).
		End()
}

func TestGetExecutionsByID_EmptyDB(t *testing.T) {
	app := echo.New()
	app.GET(apiRunsByID, func(c echo.Context) error {
		return getExecutionsByID(fakeFinder{executions: nil, err: nil}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiRunsByID).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetExecutionsByID_InternalServerError(t *testing.T) {
	app := echo.New()
	app.GET(apiRunsByID, func(c echo.Context) error {
		return getPipelines(fakeFinder{executions: executionsDB, err: errors.New("get collector: internal server error")}, c)
	})

	apitest.New().
		Handler(app).
		Get(apiRunsByID).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}
