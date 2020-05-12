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

//URLs definition
const (
	collectorsForAPI = "/alba/api/collectors"
	runsForAPI       = "/alba/api/runs/:id"
	index            = "/alba"
	runs             = "/alba/:id"
)

func newApp(dbClient *storage.DBClient) *application {
	api := api{client: dbClient, getter: dbClient}
	view := view{client: dbClient}

	app := echo.New()

	templates := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	app.Static("/static", "templates/assets")

	app.Renderer = templates

	app.GET(index, view.index)
	app.GET(runs, view.executionsByID)
	app.GET(collectorsForAPI, api.getCollectors)
	app.GET(runsForAPI, api.executionsByID)

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

	client, err := storage.NewClientDB(uri)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(); err != nil {
		log.Fatal(err)

	}
	defer client.Disconnect()

	newApp(client).start()
}
