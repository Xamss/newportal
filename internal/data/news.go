package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"newsportal.universityproject.net/internal/validator"
	"time"
)

type News struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Abstract  string    `json:"abstract"`
	Tags      []string  `json:"tags,omitempty"`
	Author    string    `json:"author,omitempty"`
	Source    string    `json:"source,omitempty"`
	Time      time.Time `json:"year,omitempty"`
	Version   uuid.UUID `json:"version"`
}

func ValidateNews(v *validator.Validator, news *News) {
	v.Check(news.Title != "", "title", "must be provided")
	v.Check(len(news.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(news.Time.Compare(time.Now()) <= 0, "time", "posted time cannot be in the future")
	v.Check(news.Tags != nil, "tags", "must be provided")
	v.Check(len(news.Tags) >= 1, "tags", "must contain at least 1 tag")
	v.Check(len(news.Tags) <= 5, "tags", "must not contain more than 5 tags")
	v.Check(validator.Unique(news.Tags), "tags", "must not contain duplicate values")
	v.Check(news.Abstract != "", "abstract", "must be provided")
	v.Check(len(news.Abstract) <= 1500, "abstract", "must not be more than 1500 bytes long")
	v.Check(news.Author != "", "author", "must be provided")
	v.Check(len(news.Author) <= 500, "author", "must be less than 500 bytes long. Please, use shorter version if possible")
	v.Check(news.Source != "", "source", "must refer to the original post")
}

type NewsModel struct {
	pool *pgxpool.Pool
}

func (m NewsModel) Insert(news *News) error {

	query := `
		INSERT INTO news (title, abstract, tags, author, source, time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version`

	args := []interface{}{news.Title, news.Abstract, news.Tags, news.Author, news.Source, news.Time}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.pool.QueryRow(ctx, query, args...).Scan(&news.ID, &news.CreatedAt, &news.Version)
}

func (m NewsModel) Get(id int64) (*News, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, abstract, tags, author, source, time, version
	FROM news
	WHERE id = $1`

	var news News

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.pool.QueryRow(ctx, query, id).Scan(
		&news.ID,
		&news.CreatedAt,
		&news.Title,
		&news.Abstract,
		&news.Tags,
		&news.Author,
		&news.Source,
		&news.Time,
		&news.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &news, nil
}

func (m NewsModel) Update(news *News) error {
	query := `
	UPDATE news
	SET title = $1, abstract = $2, tags = $3, author = $4, source = $5, time = $6, version = gen_random_uuid()
	WHERE id = $7 AND version = $8
	RETURNING version`

	args := []interface{}{
		news.Title,
		news.Abstract,
		news.Tags,
		news.Author,
		news.Source,
		news.Time,
		news.ID,
		news.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.pool.QueryRow(ctx, query, args...).Scan(&news.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m NewsModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM news
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m NewsModel) GetAll(title string, tags []string, author string, filters Filters) ([]*News, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, abstract, tags, author, source, time, version
		FROM news
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (tags @> $2 OR $2 = '{}')
		AND (to_tsvector('simple', author) @@ plainto_tsquery('simple', $3) OR $3 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, tags, author, filters.limit(), filters.offset()}

	rows, err := m.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0

	news := []*News{}
	for rows.Next() {
		var post News

		err := rows.Scan(
			&totalRecords,
			&post.ID,
			&post.CreatedAt,
			&post.Title,
			&post.Abstract,
			&post.Tags,
			&post.Author,
			&post.Source,
			&post.Time,
			&post.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		news = append(news, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return news, metadata, nil
}
