package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	conf "github.com/lightbluepoppy/gemini-api/config"
	"github.com/lightbluepoppy/gemini-api/db/sqlc"
)

type Server struct {
	config  conf.Config
	router  *gin.Engine
	Queries *sqlc.Queries
	// store  db.Store
}

func NewServer(config conf.Config) *Server {
	var router *gin.Engine
	if config.Environment == "test" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("test environment detected")
		router = gin.New()
	} else {
		router = gin.Default()
	}
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})
	server := &Server{
		config: config,
		router: router,
		// Queries: Queries,
		// store:  store,
	}
	return server
}

func (s *Server) MountHandlers() {
	todos := s.router.Group("/todos")
	todos.GET("/todos", s.GetTodos)
	todos.POST("/todos", s.CreateTodo)
	todos.GET("/todos/:id", s.GetTodoByID)
	todos.PUT("/todos/:id", s.UpdateTodo)
	todos.DELETE("/todos/:id", s.DeleteTodo)
	todos.DELETE("/todos", s.DeleteAllTodos)
}
