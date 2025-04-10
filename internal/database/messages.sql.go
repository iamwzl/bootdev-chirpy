// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: messages.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (body, user_id)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, body, user_id
`

type CreateMessageParams struct {
	Body   string
	UserID uuid.UUID
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage, arg.Body, arg.UserID)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const deleteMessage = `-- name: DeleteMessage :execrows
DELETE FROM messages
WHERE id = $1 AND user_id = $2
`

type DeleteMessageParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteMessage(ctx context.Context, arg DeleteMessageParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteMessage, arg.ID, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getMessage = `-- name: GetMessage :one
SELECT id, created_at, updated_at, body, user_id
FROM messages
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetMessage(ctx context.Context, id uuid.UUID) (Message, error) {
	row := q.db.QueryRowContext(ctx, getMessage, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Body,
		&i.UserID,
	)
	return i, err
}

const getMessages_ByAuthor_CreatedAtASC = `-- name: GetMessages_ByAuthor_CreatedAtASC :many
SELECT id, created_at, updated_at, body, user_id
FROM messages
WHERE user_id = $1
ORDER BY created_at ASC
`

func (q *Queries) GetMessages_ByAuthor_CreatedAtASC(ctx context.Context, userID uuid.UUID) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessages_ByAuthor_CreatedAtASC, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
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

const getMessages_CreatedAtASC = `-- name: GetMessages_CreatedAtASC :many
SELECT id, created_at, updated_at, body, user_id
FROM messages
ORDER BY created_at ASC
`

func (q *Queries) GetMessages_CreatedAtASC(ctx context.Context) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessages_CreatedAtASC)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
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
