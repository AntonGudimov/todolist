-- name: AddUser :one
INSERT INTO public.users(id, name, email, password, is_verified)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: VerifyUser :exec
UPDATE public.users
SET is_verified = TRUE
WHERE id = $1;

-- name: CountUsersByEmail :one
SELECT COUNT(*)
FROM public.users u
WHERE u.email = $1;

-- name: GetUserByEmail :one
SELECT *
FROM public.users u
where u.email = $1;

-- name: GetUserById :one
SELECT *
FROM public.users u
WHERE u.id = $1;


-- name: UpdatePassword :exec
UPDATE public.users
SET password = $1
WHERE id = $2;