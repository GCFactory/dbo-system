DROP TABLE IF EXISTS accounts_reserved;

CREATE TABLE accounts_reserved
(
    acc_uuid            UUID REFERENCES accounts(acc_uuid) ON DELETE CASCADE        DEFAULT uuid_generate_v4(),
    reserve_reason      VARCHAR(512)                       NOT NULL                 DEFAULT 'Unknown reason!!!'::char(512)
);