DROP TABLE IF EXISTS notification CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notification
(
    user_uuid           UUID                PRIMARY KEY     DEFAULT uuid_generate_v4(),
    email_usage         bool                NOT NULL        DEFAULT false,
    email               varchar(128)                        DEFAULT ''::varchar(128)
);

ALTER TABLE notification ADD CONSTRAINT unique_email UNIQUE (email);