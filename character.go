package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"html/template"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type PageVariables struct {
	Title   string
	Header  string
	Content template.HTML
}

type CharacterCollection struct {
	Characters []Character `json:"characters"`
}

type Character struct {
	Text            string `json:"name"`
	Shortname       string `json:"shortname"`
	Irlstart        string `json:"irlstart"`
	Irlend          string `json:"irlend"`
	Irltime         string `json:"irltime"`
	Igstart         string `json:"igstart"`
	Igend           string `json:"igend"`
	Igtime          string `json:"igtime"`
	Race            string `json:"race"`
	Class1          string `json:"class1"`
	Spec1           string `json:"spec1"`
	Class2          string `json:"class2"`
	Spec2           string `json:"spec2"`
	Totalclass      string `json:"totalclass"`
	Amalgam         string `json:"amalgam"`
	Classtype       string `json:"classtype"`
	Killer_old      string `json:"killer_old"`
	Killercr        string `json:"killercr"`
	Killer          string `json:"killer"`
	Path            string `json:"path"`
	PathNumber      string `json:"pathnumber"`
	Category        string `json:"player"`
	Died            string `json:"died"`
	Extralife       string `json:"extralife"`
	Ressadv         string `json:"ressadv"`
	Resskiller      string `json:"resskiller"`
	Crlvldiff       string `json:"crlvldiff"`
	Startlevel      string `json:"startlevel"`
	Description     string `json:"description"`
	Maxlvl          string `json:"maxlvl"`
	Maxlvl2         string `json:"maxlvl2"`
	Event           string `json:"event"`
	LevelsLived     string `json:"levelslived"`
	Info            string `json:"info"`
	Deaths          string `json:"deaths"`
	Lifescore       string `json:"lifescore"`
	Classscore      string `json:"classscore"`
	Classaverage    string `json:"classaverage"`
	Pathscore       string `json:"pathscore"`
	Pathaverage     string `json:"pathaverage"`
	Playerscore     string `json:"playerscore"`
	Playeraverage   string `json:"playeraverage"`
	Characterscore  string `json:"characterscore"`
	Comboavglvl     string `json:"comboavglvl"`
	Comboavgdays    string `json:"comboavgdays"`
	Comboavgirldays string `json:"comboavgirldays"`
	Comboavgkill    string `json:"comboavgkill"`
	Totalscore      string `json:"totalscore"`
}

func calculateScore(character Character) string {
	// Define weights for attributes
	const survivedWeight float64 = 2
	const extralifeWeight float64 = 1
	const crlvldiffWeight float64 = 0.5
	const levelsLivedWeight float64 = 0.25

	// Start with a base score
	var score float64 = 0

	// Add score if the character didn't sie at the end
	if strings.ToLower(character.Died) == "n" {
		score += survivedWeight * 2
	} else if strings.ToLower(character.Died) != "y" {
		score += survivedWeight
	}

	// Subtract score based on extralife
	if extralife, err := strconv.Atoi(character.Extralife); err == nil {
		score -= extralifeWeight * float64(extralife)
	}

	// Subtract score based on crlvldiff
	// If killed by a creature of a higher level, Crlvldiff is negative
	if crlvldiff, err := strconv.Atoi(character.Crlvldiff); err == nil {
		score -= crlvldiffWeight * float64(crlvldiff)
	}

	// Add score based on levels lived
	if levelsLived, err := strconv.Atoi(character.LevelsLived); err == nil {
		score += levelsLivedWeight * float64(levelsLived)
	}

	return strconv.Itoa(int(score))
}

