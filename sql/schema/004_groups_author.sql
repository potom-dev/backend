-- +goose Up
ALTER TABLE groups
ADD COLUMN author_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE groups
DROP COLUMN author_id;
