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

// SetCurrentHero updates sim.CurrentHero (the full Hero interface sim uses for chain-
// runner internals — Opt, OnCardPlayed, Intelligence, …). Cards reach the hero through
// the GameEngine accessors (HeroWantsLowerHealth, CurrentHeroClass) which read it on
// demand; static card methods (GoAgain() on hero-conditional cards) call the
// HeroWantsLowerHealth helper below since they have no engine handle.
func SetCurrentHero(h Hero) {
	CurrentHero = h
}

// HeroWantsLowerHealth reports whether the active hero opts into the LowerHealthWanter
// marker. Static-context helper for card.Card.GoAgain() methods that gate on the hero
// (Life for a Life, Blow for a Blow, Scar for a Scar) — those methods have no
// GameEngine handle, so they call this directly. Cards inside Play / OnHit / Block
// reach the same predicate through GameEngine.HeroWantsLowerHealth().
func HeroWantsLowerHealth() bool {
	_, ok := CurrentHero.(card.LowerHealthWanter)
	return ok
}
