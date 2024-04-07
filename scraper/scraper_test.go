package scraper

import (
	"fmt"
	"testing"
)

func TestGetF1PlanetLiveRaceResult(t *testing.T) {
	result := GetF1PlanetLiveRaceResult()

	fmt.Println(result)
	// Check if the raceHeader fields are not empty
	if result.RaceHeader.Name == "" || result.RaceHeader.Circuit == "" {
		t.Errorf("raceHeader fields are empty, got: %v", result.RaceHeader)
	}

	if result.RaceStatus != "NOT_STARTED" {
		// Check if the standings slice is not empty
		if len(result.Standings) == 0 {
			t.Errorf("No standings found")
		}
	}

	if result.RaceStatus == "FINISHED" {
		// Check if the fastestLap fields are not empty
		if result.FastestLap.Name == "" || result.FastestLap.Team == "" {
			t.Errorf("fastestLap fields are empty, got: %v", result.FastestLap)
		}
	}
}
