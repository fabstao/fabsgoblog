package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
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
var Secret string

// SITEKEY for recaptcha v3
var SITEKEY string

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
	fmt.Println("Attempting to logout user")
	cookie, err := c.Cookie("frontends1")
	if err != nil {
		fmt.Println("WARNING: Cookie could not be read: ", err)
		cookie.Name = "frontends1"
	}
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Domain = ""
	c.SetCookie(cookie)
	return c.Render(http.StatusOK, "bye", datos)
	//return c.Redirect(http.StatusMovedPermanently, "/")
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
		"alerta":       template.HTML("alert alert-success"),
	}
	var usuario models.User
	var checausuario models.User
	var rol models.Role
	user := c.FormValue("usuario")
	email := c.FormValue("email")

	if grecaptcha := validateCaptcha(c.FormValue("g-recaptcha-response")); !grecaptcha {
		datos["mensajeflash"] = "Este sistema sólo es para humanos"
		datos["alerta"] = template.HTML("alert alert-danger")
		return c.Render(http.StatusForbidden, "login", datos)
	}

	if len(user) < 4 || len(email) < 4 {
		datos["mensajeflash"] = "Usuario y correo no pueden estar vacíos o demasiado cortos"
		datos["alerta"] = template.HTML("alert alert-danger")
		return c.Render(http.StatusOK, "crearusuario", datos)
	}
	models.Dbcon.Where("username = ?", user).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos["mensajeflash"] = "Usuario ya existe: " + user
		datos["alerta"] = template.HTML("alert alert-danger")
		return c.Render(http.StatusOK, "crearusuario", datos)
	}
	models.Dbcon.Where("email = ?", email).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos["mensajeflash"] = "Dirección de correo-e ya existe: " + email
		datos["alerta"] = template.HTML("alert alert-danger")
		return c.Render(http.StatusOK, "crearusuario", datos)
	}
	models.Dbcon.Where("role = ?", "usuario").Find(&rol)
	usuario.Username = user
	usuario.Email = email
	hashed, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	usuario.Password = string(hashed)
	usuario.Role = rol

	models.Dbcon.Create(&usuario)
	datos["mensajeflash"] = "Usuario " + usuario.Username + " creado exitosamente"
	return c.Render(http.StatusOK, "crearusuario", datos)
}

// Checklogin POST verificar login
func Checklogin(c echo.Context) error {
	var usuario models.User
	email := c.FormValue("email")
	password := c.FormValue("password")

	models.Dbcon.Where("email = ?", email).Find(&usuario)
	datos := echo.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	if grecaptcha := validateCaptcha(c.FormValue("g-recaptcha-response")); !grecaptcha {
		datos["mensajeflash"] = "Este sistema sólo es para humanos"
		return c.Render(http.StatusOK, "login", datos)
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
	fmt.Println("Login: Successful")
	datos["user"] = usuario.Username
	var rol models.Role
	models.Dbcon.Where("id = ?", usuario.RoleID).Find(&rol)
	var err error
	PToken, err = CrearToken(usuario.Username, rol.Role)
	datos["role"] = rol.Role
	if err != nil {
		return err
	}
	cookie.Name = "frontends1"
	cookie.Value = strings.TrimSpace(PToken)
	cookie.Expires = time.Now().Add(15 * time.Minute)
	cookie.Domain = "localhost"
	c.SetCookie(&cookie)
	fmt.Println("Logged in as: ", usuario.Username)
	fmt.Println("Role: ", rol.Role)
	return c.Redirect(http.StatusMovedPermanently, "/")
	//return c.Render(http.StatusOK, "index", datos)
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

func validateCaptcha(grecaptcha string) bool {
	formData := url.Values{
		"secret":   {SITEKEY},
		"response": {grecaptcha},
	}
	fmt.Println("FormData - g-validation: ", formData)
	gurl := "https://www.google.com/recaptcha/api/siteverify"
	resp, err := http.PostForm(gurl, formData)
	if err != nil {
		fmt.Println("FATAL ERROR: Recaptcha validation error")
		return false
	}
	var resultado map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&resultado)
	fmt.Println(resultado["success"])
	valido := resultado["success"]
	return valido.(bool)
}
