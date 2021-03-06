package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	// Separando roles en MVC
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Dbcon para usarse en todo el paquete
var Dbcon *gorm.DB

// DbConnect llamar para conectar a base de datos
func DbConnect() {
	var err error
	dbhost := os.Getenv("DBHOST")
	dbport := os.Getenv("DBPORT")
	dbname := os.Getenv("DBNAME")
	dbuser := os.Getenv("DBUSER")
	dbpasswd := os.Getenv("DBPASSWD")
	dburl := "host=" + dbhost + " port=" + dbport + " user=" + dbuser + " dbname=" + dbname + " password=" + dbpasswd + " sslmode=disable"

	if Dbcon, err = gorm.Open("postgres", dburl); err != nil {
		panic("Falló la conexión a la base de datos")
	}
	fmt.Println("Conexión a la base de datos exitosa")
}

// MigrarModelos Checar que todos los modelos persistentes sean migrados
func MigrarModelos() {
	Dbcon.AutoMigrate(&User{})
	Dbcon.AutoMigrate(&Post{})
	Dbcon.AutoMigrate(&Comentario{})
	Dbcon.AutoMigrate(&Role{})
	var rolea Role
	var roleu Role
	var cuenta uint
	Dbcon.Where("role = ?", "admin").Find(&rolea).Count(&cuenta)
	if cuenta < 1 {
		rolea.Role = "admin"
		Dbcon.Create(&rolea)
	}
	Dbcon.Where("role = ?", "usuario").Find(&roleu).Count(&cuenta)
	if cuenta < 1 {
		roleu.Role = "usuario"
		Dbcon.Create(&roleu)
	}
	var miadmin User
	padmin := os.Getenv("PASSWD_ADMIN")
	hashed, _ := bcrypt.GenerateFromPassword([]byte(padmin), bcrypt.DefaultCost)
	cuenta = 0
	Dbcon.Where("username = ?", "admin").Find(&miadmin).Count(&cuenta)
	if cuenta < 1 {
		miadmin.Username = "admin"
		miadmin.Role = rolea
		miadmin.Password = string(hashed)
		miadmin.Email = "soporte@koalatechie.com"
		Dbcon.Create(&miadmin)
	}

	fmt.Println("Modelos migrados ")
}

// User struct para manejo básico de usuarios
type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	RoleID   uint
	Role     Role `gorm:"foreignkey:RoleID"`
}

// Role struct for RBAC
type Role struct {
	gorm.Model
	Role string
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
