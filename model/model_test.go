package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_User(t *testing.T) {
	badPwd := "12434456"
	user1 := User{Email: "a1@gmail.com", Password: badPwd}
	user2 := User{Email: "A1@gmAil.com", Password: "1KSkswef3"}
	user1.Prepare()
	user2.Prepare()

	assert.Equal(t, user1.Email, user2.Email)
	assert.NotEqual(t, user1.Password, user2.Password)

	user1.UpdatePassword("xj23iisf3xj23iisf3")
	assert.NotEqual(t, user1.Password, badPwd)
}
