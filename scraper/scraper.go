package scraper

import (
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Message struct {
	Text string `json:"text"`
	Name string `json:"name"`
}

type RaceResult struct {
	RaceHeader RaceHeader     `json:"header"`
	RaceStatus string         `json:"status"`
	RaceTime   string         `json:"time"`
	Standings  []PilotSummary `json:"standings"`
	FastestLap PilotSummary   `json:"fastestLap"`
}
type RaceHeader struct {
	Name    string `json:"name"`
	Circuit string `json:"circuit"`
}

type PilotSummary struct {
	Position string `json:"position,omitempty"`
	Name     string `json:"name,omitempty"`
	Team     string `json:"team,omitempty"`
	Time     string `json:"gap,omitempty"`
	Stops    string `json:"stops,omitempty"`
}

const layout = "1/2/2006 3:04:05 PM"
const outputLayout = "2006-1-2 3:04:05 PM"

func GetF1PlanetLiveRaceResult() RaceResult {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
	)

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	var raceHeader RaceHeader

	c.OnHTML(".pf_race_list_head", func(e *colly.HTMLElement) {
		raceHeader = RaceHeader{
			Name:    e.ChildText("h1"),
			Circuit: e.ChildText("p"),
		}
	})

	var raceDateTime time.Time
	c.OnHTML(".race_type_Id_selected_class_3182 span.race_type_list", func(e *colly.HTMLElement) {
		raceDateTime = toDateTime(e.Attr("data-sdatetime")).Local()
	})

	var raceSummary []PilotSummary
	c.OnHTML(".signalr_live_race_html", func(e *colly.HTMLElement) {
		raceSummary = make([]PilotSummary, 0)
		e.ForEach("tr", func(_ int, eh *colly.HTMLElement) {
			raceSummary = append(raceSummary, PilotSummary{
				Position: eh.ChildText("td:nth-child(1)"),
				Name:     eh.ChildText("td:nth-child(3) h2"),
				Team:     eh.ChildText("td:nth-child(3) p"),
				Time:     eh.ChildText("td:nth-child(4)"),
				Stops:    eh.ChildText("td:nth-child(5)"),
			})
		})

	})
	var fastestLap PilotSummary

	c.OnHTML(".signalr_live_race_fastestLap_html", func(e *colly.HTMLElement) {
		pilotName := e.ChildText("td:nth-child(3) h2")
		if pilotName != "" {
			fastestLap = PilotSummary{
				Name:  pilotName,
				Team:  e.ChildText("td:nth-child(3) p"),
				Time:  e.ChildText("td:nth-child(4)"),
				Stops: e.ChildText("td:nth-child(5)"),
			}
		}
	})

	var raceStatus string
	var raceResult RaceResult
	c.OnScraped(func(r *colly.Response) {
		raceStatus = getRaceStatus(raceDateTime, fastestLap.Name)
		raceResult = RaceResult{
			RaceHeader: raceHeader,
			RaceTime:   raceDateTime.Format(outputLayout),
			RaceStatus: raceStatus,
			Standings:  raceSummary,
			FastestLap: fastestLap,
		}
		log.Println("Finished", r.Request.URL)
	})

	c.Visit("https://live.planetf1.com/")

	c.Wait()

	return raceResult
}

func toDateTime(dateToParse string) time.Time {
	time, err := time.Parse(layout, dateToParse)
	if err != nil {
		log.Println("Error parsing time:", err)
	}
	return time
}

func getRaceStatus(raceDateTime time.Time, fastestLap string) string {
	now := time.Now()
	if now.Before(raceDateTime) {
		return "NOT_STARTED"
	}
	if now.After(raceDateTime) && fastestLap == "" {
		return "LIVE"
	} else {
		return "FINISHED"
	}
}
