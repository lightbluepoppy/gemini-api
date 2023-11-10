package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	conf "github.com/lightbluepoppy/gemini-api/config"
	"github.com/lightbluepoppy/gemini-api/db"
	"github.com/lightbluepoppy/gemini-api/db/sqlc"
	// "github.com/lightbluepoppy/gemini-api/api/todoHandlers"
)

// Todo構造体
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

var todos []Todo
var idCounter = 1

func main() {
	var config conf.Config
	config = conf.LoadConfig("dev", "./env")
	config.DBURL = db.DBENV()
	conn := db.Connect(config)
	defer db.Close(conn)
	q := sqlc.New(conn)

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Todo一覧を取得するエンドポイント
	router.GET("/todos", func(c *gin.Context) { getTodos(c, q) })

	// Todo作成エンドポイント
	router.POST("/todos", func(c *gin.Context) { createTodo(c, q) })

	// Todo取得エンドポイント
	router.GET("/todos/:id", func(c *gin.Context) { getTodoByID(c, q) })

	// Todo更新エンドポイント
	router.PUT("/todos/:id", func(c *gin.Context) { updateTodo(c, q) })

	// Todo削除エンドポイント
	router.DELETE("/todos/:id", func(c *gin.Context) { deleteTodo(c, q) })

	// Todoをすべて削除するエンドポイント
	router.DELETE("/todos", func(c *gin.Context) { deleteAllTodos(c, q) })

	router.Run(":8080")
}

func getTodos(c *gin.Context, q *sqlc.Queries) {
	todos, err := q.GetTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context, q *sqlc.Queries) {
	var req sqlc.Todo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := q.CreateTodo(c, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func getTodoByID(c *gin.Context, q *sqlc.Queries) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := q.GetTodoByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func updateTodo(c *gin.Context, q *sqlc.Queries) {
	var req sqlc.UpdateTodoParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	req.ID = int32(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = q.UpdateTodo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	todo, err := q.GetTodoByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context, q *sqlc.Queries) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = q.DeleteTodo(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := fmt.Sprintf("Deleted TodoID %d successfully.", id)
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func deleteAllTodos(c *gin.Context, q *sqlc.Queries) {
	err := q.DeleteAllTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := fmt.Sprintf("Deleted all Todos successfully.")
	c.JSON(http.StatusOK, gin.H{"message": message})
}
