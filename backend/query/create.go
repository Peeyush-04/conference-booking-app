package query

import (
	"backend/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx context.Context, db *pgxpool.Pool, user models.User, rawPassword string) (uint32, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO users (first_name, last_name, email, role, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var userID uint32

	err = db.QueryRow(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Role,
		string(hashedPassword),
	).Scan(&userID)

	return userID, err
}

func CreateConference(ctx context.Context, db *pgxpool.Pool, conference *models.Conference) (uint32, error) {
	query := `
		INSERT INTO conferences (
			title, description, location, event_time,
			total_tickets, available_tickets, organizer_id, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
	`

	var conferenceID uint32
	err := db.QueryRow(ctx, query,
		conference.Title,
		conference.Description,
		conference.Location,
		conference.EventTime,
		conference.TotalTickets,
		conference.AvailableTickets,
		conference.OrganizerID,
		conference.Status,
	).Scan(&conferenceID)

	return conferenceID, err
}

func CreateBooking(ctx context.Context, db *pgxpool.Pool, booking models.Booking) (uint32, error) {
	// queries
	getQuery := `
		SELECT available_tickets, status FROM conferences
		WHERE id = $1
	`
	insertQuery := `
		INSERT INTO bookings (user_id, conference_id, tickets_booked, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	updateQuery := `
		UPDATE conferences
		SET available_tickets = available_tickets - $1
		WHERE id = $2
	`

	// transaction phase
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// check conference exists and has enough tickets
	var availableTickets uint32
	var status string

	err = tx.QueryRow(ctx, getQuery, booking.ConferenceID).Scan(&availableTickets, &status)
	if err != nil {
		return 0, err
	}

	if status != "ongoing" {
		return 0, fmt.Errorf("conference is not available for booking")
	}

	if booking.TicketsBooked > availableTickets {
		return 0, fmt.Errorf("not enough tickets available")
	}

	// insert into bookings
	var bookingID uint32
	err = tx.QueryRow(ctx, insertQuery,
		booking.UserID,
		booking.ConferenceID,
		booking.TicketsBooked,
		"completed",
	).Scan(&bookingID)
	if err != nil {
		return 0, err
	}

	// update available tickets
	_, err = tx.Exec(ctx, updateQuery, booking.TicketsBooked, booking.ConferenceID)
	if err != nil {
		return 0, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return bookingID, nil
}
