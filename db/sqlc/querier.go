// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package sqlc

import (
	"context"
)

type Querier interface {
	CreateTodo(ctx context.Context, title string) (*Todo, error)
	DeleteAllTodos(ctx context.Context) error
	DeleteTodo(ctx context.Context, id int32) error
	GetTodoByID(ctx context.Context, id int32) (*Todo, error)
	GetTodos(ctx context.Context) ([]*Todo, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (*Todo, error)
}

var _ Querier = (*Queries)(nil)