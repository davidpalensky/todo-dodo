package api

import (
	"net/http"
	"todo-dodo/logic"

	"github.com/gin-gonic/gin"
)

// Api endpoint for creating a batch of tasks for a user
func TaskCreateBatchEnpoint(ctx *gin.Context) {
	args := new([]logic.TaskCreator)
	if err := ctx.Bind(args); err != nil {
		//response, _ := json.Marshal(APIResponse{Success: false, Data: nil, Err_msg: "Unable to create tasks: Invalid JSON, or incorrect fields provided."})
		ctx.JSON(http.StatusBadRequest,
			APIResponse{Success: false, Data: nil, Err_msg: "Unable to create tasks: Invalid JSON, or incorrect fields provided."})
		return
	}
	if err := logic.TaskCreateBatch(*args); err != nil {
		//log.Printf("err: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError,
			APIResponse{Success: false, Data: nil, Err_msg: "Unable to create tasks: Failed to modify database record."})
		return
	}
	ctx.JSON(http.StatusOK, APIResponse{Success: true, Data: nil, Err_msg: ""})
}

// Api endpoint for fetching all tasks for a user
func TaskFetchAllEnpoint(ctx *gin.Context) {
	args := new(logic.TaskFetchArgs)
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest,
			APIResponse{Success: false, Data: nil, Err_msg: "Unable to fetch tasks: Invalid JSON, or incorrect fields provided."})
		//log.Printf("Error: %s", err.Error())
		return
	}
	result, err := logic.TaskFetchAllWithTags(args)
	if err != nil {
		//log.Printf("Epic debugging: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError,
			APIResponse{Success: false, Data: nil, Err_msg: "Unable to create fetch tasks: Failed to read from database."})
		return
	}
	ctx.JSON(http.StatusOK,
		APIResponse{Success: true, Data: result, Err_msg: ""})
}

// Api endpoint for deleting tags by sending a list of ids [..., ..., ..., ...]
func TaskDeleteBatchEnpoint(ctx *gin.Context) {
	args := new([]uint64) // List of task ids
	if err := ctx.Bind(args); err != nil {
		ctx.JSON(http.StatusBadRequest,
			APIResponse{Success: false, Data: nil, Err_msg: "Unable to delete tasks: Invalid JSON, or incorrect fields provided."})
		//log.Printf("Error: %s", err.Error())
		return
	}
	err := logic.TaskDeleteBatch(*args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			APIResponse{Success: false, Data: nil, Err_msg: "Unable to create delete tasks: Failed to access database."})
		return
	}
	ctx.JSON(http.StatusOK, APIResponse{Success: true, Data: nil, Err_msg: ""})
}

// The endpoint for when a user clicks a checkbox on the frontend.
//func TaskToggleCompletionEndpoint(ctx *gin.Context) {
//	content, _ := io.ReadAll(ctx.Request.Body)
//	log.Printf("Content: %s\n", content)
//}
