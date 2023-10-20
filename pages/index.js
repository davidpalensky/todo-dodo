// This function is meant to be run when the completion checkbox next to a task is clicked.
function taskToggleCompletion(checkbox) {
    let task_id = parseInt(checkbox.id.match(/[0-9]+/g));
    let data = {
        state: checkbox.checked,
        task_id: task_id
    }
    console.log(data)
}