package query

import (
	"backend/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(db *pgxpool.Pool, user models.User, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if user.Role != "customer" && user.Role != "organizer" {
		return errors.New("invalid role: must be either 'customer' or 'organizer'")
	}

	query := `
		INSERT INTO users (first_name, last_name, email, role, hashed_password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	_, err := db.Exec(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Role,
		hashedPassword,
	)

	return err
}
