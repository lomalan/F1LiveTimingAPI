package httphandler

import (
	"encoding/json"
	"filivetimingapi/scraper"
	"log"
	"net/http"
)

func GetLiveDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := scraper.GetF1PlanetLiveRaceResult()
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
