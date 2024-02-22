package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getScenarioName(number string) string {
	switch number {
	case "1":
		return "01: Rise of the Runelords"
	case "2":
		return "02: Curse of the Crimson Throne"
	case "3":
		return "03: Second Darkness"
	case "4":
		return "04: Legacy of Fire"
	case "5":
		return "05: Vanguard of Hope"
	case "6":
		return "06: Into the Dark Continent"
	case "7":
		return "07: Haunted Lands"
	case "8":
		return "08: Return to Sandpoint"
	case "9":
		return "09: Skull & Shackles"
	case "10":
		return "10: Pathfinders"
	default:
		return "99: Unknown"
	}
}

type PageVariables struct {
	Title   string
	Header  string
	Content template.HTML
}

type CharacterCollection struct {
	Characters []Character `json:"characters"`
}

type Character struct {
	Text        string `json:"name"`
	Shortname   string `json:"shortname"`
	Irlstart    string `json:"irlstart"`
	Irlend      string `json:"irlend"`
	Irltime     string `json:"irltime"`
	Igstart     string `json:"igstart"`
	Igend       string `json:"igend"`
	Igtime      string `json:"igtime"`
	Race        string `json:"race"`
	Class1      string `json:"class1"`
	Spec1       string `json:"spec1"`
	Class2      string `json:"class2"`
	Spec2       string `json:"spec2"`
	Totalclass  string `json:"totalclass"`
	Amalgam     string `json:"amalgam"`
	Classtype   string `json:"classtype"`
	Killer_old  string `json:"killer_old"`
	Killercr    string `json:"killercr"`
	Killer      string `json:"killer"`
	Path        string `json:"path"`
	PathNumber  string `json:"pathnumber"`
	Category    string `json:"player"`
	Died        string `json:"died"`
	Extralife   string `json:"extralife"`
	Ressadv     string `json:"ressadv"`
	Resskiller  string `json:"resskiller"`
	Crlvldiff   string `json:"crlvldiff"`
	Startlevel  string `json:"startlevel"`
	Description string `json:"description"`
	Maxlvl      string `json:"maxlvl"`
	Maxlvl2     string `json:"maxlvl2"`
	Event       string `json:"event"`
	LevelsLived string `json:"levelslived"`
	Info        string `json:"info"`
	Deaths      string `json:"deaths"`
}

func makeUpTheStory(cha Character) string {
	character := cha.Description
	death := cha.Event
	character = strings.ReplaceAll(character, "<<", "<img id=\"imageleft\" src=\"")
	character = strings.ReplaceAll(character, ">>", "\" onclick=\"showImage(this.src)\">")
	death = strings.ReplaceAll(death, "<<", "<img id=\"imageright\" src=\"\"")
	death = strings.ReplaceAll(death, ">>", "\" onclick=\"showImage(this.src)\">")
	header := "<h2>" + cha.Text + "</h2>"

	ending := ""
	switch cha.Died {
	case "y":
		ending = "Fate: Dead"
	case "lost":
		ending = "Fate: Unknown"
	case "?":
		ending = "Fate: Retired from adventure"
	default:
		ending = "Fate: Alive"
	}

	return header + "<p>" + character + "</p><h5>" + ending + "</h5><p>" + death + "</p>"
}

