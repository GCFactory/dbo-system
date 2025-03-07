DROP TABLE IF EXISTS saga CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE saga
(
    saga_uuid       UUID PRIMARY KEY            DEFAULT  uuid_generate_v4(),
    saga_status     numeric(3)                  NOT NULL,
    saga_type       numeric(3)                  NOT NULL,
    saga_data       json                        NOT NULL,
    saga_name       varchar(64)                 NOT NULL DEFAULT 'no name',
    operation_uuid  UUID                        NOT NULL REFERENCES operation(operation_uuid) ON DELETE CASCADE
);