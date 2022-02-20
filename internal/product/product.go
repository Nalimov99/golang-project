package product

import "github.com/jmoiron/sqlx"

//List returns all known Products
func List(db *sqlx.DB) ([]Product, error) {
	list := []Product{}

	const q = `
		SELECT product_id, name, quantity, cost, date_created, date_updated
		FROM products
	`
	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}

//Retrieve returns a single Product
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	var prod Product

	const q = `
		SELECT product_id, name, quantity, cost, date_created, date_updated
		FROM products
		WHERE product_id = $1
	`

	if err := db.Get(&prod, q, id); err != nil {
		return nil, err
	}

	return &prod, nil
}
