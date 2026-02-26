-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * 
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * 
FROM users
WHERE id = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET 
    hashed_password = $1,
    updated_at = NOW()
WHERE email = $2
RETURNING *;

-- name: UpdateUserPasswordAndEmail :one
UPDATE users
SET 
    hashed_password = $1,
    email = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;
