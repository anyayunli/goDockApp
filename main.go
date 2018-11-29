package main

import (
	"goDockApp/cusmiddleware"
	"goDockApp/database"
	"goDockApp/handler"
	"goDockApp/util"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "main")

// TemplateRenderer is a custom html/template renderer
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
	dbHost := "db"
	dbPort := 5432
	dbUser := "app"
	dbName := "godockapp"
	dbPassword := "password"

	logrus.SetFormatter(util.LogFormatter{})

	db := database.GetDB(dbType, dbHost, dbPort, dbUser, dbName, dbPassword)
	defer db.Close()
	database.InitPgDb(db)

	// Echo instance
	serverPort := ":3344"
	e := echo.New()
	e.Debug = true
	e.HideBanner = true
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	logger.Infof("Start web server on http://localhost%s", serverPort)

	// Middleware
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Non-AuthGroup
	notAuthGroup := e.Group("", cusmiddleware.NotValidateMiddleware)
	notAuthGroup.GET("/signup", handler.RenderSignUpPageHandler)
	notAuthGroup.POST("/signup", handler.SignUpHandler)
	notAuthGroup.GET("/login", handler.RenderLoginPageHandler)
	notAuthGroup.POST("/login", handler.LoginHandler)

	// AuthGroup
	isAuthGroup := e.Group("", cusmiddleware.ValidateMiddleware)
	isAuthGroup.GET("/", handler.RenderIndexPageHandler)
	isAuthGroup.POST("/api/v1/tree", handler.TreeHandler)
	isAuthGroup.POST("/logout", handler.LogOutHandler)

	// Start server
	e.Logger.Fatal(e.Start(serverPort))
}
