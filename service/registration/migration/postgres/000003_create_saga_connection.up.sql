DROP TABLE IF EXISTS saga_connection CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE saga_connection
(
    current_saga_uuid   UUID    NOT NULL REFERENCES saga(saga_uuid) ON DELETE CASCADE ,
    next_saga_uuid      UUID    NOT NULL REFERENCES saga(saga_uuid) ON DELETE CASCADE
);

ALTER TABLE saga_connection ADD CONSTRAINT unique_saga_pare UNIQUE (current_saga_uuid, next_saga_uuid);
ALTER TABLE saga_connection ADD CONSTRAINT different_sagas_uuid CHECK ( current_saga_uuid::text <> next_saga_uuid::text );