-- init.cql: Cassandra schema for stats tracking

-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS goanalytics WITH REPLICATION = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};

-- Use the keyspace
USE goanalytics;

-- Table to store aggregated stats (updated per event)
CREATE TABLE IF NOT EXISTS stats_summary (
  id TEXT PRIMARY KEY,                  -- always 'global' for this app
  total_messages COUNTER,
  bot_count COUNTER,
  non_bot_count COUNTER
);

-- Table to store unique users (one row per user)
CREATE TABLE IF NOT EXISTS unique_users (
  username TEXT PRIMARY KEY
);

-- Table to store counts per server_url
CREATE TABLE IF NOT EXISTS server_counts (
  server_url TEXT PRIMARY KEY,
  count COUNTER
);

CREATE TABLE IF NOT EXISTS goanalytics.users (
    username text PRIMARY KEY,
    password text
);
INSERT INTO goanalytics.users (username, password)
VALUES ('admin', 'password123');
