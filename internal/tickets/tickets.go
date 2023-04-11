package tickets

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Ticket is a struct that represents a single ticket.
type Ticket struct {
	id            int
	name          string
	email         string
	destination   string
	departureTime time.Time
	ticketPrice   int
}

/*
ExtractTicketData extracts tickets information from a CSV file.
It takes a CSV filename and returns a slice of Ticket structs.

The CSV file must be formatted as follows:
id,name,email,destination,departure_time,ticket_price.
*/
func ExtractTicketData(filename string) ([]Ticket, error) {
	var tickets []Ticket

	// Open the CSV file
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Split the file into lines
	lines := strings.Split(string(file), "\n")

	// Remove the last line from the file (blank line) if exists
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// If the file is empty, return a nil value and an error
	if len(file) == 0 {
		return nil, errors.New("empty CSV file")
	}

	// Loop through each line
	for _, line := range lines {
		// Split the line into fields
		fields := strings.Split(line, ",")

		// Create a new ticket
		ticket := Ticket{}

		// Set the ticket ID
		ticket.id, err = strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		// Set the ticket name
		ticket.name = fields[1]

		// Set the ticket email
		ticket.email = fields[2]

		// Set the ticket destination
		ticket.destination = fields[3]

		// Set the ticket departure time
		ticket.departureTime, err = time.Parse("15:04", fields[4])
		if err != nil {
			return nil, err
		}

		// Set the ticket ticket price
		ticket.ticketPrice, err = strconv.Atoi(fields[5])
		if err != nil {
			return nil, err
		}

		// Add the ticket to the slice
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

/*
GetTotalTicketsByDestination search and count tickets based on the specified destination.
It returns the number of tickets found. If the destination is not found, it returns an error.
*/
func GetTotalTicketsByDestination(data []Ticket, destination string) (int, error) {
	totalTickets := 0
	// If the slice is empty, return an error
	if len(data) == 0 {
		return totalTickets, errors.New("no tickets found")
	}

	// Loop through each ticket
	for _, ticket := range data {
		if ticket.destination == destination {
			totalTickets++
		}
	}

	// Return a error if the destination is not found
	if totalTickets == 0 {
		return 0, errors.New("no tickets found for destination " + destination)
	}

	// Return the total number of tickets found
	return totalTickets, nil
}

/*
CheckTimeBetweenLimits is a utility function that checks if the specified target hour
is between the specified start hour and end hour. It returns true if the target hour
is between the specified start hour and end hour. Otherwise, it returns false.

If the start hour is greater than the end hour, it returns an error.
*/
func checkTimeBetweenLimits(target, start, end time.Time) (bool, error) {
	// If the start time is after the end time, return an error
	if start.After(end) {
		return false, errors.New("start time must be before end time")
	}

	// Check if the target time is between the start and end time
	if target.After(start) && target.Before(end) {
		return true, nil
	}

	// Return false otherwise
	return false, nil
}

/*
GetCountByPeriod receive a slice of Tickets structs and returns a map
containing the total number of tickets for the specified period (morning, afternoon, evening,
early morning).

The time ranges are as follows: morning: between 7:00 and 13:00, afternoon: between 13:00 and
20:00, evening: between 20:00 and 00:00 and early morning: between 00:00 and 7:00.
*/
func GetCountByPeriod(data []Ticket) (map[string]int, error) {
	var countByPeriod = map[string]int{
		"morning":       0,
		"evening":       0,
		"night":         0,
		"early_morning": 0,
	}

	// If the slice is empty, return an error
	if len(data) == 0 {
		return nil, errors.New("no tickets found")
	}

	// Definition of lower and upper limits for each period
	morningLowerLimit, _ := time.Parse("15:04:05", "6:59:59")
	morningUpperLimit, _ := time.Parse("15:04:05", "13:00:00")
	eveningLowerLimit, _ := time.Parse("15:04:05", "12:59:59")
	eveningUpperLimit, _ := time.Parse("15:04:05", "20:00:00")
	nightLowerLimit, _ := time.Parse("15:04:05", "19:59:59")
	nightUpperLimit, _ := time.Parse("15:04:05", "23:59:59")
	earlyMorningLowerLimit, _ := time.Parse("15:04:05", "0:00:00")
	earlyMorningUpperLimit, _ := time.Parse("15:04:05", "7:00:00")

	// Loop through each ticket
	for _, ticket := range data {
		departureTime := ticket.departureTime
		isMorning, _ := checkTimeBetweenLimits(
			departureTime,
			morningLowerLimit,
			morningUpperLimit,
		)
		isEvening, _ := checkTimeBetweenLimits(
			departureTime,
			eveningLowerLimit,
			eveningUpperLimit,
		)
		isNight, _ := checkTimeBetweenLimits(
			departureTime,
			nightLowerLimit,
			nightUpperLimit,
		)
		isEarlyMorning, _ := checkTimeBetweenLimits(
			departureTime,
			earlyMorningLowerLimit,
			earlyMorningUpperLimit,
		)
		if isMorning {
			countByPeriod["morning"]++
		}
		if isEvening {
			countByPeriod["evening"]++
		}
		if isNight {
			countByPeriod["night"]++
		}
		if isEarlyMorning {
			countByPeriod["early_morning"]++
		}
	}
	return countByPeriod, nil
}

/*
AverageDestination calculates the percentage of all emitted tickets that have a certain destination.

It returns the percentage of all emitted tickets that have a given destination. If the destination is not
found or if the data is empty, it returns an error.
*/
func AverageDestination(data []Ticket, destination string) (float64, error) {
	// Obtain the total amount of tickets with the specified destination
	targetTickets, err := GetTotalTicketsByDestination(data, destination)

	// If the destination is not found or the data is empty, return an error
	if err != nil {
		return 0, err
	}

	// Otherwise, calculate the percentage of all emitted tickets with the specified destination
	return float64(targetTickets) / float64(len(data)), nil
}
