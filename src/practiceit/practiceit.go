package main

import (
	"log"

	"example.com/backend"
)

func main() {
	backend, err := backend.New("sqlite3", "../../practiceit.db", ":9003")
	if err != nil {
		log.Fatal(err.Error())
	}
	backend.Run()
}
