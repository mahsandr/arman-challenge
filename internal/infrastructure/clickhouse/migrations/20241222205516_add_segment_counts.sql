-- +goose Up
-- +goose StatementBegin
-- Create materialized view for real-time counting
CREATE MATERIALIZED VIEW IF NOT EXISTS segment_counts
ENGINE = AggregatingMergeTree()
PARTITION BY tuple()
ORDER BY segment
POPULATE AS
SELECT
    segment,
    count(DISTINCT user_id) as user_count
FROM user_segments
WHERE registered_at >= (now() - INTERVAL 2 WEEK)
GROUP BY segment;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS segment_counts;
-- +goose StatementEnd
