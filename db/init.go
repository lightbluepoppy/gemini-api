package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/lightbluepoppy/gemini-api/config"
)

func Connect(config config.Config) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config.DBURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}

func Close(conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		log.Fatalf("Unable to close connection: %v\n", err)
	}
}
