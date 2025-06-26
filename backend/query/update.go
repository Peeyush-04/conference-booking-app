package query

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// universal method
func UpdateUserInfo(ctx context.Context, db *pgxpool.Pool, userID uint32, firstName, lastName, role string) error {
	// validate input
	if strings.TrimSpace(firstName) == "" || strings.TrimSpace(lastName) == "" {
		return errors.New("first name and last name cannot be empty")
	}

	role = strings.ToLower(role)
	if role != "customer" && role != "organizer" {
		return errors.New("invalid role, must be 'customer' or 'organizer'")
	}

	// update query
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, role = $3
		WHERE id = $4
	`

	cmdTag, err := db.Exec(ctx, query, firstName, lastName, role, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no user found with the given ID")
	}

	return nil
}

// only performed by organizer
func UpdateConference(
	ctx context.Context,
	db *pgxpool.Pool,
	conferenceID uint32,
	organizerID uint32,
	title, description, location string,
	eventTime time.Time,
	status string,
) error {
	// Queries
	getQuery := `
		SELECT organizer_id FROM conferences WHERE id = $1
	`
	updateQuery := `
		UPDATE conferences
		SET title = $1,
			description = $2,
			location = $3,
			event_time = $4,
			status = $5
		WHERE id = $6;
	`

	// validate input
	title = strings.TrimSpace(title)
	location = strings.TrimSpace(location)
	status = strings.ToLower(strings.TrimSpace(status))

	if title == "" || location == "" {
		return errors.New("title and location cannot be empty")
	}

	if status != "ongoing" && status != "completed" && status != "cancelled" {
		return errors.New("invalid status")
	}

	// check is conference belongs to the organizer
	var existingOrganizerID uint32
	err := db.QueryRow(ctx, getQuery, conferenceID).Scan(&existingOrganizerID)
	if err != nil {
		return errors.New("conference not found")
	}

	if existingOrganizerID != organizerID {
		return errors.New("you are not authorised to update this conference")
	}

	// update conference
	cmdTag, err := db.Exec(ctx, updateQuery,
		title,
		description,
		location,
		eventTime,
		status,
		conferenceID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no update is made")
	}

	return nil
}

// performed by customer
func UpdateBooking(
	ctx context.Context,
	db *pgxpool.Pool,
	bookingID, userID uint32,
	ticketsBooked uint32,
	status string,
) error {
	// queries
	getQuery := `
		SELECT user_id, booked_at FROM bookings
		WHERE id = $1;
	`
	updateQuery := `
		UPDATE bookings
		SET tickets_booked = $1,
			status = $2
		WHERE id = $3;
	`

	// validate inputs
	status = strings.ToLower(strings.TrimSpace(status))
	if status != "completed" && status != "cancelled" && status != "failed" {
		return errors.New("invalid status value")
	}

	if ticketsBooked <= 0 {
		return errors.New("number of tickets booked should be greater than 0")
	}

	// get booked time and user check
	var bookingAt time.Time
	var bookingUserID uint32
	err := db.QueryRow(ctx, getQuery,
		bookingID,
	).Scan(&bookingUserID, &bookingAt)
	if err != nil {
		return errors.New("booking not found")
	}

	if bookingUserID != userID {
		return errors.New("unauthorised: not your booking")
	}

	// time constraint
	if time.Since(bookingAt) > 4*time.Hour {
		return errors.New("update window expired: can only update within 4 hours after booking")
	}

	cmdTag, err := db.Exec(ctx, updateQuery,
		ticketsBooked,
		status,
		bookingID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no update performed")
	}

	return nil
}

// administration service
func UpdateTicket(
	ctx context.Context,
	db *pgxpool.Pool,
	ticketID, bookingUserID uint32,
	newCode string,
) error {
	// queries
	getQuery := `
		SELECT booking_id, issued_at FROM tickets
		WHERE id = $1;
	`
	checkQuery := `
		SELECT user_id FROM bookings
		WHERE id = $1;
	`
	updateQuery := `
		UPDATE tickets
		SET ticket_code = $1,
			issued_at = $2
		WHERE id = $3;
	`

	newCode = strings.ToLower(strings.TrimSpace(newCode))
	if newCode == "" {
		return errors.New("ticket code cannot be empty")
	}

	// get ticket's issued time and booking id
	var issuedAt time.Time
	var bookingID uint32

	err := db.QueryRow(ctx, getQuery, ticketID).Scan(&bookingID, &issuedAt)
	if err != nil {
		return errors.New("ticket not found")
	}

	// check user ID
	var userID uint32

	err = db.QueryRow(ctx, checkQuery, bookingID).Scan(&userID)
	if err != nil {
		return errors.New("unauthorised: not your ticket")
	}

	// time check
	if time.Since(issuedAt) > 4*time.Hour {
		return errors.New("cannot update ticket after 4 hours of issuance")
	}

	// updating values
	cmdTag, err := db.Exec(ctx, updateQuery,
		newCode,
		issuedAt,
		ticketID,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("no update performed")
	}

	return nil
}
