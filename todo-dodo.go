package main

import (
	"todo-dodo/api"
	"todo-dodo/db"

	"github.com/gin-gonic/gin"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	// Init db
	db.Connect()
	defer db.DB.Close()

	// Setup router
	engine := gin.Default()

	// JSON-based APIs
	engine.POST("/api/v1/task/create", api.TaskCreateBatchEnpoint)
	engine.POST("/api/v1/task/fetch", api.TaskFetchAllEnpoint)
	engine.POST("/api/v1/task/delete", api.TaskDeleteEnpoint)

	engine.POST("/api/v1/tag/delete", api.TagDeleteBatchEnpoint)

	// Web Pages
	engine.StaticFile("/", "./pages/unimplemented.html")

	// Run server
	engine.Run()
}
