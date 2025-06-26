package query

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// universal method
func DeleteUser(ctx context.Context, db *pgxpool.Pool, userID uint32) error {
	deleteQuery := `
		DELETE FROM users WHERE id = $1;
	`

	cmdTag, err := db.Exec(ctx, deleteQuery, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no user deleted")
	}

	return nil
}

// performed by customers
func DeleteBooking(ctx context.Context, db *pgxpool.Pool, bookingID, userID uint32) error {
	// queries
	getQuery := `
		SELECT user_id, booked_at
		FROM bookings
		WHERE id = $1;
	`
	deleteQuery := `
		DELETE FROM bookings WHERE id = $1;
	`

	// validate user id and booked time
	var bookingUserID uint32
	var bookedAt time.Time

	err := db.QueryRow(ctx, getQuery, bookingID).Scan(&bookingUserID, &bookedAt)
	if err != nil {
		return errors.New("booking not found")
	}

	if userID != bookingUserID {
		return errors.New("unauthorized: not your booking")
	}

	if time.Since(bookedAt) > 4*time.Hour {
		return errors.New("deletion window expired: can only delete within 4 hours of booking")
	}

	// delete booking
	cmdTag, err := db.Exec(ctx, deleteQuery, bookingID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("booking is not deleted")
	}

	return nil
}

// only performed by organizer
func DeleteConference(ctx context.Context, db *pgxpool.Pool, conferenceID, organizerID uint32) error {
	// queries
	getQuery := `
		SELECT organizer_id FROM conferences
		WHERE id = $1;
	`
	deleteQuery := `
		DELETE FROM conferences WHERE id = $1;
	`

	// validate organizer
	var existingOrganizerID uint32
	err := db.QueryRow(ctx, getQuery, conferenceID).Scan(&existingOrganizerID)
	if err != nil {
		return err
	}

	if existingOrganizerID != organizerID {
		return errors.New("unauthorized: you are not the correct organizer")
	}

	// deleting conference
	cmdTag, err := db.Exec(ctx, deleteQuery, conferenceID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("conference is not deleted")
	}

	return nil
}
