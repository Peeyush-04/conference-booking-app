package handler

import (
	"backend/auth"
	"backend/models"
	"backend/query"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{DB: db}
}

// routes for register, login
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"` // customer or organizer
	}

	// fetch creds
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, requestBodyError, http.StatusBadRequest)
		return
	}

	// validate role
	if req.Role != "customer" && req.Role != "organizer" {
		http.Error(w, roleError, http.StatusBadRequest)
		return
	}

	// user model
	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
	}

	// create user in DB
	userID, err := query.CreateUser(r.Context(), h.DB, user, req.Password)
	if err != nil {
		log.Println("CreateUser error:", err)
		http.Error(w, createUserError, http.StatusInternalServerError)
		return
	}

	// generate JWT
	token, err := auth.GenerateJWT(userID, req.Role)
	if err != nil {
		http.Error(w, generateTokenError, http.StatusInternalServerError)
		return
	}

	// respond with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// fetch creds
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, requestBodyError, http.StatusBadRequest)
		return
	}

	// fetch user by email
	user, err := query.GetUserByEmail(r.Context(), h.DB, req.Email)
	if err != nil {
		http.Error(w, emailPasswordError, http.StatusUnauthorized)
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, emailPasswordError, http.StatusUnauthorized)
		return
	}

	// generate JWT
	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		http.Error(w, generateTokenError, http.StatusInternalServerError)
		return
	}

	// respond with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// logout logic => depended on seesion or local storage
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": logoutMessage,
	})
}
