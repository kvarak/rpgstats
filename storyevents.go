package main

var defaultProphecies = []string{
	"Fate whispers of a shadow that follows unseen, a silent harbinger of a destiny yet to unfold, leading to an untimely demise.",
	"In the realm of whispers and shadows, an unnamed terror seeks the unwary, its presence a cold touch on the brink of oblivion.",
	"Beware the unseen foe, for it dances on the edge of perception, weaving a tapestry of demise with threads unseen.",
	"A chill wind carries omens of a lurking doom, a specter without name, its intentions as cryptic as the path it guards.",
	"Amongst the echoes of time, a prophecy untold speaks of a demise by forces unknown, a puzzle locked within the fabric of destiny.",
	"The stars hint at a silent threat, a void that consumes light and hope, leaving behind a trail of desolation and a tale of caution.",
}

var deathProphecies = map[string][]string{
	"Naga": {
		"Should you wander carelessly into the ancient ruins, the coils of a Naga might become your final embrace.",
		"A careless step in the marshlands could lead to a venomous end, courtesy of a lurking Naga.",
		"The whisper of a Naga in the night foretells a binding fate for those not wary of their charms.",
		"Beware the serpentine allure of a Naga, for its constricting grasp foresees a breathless doom.",
	},
	"Demon": {
		"If you dabble too freely in the dark arts, a Demon might claim your soul as its prize.",
		"The echo of demonic laughter in the darkness is a harbinger of fiery oblivion for the unwary.",
		"A pact with a Demon, forged in haste, can lead to an eternity of torment and despair.",
		"Heed the infernal temptations not, lest a Demon becomes the architect of your ruin.",
	},
	"Tiefling": {
		"In the shadow of betrayal, a Tiefling's curse could spell your doom unless you guard your back well.",
		"Trust misplaced in a Tiefling companion might lead you down a path of sorrow and darkness.",
		"A Tiefling's revenge, slow and calculated, awaits those who wrong them, a demise ensnared in deceit.",
		"Beware the heart of a Tiefling scorned, for from such wounds, fatal plots are born.",
	},
	"Undead": {
		"The night's silence hides the approach of the Undead, craving the warmth of the living, a fate you might share if not vigilant.",
		"A crypt's false treasure lures the greedy, only to find their company among the restless Undead.",
		"In the fog of the graveyard, the Undead whisper your name, predicting a cold embrace unless you heed their call.",
		"An unprepared visit to haunted grounds could see you join the ranks of the Undead, forever bound to the earth.",
	},
	"Guild of Thieves": {
		"A slip of attention could see your fortunes stolen by the Guild of Thieves, or worse, your life.",
		"The Guild of Thieves watches from the shadows, ready to cut short the lives of those who flash their wealth carelessly.",
		"Entanglement in the Guild of Thieves' schemes promises a blade in the dark for those not cautious.",
		"A misplaced trust in the city's underbelly might find you at the wrong end of a Thieves' Guild blade.",
	},
	"Barbarian": {
		"The thunderous charge of a Barbarian horde on the horizon spells doom for those not fortified against their rage.",
		"In the wilderness, a lone Barbarian's cry might be the last thing you hear if you tread without caution.",
		"The clash of steel heralds the arrival of a Barbarian's wrath, a storm you may not weather.",
		"Beware the lands ruled by Barbarian clans, for their justice is swift and brutal to outsiders.",
	},
	"Monstrosity": {
		"The shadow of a Monstrosity looms over those who explore the unknown without heed.",
		"In the depths of the forest, a Monstrosity's roar foretells a gruesome end for the unwary traveler.",
		"Tales of Monstrosities are not just stories to scare children; ignore them, and you may become a cautionary tale.",
		"The mark of a true adventurer is caution, for a Monstrosity awaits those who boast without substance.",
	},
	"Whispering Ways": {
		"The Whispering Ways seduce with forbidden knowledge, leading down a path of ruin for the incautious.",
		"Beware the allure of the Whispering Ways; their secrets ensnare souls in a web from which there is no escape.",
		"A mind open to the Whispering Ways without guard can find itself shattered by truths not meant for mortals.",
		"The call of the Whispering Ways is a siren's song, leading only to madness and despair for those who listen.",
	},
	"Aberration": {
		"The gaze of an Aberration promises a reality undone, a future untold for those caught in its sight.",
		"An Aberration's whisper across the void seeks to unravel the sanity of those who dare to listen.",
		"To face an Aberration without caution is to invite a fate of being torn asunder by the fabric of reality itself.",
		"The mere presence of an Aberration distorts destiny, a twisted path awaiting those who approach unwarily.",
	},
	"PK (vildmagi)": {
		"A chance encounter with wild magic gone awry might end your journey in an unexpected twist of fate.",
		"Beware the unpredictable nature of vildmagi; its volatile power can turn an ally into an unwitting executioner.",
		"The chaos of wild magic, unbound and untamed, holds a perilous end for those who gamble with its forces.",
		"In a world where vildmagi roams free, even the most cautious adventurer can be undone by magic's whims.",
	},
	"Assassin": {
		"An Assassin's contract is a death sentence, and your name is the ink that seals your fate.",
		"Beware the unseen hand of an Assassin, for its touch is the last you will feel.",
		"An Assassin's blade strikes without warning, leaving behind only the echo of a life cut short.",
		"An Assassin's mark is a promise of a swift and silent end, a fate sealed in blood.",
	},
	"Cleric": {
		"A Cleric's divine wrath is a force to be reckoned with, a fate that awaits those who dare to challenge the gods.",
		"Under a Cleric's gaze, the heavens themselves might open, raining down divine retribution upon the unworthy.",
		"A Cleric's faith is a shield against the darkness, and those who stand against it find themselves consumed by its light.",
		"The divine fury of a Cleric is a testament to the power of the gods, a force that leaves no room for mercy.",
	},
	"Devil": {
		"A pact with a Devil is a bargain with the end, a fate sealed in blood and fire.",
		"The whispers of a Devil promise power, but the price is a soul lost to the infernal flames.",
		"A Devil's contract is a promise of eternal torment, a fate that awaits those who seek power at any cost.",
		"Beware the allure of a Devil's bargain, for the price is a fate worse than death.",
	},
	"Drows": {
		"In the darkness of the Underdark, the Drows' blades are swift and merciless, a fate that awaits those who dare to trespass.",
		"The Drows' whispers in the night are a promise of a fate worse than death, a life lost to the shadows.",
		"A Drow's vengeance is a fate sealed in darkness, a fate that awaits those who cross their path.",
		"Beware the Drows' treachery, for their blades are as sharp as their wit.",
	},
	"Druid": {
		"A Druid's wrath is a force of nature, a fate that awaits those who dare to desecrate the wilds.",
		"In the heart of the forest, a Druid's curse is a promise of a fate entwined with the earth itself.",
		"A Druid's bond with nature is a shield against the darkness, and those who stand against it find themselves consumed by its fury.",
		"The elemental fury of a Druid is a testament to the power of the wilds, a force that leaves no room for mercy.",
	},
	"Elemental": {
		"The fury of an Elemental is a force to be reckoned with, a fate that awaits those who dare to challenge the elements.",
		"Under an Elemental's gaze, the very earth itself might open, swallowing the unworthy whole.",
		"An Elemental's wrath is a shield against the darkness, and those who stand against it find themselves consumed by its fury.",
		"The elemental fury of an Elemental is a testament to the power of the elements, a force that leaves no room for mercy.",
	},
	"Environment": {
		"The unforgiving environment is a force to be reckoned with, a fate that awaits those who dare to challenge the elements.",
		"Under the environment's gaze, the very earth itself might open, swallowing the unworthy whole.",
		"The environment's wrath is a shield against the darkness, and those who stand against it find themselves consumed by its fury.",
		"The elemental fury of the environment is a testament to the power of the elements, a force that leaves no room for mercy.",
	},
	"Fey": {
		"In the Fey's enchanting gaze, reality slips away, leading you down a path of no return, lost to the whims of caprice.",
		"The laughter of the Fey marks the end, a playful jest with a fatal conclusion, as you fade into the realm of forgotten tales.",
	},
	"Giant": {
		"In the shadow of a Giant, you find your end, crushed beneath a weight too great to bear, a story too short to tell.",
		"A Giant's roar is the last sound you hear, as the world blurs into darkness, your story lost in the echo of their might.",
	},
	"Outsider": {
		"From beyond the veil, an Outsider reaches forth, pulling you into a void from which there is no return.",
		"The touch of an Outsider is cold, its alien thoughts a maze with no exit, leading you into the abyss.",
	},
	"Wizard": {
		"Caught in a Wizard's spell, you are but a pawn in a game of power, your end written in runes of fire and ice.",
		"A Wizard's ambition knows no bounds, and in their quest for knowledge, your fate is but a footnote, quickly forgotten.",
	},
}

var storyStart = []string{
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
