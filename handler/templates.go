package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func RenderIndexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

func renderIndexPageWithData(data map[string]interface{}, c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", data)
}

func RenderSignUpPage(c echo.Context) error {
	return c.Render(http.StatusOK, "signup.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func RenderLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"name": "Dolly!",
	})
}
