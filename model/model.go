package model

type (
	User struct {
		ID       int    `json:"-"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Token    string `json:"token,omitempty"`
	}
	TreeSerialized struct {
		Data           string   `json:"data"`
		Max            int      `json:"max"`
		LengthestPaths []string `json:"longest_paths"`
	}
)
