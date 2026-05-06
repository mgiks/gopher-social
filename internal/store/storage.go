package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (Post, error)
		DeleteByID(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (User, error)
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followeeID, followerID int64) error
		Unfollow(ctx context.Context, followeeID, followerID int64) error
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Posts:     PostStore{db: db},
		Users:     UserStore{db: db},
		Comments:  CommentStore{db: db},
		Followers: FollowerStore{db: db},
	}
}