func createClassbox(cha Character) string {
	classbox := ""
	classbox += "<div class=\"classbox frame\">"
	classbox += "<table class=\"classbox\">"
	classbox += "<thead><tr><th>" + cha.Text + "</th></tr></thead>"
	classbox += "<tbody>"
	classbox += "<tr><td><i>" + cha.Race + "</i></td></tr>"
	classbox += "<tr><td><hr></td></tr>"

	// Add Path if not empty
	if cha.Path != "" {
		classbox += "<tr><td markdown=\"1\" class=\"classbox\">"
		classbox += "<strong>Path</strong> " + cha.Path + ""
		classbox += "</td></tr>"
		classbox += "<tr><td><hr></td></tr>"
	}

	classbox += "<tr><td><i>" + cha.Irlstart + " - " + cha.Irlend + " (" + cha.Irltime + " days)</i></td></tr>"
	if cha.Igend != "" {
		classbox += "<tr><td><i>" + cha.Igstart + " - " + cha.Igend + " (" + cha.Igtime + " days)</i></td></tr>"
	}
	classbox += "<tr><td><hr></td></tr>"
	classbox += "<tr><td><i>Start level: " + cha.Startlevel + "</i></td></tr>"
	classbox += "<tr><td markdown=\"1\" class=\"classbox\">"
	classbox += "<strong>" + cha.Class1
	if cha.Spec1 != "" && cha.Spec1 != "<low lvl>" {
		classbox += " (" + cha.Spec1 + ")"
	}
	if cha.Maxlvl2 != "" {
		maxlvl, _ := strconv.Atoi(cha.Maxlvl)
		maxlvl2, _ := strconv.Atoi(cha.Maxlvl2)
		newmaxlvl := strconv.Itoa(maxlvl - maxlvl2)
		classbox += "</strong> " + newmaxlvl
	} else {
		classbox += "</strong> " + cha.Maxlvl
	}
	// Add class2 if not empty
	if cha.Class2 != "" {
		if cha.Spec2 != "" && cha.Spec2 != "<low lvl>" {
			classbox += "<br/><strong>" + cha.Class2 + " (" + cha.Spec2 + ")</strong> " + cha.Maxlvl2
		} else {
			classbox += "<br/><strong>" + cha.Class2 + "</strong> " + cha.Maxlvl2
		}
	}
	classbox += "</td></tr>"
	classbox += "<tr><td><hr></td></tr>"

	// Add Resurrection if not empty
	if cha.Extralife != "" {
		classbox += "<tr><td markdown=\"1\" class=\"classbox\">"
		classbox += "<strong>Resurrections:</strong> " + cha.Extralife + "<br/>"
		classbox += "Killed by: " + cha.Resskiller + "<br/>"
		classbox += "In adventure(s): " + cha.Ressadv
		classbox += "</td></tr>"
		classbox += "<tr><td><hr></td></tr>"
	}

	// Add Killer if not empty
	if cha.Died == "y" {
		classbox += "<tr><td markdown=\"1\" class=\"classbox\">"
		if cha.Killer_old != "" {
			classbox += "<strong>Final death by:</strong><br/>" + cha.Killer + "/" + cha.Killer_old + " (CR: " + cha.Killercr + ")<br/>"
		} else {
			classbox += "<strong>Final death by:</strong><br/>" + cha.Killer + " (CR: " + cha.Killercr + ")<br/>"
		}
		classbox += "</td></tr>"
		classbox += "<tr><td><hr></td></tr>"
	}

	// Add Scores
	classbox += "<tr><td markdown=\"1\" class=\"classbox\">"
	classbox += "<strong>Scores:</strong><br/>"
	classbox += "<strong>Character total score:</strong> " + cha.Characterscore + "<br/>"
	classbox += "<i>Life: " + cha.Lifescore + " / Class: " + cha.Classscore + " / Path: " + cha.Pathscore + "</i><br/>"
	classbox += "</td></tr>"

	classbox += "</tbody></table></div>"

	return classbox
}

func makeUpTheStory(cha Character) string {
	character := cha.Description
	death := cha.Event
	character = strings.ReplaceAll(character, "<<", "<img class=\"imageleft\" src=\"")
	character = strings.ReplaceAll(character, ">>", "\" onclick=\"showImage(this.src)\">")
	death = strings.ReplaceAll(death, "<<", "<img class=\"imageright\" src=\"\"")
	death = strings.ReplaceAll(death, ">>", "\" onclick=\"showImage(this.src)\">")
	header := "<h2>" + cha.Shortname + "</h2>"

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

	classbox := createClassbox(cha)

	return classbox + " " + header + "<p>" + character + "</p><h5>" + ending + "</h5><p>" + death + "</p>"
}

