package main

import (
	"os"

	"github.com/sakjur/rails.tid.app/trafikverket"
)

func main() {
	client, err := trafikverket.NewClient()
	if err != nil {
		panic(err)
	}

	trains, err := client.Trains()
	if err != nil {
		panic(err)
	}

	trafikverket.SVG(os.Stdout, trains)
}
