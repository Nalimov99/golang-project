package schema

import "github.com/jmoiron/sqlx"

const seeds = `
	INSERT INTO products (name, cost, quantity) VALUES
	('Book', 12, 3);

	INSERT INTO users
	(user_id, name, email, roles, password_hash, date_created, date_updated)
	VALUES
	('3888f3a5-7b53-4674-96b7-ef5c1c758ef8', 'Ilia', 'nalimvp@gmail.com', '{ADMIN,USER}', '$2a$10$6SkK1KmGpnQ6uesDqVkLYO/h8LjEKSFM9XkXJ2Djiv3IAejneH6Mu', '2022-01-01 12:00:00', '2022-01-01 12:00:00')
	ON CONFLICT DO NOTHING;
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
