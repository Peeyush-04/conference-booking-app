package handlers

import (
	"booking-database/api/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BookTicketHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check both method request are same
		if r.Method != http.MethodPost {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var bookingReq models.BookingRequest

		// decode incoming json body
		err := json.NewDecoder(r.Body).Decode(&bookingReq)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// input validation
		if bookingReq.FirstName == "" ||
			bookingReq.LastName == "" ||
			bookingReq.Email == "" ||
			bookingReq.Tickets == 0 {
			http.Error(w, "All fields must be provided and tickets > 0", http.StatusBadRequest)
			return
		}

		// context timeout for DB interaction
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// inserts data into postgres
		query := `INSERT INTO bookings (first_name, last_name, email, tickets) VALUES ($1, $2, $3, $4)`
		_, err = db.Exec(ctx, query,
			bookingReq.FirstName,
			bookingReq.LastName,
			bookingReq.Email,
			bookingReq.Tickets)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		// success interaction
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Booking successful!",
			"data":    bookingReq,
		})
	}
}
