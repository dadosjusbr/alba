package pipeline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dadosjusbr/alba/storage"
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

		url := fmt.Sprintf("https://%s", pip.Repo)
		resp, err := http.Head(url)
		if err != nil {
			return fmt.Errorf("error adding pipeline. error http.Head(): %q", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("error adding pipeline. error reaching repo url [%s]: %q", url, resp.Status)
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
func NewAddCommand() *cli.Command {
	var client *storage.DBClient
	c := newAddCommand(client)
	c.Before = func(c *cli.Context) error {
		uri := os.Getenv("MONGODB")
		if uri == "" {
			log.Fatal("[add command] error trying get environment variable: $MONGODB is empty")
		}
		var err error
		client, err = storage.NewDBClient(uri)
		if err != nil {
			log.Fatal(err)
		}
		return client.Connect()
	}
	c.After = func(c *cli.Context) error {
		return client.Disconnect()
	}
	return c
}

func newAddCommand(i inserter) *cli.Command {
	a := addCommand{i}
	return &cli.Command{Name: "add",
		Usage: "Register one or more pipelines from a json file.",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "from-file", Usage: "File path containing the spec of the pipeline to be added.", Required: true},
		},
		Action: a.do,
	}
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
