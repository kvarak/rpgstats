package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"html/template"
	"log"
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
	Text           string `json:"name"`
	Shortname      string `json:"shortname"`
	Irlstart       string `json:"irlstart"`
	Irlend         string `json:"irlend"`
	Irltime        string `json:"irltime"`
	Igstart        string `json:"igstart"`
	Igend          string `json:"igend"`
	Igtime         string `json:"igtime"`
	Race           string `json:"race"`
	Class1         string `json:"class1"`
	Spec1          string `json:"spec1"`
	Class2         string `json:"class2"`
	Spec2          string `json:"spec2"`
	Totalclass     string `json:"totalclass"`
	Amalgam        string `json:"amalgam"`
	Classtype      string `json:"classtype"`
	Killer_old     string `json:"killer_old"`
	Killercr       string `json:"killercr"`
	Killer         string `json:"killer"`
	Path           string `json:"path"`
	PathNumber     string `json:"pathnumber"`
	Category       string `json:"player"`
	Died           string `json:"died"`
	Extralife      string `json:"extralife"`
	Ressadv        string `json:"ressadv"`
	Resskiller     string `json:"resskiller"`
	Crlvldiff      string `json:"crlvldiff"`
	Startlevel     string `json:"startlevel"`
	Description    string `json:"description"`
	Maxlvl         string `json:"maxlvl"`
	Maxlvl2        string `json:"maxlvl2"`
	Event          string `json:"event"`
	LevelsLived    string `json:"levelslived"`
	Info           string `json:"info"`
	Deaths         string `json:"deaths"`
	Lifescore      string `json:"lifescore"`
	Classscore     string `json:"classscore"`
	Classaverage   string `json:"classaverage"`
	Pathscore      string `json:"pathscore"`
	Pathaverage    string `json:"pathaverage"`
	Playerscore    string `json:"playerscore"`
	Playeraverage  string `json:"playeraverage"`
	Characterscore string `json:"characterscore"`
	Totalscore     string `json:"totalscore"`
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
		character.Info = makeUpTheStory(character)
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

	return CharacterCollection{
		Characters: characters,
	}
}

