// Package simstate holds process-wide simulation state that card effects read. Leaf package
// (only depends on internal/card) so any card implementation can import it without cycling
// through the hand/deck/cards stack.
package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// CurrentHero is the hero playing the current simulation. Set once at the start of a run; card
// effects read profile info like Intelligence without plumbing it through TurnState.
var CurrentHero Hero

// OptDebug, when true, makes TurnState.Opt print every Opt outcome to stdout
// (input cards, top split, bottom split). Set by fabsim's -debug flag at the top of a
// run. Off in production. Not synchronised — runs that flip it during a parallel section
// can interleave; today the sim is single-goroutine, so a plain bool is fine.
var OptDebug bool

// HeroWantsLowerHealth reports whether the current hero opts into the LowerHealthWanter
// marker — the proxy for "this hero's strategy keeps them at lower {h} than the opponent". Cards
// with a "less {h} than an opposing hero" rider credit the rider when this returns true. Returns
// false when no hero is set (tests, startup).
func HeroWantsLowerHealth() bool {
	_, ok := CurrentHero.(LowerHealthWanter)
	return ok
}

// Universal returns the current hero's class as a TypeSet — the bit a Universal card folds
// into its Types() so class-gated triggers (e.g. Viserai's "Runeblade card" rider) fire on
// Universal cards. Zero when no hero is set.
func Universal() card.TypeSet {
	if CurrentHero == nil {
		return 0
	}
	return card.NewTypeSet(CurrentHero.Class())
}
