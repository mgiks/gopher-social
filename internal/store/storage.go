package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStore(db *sql.DB) Store {
	return Store{
		Posts: PostStore{db: db},
		Users: UserStore{db: db},
	}
}
