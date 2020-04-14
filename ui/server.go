package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

//Message represents a information for the user
type Message struct {
	Message string `json:"message"`
}

const msgNoResults = "No results"

//Template represents the html/template
type Template struct {
	templates *template.Template
}

//Render implements echo Renderer
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//Execution represents a execution
type Execution struct {
	Date   string
	Status string
	Result string
}

//ExecutionDetails represents the details of the execution from collector
type ExecutionDetails struct {
	Entity     string
	Executions []Execution
}

func main() {
	e := echo.New()

	templates := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Static("/static", "templates/assets")

	e.Renderer = templates

	e.GET("/alba", index)
	e.GET("/alba/:id", viewExectuionByID)
	e.GET("/alba/api/coletores", viewAPIAllCollectors)
	e.GET("/alba/api/coletores/:id", viewAPIExecutionByID)

	e.Logger.Fatal(e.Start(":8080"))
}

//e.GET("/alba", index)
func index(c echo.Context) error {
	result, err := storage.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	//TODO: Retornar página hmtl sem resultados
	if result == nil {
		return c.JSON(http.StatusOK, Message{msgNoResults})
	}

	collectors := []storage.Collector{}
	err = json.Unmarshal(result, &collectors)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	data := struct {
		Collectors []storage.Collector
	}{
		collectors,
	}

	return c.Render(http.StatusOK, "home.html", data)
}

//TODO: Implementar a busca da informação das execuções
//e.GET("/alba/:id", viewExectuionByID)
func viewExectuionByID(c echo.Context) error {
	data := ExecutionDetails{
		Entity: "Tribunal Regional do Trabalho 13ª Região",
		Executions: []Execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.Render(http.StatusOK, "executionsDetails.html", data)
}

//e.GET("/alba/api/coletores", viewAPIAllCollectors)
func viewAPIAllCollectors(c echo.Context) error {
	result, err := storage.GetCollectors()
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	if result == nil {
		return c.JSON(http.StatusOK, Message{msgNoResults})
	}

	return c.JSONBlob(http.StatusOK, result)
}

//TODO: Implementar a busca da informação das execuções
//e.GET("/alba/api/coletores/:id", viewAPIExecutionByID)
func viewAPIExecutionByID(c echo.Context) error {
	data := ExecutionDetails{
		Entity: "Tribunal Regional do Trabalho 13ª Região",
		Executions: []Execution{
			{Date: "10/01/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/02/2020", Status: "Finalizado com sucesso", Result: "link"},
			{Date: "10/03/2020", Status: "Finalizado com erro", Result: "link"},
			{Date: "11/03/2020", Status: "Finalizado com sucesso", Result: "link"},
		},
	}

	return c.JSON(http.StatusOK, data)
}
