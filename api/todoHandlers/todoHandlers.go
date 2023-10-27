package todoHandlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/lightbluepoppy/gemini-api/db"
	"github.com/lightbluepoppy/gemini-api/db/dbModules"
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
	Todos     []Todo
}

func TodoHandlerFunc() *TodoHandler {
	dbURL := db.DBENV()
	conn, err := pgx.Connect(context.Background(), dbURL)
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

func (t *TodoHandler) CreateTodo(ctx *gin.Context) {
	var newTodo Todo
	if err := ctx.ShouldBindJSON(&newTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodo.ID = t.IDCounter
	t.IDCounter++
	newTodo.CreatedTime = time.Now()
	newTodo.UpdatedTime = time.Now()
	t.Todos = append(t.Todos, newTodo)
	ctx.JSON(http.StatusCreated, newTodo)
}
