CREATE DATABASE IF NOT EXISTS azartio;

CREATE TABLE IF NOT EXISTS azartio.logs (
    date        Date,
    time        DateTime,
    level       String,
    message     String,
    event       String,
    user_id     String,
    chat_id	    String
) ENGINE = MergeTree(date, (level, event, user_id), 8192);
