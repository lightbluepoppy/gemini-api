// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package sqlc

import (
	"context"
)

type Querier interface {
	CreateTodo(ctx context.Context, arg CreateTodoParams) (*Todo, error)
	DeleteAllTodos(ctx context.Context) error
	DeleteTodoByID(ctx context.Context, arg DeleteTodoByIDParams) error
	GetAllTodos(ctx context.Context) ([]*Todo, error)
	GetTodoByID(ctx context.Context, arg GetTodoByIDParams) (*Todo, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (*Todo, error)
}

var _ Querier = (*Queries)(nil)
