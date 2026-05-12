// Package simstate holds process-wide simulation state that card effects read. Leaf package
// (only depends on internal/card) so any card implementation can import it without cycling
// through the hand/deck/cards stack.
package sim

// CurrentHero is the hero playing the current simulation. Set once at the start of a run; card
// effects read profile info like Intelligence without plumbing it through TurnState.
var CurrentHero Hero

// OptDebug, when true, makes TurnState.Opt print every Opt outcome to stdout
// (input cards, top split, bottom split). Set by fabsim's -debug flag at the top of a
// run. Off in production. Not synchronised — runs that flip it during a parallel section
// can interleave; today the sim is single-goroutine, so a plain bool is fine.
var OptDebug bool

// SetCurrentHero updates sim.CurrentHero (the full Hero interface sim uses for chain-
// runner internals — Opt, OnCardPlayed, Intelligence, …). Cards reach the hero
// exclusively through the GameEngine accessors (HeroWantsLowerHealth,
// CurrentHeroClass), which read this on demand.
func SetCurrentHero(h Hero) {
	CurrentHero = h
}
