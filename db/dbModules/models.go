// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package dbModules

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Todo struct {
	ID          int32
	Title       string
	CreatedTime pgtype.Timestamp
	UpdatedTime pgtype.Timestamp
}