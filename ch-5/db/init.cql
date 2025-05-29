CREATE KEYSPACE IF NOT EXISTS goanalytics
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE goanalytics;

CREATE TABLE IF NOT EXISTS stats_by_domain (
    domain TEXT PRIMARY KEY,
    count INT
);

CREATE TABLE IF NOT EXISTS stats_by_user (
    user TEXT PRIMARY KEY,
    count INT
);
