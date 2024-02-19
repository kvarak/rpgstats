package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// --------------------------------------------------------
// HELPER
// --------------------------------------------------------

func processCharacterData(collection CharacterCollection, keySelector, valueSelector func(char Character) string) map[string]map[string]int {
	resultMap := make(map[string]map[string]int)
	allKeys := make(map[string]bool)

	for _, char := range collection.Characters {
		key := keySelector(char)
		value := valueSelector(char)
		if _, ok := resultMap[key]; !ok {
			resultMap[key] = make(map[string]int)
		}
		resultMap[key][value]++
		allKeys[value] = true
	}

	// Ensure all keys exist in each sub-map, set to 0 if missing
	for _, subMap := range resultMap {
		for value := range allKeys {
			if _, ok := subMap[value]; !ok {
				subMap[value] = 0
			}
		}
	}

	return resultMap
}

// --------------------------------------------------------
// TIME
// --------------------------------------------------------

func handleTimeSeries(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	simplifiedData := make([]map[string]string, len(collection.Characters))
	for i, char := range collection.Characters {
		simplifiedData[i] = map[string]string{
			"Irlstart": char.Irlstart,
			// "Died": char.Died, // Uncomment if tracking deaths
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(simplifiedData)
}

// --------------------------------------------------------
// HANDLE
// --------------------------------------------------------

func handleLevelsLived(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	keySelector := func(char Character) string { return char.Category }
	valueSelector := func(char Character) string { return char.LevelsLived }
	categoryCounts := processCharacterData(collection, keySelector, valueSelector)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryCounts)
}

// --------------------------------------------------------

func handleRacesByPlayer(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	keySelector := func(char Character) string { return char.Race }
	valueSelector := func(char Character) string { return char.Category }
	categoryCounts := processCharacterData(collection, keySelector, valueSelector)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryCounts)
}

// --------------------------------------------------------

func handleClassesByPlayer(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	keySelector := func(char Character) string { return char.Category }
	valueSelector := func(char Character) string { return char.Class1 }
	categoryCounts := processCharacterData(collection, keySelector, valueSelector)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryCounts)
}

// --------------------------------------------------------

func handleClassesByPlayerMulti(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	keySelector := func(char Character) string { return char.Category }
	valueSelector1 := func(char Character) string { return char.Class1 }
	valueSelector2 := func(char Character) string { return char.Class2 }

	categoryCounts1 := processCharacterData(collection, keySelector, valueSelector1)
	categoryCounts2 := processCharacterData(collection, keySelector, valueSelector2)

	mergedCounts := categoryCounts1

	for key, valueMap := range mergedCounts {
		for value, _ := range valueMap {
			mergedCounts[key][value] += categoryCounts2[key][value]
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mergedCounts)
}

// --------------------------------------------------------

func handleCharacterCount(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	keySelector := func(char Character) string { return char.Category }
	valueSelector := func(char Character) string { return char.Path }
	categoryPathCounts := processCharacterData(collection, keySelector, valueSelector)

	// reset as processCharacterData counts existance
	for _, catMap := range categoryPathCounts {
		for _, char := range collection.Characters {
			catMap[char.Path] = 0
		}
	}

	// Adjust the logic to include the "Died" and "Extralife" handling
	for category, catMap := range categoryPathCounts {
		for _, char := range collection.Characters {
			if char.Category == category {
				if char.Died == "y" {
					catMap[char.Path]++
				}
				extralifeDeaths, err := strconv.Atoi(char.Extralife)
				if err == nil {
					catMap[char.Path] += extralifeDeaths
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryPathCounts)
}

// --------------------------------------------------------

func handleDeathsByPath(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()
	keySelector := func(char Character) string { return char.Path }
	valueSelector := func(char Character) string { return char.Category }
	deathsByPath := processCharacterData(collection, keySelector, valueSelector)

	// reset as processCharacterData counts existance
	for _, catMap := range deathsByPath {
		for _, char := range collection.Characters {
			catMap[char.Category] = 0
		}
	}

	// Adjust the logic to include the "Died" and "Extralife" handling
	for path, catMap := range deathsByPath {
		for _, char := range collection.Characters {
			if char.Path == path {
				if char.Died == "y" {
					catMap[char.Category]++
				}
				extralifeDeaths, err := strconv.Atoi(char.Extralife)
				if err == nil {
					catMap[char.Category] += extralifeDeaths
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deathsByPath)
}