// Function that takes the spreadsheet values and creates a list of Characters
func getSheetData() CharacterCollection {
	// Create a new Sheets service
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(CredentialsFile))

	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(SpreadsheetID, ReadRange).Do()
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}

	characters := make([]Character, 0) // Initialize an empty slice for characters
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

		// Safe access and conversion function for integers
		safeGetInt := func(index int) int {
			if index < len(row) {
				if val, err := strconv.Atoi(fmt.Sprint(row[index])); err == nil {
					return val
				}
			}
			return 0 // Return 0 or a sensible default for integers
		}

		totalDeathsInt := safeGetInt(22)
		didDie := safeGetString(21)
		if didDie == "y" {
			totalDeathsInt += 1
		}
		totalDeaths := strconv.Itoa(totalDeathsInt)

		startLevel := safeGetInt(26)
		maxLevel := safeGetInt(28)
		levelsLivedInt := 1 + maxLevel - startLevel
		levelsLived := strconv.Itoa(levelsLivedInt)
		if levelsLivedInt < 10 {
			levelsLived = "0" + levelsLived // Prefix with 0 if less than 10
		}

		// Construct the Character struct safely
		character := Character{
			Text:        safeGetString(0),
			Shortname:   safeGetString(1),
			Igstart:     safeGetString(2),
			Igend:       safeGetString(3),
			Igtime:      safeGetString(4),
			Irlstart:    safeGetString(5),
			Irlend:      safeGetString(6),
			Irltime:     safeGetString(7),
			Race:        safeGetString(8),
			Class1:      safeGetString(9),
			Spec1:       safeGetString(10),
			Class2:      safeGetString(11),
			Spec2:       safeGetString(12),
			Totalclass:  safeGetString(13),
			Amalgam:     safeGetString(14),
			Classtype:   safeGetString(15),
			Killer_old:  safeGetString(16),
			Killercr:    safeGetString(17),
			Killer:      safeGetString(18),
			Path:        getScenarioName(safeGetString(19)),
			PathNumber:  safeGetString(19),
			Category:    safeGetString(20),
			Died:        safeGetString(21),
			Extralife:   safeGetString(22),
			Ressadv:     safeGetString(23),
			Resskiller:  safeGetString(24),
			Crlvldiff:   safeGetString(25),
			Startlevel:  safeGetString(26),
			Description: safeGetString(27),
			Maxlvl:      safeGetString(28),
			Maxlvl2:     safeGetString(29),
			Event:       safeGetString(30),
			LevelsLived: levelsLived,
			Deaths:      totalDeaths,
		}
		character.Info = makeUpTheStory(character)
		characters = append(characters, character)
	}

	return CharacterCollection{
		Characters: characters,
	}
}

func getGoogleSheetData() PageVariables {

	data := getSheetData()

	MyPageVariables := PageVariables{
		Title:   MainTitle,
		Header:  MainHeader,
		Content: "",
	}

	MyPageVariables.Content += "<table>"
	for _, row := range data.Characters {
		MyPageVariables.Content += "<tr>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Path) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Text) + "</td>"
		// MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Shortname) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreakcenter\">" + template.HTML(row.Irlstart) + "<br/>-<br/>"
		MyPageVariables.Content += template.HTML(row.Irlend) + "<br/><br/>"
		MyPageVariables.Content += template.HTML(row.Irltime) + " days IRL</td>"
		MyPageVariables.Content += "<td class=\"td-nobreakcenter\">" + template.HTML(row.Igstart) + "<br/>-<br/>"
		MyPageVariables.Content += template.HTML(row.Igend) + "<br/><br/>"
		MyPageVariables.Content += template.HTML(row.Igtime) + " days in-game</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Race) + "<br/>"
		MyPageVariables.Content += template.HTML(row.Class1) + "(" + template.HTML(row.Spec1) + ")<br/>"
		MyPageVariables.Content += template.HTML(row.Class2) + "(" + template.HTML(row.Spec2) + ")<br/>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">Killer<br/>"
		MyPageVariables.Content += template.HTML(row.Killer_old) + "<br/>"
		MyPageVariables.Content += template.HTML(row.Killercr) + "<br/>"
		MyPageVariables.Content += template.HTML(row.Killer) + "<br/>"
		MyPageVariables.Content += template.HTML(row.Crlvldiff) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Category) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Died) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Extralife) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Ressadv) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Resskiller) + "</td>"
		MyPageVariables.Content += "<td class=\"td-nobreak\">" + template.HTML(row.Startlevel) + "</td>"
		MyPageVariables.Content += "<td>" + template.HTML(row.Description) + "</td>"
		MyPageVariables.Content += "<td>" + template.HTML(row.Maxlvl) + "</td>"
		MyPageVariables.Content += "<td>" + template.HTML(row.Maxlvl2) + "</td>"
		MyPageVariables.Content += "<td>" + template.HTML(row.Event) + "</td>"
		MyPageVariables.Content += "</tr>"
	}
	MyPageVariables.Content += "</table>"
	return MyPageVariables
}

func handleAllTheData(w http.ResponseWriter, r *http.Request) {
	collection := getSheetData()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(collection); err != nil {
		// Handle error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
