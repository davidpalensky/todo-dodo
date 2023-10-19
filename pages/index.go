package pages

import (
	"log"
	"net/http"
	"todo-dodo/orchestration"

	"github.com/gin-gonic/gin"
)

type TaskView struct {
	Title     string
	Content   string
	Completed bool
}

func Index(ctx *gin.Context) {
	data, err := orchestration.TaskFetchAllWithTags(&orchestration.TaskFetchArgs{User_id: 1})
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", nil)
		//log.Println("Could not render html: ", err.Error())
	}
	log.Println("data: ", data)

	var view []TaskView
	for i := 0; i < len(data.Tasks); i++ {
		view = append(view, TaskView{
			Title:     data.Tasks[i].Task_Data.Title,
			Content:   data.Tasks[i].Task_Data.Content,
			Completed: bool(data.Tasks[i].Task_Data.Completed),
		})
	}

	log.Println("view: ", view)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Tasks": view,
	})
}
