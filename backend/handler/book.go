package handler

import (
	"backend/middleware"
	"backend/models"
	"backend/query"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingHandler struct {
	DB *pgxpool.Pool
}

func NewBookingHandler(db *pgxpool.Pool) *BookingHandler {
	return &BookingHandler{DB: db}
}

// routes
func (h *BookingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/booking", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)

		r.With(middleware.RequireRole("customer")).Post("/", (h.CreateBooking))
		r.Get("/{id}", h.GetBooking)
		r.With(middleware.RequireRole("customer")).Put("/{id}", (h.UpdateBooking))
		r.With(middleware.RequireRole("customer")).Delete("/{id}", (h.DeleteBooking))
	})
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	type bookingRequest struct {
		ConferenceID  uint32 `json:"conference_id"`
		TicketsBooked uint32 `json:"tickets_booked"`
	}

	// parse request body
	var req bookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get user ID from context
	userIDVal := r.Context().Value("user_id")
	userID, ok := userIDVal.(uint32)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// create booking record
	booking := models.Booking{
		UserID:        userID,
		ConferenceID:  req.ConferenceID,
		TicketsBooked: req.TicketsBooked,
	}

	bookingID, err := query.CreateBooking(r.Context(), h.DB, booking)
	if err != nil {
		http.Error(w, "Failed to create booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	// generate tickets for the booking
	err = query.GenerateTickets(r.Context(), h.DB, bookingID)
	if err != nil {
		http.Error(w, "Booking made but failed to generate tickets", http.StatusInternalServerError)
		return
	}

	// return response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"booking_id": bookingID,
	})
}

// get booking => only customer or organizer requester
func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	// extract booking id from url
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, bookingIDError, http.StatusBadRequest)
		return
	}

	// fecth booking from database
	booking, err := query.GetBookingByID(r.Context(), h.DB, uint32(id))
	if err != nil {
		http.Error(w, bookingError, http.StatusNotFound)
		return
	}

	// fetch user identity from JWT claims
	userID, ok1 := r.Context().Value("user_id").(uint32)
	role, ok2 := r.Context().Value("role").(string)

	if !ok1 || !ok2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// role based perimissions
	if role == "customer" && booking.UserID != userID {
		http.Error(w, bookingAuthError, http.StatusForbidden)
		return
	}

	if role == "organizer" {
		// get conference to validate user id
		conference, err := query.GetConferenceByID(r.Context(), h.DB, userID)
		if err != nil || conference.OrganizerID != userID {
			http.Error(w, conferenceAuthError, http.StatusForbidden)
			return
		}
	}

	// return booking in json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// update booking => customer
func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	// extract booking id from url
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, bookingIDError, http.StatusBadRequest)
		return
	}

	// fetch booking to validate ownership
	booking, err := query.GetBookingByID(r.Context(), h.DB, uint32(id))
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint32)
	if !ok || booking.UserID != userID {
		http.Error(w, bookingAuthError, http.StatusForbidden)
		return
	}

	// decode request body
	type updateBookingRequest struct {
		Tickets uint32 `json:"tickets_booked"`
		Status  string `json:"status"`
	}

	var req updateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, requestBodyError, http.StatusBadRequest)
		return
	}

	// update
	err = query.UpdateBooking(r.Context(), h.DB, uint32(id), userID, req.Tickets, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Booking updates successfully"))
}

// delete booking => customer
func (h *BookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	// extract booking ID
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, bookingIDError, http.StatusBadRequest)
		return
	}

	// extract id from JWT claims
	userID, ok := r.Context().Value("user_id").(uint32)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// delete booking
	err = query.DeleteBooking(r.Context(), h.DB, uint32(id), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
