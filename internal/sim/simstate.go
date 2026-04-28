// Package simstate holds process-wide simulation state that card effects read. Leaf package
// (only depends on internal/card) so any card implementation can import it without cycling
// through the hand/deck/cards stack.
package sim

// CurrentHero is the hero playing the current simulation. Set once at the start of a run; card
// effects read profile info like Intelligence without plumbing it through TurnState.
var CurrentHero Hero

// HeroWantsLowerHealth reports whether the current hero opts into the LowerHealthWanter
// marker — the proxy for "this hero's strategy keeps them at lower {h} than the opponent". Cards
// with a "less {h} than an opposing hero" rider credit the rider when this returns true. Returns
// false when no hero is set (tests, startup).
func HeroWantsLowerHealth() bool {
	_, ok := CurrentHero.(LowerHealthWanter)
	return ok
}
