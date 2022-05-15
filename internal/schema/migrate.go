package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script: `
CREATE TABLE products (
	product_id SERIAL,
	name TEXT,
	cost INT,
	quantity INT,
	date_created TIMESTAMP NULL,
	date_updated TIMESTAMP NULL,

	PRIMARY KEY (product_id)
);
		`,
	},
	{
		Version:     2,
		Description: "Add sales",
		Script: `
		CREATE TABLE sales (
			sale_id UUID,
			product_id int,
			quantity int,
			paid int,
			date_created TIMESTAMP,
			
			PRIMARY KEY (sale_id),
			FOREIGN KEY (product_id) REFERENCES products (product_id) ON DELETE CASCADE
		);
		`,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
