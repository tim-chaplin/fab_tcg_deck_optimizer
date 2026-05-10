package registry

// Pool-exclusion marker interfaces. The registry filters cards and weapons carrying either
// marker out of LegalCards / LegalWeapons before deck construction sees them, so the deck
// builder never has to know about either.

// NotImplemented opts a card or weapon out of the deck-construction pool. Cards whose
// printed effect references mechanics the simulator doesn't model (e.g. Gold / Silver /
// Copper token economies, Landmarks) implement it so random deck generation and mutation
// pools skip them. A NotImplemented card is still valid in pre-built hands — it evaluates
// using its static Attack / Pitch / Defense — but the optimizer won't introduce it into a
// new deck or swap one in via mutation. Orthogonal to NotSilverAgeLegal: a card can be
// format-legal yet unimplemented, banned yet fully implemented, or both.
type NotImplemented interface {
	NotImplemented()
}

// Unplayable opts a card or weapon out of the deck-construction pool with identical
// semantics to NotImplemented. The distinction is intent — NotImplemented means "we
// haven't gotten around to modelling this card's effect"; Unplayable means "this card's
// effect is so weak the optimizer would never pick it even if fully modelled, so don't
// bother". Both still satisfy Card / Weapon and remain valid in pre-built hands; only
// the deck-construction pipeline filters them out.
type Unplayable interface {
	Unplayable()
}
