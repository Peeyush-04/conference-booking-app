package main

import "fmt"

func main() {
	var conference_name = "Go Conference"
	const conference_tickets = 50
	var remaining_tickets = 50

	fmt.Printf("Welcome to %v booking appilication\n", conference_name)
	fmt.Printf("We have total of %d titckets and %d are available!\n", conference_tickets, remaining_tickets)
	fmt.Printf("Get your tickets here to attend\n")
}
