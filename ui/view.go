package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/dadosjusbr/alba/storage"
	"github.com/dadosjusbr/alba/ui/api"

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

func newApp() *application {
	app := echo.New()

	templates := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	app.Static("/static", "templates/assets")

	app.Renderer = templates

	app.GET("/alba", index)
	app.GET("/alba/:id", viewExecutionsByID)
	app.GET("/alba/api/collectors", api.AddGetCollector)
	app.GET("/alba/api/collectors/executions/:id", api.ExecutionByID)

	return &application{
		app: app,
	}
}

func (a *application) start() {
	a.app.Logger.Fatal(a.app.Start(":8080"))
}

func main() {
	newApp().start()
}

//e.GET("/alba", index)
func index(c echo.Context) error {
	results, err := storage.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	//TODO: Retornar página hmtl sem resultados
	if len(results) == 0 {
		return c.String(http.StatusOK, "No results")
	}

	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	data := struct {
		Collectors []storage.Collector
	}{
		results,
	}

	return c.Render(http.StatusOK, "home.html", data)
}

//e.GET("/alba/:id", viewExecutionsByID)
func viewExecutionsByID(c echo.Context) error {
	//mockup
	data := api.ExecutionDetails{
		Entity: "Nome do órgão",
		Executions: []api.Execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.Render(http.StatusOK, "executionsDetails.html", data)
}
