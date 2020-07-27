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

// Inicio Handler
func Inicio(c echo.Context) error {
	token := "null"
	if vcookie, err := c.Cookie("frontends1"); err != nil {
		fmt.Println("ERROR reading cookie: ", err)
	} else {
		token = vcookie.Value
		fmt.Println("TOKEN - index: ", token)
		vcookie.Expires = time.Now().Add(15 * time.Minute)
		c.SetCookie(vcookie)
	}

	clamas := ValidateToken(token).(map[string]string)
	fmt.Println("Empezando Blog...")
	fmt.Println(clamas)
	var datos struct {
		Title  string
		Header Titulo
		Body   []string
		Footer string
		User   string
		Role   string
	}
	datos.Header = Titulo{
		"Fabs BLOG",
		"",
		"",
	}
	datos.Header.Title = views.Comunes.Title
	datos.Footer = views.Comunes.Footer
	datos.Header.User = clamas["User"]
	datos.Header.Role = clamas["Role"]
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
