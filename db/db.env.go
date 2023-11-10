package db

import (
	"fmt"
	"os"

	conf "github.com/lightbluepoppy/gemini-api/config"
)

func DBENV() string {
	var (
		config conf.Config
	)

	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "dev" {

		// set up logger for dev
		env = "dev"
		config = conf.LoadConfig(env, "./env")
	}

	config.DBURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	return config.DBURL
}
