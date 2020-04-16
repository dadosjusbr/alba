package main

import (
	"net/http"

	"github.com/dadosjusbr/alba/storage"

	"github.com/labstack/echo/v4"
)

//Message represents a information for the user
type Message struct {
	Message string `json:"message"`
}

const msgNoResults = "No results"

func main() {
	e := echo.New()

	e.GET("/alba", index)
	e.GET("/alba/api/coletores", viewAPIAllCollectors)
	e.GET("/alba/api/coletores/execucoes/:id", viewAPIExecutionsByID)

	e.Logger.Fatal(e.Start(":8080"))
}

//e.GET("/alba", index)
func index(c echo.Context) error {
	return c.String(http.StatusOK, "HTML com lista dos coletores")
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

	return c.JSON(http.StatusOK, result)
}

//e.GET("/alba/api/coletores/execucoes/:id", viewAPIExecutionsByID)
func viewAPIExecutionsByID(c echo.Context) error {
	id := c.Param("id")
	msg := "Lista de execuções " + id

	return c.JSON(http.StatusOK, Message{msg})
}
