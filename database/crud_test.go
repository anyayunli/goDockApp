package database

import (
	"database/sql"
	"fmt"
	"goDockApp/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		"localhost", 5432, "api_rw", "godockapp", "")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Panic(err)
	}
	InitPgDb(db)
	result := m.Run()
	db.Exec("TRUNCATE users;")
	os.Exit(result)
}

func cleanDB() {
	DB.Exec("TRUNCATE users;")
}

func Test_CreateUser(t *testing.T) {
	user := model.User{Email: "a1@gmail.com", Password: "132352523"}
	err := CreateUser(&user)
	assert.Nil(t, err)
	dbPwd := getDBPassword(&user)
	assert.NotEqual(t, dbPwd, user.Password)

	err = CreateUser(&user)
	assert.NotNil(t, err)
}

func Test_IsUserExists(t *testing.T) {
	user1 := model.User{Email: "a1@gmail.com", Password: "132352523"}
	user2 := model.User{Email: "a1@gmail.com", Password: "111111111"}
	user3 := model.User{Email: "a3@gmail.com", Password: "132352523"}
	CreateUser(&user1)
	assert.True(t, IsUserExists(&user1))
	assert.False(t, IsUserExists(&user2))
	assert.False(t, IsUserExists(&user3))
}
