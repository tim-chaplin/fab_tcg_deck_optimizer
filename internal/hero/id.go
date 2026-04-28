package hero

// ID uniquely identifies a hero printing. The zero value (Invalid) is reserved so zero-valued
// IDs in other data structures read as "unset".
//
// IDs are stable within a build but NOT a persistence format: adding or removing heroes may
// renumber existing entries. Treat IDs as opaque in-process handles.
//
// Same width as ids.CardID so (hero, card) tuples stay fixed-size integer structs rather than
// string-keyed by display name.
type ID uint16

// Invalid is the sentinel zero value. Valid IDs start at 1.
const Invalid ID = 0

// Hero IDs. Suffixed with "ID" to distinguish from the hero struct types of the same display
// name (e.g. hero.Viserai the struct vs. hero.ViseraiID the ID constant).
const (
	ViseraiID ID = iota + 1
)
