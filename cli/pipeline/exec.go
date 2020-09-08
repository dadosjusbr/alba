package pipeline

import (
	"github.com/dadosjusbr/executor"
	"github.com/urfave/cli/v2"
)

type managerExecution interface {
	GetPipeline(string) error
	InsertExecution(executor.PipelineResult) error
}

type execCommand struct {
	manager managerExecution
}

// NewExecCommand creates a new command to add a pipeline to the database.
func NewExecCommand(m managerExecution) *cli.Command {
	e := execCommand{manager: m}
	return &cli.Command{Name: "exec",
		Usage:  "Run a pipeline registered in the database.",
		Action: e.do,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Usage: "Pipeline ID.", Required: true},
			&cli.StringFlag{Name: "month", Usage: "Month to be collected.", Required: true},
			&cli.StringFlag{Name: "year", Usage: "Year to be collected.", Required: true},
		}}
}

func (e execCommand) do(c *cli.Context) error {
	// parameters validation
	// config pipeline from database

	return nil
}
