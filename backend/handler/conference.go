package handler

import (
	"backend/middleware"
	"backend/models"
	"backend/query"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConferenceHandler struct {
	DB *pgxpool.Pool
}

func NewConferenceHandler(db *pgxpool.Pool) *ConferenceHandler {
	return &ConferenceHandler{DB: db}
}

func (h *ConferenceHandler) RegisterRoutes(r chi.Router) {
	r.Route("/conference", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)

		r.With(middleware.JWTAuthMiddleware).Post("/", h.CreateConference)       // organizer only
		r.Get("/upcoming", h.GetUpcomingConferences)                             // public
		r.Get("/{id}", h.GetConferenceByID)                                      // public
		r.With(middleware.JWTAuthMiddleware).Put("/{id}", h.UpdateConference)    // organizer only
		r.With(middleware.JWTAuthMiddleware).Delete("/{id}", h.DeleteConference) // organizer only
	})
}

// create conference => organizer
func (h *ConferenceHandler) CreateConference(w http.ResponseWriter, r *http.Request) {
	type createConferenceRequest struct {
		Title        string `json:"title"`
		Description  string `json:"description"`
		Location     string `json:"location"`
		EventTime    string `json:"event_time"`
		TotalTickets uint32 `json:"total_tickets"`
	}

	// fetch user id and role from context
	userIDVal := r.Context().Value(middleware.UserIDKey)
	roleVal := r.Context().Value(middleware.RoleKey)

	// validate context presence
	userID, ok1 := userIDVal.(uint32)
	role, ok2 := roleVal.(string)
	if !ok1 || !ok2 || role != "organizer" {
		http.Error(w, notOrganizerError, http.StatusUnauthorized)
		return
	}

	// parse json body
	var req createConferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, requestBodyError, http.StatusBadRequest)
		return
	}

	// parse event time in time.Time format
	eventTime, err := time.Parse(time.RFC3339, req.EventTime)
	if err != nil {
		http.Error(w, eventTimeError, http.StatusBadRequest)
		return
	}

	// creates conference
	conference := models.Conference{
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		EventTime:        eventTime,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		OrganizerID:      userID,
		Status:           "ongoing",
	}

	conferenceID, err := query.CreateConference(r.Context(), h.DB, &conference)
	if err != nil {
		http.Error(w, createConferenceError, http.StatusInternalServerError)
		return
	}

	// respond with conference id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"conference_id": conferenceID,
	})
}

// get all conferences
func (h *ConferenceHandler) GetUpcomingConferences(w http.ResponseWriter, r *http.Request) {
	// parse query param ?days=
	days := 30 // default
	if val := r.URL.Query().Get("days"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			days = parsed
		}
	}

	// fetched upcoming conferences
	confs, err := query.GetUpcomingConferences(r.Context(), h.DB, days)
	if err != nil {
		http.Error(w, conferencesFetchError+err.Error(), http.StatusInternalServerError)
		return
	}

	// return as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(confs)
}

// get conference by id
func (h *ConferenceHandler) GetConferenceByID(w http.ResponseWriter, r *http.Request) {
	// extract id from url
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, conferenceIDError, http.StatusBadRequest)
		return
	}

	// fetch the conference from data base
	conf, err := query.GetConferenceByID(r.Context(), h.DB, uint32(id))
	if err != nil {
		http.Error(w, conferenceNotFoundError, http.StatusNotFound)
		return
	}

	// return as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conf)
}

// update conference => organizer
func (h *ConferenceHandler) UpdateConference(w http.ResponseWriter, r *http.Request) {
	// update struct
	type updateRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		EventTime   string `json:"EventTime"`
		Status      string `json:"status"`
	}

	// get conference id
	idString := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, conferenceIDError, http.StatusBadRequest)
		return
	}

	// extract user id and role
	userID, ok1 := r.Context().Value(middleware.UserIDKey).(uint32)
	role, ok2 := r.Context().Value(middleware.RoleKey).(string)
	if !ok1 || !ok2 || role != "organizer" {
		http.Error(w, notOrganizerError, http.StatusUnauthorized)
		return
	}

	// parse json body
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, requestBodyError, http.StatusBadRequest)
		return
	}

	// parse event time in time.Time format
	eventTime, err := time.Parse(time.RFC3339, req.EventTime)
	if err != nil {
		http.Error(w, eventTimeError, http.StatusBadRequest)
		return
	}

	// update
	err = query.UpdateConference(
		r.Context(),
		h.DB,
		uint32(id),
		userID,
		req.Title,
		req.Description,
		req.Location,
		eventTime,
		req.Status,
	)
	if err != nil {
		http.Error(w, updateConferenceError+err.Error(), http.StatusBadRequest)
		return
	}

	// response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Conference updated successfully"))
}

// delete conference => organizer
func (h *ConferenceHandler) DeleteConference(w http.ResponseWriter, r *http.Request) {
	// get user id and role from context
	userID, ok1 := r.Context().Value("user_id").(uint32)
	role, ok2 := r.Context().Value("role").(string)

	if !ok1 || !ok2 || role != "organizer" {
		http.Error(w, notOrganizerError, http.StatusUnauthorized)
		return
	}

	// get conference id from url
	idString := chi.URLParam(r, "id")
	confID, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		http.Error(w, conferenceIDError, http.StatusBadRequest)
		return
	}

	// perform delete operation
	err = query.DeleteConference(r.Context(), h.DB, uint32(confID), userID)
	if err != nil {
		http.Error(w, deleteConferenceError+err.Error(), http.StatusForbidden)
		return
	}

	// success response
	w.WriteHeader(http.StatusNoContent)
}
