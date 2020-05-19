package main

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

//Template represents the html/template
type Template struct {
	templates *template.Template
}

//Render implements echo Renderer
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type application struct {
	app *echo.Echo
}

type finder interface {
	GetCollectors() ([]storage.Collector, error)
}

//URLs definition
const (
	apiCollectors = "/alba/api/collectors"
	apiRuns       = "/alba/api/runs/:id"
	home          = "/alba"
	runs          = "/alba/:id"
)

func newApp(dbClient *storage.DBClient) *application {
	app := echo.New()

	templates := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	app.Static("/static", "templates/assets")

	app.Renderer = templates

	app.GET(home, func(c echo.Context) error {
		return index(dbClient, c)
	})
	app.GET(runs, func(c echo.Context) error {
		return viewExecutions(dbClient, c)
	})
	app.GET(apiCollectors, func(c echo.Context) error {
		return getCollectors(dbClient, c)
	})
	app.GET(apiRuns, func(c echo.Context) error {
		return getExecutions(dbClient, c)
	})

	return &application{
		app: app,
	}
}

func (a *application) start() {
	a.app.Logger.Fatal(a.app.Start(":8080"))
}

func main() {
	var client *storage.DBClient

	uri := os.Getenv("MONGODB")
	if uri == "" {
		log.Fatal("error trying get environment variable: $MONGODB is empty")
	}

	client, err := storage.NewDBClient(uri)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(); err != nil {
		log.Fatal(err)

	}
	defer client.Disconnect()

	newApp(client).start()
}
