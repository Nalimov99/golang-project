package database

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register postgres database
)

type Config struct {
	User     string `default:"postgres"`
	Password string `default:"1234"`
	Host     string `default:"localhost"`
	Path     string `default:"postgres"`
	SslMode  bool   `default:"false"`
}

func Open(c Config) (*sqlx.DB, error) {
	q := url.Values{}

	q.Set("sslmode", "disable")
	if c.SslMode {
		q.Set("sslmode", "require")
	}
	q.Set("timezone", "utc")
	q.Set("port", "5434")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     c.Host,
		Path:     c.Path,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
