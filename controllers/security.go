package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.com/fabstao/fabsgoblog/models"
	"gitlab.com/fabstao/fabsgoblog/views"
	"golang.org/x/crypto/bcrypt"
)

var cookie http.Cookie

// PToken share token with middleware
var PToken string

// Login Handler
func Login(c echo.Context) error {
	datos := PageRender{
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	return c.Render(http.StatusOK, "login.html", datos)
}

// Logout - destroy session
func Logout(c echo.Context) error {
	datos := PageRender{
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	cookie.Name = "jsessionid"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Domain = ""
	return c.Render(http.StatusOK, "index.html", datos)
}

// Nuevo - controlador para formulario nuevo usuario
func Nuevo(c echo.Context) error {
	datos := PageRender{
		Username: "",
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}

// Crear usuario
func Crear(c echo.Context) error {
	datos := PageRender{
		Username: "",
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
		Alerta:       "success",
	}
	var usuario models.User
	var checausuario models.User
	var rol models.Role
	user := c.FormValue("usuario")
	email := c.FormValue("email")
	if len(user) < 4 || len(email) < 4 {
		datos.MensajeFlash = "Usuario y correo no pueden estar vacíos o demasiado cortos"
		datos.Alerta = "danger"
		return c.Render(http.StatusOK, "crearusuario.html", datos)
	}
	models.Dbcon.Where("username = ?", user).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos.MensajeFlash = "Usuario ya existe: " + user
		datos.Alerta = "danger"
		return c.Render(http.StatusOK, "crearusuario.html", datos)
	}
	models.Dbcon.Where("email = ?", email).Find(&checausuario)
	if len(checausuario.Email) > 0 {
		datos.MensajeFlash = "Dirección de correo-e ya existe: " + email
		datos.Alerta = "danger"
		return c.Render(http.StatusOK, "crearusuario.html", datos)
	}
	models.Dbcon.Where("role = ?", "usuario").Find(&rol)
	usuario.Username = user
	usuario.Email = email
	hashed, _ := bcrypt.GenerateFromPassword([]byte(c.FormValue("password")), bcrypt.DefaultCost)
	usuario.Password = string(hashed)
	usuario.Role = rol

	fmt.Println(usuario)
	models.Dbcon.Create(&usuario)
	datos.MensajeFlash = "Usuario " + usuario.Username + " creado exitosamente"
	return c.Render(http.StatusOK, "crearusuario.html", datos)
}

// Checklogin POST verificar login
func Checklogin(c echo.Context) error {
	var usuario models.User
	email := c.FormValue("email")
	password := c.FormValue("password")
	fmt.Print("Login: ", email)
	models.Dbcon.Where("email = ?", email).Find(&usuario)
	datos := PageRender{
		Header: Titulo{
			views.Comunes.Title,
			"",
			"",
		},
		Title:        views.Comunes.Header,
		Footer:       views.Comunes.Footer,
		MensajeFlash: "",
	}
	fmt.Println(datos)
	if len(usuario.Email) < 1 {
		datos.MensajeFlash = "Correo-e o contraseña incorrectos"
		return c.Render(http.StatusOK, "login.html", datos)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password)); err != nil {
		datos.MensajeFlash = "Correo-e o contraseña incorrectos"
		return c.Render(http.StatusOK, "login.html", datos)
	}
	fmt.Print("Login: Successful")
	datos.Header.User = usuario.Username
	var rol models.Role
	models.Dbcon.Where("id = ?", usuario.RoleID).Find(&rol)
	token := jwt.New(jwt.SigningMethodHS256)
	clamas := token.Claims.(jwt.MapClaims)
	clamas["user"] = usuario.Username
	clamas["rol"] = rol.Role
	clamas["exp"] = time.Now().Add(time.Hour).Unix()
	var err error
	PToken, err = token.SignedString([]byte("el_secreto"))
	if err != nil {
		return err
	}
	cookie.Name = "jsessionid"
	cookie.Value = strings.TrimSpace(PToken)
	cookie.Expires = time.Now().Add(5 * time.Minute)
	cookie.Domain = "localhost"
	c.SetCookie(&cookie)
	fmt.Println(usuario)
	return c.Render(http.StatusOK, "index.html", datos)
}

// ValidateToken to check JWT cookie
func ValidateToken(vtoken string) (interface{}, error) {
	token, err := jwt.Parse(strings.TrimSpace(vtoken), func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		fmt.Println("Validando firma HMAC: ", vtoken)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de firmado inválido: %v", token.Header["alg"])
		}

		var verifyKey []byte
		// verifyKey is a []byte containing your secret, e.g. []byte("my_secret_key")
		return verifyKey, nil
	})
	if err != nil {
		fmt.Println("Error: ", err)
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				fmt.Println("Token expirado")
				return nil, err
			case jwt.ValidationErrorSignatureInvalid:
				fmt.Println("Token con firma inválida")
				return nil, err
			default:
				fmt.Println("Token inválido")
				return nil, err
			}
		default:
			fmt.Println("Token inválido")
			return nil, err
		}
	}

	if clamas, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(clamas["user"], clamas["rol"])
		return clamas, nil
	}

	fmt.Println(err)
	return nil, fmt.Errorf("ERROR: Error desconocido")
}
