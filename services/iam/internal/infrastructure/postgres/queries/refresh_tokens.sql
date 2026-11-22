-- name: AddRefreshToken :one
INSERT INTO public.refresh_tokens(id, user_id, family_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM public.refresh_tokens
WHERE id = $1;

-- name: GetRefreshTokenByHash :one
SELECT *
FROM public.refresh_tokens
WHERE token_hash = $1;

-- name: RevokeRefreshTokenByHash :exec
UPDATE public.refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = $1 AND revoked_at IS NULL;

-- name: RevokeRefreshTokensByFamilyId :exec
UPDATE public.refresh_tokens
SET revoked_at = NOW()
WHERE family_id = $1 AND revoked_at is NULL;

-- name: RevokeRefreshTokensByUserId :exec
UPDATE public.refresh_tokens
SET revoked_at = NOW()
WHERE user_id = $1 AND revoked_at is NULL;