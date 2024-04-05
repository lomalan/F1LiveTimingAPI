package scraper

import (
	"testing"
)

func TestGetF1PlanetLiveRaceResult(t *testing.T) {
	result := GetF1PlanetLiveRaceResult()

	// Check if the raceHeader fields are not empty
	if result.raceHeader.name == "" || result.raceHeader.circuit == "" {
		t.Errorf("raceHeader fields are empty, got: %v", result.raceHeader)
	}

	// Check if the standings slice is not empty
	if len(result.standings) == 0 {
		t.Errorf("No standings found")
	}
}
