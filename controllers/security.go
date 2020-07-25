package controllers

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
	"golang.org/x/crypto/bcrypt"
)

var cookie http.Cookie

// Login Handler
func Login(c echo.Context) error {
	datos := PageRender{
		Title:        views.Comunes.Title,
		Header:       views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	return c.Render(http.StatusOK, "login.html", datos)
}

// Nuevo - controlador para formulario nuevo usuario
func Nuevo(c echo.Context) error {
	datos := PageRender{
		Username:     "",
		Title:        views.Comunes.Title,
		Header:       views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}

// Crear usuario
func Crear(c echo.Context) error {
	datos := PageRender{
		Username:     "",
		Title:        views.Comunes.Title,
		Header:       views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
		Alerta:       "success",
	}
	var usuario models.User
	var checausuario models.User
	var rol models.Role
	user := c.FormValue("usuario")
	email := c.FormValue("email")
	if len(user) < 4 || len(email) < 4 {
		datos.MensajeFlash = "Usuario y correo no pueden estar vacíos o demasiado cortos"
		datos.Alerta = "danger"
		return c.Render(http.StatusOK, "crearusuario.html", datos)
	}
	models.Dbcon.Where("username = ?", user).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos.MensajeFlash = "Usuario ya existe: " + user
		datos.Alerta = "danger"
		return c.Render(http.StatusOK, "crearusuario.html", datos)
	}
	models.Dbcon.Where("email = ?", email).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos.MensajeFlash = "Dirección de correo-e ya existe: " + email
		datos.Alerta = "danger"
		return c.Render(http.StatusOK, "crearusuario.html", datos)
	}
	models.Dbcon.Where("role = ?", "usuario").Find(&rol)
	usuario.Username = user
	usuario.Email = email
	hashed, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	usuario.Password = string(hashed)
	usuario.Role = rol

	fmt.Println(usuario)
	models.Dbcon.Create(&usuario)
	datos.MensajeFlash = "Usuario " + usuario.Username + " creado exitosamente"
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}

// Checklogin POST verificar login
func Checklogin(c echo.Context) error {
	var usuario models.User
	email := c.FormValue("email")
	password := c.FormValue("password")
	fmt.Print("Login: ", email)
	models.Dbcon.Where("email = ?", email).Find(&usuario)
	fmt.Print("Login: Successful")
	datos := PageRender{
		Title:        views.Comunes.Title,
		Header:       views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	fmt.Println(datos)
	if len(usuario.Email) < 1 {
		datos.MensajeFlash = "Correo-e o contraseña incorrectos"
		return c.Render(http.StatusOK, "login.html", datos)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password)); err != nil {
		datos.MensajeFlash = "Correo-e o contraseña incorrectos"
		return c.Render(http.StatusOK, "login.html", datos)
	}
	var rol models.Role
	models.Dbcon.Where("id = ?", usuario.RoleID).Find(&rol)
	token := jwt.New(jwt.SigningMethodHS256)
	clamas := token.Claims.(jwt.MapClaims)
	clamas["user"] = usuario.Username
	clamas["rol"] = rol.Role
	clamas["exp"] = time.Now().Add(time.Hour).Unix()
	t, err := token.SignedString([]byte("tequisquiapan"))
	if err != nil {
		return err
	}
	cookie.Name = "jsessionid"
	cookie.Value = t
	cookie.Expires = time.Now().Add(5 * time.Minute)
	cookie.Domain = "localhost"
	c.SetCookie(&cookie)
	fmt.Println(usuario)
	return c.Render(http.StatusOK, "index.html", datos)
}
