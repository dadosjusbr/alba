package pipeline

import (
	"testing"

	"github.com/dadosjusbr/alba/storage"

	"github.com/matryer/is"
	"github.com/urfave/cli/v2"
)

type fakeInserter struct {
}

func (fakeInserter) InsertPipeline(storage.Pipeline) error {
	return nil
}

func ExampleNewAddCommand() {
	app := newAddApp()
	app.Run([]string{
		"alba",
		"add",
		"--from-file=pipeline-example.json",
	})
	// Output:
	//Pipeline ID: stagego, Repo: github.com/dadosjusbr/executor
}

func TestAdd_RequiredParam(t *testing.T) {
	is := is.New(t)
	app := newAddApp()
	args := append([]string{"alba", "add", ""})
	is.True(app.Run(args) != nil)
}

func TestAdd_NoexistentURL(t *testing.T) {
	is := is.New(t)
	app := newAddApp()
	args := append([]string{"alba", "add", "--from-file=pipeline-example-error.json"})
	is.True(app.Run(args) != nil)
}

func newAddApp() *cli.App {
	app := cli.NewApp()
	app.Commands = []*cli.Command{newAddCommand(fakeInserter{})}
	return app
}
