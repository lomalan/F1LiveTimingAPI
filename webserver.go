package main

import (
	"filivetimingapi/httphandler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/live/data", httphandler.GetLiveDataHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
