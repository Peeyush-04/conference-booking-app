package models

// user input data structure
type BookingRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Tickets   uint8  `json:"tickets"`
}

// storing in PostgreSQL data structure
type Booking struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Tickets   uint8
}
