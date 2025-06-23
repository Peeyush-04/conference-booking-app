package main

import (
	"fmt"
	"strings"
)

type User struct {
	FirstName     string
	LastName      string
	Email         string
	BookedTickets uint8
}

func main() {
	conferenceName := "Go Conference"
	const conferenceTickets uint8 = 50
	remainingTickets := uint8(50)

	// Welcome
	fmt.Printf("Welcome to the %v booking application\n", conferenceName)
	fmt.Printf("We have a total of %d tickets and %d are available!\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend!")

	// Input variables
	var firstName, lastName, email string
	var userTickets uint8

	// Booking storage (index → User)
	index := uint8(0)
	bookings := make(map[uint8]User)

	for {
		// Input name
		fmt.Print("\nFirst Name: ")
		fmt.Scan(&firstName)
		fmt.Print("Last Name: ")
		fmt.Scan(&lastName)

		// Valid name check
		for len(firstName) < 2 || len(lastName) < 2 {
			fmt.Printf("\n%s %s is not a valid name.\nFirst and last names should be at least 2 characters long.\n", firstName, lastName)
			fmt.Print("First Name: ")
			fmt.Scan(&firstName)
			fmt.Print("Last Name: ")
			fmt.Scan(&lastName)
		}

		// Input email
		fmt.Print("Email: ")
		fmt.Scan(&email)

		// Valid email check
		for !strings.HasSuffix(email, "@gmail.com") {
			fmt.Printf("\n%s is not a valid email. Please use an @gmail.com email.\n", email)
			fmt.Print("Email: ")
			fmt.Scan(&email)
		}

		// Input booking
		fmt.Print("Booking Tickets: ")
		fmt.Scan(&userTickets)

		// Valid ticket count
		for userTickets == 0 || userTickets > remainingTickets {
			fmt.Printf("\n%d is an invalid number of tickets.\nRemaining: %d\nEnter a value between 1 and %d\n", userTickets, remainingTickets, remainingTickets)
			fmt.Print("Booking Tickets: ")
			fmt.Scan(&userTickets)
		}

		// Store booking
		user := User{
			FirstName:     firstName,
			LastName:      lastName,
			Email:         email,
			BookedTickets: userTickets,
		}
		bookings[index] = user
		index++
		remainingTickets -= userTickets

		// Confirmation
		fmt.Printf("\nThank you %s %s for booking %d ticket(s). A confirmation email has been sent to %s.\n", user.FirstName, user.LastName, user.BookedTickets, user.Email)
		fmt.Printf("Remaining Tickets: %d\n", remainingTickets)

		// Break if sold out
		if remainingTickets == 0 {
			fmt.Println("\nAll tickets are sold out. Booking closed.")
			break
		}
	}

	// Final list
	fmt.Println("\nFinal Booking List:")
	for i := uint8(0); i < index; i++ {
		user := bookings[i]
		fmt.Printf("%d. %s %s (%d tickets) - %s\n", i+1, user.FirstName, user.LastName, user.BookedTickets, user.Email)
	}
}
