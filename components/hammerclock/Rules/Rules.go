package Rules

// Rules defines the rules for a specific game, including the name, phases, and whether players are only taking
// one turn (in that case, phases are being ignored).
type Rules struct {
	Name                 string   `json:"name"`
	Phases               []string `json:"phases"`
	OneTurnForAllPlayers bool     `json:"oneTurnForAllPlayers"`
}

// AllRules contains all the rules available in the application
var AllRules = []Rules{
	WarhammerRules,
	KillTeamRules,
	NecromundaRules,
	AgeOfSigmarRules,
	WarcryRules,
	BloodBowlRules,
	BunnyKingdomRules,
	ChessRules,
}

// WarhammerRules Warhammer rules
var WarhammerRules = Rules{
	Name:                 "Warhammer 40K (10th Edition)",
	Phases:               []string{"Command Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Fight Phase", "End Phase"},
	OneTurnForAllPlayers: false,
}

// KillTeamRules Kill Team rules
var KillTeamRules = Rules{
	Name:                 "Kill Team (2021)",
	Phases:               []string{"Initiative Phase", "Movement Phase", "Shooting Phase", "Fight Phase", "Morale Phase"},
	OneTurnForAllPlayers: false,
}

// NecromundaRules Necromunda rules
var NecromundaRules = Rules{
	Name:                 "Necromunda (2022 edition)",
	Phases:               []string{"Recovery Phase", "Action Phase", "End Phase"},
	OneTurnForAllPlayers: false,
}

// AgeOfSigmarRules Age of Sigmar rules
var AgeOfSigmarRules = Rules{
	Name:                 "Age of Sigmar (4th Edition)",
	Phases:               []string{"Start of Turn Phase", "Hero Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Combat Phase", "End of Turn Phase"},
	OneTurnForAllPlayers: false,
}

// WarcryRules Warcry rules
var WarcryRules = Rules{
	Name:                 "Warcry (3rd edition)",
	Phases:               []string{"Set Up Phase", "Players' Phase (activating models alternately)", "End Phase"},
	OneTurnForAllPlayers: false,
}

// BloodBowlRules Blood Bowl rules
var BloodBowlRules = Rules{
	Name:                 "Blood Bowl (2020 edition)",
	Phases:               []string{"Pre-Match Phase", "Kick-Off Phase", "Team Turn (both teams alternate)", "End of Turn Phase", "Post-Match Phase"},
	OneTurnForAllPlayers: false,
}

// BunnyKingdomRules Bunny Kingdom rules
var BunnyKingdomRules = Rules{
	Name: "Bunny Kingdom",
	Phases: []string{"Draft Phase (players select cards)",
		"Build Phase (place cards on the board)",
		"Scoring Phase (calculate points based on card placement)"},
	OneTurnForAllPlayers: false,
}

// ChessRules Chess rules
var ChessRules = Rules{
	Name:                 "Chess",
	Phases:               []string{},
	OneTurnForAllPlayers: true,
}

// RulesetNames returns the names of the rulesets
func RulesetNames(rules []Rules) []string {
	names := make([]string, len(rules))
	for i, ruleset := range rules {
		names[i] = ruleset.Name
	}
	return names
}
