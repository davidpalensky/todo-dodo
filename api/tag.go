package api

import (
	"net/http"
	"todo-dodo/logic"

	"github.com/gin-gonic/gin"
)

// Api endpoint for deleting tags by sending a list of ids [..., ..., ..., ...]
func TagDeleteBatchEnpoint(ctx *gin.Context) {
	args := new([]uint64) // List of task ids
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest,
			APIResponse{Success: false, Data: nil, ErrMsg: "Unable to delete tags: Invalid JSON, or incorrect fields provided."})
		//log.Printf("Error: %s", err.Error())
		return
	}
	err := logic.TagDeleteBatch(*args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			APIResponse{Success: false, Data: nil, ErrMsg: "Unable to create fetch tasks: Failed to write to database."})
		return
	}
	ctx.JSON(http.StatusOK, APIResponse{Success: true, Data: nil, ErrMsg: ""})
}
