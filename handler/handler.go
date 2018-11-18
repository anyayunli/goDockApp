package handler

import (
	"goDockApp/core"
	"goDockApp/database"
	"goDockApp/model"
	"goDockApp/util"
	"net/http"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "handler")

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

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
	resp := &model.LoginResponse{}
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		logrus.Error(err)
		resp.Error = err.Error()
		return c.Render(http.StatusOK, "login.html", resp)
	}
	user.Prepare()
	resp.Email = user.Email
	resp.Password = user.Password

	// Validate
	if !inValidCredentials(user.Email, user.Password) || !database.IsUserExists(user) {
		resp.Error = util.InValidCredentialsMsg
		logger.Errorf(resp.Error)
		return c.Render(http.StatusOK, "login.html", resp)
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

// SignUpHandler creates a new user when credentials are valid
// returns error if error happens
func SignUpHandler(c echo.Context) error {
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
		resp.Error = util.InValidCredentialsMsg
		logger.Error(resp.Error)
		return c.Render(http.StatusOK, "signup.html", resp)
	}

	// Save user
	if err := database.CreateUser(user); err != nil {
		resp.Error = err.Error()
		logger.Error(resp.Error)
		return c.Render(http.StatusOK, "signup.html", resp)
	}

	return c.Redirect(http.StatusSeeOther, "/index")
}
