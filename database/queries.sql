-- name: GetPetByID :one
SELECT * FROM pets
WHERE id = sqlc.arg(ids) LIMIT 1;

-- name: GetPetsByIDs :many
SELECT * FROM pets
WHERE id = ANY(sqlc.arg(ids)::uuid[]);

-- name: UpdatePetDateOfBirth :exec
UPDATE pets SET date_of_birth = $1;
