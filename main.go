package main

import (
	"fmt"

	"github.com/bootcamp-go/desafio-go-bases/internal/tickets"
)

func main() {
	data, err := tickets.ExtractTicketData("./desafio-go-bases/tickets.csv")
	if err != nil {
		fmt.Println(err)
	}

	result, err := tickets.AverageDestination(data, "China")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