func normalizeLifescores(characters []Character) ([]Character, error) {
	var min, max float64
	first := true

	// Find min and max Lifescore values
	for _, c := range characters {
		score, err := strconv.ParseFloat(c.Lifescore, 64)
		if err != nil {
			return nil, err // Handle the error appropriately
		}

		if first {
			min, max = score, score
			first = false
		} else {
			if score < min {
				min = score
			}
			if score > max {
				max = score
			}
		}
	}

	// Normalize Lifescores to a 1-100 scale
	for i, c := range characters {
		score, err := strconv.ParseFloat(c.Lifescore, 64)
		if err != nil {
			return nil, err // Handle the error appropriately
		}

		normalized := int(1 + ((score - min) * 99 / (max - min)))
		characters[i].Lifescore = strconv.Itoa(normalized)
	}

	return characters, nil
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
		character.Lifescore = calculateScore(character)
		characters = append(characters, character)
	}

	characters, err = normalizeLifescores(characters)
	if err != nil {
		fmt.Println("Error normalizing Lifescores:", err)
	}

	// Calculate the average character.Lifescore for each character.Path
	// And save it into character.Pathaverage
	scores := make(map[string]float64)
	counts := make(map[string]int)
	for _, character := range characters {
		lifescore, _ := strconv.ParseFloat(character.Lifescore, 64)
		scores[character.Path] += lifescore
		counts[character.Path]++
	}
	for i := range characters {
		lookup := characters[i].Path
		average := int(scores[lookup] / float64(counts[lookup]))
		characters[i].Pathaverage = strconv.Itoa(average)
		lifescore, _ := strconv.Atoi(characters[i].Lifescore)
		characters[i].Pathscore = strconv.Itoa((lifescore - average + 100) / 2)
	}

	// Calculate the average character.Lifescore for each character.Totalclass
	// And save it into character.Classaverage
	scores = make(map[string]float64)
	counts = make(map[string]int)
	for _, character := range characters {
		lifescore, _ := strconv.ParseFloat(character.Lifescore, 64)
		scores[character.Totalclass] += lifescore
		counts[character.Totalclass]++
	}
	for i := range characters {
		lookup := characters[i].Totalclass
		average := int(scores[lookup] / float64(counts[lookup]))
		characters[i].Classaverage = strconv.Itoa(average)
		lifescore, _ := strconv.Atoi(characters[i].Lifescore)
		characters[i].Classscore = strconv.Itoa((lifescore - average + 100) / 2)
	}

	// Calculate the average character.Lifescore for each character.Category
	// And save it into character.Playeraverage
	scores = make(map[string]float64)
	counts = make(map[string]int)
	for _, character := range characters {
		lifescore, _ := strconv.ParseFloat(character.Lifescore, 64)
		scores[character.Category] += lifescore
		counts[character.Category]++
	}
	for i := range characters {
		lookup := characters[i].Category
		average := int(scores[lookup] / float64(counts[lookup]))
		characters[i].Playeraverage = strconv.Itoa(average)
		lifescore, _ := strconv.Atoi(characters[i].Lifescore)
		characters[i].Playerscore = strconv.Itoa((lifescore - average + 100) / 2)
	}

	for i := range characters {
		class, _ := strconv.Atoi(characters[i].Classscore)
		// player, _ := strconv.Atoi(characters[i].Playerscore)
		path, _ := strconv.Atoi(characters[i].Pathscore)
		life, _ := strconv.Atoi(characters[i].Lifescore)
		characters[i].Characterscore = strconv.Itoa(class + path + life)
	}

	//	Comboavglvl    string `json:"comboavglvl"`
	//	Comboavgdays   string `json:"comboavgdays"`
	//	Comboavgkill   string `json:"comboavgkill"`

	today := time.Now().Format("2006-01-02")
	comboTotalLevels := 0.0
	comboTotalDays := 0.0
	comboTotalIrlDays := 0.0
	comboTotalKill := ""
	comboCounts := 0.0

	for i, current := range characters {
		comboTotalLevels = 0.0
		comboTotalDays = 0.0
		comboTotalIrlDays = 0.0
		comboTotalKill = ""
		comboCounts = 0.0
		if current.Irlend == today {
			for _, character := range characters {
				if character.Totalclass == current.Totalclass || character.Category == current.Category {
					levelsLived, _ := strconv.Atoi(character.LevelsLived)
					igtime, _ := strconv.Atoi(character.Igtime)
					irltime, _ := strconv.Atoi(character.Irltime)
					killer := character.Killer
					killer_old := character.Killer_old
					resskiller := character.Resskiller
					comboTotalLevels += float64(levelsLived)
					comboTotalDays += float64(igtime)
					comboTotalIrlDays += float64(irltime)
					comboCounts++
					if character.Died == "y" {
						// only add , if there is already a killer
						if comboTotalKill == "" {
							comboTotalKill += killer
						} else {
							comboTotalKill += ", " + killer
						}
						if comboTotalKill == "" {
							comboTotalKill += killer_old
						} else {
							comboTotalKill += ", " + killer_old
						}
					}
					if character.Extralife != "" {
						if comboTotalKill == "" {
							comboTotalKill += resskiller
						} else {
							comboTotalKill += ", " + resskiller
						}
					}
				}
			}
			characters[i].Comboavglvl = strconv.FormatFloat(comboTotalLevels/comboCounts, 'f', 2, 64)
			characters[i].Comboavgdays = strconv.FormatFloat(comboTotalDays/comboCounts, 'f', 2, 64)
			characters[i].Comboavgirldays = strconv.FormatFloat(comboTotalIrlDays/comboCounts, 'f', 2, 64)
			characters[i].Comboavgkill = comboTotalKill
		}

		// For each character, make up the story
		characters[i].Info = makeUpTheStory(characters[i])

	}

	return CharacterCollection{
		Characters: characters,
	}
}

