package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	engine := gin.Default()
	engine.GET("/", HomeHandler)
	engine.Run()
}

func HomeHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"status":  "not implemented",
		"message": "the website has not been created yet",
	})
}
