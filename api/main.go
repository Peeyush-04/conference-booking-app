package main

import (
	"booking-database/api/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// connects to postgresql
	dbPool := db.ConnectDB()
	defer dbPool.Close()

	// registering routes
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is Alive!!")
	})
	// http.HandleFunc("/book", handlers.BookTicketHandler(dbPool))

	fmt.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server Failed.", err)
	}
}
