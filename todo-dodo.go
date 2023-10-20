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

	// JSON-based APIs
	engine.POST("/api/v1/task/create", api.TaskCreateBatchEnpoint)
	engine.POST("/api/v1/task/fetch", api.TaskFetchAllEnpoint)
	engine.POST("/api/v1/task/delete", api.TaskDeleteBatchEnpoint)

	engine.POST("/api/v1/tag/delete", api.TagDeleteBatchEnpoint)

	// Web Pages
	engine.LoadHTMLGlob("./templates/*")

	engine.GET("/", pages.Index)

	if os.Getenv("TODO_DODO_DEV") == "1" {
		engine.GET("/test.html", pages.Test)
	}

	// Js files for client side behaviour
	engine.StaticFile("/pages/index.js", "./pages/index.js")

	// Run server
	engine.Run()
}
