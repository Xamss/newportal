package data

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"newsportal.universityproject.net/internal/validator"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	NewsId    int64     `json:"news_id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	Version   uuid.UUID `json:"-"`
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Content != "", "content", "must be provided")
	v.Check(len(comment.Content) <= 1000, "content", "must not be more than 1000 bytes long")
}

type CommentModel struct {
	pool *pgxpool.Pool
}

func (c CommentModel) Insert(comment *Comment) error {
	query := `INSERT INTO comments(user_id, news_id, content)
	VALUES($1, $2, $3)
	RETURNING id, created_at, version`

	args := []interface{}{comment.UserId, comment.NewsId, comment.Content}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.pool.QueryRow(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt, &comment.Version)
	if err != nil {
		return err
	}
	return nil
}

func (c CommentModel) ShowAll(id int64, filters Filters) ([]*Comment, Metadata, error) {
	query := `SELECT count(*) OVER(), comments.id,  comments.created_at, user_id, news_id, content, comments.version 
	FROM comments inner join news on comments.news_id = news.id
	WHERE news.id = $1
	LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{id, filters.limit(), filters.offset()}
	rows, err := c.pool.Query(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return nil, Metadata{}, err
	}

	totalRecords := 0
	comments := []*Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&totalRecords,
			&comment.ID,
			&comment.CreatedAt,
			&comment.UserId,
			&comment.NewsId,
			&comment.Content,
			&comment.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return comments, metadata, nil
}
