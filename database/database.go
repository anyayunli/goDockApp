package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

// InitPgDb initialize database instance
func InitPgDb(db *sql.DB) {
	DB = db
}

// GetDB gets postgresql database or crash.
func GetDB(dbType, dbHost string, dbPort int, dbUser, dbName, dbPassword string) *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := sql.Open(dbType, connString)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"db_type": dbType,
			"db_host": dbHost,
			"db_port": dbPort,
			"db_user": dbUser,
			"db_name": dbName,
		}).Fatal("Database connection failed : ", err)
	}
	err = db.Ping()
	if err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.Info("SUCCESS")
	}
	return db
}
