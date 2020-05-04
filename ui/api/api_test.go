package api

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
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

type fakeCollectorsGetter struct {
	resp []storage.Collector
	err  error
}

func (fake fakeCollectorsGetter) GetCollectors() ([]storage.Collector, error) {
	return fake.resp, fake.err
}

type application struct {
	app *echo.Echo
}

func newAppTest() *application {
	app := echo.New()
	app.GET(CollectorsURL, func(c echo.Context) error {
		//fake database result
		return getCollectors(c, fakeCollectorsGetter{resultDB, nil})
	})

	return &application{
		app: app,
	}
}

func TestGetCollector_NotFoundURL(t *testing.T) {
	apitest.New().
		Handler(newAppTest().app).
		Get("/alba/api/collectors/").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetCollector_Sucess(t *testing.T) {
	apitest.New().
		Handler(newAppTest().app).
		Get(CollectorsURL).
		Expect(t).
		Body(`[{"city":"João Pessoa", "entity":"Tribunal Regional do Trabalho 13ª Região", "frequency":30, "fu":"PB", "id":"trt13", "limit-month-backward":2, "limit-year-backward":2018, "path":"github.com/dadosjusbr/coletores/trt13", "start-day":5, "update-date":"2009-10-17T20:34:58.651387237Z"}]`).
		Status(http.StatusOK).
		End()
}

func TestGetCollector_EmptyDB(t *testing.T) {
	app := echo.New()
	app.GET(CollectorsURL, func(c echo.Context) error {
		//result is an empty list
		return getCollectors(c, fakeCollectorsGetter{[]storage.Collector{}, nil})
	})

	apitest.New().
		Handler(app).
		Get(CollectorsURL).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetCollector_InternalServerError(t *testing.T) {
	app := echo.New()
	app.GET(CollectorsURL, func(c echo.Context) error {
		//result is an empty list and the function raise an error
		return getCollectors(c, fakeCollectorsGetter{[]storage.Collector{}, errors.New("get collector: internal server error")})
	})

	apitest.New().
		Handler(app).
		Get(CollectorsURL).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}
