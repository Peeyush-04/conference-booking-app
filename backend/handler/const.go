package handler

// user errors: includes auth, role, email, password and others
const (
	userError          string = "User not found"
	invalidUserError   string = "Invalid user ID"
	createUserError    string = "Error creating user"
	notOrganizerError  string = "Unauthorized: Only organizers can create conferences"
	roleError          string = "Invalid error"
	emailPasswordError string = "Invalid Email or Password"
)

// conference errors
const (
	eventTimeError          string = "Invalid event time format"
	createConferenceError   string = "Failed to create conference"
	conferencesFetchError   string = "Error fecthing upcoming conferences: "
	conferenceIDError       string = "Invalid conference ID"
	conferenceNotFoundError string = "Conference not found"
	updateConferenceError   string = "Error updating conference: "
	deleteConferenceError   string = "Failed to delete conference: "
	conferenceAuthError     string = "Unauthorized: Not your conference"
)

// booking error
const (
	bookingIDError     string = "Invalid booking ID"
	bookingError       string = "Booking not found"
	bookingAuthError   string = "Unauthorized: Not your booking"
	bookingAccessError string = "Access denied: not your booking"
)

// ticket error
const (
	fetchTicketsError string = "Failed to fetch tickets"
)

// JSON related errors: includes json, jwt
const (
	requestBodyError   string = "Invalid request body"
	invalidJSONRequest string = "Invalid JSON request"
	generateTokenError string = "Error generating token"
)

// Server related error
const (
	internalServerError string = "Internal server error"
)

// Successful messages
const (
	logoutMessage string = "Logged out successfully. Please delete the token on client side."
)