func generateStory(character Character) string {
	templates := storyStart

	// Generate a hash from the character's name to select a template
	h := fnv.New32a()
	h.Write([]byte(character.Text + character.Maxlvl))
	hash := h.Sum32()
	templateIndex := int(hash) % len(templates)

	// Replace placeholders in the selected template
	story := templates[templateIndex]
	story = strings.Replace(story, "{name}", "<b>"+character.Shortname+"</b>", -1)
	story = strings.Replace(story, "{class}", character.Totalclass, -1)
	race := character.Race
	if strings.Contains(race, "(") && strings.Contains(race, ")") {
		startIndex := strings.Index(race, "(")
		endIndex := strings.Index(race, ")")
		race = race[startIndex+1 : endIndex]
	}
	story = strings.Replace(story, "{race}", race, -1)
	story = strings.Replace(story, "{date}", character.Igstart, -1)

	return story
}

func addDaysToDate(dateStr, daysStr string) string {
	date, _ := time.Parse("2006-01-02", dateStr)
	daysFloat, _ := strconv.ParseFloat(daysStr, 64)
	days := int(daysFloat)
	newDate := date.AddDate(0, 0, days)
	return newDate.Format("2006-01-02")
}

func GenerateLevelOutcome(currentLvlStr, willDieStr, uniqueId string) string {
	currentLvl, err1 := strconv.Atoi(currentLvlStr)
	willDie, err2 := strconv.ParseFloat(willDieStr, 64)
	if err1 != nil || err2 != nil {
		return "Invalid input."
	}

	diff := math.Abs(float64(currentLvl) - willDie)

	// Use a hash of the uniqueId to seed the random number generator
	hash := fnv.New32a()
	_, err := hash.Write([]byte(uniqueId))
	if err != nil {
		return "Error generating hash."
	}
	seed := int64(hash.Sum32())
	rand.Seed(seed)

	if float64(currentLvl) < willDie {
		messagesLiveLonger := []string{
			fmt.Sprintf("seems destined for great things, with at least %d more levels to conquer.", int(math.Ceil(diff))),
			fmt.Sprintf("has a bright future ahead, promising %d more levels of adventure.", int(math.Ceil(diff))),
			fmt.Sprintf("is on a path of glory, with %d levels yet to be explored.", int(math.Ceil(diff))),
			fmt.Sprintf("holds the promise of enduring %d more levels, according to the stars.", int(math.Ceil(diff))),
			fmt.Sprintf("is blessed by fate, with a journey extending %d levels further.", int(math.Ceil(diff))),
			fmt.Sprintf("carries an aura of resilience, likely surviving %d additional levels.", int(math.Ceil(diff))),
			fmt.Sprintf("is shielded by unseen forces for another %d levels at least.", int(math.Ceil(diff))),
			fmt.Sprintf("has tales unwritten for %d more levels, as foretold by the ancients.", int(math.Ceil(diff))),
			fmt.Sprintf("is whispered in prophecies to stride across %d levels more.", int(math.Ceil(diff))),
			fmt.Sprintf("has a destiny that spans %d additional levels, so say the oracles.", int(math.Ceil(diff))),
			fmt.Sprintf("is marked by the gods to journey through %d more levels.", int(math.Ceil(diff))),
			fmt.Sprintf("will tread %d more levels, untouched by shadow.", int(math.Ceil(diff))),
			fmt.Sprintf("is fated to rise above challenges for at least %d more levels.", int(math.Ceil(diff))),
			fmt.Sprintf("will outlast %d more levels, as the winds of fortune suggest.", int(math.Ceil(diff))),
			fmt.Sprintf("has %d levels of untold stories waiting to unfold.", int(math.Ceil(diff))),
			fmt.Sprintf("is surrounded by a light that promises %d more levels of life.", int(math.Ceil(diff))),
			fmt.Sprintf("will dance through %d more levels, guided by luck.", int(math.Ceil(diff))),
			fmt.Sprintf("is foreseen to have %d more levels of laughter and joy.", int(math.Ceil(diff))),
			fmt.Sprintf("will continue to weave tales of heroism for %d more levels.", int(math.Ceil(diff))),
			fmt.Sprintf("is destined to leave footprints across %d additional levels.", int(math.Ceil(diff))),
		}
		return messagesLiveLonger[rand.Intn(len(messagesLiveLonger))]
	} else {
		messagesLivedLonger := []string{
			fmt.Sprintf("has cheated fate, living %d levels beyond the foretold end.", int(math.Floor(diff))),
			fmt.Sprintf("carries the whispers of survival, outlasting %d levels against all odds.", int(math.Floor(diff))),
			fmt.Sprintf("has defied the grim predictions by %d levels, and still stands strong.", int(math.Floor(diff))),
			fmt.Sprintf("is a living legend, having surpassed expected demise by %d levels.", int(math.Floor(diff))),
			fmt.Sprintf("basks in borrowed time, %d levels past the shadow of death.", int(math.Floor(diff))),
			fmt.Sprintf("walks with the ghosts of %d levels past, yet breathes life.", int(math.Floor(diff))),
			fmt.Sprintf("has seen %d levels more than destiny intended.", int(math.Floor(diff))),
			fmt.Sprintf("rides the winds of chance, %d levels past the destined fall.", int(math.Floor(diff))),
			fmt.Sprintf("wears a cloak of defiance, surviving %d levels past doom.", int(math.Floor(diff))),
			fmt.Sprintf("is a beacon of hope, shining %d levels past darkness.", int(math.Floor(diff))),
			fmt.Sprintf("is the master of their fate, living %d levels beyond the written end.", int(math.Floor(diff))),
			fmt.Sprintf("stands tall, %d levels beyond the whisper of death.", int(math.Floor(diff))),
			fmt.Sprintf("holds stories of %d levels that were never meant to be.", int(math.Floor(diff))),
			fmt.Sprintf("is the anomaly in the tapestry of fate, %d levels beyond their time.", int(math.Floor(diff))),
			fmt.Sprintf("is the unexpected chapter, continuing %d levels past the finale.", int(math.Floor(diff))),
			fmt.Sprintf("has outpaced destiny by %d levels, forging a new path.", int(math.Floor(diff))),
			fmt.Sprintf("is a tale of resilience, %d levels beyond the end.", int(math.Floor(diff))),
			fmt.Sprintf("lives as if borrowed time is infinite, %d levels beyond reckoning.", int(math.Floor(diff))),
			fmt.Sprintf("has unwritten %d levels of history, beyond what was foreseen.", int(math.Floor(diff))),
			fmt.Sprintf("is the unexpected hero, surviving %d levels past prophecy.", int(math.Floor(diff))),
		}
		return messagesLivedLonger[rand.Intn(len(messagesLivedLonger))]
	}
}

