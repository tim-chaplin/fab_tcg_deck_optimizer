package token

// GameEngine is the slice of engine surface the aura-token fire closures need. Holds the
// union of every method any built-in token fire calls — Runechant uses LikelyDamageHits +
// SetArcaneDamageDealt, Ponder uses PonderDrawOne. *gameengine.GameEngine satisfies it
// structurally.
type GameEngine interface {
	LikelyDamageHits(n int, dominate bool) bool
	SetArcaneDamageDealt(bool)
	PonderDrawOne() bool
}

// Aura is the slice of the firing aura's surface the aura-token fire closures need. Both
// built-in aura tokens read Count and call Destroy on completion.
type Aura interface {
	Count() int
	Destroy(addToGraveyard bool)
}
