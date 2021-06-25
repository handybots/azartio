CREATE DATABASE IF NOT EXISTS azartio;

CREATE TABLE IF NOT EXISTS azartio.logs (
    date        Date,
    time        DateTime,
    level       String,
    message     String,
    event       String,
    user_id     UInt32,
	chat_id		UInt32
) ENGINE = MergeTree(date, (level, event, user_id), 8192);
