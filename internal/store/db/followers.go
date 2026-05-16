package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	FolloweeID int64  `json:"followee_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s FollowerStore) Follow(ctx context.Context, followeeID, followerID int64) error {
	query := `
		INSERT INTO followers (followee_id, follower_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followeeID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return nil
}

func (s FollowerStore) Unfollow(ctx context.Context, followeeID, followerID int64) error {
	query := `
		DELETE FROM followers WHERE followee_id = $1 AND follower_id = $2	
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followeeID, followerID)
	return err
}
