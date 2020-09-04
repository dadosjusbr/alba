package pipeline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dadosjusbr/alba/storage"

	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

type inserter interface {
	InsertPipeline(storage.Pipeline) error
}

type addCommand struct {
	inserter inserter
}

func (a addCommand) do(c *cli.Context) error {
	pipelines, err := fromFile(c.String("from-file"))
	if err != nil {
		return fmt.Errorf("error adding pipeline: %q", err)
	}
	for _, pip := range pipelines {
		err := validate(pip)
		if err != nil {
			return fmt.Errorf("error adding pipeline. invalid pipeline descriptor: %q", err)
		}
		pip.Pipeline.DefaultBaseDir, err = cloneRepo(pip.Repo)
		if err != nil {
			return fmt.Errorf("error adding pipeline: %q", err)
		}
		pip.UpdateDate = time.Now()
		if err := a.inserter.InsertPipeline(pip); err != nil {
			return fmt.Errorf("error adding pipeline. error updating database: %q", err)
		}
		fmt.Printf("Pipeline ID: %s, Repo: %s\n", pip.ID, pip.Repo)
	}

	return nil
}

// NewAddCommand creates a new command to add a pipeline to the database.
func NewAddCommand(i inserter) *cli.Command {
	a := addCommand{inserter: i}
	return &cli.Command{Name: "add",
		Usage:  "Register one or more pipelines from a json file.",
		Action: a.do,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "from-file", Usage: "File path containing the spec of the pipeline to be added.", Required: true},
		}}
}

func fromFile(path string) ([]storage.Pipeline, error) {
	f, err := os.Open(path)
	if err != nil {
		return []storage.Pipeline{}, fmt.Errorf("error opening file [%s]:{%q}", path, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return []storage.Pipeline{}, fmt.Errorf("{error reading file [%s]:{%q}", path, err)
	}
	var pipelines []storage.Pipeline
	if err := json.Unmarshal(b, &pipelines); err != nil {
		return []storage.Pipeline{}, fmt.Errorf("error parsing pipeline descriptor [path:%s \n desc:'%s']:{%q}", path, string(b), err)
	}
	return pipelines, nil
}

func validate(p storage.Pipeline) error {
	if p.ID == "" {
		return fmt.Errorf("id were not provided. Please provide all mandatory parameters to continue")
	}
	if p.Entity == "" {
		return fmt.Errorf("entity were not provided. Please provide all mandatory parameters to continue")
	}
	if p.City == "" {
		return fmt.Errorf("city were not provided. Please provide all mandatory parameters to continue")
	}
	if p.FU == "" {
		return fmt.Errorf("fu were not provided. Please provide all mandatory parameters to continue")
	}
	if p.Repo == "" {
		return fmt.Errorf("repo were not provided. Please provide all mandatory parameters to continue")
	}
	if p.Frequency == 0 {
		return fmt.Errorf("frequency were not provided. Please provide all mandatory parameters to continue")
	}
	if p.StartDay == 0 {
		return fmt.Errorf("start-day were not provided. Please provide all mandatory parameters to continue")
	}
	if p.LimitMonthBackward == 0 {
		return fmt.Errorf("limit-month-backward were not provided. Please provide all mandatory parameters to continue")
	}
	if p.LimitYearBackward == 0 {
		return fmt.Errorf("limit-year-backward were not provided. Please provide all mandatory parameters to continue")
	}

	return nil
}

func cloneRepo(repo string) (string, error) {
	reposDir := os.Getenv("REPOSDIR")
	if reposDir == "" {
		return "", fmt.Errorf("error cloning the repository. REPOSDIR env var can not be empty")
	}
	reposDir = fmt.Sprintf("%s/%s", reposDir, repo)
	url := fmt.Sprintf("https://%s", repo)
	fmt.Println(url)
	_, err := git.PlainClone(reposDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", fmt.Errorf("error cloning the repository: %q", err)
	}

	return reposDir, nil
}
