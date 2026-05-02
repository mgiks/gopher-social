package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (Post, error)
		DeleteByID(context.Context, int64) error
		UpdateByID(ctx context.Context, id int64, title, content *string, tags []string) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Posts:    PostStore{db: db},
		Users:    UserStore{db: db},
		Comments: CommentStore{db: db},
	}
}
