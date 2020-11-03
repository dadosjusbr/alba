package pipeline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dadosjusbr/alba/storage"
	"github.com/urfave/cli/v2"
)

type adhocStdinManager struct {
}

func (m adhocStdinManager) GetPipeline(string) (storage.Pipeline, error) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return storage.Pipeline{}, fmt.Errorf("error reading pipeline description from stdin:%q", err)
	}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var p storage.Pipeline
	if err := json.Unmarshal(b, &p); err != nil {
		return storage.Pipeline{}, fmt.Errorf("error unmarshaling pipeline description from stdin. Err:\"%q\"\nDescription:\"%s\"", err, string(b))
	}

	h, _ := json.MarshalIndent(p, "", "\t")
	fmt.Println(string(h))

	p.Pipeline.DefaultBaseDir = fmt.Sprintf("%s/%s", dir, p.Repo)
	fmt.Println(p.Pipeline.DefaultBaseDir)

	return p, nil
}

func (m adhocStdinManager) InsertExecution(storage.Execution) error {
	return nil // Not implemented in the adhoc execution.
}

// NewRunAdhocCommand creates a new command to run a pipeline from the description passed-in via standard input.
func NewRunAdhocCommand() *cli.Command {
	e := runCommand{manager: &adhocStdinManager{}}
	return &cli.Command{Name: "run-adhoc",
		Usage:  "Run a pipeline registered described from the standard input",
		Action: e.do,
	}
}
