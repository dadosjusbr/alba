package pipeline

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/dadosjusbr/alba/storage"
	"github.com/dadosjusbr/alba/worker"
	"github.com/dadosjusbr/executor"
)

type managerExecution interface {
	GetPipeline(string) (storage.Pipeline, error)
	InsertExecution(executor.PipelineResult) error
}

type runCommand struct {
	manager managerExecution
}

// NewRunCommand creates a new command to run a pipeline.
func NewRunCommand(m managerExecution) *cli.Command {
	e := runCommand{manager: m}
	return &cli.Command{Name: "run",
		Usage:  "Run a pipeline registered in the database.",
		Action: e.do,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Usage: "Pipeline ID.", Required: true},
			&cli.BoolFlag{Name: "dadosjusbr", Usage: "Flag to identify if it is a pipeline dadosjusbr.", Required: true},
			&cli.StringFlag{Name: "month", Usage: "Month to be collected."},
			&cli.StringFlag{Name: "year", Usage: "Year to be collected."},
		}}
}

func (r runCommand) do(c *cli.Context) error {
	id := c.String("id")
	p, err := r.manager.GetPipeline(id)
	if err != nil {
		return fmt.Errorf("error running pipeline. error getting pipeline from database: %q", err)
	}
	if len(p.Repo) == 0 {
		return fmt.Errorf("there is no pipeline registered for id: %s", id)
	}

	if c.Bool("id") == true {
		month := c.String("month")
		year := c.String("year")
		// Todo: Função que reúne as regras de negócio para um pipeline DadosJusBR,
		// como configuração das variáveis commit, mes e ano
		p, err = worker.ConfigureDadosjusBR(p, month, year)
		if err != nil {
			return fmt.Errorf("error running pipeline. error setting dadosjusbr pipeline: %q", err)
		}
	}

	p.Pipeline.DefaultBaseDir, _, err = worker.CloneRepository(p.Repo)
	if err != nil {
		return fmt.Errorf("error running pipeline: %q", err)
	}

	result, _ := p.Pipeline.Run()
	r.manager.InsertExecution(result)

	return nil
}
