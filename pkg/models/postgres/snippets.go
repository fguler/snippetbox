package postgres

import (
	"context"
	"errors"

	"github.com/fguler/snippetbox/pkg/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SnippetRepo struct {
	DB *pgxpool.Pool
}

func (s *SnippetRepo) Insert(title string, content string, expires string) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES ($1,$2,NOW(),NOW() + ( $3 || ' days')::interval) RETURNING id`

	// we use QueryRow here instead of Exec because of RETURNING statment at the end of query string above
	row := s.DB.QueryRow(context.Background(), stmt, title, content, expires)

	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SnippetRepo) Get(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() AND id=$1`

	row := s.DB.QueryRow(context.Background(), stmt, id)

	sn := &models.Snippet{}

	err := row.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return sn, nil
}

func (s *SnippetRepo) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() ORDER BY created DESC LIMIT 10`

	// multiple records are expected
	rows, err := s.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {

		sn := &models.Snippet{}

		err := rows.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)

		if err != nil {
			return nil, err
		}
		snippets = append(snippets, sn)
	}

	// we should check again for error here
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil

}

// https://pgdash.io/blog/date-time-interval-postgres.html
