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
	{
		Version:     3,
		Description: "Add users",
		Script: `
		CREATE TABLE users (
			user_id UUID,
			name TEXT,
			email TEXT UNIQUE,
			roles TEXT[],
			password_hash TEXT,
			date_created TIMESTAMP,
			date_updated TIMESTAMP,

			PRIMARY KEY (user_id)
		);
		`,
	},
	{
		Version:     4,
		Description: "Add user column to products",
		Script: `
		ALTER TABLE products
		ADD COLUMN user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
		`,
	},
	{
		Version:     5,
		Description: "Add fk user_id to products",
		Script: `
		ALTER TABLE products
		ADD CONSTRAINT FK_user_id
		FOREIGN KEY (user_id) REFERENCES users(user_id);
		`,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
