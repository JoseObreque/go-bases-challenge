package tickets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractTickedData(t *testing.T) {
	t.Run("Open inexistent tickets file", func(t *testing.T) {
		filename := "./inexistent_file.csv"

		tickets, err := ExtractTicketData(filename)

		assert.Nil(t, tickets)
		assert.Error(t, err)
	})

	t.Run("Open a empty tickets file", func(t *testing.T) {
		filename := "./empty_ticket_test.csv"

		tickets, err := ExtractTicketData(filename)

		assert.Nil(t, tickets)
		assert.Error(t, err)
	})

	t.Run("Open a valid tickets file", func(t *testing.T) {
		filename := "./ticket_test.csv"

		expectedTicketTime, _ := time.Parse("15:04", "17:11")
		expectedData := []Ticket{
			{
				1,
				"Tait Mc Caughan",
				"tmc0@scribd.com",
				"Finland",
				expectedTicketTime,
				785,
			},
		}

		data, err := ExtractTicketData(filename)

		assert.Equal(t, expectedData, data)
		assert.NoError(t, err)
	})
}

func TestGetTotalTicketsByDestination(t *testing.T) {
	t.Run("Search in empty ticket slice", func(t *testing.T) {
		destination := "China"
		var ticketSlice []Ticket
		expectedTotalTickets := 0

		totalTickets, err := GetTotalTicketsByDestination(ticketSlice, destination)

		assert.Equal(t, expectedTotalTickets, totalTickets)
		assert.Error(t, err)
	})

	t.Run("Search in valid ticket slice", func(t *testing.T) {
		filename := "./ticket_test.csv"
		destination := "Finland"
		ticketSlice, _ := ExtractTicketData(filename)
		expectedTotalTickets := 1

		totalTickets, err := GetTotalTicketsByDestination(ticketSlice, destination)
		assert.Equal(t, expectedTotalTickets, totalTickets)
		assert.NoError(t, err)
	})

	t.Run("Search in valid ticket slice (0 results)", func(t *testing.T) {
		filename := "./ticket_test.csv"
		destination := "The Moon"
		ticketSlice, _ := ExtractTicketData(filename)
		expectedTotalTickets := 0

		totalTickets, err := GetTotalTicketsByDestination(ticketSlice, destination)
		assert.Equal(t, expectedTotalTickets, totalTickets)
		assert.Error(t, err)
	})
}

func TestCheckTimeBetweenLimits(t *testing.T) {
	t.Run("Lower limit greater than upper limit", func(t *testing.T) {
		targetTime, _ := time.Parse("15:04", "17:11")
		lowerLimit, _ := time.Parse("15:04", "18:11")
		upperLimit, _ := time.Parse("15:04", "16:11")

		result, err := checkTimeBetweenLimits(targetTime, lowerLimit, upperLimit)
		assert.False(t, result)
		assert.Error(t, err)
	})

	t.Run("Target is between lower and upper limit", func(t *testing.T) {
		targetTime, _ := time.Parse("15:04", "17:11")
		lowerLimit, _ := time.Parse("15:04", "16:11")
		upperLimit, _ := time.Parse("15:04", "18:11")

		result, err := checkTimeBetweenLimits(targetTime, lowerLimit, upperLimit)

		assert.True(t, result)
		assert.NoError(t, err)
	})

	t.Run("Target is before lower limit", func(t *testing.T) {
		targetTime, _ := time.Parse("15:04", "15:11")
		lowerLimit, _ := time.Parse("15:04", "16:11")
		upperLimit, _ := time.Parse("15:04", "18:11")

		result, err := checkTimeBetweenLimits(targetTime, lowerLimit, upperLimit)

		assert.False(t, result)
		assert.NoError(t, err)
	})

	t.Run("Target is after upper limit", func(t *testing.T) {
		targetTime, _ := time.Parse("15:04", "19:11")
		lowerLimit, _ := time.Parse("15:04", "16:11")
		upperLimit, _ := time.Parse("15:04", "18:11")

		result, err := checkTimeBetweenLimits(targetTime, lowerLimit, upperLimit)

		assert.False(t, result)
		assert.NoError(t, err)
	})
}

func TestGetCountByPeriod(t *testing.T) {
	t.Run("Search in empty ticket slice", func(t *testing.T) {
		var ticketSlice []Ticket

		count, err := GetCountByPeriod(ticketSlice)

		assert.Nil(t, count)
		assert.Error(t, err)
	})

	t.Run("Search in valid ticket slice", func(t *testing.T) {
		filename := "./ticket_test_2.csv"
		ticketSlice, _ := ExtractTicketData(filename)

		expectedCount := map[string]int{
			"morning":       1,
			"evening":       1,
			"night":         1,
			"early_morning": 1,
		}

		count, err := GetCountByPeriod(ticketSlice)

		assert.Equal(t, expectedCount, count)
		assert.NoError(t, err)
	})
}

func TestAverageDestination(t *testing.T) {
	t.Run("Search in empty ticket slice", func(t *testing.T) {
		var ticketSlice []Ticket
		expectedAvg := float64(0)

		avg, err := AverageDestination(ticketSlice, "China")

		assert.Equal(t, expectedAvg, avg)
		assert.Error(t, err)
	})

	t.Run("Search in valid ticket slice", func(t *testing.T) {
		filename := "./ticket_test_2.csv"
		ticketSlice, _ := ExtractTicketData(filename)

		// The test file contains 2 of 4 registered tickets with destination "China".
		expectedAvg := 0.50

		avg, err := AverageDestination(ticketSlice, "China")

		assert.Equal(t, expectedAvg, avg)
		assert.NoError(t, err)
	})
}
