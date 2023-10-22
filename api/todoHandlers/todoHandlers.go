package todoHandlers

import (
	"context"
	"database/dbModules"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}

type TodoHandler struct {
	DB        *dbModules.Queries
	IDCounter int
}

func TodoHandlerFunc() *TodoHandler {
	godotenv.Load(".env")
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	q := dbModules.New(conn)
	return &TodoHandler{
		DB:        q,
		IDCounter: 0,
	}
}

func (t *TodoHandler) GetTodos(ctx *gin.Context) {
	todos, err := t.DB.GetTodos(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, todos)
}
