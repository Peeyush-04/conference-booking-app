package handler

import (
	"backend/query"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	invalidUserError    string = "Invalid user ID"
	userError           string = "User not found"
	internalServerError string = "Internal server error"
	invalidJSONRequest  string = "Invalid JSON request"
)

type UserHandler struct {
	DB *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{DB: db}
}

// user request structure
type updateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// includes register routes in UserHandler function
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)
		r.Delete("/{id}", h.DeleteUser)
	})
}

// includes get user handler in userhandler
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// get user id from url
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, invalidUserError, http.StatusBadRequest)
		return
	}

	// fetch user data from db by user id
	user, err := query.GetUserByID(r.Context(), h.DB, uint32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, userError, http.StatusNotFound)
		} else {
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}

	// return user as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// update user handler
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// fetch id from url
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, invalidUserError, http.StatusBadRequest)
		return
	}

	// decode json
	var request updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, invalidJSONRequest, http.StatusBadRequest)
		return
	}

	// update user info in DB
	err = query.UpdateUserInfo(r.Context(), h.DB, uint32(id), request.FirstName, request.LastName, request.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// delete user handler
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// fetch user ID from url
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, invalidUserError, http.StatusBadRequest)
		return
	}

	// delete user data
	err = query.DeleteUser(r.Context(), h.DB, uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// response with status
	w.WriteHeader(http.StatusNoContent)
}
