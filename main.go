package main

import "fmt"

func main() {
	var conference_name string = "Go Conference"
	const conference_tickets uint8 = 50
	var remaining_tickets uint8 = 50

	fmt.Printf("Welcome to %v booking appilication\n", conference_name)
	fmt.Printf("We have total of %d titckets and %d are available!\n", conference_tickets, remaining_tickets)
	fmt.Printf("Get your tickets here to attend\n")

	var user_name string
	var user_tickets uint8

	// User creds
	user_name = "Tom"
	user_tickets = 2
	fmt.Printf("User %s booked %d tickets.\n", user_name, user_tickets)

}
