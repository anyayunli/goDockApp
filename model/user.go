package model

type (
	User struct {
		ID       int    `json:"-"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Token    string `json:"token,omitempty"`
	}
)