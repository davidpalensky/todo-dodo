// taskToggleCompleted contacts the API to update completion status
function taskToggleCompleted(checkbox) {
    let task_id = parseInt(checkbox.id.match(/[0-9]+/g));
    let update_args = JSON.stringify({
        completed: checkbox.checked,
        task_id: task_id,
        deadline: null,
        tag_ids: []
    });
    //console.log("Updating task completed status.")
    //console.log(update_args)

    let url = "http:\/\/" + self.location.host + "/api/v1/task/update"
    fetch(url, {
        method: "POST",
        body: update_args,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    })
}
// taskDelete conducts an API req to delete that task, and removes it form the DOM
function taskDelete(deleteButtonElem) {
    let e = deleteButtonElem
    let task_id = parseInt(e.id.match(/[0-9]+/g));
    //console.log(e)

    let delete_args = JSON.stringify({
        task_ids: [task_id]
    });
    //console.log("Deleting task.");
    //console.log(delete_args);

    let url = "http:\/\/" + self.location.host + "/api/v1/task/delete";
    //console.log(url);
    let response = Promise.resolve(fetch(url, {
        method: "POST",
        body: delete_args,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }));
    //console.log(response);

    // Remove entire task element
    e.parentNode.parentNode.outerHTML = "";
}
