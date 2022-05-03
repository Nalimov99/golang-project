package databasetest

import (
	"garagesale/internal/platform/database"
	"garagesale/internal/schema"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

// Setup creates a db connection inside docker container
// It returns the db to use and function to call at the end of the tests
func Setup(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := startContainer(t)

	db, err := database.Open(database.Config{
		User:     "postgres",
		Password: "postgres",
		Host:     c.Ports,
		Path:     "postgres",
		SslMode:  false,
	})
	if err != nil {
		t.Fatalf("db connection failed: %v", err)
	}

	t.Log("waiting for db ready")

	var pingError error
	// wait for db to be ready
	for attempts := 1; attempts < 20; attempts++ {
		pingError = db.Ping()

		if pingError == nil {
			break
		}

		time.Sleep(time.Second)
	}

	if pingError != nil {
		stopContainer(t, c)
		t.Fatalf("db is not ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		t.Fatalf("could not migrate: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with database
	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c)
	}

	return db, teardown
}
