package main

import (
	"goDockApp/server/database"
	"goDockApp/server/model"
	"goDockApp/server/util"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "main")

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
	logger.Infof("Start web server on http://localhost%d", serverPort)

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
	e.POST("/signup", signUpSumHandler)

	// Start server
	e.Logger.Fatal(e.Start(serverPort))
}

func getMaxSumHandler(c echo.Context) error {
	data := c.Request().PostFormValue("data")
	max := util.GetMaxSum(data)
	return c.JSON(http.StatusOK, max)
}

func loginHandler(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	// Validate
	if user.Email == "" || user.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	}

	rows, err := database.DB.Query("select * from users where email=$1 and password=$2", user.Email, user.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid email or password")
	}
	println(rows)
	return c.JSON(http.StatusCreated, user)
}

func signUpSumHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