func GenerateMessageBasedOnDates(currentDate, dieDate, hash string) string {
	// Parse the input dates
	layout := "2006-01-02"
	current, err := time.Parse(layout, currentDate)
	if err != nil {
		fmt.Println("Error parsing currentDate:", err)
		return ""
	}
	die, err := time.Parse(layout, dieDate)
	if err != nil {
		fmt.Println("Error parsing dieDate:", err)
		return ""
	}

	// Calculate the difference in days
	diff := int(die.Sub(current).Hours() / 24)

	// Seed the random generator with the hash for deterministic output
	seed := int64(0)
	for _, c := range hash {
		seed = (seed*31 + int64(c)) % 1e9
	}
	rand.Seed(seed)

	// Message templates
	if diff > 0 {
		templates := []string{
			fmt.Sprintf("Legend has it that the final day approaches, marked on the calendar as %s.", dieDate),
			fmt.Sprintf("Whispers in the wind speak of a fateful day drawing near: %s.", dieDate),
			fmt.Sprintf("Legends whisper of %s, a day shrouded in destiny, awaiting to unfold.", dieDate),
			fmt.Sprintf("Fate's loom weaves a tale of %s, where destiny's crossroads shall meet.", dieDate),
			fmt.Sprintf("A prophecy carved in the stars speaks of %s, marking a pivotal turn in the saga.", dieDate),
			fmt.Sprintf("Ancient scrolls foretell %s as a day of reckoning, where fate shall be sealed.", dieDate),
			fmt.Sprintf("Seers have long envisioned %s, where destiny's threads will intertwine.", dieDate),
			fmt.Sprintf("A harbinger of destiny, %s looms over the horizon, shrouded in mystery.", dieDate),
			fmt.Sprintf("The oracle's vision reveals %s as a crucial juncture in the fabric of time.", dieDate),
			fmt.Sprintf("Tales of yore speak of %s, a day foreseen to change the course of history.", dieDate),
			fmt.Sprintf("A celestial alignment on %s portends a momentous event, written in the cosmos.", dieDate),
			fmt.Sprintf("The echo of fate resounds towards %s, a day that will decide all.", dieDate),
			fmt.Sprintf("On %s, the scales of destiny will tip, marking a new chapter in the annals of time.", dieDate),
			fmt.Sprintf("The shadow of %s casts forth, promising a turning point of cosmic significance.", dieDate),
			fmt.Sprintf("Eldritch powers have long marked %s as a day of unparalleled importance.", dieDate),
			fmt.Sprintf("The chorus of the ancients crescendos towards %s, a day of fated encounters.", dieDate),
			fmt.Sprintf("A cosmic confluence on %s signals a convergence of fateful paths.", dieDate),
			fmt.Sprintf("The veil of time thins as %s approaches, a day foretold to bear witness to destiny.", dieDate),
			fmt.Sprintf("A tempest of fate brews towards %s, promising upheaval and rebirth.", dieDate),
			fmt.Sprintf("The stars align in anticipation of %s, a day that will illuminate the fates.", dieDate),
			fmt.Sprintf("As the dawn of %s nears, destiny's hand readies to pen its indelible mark.", dieDate),
		}
		return templates[rand.Intn(len(templates))]
	} else {
		diff = -diff // Make positive since dieDate is in the past
		templates := []string{
			fmt.Sprintf("Against all odds, the hero has lived %d days beyond the foretold end date.", diff),
			fmt.Sprintf("Fate's design was defied, with %d extra days already lived beyond %s.", diff, dieDate),
			fmt.Sprintf("Beyond %s, they tread in borrowed time, each day a gift from the fates.", dieDate),
			fmt.Sprintf("Having eclipsed %s, their story weaves on, defying the edicts of destiny.", dieDate),
			fmt.Sprintf("The shadow of %s lies behind, as they step into uncharted epochs, masters of their fate.", dieDate),
			fmt.Sprintf("Surpassing %s, they now walk a path not charted by stars but by will.", dieDate),
			fmt.Sprintf("The date %s now a tale of the past, as they carve a future untold by oracles.", dieDate),
			fmt.Sprintf("Each day post %s is a testament to their defiance against the preordained.", dieDate),
			fmt.Sprintf("With %s in the annals of history, their journey continues, unwritten and undefined.", dieDate),
			fmt.Sprintf("Beyond the bounds of %s, they forge a legacy beyond the confines of fate.", dieDate),
			fmt.Sprintf("Past %s, they stand, a beacon of resilience in the face of destiny's decree.", dieDate),
			fmt.Sprintf("Having outlasted %s, their saga unfolds in realms unbound by prophecy.", dieDate),
			fmt.Sprintf("The echoes of %s fade into legend, as their steps forge new myths.", dieDate),
			fmt.Sprintf("Beyond the foretold end of %s, they emerge, undiminished and unconquered.", dieDate),
			fmt.Sprintf("With %s receding into the tapestry of time, their story enters unscripted territories.", dieDate),
			fmt.Sprintf("Outliving %s, they redefine the contours of fate, etching a path of their own design.", dieDate),
			fmt.Sprintf("The milestone of %s now surpassed, they journey beyond the realm of prophecy.", dieDate),
			fmt.Sprintf("Having transcended %s, their tale evolves, unbounded by the chains of destiny.", dieDate),
			fmt.Sprintf("Past the threshold of %s, they venture into a future unfettered by foresight.", dieDate),
			fmt.Sprintf("The legacy of %s now a relic, as they pioneer into epochs anew.", dieDate),
			fmt.Sprintf("With %s but a memory, they stride into an era not foreseen, shaping destiny anew.", dieDate),
		}
		return templates[rand.Intn(len(templates))]
	}
}

