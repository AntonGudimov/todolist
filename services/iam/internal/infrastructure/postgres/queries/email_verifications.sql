-- name: AddEmailVerification :one
INSERT INTO public.email_verifications(id, user_id, token, expires_at, used_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetEmailVerificationByToken :one
SELECT *
FROM public.email_verifications ev
WHERE ev.token = $1; 

-- name: UseEmailVerification :exec
UPDATE public.email_verifications
SET used_at = $1
WHERE id = $2;

-- name: CleanupEmailVerifications :many
WITH batch AS (
    SELECT id
    FROM public.email_verifications ev
    WHERE ev.expires_at <= $1 OR ev.used_at IS NOT NULL
    ORDER BY ev.expires_at
    LIMIT $2
    FOR UPDATE SKIP LOCKED
)
DELETE FROM public.email_verifications
WHERE id IN (SELECT id FROM batch)
RETURNING id;