-- name: AddPasswordResetToken :one
INSERT INTO public.password_reset_tokens(id, user_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: RevokePasswordResetTokensByUserId :exec
UPDATE public.password_reset_tokens
SET revoked_at = NOW()
WHERE user_id = $1 AND used_at IS NULL AND revoked_at IS NULL;

-- name: GetPasswordResetTokenByTokenHash :one
SELECT *
FROM public.password_reset_tokens
WHERE token_hash = $1;

-- name: UsePasswordResetToken :exec
UPDATE public.password_reset_tokens
SET used_at = NOW()
WHERE id = $1;