func GenerateDeathProphecy(cause, hash string) string {
	seed := int64(0)
	for _, c := range hash {
		seed = (seed*31 + int64(c)) % 1e9
	}
	rand.Seed(seed)

	// Get a random prophecy from the selected cause if it exists
	if prophecies, ok := deathProphecies[cause]; ok {
		return prophecies[rand.Intn(len(prophecies))]
	}

	// If the cause is not found, return a random message from defaultProphecies
	return defaultProphecies[rand.Intn(len(defaultProphecies))]
}

func mostCommonItemAll(input string) string {
	// Check if the input is empty or contains only "End boss" after trimming spaces
	if input = strings.TrimSpace(input); input == "" || strings.ReplaceAll(input, "End boss,", "") == "" {
		return "" // Return empty if input is empty or only contains "End boss"
	}

	// Remove "End boss"
	input = strings.Replace(input, "End boss", "", -1)
	input = strings.Replace(input, ", ,", ", ", -1)

	// Split the cleaned input string into items
	items := strings.Split(input, ", ")
	itemCount := make(map[string]int)

	// Count occurrences of each non-empty item
	for _, item := range items {
		item = strings.TrimSpace(item) // Trim spaces to avoid counting empty items
		if item != "" {
			itemCount[item]++
		}
	}

	// Find the max count
	var maxCount int
	for _, count := range itemCount {
		if count > maxCount {
			maxCount = count
		}
	}

	// Collect items that have the max count
	var candidates []string
	for item, count := range itemCount {
		if count == maxCount {
			candidates = append(candidates, item)
		}
	}

	// Returns all candidates in a string with the format "item1, item2 or item3"
	if len(candidates) == 1 {
		return candidates[0]
	}
	if len(candidates) == 2 {
		return candidates[0] + " or " + candidates[1]
	}
	if len(candidates) > 2 {
		return strings.Join(candidates[:len(candidates)-1], ", ") + " or " + candidates[len(candidates)-1]
	}
	return "<unknown>"
}

