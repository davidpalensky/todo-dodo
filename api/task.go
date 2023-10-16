package api

import (
	"encoding/json"
	"net/http"
	"todo-dodo/action"

	"github.com/gin-gonic/gin"
)

// Api endpoint for creating a batch of tasks for a user
func TaskCreateEnpoint(ctx *gin.Context) {
	args := new([]action.TaskCreateArgs)
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		return
	}
	if err := action.TaskCreateImpl(*args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not update record",
		})
		return
	}
	return
}

// Api endpoint for fetching all tasks for a user
func TaskFetchEnpoint(ctx *gin.Context) {
	args := new(action.TaskFetchArgs)
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid arguements or incorrectly encoded json provided",
		})
		return
	}
	result, err := action.TaskFetch(args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not update record",
		})
		return
	}
	b, _ := json.Marshal(result)
	ctx.JSON(http.StatusOK, string(b))
	return
}
