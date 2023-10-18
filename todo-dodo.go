package main

import (
	"net/http"
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
	engine.POST("/api/v2/task/create", api.TaskCreateBatchEnpoint)
	engine.POST("/api/v2/task/fetch", api.TaskFetchAllEnpoint)
	engine.POST("/api/v2/task/delete", api.TaskDeleteEnpoint)

	engine.POST("/api/v2/tag/delete", api.TagDeleteBatchEnpoint)

	// Run server
	engine.Run()
}

func HomeHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"status":  "not implemented",
		"message": "the website has not been created yet",
	})
}
