package handler

import (
	"fmt"
	"goDockApp/config"
	"goDockApp/core"
	"goDockApp/database"
	"goDockApp/model"
	"net/http"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "handler")

type H map[string]interface{}

func inValidCredentials(email, password string) bool {
	return govalidator.IsEmail(email) && len(password) > 7
}

func isValidTreeData(data string) bool {
	isValid := regexp.MustCompile(`^[0-9#,]+$`).MatchString
	return isValid(data)
}

func errorMsg(err string) map[string]interface{} {
	return map[string]interface{}{"error": err}
}

// RenderSignUpPage ...
func RenderSignUpPage(c echo.Context) error {
	if isUserLogin(c) {
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Render(http.StatusOK, "signup.html", map[string]interface{}{})
}

// RenderLoginPage ...
func RenderLoginPage(c echo.Context) error {
	if isUserLogin(c) {
		return c.Redirect(http.StatusFound, "/")
	}
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{})
}

func isUserLogin(c echo.Context) bool {
	authCookie, err := c.Cookie(config.AuthTokenName)
	if err != nil || authCookie == nil {
		logger.Errorf(err.Error())
		return false
	}
	token, err := jwt.Parse(authCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		logger.Errorf(err.Error())
		return false
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}
	return false
}

// IndexHandler sends user to login page if he failed the jwt authentication
func IndexHandler(c echo.Context) error {
	if !isUserLogin(c) {
		return c.Redirect(http.StatusFound, "/login")
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

// TreeHandler returns max sum of the longest path
// returns error when input tree is not valid
func TreeHandler(c echo.Context) error {
	tree := &model.TreeSerialized{}
	if err := c.Bind(tree); err != nil {
		return c.JSON(http.StatusBadRequest, errorMsg(err.Error()))
	}
	if !isValidTreeData(tree.Data) {
		return c.JSON(http.StatusBadRequest, errorMsg("Invalid Tree!"))
	}
	var err error
	if tree.Max, err = core.GetMaxSum(tree.Data); err != nil {
		return c.JSON(http.StatusBadRequest, errorMsg(err.Error()))
	}
	return c.JSON(http.StatusOK, tree)
}

// LoginHandler handles user's login
// returns error when user does not exist or login info is invalid
func LoginHandler(c echo.Context) error {
	if isUserLogin(c) {
		return c.Redirect(http.StatusFound, "/")
	}
	resp := &model.LoginResponse{}
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		logger.Error(err)
		resp.Error = err.Error()
		return c.Render(http.StatusOK, "login.html", resp)
	}
	user.Prepare()
	resp.Email = user.Email
	resp.Password = user.Password

	// Validate
	if !inValidCredentials(user.Email, user.Password) || !database.IsUserExists(user) {
		resp.Error = config.InValidCredentialsMsg
		logger.Errorf(resp.Error)
		return c.Render(http.StatusOK, "login.html", resp)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and set it in the cookie
	user.Token, _ = token.SignedString([]byte(config.Key))
	user.Password = ""

	cookie := &http.Cookie{
		Name:     config.AuthTokenName,
		Value:    user.Token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   config.SevenDaysInSeconds,
		Expires:  time.Now().Add(time.Duration(config.SevenDaysInSeconds) * time.Second),
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

// SignUpHandler creates a new user when credentials are valid
// returns error if error happens
func SignUpHandler(c echo.Context) error {
	if isUserLogin(c) {
		return c.Redirect(http.StatusFound, "/")
	}
	user := &model.User{}
	resp := &model.LoginResponse{}
	if err := c.Bind(user); err != nil {
		resp.Error = err.Error()
		logger.Error(resp.Error)
		return c.Render(http.StatusOK, "signup.html", resp)
	}
	user.Prepare()
	resp.Email = user.Email
	resp.Password = user.Password

	// Validate
	if !inValidCredentials(user.Email, user.Password) {
		resp.Error = config.InValidCredentialsMsg
		logger.Error(resp.Error)
		return c.Render(http.StatusOK, "signup.html", resp)
	}

	// Save user
	if err := database.CreateUser(user); err != nil {
		resp.Error = err.Error()
		logger.Error(resp.Error)
		return c.Render(http.StatusOK, "signup.html", resp)
	}

	user.Password = ""
	return c.Redirect(http.StatusSeeOther, "/")
}
