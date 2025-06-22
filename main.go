package main

import "fmt"

func main() {
	var conference_name = "Go Conference"
	const conference_tickets = 50
	var remaining_tickets = 50

	fmt.Println("Welcome to", conference_name, "booking appilication")
	fmt.Println("We have total of", conference_tickets, "titckets and", remaining_tickets, "are available!")
	fmt.Println("Get your tickets here to attend")
}
