package schema

import "github.com/jmoiron/sqlx"

const seeds = `
INSERT INTO products (name, cost, quantity) VALUES
('Book', 12, 3)
`

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
