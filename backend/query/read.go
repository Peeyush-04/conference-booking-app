package query

import (
	"backend/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// fetch user by user id
func GetUserByID(ctx context.Context, db *pgxpool.Pool, userID uint32) (*models.User, error) {
	// query
	getQuery := `
		SELECT id, first_name, last_name, email, role, created_at
		FROM users
		WHERE id = $1;
	`

	// fetches user details
	var user models.User
	err := db.QueryRow(ctx, getQuery, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// fetches user by email
func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*models.User, error) {
	// query
	getQuery := `
		SELECT id, first_name, last_name, email, role, password_hash, created_at
		FROM users
		WHERE email = $1;
	`

	// fetches from DB
	var user models.User
	err := db.QueryRow(ctx, getQuery, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// fetch conference by conference id
func GetConferenceByID(ctx context.Context, db *pgxpool.Pool, conferenceID uint32) (*models.Conference, error) {
	// query
	getQuery := `
		SELECT id, title, description, location, event_time, total_tickets, available_tickets, organizer_id, status, created_at
		FROM conferences
		WHERE id = $1;
	`

	// fetches conference details
	var conf models.Conference
	err := db.QueryRow(ctx, getQuery, conferenceID).Scan(
		&conf.ID,
		&conf.Title,
		&conf.Description,
		&conf.Location,
		&conf.EventTime,
		&conf.TotalTickets,
		&conf.AvailableTickets,
		&conf.OrganizerID,
		&conf.Status,
		&conf.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// fetchs booking details from booking id
func GetBookingByID(ctx context.Context, db *pgxpool.Pool, bookingID uint32) (*models.Booking, error) {
	// query
	getQuery := `
		SELECT id, user_id, conference_id, tickets_booked, status, booked_at
		FROM bookings
		WHERE id = $1;
	`

	// fetchs booking details
	var booking models.Booking
	err := db.QueryRow(ctx, getQuery, bookingID).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.ConferenceID,
		&booking.TicketsBooked,
		&booking.Status,
		&booking.BookedAt,
	)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// fetches array of tickets from booking id
func GetTicketsByBookingID(ctx context.Context, db *pgxpool.Pool, bookingID uint32) ([]models.Ticket, error) {
	// query
	getQuery := `
	SELECT id, booking_id, ticket_code, issued_at
	FROM tickets
	WHERE booking_id = $1;
	`

	// fetch number of rows
	rows, err := db.Query(ctx, getQuery, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// creates array of tickets from rows to tickets
	tickets := []models.Ticket{}
	for rows.Next() {
		var ticket models.Ticket
		err := rows.Scan(&ticket.ID, &ticket.BookingID, &ticket.TicketCode, &ticket.IssuedAt)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// fetches upcoming conferences
func GetUpcomingConferences(ctx context.Context, db *pgxpool.Pool, days int) ([]models.Conference, error) {
	// time validate
	if days <= 0 || days > 90 {
		return nil, errors.New("invalid period: must be between 1 and 90 days")
	}

	// get query
	getQuery := `
		SELECT id, title, description, location, event_time, total_tickets, available_tickets, organizer_id, status
		FROM conferences
		WHERE event_time BETWEEN NOW() AND NOW() + ($1 * INTERVAL '1 day');
	`

	// fetches available conferences
	rows, err := db.Query(ctx, getQuery, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conferences := []models.Conference{}
	for rows.Next() {
		var conference models.Conference
		err := rows.Scan(
			&conference.ID,
			&conference.Title,
			&conference.Description,
			&conference.Location,
			&conference.EventTime,
			&conference.TotalTickets,
			&conference.AvailableTickets,
			&conference.OrganizerID,
			&conference.Status,
		)
		if err != nil {
			return nil, err
		}
		conferences = append(conferences, conference)
	}

	return conferences, nil
}
