package controllers

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

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

// FrontData fiber.Map type trying to apply DRY
var FrontData fiber.Map

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
func Inicio(c *fiber.Ctx) {
	token := ""
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
	} else {
		fmt.Println("TOKEN - index: ", token)
		RefreshFCookie(token)
		c.Cookie(&cookie)
	}

	clamas := ValidateToken(token).(map[string]string)
	fmt.Println("Empezando Blog...")
	datos := fiber.Map{
		"title": "Fabs Blog",
		"user":  "",
		"role":  "",
	}

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	type tsposts struct {
		ID     uint
		Fecha  string
		Titulo string
		Autor  string
		User   string
	}
	spost := tsposts{
		ID:     0,
		Fecha:  "",
		Titulo: "",
		Autor:  "",
		User:   clamas["User"],
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
	fmt.Println("Empezando render")
	c.Status(fiber.StatusOK).Render("index", datos, "layouts/main")
}

// Show - Show post form - anyone can read
func Show(c *fiber.Ctx) {
	token := ""
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
	} else {
		fmt.Println("TOKEN - index: ", token)
		RefreshFCookie(token)
		c.Cookie(&cookie)
	}
	clamas := ValidateToken(token).(map[string]string)
	datos := fiber.Map{
		"title":        "Fabs Blog",
		"mensajeflash": "",
	}
	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	var post models.Post
	var autor models.User
	//pid := c.ParamNames()
	pidv := c.Params("id")
	models.Dbcon.Where("id = ?", pidv).Find(&post)
	models.Dbcon.Where("id = ?", post.UserID).Find(&autor)
	datos["titulo"] = post.Titulo
	datos["mitexto"] = post.Texto
	datos["autor"] = autor.Username
	datos["fecha"] = post.UpdatedAt.String()
	datos["pid"] = pidv
	if datos["user"] == autor.Username {
		datos["editar"] = true
	}
	c.Status(fiber.StatusOK).Render("show", datos, "layouts/main")
}

// Post - New post form
func Post(c *fiber.Ctx) {
	token := "null"
	datos := fiber.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
		c.Status(fiber.StatusForbidden).Render("login", datos, "layouts/main")
		return
	}
	fmt.Println("TOKEN - index: ", token)
	RefreshFCookie(token)
	c.Cookie(&cookie)

	clamas := ValidateToken(token).(map[string]string)
	fmt.Println("Empezando Blog...")
	fmt.Println(clamas)

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	c.Render("new", datos, "layouts/main")
}

// New post
func New(c *fiber.Ctx) {
	token := "null"
	datos := fiber.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
		c.Status(fiber.StatusForbidden).Render("login", datos, "layouts/main")
		return
	}
	fmt.Println("TOKEN - index: ", token)
	RefreshFCookie(token)
	c.Cookie(&cookie)

	clamas := ValidateToken(token).(map[string]string)

	fmt.Println("User: ", clamas["User"])
	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]

	var usuario models.User
	var post models.Post
	titulo := c.FormValue("titulo")
	texto := c.FormValue("texto")

	if grecaptcha := validateCaptcha(c.FormValue("g-recaptcha-response")); !grecaptcha {
		datos["mensajeflash"] = "Este sistema sólo es para humanos"
		datos["alerta"] = template.HTML("alert alert-danger")
		c.Status(fiber.StatusForbidden).Render("login", datos, "layouts/main")
		return
	}

	if len(titulo) < 4 || len(texto) < 4 {
		datos["mensajeflash"] = "Título o texto demasiado cortos"
		datos["alerta"] = template.HTML("alert alert-danger")
		//c.SendStatus(http.StatusForbidden)
		c.Render("new", datos, "layouts/main")
		return
	}

	post.Titulo = titulo
	post.Texto = texto
	models.Dbcon.Where("username = ?", datos["user"].(string)).Find(&usuario)
	post.User = usuario
	models.Dbcon.Create(&post)
	datos["mensajeflash"] = "Entrada " + post.Titulo + " creada correctamente"
	datos["titulo"] = titulo
	datos["texto"] = texto
	c.Render("new", datos, "layouts/main")
}

// Edit - New post form
func Edit(c *fiber.Ctx) {
	token := "null"
	datos := fiber.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}
	fmt.Println("TOKEN - index: ", token)
	RefreshFCookie(token)
	c.Cookie(&cookie)

	clamas := ValidateToken(token).(map[string]string)
	fmt.Println("Empezando Blog...")
	fmt.Println(clamas)

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	var post models.Post
	var autor models.User
	pidv := c.Params("id")
	models.Dbcon.Where("id = ?", pidv).Find(&post)
	models.Dbcon.Where("id = ?", post.UserID).Find(&autor)
	if autor.Username != datos["user"] {
		fmt.Println("SECURITY ERROR: Unauthorized access")
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}
	datos["titulo"] = post.Titulo
	datos["mitexto"] = post.Texto
	datos["autor"] = autor.Username
	datos["pid"] = pidv[0]
	c.Render("edit", datos, "layouts/main")
}

// Update post
func Update(c *fiber.Ctx) {
	token := "null"
	datos := fiber.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}
	fmt.Println("TOKEN - index: ", token)
	RefreshFCookie(token)
	c.Cookie(&cookie)

	clamas := ValidateToken(token).(map[string]string)

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]

	var usuario models.User
	var post models.Post
	titulo := c.FormValue("titulo")
	texto := c.FormValue("texto")
	pid := c.FormValue("pid")

	if grecaptcha := validateCaptcha(c.FormValue("g-recaptcha-response")); !grecaptcha {
		datos["mensajeflash"] = "Este sistema sólo es para humanos"
		datos["alerta"] = template.HTML("alert alert-danger")
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}

	if len(titulo) < 4 || len(texto) < 4 {
		datos["mensajeflash"] = "Título o texto demasiado cortos"
		datos["alerta"] = template.HTML("alert alert-danger")
		c.Render("new", datos, "layouts/main")
		return
	}

	models.Dbcon.Where("username = ?", datos["user"].(string)).Find(&usuario)
	models.Dbcon.Where("id = ?", pid).Find(&post)
	post.Titulo = titulo
	post.Texto = texto
	models.Dbcon.Save(&post)
	datos["mensajeflash"] = "Entrada " + post.Titulo + " actualizada correctamente"
	datos["titulo"] = titulo
	datos["mitexto"] = texto
	datos["pid"] = pid
	c.Render("edit", datos, "layouts/main")
}

// Delete - danger, take care
func Delete(c *fiber.Ctx) {
	token := "null"
	datos := fiber.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if token = c.Cookies("frontends1"); token == "" {
		fmt.Println("ERROR reading cookie ")
		c.SendStatus(http.StatusForbidden)
		c.Render("login", datos, "layouts/main")
		return
	}
	fmt.Println("TOKEN - index: ", token)
	RefreshFCookie(token)
	c.Cookie(&cookie)

	clamas := ValidateToken(token).(map[string]string)
	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	var post models.Post
	//pid := c.ParamNames()
	pidv := c.Params("id")
	models.Dbcon.Where("id = ?", pidv).Find(&post)
	models.Dbcon.Unscoped().Delete(&post)
	datos["mensajeflash"] = "Entrada borrada correctamente"
	datos["alerta"] = template.HTML("alert alert-success")
	c.Render("index", datos, "layouts/main")
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
