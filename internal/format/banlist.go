package format

// silverAgeBanlist is the set of card names banned in the Silver Age format — the single
// source of truth for Silver Age legality. Update it manually when the official banlist
// changes (it changes infrequently). The authoritative list is the official card-legality
// policy page; the official cardvault API (legal_formats=Silver Age) reflects the same bans
// and can validate this set — a banned name should be absent from that query's results.
//
// Names must match card.Card.Name() exactly: the base printed name, no pitch suffix, and a
// straight ASCII apostrophe (') to match this repo's card-naming convention — a mismatch
// would silently leak a banned card back into the deck-construction pool. Most entries name
// cards we haven't implemented yet; they're kept so the list mirrors the official banlist in
// full. The implemented ones that this set currently gates out of the pool are Sirens of
// Safe Harbor, Plunder Run, Fiddler's Green, and Fate Foreseen (the unplayable/ banned cards
// and Rosetta Thorn are already excluded by their Unplayable / NotImplemented markers).
var silverAgeBanlist = map[string]struct{}{
	"Aether Flare":               {},
	"Aether Ironweave":           {},
	"Ball Lightning":             {},
	"Belittle":                   {},
	"Bonds of Ancestry":          {},
	"Bracers of Belief":          {},
	"Cash In":                    {},
	"Count Your Blessings":       {},
	"Deadwood Dirge":             {},
	"Drone of Brutality":         {},
	"Electromagnetic Somersault": {},
	"Fate Foreseen":              {},
	"Fiddler's Green":            {},
	"Flic Flak":                  {},
	"Goliath Gauntlet":           {},
	"Heartened Cross Strap":      {},
	"Honing Hood":                {},
	"Mask of Three Tails":        {},
	"Nimby":                      {},
	"Old Knocker":                {},
	"Plunder Run":                {},
	"Ragamuffin's Hat":           {},
	"Reality Refractor":          {},
	"Rosetta Thorn":              {},
	"Seeds of Agony":             {},
	"Sigil of Solace":            {},
	"Sink Below":                 {},
	"Sirens of Safe Harbor":      {},
	"Snapdragon Scalers":         {},
	"Steelblade Shunt":           {},
	"Stubby Hammerers":           {},
	"Vest of the First Fist":     {},
	"Vigorous Smashup":           {},
	"Waning Moon":                {},
	"Zephyr Needle":              {},
}
