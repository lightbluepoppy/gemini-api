package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	conf "github.com/lightbluepoppy/gemini-api/config"
	"github.com/lightbluepoppy/gemini-api/db"
)

type Server struct {
	config conf.Config
	router *gin.Engine
	store  db.Store
}

func NewServer(config conf.Config, store db.Store) *Server {
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
		store:  store,
	}
	return server
}

func (s *Server) MountHandlers() {
	todos := s.router.Group("/todos")
	todos.GET("", s.GetTodos)
	todos.POST("", s.CreateTodo)

	todos.GET("/:id", s.GetTodoByID)
	todos.PUT("/:id", s.UpdateTodo)
	todos.DELETE("/:id", s.DeleteTodo)

	todos.DELETE("", s.DeleteAllTodos)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
