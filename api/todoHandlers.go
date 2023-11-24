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
	todos, err := s.store.GetAllTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (s *Server) CreateTodo(c *gin.Context) {
	var req sqlc.CreateTodoParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo, err := s.store.CreateTodo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (s *Server) GetTodoByID(c *gin.Context) {
	var req sqlc.GetTodoByIDParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
		// 	// handledReq := handleRequest[sqlc.DeleteTodoParams](c)
		// 	// err := s.store.DeleteTodo(c, *handledReq)
		// 	// if err != nil {
		// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	// 	return
		// 	// }
		// 	// message := fmt.Sprintf("Deleted TodoID %d successfully.", handledReq.ID)

	}
	req.ID = int32(id)
	todo, err := s.store.GetTodoByID(c, req)
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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = int32(id)
	todo, err := s.store.UpdateTodoTx(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}
	c.JSON(http.StatusOK, todo)
}

// func (s *Server) DeleteTodo(c *gin.Context) {
// 	var req sqlc.DeleteTodoByIDParams
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	handledReq, err := handleRequest[sqlc.DeleteTodoByIDParams](c)
// 	if err != nil {
// 		return
// 	}
// 	err = s.store.DeleteTodoByID(c, *handledReq)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	message := fmt.Sprintf("Deleted TodoID %d successfully.", handledReq.ID)
// 	c.JSON(http.StatusOK, gin.H{"message": message})
// }

func (s *Server) DeleteTodo(c *gin.Context) {
	var req sqlc.DeleteTodoByIDParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = int32(id)
	err = s.store.DeleteTodoByID(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	message := fmt.Sprintf("Deleted TodoID %d successfully.", req.ID)
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

// func handleRequest[T db.IDParams](
// 	c *gin.Context,
// ) (*T, error) {
// 	idParam := c.Param("id")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return nil, err
// 	}
// 	req := &T{
// 		ID: int32(id),
// 	}
// 	return req, nil
// }
