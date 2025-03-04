DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    user_uuid           UUID                PRIMARY KEY     DEFAULT uuid_generate_v4(),
    passport_uuid 	    UUID                REFERENCES passport(passport_uuid) ON DELETE CASCADE NOT NULL,
    user_inn 	        char(20)            NOT NULL,
    user_accounts       json                NOT NULL,
    user_login          varchar(64)            NOT NULL,
    user_password       char(128)           NOT NULL
);


ALTER TABLE users ADD CONSTRAINT unique_passport UNIQUE (passport_uuid);
ALTER TABLE users ADD CONSTRAINT unique_inn UNIQUE (user_inn);
ALTER TABLE users ADD CONSTRAINT unique_login UNIQUE (user_login);
ALTER TABLE users ADD CONSTRAINT unique_password UNIQUE (user_password);

ALTER TABLE users ADD CONSTRAINT inn_check CHECK (user_inn ~* '^[0-9]{20}$');