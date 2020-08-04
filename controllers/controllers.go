package controllers

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gitlab.com/fabstao/fabsgoblog/models"
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
	datos := echo.Map{
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
	return c.Render(http.StatusOK, "index", datos)
}

// Show - Show post form - anyone can read
func Show(c echo.Context) error {
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
	datos := echo.Map{
		"title":        "Fabs Blog",
		"mensajeflash": "",
	}
	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	var post models.Post
	var autor models.User
	//pid := c.ParamNames()
	pidv := c.ParamValues()
	models.Dbcon.Where("id = ?", pidv[0]).Find(&post)
	models.Dbcon.Where("id = ?", post.UserID).Find(&autor)
	datos["titulo"] = post.Titulo
	datos["mitexto"] = post.Texto
	datos["autor"] = autor.Username
	datos["fecha"] = post.UpdatedAt.String()
	datos["pid"] = pidv[0]
	if datos["user"] == autor.Username {
		datos["editar"] = true
	}
	return c.Render(http.StatusOK, "show", datos)
}

// Post - New post form
func Post(c echo.Context) error {
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
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	return c.Render(http.StatusOK, "new", datos)
}

// New post
func New(c echo.Context) error {
	token := "null"
	datos := echo.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if vcookie, err := c.Cookie("frontends1"); err != nil {
		fmt.Println("ERROR reading cookie: ", err)
		return c.Render(http.StatusForbidden, "login", datos)
	} else {
		token = vcookie.Value
		fmt.Println("TOKEN - index: ", token)
		vcookie.Expires = time.Now().Add(15 * time.Minute)
		c.SetCookie(vcookie)
	}

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
		return c.Render(http.StatusForbidden, "login", datos)
	}

	if len(titulo) < 4 || len(texto) < 4 {
		datos["mensajeflash"] = "Título o texto demasiado cortos"
		datos["alerta"] = template.HTML("alert alert-danger")
		return c.Render(http.StatusOK, "new", datos)
	}

	post.Titulo = titulo
	post.Texto = texto
	models.Dbcon.Where("username = ?", datos["user"].(string)).Find(&usuario)
	post.User = usuario
	models.Dbcon.Create(&post)
	datos["mensajeflash"] = "Entrada " + post.Titulo + " creada correctamente"
	datos["titulo"] = titulo
	datos["texto"] = texto
	return c.Render(http.StatusOK, "new", datos)
}

// Edit - New post form
func Edit(c echo.Context) error {
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
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
	}

	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	var post models.Post
	var autor models.User
	pid := c.ParamNames()
	fmt.Println(pid)
	pidv := c.ParamValues()
	models.Dbcon.Where("id = ?", pidv[0]).Find(&post)
	models.Dbcon.Where("id = ?", post.UserID).Find(&autor)
	if autor.Username != datos["user"] {
		fmt.Println("SECURITY ERROR: Unauthorized access")
		return c.Render(http.StatusForbidden, "index", datos)
	}
	datos["titulo"] = post.Titulo
	datos["mitexto"] = post.Texto
	datos["autor"] = autor.Username
	datos["pid"] = pidv[0]
	return c.Render(http.StatusOK, "edit", datos)
}

// Update post
func Update(c echo.Context) error {
	token := "null"
	datos := echo.Map{
		"title":        "Fabs Blog",
		"user":         "",
		"role":         "",
		"mensajeflash": "",
		"alerta":       template.HTML("alert alert-success"),
	}
	if vcookie, err := c.Cookie("frontends1"); err != nil {
		fmt.Println("ERROR reading cookie: ", err)
		return c.Render(http.StatusForbidden, "login", datos)
	} else {
		token = vcookie.Value
		vcookie.Expires = time.Now().Add(15 * time.Minute)
		c.SetCookie(vcookie)
	}

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
		return c.Render(http.StatusForbidden, "login", datos)
	}

	if len(titulo) < 4 || len(texto) < 4 {
		datos["mensajeflash"] = "Título o texto demasiado cortos"
		datos["alerta"] = template.HTML("alert alert-danger")
		return c.Render(http.StatusOK, "new", datos)
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
	return c.Render(http.StatusOK, "edit", datos)
}

// Delete - danger, take care
func Delete(c echo.Context) error {
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
	datos := echo.Map{
		"title":        "Fabs Blog",
		"mensajeflash": "ERROR borrando entrada",
		"alerta":       template.HTML("alert alert-danger"),
	}
	datos["user"] = clamas["User"]
	datos["role"] = clamas["Role"]
	var post models.Post
	//pid := c.ParamNames()
	pidv := c.ParamValues()
	models.Dbcon.Where("id = ?", pidv[0]).Find(&post)
	models.Dbcon.Unscoped().Delete(&post)
	datos["mensajeflash"] = "Entrada borrada correctamente"
	datos["alerta"] = template.HTML("alert alert-success")
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
