package main

import (
	"goDock/server/util"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "main")

func main() {
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	logger.Infof("Start web server on %s", "http://localhost:2819")

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

	// Start server
	e.Logger.Fatal(e.Start(":2819"))
}

func getMaxSumHandler(c echo.Context) error {
	data := c.Request().PostFormValue("data")
	max := util.GetMaxSum(data)
	// if err != nil {
	// 	return c.JSON(http.StatusNotFound, err)
	// }
	return c.JSON(http.StatusOK, max)
}
