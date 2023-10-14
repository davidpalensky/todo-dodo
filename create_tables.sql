CREATE TABLE tags (
    tag_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL
);

CREATE TABLE tasks (
    task_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

// Not yet implemented.
CREATE TABLE tag_task_links (
    FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    FOREIGN KEY (tag_id) REFERENCES tags (tag_id)
);
