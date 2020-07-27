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
	e.GET("/", controllers.Inicio)
	e.GET("/login", controllers.Login)
	e.GET("/logout", controllers.Logout)
	e.GET("/cuenta", controllers.Nuevo)
	e.POST("/cuenta", controllers.Crear)
	e.POST("/login", controllers.Checklogin)

	//h := e.Group("/admin")
	//h.GET("/index", controllers.Inicio)

	api := e.Group("/api")
	api.GET("/api", controllers.Hello)

	sapi := e.Group("/sapi")
	sapi.Use(middleware.JWT([]byte(controllers.Secret)))

	// Go echo server!
	e.Logger.Fatal(e.Start(":8019"))
	defer models.Dbcon.Close()
}
