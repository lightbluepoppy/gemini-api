
-- name: GetTodos :many
SELECT id, title, created_time, updated_time
FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (title, created_time, updated_time)
VALUES ($1, $2, $3)
RETURNING id, title, created_time, updated_time;

-- name: GetTodoByID :one
SELECT id, title, created_time, updated_time
FROM todos
WHERE id = $1;

-- name: UpdateTodo :exec
UPDATE todos
SET title = $1, updated_time = $2
WHERE id = $3;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;

-- name: DeleteAllTodos :exec
DELETE FROM todos;