package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"strconv"

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
	Text        string `json:"text"`
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
	Category    string `json:"category"`
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

	characters := make([]Character, len(resp.Values)-1)
	for i, row := range resp.Values {
		if i == 0 {
			continue // Skip the first line
		}
		startLevel, _ := strconv.Atoi(fmt.Sprint(row[26]))
		maxLevel, _ := strconv.Atoi(fmt.Sprint(row[28]))
		levelsLivedInt := maxLevel - startLevel
		var levelsLived string
		if levelsLivedInt < 10 {
			levelsLived = "0" + strconv.Itoa(levelsLivedInt)
		} else {
			levelsLived = strconv.Itoa(levelsLivedInt)
		}
		character := Character{
			Text:        fmt.Sprint(row[0]),
			Shortname:   fmt.Sprint(row[1]),
			Igstart:     fmt.Sprint(row[2]),
			Igend:       fmt.Sprint(row[3]),
			Igtime:      fmt.Sprint(row[4]),
			Irlstart:    fmt.Sprint(row[5]),
			Irlend:      fmt.Sprint(row[6]),
			Irltime:     fmt.Sprint(row[7]),
			Race:        fmt.Sprint(row[8]),
			Class1:      fmt.Sprint(row[9]),
			Spec1:       fmt.Sprint(row[10]),
			Class2:      fmt.Sprint(row[11]),
			Spec2:       fmt.Sprint(row[12]),
			Totalclass:  fmt.Sprint(row[13]),
			Amalgam:     fmt.Sprint(row[14]),
			Classtype:   fmt.Sprint(row[15]),
			Killer_old:  fmt.Sprint(row[16]),
			Killercr:    fmt.Sprint(row[17]),
			Killer:      fmt.Sprint(row[18]),
			Path:        getScenarioName(fmt.Sprint(row[19])),
			Category:    fmt.Sprint(row[20]),
			Died:        fmt.Sprint(row[21]),
			Extralife:   fmt.Sprint(row[22]),
			Ressadv:     fmt.Sprint(row[23]),
			Resskiller:  fmt.Sprint(row[24]),
			Crlvldiff:   fmt.Sprint(row[25]),
			Startlevel:  fmt.Sprint(row[26]),
			Description: fmt.Sprint(row[27]),
			Maxlvl:      fmt.Sprint(row[28]),
			LevelsLived: fmt.Sprint(levelsLived),
		}
		characters[i-1] = character
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
		MyPageVariables.Content += template.HTML(row.Classtype) + "</td>"
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
