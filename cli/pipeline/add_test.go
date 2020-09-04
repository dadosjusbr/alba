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
		"-from-file=pipelines.json",
	})
	// Output:
	//Pipeline ID: trt13, Repo: github.com/dadosjusbr/coletores
	//Pipeline ID: mppb, Repo: github.com/dadosjusbr/coletores
}

func TestAdd_InvalidParam(t *testing.T) {
	is := is.New(t)
	app := newAddApp()
	args := append([]string{"alba", "add", ""})
	is.True(app.Run(args) != nil)
}

func newAddApp() *cli.App {
	app := cli.NewApp()
	app.Commands = []*cli.Command{NewAddCommand(fakeInserter{})}
	return app
}
