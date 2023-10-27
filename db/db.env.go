package db

import (
	"fmt"
	"os"
)

func DBENV() string {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDBName := os.Getenv("DB_DBNAME")

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbDBName)

	return dbURL
}
