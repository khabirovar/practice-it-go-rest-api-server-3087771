package library

import (
	"fmt"
	"log"
	"net/http"
)

func Test() {
	fmt.Println("Success call Test function from librarly module")
}

func Response(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Success response from library\n")
}

func ListenAndServe() {
	http.HandleFunc("/", Response)
	fmt.Printf("Server start on :9004\n")
	log.Fatal(http.ListenAndServe(":9004", nil))
}