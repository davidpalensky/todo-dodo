CREATE TABLE users (
    user_id INTEGER PRIMARY KEY,
    username TEXT NOT NULL
);


CREATE TABLE tags (
    tag_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    color TEXT NOT NULL
);

CREATE TABLE tasks (
    task_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    creation INTEGER DEFAULT (unixepoch('now')),
    deadline INTEGER DEFAULT (unixepoch('now')),
    completed BOOLEAN DEFAULT 0,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE task_tag_links (
    task_id INTEGER,
    tag_id INTEGER,
    FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    FOREIGN KEY (tag_id) REFERENCES tags (tag_id)
);

