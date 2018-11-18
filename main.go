package main

import (
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
	dbHost := "127.0.0.1"
	dbPort := 5432
	dbUser := "api_rw"
	dbName := "godockapp"
	dbPassword := ""

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

	// Routes
	e.POST("/login", handler.LoginHandler)
	e.POST("/signup", handler.SignUpHandler)
	e.POST("/api/v1/tree", handler.TreeHandler)

	// Templates
	e.GET("/index", handler.RenderIndexPage)
	e.GET("/signup", handler.RenderSignUpPage)
	e.GET("/login", handler.RenderLoginPage)

	// Start server
	e.Logger.Fatal(e.Start(serverPort))
}
