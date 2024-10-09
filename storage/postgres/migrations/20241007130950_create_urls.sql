-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(255),
    url VARCHAR(2048) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    group_id BIGINT NOT NULL,
    FOREIGN KEY (group_id)
        REFERENCES groups(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
