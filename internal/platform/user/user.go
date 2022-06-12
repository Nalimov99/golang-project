package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Create insert new user into the database
func Create(ctx context.Context, db *sqlx.DB, nu NewUser, now time.Time) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generate password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `
		INSERT INTO users
		(user_id, name, email, roles, password_hash, date_created, date_updated)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = db.ExecContext(
		ctx, q,
		u.ID, u.Name, u.Email, u.Roles, u.PasswordHash, u.DateCreated, u.DateUpdated,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
