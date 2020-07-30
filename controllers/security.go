package controllers

import (
	"fmt"
	"html/template"
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

// Secret : change here for all JWT signing
var Secret = "Tequ1squiapan"

// PToken share token with middleware
var PToken string

// Login Handler
func Login(c echo.Context) error {
	datos := echo.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	return c.Render(http.StatusOK, "login", datos)
}

// Logout - destroy session
func Logout(c echo.Context) error {
	datos := echo.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	dcookie, _ := c.Cookie("frontends1")
	dcookie.Value = ""
	dcookie.Expires = time.Now()
	dcookie.Domain = ""
	c.SetCookie(dcookie)
	return c.Render(http.StatusOK, "index", datos)
}

// Nuevo - controlador para formulario nuevo usuario
func Nuevo(c echo.Context) error {
	datos := echo.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	return c.Render(http.StatusOK, "crearusuario", datos)
}

// Crear usuario
func Crear(c echo.Context) error {
	datos := echo.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("\"alert alert-success\""),
	}
	var usuario models.User
	var checausuario models.User
	var rol models.Role
	user := c.FormValue("usuario")
	email := c.FormValue("email")
	if len(user) < 4 || len(email) < 4 {
		datos["mensajeflash"] = "Usuario y correo no pueden estar vacíos o demasiado cortos"
		datos["alerta"] = "\"alert alert-danger\""
		return c.Render(http.StatusOK, "crearusuario", datos)
	}
	models.Dbcon.Where("username = ?", user).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos["mensajeflash"] = "Usuario ya existe: " + user
		datos["alerta"] = "\"alert alert-danger\""
		return c.Render(http.StatusOK, "crearusuario", datos)
	}
	models.Dbcon.Where("email = ?", email).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos["mensajeflash"] = "Dirección de correo-e ya existe: " + email
		datos["alerta"] = "\"alert alert-danger\""
		return c.Render(http.StatusOK, "crearusuario", datos)
	}
	models.Dbcon.Where("role = ?", "usuario").Find(&rol)
	usuario.Username = user
	usuario.Email = email
	hashed, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	usuario.Password = string(hashed)
	usuario.Role = rol

	fmt.Println(usuario)
	models.Dbcon.Create(&usuario)
	datos["mensajeflash"] = "Usuario " + usuario.Username + " creado exitosamente"
	return c.Render(http.StatusOK, "crearusuario", datos)
}

// Checklogin POST verificar login
func Checklogin(c echo.Context) error {
	var usuario models.User
	email := c.FormValue("email")
	password := c.FormValue("password")
	fmt.Print("Login: ", email)
	models.Dbcon.Where("email = ?", email).Find(&usuario)
	datos := echo.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	fmt.Println(datos)
	if len(usuario.Email) < 1 {
		datos["mensajeflash"] = "Correo-e o contraseña incorrectos"
		return c.Render(http.StatusOK, "login", datos)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password)); err != nil {
		datos["mensajeflash"] = "Correo-e o contraseña incorrectos"
		return c.Render(http.StatusOK, "login", datos)
	}
	fmt.Print("Login: Successful")
	datos["user"] = usuario.Username
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
	return c.Render(http.StatusOK, "index", datos)
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
