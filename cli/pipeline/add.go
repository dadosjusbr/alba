package pipeline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dadosjusbr/alba/storage"

	"github.com/urfave/cli/v2"
)

const fromFileParam = "from-file"

type inserter interface {
	InsertPipeline(storage.Pipeline) error
}

type addCommand struct {
	inserter inserter
}

func (a addCommand) do(c *cli.Context) error {
	var pipeline storage.Pipeline
	p := c.String(fromFileParam)
	if p != "" { // From file has priority over passing parameters.
		pip, err := fromFile(p)
		if err != nil {
			return fmt.Errorf("error adding pipeline:{%q}", err)
		}
		pipeline = pip
	} else {
		pipeline = fromContext(c)
	}
	err := validate(pipeline)
	if err != nil {
		return fmt.Errorf("invalid pipeline descriptor:{%q}", err)
	}

	pipeline.UpdateDate = time.Now()

	if err := a.inserter.InsertPipeline(pipeline); err != nil {
		return fmt.Errorf("error updating database:{%q}", err)
	}
	fmt.Printf("Collector ID: %s, Path: %s", pipeline.ID, pipeline.Path)
	return nil
}

// NewAddCommand creates a new command to add a pipeline to the database.
func NewAddCommand(i inserter) *cli.Command {
	a := addCommand{inserter: i}
	return &cli.Command{Name: "add",
		Usage:  "Register a pipeline from parameters",
		Action: a.do,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Usage: "Initials entity like 'trt13'"},
			&cli.StringFlag{Name: "entity", Usage: "Entity from which the pipeline extracts data like 'Tribunal Regional do Trabalho 13° Região'"},
			&cli.StringFlag{Name: "city", Usage: "City of the entity from which the pipeline extracts data"},
			&cli.StringFlag{Name: "fu", Usage: "Federation unit of the entity from which the pipeline extracts data"},
			&cli.StringFlag{Name: "path", Usage: "Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'"},
			&cli.IntFlag{Name: "frequency", Usage: "Frequency of the pipeline execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'"},
			&cli.IntFlag{Name: "start-day", Usage: "Day of the month for the pipeline execution. Values must be between 1 and 30"},
			&cli.IntFlag{Name: "limit-month-backward", Usage: "The limit month to which the pipeline must be executed in its historical execution"},
			&cli.IntFlag{Name: "limit-year-backward", Usage: "The limit year until which the pipeline must be executed in its historical execution"},
			&cli.StringFlag{Name: fromFileParam, Usage: "File path containing the spec of the collection to be added."},
		}}
}

func fromFile(path string) (storage.Pipeline, error) {
	f, err := os.Open(path)
	if err != nil {
		return storage.Pipeline{}, fmt.Errorf("error opening file [%s]:{%q}", path, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return storage.Pipeline{}, fmt.Errorf("{error reading file [%s]:{%q}", path, err)
	}
	var c storage.Pipeline
	if err := json.Unmarshal(b, &c); err != nil {
		return storage.Pipeline{}, fmt.Errorf("error parsing pipeline descriptor [path:%s \n desc:'%s']:{%q}", path, string(b), err)
	}
	return c, nil
}

func validate(col storage.Pipeline) error {
	if col.ID == "" {
		return fmt.Errorf("--id were not provided completely. Please provide all parameters to continue")
	}
	if col.Entity == "" {
		return fmt.Errorf("--entity were not provided completely. Please provide all parameters to continue")
	}
	if col.City == "" {
		return fmt.Errorf("--city were not provided completely. Please provide all parameters to continue")
	}
	if col.FU == "" {
		return fmt.Errorf("--fu were not provided completely. Please provide all parameters to continue")
	}
	if col.Path == "" {
		return fmt.Errorf("--path were not provided completely. Please provide all parameters to continue")
	}
	if col.Frequency == 0 {
		return fmt.Errorf("--frequency were not provided completely. Please provide all parameters to continue")
	}
	if col.StartDay == 0 {
		return fmt.Errorf("--start-day were not provided completely. Please provide all parameters to continue")
	}
	if col.LimitMonthBackward == 0 {
		return fmt.Errorf("--limit-month-backward were not provided completely. Please provide all parameters to continue")
	}
	if col.LimitYearBackward == 0 {
		return fmt.Errorf("--limit-year-backward were not provided completely. Please provide all parameters to continue")
	}

	return nil
}

func fromContext(c *cli.Context) storage.Pipeline {
	return storage.Pipeline{
		ID:                 c.String("id"),
		Entity:             c.String("entity"),
		City:               c.String("city"),
		FU:                 c.String("fu"),
		Path:               c.String("path"),
		Frequency:          c.Int("frequency"),
		StartDay:           c.Int("start-day"),
		LimitMonthBackward: c.Int("limit-month-backward"),
		LimitYearBackward:  c.Int("limit-year-backward"),
	}
}
