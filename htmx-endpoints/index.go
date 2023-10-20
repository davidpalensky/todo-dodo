package htmx

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// The endpoint for when a user clicks a checkbox on the frontend.
func TaskToggleCompletionEndpoint(ctx *gin.Context) {
	content, _ := io.ReadAll(ctx.Request.Body)
	log.Printf("Content: %s\n", content)
}

func TaskDeleteEnpoint(ctx *gin.Context) {
	content, _ := io.ReadAll(ctx.Request.Body)
	log.Printf("Content: %s\n", content)
}
