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
	"github.com/lightbluepoppy/gemini-api/src/db/db"
)

// Todo構造体
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}

// type TaskHandler struct {
// 	todos []Todo
// 	idCounter int
// }

var todos []Todo
var idCounter = 1

func run() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	q := db.New(conn)

	todos, err := q.GetTodos(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAuthor failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(todos)
}

func main() {
	router := gin.Default()

	// Todo一覧を取得するエンドポイント
	router.GET("/todos", getTodos)

	// Todo作成エンドポイント
	router.POST("/todos", createTodo)

	// Todo取得エンドポイント
	router.GET("/todos/:id", getTodoByID)

	// Todo更新エンドポイント
	router.PUT("/todos/:id", updateTodo)

	// Todo削除エンドポイント
	router.DELETE("/todos/:id", deleteTodo)

	// Todoをすべて削除するエンドポイント
	router.DELETE("/todos", deleteAllTodos)

	router.Run(":8080")

	run()
}

// Todo一覧を取得するハンドラ
func getTodos(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, todos)
}

// Todoを作成するハンドラ
func createTodo(ctx *gin.Context) {
	var newTodo Todo
	if err := ctx.ShouldBindJSON(&newTodo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodo.ID = idCounter
	idCounter++
	newTodo.CreatedTime = time.Now()
	newTodo.UpdatedTime = time.Now()
	todos = append(todos, newTodo)
	ctx.JSON(http.StatusCreated, newTodo)
}

// Todoを取得するハンドラ
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

// Todoを更新するハンドラ
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

	// Todoの更新
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

// Todoを削除するハンドラ
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
