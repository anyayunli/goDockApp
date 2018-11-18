package database

import (
	"errors"
	"goDockApp/config"
	"goDockApp/model"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates user
// will return error when user already exists
func CreateUser(user *model.User) error {
	// Encrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return errors.New(config.InValidCredentialsMsg)
	}

	// Save
	var lastInsertId int
	DB.QueryRow("INSERT INTO users(email,password) VALUES($1,$2) returning id;",
		user.Email, hashedPassword).Scan(&lastInsertId)
	if lastInsertId == 0 {
		return errors.New(config.AlreadyExistsMsg)
	}
	return nil
}

func getDBPassword(user *model.User) string {
	var dbPassword string
	DB.QueryRow("select password from users where email=$1", user.Email).Scan(&dbPassword)
	return dbPassword
}

// IsUserExists checks if a user exist
func IsUserExists(user *model.User) bool {
	dbPassword := getDBPassword(user)
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password)); err != nil {
		return false
	}
	return true
}
