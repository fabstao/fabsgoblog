package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Inicio Handler
func Inicio(c echo.Context) error {
	fmt.Println("Empezando Inicio...")
	var datos struct {
		Head   map[string]string
		Header map[string]string
		Body   []string
		Footer map[string]string
	}
	datos.Head = map[string]string{
		"Title": "Fabs BLOG",
	}
	datos.Header = map[string]string{
		"Title": "Fabs BLOG",
	}
	datos.Body = append(datos.Body, "Item 1")
	datos.Body = append(datos.Body, "Item 2")
	datos.Body = append(datos.Body, "Item 3")
	datos.Footer = map[string]string{
		"Footer": "Page 1",
	}
	return c.Render(http.StatusOK, "index.html", datos)
}

// Hello REST example
func Hello(c echo.Context) error {
	nombre := c.QueryParam("nombre")
	var content struct {
		Response  string    `json:"response"`
		Timestamp time.Time `json:"timestamp"`
		Random    int       `json:"random"`
	}
	content.Response = "Hola " + nombre
	content.Timestamp = time.Now().UTC()
	content.Random = rand.Intn(1000)
	return c.JSON(http.StatusOK, content)
}
