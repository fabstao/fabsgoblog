package main

import (
	"html/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/fabstao/fabsgoblog/controllers"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
)

func main() {
	// Initial vars
	templatesDir := "views/templates/*.html"

	// Inicializar capa de datos
	models.DbConnect()
	models.MigrarModelos()
	defer models.Dbcon.Close()

	// Iniciar echo web framework
	e := echo.New()

	// Static assets
	e.Static("/static", "views/assets")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Views/templates
	t := &views.Template{
		Templates: template.Must(template.ParseGlob(templatesDir)),
	}

	e.Renderer = t

	// echo ROUTER (declare HTTP verbs here: GET, PUT, POST, DELETE)
	e.GET("/", controllers.Login)

	e.GET("/cuenta", controllers.Nuevo)

	e.POST("/cuenta", controllers.Crear)

	h := e.Group("/pages")
	h.POST("/index", controllers.Checklogin)
	h.GET("/index", controllers.Inicio)

	api := e.Group("/api")
	api.GET("/api", controllers.Hello)

	// Go echo server!
	e.Logger.Fatal(e.Start(":8019"))
}
