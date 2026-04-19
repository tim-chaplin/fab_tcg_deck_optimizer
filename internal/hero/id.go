package hero

// ID uniquely identifies a hero printing. The zero value (Invalid) is reserved so zero-valued
// IDs in other data structures can be detected as "unset".
//
// IDs are stable within a build but NOT a persistence format: adding or removing heroes may
// renumber existing entries. Treat IDs as opaque in-process handles.
//
// Mirrors card.ID so code keying on (hero, …) tuples can use a fixed-size integer struct
// instead of the display name.
type ID uint16

// Invalid is the sentinel zero value. Valid IDs start at 1.
const Invalid ID = 0

// Hero IDs. Suffixed with "ID" to distinguish from the hero struct types of the same display
// name (e.g. hero.Viserai the struct vs. hero.ViseraiID the ID constant).
const (
	ViseraiID ID = iota + 1
)
