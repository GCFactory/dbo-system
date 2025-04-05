CREATE SCHEMA IF NOT EXISTS totp_service;
DROP TABLE IF EXISTS totp_codes CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE totp_codes
(
    totp_id      UUID PRIMARY KEY                       DEFAULT uuid_generate_v4(),
    issuer       VARCHAR(32)                            NOT NULL CHECK ( issuer <> '' ),        -- кто выдаёт
    secret       VARCHAR(32)                            NOT NULL CHECK ( secret <> '' ),        -- сам секрет
    account_name VARCHAR(64)                            NOT NULL CHECK( account_name <> '' ),   -- владелец секрета
    created_at   TIMESTAMP WITH TIME ZONE               NOT NULL DEFAULT NOW(),                 -- когда выдали
    updated_at   TIMESTAMP WITH TIME ZONE               DEFAULT CURRENT_TIMESTAMP               -- когда обновили
);