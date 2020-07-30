package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

//var cookie http.Cookie

// FrontData echo.Map type trying to apply DRY
var FrontData echo.Map

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
	datos := echo.Map{
		"title": "Fabs Blog",
		"user":  "",
		"role":  "",
	}

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	return c.Render(http.StatusOK, "index", datos)
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
