package product

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors for known failure scenarios
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidId = errors.New("ID provides was not a valid ID")
)

//List returns all known Products
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	list := []Product{}

	const q = `
		SELECT 
			p.product_id, p.name, p.quantity, p.cost,
			COALESCE(SUM(s.quantity), 0) AS sold,
			COALESCE(SUM(s.paid), 0) AS revenue,
			p.date_created, p.date_updated
		FROM products AS p
		LEFT JOIN sales AS s ON s.product_id = p.product_id
		GROUP BY p.product_id
	`
	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return list, nil
}

//Retrieve returns a single Product
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {
	var prod Product

	if _, err := strconv.Atoi(id); err != nil {
		return nil, ErrInvalidId
	}

	const q = `
		SELECT 
		p.product_id, p.name, p.quantity, p.cost,
		COALESCE(SUM(s.quantity), 0) AS sold,
		COALESCE(SUM(s.paid), 0) AS revenue,
		p.date_created, p.date_updated
	FROM products AS p
	LEFT JOIN sales AS s ON s.product_id = p.product_id
	WHERE p.product_id = $1
	GROUP BY p.product_id
	`

	if err := db.GetContext(ctx, &prod, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &prod, nil
}

// Create makes a new product
func Create(ctx context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	var p Product

	const q = `
		INSERT INTO products
		(name, cost, quantity, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`
	if err := db.QueryRowxContext(ctx, q, np.Name, np.Cost, np.Quantity, now.UTC(), now.UTC()).StructScan(&p); err != nil {
		return nil, errors.Wrapf(err, "inserting products: %v \nNow: %v", p, now)
	}

	return &p, nil
}
