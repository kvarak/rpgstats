package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
