window.onload = init()

function init() {
    const taskTitles = document.getElementsByClassName('openTask');
    console.log(taskTitles);
    console.log("Hello, hello");
    const taskTitles_as_arr = [...taskTitles];
    console.log(taskTitles_as_arr);
    console.log("Hello, hello");
}

// openTask opens a task in the side view of the page
function openTask(event) {
    console.log("openTask is not yet implemented");
    console.log(event);
}


// taskToggleCompleted contacts the API to update completion status
// TODO: Put this into init() not into html directly
function taskToggleCompleted(checkbox) {
    let taskId = parseInt(checkbox.id.match(/[0-9]+/g));
    let update_args = JSON.stringify({
        completed: checkbox.checked,
        task_id: taskId,
        deadline: null,
        tag_ids: []
    });
    //console.log("Updating task completed status.")
    //console.log(update_args)

    let url = "http:\/\/" + self.location.host + "/api/v1/task/update";
    fetch(url, {
        method: "POST",
        body: update_args,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    });
}

// taskDelete conducts an API req to delete that task, and removes it form the DOM
// TODO: Put this into init() not into html directly
function taskDelete(deleteButtonElem) {
    let e = deleteButtonElem;
    let taskId = parseInt(e.id.match(/[0-9]+/g));
    //console.log(e)

    let delete_args = JSON.stringify({
        task_ids: [taskId]
    });
    //console.log("Deleting task.");
    //console.log(delete_args);

    let url = "http:\/\/" + self.location.host + "/api/v1/task/delete";
    //console.log(url);
    fetch(url, {
        method: "POST",
        body: delete_args,
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    });
    //console.log(response);

    // Remove entire task element
    e.parentNode.parentNode.outerHTML = "";
}

// sortTasks sorts the task table.
//function sortTasks(sortBy, Asc) {
//
//}

