-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_segments (
    user_id LowCardinality(String),
    segment LowCardinality(String),
    registered_at DateTime CODEC(DoubleDelta, LZ4)
) ENGINE = ReplacingMergeTree() PARTITION BY toYYYYMM(registered_at)
ORDER BY (segment, user_id, registered_at) TTL toDateTime(registered_at) + INTERVAL 2 WEEK DELETE;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_segments;
-- +goose StatementEnd