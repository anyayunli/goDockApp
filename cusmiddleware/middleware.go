package cusmiddleware

import (
	"fmt"
	"goDockApp/config"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func isUserLogin(c echo.Context) bool {
	authCookie, err := c.Cookie(config.AuthTokenName)
	if err != nil || authCookie == nil {
		return false
	}
	token, err := jwt.Parse(authCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(config.InvalidAuthTokenMsg)
		}
		return []byte(config.Key), nil
	})
	if err != nil {
		logrus.Errorf(err.Error())
		return false
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}
	return false
}

// NotValidateMiddleware redirects user to / if sign in
func NotValidateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if isUserLogin(c) {
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}

// ValidateMiddleware redirects user to /login if not sign in
func ValidateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isUserLogin(c) {
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}
