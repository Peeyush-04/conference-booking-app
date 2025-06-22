package main

import "fmt"

func main() {
	var conference_name string = "Go Conference"
	const conference_tickets uint8 = 50
	var remaining_tickets uint8 = 50

	fmt.Printf("Welcome to %v booking appilication\n", conference_name)
	fmt.Printf("We have total of %d titckets and %d are available!\n", conference_tickets, remaining_tickets)
	fmt.Printf("Get your tickets here to attend\n")

	var bookings []string
	var first_name string
	var last_name string
	var email string
	var user_tickets uint8
	index := uint8(0)
	for {
		// User creds
		fmt.Printf("\nFirst Name: ")
		fmt.Scan(&first_name)
		fmt.Printf("Last Name: ")
		fmt.Scan(&last_name)
		fmt.Printf("Email: ")
		fmt.Scan(&email)
		fmt.Printf("Booking Tickets: ")
		fmt.Scan(&user_tickets)
		fmt.Printf("\n")

		if user_tickets > remaining_tickets {
			fmt.Printf("Booking limit exceeded, only %d are currently available!\n", remaining_tickets)
			break
		} else {
			fmt.Printf("User %s booked %d tickets.\n\n", first_name, user_tickets)
			remaining_tickets -= user_tickets
			bookings = append(bookings, first_name+" "+last_name)
			fmt.Printf("Remaing Tickets: %d\n\n", remaining_tickets)
			fmt.Printf("Thank you %s for booking %d tickets.\nYou will recieve the confirmantion email at %s\n", bookings[index], user_tickets, email)
			if remaining_tickets == 0 {
				break
			}
		}
		index++
	}

	// Printing Final Booking List
	fmt.Println("Final Booking list: ", bookings)
}
