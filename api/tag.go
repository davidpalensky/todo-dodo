package api

import (
	"log"
	"net/http"
	"todo-dodo/action"

	"github.com/gin-gonic/gin"
)

// Api endpoint for deleting tags by sending a list of ids [..., ..., ..., ...]
func TagDeleteBatchEnpoint(ctx *gin.Context) {
	args := new([]uint64) // List of task ids
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		log.Printf("Error: %s", err.Error())
		return
	}
	err := action.TaskDeleteBatch(*args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not update record",
		})
		return
	}
	return
}
