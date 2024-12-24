CREATE TABLE IF NOT EXISTS testsegments (
		user_id LowCardinality(String),
   	    segment LowCardinality(String),
   		registered_at DateTime CODEC(DoubleDelta, LZ4)
	) ENGINE = Memory;