package main

import (
	"goDockApp/database"
	"goDockApp/model"
	"goDockApp/util"
	"html/template"
	"io"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var logger = logrus.WithField("package", "main")

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Postgres
	dbType := "postgres"
	dbHost := "127.0.0.1"
	dbPort := 5432
	dbUser := "api_rw"
	dbName := "godockapp"
	dbPassword := ""

	db := database.GetDB(dbType, dbHost, dbPort, dbUser, dbName, dbPassword)
	defer db.Close()
	database.InitPgDb(db)

	// Echo instance
	serverPort := ":2819"
	e := echo.New()
	e.HideBanner = true
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	logger.Infof("Start web server on http://localhost%s", serverPort)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.POST("/maxsum", getMaxSumHandler)
	e.POST("/login", loginHandler)
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"name": "Dolly!",
		})
	})
	e.POST("/signup", signUpHandler)
	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup.html", map[string]interface{}{
			"name": "Dolly!",
		})
	})

	e.GET("/index", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name": "Dolly!",
		})
	})

	// Start server
	e.Logger.Fatal(e.Start(serverPort))
}

func getMaxSumHandler(c echo.Context) error {
	data := c.Request().PostFormValue("data")
	max := util.GetMaxSum(data)
	return c.JSON(http.StatusOK, max)
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func loginHandler(c echo.Context) error {
	// return c.Render(http.StatusOK, "something.html", map[string]interface{}{
	// 	"name": "Dolly!",
	// })
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	// Validate
	if user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, "invalid email or password")
	}

	var dbPassword string
	err := database.DB.QueryRow("select password from users where email=$1", user.Email).Scan(&dbPassword)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	user.Token, _ = token.SignedString([]byte(Key))
	return c.Redirect(http.StatusSeeOther, "/index")
}

func signUpHandler(c echo.Context) error {
	// Bind
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	// Validate
	if user.Email == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, "invalid email or password")
	}

	// Encrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad password")
	}

	// Save user
	var lastInsertId int
	database.DB.QueryRow("INSERT INTO users(email,password) VALUES($1,$2) returning id;",
		user.Email, hashedPassword).Scan(&lastInsertId)
	return c.Redirect(http.StatusSeeOther, "/index")
}
