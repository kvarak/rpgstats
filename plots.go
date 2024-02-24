package main

import (
	"encoding/json"
	"net/http"
)

// --------------------------------------------------------
// HANDLE CHARACTERS
// --------------------------------------------------------

func handleAllTheData(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(collection); err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
