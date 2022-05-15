package product

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddSale records a Sale transaction for a single Product.
func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, productID string, now time.Time) (*Sale, error) {
	id, err := strconv.Atoi(productID)
	if err != nil {
		return nil, ErrInvalidId
	}

	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   id,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now.UTC(),
	}

	q := `
	INSERT INTO sales
	(sale_id, product_id, quantity, paid, date_created)
	VALUES
	($1, $2, $3, $4, $5)
	RETURNING *
	`

	var result Sale
	if err := db.QueryRowxContext(ctx, q, s.ID, s.ProductID, s.Quantity, s.Paid, s.DateCreated).StructScan(&result); err != nil {
		return nil, errors.Wrapf(err, "inserting sales: %v", s)
	}

	return &result, nil
}

// ListSales gives all Sales for a Product
func ListSales(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {
	id, err := strconv.Atoi(productID)
	if err != nil {
		return nil, ErrInvalidId
	}

	sales := []Sale{}

	const q = `SELECT * FROM sales WHERE product_id = $1`
	if err := db.SelectContext(ctx, &sales, q, id); err != nil {
		return nil, errors.Wrapf(err, "selecting sales. Product id: %v", id)
	}

	return sales, nil
}
