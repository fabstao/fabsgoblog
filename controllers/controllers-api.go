package controllers

import (
	"fmt"
	"math/rand"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"gitlab.com/fabstao/fabsgoblog/models"
)

//var cookie http.Cookie

/*
type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Path     string    `json:"path"`
	Domain   string    `json:"domain"`
	Expires  time.Time `json:"expires"`
	Secure   bool      `json:"secure"`
	HTTPOnly bool      `json:"http_only"`
	SameSite string    `json:"same_site"`
}
*/

// AllPosts Handler
func AllPosts(c *fiber.Ctx) {

	datos := fiber.Map{
		"title": "Fabs Blog",
	}

	type tsposts struct {
		ID     uint
		Fecha  string
		Titulo string
		Autor  string
	}
	spost := tsposts{
		ID:     0,
		Fecha:  "",
		Titulo: "",
		Autor:  "",
	}
	var sposts []tsposts
	var posts []models.Post
	models.Dbcon.Find(&posts).Limit(10)
	var autor models.User
	for k, v := range posts {
		sposts = append(sposts, spost)
		models.Dbcon.Where("id = ?", v.UserID).Find(&autor)
		sposts[k].ID = v.ID
		sposts[k].Titulo = v.Titulo
		sposts[k].Autor = autor.Username
		sposts[k].Fecha = v.UpdatedAt.String()
		autor = models.User{}
	}
	datos["entradas"] = sposts
	fmt.Println("Mandando Entradas")
	c.JSON(datos)
}

// RESTShow - Show post form - anyone can read
func RESTShow(c *fiber.Ctx) {
	datos := fiber.Map{
		"title": "Fabs Blog",
	}
	var post models.Post
	var autor models.User
	//pid := c.ParamNames()
	pidv := c.Params("id")
	models.Dbcon.Where("id = ?", pidv).Find(&post)
	models.Dbcon.Where("id = ?", post.UserID).Find(&autor)
	datos["titulo"] = post.Titulo
	datos["texto"] = post.Texto
	datos["autor"] = autor.Username
	datos["fecha"] = post.UpdatedAt.String()
	datos["pid"] = pidv
	c.JSON(datos)
}

// RESTPost post
func RESTPost(c *fiber.Ctx) {
	datos := fiber.Map{
		"title": "Fabs Blog",
		"user":  "",
		"error": "",
	}
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		fmt.Println("ERROR reading cookie ")
		datos["error"] = "TOKEN authentication error"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}
	//fmt.Println("TOKEN - index: ", user)
	fmt.Println("Manejando TOKEN")
	clamas := user.Claims.(jwt.MapClaims)
	fmt.Printf("%+v\n", clamas)

	datos["user"] = clamas["user"]

	var usuario models.User
	var post models.Post

	npost := new(struct {
		Titulo string `json:"titulo" form:"titulo"`
		Texto  string `json:"texto" form:"texto"`
	})

	if err := c.BodyParser(npost); err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err.Error())
		datos["error"] = "Error en los datos mandados"
		c.Status(fiber.StatusInternalServerError).JSON(datos)
		return
	}

	if len(npost.Titulo) < 4 || len(npost.Texto) < 4 {
		datos["error"] = "Título o texto demasiado cortos"
		c.Status(fiber.StatusInternalServerError).JSON(datos)
		return
	}

	post.Titulo = npost.Titulo
	post.Texto = npost.Texto
	models.Dbcon.Where("username = ?", datos["user"].(string)).Find(&usuario)
	post.User = usuario
	models.Dbcon.Create(&post)
	datos["error"] = ""
	datos["status"] = "Entrada creada exitosamente"
	datos["titulo"] = npost.Titulo
	c.JSON(datos)
}

// RESTUpdate post
func RESTUpdate(c *fiber.Ctx) {
	datos := fiber.Map{
		"title": "Fabs Blog",
		"user":  "",
		"error": "",
	}
	token := c.Locals("token").(*jwt.Token)
	if token == nil {
		fmt.Println("ERROR leyendo token o token inválido")
		datos["error"] = "ERROR leyendo token o token inválido"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}
	fmt.Println("TOKEN - index: ", token)

	clamas := token.Claims.(jwt.MapClaims)
	fmt.Println(clamas)

	datos["user"] = clamas["user"]

	var usuario models.User
	var post models.Post

	epost := new(struct {
		ID     uint   `json:"id" form:"id"`
		Titulo string `json:"titulo" form:"titulo"`
		Texto  string `json:"texto" form:"texto"`
	})

	if err := c.BodyParser(epost); err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err.Error())
		datos["error"] = "Error en los datos mandados"
		c.Status(fiber.StatusInternalServerError).JSON(datos)
		return
	}

	models.Dbcon.Where("username = ?", datos["user"].(string)).Find(&usuario)

	if datos["user"] != usuario.Username {
		fmt.Println("ERROR de Seguridad")
		datos["error"] = "ERROR de seguridad, no se puede editar entrada de otro usuario"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}

	if err := models.Dbcon.Where("id = ?", epost.ID).Find(&post); err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err)
		datos["error"] = "Entrada no encontrada"
		c.Status(fiber.StatusInternalServerError).JSON(datos)
		return
	}

	if len(epost.Titulo) < 4 || len(epost.Texto) < 4 {
		datos["error"] = "Título o texto demasiado cortos"
		c.Status(fiber.StatusInternalServerError).JSON(datos)
		return
	}

	post.Titulo = epost.Titulo
	post.Texto = epost.Texto
	models.Dbcon.Save(&post)
	datos["mensajeflash"] = "Entrada " + post.Titulo + " actualizada correctamente"
	datos["titulo"] = epost.Titulo
	datos["texto"] = epost.Texto
	datos["id"] = epost.ID
	c.JSON(datos)
}

// RESTDelete - danger, take care
func RESTDelete(c *fiber.Ctx) {
	datos := fiber.Map{
		"title": "Fabs Blog",
		"user":  "",
		"error": "",
	}
	token := c.Locals("token").(*jwt.Token)
	if token == nil {
		fmt.Println("ERROR leyendo token o token inválido")
		datos["error"] = "ERROR leyendo token o token inválido"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}
	fmt.Println("TOKEN - index: ", token)

	clamas := token.Claims.(jwt.MapClaims)
	fmt.Println(clamas)

	datos["user"] = clamas["user"]

	var usuario models.User
	var post models.Post
	pidv := c.Params("id")
	models.Dbcon.Where("id = ?", pidv).Find(&post)
	models.Dbcon.Where("username = ?", post.UserID).Find(&usuario)
	if usuario.Username != datos["user"] {
		fmt.Println("ERROR de Seguridad")
		datos["error"] = "ERROR de seguridad, no se puede borrar entrada de otro usuario"
		c.Status(fiber.StatusForbidden).JSON(datos)
		return
	}

	models.Dbcon.Unscoped().Delete(&post)
	datos["status"] = "Entrada borrada correctamente"
	c.JSON(datos)
}

// Hello REST example
func Hello(c *fiber.Ctx) {
	nombre := c.Query("nombre")
	var content struct {
		Response  string    `json:"response"`
		Timestamp time.Time `json:"timestamp"`
		Random    int       `json:"random"`
	}
	content.Response = "Hola " + nombre
	content.Timestamp = time.Now().UTC()
	content.Random = rand.Intn(1000)
	c.JSON(content)
}
