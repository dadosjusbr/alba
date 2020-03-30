package collector

import (
	"testing"

	"github.com/dadosjusbr/alba/storage"

	"github.com/matryer/is"
	"github.com/urfave/cli/v2"
)

type fakeInserter struct {
}

func (fakeInserter) InsertCollector(storage.Collector) error {
	return nil
}

func ExampleAddCommand() {
	app := newAddApp()
	app.Run([]string{
		"alba",
		"add-collector",
		"-from-file=input.json",
	})
	// Output: {ID:trt13 Entity:Tribunal Regional do Trabalho 13ª Região City:João Pessoa FU:PB UpdateDate:0001-01-01 00:00:00 +0000 UTC Path:github.com/dadosjusbr/coletores/trt13 Frequency:30 StartDay:5 LimitMonthBackward:2 LimitYearBackward:2018}
}

func TestAdd_Sucess(t *testing.T) {
	is := is.New(t)
	app := newAddApp()
	args := []string{"alba", "add-collector", "-id=1", "-entity=1", "-city=1", "-fu=1", "-path=1", "-frequency=1", "-start-day=1", "-limit-month-backward=1", "-limit-year-backward=1"}
	is.NoErr(app.Run(args))
}

func TestAdd_InvalidParams(t *testing.T) {
	is := is.New(t)
	app := newAddApp()
	tests := []struct {
		desc string
		args []string
	}{
		{"EmptyID", []string{"-id"}},
		{"EmptyEntity", []string{"-id=1", "-entity"}},
		{"EmptyCity", []string{"-id=1", "-entity=1", "-city"}},
		{"EmptyFU", []string{"-id=1", "-entity=1", "-city=1", "-fu"}},
		{"EmptyPath", []string{"-id=1", "-entity=1", "-city=1", "-fu=1", "-path"}},
		{"EmptyFrequency", []string{"-id=1", "-entity=1", "-city=1", "-fu=1", "-path=1", "-frequency"}},
		{"EmptyStartDay", []string{"-id=1", "-entity=1", "-city=1", "-fu=1", "-path=1", "-frequency=1", "-start-day"}},
		{"EmptyLimitMonthBackward", []string{"-id=1", "-entity=1", "-city=1", "-fu=1", "-path=1", "-frequency=1", "-start-day=1", "-limit-month-backward"}},
		{"EmptyLimitYearBackward", []string{"-id=1", "-entity=1", "-city=1", "-fu=1", "-path=1", "-frequency=1", "-start-day=1", "-limit-month-backward=1", "-limit-year-backward"}},
	}
	for _, test := range tests {
		args := append([]string{"alba", "add-collector"}, test.args...)
		is.True(app.Run(args) != nil)
	}
}

func newAddApp() *cli.App {
	app := cli.NewApp()
	app.Commands = []*cli.Command{addCommand(&fakeInserter{})}
	return app
}
