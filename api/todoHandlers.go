package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lightbluepoppy/gemini-api/db/sqlc"
)

type Handler struct {
	Queries *sqlc.Queries
}

func (s *Server) GetTodos(c *gin.Context) {
	todos, err := s.store.GetTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (s *Server) CreateTodo(c *gin.Context) {
	var req sqlc.Todo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := s.store.CreateTodo(c, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (s *Server) GetTodoByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := s.store.GetTodoByID(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (s *Server) UpdateTodo(c *gin.Context) {
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
	todo, err := s.store.UpdateTodo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// todo, err := s.Queries.GetTodoByID(c, int32(id))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, todo)
}

func (s *Server) DeleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.store.DeleteTodo(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := fmt.Sprintf("Deleted TodoID %d successfully.", id)
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (s *Server) DeleteAllTodos(c *gin.Context) {
	err := s.store.DeleteAllTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := "Deleted all Todos successfully."
	c.JSON(http.StatusOK, gin.H{"message": message})
}
