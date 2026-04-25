package weapon

// Weapon roster + name lookup. Weapons are equipment, not deck cards, so they aren't ID-indexed
// like cards; this registry is a small standalone lookup keyed by display name. Lives alongside
// the Weapon implementations so the roster and the impls are maintained together.

// All lists every implemented weapon. Used by deck-search code to enumerate loadouts.
var All = []Weapon{
	NebulaBlade{},
	ReapingBlade{},
	ScepterOfPain{},
	Talishar{},
}

// byName maps Weapon.Name() → Weapon for reverse lookup. Built once at init.
var byName = func() map[string]Weapon {
	m := make(map[string]Weapon, len(All))
	for _, w := range All {
		m[w.Name()] = w
	}
	return m
}()

// ByName returns the registered Weapon for the given display name. Returns (nil, false) when no
// such weapon exists — serialization callers surface that to the user.
func ByName(name string) (Weapon, bool) {
	w, ok := byName[name]
	return w, ok
}
