package handler

import (
	"backend/middleware"
	"backend/query"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketHandler struct {
	DB *pgxpool.Pool
}

func NewTicketHandler(db *pgxpool.Pool) *TicketHandler {
	return &TicketHandler{DB: db}
}

func (h *TicketHandler) RegisterRoutes(r chi.Router) {
	r.Route("/ticket", func(r chi.Router) {
		r.Get("/booking/{bookingID}", h.GetTicketsByBookingID)
	})
}

// get tickets by booking ID
func (h *TicketHandler) GetTicketsByBookingID(w http.ResponseWriter, r *http.Request) {
	// get booking id from url
	bookingIDString := chi.URLParam(r, "bookingID")
	bookingID, err := strconv.ParseUint(bookingIDString, 10, 32)
	if err != nil {
		http.Error(w, bookingIDError, http.StatusBadRequest)
		return
	}

	// get user id from JWT claims
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint32)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// check booking ownership
	booking, err := query.GetBookingByID(r.Context(), h.DB, uint32(bookingID))
	if err != nil || booking.UserID != userID {
		http.Error(w, bookingAccessError, http.StatusForbidden)
		return
	}

	// get tickets
	tickets, err := query.GetTicketsByBookingID(r.Context(), h.DB, uint32(bookingID))
	if err != nil {
		http.Error(w, fetchTicketsError, http.StatusInternalServerError)
		return
	}

	// return tickets in json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}
