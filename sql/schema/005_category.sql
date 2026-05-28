-- +goose Up
CREATE TABLE categories (
    name TEXT NOT NULL,
    post_id UUID NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE categories;
