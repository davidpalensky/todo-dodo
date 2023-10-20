// This function is meant to be run when the completed checkbox next to a task is clicked.
function taskToggleCompleted(checkbox) {
    let task_id = parseInt(checkbox.id.match(/[0-9]+/g));
    let update_args = JSON.stringify({
        completed: checkbox.checked,
        task_id: task_id,
        deadline: null,
        tag_ids: []
    });
    console.log("Updating task completed status.")
    console.log(update_args)

    let url = "http://" + self.location.host + "/api/v1/task/update"
    fetch(url, {
        method: "POST",
        body: update_args,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    })
}