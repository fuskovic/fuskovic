package store

import (
	"context"
	"database/sql"
	"net/http"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	query  = "SELECT url FROM urls WHERE path = ? LIMIT 1;"
	create = "create table if not exists urls (id integer PRIMARY KEY,path text NOT NULL,url text NOT NULL)"
	insert = "INSERT INTO urls (path, url) VALUES(?,?);"
)

// Store represents a new store instance
type Store struct{ db *sql.DB }

// New returns a new store instance
func New(ctx context.Context, db *sql.DB, m map[string]string) (*Store, error) {
	store := &Store{db: db}
	err := store.InitTable(ctx, m)
	return store, err
}

// InitTable initializes the urls table.
func (s *Store) InitTable(ctx context.Context, m map[string]string) error {
	if _, err := s.db.ExecContext(ctx, create); err != nil {
		return err
	}
	for path, url := range m {
		if err := s.setURL(ctx, path, url); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) setURL(ctx context.Context, path, url string) error {
	_, err := s.db.ExecContext(ctx, insert, path, url)
	return err
}

func (s *Store) getURL(ctx context.Context, path string) string {
	var url string
	s.db.QueryRowContext(ctx, query, path).Scan(&url)
	return url
}

// Handler redirects requests to the url assigned to the path.
func (s *Store) Handler(ctx context.Context, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := s.getURL(ctx, r.URL.String())
		if url == "" {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	}
}
