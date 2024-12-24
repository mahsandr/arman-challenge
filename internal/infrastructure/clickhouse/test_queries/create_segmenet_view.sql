CREATE MATERIALIZED VIEW testsegmentsview TO temptable AS
SELECT segment,count(DISTINCT user_id) as user_count
FROM testsegments
Group by segment;