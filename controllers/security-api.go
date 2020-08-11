package controllers

import (
	"fmt"

	"github.com/gofiber/fiber"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
	"golang.org/x/crypto/bcrypt"
)

// RESTChecklogin POST verificar login - JSON payload
func RESTChecklogin(c *fiber.Ctx) {
	var usuario models.User
	type suserpload struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	userpload := new(suserpload)
	if err := c.BodyParser(userpload); err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err.Error())
	}
	fmt.Println(c.Body())
	fmt.Printf("Login payload: %+v\n", userpload)

	models.Dbcon.Where("email = ?", userpload.Email).Find(&usuario)
	datos := fiber.Map{
		"title": views.Comunes.Title,
		"error": "",
	}

	if len(usuario.Email) < 1 {
		datos["error"] = "Correo-e o contraseña incorrectos"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(userpload.Password)); err != nil {
		datos["error"] = "Correo-e o contraseña incorrectos 2"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}
	fmt.Println("API Login: Successful")
	datos["user"] = usuario.Username
	var rol models.Role
	models.Dbcon.Where("id = ?", usuario.RoleID).Find(&rol)
	var err error
	PToken, err = CrearToken(usuario.Username, rol.Role)
	if err != nil {
		fmt.Println("ERROR creating token")
		return
	}
	datos["token"] = PToken

	fmt.Println("API Logged in as: ", usuario.Username)
	fmt.Println("API Role: ", rol.Role)
	c.JSON(datos)
}
