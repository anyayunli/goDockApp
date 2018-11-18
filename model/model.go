package model

import "strings"

type (
	// User ...
	User struct {
		ID       int    `json:"-"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Token    string `json:"token,omitempty"`
	}
	// TreeSerialized ...
	TreeSerialized struct {
		Data           string   `json:"data"`
		Max            int64    `json:"max"`
		LengthestPaths []string `json:"longest_paths"`
	}
	// LoginResponse ...
	LoginResponse struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Error    string `json:"error"`
	}
)

// Prepare makes user's email case-insensitive
func (user *User) Prepare() {
	user.Email = strings.ToLower(user.Email)
}
