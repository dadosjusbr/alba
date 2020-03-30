package collector

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

// AddCommand is the command which adds a collector to the database.
var AddCommand = &cli.Command{
	Name:   "add-collector",
	Usage:  "Register a collector from parameters",
	Action: add,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "Initials entity like 'trt13'"},
		&cli.StringFlag{Name: "entity", Usage: "Entity from which the collector extracts data like 'Tribunal Regional do Trabalho 13° Região'"},
		&cli.StringFlag{Name: "city", Usage: "City of the entity from which the collector extracts data"},
		&cli.StringFlag{Name: "fu", Usage: "Federation unit of the entity from which the collector extracts data"},
		&cli.StringFlag{Name: "path", Usage: "Collector repository path. Using the import pattern in golang like 'github.com/dadosjusbr/coletores/trt13'"},
		&cli.IntFlag{Name: "frequency", Usage: "Frequency of the collector execution in days. Values must be between 1 and 30. To be executed monthly it must be filled with '30'"},
		&cli.IntFlag{Name: "startDay", Usage: "Day of the month for the collector execution. Values must be between 1 and 30"},
		&cli.IntFlag{Name: "limitMonthBackward", Usage: "The limit month to which the collector must be executed in its historical execution"},
		&cli.IntFlag{Name: "limitYearBackward", Usage: "The limit year until which the collector must be executed in its historical execution"},
		&cli.StringFlag{Name: fromFileParam, Usage: "File path containing the spec of the collection to be added."},
	},
}

// Add adds a collector to the database.
func add(c *cli.Context) error {
	var collector storage.Collector
	p := c.String(fromFileParam)
	if p != "" { // From file has priority over passing parameters.
		col, err := fromFile(p)
		if err != nil {
			return fmt.Errorf("error adding collector:{%q}", err)
		}
		collector = col
	} else {
		collector = fromContext(c)
	}
	if err := validate(collector); err != nil {
		return fmt.Errorf("invalid collector descriptor:{%q}", err)
	}
	if err := storage.InsertCollector(collector); err != nil {
		return fmt.Errorf("error updating database:{%q}", err)
	}
	return nil
}

func fromFile(path string) (storage.Collector, error) {
	f, err := os.Open(path)
	if err != nil {
		return storage.Collector{}, fmt.Errorf("error opening file [%s]:{%q}", path, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return storage.Collector{}, fmt.Errorf("{error reading file [%s]:{%q}", path, err)
	}
	var c storage.Collector
	if err := json.Unmarshal(b, &c); err != nil {
		return storage.Collector{}, fmt.Errorf("error parsing collector descriptor [path:%s \n desc:'%s']:{%q}", path, string(b), err)
	}
	return c, nil
}

func validate(col storage.Collector) error {
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
		return fmt.Errorf("--startDay were not provided completely. Please provide all parameters to continue")

	}
	if col.LimitMonthBackward == 0 {
		return fmt.Errorf("--limitMonthBackward were not provided completely. Please provide all parameters to continue")

	}
	if col.LimitYearBackward == 0 {
		return fmt.Errorf("--limitYearBackward were not provided completely. Please provide all parameters to continue")
	}
	return nil
}

func fromContext(c *cli.Context) storage.Collector {
	return storage.Collector{
		ID:                 c.String("id"),
		Entity:             c.String("entity"),
		City:               c.String("city"),
		FU:                 c.String("fu"),
		UpdateDate:         time.Now(),
		Path:               c.String("path"),
		Frequency:          c.Int("frequency"),
		StartDay:           c.Int("startDay"),
		LimitMonthBackward: c.Int("limitMonthBackward"),
		LimitYearBackward:  c.Int("limitYearBackward"),
	}
}
