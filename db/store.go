package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lightbluepoppy/gemini-api/db/sqlc"
)

type Store interface {
	sqlc.Querier
	UpdateTodoTx(
		ctx context.Context,
		arg UpdateTodoTxParams,
	) (*UpdateTodoTxResult, error)
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

type UpdateTodoTxParams struct {
	sqlc.UpdateTodoParams
}

type UpdateTodoTxResult struct {
	Todo *sqlc.Todo
}

func (store *DatabaseStore) UpdateTodoTx(
	ctx context.Context,
	arg UpdateTodoTxParams,
) (*UpdateTodoTxResult, error) {
	tx, err := store.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	qtx := store.Queries.WithTx(tx)
	updateTodo, err := qtx.UpdateTodo(ctx, arg.UpdateTodoParams)
	if err != nil {
		return nil, err
	}
	todo, err := qtx.GetTodoByID(ctx, updateTodo.ID)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return &UpdateTodoTxResult{
		Todo: todo,
	}, nil
}

func NullableID(row string, err error) (string, error) {
	if err == nil {
		return row, nil
	}
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return "", err
}
