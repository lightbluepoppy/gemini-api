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
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

type TodoHandler struct {
	DB        *dbModules.Queries
	IDCounter int
	Todos     []Todo
}

type HandlerStrategy struct {
	// GetTodos(c *gin.Context) ([]Todo, error)
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

func (t *TodoHandler) GetTodos(c *gin.Context) {
	todos, err := t.DB.GetTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) CreateTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodo.ID = t.IDCounter
	t.IDCounter++
	newTodo.CreatedTime = time.Now()
	newTodo.UpdatedTime = time.Now()
	t.Todos = append(t.Todos, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}
