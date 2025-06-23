package main

import (
	"fmt"
)

func booking_logic(
	bookings *[]string,
	index *uint8,
	first_name *string,
	last_name *string,
	user_tickets *uint8,
	remaining_tickets *uint8,
	email *string,
) bool {
	full_name := *first_name + " " + *last_name
	if *user_tickets > *remaining_tickets {
		fmt.Printf("Booking limit exceeded, only %d are currently avaiable!\n", *remaining_tickets)
		return true
	} else {
		fmt.Printf("%s has booked %d tickets.\n\n", *first_name, *user_tickets)
		*remaining_tickets -= *user_tickets
		*bookings = append(*bookings, full_name)
		fmt.Printf("Remaining tickets: %d\n\n", *remaining_tickets)
		fmt.Printf("Thank you %s for booking %d tickets.\nYou will receive the confirmation email at %s.\n", full_name, *user_tickets, *email)
		if *remaining_tickets == 0 {
			return false
		}
	}
	*index += 1
	return true
}
