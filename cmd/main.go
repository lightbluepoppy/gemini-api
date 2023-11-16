package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lightbluepoppy/gemini-api/api"
	conf "github.com/lightbluepoppy/gemini-api/config"
	"github.com/lightbluepoppy/gemini-api/db"
	"github.com/lightbluepoppy/gemini-api/db/sqlc"
)

func main() {
	var config conf.Config
	config = conf.LoadConfig("dev", "./env")
	config.DBURL = db.DBENV()
	conn := db.Connect(config)
	defer db.Close(conn)
	q := sqlc.New(conn)

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	h := &api.Handler{Queries: q}

	r.GET("/todos", h.GetTodos)
	r.POST("/todos", h.CreateTodo)
	r.GET("/todos/:id", h.GetTodoByID)
	r.PUT("/todos/:id", h.UpdateTodo)
	r.DELETE("/todos/:id", h.DeleteTodo)
	r.DELETE("/todos", h.DeleteAllTodos)

	r.Run(":8080")
}
