-- name: GetAllTodos :many
SELECT id, title, created_time, updated_time
FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (title, created_time, updated_time)
VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, title, created_time, updated_time;

-- name: GetTodoByID :one
SELECT id, title, created_time, updated_time
FROM todos
WHERE id = $1;

-- name: UpdateTodo :one
UPDATE todos
SET title = $1, updated_time = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING id, title, created_time, updated_time;

-- name: DeleteTodoByID :exec
DELETE FROM todos
WHERE id = $1;

-- name: DeleteAllTodos :exec
DELETE FROM todos;
