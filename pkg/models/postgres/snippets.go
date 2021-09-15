package postgres

import (
	"context"

	"github.com/fguler/snippetbox/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SnippetRepo struct {
	DB *pgxpool.Pool
}

func (m *SnippetRepo) Insert(title string, content string, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1,$2,NOW(),NOW() + ( $3 || ' days')::interval) RETURNING id`

	// we use QueryRow here instead of Exec because of RETURNING statment at the end of query string
	row := m.DB.QueryRow(context.Background(), stmt, title, content, expires)

	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetRepo) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (m *SnippetRepo) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

// https://pgdash.io/blog/date-time-interval-postgres.html
