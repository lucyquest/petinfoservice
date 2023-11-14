// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Idempotency struct {
	UserID     uuid.UUID
	Key        string
	MethodPath string
	Request    []byte
	Response   []byte
	StatusCode int32
	StatusText sql.NullString
	LockedAt   time.Time
	Step       string
}

type Pet struct {
	ID          uuid.UUID
	Name        string
	DateOfBirth time.Time
}
