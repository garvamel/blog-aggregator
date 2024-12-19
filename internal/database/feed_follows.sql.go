// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollows = `-- name: CreateFeedFollows :one
WITH inserted_feed_follows AS (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, user_id, feed_id)
SELECT inserted_feed_follows.id, inserted_feed_follows.created_at, inserted_feed_follows.updated_at, inserted_feed_follows.user_id, inserted_feed_follows.feed_id,
        feeds.name AS feed_name,
        users.name AS user_name
FROM inserted_feed_follows
INNER JOIN feeds
ON inserted_feed_follows.feed_id = feeds.id
INNER JOIN users
ON inserted_feed_follows.user_id = users.id
`

type CreateFeedFollowsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowsRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollows(ctx context.Context, arg CreateFeedFollowsParams) (CreateFeedFollowsRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollows,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const deleteFeedFollowRecord = `-- name: DeleteFeedFollowRecord :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2
`

type DeleteFeedFollowRecordParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) DeleteFeedFollowRecord(ctx context.Context, arg DeleteFeedFollowRecordParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowRecord, arg.UserID, arg.FeedID)
	return err
}

const getFeedFollowsForUserID = `-- name: GetFeedFollowsForUserID :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, feed_id, feeds.id, feeds.created_at, feeds.updated_at, feeds.name, url, feeds.user_id, users.id, users.created_at, users.updated_at, users.name, feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1
`

type GetFeedFollowsForUserIDRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	FeedID      uuid.UUID
	ID_2        uuid.UUID
	CreatedAt_2 time.Time
	UpdatedAt_2 time.Time
	Name        string
	Url         string
	UserID_2    uuid.UUID
	ID_3        uuid.UUID
	CreatedAt_3 time.Time
	UpdatedAt_3 time.Time
	Name_2      string
	FeedName    string
	UserName    string
}

func (q *Queries) GetFeedFollowsForUserID(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsForUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserIDRow
	for rows.Next() {
		var i GetFeedFollowsForUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.Name,
			&i.Url,
			&i.UserID_2,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.Name_2,
			&i.FeedName,
			&i.UserName,
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
