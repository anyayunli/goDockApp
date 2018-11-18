package handler

import (
	"goDockApp/database"
	"goDockApp/model"
	"goDockApp/util"
	"net/http"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func TreeHandler(c echo.Context) error {
	tree := &model.TreeSerialized{}
	if err := c.Bind(tree); err != nil {
		return err
	}
	isValid := regexp.MustCompile(`^[0-9#,]+$`).MatchString
	if !isValid(tree.Data) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Tree!"})
	}
	tree.Max = util.GetMaxSum(tree.Data)
	return c.JSON(http.StatusOK, tree)
}

func LoginHandler(c echo.Context) error {
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

func SignUpHandler(c echo.Context) error {
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
