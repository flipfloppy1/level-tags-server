package main

// LevelDifficulty represents the difficulty of a GD level.
type LevelDifficulty uint

const (
	DifficultyUnrated LevelDifficulty = iota
	DifficultyAuto
	DifficultyEasy
	DifficultyNormal
	DifficultyHardish
	DifficultyHard
	DifficultyHarderish
	DifficultyHarder
	DifficultyInsaneish
	DifficultyInsane
	DifficultyDemon
)

func (d LevelDifficulty) String() string {
	switch d {
	case DifficultyUnrated:
		return "unrated"
	case DifficultyAuto:
		return "auto"
	case DifficultyEasy:
		return "easy"
	case DifficultyNormal:
		return "normal"
	case DifficultyHardish:
		return "hard"
	case DifficultyHard:
		return "hard"
	case DifficultyHarderish:
		return "harder"
	case DifficultyHarder:
		return "harder"
	case DifficultyInsaneish:
		return "insane"
	case DifficultyInsane:
		return "insane"
	case DifficultyDemon:
		return "demon"
	}
	return "unrated"
}

// LevelRating represents the rating given for a GD level.
type LevelRating uint

const (
	RatingNone LevelRating = iota + 1
	RatingFeatured
	RatingEpic
	RatingLegendary
	RatingMythic
)

func (r LevelRating) String() string {
	switch r {
	case RatingNone:
		return ""
	case RatingFeatured:
		return "featured"
	case RatingEpic:
		return "epic"
	case RatingLegendary:
		return "legendary"
	case RatingMythic:
		return "mythic"
	}
	return ""
}
