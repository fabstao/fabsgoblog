package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Separando roles en MVC
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Dbcon para usarse en todo el paquete
var Dbcon *gorm.DB

// DbConnect llamar para conectar a base de datos
func DbConnect() {
	var err error
	if Dbcon, err = gorm.Open("sqlite3", "fabsgoblog.db"); err != nil {
		panic("Falló la conexión a la base de datos")
	}
	fmt.Println("Conexión a la base de datos exitosa")
}

// MigrarModelos Checar que todos los modelos persistentes sean migrados
func MigrarModelos() {
	Dbcon.AutoMigrate(&User{})
	Dbcon.AutoMigrate(&Post{})
	Dbcon.AutoMigrate(&Comentario{})
	fmt.Println("Modelos migrados ")
}

// User struct para manejo básico de usuarios
type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
}

// Post en el blog, pertenece al usuario y tiene muchos comentarios
type Post struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignkey:UserID"`
	Titulo string
	Texto  string
	Likes  uint
}

// Comentario pertenece a Post y al autor
type Comentario struct {
	gorm.Model
	PostID     uint
	Post       Post `gorm:"foreignkey:PostID"`
	UserID     uint
	User       User `gorm:"foreignkey:UserID"`
	Comentario string
	Likes      uint
}
