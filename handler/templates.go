package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// RenderIndexPage ...
func RenderIndexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

// RenderSignUpPage ...
func RenderSignUpPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup.html", map[string]interface{}{})
}

// RenderLoginPage ...
func RenderLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}
