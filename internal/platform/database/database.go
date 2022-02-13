package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register postgres database
)

func Open() (*sqlx.DB, error) {
	q := url.Values{}

	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")
	q.Set("port", "5434")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "1234"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
