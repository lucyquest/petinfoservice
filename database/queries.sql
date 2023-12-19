-- name: GetIdempotencyEntry :one
SELECT response
FROM idempotency WHERE user_id = $1 AND key = $2 AND method_path = $3 AND request = $4;

-- name: UpdateIdempotencyEntry :exec
UPDATE idempotency
SET 
  response    = $1
WHERE 
  user_id     = $2 AND
  key         = $3 AND
  method_path = $4 AND
  request     = $5;


-- name: AddIdempotencyEntry :exec
INSERT INTO idempotency (
  user_id,
  key,
  method_path,
  request,
  response
)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
);

-- TODO: these queries that select * get more data than needed

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

-- name: AddPet :one
INSERT INTO pets (name, date_of_birth) VALUES ($1, $2) RETURNING id;