func generateStory(character Character) string {
	templates := []string{
		"From {date}, tales of {name}, a {race} {class}, spread like wildfire, a testament to their indomitable spirit and unyielding resolve.",
		"On {date}, {name} the {race} {class} set sail towards destiny, their eyes set on horizons that many deemed unreachable.",
		"Starting on {date}, {name}'s name became synonymous with courage as the {race} {class} ventured into the unknown, a beacon of hope in dark times.",
		"Armed with ancient knowledge and the heart of a warrior, {name}, the {race} {class}, embarked on their fabled journey on {date}.",
		"As the stars aligned on {date}, so did the fate of {name}, the {race} {class}, whose name would become etched in the eternal chronicles of time.",
		"The saga of {name} began on {date}, a {race} {class} whose bravery transcended the tales of old, venturing into realms where angels dare not tread.",
		"On {date}, {name} the {race} {class}, took an oath beneath the ancient oaks, a vow that would steer them through battles unseen and foes unknown.",
		"A whisper on the wind on {date} spoke of {name}, a {race} {class}, chosen by destiny to walk a path lined with shadow and light.",
		"None could have predicted on {date}, how {name}, a simple {race} {class}, would shake the very foundations of the world with their deeds.",
		"It was on {date} that {name}, the {race} {class}, emerged from the mists of legend, a force to be reckoned with, in pursuit of a destiny foretold.",
		"With a heart as fierce as dragons and a will unbreakable, {name}, the {race} {class}, embarked on their legendary quest on {date}, a journey that would echo through ages.",
		"On {date}, {name} the {race} {class}, cast aside all doubt, embracing a destiny that would intertwine their name with the essence of adventure itself.",
		"The legacy of {name}, a {race} {class}, was forged on {date}, amidst the flames of destiny, to become a beacon for those who dare to dream.",
		"From the tranquil shores of {date}, {name}, a {race} {class}, set sail towards the tempest of the unknown, their story a canvas for the ages.",
		"Every tavern's tale and every child's bedtime story held the feats of {name}, the {race} {class}, in awe. It's said that fate's wheel began to turn on {date}.",
		"In the heart of the dense forests, whispers tell of {name}, the {race} {class}. Their journey is said to have begun on {date}, seeking what lies beyond the known paths.",
		"Legends echo through the halls of time about {name}, a {class} of {race} origin, whose adventures embarked from the ancient ruins on {date}.",
		"Fate's hand guided {name}, the {race} {class}, from a humble beginning on {date} to the cusp of reshaping history itself.",
		"Under the silver moon on {date}, {name} the {race} {class} pledged to unravel the mysteries that bind the realms.",
		"It is said that on {date}, the winds of fortune whispered to {name}, a {class} of {race}, propelling them on a quest that would enter the annals of legend.",
		"Beneath the gaze of the crimson moon on {date}, {name} the {race} {class}, embarked on a journey that would etch their name across the heavens and into the annals of history.",
		"On {date}, in the shadow of ancient ruins, {name}, a {race} {class}, uncovered a truth so powerful it threatened to unravel the very fabric of existence.",
		"{name}, the {race} {class}, found their destiny intertwined with the fate of the world on {date}, when a prophecy whispered from the lips of the dying sun came to pass.",
		"The echoes of {name}'s valorous deeds, a {race} {class}, begun on {date}, still resound through the hallowed halls where heroes are remembered.",
		"On {date}, {name}, a {race} {class}, stood at the crossroads of fate, their decision a beacon that would guide the lost through the darkness.",
		"It is said that on {date}, {name} the {race} {class}, danced with the stars, their steps a melody that brought balance to the chaos of the cosmos.",
		"The legend of {name}, a {race} {class}, was born on {date}, from the ashes of a world consumed by darkness, a light fierce enough to challenge the night.",
		"On {date}, whispers of {name}'s arrival, the {race} {class}, stirred the ancient guardians from their slumber, heralding the dawn of a new era.",
		"{name}, the {race} {class}, chose on {date} to walk the path less traveled, their journey a tapestry woven from the threads of countless destinies.",
		"As the first light of {date} pierced the veil of night, {name}, a {race} {class}, set forth to claim their place among the constellations as a paragon of virtue and valor.",
	}

	// Generate a hash from the character's name to select a template
	h := fnv.New32a()
	h.Write([]byte(character.Text + character.Maxlvl))
	hash := h.Sum32()
	templateIndex := int(hash) % len(templates)

	// Replace placeholders in the selected template
	story := templates[templateIndex]
	story = strings.Replace(story, "{name}", "<b>"+character.Shortname+"</b>", -1)
	story = strings.Replace(story, "{class}", character.Totalclass, -1)
	story = strings.Replace(story, "{race}", character.Race, -1)
	story = strings.Replace(story, "{date}", character.Igstart, -1)

	return story
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

	for _, row := range data.Characters {

		if row.Irlend == todaysDate {
			MyPageVariables.Content += "<h4>" + template.HTML(row.Text) + "</h4>"
			MyPageVariables.Content += "<p class=\"story\">"
			MyPageVariables.Content += template.HTML(generateStory(row))
			MyPageVariables.Content += "</p>"

			extralife, _ := strconv.Atoi(row.Extralife)
			if extralife > 0 {
				MyPageVariables.Content += "<p class=\"story\">"
				MyPageVariables.Content += "As an adventurer, " + template.HTML(row.Shortname) + " has been granted " + template.HTML(row.Extralife) + " extra lives. "
				MyPageVariables.Content += "</p>"
			}

			story := row.Description
			story = strings.ReplaceAll(story, "<<", "<img id=\"imageleft\" src=\"")
			story = strings.ReplaceAll(story, ">>", "\" onclick=\"showImage(this.src)\">")
			MyPageVariables.Content += "<p class=\"story\">" + template.HTML(story) + "</p>"

			MyPageVariables.Content += "<p class=\"story\"><i>IRL: Started at "
			MyPageVariables.Content += template.HTML(row.Irlstart) + "; "
			MyPageVariables.Content += template.HTML(row.Irltime) + " days ago."
			MyPageVariables.Content += "</i></p>"
		}
	}
	MyPageVariables.Content += "</table>"
	return MyPageVariables
}
