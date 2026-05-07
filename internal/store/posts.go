package store

import (
	"context"
	"database/sql"
	"errors"

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
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s PostStore) DeleteByID(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		DELETE FROM comments WHERE id = $1
	`
	if err := execQueryWithTx(ctx, tx, query, id); err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}
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

func (s PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, tags = $3, updated_at = NOW()
		WHERE id = $4 AND updated_at = $5
		RETURNING updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
		post.UpdatedAt,
	).Scan(&post.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.tags, u.username, COUNT(c.id) AS comments_count 
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN followers f ON f.followee_id = p.user_id
		WHERE (f.follower_id = $1 OR p.user_id = $1) 
			AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') 
			AND ($5 <@ p.tags OR $5 = '{}')
			AND ($6::timestamptz IS NULL OR $6::timestamptz <= p.created_at)
			AND ($7::timestamptz IS NULL OR p.created_at < $7::timestamptz)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + ` 
		LIMIT $2
		OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
		fq.Limit,
		fq.Offset,
		fq.Search,
		pq.Array(fq.Tags),
		fq.Since,
		fq.Until,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	feed := []PostWithMetadata{}
	for rows.Next() {
		var post PostWithMetadata
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}

	return feed, nil
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
