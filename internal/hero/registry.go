package hero

// byName resolves a Hero's display name to the concrete Hero value. Built once at init.
// Serialization packages (deckio, fabrary) look up heroes here so the list is maintained in
// one place as new heroes come online.
var byName = func() map[string]Hero {
	all := []Hero{
		Viserai{},
	}
	m := make(map[string]Hero, len(all))
	for _, h := range all {
		m[h.Name()] = h
	}
	return m
}()

// ByName returns the Hero for the given display name. Returns (nil, false) when the name isn't
// registered — callers surface that to the user rather than falling back to a default.
func ByName(name string) (Hero, bool) {
	h, ok := byName[name]
	return h, ok
}
