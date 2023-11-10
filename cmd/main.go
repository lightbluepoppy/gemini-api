package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	conf "github.com/lightbluepoppy/gemini-api/config"
	db "github.com/lightbluepoppy/gemini-api/db/dbModules"
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

	config.DBURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	conn, err := pgx.Connect(context.Background(), config.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	q := db.New(conn)

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// Todo一覧を取得するエンドポイント
	router.GET("/todos", func(c *gin.Context) { getTodos(c, q) })

	// Todo作成エンドポイント
	router.POST("/todos", func(c *gin.Context) { createTodo(c, q) })
	// Todo取得エンドポイント
	router.GET("/todos/:id", getTodoByID)

	// Todo更新エンドポイント
	router.PUT("/todos/:id", updateTodo)

	// Todo削除エンドポイント
	router.DELETE("/todos/:id", deleteTodo)

	// Todoをすべて削除するエンドポイント
	router.DELETE("/todos", deleteAllTodos)

	router.Run(":8080")
}

// func getTodos(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, todos)
// }

// func getTodos(q) {
// }

// func createTodo(ctx *gin.Context) {
// var newTodo Todo
// if err := ctx.ShouldBindJSON(&newTodo); err != nil {
// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }

// newTodo.ID = idCounter
// idCounter++
// newTodo.CreatedTime = time.Now()
// newTodo.UpdatedTime = time.Now()
// todos = append(todos, newTodo)
// ctx.JSON(http.StatusCreated, newTodo)
// }

func getTodos(c *gin.Context, q *db.Queries) {
	todos, err := q.GetTodos(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context, q *db.Queries) {
	var req db.Todo
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

func getTodoByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			ctx.JSON(http.StatusOK, todo)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func updateTodo(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTodo Todo
	if err := ctx.ShouldBindJSON(&updatedTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for index := range todos {
		if todos[index].ID == id {
			todos[index].Title = updatedTodo.Title
			todos[index].UpdatedTime = time.Now()
			ctx.JSON(http.StatusOK, todos[index])
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func deleteTodo(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for index, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:index], todos[index+1:]...)
			ctx.JSON(http.StatusOK, todos)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func deleteAllTodos(ctx *gin.Context) {
	todos = todos[:0]
	ctx.JSON(http.StatusOK, todos)
}
