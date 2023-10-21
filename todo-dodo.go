package main

import (
	"os"
	"todo-dodo/api"
	"todo-dodo/db"
	"todo-dodo/pages"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init db
	db.Open()
	defer db.Close()

	// Setup router
	engine := gin.Default()

	// Middleware
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// JSON-based APIs
	api_g := engine.Group("/api/v1")
	api_g.POST("/task/create", api.TaskCreateBatchEnpoint)
	api_g.POST("/task/fetch", api.TaskFetchAllEnpoint)
	api_g.POST("/task/delete", api.TaskDeleteBatchEnpoint)
	api_g.POST("/task/update", api.TaskUpdateEndpoint)

	api_g.POST("/tag/delete", api.TagDeleteBatchEnpoint)

	// Load Templates
	engine.LoadHTMLGlob("./templates/*")

	// Make Templates avaliable
	engine.GET("/", pages.Index)
	if os.Getenv("TODO_DODO_DEV") == "1" {
		engine.GET("/test.html", pages.Test)
	}

	// Js files
	engine.StaticFile("index.js", "./pages/index.js")

	// Run server
	engine.Run()
}
