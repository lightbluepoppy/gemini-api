package main

import (
	"log"

	"github.com/lightbluepoppy/gemini-api/api"
	"github.com/lightbluepoppy/gemini-api/config"
	"github.com/lightbluepoppy/gemini-api/db"
)

func main() {
	config := config.LoadConfig("dev", "./env")
	config.DBURL = db.DBENV()

	conn := db.Connect(config)
	defer db.Close(conn)
	store := db.NewDatabaseStore(conn)

	server := api.NewServer(
		config,
		store,
	)

	server.MountHandlers()
	err := server.Start("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
