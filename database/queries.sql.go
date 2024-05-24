// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addIdempotencyEntry = `-- name: AddIdempotencyEntry :exec
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
)
`

type AddIdempotencyEntryParams struct {
	UserID     uuid.UUID
	Key        string
	MethodPath string
	Request    []byte
	Response   []byte
}

func (q *Queries) AddIdempotencyEntry(ctx context.Context, arg AddIdempotencyEntryParams) error {
	_, err := q.db.ExecContext(ctx, addIdempotencyEntry,
		arg.UserID,
		arg.Key,
		arg.MethodPath,
		arg.Request,
		arg.Response,
	)
	return err
}

const addPet = `-- name: AddPet :one
INSERT INTO pets (name, date_of_birth) VALUES ($1, $2) RETURNING id
`

type AddPetParams struct {
	Name        string
	DateOfBirth time.Time
}

func (q *Queries) AddPet(ctx context.Context, arg AddPetParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, addPet, arg.Name, arg.DateOfBirth)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deletePet = `-- name: DeletePet :exec
DELETE FROM pets WHERE id = $1
`

func (q *Queries) DeletePet(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePet, id)
	return err
}

const getIdempotencyEntry = `-- name: GetIdempotencyEntry :one
SELECT response
FROM idempotency WHERE user_id = $1 AND key = $2 AND method_path = $3 AND request = $4
`

type GetIdempotencyEntryParams struct {
	UserID     uuid.UUID
	Key        string
	MethodPath string
	Request    []byte
}

func (q *Queries) GetIdempotencyEntry(ctx context.Context, arg GetIdempotencyEntryParams) ([]byte, error) {
	row := q.db.QueryRowContext(ctx, getIdempotencyEntry,
		arg.UserID,
		arg.Key,
		arg.MethodPath,
		arg.Request,
	)
	var response []byte
	err := row.Scan(&response)
	return response, err
}

const getPetByID = `-- name: GetPetByID :one

SELECT row_id, id, name, date_of_birth FROM pets
WHERE id = $1 LIMIT 1
`

// TODO: these queries that select * get more data than needed
func (q *Queries) GetPetByID(ctx context.Context, ids uuid.UUID) (Pet, error) {
	row := q.db.QueryRowContext(ctx, getPetByID, ids)
	var i Pet
	err := row.Scan(
		&i.RowID,
		&i.ID,
		&i.Name,
		&i.DateOfBirth,
	)
	return i, err
}

const getPetsByIDs = `-- name: GetPetsByIDs :many
SELECT row_id, id, name, date_of_birth FROM pets
WHERE id = ANY($1::uuid[])
`

func (q *Queries) GetPetsByIDs(ctx context.Context, ids []uuid.UUID) ([]Pet, error) {
	rows, err := q.db.QueryContext(ctx, getPetsByIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pet
	for rows.Next() {
		var i Pet
		if err := rows.Scan(
			&i.RowID,
			&i.ID,
			&i.Name,
			&i.DateOfBirth,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateIdempotencyEntry = `-- name: UpdateIdempotencyEntry :exec
UPDATE idempotency
SET 
  response    = $1
WHERE 
  user_id     = $2 AND
  key         = $3 AND
  method_path = $4 AND
  request     = $5
`

type UpdateIdempotencyEntryParams struct {
	Response   []byte
	UserID     uuid.UUID
	Key        string
	MethodPath string
	Request    []byte
}

func (q *Queries) UpdateIdempotencyEntry(ctx context.Context, arg UpdateIdempotencyEntryParams) error {
	_, err := q.db.ExecContext(ctx, updateIdempotencyEntry,
		arg.Response,
		arg.UserID,
		arg.Key,
		arg.MethodPath,
		arg.Request,
	)
	return err
}

const updatePetDateOfBirth = `-- name: UpdatePetDateOfBirth :exec
UPDATE pets SET date_of_birth = $2 WHERE id = $1
`

type UpdatePetDateOfBirthParams struct {
	ID          uuid.UUID
	DateOfBirth time.Time
}

func (q *Queries) UpdatePetDateOfBirth(ctx context.Context, arg UpdatePetDateOfBirthParams) error {
	_, err := q.db.ExecContext(ctx, updatePetDateOfBirth, arg.ID, arg.DateOfBirth)
	return err
}

const updatePetName = `-- name: UpdatePetName :exec
UPDATE pets SET name = $2 WHERE id = $1
`

type UpdatePetNameParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) UpdatePetName(ctx context.Context, arg UpdatePetNameParams) error {
	_, err := q.db.ExecContext(ctx, updatePetName, arg.ID, arg.Name)
	return err
}
