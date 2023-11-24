package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lightbluepoppy/gemini-api/db/sqlc"
)

type Store interface {
	sqlc.Querier
	CreateTodoTx(ctx context.Context, arg sqlc.CreateTodoParams) (*TodoTxResult, error)
	UpdateTodoTx(ctx context.Context, arg sqlc.UpdateTodoParams) (*TodoTxResult, error)
}

type DatabaseStore struct {
	db *pgx.Conn
	*sqlc.Queries
}

func NewDatabaseStore(db *pgx.Conn) Store {
	return &DatabaseStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

type TodoTxResult struct {
	Todo *sqlc.Todo
}

func (store *DatabaseStore) CreateTodoTx(
	ctx context.Context,
	arg sqlc.CreateTodoParams,
) (*TodoTxResult, error) {
	tx, err := store.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	qtx := store.Queries.WithTx(tx)
	todo, err := qtx.CreateTodo(ctx, arg)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return &TodoTxResult{
		todo,
	}, nil
}

func (store *DatabaseStore) UpdateTodoTx(
	ctx context.Context,
	arg sqlc.UpdateTodoParams,
) (*TodoTxResult, error) {
	tx, err := store.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	qtx := store.Queries.WithTx(tx)
	updateTodo, err := qtx.UpdateTodo(ctx, arg)
	if err != nil {
		return nil, err
	}
	var req sqlc.GetTodoByIDParams
	req.ID = updateTodo.ID
	todo, err := qtx.GetTodoByID(ctx, req)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return &TodoTxResult{
		todo,
	}, nil
}

type IDParams interface {
	sqlc.GetTodoByIDParams | sqlc.DeleteTodoByIDParams
}
