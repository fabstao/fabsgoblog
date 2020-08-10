package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/template/django"
	"github.com/joho/godotenv"

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

	// Inicializar capa de datos
	models.DbConnect()
	models.MigrarModelos()

	engine := django.New("./views", ".html")

	// Iniciar echo web framework
	f := fiber.New(&fiber.Settings{
		Views: engine,
	})

	// Static assets
	f.Static("/static", "views/assets")

	// Middleware
	f.Use(middleware.Logger())
	f.Use(middleware.Recover())

	// echo ROUTER (declare HTTP verbs here: GET, PUT, Post, DELETE)
	f.Get("/", controllers.Inicio)
	f.Get("/login", controllers.Login)
	f.Get("/logout", controllers.Logout)
	f.Get("/cuenta", controllers.Nuevo)
	f.Get("/new", controllers.Post)
	f.Get("/show/:id", controllers.Show)
	f.Get("/edit/:id", controllers.Edit)
	f.Get("/delete/:id", controllers.Delete)

	f.Post("/cuenta", controllers.Crear)
	f.Post("/login", controllers.Checklogin)
	f.Post("/new", controllers.New)
	f.Post("/edit", controllers.Update)

	api := f.Group("/api")
	api.Get("/api", controllers.Hello)

	sapi := e.Group("/sapi")
	sapi.Use(middleware.JWT([]byte(controllers.Secret)))

	// Go echo server!
	e.Logger.Fatal(e.Start(":8019"))
	defer models.Dbcon.Close()
}
