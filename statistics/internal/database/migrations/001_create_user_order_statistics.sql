-- +goose Up
CREATE TABLE IF NOT EXISTS user_order_statistics (
    user_id UUID PRIMARY KEY,
    order_count INTEGER NOT NULL DEFAULT 0
);
-- +goose Down
DROP TABLE IF EXISTS user_order_statistics;
