package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fetch data from Google Sheets
		collection := getGoogleSheetData()

		t, err := template.ParseFiles(MainTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html") // Set the content type to "text/html"
		err = t.Execute(w, collection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	// collection := getSheetData()
	// generateStatistics(collection)

	fsFonts := http.FileServer(http.Dir("fonts"))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", fsFonts))
	fsScript := http.FileServer(http.Dir("script"))
	http.Handle("/script/", http.StripPrefix("/script/", fsScript))
	fsImages := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fsImages))
	fsCss := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fsCss))

	http.HandleFunc("/data/levelsLived", handleLevelsLived)
	http.HandleFunc("/data/racesByPlayer", handleRacesByPlayer)
	http.HandleFunc("/data/classesByPlayer", handleClassesByPlayer)
	http.HandleFunc("/data/classesByPlayerMulti", handleClassesByPlayerMulti)
	http.HandleFunc("/data/deathsByPath", handleDeathsByPath)
	http.HandleFunc("/data/timeseries", handleTimeSeries)
	http.HandleFunc("/data/charactercount", handleCharacterCount)
	http.ListenAndServe(":8080", nil)
}
