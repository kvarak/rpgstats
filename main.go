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

	fsCampaigns := http.FileServer(http.Dir("campaigns"))
	http.Handle("/campaigns/", http.StripPrefix("/campaigns/", fsCampaigns))
	fsFonts := http.FileServer(http.Dir("fonts"))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", fsFonts))
	fsScript := http.FileServer(http.Dir("script"))
	http.Handle("/script/", http.StripPrefix("/script/", fsScript))
	fsImages := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fsImages))
	fsCss := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fsCss))

	http.HandleFunc("/data/allThePaths", handleTimeAdventures)
	http.HandleFunc("/data/allTheData", handleAllTheData)
	http.ListenAndServe(":8089", nil)
}
