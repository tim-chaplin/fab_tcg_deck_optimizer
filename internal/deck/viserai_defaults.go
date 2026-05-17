package deck

// ViseraiDefaults is the equipment + sideboard loadout the persistence layer tops every
// saved Viserai deck up toward — the common slots that ride alongside the optimizer's
// choices regardless of the climb. fabsim's writeDeck applies it before serialising so
// both the JSON and the sibling fabrary .txt always carry the full loadout.
//
// Card names in Sideboard must match Card.DisplayName() format ("Read the Runes [R]")
// since ApplyDefaults dedupes by display name. Counts in [1, SideboardCopyCap].
var ViseraiDefaults = Defaults{
	Equipment: []string{
		"Beckoning Haunt",
		"Blade Beckoner Boots",
		"Blade Beckoner Helm",
		"Blossom of Spring",
	},
	Sideboard: []SideboardDefault{
		{Name: "Crown of Dichotomy", Count: 1},
		{Name: "Nullrune Boots", Count: 1},
		{Name: "Nullrune Gloves", Count: 1},
		{Name: "Runebleed Robe", Count: 1},
		{Name: "Read the Runes [R]", Count: 2},
		{Name: "Reduce to Runechant [R]", Count: 2},
		{Name: "Sigil of Suffering [R]", Count: 2},
	},
}
