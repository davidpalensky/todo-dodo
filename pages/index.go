package pages

import (
	"log"
	"net/http"
	"time"
	"todo-dodo/logic"

	"github.com/gin-gonic/gin"
)

// Golang, what about #[allow(unused_imports)]? We've solved this problem already.
var _ = log.Printf

type TaskView struct {
	Title     string
	Creation  string
	Deadline  string
	Completed bool
	Task_id   uint64
}

// Generate index file and fill in the data
func Index(ctx *gin.Context) {
	data, err := logic.TaskFetchAllWithTags(&logic.TaskFetchArgs{User_id: 1})
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		//log.Println("Could not render html: ", err.Error())
	}
	//log.Println("data: ", data)

	var task_view []TaskView
	for i := 0; i < len(data.Tasks); i++ {
		task_view = append(task_view, TaskView{
			Title:     data.Tasks[i].TaskData.Title,
			Creation:  time.Unix(int64(data.Tasks[i].TaskData.Creation), 0).Format("2 Jan 2006"),
			Deadline:  time.Unix(int64(data.Tasks[i].TaskData.Deadline), 0).Format("2 Jan 2006"),
			Completed: bool(data.Tasks[i].TaskData.Completed),
			Task_id:   data.Tasks[i].TaskData.TaskId,
		})
	}
	//log.Println("view: ", task_view)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Tasks": task_view,
	})
}
