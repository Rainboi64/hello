-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY first_name;

-- name: CreateUser :one
INSERT INTO users (
  first_name, last_name, email, passhash, salt, phone_number
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  set first_name = $2,
  last_name = $3,
  email = $4,
  passhash = $5,
  salt = $6,
  phone_number = $7
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateLink :one
INSERT INTO links (
    destination, createdBy 
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetLink :one 
SELECT * FROM links 
WHERE id = $1 LIMIT 1;