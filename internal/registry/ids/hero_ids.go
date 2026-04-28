package ids

// HeroID uniquely identifies a hero printing. The zero value (InvalidHero) is reserved so
// zero-valued IDs in other data structures read as "unset".
//
// HeroIDs are stable within a build but NOT a persistence format: adding or removing heroes
// may renumber existing entries. Treat IDs as opaque in-process handles.
//
// Same width as CardID so (hero, card) tuples stay fixed-size integer structs rather than
// string-keyed by display name.
type HeroID uint16

// InvalidHero is the sentinel zero value. Valid IDs start at 1.
const InvalidHero HeroID = 0

// Hero IDs. Suffixed with "ID" to distinguish from the hero struct types of the same display
// name (e.g. heroes.Viserai the struct vs. ids.ViseraiID the ID constant).
const (
	ViseraiID HeroID = iota + 1
)
