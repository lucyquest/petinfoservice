-- name: GetIdempotencyEntry :one
SELECT response, status_code, status_text, locked_at, step
FROM idempotency WHERE user_id = $1 AND key = $2 AND method_path = $3 AND request = $4;

-- name: UpdateIdempotencyEntry :exec
UPDATE idempotency
SET 
  response    = $5 AND
  status_code = $6 AND
  status_text = $7 AND
  locked_at   = $8 AND
  step        = $9
WHERE 
  user_id     = $1 AND
  key         = $2 AND
  method_path = $3 AND
  request     = $4;


-- name: AddIdempotencyEntry :exec
INSERT INTO idempotency
  user_id,
  key,
  method_path,
  request,
  response,
  status_code,
  status_text,
  locked_at,
  step
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9
);

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
