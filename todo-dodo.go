package main

import (
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

	// JSON-based APIs
	engine.POST("/api/v1/task/create", api.TaskCreateBatchEnpoint)
	engine.POST("/api/v1/task/fetch", api.TaskFetchAllEnpoint)
	engine.POST("/api/v1/task/delete", api.TaskDeleteEnpoint)

	engine.POST("/api/v1/tag/delete", api.TagDeleteBatchEnpoint)

	// Web Pages
	engine.LoadHTMLGlob("./pages/*")

	engine.GET("/", pages.Index)

	// Run server
	engine.Run()
}
