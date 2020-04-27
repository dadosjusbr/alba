package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

var resultDB []storage.Collector = []storage.Collector{
	{
		ID:                 "trt13",
		Entity:             "Tribunal Regional do Trabalho 13ª Região",
		City:               "João Pessoa",
		FU:                 "PB",
		UpdateDate:         time.Date(2009, 10, 17, 20, 34, 58, 651387237, time.UTC),
		Path:               "github.com/dadosjusbr/coletores/trt13",
		Frequency:          30,
		StartDay:           5,
		LimitMonthBackward: 2,
		LimitYearBackward:  2018,
	},
}

type fakeGetterCollector struct {
}

func (fake fakeGetterCollector) GetCollectors() ([]storage.Collector, error) {
	return resultDB, nil
}

func addFakeGetCollector(c echo.Context) error {
	return getCollectors(c, fakeGetterCollector{})
}

type application struct {
	app *echo.Echo
}

func newApp() *application {
	app := echo.New()

	app.GET("/teste", addFakeGetCollector)

	return &application{
		app: app,
	}
}

func TestGetCollectors(t *testing.T) {
	apitest.New().
		Handler(newApp().app).
		Get("/teste").
		Expect(t).
		Assert(jsonpath.Equal(`{"city":"João Pessoa", "entity":"Tribunal Regional do Trabalho 13ª Região", "frequency":30, "fu":"PB", "id":"trt13", "limit-month-backward":2, "limit-year-backward":2018, "path":"github.com/dadosjusbr/coletores/trt13", "start-day":5, "update-date":"2009-10-17T20:34:58.651387237Z"}`)).
		Status(http.StatusOK).
		End()
}
