// xid reads the dom element id of domElement, and returns the first match
// of subsequent numeric characters.
// 
// Example:
// if domElement.id = "task-6789-completion" then
// xid(domElement) --> 6789
//
// May throw.
function xid(domElement) {
    return parseInt(domElement.id.match(/[0-9]+/g));
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

// taskFilter should be called when the user selects a filter option.
// It will directly manipulate the dom, and uses the 
function taskFilter(filterButtonElem) {
    let e = filterButtonElem;
    console.log("Hello from taskFilter");
    console.log("Tag dom id: " + e.id);
    console.log("checkbox.checked: " + e.checked);
    console.log("\n");
}

// sortTasks sorts the task table.
//function sortTasks(sortBy, Asc) {
//
//}

