SELECT
    tasks.*, tags.*
FROM
    tasks
    INNER JOIN task_tag_links
        ON
            task_tag_links.task_id = tasks.task_id
    INNER JOIN tags
        ON
            task_tag_links.tag_id = tags.tag_id
WHERE
    tasks.user_id = 1;
