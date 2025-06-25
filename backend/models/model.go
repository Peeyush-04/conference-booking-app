package models

import "time"

// User Model
type User struct {
	ID        uint32    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Conference Model
type Conference struct {
	ID               uint32    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Location         string    `json:"location"`
	EventTime        time.Time `json:"event_time"`
	TotalTickets     uint32    `json:"total_tickets"`
	AvailableTickets uint32    `json:"available_tickets"`
	OrganizerID      uint32    `json:"organizer_id"`
	CreatedAt        time.Time `json:"created_at"`
}

// Booking Model
type Booking struct {
	ID            uint32    `json:"id"`
	UserID        uint32    `json:"user_id"`
	ConferenceID  uint32    `json:"conference_id"`
	TicketsBooked uint32    `json:"tickets_booked"`
	BookedAt      time.Time `json:"booked_at"`
}

// Ticket Model
type Ticket struct {
	ID         uint32    `json:"id"`
	BookingID  uint32    `json:"booking_id"`
	TicketCode string    `json:"ticket_code"`
	IssuedAt   time.Time `json:"issued_at"`
}
