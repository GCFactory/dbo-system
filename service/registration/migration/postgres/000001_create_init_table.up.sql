CREATE SCHEMA IF NOT EXISTS registration_service;
DROP TABLE IF EXISTS registration_service.saga_table CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE registration_service.saga_table
(
    saga_uuid       UUID PRIMARY KEY            DEFAULT  uuid_generate_v4(),
    saga_status     numeric(3)                  DEFAULT 1,  -- создано,
    saga_type       numeric(3)                  DEFAULT 0,
    saga_name       char(50)                    DEFAULT 'no name',
    list_of_events  json                        DEFAULT '{"list of events":[]}'::json
);