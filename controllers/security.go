package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/robbert229/jwt"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
	"golang.org/x/crypto/bcrypt"
)

var cookie http.Cookie
var Secret = "Tequ1squiapan"

// PToken share token with middleware
var PToken string

// Login Handler
func Login(c echo.Context) error {
	datos := PageRender{
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	return c.Render(http.StatusOK, "login.html", datos)
}

// Logout - destroy session
func Logout(c echo.Context) error {
	datos := PageRender{
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	cookie.Name = "jsessionid"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Domain = ""
	return c.Render(http.StatusOK, "index.html", datos)
}

// Nuevo - controlador para formulario nuevo usuario
func Nuevo(c echo.Context) error {
	datos := PageRender{
		Username: "",
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}

// Crear usuario
func Crear(c echo.Context) error {
	datos := PageRender{
		Username: "",
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
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
	datos := PageRender{
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
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
	fmt.Print("Login: Successful")
	datos.Header.User = usuario.Username
	var rol models.Role
	models.Dbcon.Where("id = ?", usuario.RoleID).Find(&rol)
	var err error
	PToken, err = CrearToken(usuario.Username, rol.Role)
	if err != nil {
		return err
	}
	cookie.Name = "frontends1"
	cookie.Value = strings.TrimSpace(PToken)
	cookie.Expires = time.Now().Add(15 * time.Minute)
	cookie.Domain = "localhost"
	c.SetCookie(&cookie)
	fmt.Println(usuario)
	return c.Render(http.StatusOK, "index.html", datos)
}

// CrearToken debe ser reusable
func CrearToken(usuario, rol string) (string, error) {
	algorithm := jwt.HmacSha256(Secret)

	claims := jwt.NewClaim()
	claims.Set("Role", rol)
	claims.Set("User", usuario)
	claims.SetTime("exp", time.Now().Add(time.Hour*2))

	token, err := algorithm.Encode(claims)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return "", err
	}

	fmt.Printf("Token: %s\n", token)
	return token, nil
}

// ValidateToken to check JWT cookie
func ValidateToken(token string) interface{} {
	fmt.Println("Starting token validation...")
	algorithm := jwt.HmacSha256(Secret)
	if algorithm.Validate(token) != nil {
		fmt.Println("ERROR: Invalid token")
		return map[string]string{"User": "", "Role": ""}
	}

	loadedClaims, err := algorithm.Decode(token)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return map[string]string{"User": "", "Role": ""}
	}

	role, err := loadedClaims.Get("Role")
	if err != nil {
		fmt.Println("ERROR: ", err)
		return map[string]string{"User": "", "Role": ""}
	}

	user, err := loadedClaims.Get("User")
	if err != nil {
		fmt.Println("ERROR: ", err)
		return map[string]string{"User": "", "Role": ""}
	}

	roleString, ok := role.(string)
	if !ok {
		fmt.Println("ERROR: ", err)
		return map[string]string{"User": "", "Role": ""}
	}

	userString, ok := user.(string)
	if !ok {
		fmt.Println("ERROR: ", err)
		return map[string]string{"User": "", "Role": ""}
	}

	udatos := make(map[string]string)
	udatos["User"] = userString
	udatos["Role"] = roleString

	return udatos
}
