DROP TABLE IF EXISTS event CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE event
(
    event_uuid              UUID PRIMARY KEY            DEFAULT  uuid_generate_v4(),
    saga_uuid               UUID                        NOT NULL REFERENCES saga(saga_uuid) ON DELETE CASCADE,
    event_status            numeric(3)                  NOT NULL,
    event_name              varchar(64)                 NOT NULL DEFAULT 'No name'::varchar,
    event_is_roll_back      bool                        NOT NULL default false,
    event_result            json                        NOT NULL,
    event_rollback_uuid     UUID                        default null
);