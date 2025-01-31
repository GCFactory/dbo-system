DROP TABLE IF EXISTS passport CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE passport
(
    passport_uuid           UUID                PRIMARY KEY     DEFAULT uuid_generate_v4(),
    passport_series         char(4)             NOT NULL,
    passport_number         char(6)             NOT NULL,
    name                    varchar(32)         NOT NULL,
    surname 	            varchar(32)         NOT NULL,
    patronimic              varchar(32)         NOT NULL,
    birth_date              timestamp           NOT NULL,
    birth_location 	        varchar(256)        NOT NULL,
    pick_up_point 	        varchar(256)        NOT NULL,
    authority               char(7)             NOT NULL,
    authority_date          timestamp           NOT NULL,
    registration_adress     varchar(256)        NOT NULL
);

ALTER TABLE passport ADD CONSTRAINT unique_passport_data UNIQUE (passport_series, passport_number, authority);

ALTER TABLE passport ADD CONSTRAINT passport_series_check CHECK (passport_series ~* '^[0-9]{4}$');
ALTER TABLE passport ADD CONSTRAINT passport_number_check CHECK (passport_number ~* '^[0-9]{6}$');
ALTER TABLE passport ADD CONSTRAINT authority_check CHECK (authority ~* '^[0-9]{3}-[0-9]{3}$');