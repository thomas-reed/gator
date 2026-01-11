-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  user_id UUID NOT NULL,
  CONSTRAINT fk_users
  FOREIGN KEY (user_id)
  REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
