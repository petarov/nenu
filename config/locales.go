package config

// Locales A nested map of localized texts
var Locales = map[string]map[string]string{
	// English
	"en_US": {
		"latest":     "latest",
		"archive":    "archive",
		"light_mode": "lights on",
		"dark_mode":  "lights off",
		"Archive":    "Archive",
		"gen_by":     "gen. by",
	},

	// Български
	"bg_BG": {
		"latest":     "последно",
		"archive":    "архив",
		"light_mode": "светни",
		"dark_mode":  "загаси",
		"Archive":    "Архив",
		"gen_by":     "ген. от",
	},
}