func mostCommonItem(input, hash string) string {
	// Check if the input is empty or contains only "End boss" after trimming spaces
	if input = strings.TrimSpace(input); input == "" || strings.ReplaceAll(input, "End boss,", "") == "" {
		return "" // Return empty if input is empty or only contains "End boss"
	}

	// Replace "End boss" with an empty string and trim spaces around commas
	input = strings.Replace(input, "End boss", "", -1)
	input = strings.Replace(input, ", ,", ", ", -1)

	// Split the cleaned input string into items
	items := strings.Split(input, ", ")
	itemCount := make(map[string]int)

	// Count occurrences of each non-empty item
	for _, item := range items {
		item = strings.TrimSpace(item) // Trim spaces to avoid counting empty items
		if item != "" {
			itemCount[item]++
		}
	}

	// Find the max count
	var maxCount int
	for _, count := range itemCount {
		if count > maxCount {
			maxCount = count
		}
	}

	// Collect items that have the max count
	var candidates []string
	for item, count := range itemCount {
		if count == maxCount {
			candidates = append(candidates, item)
		}
	}

	if len(candidates) == 0 {
		return "<unknown>"
	}

	// Use hash string to deterministically select a candidate
	hashValue := hashStringToInt(hash)
	selectedIndex := hashValue % len(candidates)

	return candidates[selectedIndex]
}

