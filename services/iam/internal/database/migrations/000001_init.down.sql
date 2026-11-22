DROP INDEX IF EXISTS idx_email_verifications_user_id;
DROP TABLE IF EXISTS public.email_verifications;

DROP INDEX IF EXISTS idx_refresh_tokens_family_id;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS public.refresh_tokens;

DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;
DROP TABLE IF EXISTS public.password_reset_tokens;

DROP INDEX IF EXISTS idx_users_created_at;
DROP TABLE IF EXISTS public.users;
