package scraper

import (
	"log"

	"github.com/gocolly/colly"
)

type RaceResult struct {
	raceHeader RaceHeader
	standings  []PilotSummary
	fastestLap PilotSummary
}
type RaceHeader struct {
	name    string
	circuit string
}

type PilotSummary struct {
	position string
	name     string
	team     string
	time     string
	stops    string
}

func GetF1PlanetLiveRaceResult() RaceResult {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
	})

	var raceHeader RaceHeader

	c.OnHTML(".pf_race_list_head", func(e *colly.HTMLElement) {
		raceHeader = RaceHeader{
			name:    e.ChildText("h1"),
			circuit: e.ChildText("p"),
		}
	})

	var raceSummary []PilotSummary
	c.OnHTML(".signalr_live_race_html", func(e *colly.HTMLElement) {
		raceSummary = make([]PilotSummary, 0)
		e.ForEach("tr", func(_ int, eh *colly.HTMLElement) {
			raceSummary = append(raceSummary, PilotSummary{
				position: eh.ChildText("td:nth-child(1)"),
				name:     eh.ChildText("td:nth-child(3) h2"),
				team:     eh.ChildText("td:nth-child(3) p"),
				time:     eh.ChildText("td:nth-child(4)"),
				stops:    eh.ChildText("td:nth-child(5)"),
			})
		})

	})
	var fastestLap PilotSummary

	c.OnHTML(".signalr_live_race_fastestLap_html", func(e *colly.HTMLElement) {
		fastestLap = PilotSummary{
			name:  e.ChildText("td:nth-child(3) h2"),
			team:  e.ChildText("td:nth-child(3) p"),
			time:  e.ChildText("td:nth-child(4)"),
			stops: e.ChildText("td:nth-child(5)"),
		}
	})

	var raceResult RaceResult
	c.OnScraped(func(r *colly.Response) {
		raceResult = RaceResult{
			raceHeader: raceHeader,
			standings:  raceSummary,
			fastestLap: fastestLap,
		}
		log.Println("Finished", r.Request.URL)
	})

	c.Visit("https://live.planetf1.com/")

	c.Wait()

	return raceResult
}
