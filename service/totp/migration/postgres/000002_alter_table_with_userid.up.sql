ALTER TABLE totp_codes
    ADD COLUMN user_id UUID DEFAULT uuid_generate_v4(),
    ADD COLUMN is_active boolean DEFAULT false;

DROP INDEX IF EXISTS totp_codes_user_id_index ;
CREATE INDEX totp_codes_user_id_index
    ON totp_codes (user_id);