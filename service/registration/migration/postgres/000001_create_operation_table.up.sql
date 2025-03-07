DROP TABLE IF EXISTS saga CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE operation
(
    operation_uuid      UUID PRIMARY KEY            DEFAULT  uuid_generate_v4(),
    operation_name      varchar(128)                NOT NULL DEFAULT 'no name',
    create_time         timestamp                   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_time_update    timestamp                   NOT NULL DEFAULT CURRENT_TIMESTAMP
);