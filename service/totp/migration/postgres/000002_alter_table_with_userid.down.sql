DROP INDEX IF EXISTS totp_codes_user_id_index ;
ALTER TABLE totp_service.totp_codes
    DROP COLUMN IF EXISTS user_id,
    DROP COLUMN IF EXISTS is_active;