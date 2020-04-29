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

type fakeGetterCollector struct {
}

func (fake fakeGetterCollector) GetCollectors() ([]storage.Collector, error) {
	return resultDB, nil
}

func addFakeGetCollector(c echo.Context) error {
	return getCollector(c, fakeGetterCollector{})
}

type emptyGetterCollector struct {
}

func (empty emptyGetterCollector) GetCollectors() ([]storage.Collector, error) {
	return []storage.Collector{}, nil
}

func addNilGetCollector(c echo.Context) error {
	return getCollector(c, emptyGetterCollector{})
}

type errGetterCollector struct {
}

func (fake errGetterCollector) GetCollectors() ([]storage.Collector, error) {
	return []storage.Collector{}, errors.New("get collector: internal server error")
}

func addErrGetCollector(c echo.Context) error {
	return getCollector(c, errGetterCollector{})
}

type application struct {
	app *echo.Echo
}

func newAppTest() *application {
	app := echo.New()
	app.GET("/alba/api/coletores", addFakeGetCollector)

	return &application{
		app: app,
	}
}

func TestGetCollector_NotFound(t *testing.T) {
	apitest.New().
		Handler(newAppTest().app).
		Get("/alba/api/coletores/").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetCollector_Sucess(t *testing.T) {
	apitest.New().
		Handler(newAppTest().app).
		Get("/alba/api/coletores").
		Expect(t).
		Body(`[{"city":"João Pessoa", "entity":"Tribunal Regional do Trabalho 13ª Região", "frequency":30, "fu":"PB", "id":"trt13", "limit-month-backward":2, "limit-year-backward":2018, "path":"github.com/dadosjusbr/coletores/trt13", "start-day":5, "update-date":"2009-10-17T20:34:58.651387237Z"}]`).
		Status(http.StatusOK).
		End()
}

func TestGetCollector_Documentation(t *testing.T) {
	app := echo.New()
	app.GET("/alba/api/coletores", addNilGetCollector)

	apitest.New().
		Handler(app).
		Get("/alba/api/coletores").
		Expect(t).
		Body("{\"message\":\"" + msgNotFound + "\", \"docmentation_url\":\"" + docmentationURL + "\"}").
		Status(http.StatusOK).
		End()
}

func TestGetCollector_InternalServerError(t *testing.T) {
	app := echo.New()
	app.GET("/alba/api/coletores", addErrGetCollector)

	apitest.New().
		Handler(app).
		Get("/alba/api/coletores").
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}
