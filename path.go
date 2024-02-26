package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type AdventureCollection struct {
	Adventures []Adventure `json:"adventures"`
}

type Adventure struct {
	Pathnr              string `json:"pathnr"`
	Path                string `json:"path"`
	Advnr               string `json:"advnr"`
	Adventure           string `json:"adventure"`
	Igstart             string `json:"igstart"`
	Igend               string `json:"igend"`
	Irlstart            string `json:"irlstart"`
	Irlend              string `json:"irlend"`
	StartLevel          string `json:"startLevel"`
	EndLevel            string `json:"endLevel"`
	FinalBoss           string `json:"finalBoss"`
	PercentCompleted    string `json:"percentCompl"`
	ShortIntro          string `json:"shortIntro"`
	AdventureBackground string `json:"advBkgrnd"`
	OtherBackground     string `json:"otherBckgrnd"`
	Igtime              string `json:"igtime"`
	Irltime             string `json:"irltime"`
	Info                string `json:"info"`
}

//${adventure.shortIntro}<hr>${adventure.adventureBackground}<hr>${adventure.otherBackground}`

func pathInfo(adv Adventure) string {
	intro := ""
	if adv.ShortIntro != "" {
		intro = adv.ShortIntro
		intro = strings.ReplaceAll(intro, "<<", "<img class=\"imageleft\" src=\"")
		intro = strings.ReplaceAll(intro, ">>", "\" onclick=\"showImage(this.src)\">")
		intro = "<a name=\"intro\"><h4>Short Introduction</h4></a>" + intro
		intro += "<hr>"
	}
	background := ""
	if adv.AdventureBackground != "" {
		background = adv.AdventureBackground
		background = strings.ReplaceAll(background, "<<", "<img class=\"imageright\" src=\"")
		background = strings.ReplaceAll(background, ">>", "\" onclick=\"showImage(this.src)\">")
		background = "<a name=\"background\"><h4>Adventure Background</h4></a>" + background
		background += "<hr>"
	}
	other := ""
	if adv.OtherBackground != "" {
		other = adv.OtherBackground
		other = strings.ReplaceAll(other, "<<", "<img class=\"imageleft\" src=\"")
		other = strings.ReplaceAll(other, ">>", "\" onclick=\"showImage(this.src)\">")
		other = "<a name=\"other\"><h4>Other Background</h4></a>" + other
		other += "<hr>"
	}
	boss := ""
	if adv.FinalBoss != "" {
		boss = adv.FinalBoss
		boss = strings.ReplaceAll(boss, "<<", "<img class=\"imageright\" src=\"")
		boss = strings.ReplaceAll(boss, ">>", "\" onclick=\"showImage(this.src)\">")
		boss = "<a name=\"boss\"><h4>Final Boss</h4></a>" + boss
	}

	header := "<h2>" + adv.Adventure + "</h2>"
	header += "<p><b><i>Path " + adv.Path + "</i></b></p>"

	pathbox := ""
	pathbox += "<div class=\"classbox frame\">"
	pathbox += "<table class=\"classbox\">"
	pathbox += "<thead><tr><th>Path " + adv.Path + "</th></tr></thead>"
	pathbox += "<tbody>"
	pathbox += "<tr><td><i>Adventure " + adv.Advnr + ": " + adv.Adventure + "</i></td></tr>"
	pathbox += "<tr><td><hr></td></tr>"

	// Add Dates
	pathbox += "<tr><td markdown=\"1\" class=\"classbox\">"
	pathbox += adv.Igstart + " - " + adv.Igend + " (" + adv.Igtime + " days)<br/>"
	pathbox += adv.Irlstart + " - " + adv.Irlend + " (" + adv.Irltime + " days)<br/>"
	pathbox += "</td></tr>"
	pathbox += "<tr><td><hr></td></tr>"

	// Add links to the sections
	if intro+background+other+boss != "" {
		pathbox += "<tr><td markdown=\"1\" class=\"classbox\">"
		if intro != "" {
			pathbox += "<a href=\"#intro\">Introduction</a><br/>"
		}
		if background != "" {
			pathbox += "<a href=\"#background\">Background</a><br/>"
		}
		if other != "" {
			pathbox += "<a href=\"#other\">Other Background</a><br/>"
		}
		if boss != "" {
			pathbox += "<a href=\"#boss\">Final Boss</a><br/>"
		}
		pathbox += "</td></tr>"
	}

	pathbox += "</tbody></table></div>"
	return pathbox + header + intro + background + other + boss
}

// Function that returns the days between two dates in "YYYY-MM-DD" format
func daysBetweenDates(startDate, endDate string) int {
	const dateFormat = "2006-01-02" // Go's reference date format
	start, err := time.Parse(dateFormat, startDate)
	if err != nil {
		return 0
	}
	end, err := time.Parse(dateFormat, endDate)
	if err != nil {
		return 0
	}

	// Calculate difference and convert to whole days
	days := end.Sub(start).Hours() / 24
	return int(days)
}

// Function that takes the spreadsheet values and creates a list of Adventures
func getAdventureSheetData() AdventureCollection {
	// Create a new Sheets service
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(CredentialsFile))

	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(SpreadsheetID, ReadRangePath).Do()
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}

	adventures := make([]Adventure, 0) // Initialize an empty slice for adventures
	for i, row := range resp.Values {
		if i == 0 {
			continue // Skip the first line if it's a header
		}

		// Safe access function
		safeGetString := func(index int) string {
			if index < len(row) {
				return fmt.Sprint(row[index])
			}
			return ""
		}

		// Construct the Adventure struct safely
		adventure := Adventure{
			Pathnr:              safeGetString(0),
			Path:                getScenarioName(safeGetString(0)),
			Advnr:               safeGetString(2),
			Adventure:           safeGetString(3),
			Igstart:             safeGetString(4),
			Igend:               safeGetString(5),
			Irlstart:            safeGetString(6),
			Irlend:              safeGetString(7),
			StartLevel:          safeGetString(8),
			EndLevel:            safeGetString(9),
			FinalBoss:           safeGetString(10),
			PercentCompleted:    safeGetString(11),
			ShortIntro:          safeGetString(12),
			AdventureBackground: safeGetString(13),
			OtherBackground:     safeGetString(14),
			Igtime:              fmt.Sprint(daysBetweenDates(safeGetString(4), safeGetString(5))),
			Irltime:             fmt.Sprint(daysBetweenDates(safeGetString(6), safeGetString(7))),
		}
		adventures = append(adventures, adventure)
	}

	// For each adventure, make up the story
	for i, adv := range adventures {
		adventures[i].Info = pathInfo(adv)
	}

	return AdventureCollection{
		Adventures: adventures,
	}
}

// --------------------------------------------------------
// HANDLE ADVENTURES
// --------------------------------------------------------

func handleTimeAdventures(w http.ResponseWriter, r *http.Request) {
	adventures := getAdventureSheetData()
	// keySelector := func(char Adventure) string { return char.Path }
	// valueSelector := func(char Adventure) string { return char.LevelsLived }
	// categoryCounts := processAdventureData(adventures, keySelector, valueSelector)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adventures)
}
