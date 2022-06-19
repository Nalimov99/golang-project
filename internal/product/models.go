package product

import (
	"database/sql"
	"time"
)

// Product is something we sell
type Product struct {
	ID          int          `db:"product_id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Quantity    int          `db:"quantity" json:"quantity"`
	UserID      string       `db:"user_id" json:"user_id"`
	Cost        int          `db:"cost" json:"cost"`
	Sold        int          `db:"sold" json:"sold"`
	Revenue     int          `db:"revenue" json:"revenue"`
	DateCreated sql.NullTime `db:"date_created" json:"date_created"`
	DateUpdated sql.NullTime `db:"date_updated" json:"date_updated"`
}

// NewProduct is what we reqiered from clients to make new product
type NewProduct struct {
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"gte=0"`
	Cost     int    `json:"cost" validate:"gt=0"`
}

// UpdateProduct defines what information can be provided to modify
// an existing Product. All fields are optional so client can send
// just the fields they want changed.
type UpdateProduct struct {
	Name     *string `json:"name"`
	Quantity *int    `json:"quantity" validate:"omitempty,gte=0"`
	Cost     *int    `json:"cost" validate:"omitempty,gt=0"`
}

// Sale reperesents one item of a transaction where some amount of product
// was sold. Quantity is the number of units sold and Paid is the total
// price paid
type Sale struct {
	ID          string    `db:"sale_id" json:"id"`
	ProductID   int       `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Paid        int       `db:"paid" json:"paid"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}

// NewSale is what we required from the clients for recording new transactions
type NewSale struct {
	Quantity int `json:"quantity" validate:"gt=0"`
	Paid     int `json:"paid" validate:"gt=0"`
}