// hashStringToInt converts a string to an int using FNV hashing.
func hashStringToInt(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func getGoogleSheetData() PageVariables {

	data := getSheetData()

	MyPageVariables := PageVariables{
		Title:   MainTitle,
		Header:  MainHeader,
		Content: "",
	}

	todaysDate := time.Now().Format("2006-01-02")

	MyPageVariables.Content += "<h1>Our current adventurers</h1>"
	deathDate := ""
	for _, row := range data.Characters {

		if row.Irlend == todaysDate {
			MyPageVariables.Content += "<h3>" + template.HTML(row.Text) + "</h3>"
			MyPageVariables.Content += "<p class=\"story\">"
			MyPageVariables.Content += template.HTML(generateStory(row)) + " "
			deathDate = addDaysToDate(row.Igstart, row.Comboavgdays)
			MyPageVariables.Content += template.HTML(GenerateMessageBasedOnDates(row.Igstart, deathDate, row.Text+row.Maxlvl)) + " "
			MyPageVariables.Content += template.HTML(row.Text) + " " + template.HTML(GenerateLevelOutcome(row.LevelsLived, row.Comboavglvl, row.Shortname+row.Maxlvl)) + " "
			MyPageVariables.Content += template.HTML(GenerateDeathProphecy(mostCommonItem(row.Comboavgkill, row.Text+row.Maxlvl), row.Text+row.Maxlvl)) + " "
			MyPageVariables.Content += "</p>"

			extralife, _ := strconv.Atoi(row.Extralife)
			if extralife > 0 {
				MyPageVariables.Content += "<p class=\"story\">"
				MyPageVariables.Content += "As an adventurer, " + template.HTML(row.Shortname) + " has been granted " + template.HTML(row.Extralife) + " extra lives. "
				MyPageVariables.Content += "</p>"
			}

			MyPageVariables.Content += "<p class=\"story\"><i>IRL: Started at "
			MyPageVariables.Content += template.HTML(row.Irlstart) + "; "
			MyPageVariables.Content += template.HTML(row.Irltime) + " days ago. "
			MyPageVariables.Content += template.HTML(addDaysToDate(row.Irlstart, row.Comboavgirldays)) + " is the statistical last day. "
			MyPageVariables.Content += " Killed by " + template.HTML(mostCommonItemAll(row.Comboavgkill)) + ", "
			MyPageVariables.Content += " at level " + template.HTML(row.Comboavglvl) + ". "
			MyPageVariables.Content += " This is based on the average of all adventurers with the same class or player. "
			MyPageVariables.Content += "</i></p>"

			story := row.Description
			if story != "" {
				MyPageVariables.Content += "<h5>Background</h5>"
				story = strings.ReplaceAll(story, "<<", "<img class=\"imageleft\" src=\"")
				story = strings.ReplaceAll(story, ">>", "\" onclick=\"showImage(this.src)\">")
				MyPageVariables.Content += "<p class=\"story\">" + template.HTML(story) + "</p>"
			}

		}
	}
	return MyPageVariables
}
