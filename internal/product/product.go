package product

import "github.com/jmoiron/sqlx"

func List(db *sqlx.DB) ([]Products, error) {
	list := []Products{}

	const q = `SELECT * FROM products`
	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}
