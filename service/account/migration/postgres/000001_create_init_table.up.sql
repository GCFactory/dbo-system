DROP TABLE IF EXISTS accounts CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE accounts
(
    acc_uuid            UUID PRIMARY KEY                                    DEFAULT uuid_generate_v4(),
    acc_status          numeric(3)                  NOT NULL                DEFAULT 0,
    acc_culc_number     CHAR(20)                    NOT NULL,
    acc_corr_number     CHAR(20)                    NOT NULL,
    acc_bic             CHAR(9)                     NOT NULL,
    acc_cio             CHAR(9)                     NOT NULL,
    acc_money_value     NUMERIC(3)                  NOT NULL ,
    acc_money_amount    NUMERIC(34,4)               NOT NULL                DEFAULT 0
);