package main

import (
	"testing"

	"github.com/dadosjusbr/alba/storage"
	"github.com/dadosjusbr/executor"
	"github.com/matryer/is"
)

var pipelinesDB = []storage.Pipeline{
	{
		Pipeline: executor.Pipeline{
			Name: "Pipeline 01: Finalizado com sucesso",
		},
		ID:                 "pipeline1",
		Entity:             "Pipeline 01",
		City:               "João Pessoa",
		FU:                 "PB",
		Repo:               "github.com/dadosjusbr/coletores",
		Frequency:          30,
		StartDay:           5,
		LimitMonthBackward: 2,
		LimitYearBackward:  2021,
	},
	{
		Pipeline: executor.Pipeline{
			Name: "Pipeline 02: Tolerância a falha",
		},
		ID:                 "pipeline2",
		Entity:             "Tribunal Regional do Trabalho 13ª Região",
		City:               "João Pessoa",
		FU:                 "PB",
		Repo:               "github.com/dadosjusbr/coletores",
		Frequency:          30,
		StartDay:           5,
		LimitMonthBackward: 2,
		LimitYearBackward:  2021,
	},
	{
		Pipeline: executor.Pipeline{
			Name: "Pipeline 03: Histórico",
		},
		ID:                 "pipeline3",
		Entity:             "Pipeline 03",
		City:               "João Pessoa",
		FU:                 "PB",
		Repo:               "github.com/dadosjusbr/coletores",
		Frequency:          30,
		StartDay:           5,
		LimitMonthBackward: 2,
		LimitYearBackward:  2021,
	},
}

type fakeFinder struct {
	pipelines  []storage.Pipeline
	executions []storage.Execution
	err        error
}

func (fake fakeFinder) GetPipelinesByDay(day int) ([]storage.Pipeline, error) {
	return fake.pipelines, fake.err
}
func (fake fakeFinder) InsertExecution(e storage.Execution) error {
	return fake.err
}

// TODO: Retornar as últimas execuções de cada pipeline
func (fake fakeFinder) GetLastExecutionsForAllPipelines() ([]storage.Execution, error) {
	return fake.executions, fake.err
}

// RODO: Retornar as <limit> ultimas execuções daquele pipeline
func (fake fakeFinder) GetLastExecutionsByPipelineID(limit, id int) ([]storage.Execution, error) {
	return fake.executions, fake.err
}

// PrioritizeAndLimit
// 1º descarta se o pipeline já tem uma execução finalizada com sucesso hoje
// 2º se o número se execuções com erro for igual ou maior que 3
// por ultimo aplicar filtro de tamanho
func TestPrioritizeAndLimit_LimitTest(t *testing.T) {
	is := is.New(t)
	list := pipelinesDB
	limit := 2

	finalPipelines := prioritizeAndLimit(fakeFinder{pipelines: pipelinesDB, err: nil}, list, limit)

	is.True(len(finalPipelines) <= limit)
}

func TestPrioritizeAndLimit_PipelineWithMaxLimitOfFailure(t *testing.T) {
	is := is.New(t)
	list := pipelinesDB
	limit := 3

	executions := []storage.Execution{
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 02",
				Status: "Setup Error",
			},
			ID:     "pipeline2",
			Entity: "Pipeline 02",
		},
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 02",
				Status: "Connection Error",
			},
			ID:     "pipeline2",
			Entity: "Pipeline 02",
		},
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 02",
				Status: "Run Error",
			},
			ID:     "pipeline2",
			Entity: "Pipeline 02",
		},
	}

	finalPipelines := prioritizeAndLimit(fakeFinder{executions: executions, err: nil}, list, limit)

	is.True(len(finalPipelines) == 2)
}

func TestPrioritizeAndLimit_PipelineWithTwoFailures(t *testing.T) {
	is := is.New(t)
	list := pipelinesDB
	limit := 3

	executions := []storage.Execution{
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 02",
				Status: "Setup Error",
			},
			ID:     "pipeline2",
			Entity: "Pipeline 02",
		},
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 02",
				Status: "Connection Error",
			},
			ID:     "pipeline2",
			Entity: "Pipeline 02",
		},
	}

	finalPipelines := prioritizeAndLimit(fakeFinder{executions: executions, err: nil}, list, limit)

	is.True(len(finalPipelines) == 3)
}

func TestPrioritizeAndLimit_PipelineWithSuccessfulExecution(t *testing.T) {
	is := is.New(t)
	list := pipelinesDB
	limit := 3

	executions := []storage.Execution{
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 01",
				Status: "OK",
			},
			ID:     "pipeline1",
			Entity: "Pipeline 01",
		},
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 02",
				Status: "Setup Error",
			},
			ID:     "pipeline2",
			Entity: "Pipeline 02",
		},
		{
			PipelineResult: executor.PipelineResult{
				Name:   "Pipeline 03",
				Status: "OK",
			},
			ID:     "pipeline3",
			Entity: "Pipeline 03",
		},
	}

	finalPipelines := prioritizeAndLimit(fakeFinder{executions: executions, err: nil}, list, limit)

	is.True(len(finalPipelines) == 1)
}
