package product

import "database/sql"

type Product struct {
	ID          int          `db:"product_id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Quantity    int          `db:"quantity" json:"quantity"`
	Cost        int          `db:"cost" json:"cost"`
	DateCreated sql.NullTime `db:"date_created" json:"date_created"`
	DateUpdated sql.NullTime `db:"date_updated" json:"date_updated"`
}
