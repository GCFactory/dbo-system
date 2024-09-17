DROP TABLE IF EXISTS saga CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE saga
(
    saga_uuid       UUID PRIMARY KEY            DEFAULT  uuid_generate_v4(),
    saga_status     numeric(3)                  NOT NULL,
    saga_type       numeric(3)                  NOT NULL,
    saga_name       varchar(64)                    NOT NULL DEFAULT 'no name',
    list_of_events  json                        NOT NULL DEFAULT '{"list of events":[]}'::json
);