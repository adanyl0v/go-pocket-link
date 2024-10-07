-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    refresh_tok TEXT NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth;
-- +goose StatementEnd
