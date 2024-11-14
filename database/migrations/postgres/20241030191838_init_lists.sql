-- +goose Up
-- +goose StatementBegin
CREATE TABLE lists (
    user_id uuid NOT NULL,
    title VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lists;
-- +goose StatementEnd
