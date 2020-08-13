package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
	"golang.org/x/crypto/bcrypt"
)

// Security types, vars and configuration

// Secret : change here for all JWT signing
var Secret string

// SITEKEY for recaptcha v3
var SITEKEY string

// PToken share token with middleware
var PToken string

// Cdomain - domain for session cookie
var Cdomain string

var cookie fiber.Cookie

// RefreshFCookie used to keep session alive
func RefreshFCookie(token string) {
	cookie.Name = "frontends1"
	cookie.Value = token
	cookie.Expires = time.Now().Add(15 * time.Minute)
	cookie.Domain = Cdomain
}

// Login Handler
func Login(c *fiber.Ctx) {
	datos := fiber.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	c.Render("login", datos, "layouts/main")
}

// Logout - destroy session
func Logout(c *fiber.Ctx) {
	datos := fiber.Map{
		"title": views.Comunes.Title,
		"user":  "",
		"role":  "",
	}
	fmt.Println("Attempting to logout user")
	c.ClearCookie()
	c.ClearCookie("frontends1")
	c.Render("bye", datos, "layouts/main")
}

// Nuevo - controlador para formulario nuevo usuario
func Nuevo(c *fiber.Ctx) {
	datos := fiber.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	c.Render("crearusuario", datos, "layouts/main")
}

// Crear usuario
func Crear(c *fiber.Ctx) {
	datos := fiber.Map{
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
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}

	if len(user) < 4 || len(email) < 4 {
		datos["mensajeflash"] = "Usuario y correo no pueden estar vacíos o demasiado cortos"
		datos["alerta"] = template.HTML("alert alert-danger")
		c.Render("crearusuario", datos, "layouts/main")
		return
	}
	models.Dbcon.Where("username = ?", user).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos["mensajeflash"] = "Usuario ya existe: " + user
		datos["alerta"] = template.HTML("alert alert-danger")
		c.Render("crearusuario", datos, "layouts/main")
		return
	}
	models.Dbcon.Where("email = ?", email).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos["mensajeflash"] = "Dirección de correo-e ya existe: " + email
		datos["alerta"] = template.HTML("alert alert-danger")
		c.Render("crearusuario", datos, "layouts/main")
		return
	}
	models.Dbcon.Where("role = ?", "usuario").Find(&rol)
	usuario.Username = user
	usuario.Email = email
	hashed, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	usuario.Password = string(hashed)
	usuario.Role = rol

	models.Dbcon.Create(&usuario)
	datos["mensajeflash"] = "Usuario " + usuario.Username + " creado exitosamente"
	c.Render("crearusuario", datos, "layouts/main")
}

// Checklogin POST verificar login
func Checklogin(c *fiber.Ctx) {
	var usuario models.User
	email := c.FormValue("email")
	password := c.FormValue("password")

	models.Dbcon.Where("email = ?", email).Find(&usuario)
	datos := fiber.Map{
		"title":        views.Comunes.Title,
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}
	if grecaptcha := validateCaptcha(c.FormValue("g-recaptcha-response")); !grecaptcha {
		datos["mensajeflash"] = "Este sistema sólo es para humanos"
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}

	fmt.Println(datos, "layouts/main")
	if len(usuario.Email) < 1 {
		datos["mensajeflash"] = "Correo-e o contraseña incorrectos"
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password)); err != nil {
		datos["mensajeflash"] = "Correo-e o contraseña incorrectos"
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}
	fmt.Println("Login: Successful")
	datos["user"] = usuario.Username
	var rol models.Role
	models.Dbcon.Where("id = ?", usuario.RoleID).Find(&rol)
	var err error
	PToken, err = CrearToken(usuario.Username, rol.Role)
	datos["role"] = rol.Role
	if err != nil {
		fmt.Println("ERROR creating token")
		return
	}
	cookie.Name = "frontends1"
	cookie.Value = strings.TrimSpace(PToken)
	cookie.Expires = time.Now().Add(15 * time.Minute)
	cookie.Domain = Cdomain
	c.Cookie(&cookie)
	fmt.Println("Logged in as: ", usuario.Username)
	fmt.Println("Role: ", rol.Role)
	c.Redirect("/", http.StatusMovedPermanently)
}

// CrearToken debe ser reusable
func CrearToken(usuario, rol string) (string, error) {
	// Create token
	btoken := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := btoken.Claims.(jwt.MapClaims)
	claims["user"] = usuario
	claims["role"] = rol
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	// Generate encoded token and send it as response.
	token, err := btoken.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	fmt.Printf("Token: %s\n", token)
	return token, nil
}

// ValidateToken to check JWT cookie
func ValidateToken(token string) interface{} {
	fmt.Println("Starting token validation...")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		fmt.Println(err)
		return map[string]string{"User": "", "Role": ""}
	}

	if err != nil {
		fmt.Println("ERROR: ", err)
		return map[string]string{"User": "", "Role": ""}
	}

	fmt.Println(claims)

	udatos := make(map[string]string)
	udatos["User"] = claims["user"].(string)
	udatos["Role"] = claims["role"].(string)

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
