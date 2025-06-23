package main

import (
	"fmt"
	"strings"
)

func main() {
	conference_name := string("Go Conference")
	const conference_tickets uint8 = 50 // conference limit
	remaining_tickets := uint8(50)

	// Welcome
	fmt.Printf("Welcome to %v booking appilication\n", conference_name)
	fmt.Printf("We have total of %d titckets and %d are available!\n", conference_tickets, remaining_tickets)
	fmt.Printf("Get your tickets here to attend\n")

	// input variables =>(values) default
	var first_name string
	var last_name string
	var email string
	var user_tickets uint8

	// data containers
	index := uint8(0)
	bookings := []string{}

	for {
		// input name
		fmt.Printf("\nFirst Name: ")
		fmt.Scan(&first_name)
		fmt.Printf("Last Name: ")
		fmt.Scan(&last_name)

		// valid name check
		valid_name := len(first_name) >= 2 && len(last_name) >= 2
		for {
			if !valid_name {
				fmt.Printf("\n%s %s is not a valid name.\nFirst name and Last name should be greater that equal to 2.\n", first_name, last_name)
				// input valid name
				fmt.Printf("\nFirst Name: ")
				fmt.Scan(&first_name)
				fmt.Printf("Last Name: ")
				fmt.Scan(&last_name)
			} else {
				break // valid name found => exit
			}
		}

		// input email
		fmt.Printf("Email: ")
		fmt.Scan(&email)

		// valid email check
		valid_email := strings.Contains(email, "@gmail.com")
		for {
			if !valid_email {
				fmt.Printf("\n%s is not a valid email.\nUse @gmail.com at last.", email)
				// input valid email
				fmt.Printf("Email: ")
				fmt.Scan(&email)
			} else {
				break // valid email found => exit
			}
		}

		// input user booked tickets
		fmt.Printf("Booking Tickets: ")
		fmt.Scan(&user_tickets)

		// valid user tickets check
		valid_user_tickets := user_tickets > 0 && user_tickets <= remaining_tickets
		for {
			if !valid_user_tickets {
				fmt.Printf("\n%d are invalid number of tickets.\nRemaining tickets: %d\nValue should range from 1 to %d", user_tickets, remaining_tickets, remaining_tickets)
				// input valid user tickets
				fmt.Printf("\nBooking Tickets: ")
				fmt.Scan(&user_tickets)
			} else {
				break
			}
		}
		fmt.Printf("\n")

		// Booking Logic
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
		// End booking logic
	}

	// Printing Final Booking List
	fmt.Printf("\nFinal Booking list:\n")
	for i := uint8(0); i <= index; i++ {
		fmt.Printf("%d. %s\n", i+1, bookings[i])
	}
	fmt.Println("")
}
