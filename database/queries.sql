-- name: GetPetByID :one
SELECT * FROM pets
WHERE id = sqlc.arg(ids) LIMIT 1;

-- name: GetPetsByIDs :many
SELECT * FROM pets
WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: UpdatePetName :exec
UPDATE pets SET name = $2 WHERE id = $1;

-- name: UpdatePetDateOfBirth :exec
UPDATE pets SET date_of_birth = $2 WHERE id = $1;
