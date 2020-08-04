package main

import (
	"log"
	"os"

	"github.com/foolin/goview/supports/echoview"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.com/fabstao/fabsgoblog/controllers"
	"gitlab.com/fabstao/fabsgoblog/models"
)

func main() {
	// Initial vars
	err := godotenv.Load()
	if err != nil {
		//log.Fatal("Error loading .env file")
		log.Println("Error loading .env file")
	}
	controllers.SITEKEY = os.Getenv("SITEKEY")
	controllers.Secret = os.Getenv("FGOSECRET")
	controllers.Cdomain = os.Getenv("CDOMAIN")

	//templatesDir := "views/templates/*.html"

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
	/*t := &views.Template{
		Templates: template.Must(template.ParseGlob(templatesDir)),
	}

	e.Renderer = t */

	e.Renderer = echoview.Default()

	// echo ROUTER (declare HTTP verbs here: GET, PUT, POST, DELETE)
	e.GET("/", controllers.Inicio)
	e.GET("/login", controllers.Login)
	e.GET("/logout", controllers.Logout)
	e.GET("/cuenta", controllers.Nuevo)
	e.GET("/new", controllers.Post)
	e.GET("/show/:id", controllers.Show)
	e.GET("/edit/:id", controllers.Edit)
	e.GET("/delete/:id", controllers.Delete)

	e.POST("/cuenta", controllers.Crear)
	e.POST("/login", controllers.Checklogin)
	e.POST("/new", controllers.New)
	e.POST("/edit", controllers.Update)

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
