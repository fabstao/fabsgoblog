package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	jwtware "github.com/gofiber/jwt"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"

	"gitlab.com/fabstao/fabsgoblog/controllers"
	"gitlab.com/fabstao/fabsgoblog/models"
)

func main() {
	// Initial vars
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	controllers.SITEKEY = os.Getenv("SITEKEY")
	controllers.Secret = os.Getenv("FGOSECRET")
	controllers.Cdomain = os.Getenv("CDOMAIN")
	port := os.Getenv("PORT")

	// Inicializar capa de datos
	models.DbConnect()
	models.MigrarModelos()

	engine := html.New("./views", ".html")

	// Iniciar echo web framework
	f := fiber.New(&fiber.Settings{
		Views: engine,
	})

	// Static assets
	f.Static("/static", "views/assets")

	// Middleware
	f.Use(middleware.Logger(middleware.LoggerConfig{
		Format:     "[ ${pid} | ${ip} | ${ips} | ${method} | ${time} | ${path} | ${status} | ${error} | ${header:<Authorization>} ]",
		TimeFormat: "15:04:05",
		TimeZone:   "America/Mexico_City",

		Output: os.Stdout,
	}))
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
	api.Get("/", controllers.AllPosts)
	api.Post("/login", controllers.RESTChecklogin)
	api.Get("/:id", controllers.RESTShow)

	sapi := f.Group("/sapi")
	sapi.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(controllers.Secret),
	}))
	sapi.Put("/", controllers.RESTPost)
	sapi.Get("/:id", controllers.RESTShow)
	sapi.Delete("/:id", controllers.RESTDelete)

	// Go fiber server!
	f.Listen(port)
	defer models.Dbcon.Close()
}
