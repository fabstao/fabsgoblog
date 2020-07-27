package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gitlab.com/fabstao/fabsgoblog/views"
)

//var cookie http.Cookie

// Titulo for PageRender
type Titulo struct {
	Title string
	User  string
	Role  string
}

// PageRender struct helper for security rendering
type PageRender struct {
	Username     string
	Title        string
	Header       Titulo
	Footer       string
	MensajeFlash string
	Alerta       string
}

// ObtenClamas read from cookie
func ObtenClamas() (map[string]string, error) {
	//sdatos := make(map[string]string)
	sdatos, err := ValidateToken(cookie.Value)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return sdatos.(map[string]string), nil
}

// Inicio Handler
func Inicio(c echo.Context) error {
	fmt.Println("Empezando Blog...")
	sesion, err := ObtenClamas()
	if err != nil {
		sesion := make(map[string]string)
		sesion["user"] = ""
		sesion["rol"] = ""
	}
	var datos struct {
		Title  string
		Header Titulo
		Body   []string
		Footer string
		User   string
		Rol    string
	}
	datos.Header = Titulo{
		"Fabs BLOG",
		"",
		"",
	}
	datos.Title = views.Comunes.Title
	datos.Footer = views.Comunes.Footer
	datos.User = sesion["user"]
	datos.Rol = sesion["rol"]
	datos.Body = append(datos.Body, "Item 1")
	datos.Body = append(datos.Body, "Item 2")
	datos.Body = append(datos.Body, "Item 3")
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
