package api

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-dodo/orchestration"

	"github.com/gin-gonic/gin"
)

// Api endpoint for creating a batch of tasks for a user
func TaskCreateBatchEnpoint(ctx *gin.Context) {
	args := new([]orchestration.TaskCreateArgs)
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		return
	}
	if err := orchestration.TaskCreateBatch(*args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could create tasks",
		})
		return
	}
}

// Api endpoint for fetching all tasks for a user
func TaskFetchAllEnpoint(ctx *gin.Context) {
	args := new(orchestration.TaskFetchArgs)
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		log.Printf("Error: %s", err.Error())
		return
	}
	result, err := orchestration.TaskFetchAllWithTags(args)
	if err != nil {
		//log.Printf("Epic debugging: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not fetch tasks",
		})
		return
	}
	b, _ := json.Marshal(result)
	ctx.JSON(http.StatusOK, string(b))
}

// Api endpoint for deleting tags by sending a list of ids [..., ..., ..., ...]
func TaskDeleteEnpoint(ctx *gin.Context) {
	args := new([]uint64) // List of task ids
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		log.Printf("Error: %s", err.Error())
		return
	}
	err := orchestration.TaskDeleteBatch(*args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not delete tasks",
		})
		return
	}
}
