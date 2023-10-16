package main

import (
	"net/http"
	"todo-dodo/api2"
	"todo-dodo/db"

	"github.com/gin-gonic/gin"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	// Init db
	db.DBConnect()
	defer db.DB.Close()

	// Setup router
	engine := gin.Default()
	engine.GET("/", HomeHandler)
	engine.POST("/api/v2/task/create", api2.TaskCreateEnpoint)
	engine.Run()
}

func HomeHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"status":  "not implemented",
		"message": "the website has not been created yet",
	})
}
