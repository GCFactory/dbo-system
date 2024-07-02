DROP TABLE IF EXISTS accounts_reserved;

CREATE TABLE accounts_reserved
(
    acc_uuid            UUID references accounts(acc_uuid)                          DEFAULT uuid_generate_v4(),
    reserve_reason      CHAR(512)                           NOT NULL                DEFAULT 'Unknown reason!!!'::char(512)
);