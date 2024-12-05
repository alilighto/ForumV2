CREATE TABLE IF NOT EXISTS tag_and_post(
    tag_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
    UNIQUE(tag_id, post_id)
);