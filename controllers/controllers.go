package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
	"golang.org/x/crypto/bcrypt"
)

// Login Handler
func Login(c echo.Context) error {
	var datos struct {
		Title  string
		Header string
		Body   []string
		Footer string
	}
	datos.Header = views.Comunes.Header
	datos.Title = views.Comunes.Title
	datos.Footer = views.Comunes.Footer
	return c.Render(http.StatusOK, "login.html", datos)
}

// Signup : Registrarse
//func Signup(c echo.Context) error {

//}

// Inicio Handler
func Inicio(c echo.Context) error {
	fmt.Println("Empezando Inicio...")
	var datos struct {
		Title  string
		Header string
		Body   []string
		Footer string
	}
	datos.Header = views.Comunes.Header
	datos.Title = views.Comunes.Title
	datos.Footer = views.Comunes.Footer
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

// Nuevo - controlador para formulario nuevo usuario
func Nuevo(c echo.Context) error {
	datos := struct {
		Username string
		Title    string
		Header   string
		Footer   string
	}{
		Username: "",
		Title:    views.Comunes.Title,
		Header:   views.Comunes.Header,
		Footer:   views.Comunes.Footer,
	}
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}

// Crear usuario
func Crear(c echo.Context) error {
	var usuario models.User
	usuario.Username = c.FormValue("usuario")
	hashed, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	usuario.Password = string(hashed)
	usuario.Email = c.FormValue("email")
	models.Dbcon.Create(&usuario)
	datos := struct {
		Username string
		Title    string
		Header   string
		Footer   string
	}{
		Username: usuario.Username,
		Title:    views.Comunes.Title,
		Header:   views.Comunes.Header,
		Footer:   views.Comunes.Footer,
	}
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}
