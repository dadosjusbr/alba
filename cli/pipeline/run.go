package pipeline

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/dadosjusbr/alba/git"
	"github.com/dadosjusbr/alba/storage"
)

type managerExecution interface {
	GetPipeline(string) (storage.Pipeline, error)
	InsertExecution(storage.Execution) error
}

type runCommand struct {
	manager managerExecution
}

// NewRunCommand creates a new command to run a pipeline.
func NewRunCommand() *cli.Command {
	uri := os.Getenv("MONGODB")
	if uri == "" {
		log.Fatal("error trying get environment variable: $MONGODB is empty")
	}
	client, err := storage.NewDBClient(uri)
	if err != nil {
		log.Fatal(err)
	}
	e := runCommand{manager: client}
	return &cli.Command{Name: "run",
		Usage:  "Run a pipeline registered in the database.",
		Action: e.do,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Usage: "Pipeline ID.", Required: true},
		},
		Before: func(c *cli.Context) error {
			return client.Connect()
		},
		After: func(c *cli.Context) error {
			return client.Disconnect()
		},
	}
}

func (r runCommand) do(c *cli.Context) error {
	id := c.String("id")
	p, err := r.manager.GetPipeline(id)
	if err != nil {
		return fmt.Errorf("error running pipeline. error getting pipeline from database: %q", err)
	}
	if len(p.Repo) == 0 {
		return fmt.Errorf("error running pipeline. there is no pipeline registered for id: %s", id)
	}

	baseDir := os.Getenv("BASEDIR")
	if baseDir == "" {
		return fmt.Errorf("error running pipeline. BASEDIR env var can not be empty")
	}
	month := os.Getenv("MONTH")
	if month == "" {
		return fmt.Errorf("error running pipeline. MONTH env var can not be empty")
	}
	year := os.Getenv("YEAR")
	if year == "" {
		return fmt.Errorf("error running pipeline. YEAR env var can not be empty")
	}
	var commit string
	p.Pipeline.DefaultBaseDir = fmt.Sprintf("%s/%s", baseDir, p.Repo)
	commit, err = git.CloneRepository(p.Pipeline.DefaultBaseDir, fmt.Sprintf("https://%s", p.Repo))
	if err != nil {
		return fmt.Errorf("error running pipeline: %q", err)
	}

	defaultBuildEnv := map[string]string{
		"GIT_COMMIT": commit,
	}

	defaultRunEnv := map[string]string{
		"OUTPUT_FOLDER": "/output",
		"MONTH":         month,
		"YEAR":          year,
	}

	for pos, stage := range p.Pipeline.Stages {
		stage.BuildEnv = mergeEnv(defaultBuildEnv, stage.BuildEnv)
		stage.RunEnv = mergeEnv(defaultRunEnv, stage.RunEnv)
		p.Pipeline.Stages[pos] = stage
	}

	result, _ := p.Pipeline.Run()
	e := storage.Execution{
		PipelineResult: result,
		Entity:         p.Entity,
		ID:             p.ID,
	}
	r.manager.InsertExecution(e)
	return nil
}

func mergeEnv(defaultEnv, stageEnv map[string]string) map[string]string {
	env := make(map[string]string)

	for k, v := range defaultEnv {
		env[k] = v
	}
	for k, v := range stageEnv {
		env[k] = v
	}
	return env
}
