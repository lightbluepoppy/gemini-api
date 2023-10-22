package main

import (
	"net/http"
	"strconv"
	"time"

	"api/todoHandlers"

	"github.com/gin-gonic/gin"
)

// Todo構造体
type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}

var todos []Todo
var idCounter = 1

func main() {
	h := todoHandlers.TodoHandlerFunc()
	router := gin.Default()

	// Todo一覧を取得するエンドポイント
	router.GET("/todos", h.GetTodos)

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

// func getTodos(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, todos)
// }

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
