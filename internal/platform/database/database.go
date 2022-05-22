package database

import (
	"context"
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
	q.Set("port", "5432")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     c.Host,
		Path:     c.Path,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to
// the database. It returns non-nil error otherwise
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	// Run a simple query to determine connectivity. The db has Ping method
	// but it can be false-positive when it was previously able talk to DB
	// but the DB has since gone away. Running this query forces a round trip
	// to the DB.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
