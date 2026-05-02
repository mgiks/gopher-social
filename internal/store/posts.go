package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (s PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s PostStore) GetByID(ctx context.Context, id int64) (Post, error) {
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts
		WHERE id = $1
	`

	var post Post
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Post{}, ErrNotFound
		default:
			return Post{}, err
		}
	}

	return post, nil
}

func execQueryWithTx(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count < 1 {
		return ErrNotFound
	}

	return nil
}

func (s PostStore) DeleteByID(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
		DELETE FROM comments WHERE id = $1
	`
	if err := execQueryWithTx(ctx, tx, query, id); err != nil {
		return err
	}

	query = `
		DELETE FROM posts WHERE id = $1
	`
	if err := execQueryWithTx(ctx, tx, query, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s PostStore) UpdateByID(
	ctx context.Context,
	id int64,
	title, content *string,
	tags []string,
) error {
	if title == nil && content == nil && tags == nil {
		return nil
	}

	query := "UPDATE posts SET "
	args := []any{}

	if title != nil {
		args = append(args, *title)
		query += fmt.Sprintf("title = $%d", len(args))
	}
	if content != nil {
		args = append(args, *content)
		if len(args) > 1 {
			query += ", "
		}
		query += fmt.Sprintf("content = $%d", len(args))
	}
	if tags != nil {
		args = append(args, pq.Array(tags))
		if len(args) > 1 {
			query += ", "
		}
		query += fmt.Sprintf("tags = $%d", len(args))
	}

	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)

	args = append(args, id)

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count < 1 {
		return ErrNotFound
	}

	return nil
}
