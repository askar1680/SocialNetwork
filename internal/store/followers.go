package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserId     int64  `json:"user_id"`
	FollowerId int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (store *FollowerStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := store.db.ExecContext(ctx, query, userID, followerID)
	return err
}

func (store *FollowerStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		DELETE FROM followers WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := store.db.ExecContext(ctx, query, userID, followerID)
	return err
}
