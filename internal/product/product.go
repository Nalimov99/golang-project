package product

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//List returns all known Products
func List(db *sqlx.DB) ([]Product, error) {
	list := []Product{}

	const q = `
		SELECT product_id, name, quantity, cost, date_created, date_updated
		FROM products
	`
	if err := db.Select(&list, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
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
		return nil, errors.Wrap(err, "selecting product")
	}

	return &prod, nil
}

// Create makes a new product
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	var p Product

	const q = `
		INSERT INTO products
		(name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`
	if err := db.QueryRowx(q, np.Name, np.Cost, np.Quantity, now, now).StructScan(&p); err != nil {
		return nil, errors.Wrapf(err, "inserting products: %v \nNow: %v", p, now)
	}

	return &p, nil
}
