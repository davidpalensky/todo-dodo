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
function taskDelete(domIdOfElemOfEntireTask) {
    let dom_id = domIdOfElemOfEntireTask
    let task_id = parseInt(dom_id.match(/[0-9]+/g));
    // Task element
    const e = document.getElementById(dom_id)

    let delete_args = JSON.stringify({
        task_ids: [task_id]
    });
    console.log("Deleting task.");
    console.log(delete_args);

    let url = "http:\/\/" + self.location.host + "/api/v1/task/delete";
    console.log(url);
    let response = Promise.resolve(fetch(url, {
        method: "POST",
        body: delete_args,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }));
    console.log(response);

    // Remove element, lets hope this is like correct and shit
    dom_id.remove();
}