package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Todo構造体
type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var todos []Todo
var idCounter = 1

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
			todos[index].Task = updatedTodo.Task
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
