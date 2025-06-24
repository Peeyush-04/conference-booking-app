package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BookTicketHandler(db *pgxpool.Pool) http.HandlerFunc
