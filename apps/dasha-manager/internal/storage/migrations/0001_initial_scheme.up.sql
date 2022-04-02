CREATE TABLE IF NOT EXISTS db_migrations (
    name VARCHAR (256) PRIMARY KEY,
    applied_time TIMESTAMP
);