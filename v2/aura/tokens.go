package aura

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// FaB's two aura-flavored tokens. Production code calls ge.CreateRunechants /
// ge.CreatePonders — those bump an existing entry's Count and credit the appropriate
// value. The factories here are for tests / Spec seeding that want a fresh aura without
// the engine-side side effects.

const (
	NameRunechant = "Runechant"
	NamePonder    = "Ponder"
)

// NewRunechant returns a Runechant token aura at count n.
func NewRunechant(n int) *Aura {
	return NewToken(NameRunechant, ids.RunechantTokenID, triggertype.Attack, runechantFire, n)
}

// NewPonder returns a Ponder token aura at count n.
func NewPonder(n int) *Aura {
	return NewToken(NamePonder, ids.PonderTokenID, triggertype.EndOfTurn, ponderFire, n)
}

// runechantEngine is the slice of engine surface the runechant handler needs.
// *gameengine.GameEngine satisfies it structurally.
type runechantEngine interface {
	SetArcaneDamageDealt(bool)
}

// runechantFire is the triggertype.Attack handler shared by every Runechant aura. Fires
// before each attack / weapon swing: flips ArcaneDamageDealt when aura.Count clears the
// damage-likely-to-hit window and destroys the aura. Damage was credited at creation time
// in CreateRunechants — this handler is pure state cleanup.
func runechantFire(engine, _ any, a Ctx) {
	if likelyDamageHits(a.Count(), false) {
		engine.(runechantEngine).SetArcaneDamageDealt(true)
	}
	a.Destroy(false)
}

// ponderEngine is the slice of engine surface the ponder handler needs. PonderDrawOne
// pops the top of the deck into the hand without bumping CardsDrawn — defined narrowly
// so aura doesn't reference the engine's card-typed PopDeckTop / AppendHandRaw.
type ponderEngine interface {
	PonderDrawOne() bool
}

// ponderFire is the triggertype.EndOfTurn handler shared by every Ponder aura. For each
// token in play it pops the deck top into the hand, letting the post-hoc arsenal-promotion
// step fill an otherwise-empty arsenal slot. Pops past deck-end are silently skipped.
func ponderFire(engine, _ any, a Ctx) {
	eng := engine.(ponderEngine)
	for i := 0; i < a.Count(); i++ {
		if !eng.PonderDrawOne() {
			break
		}
	}
	a.Destroy(false)
}

// likelyDamageHits is the runechant-specific copy of the engine's damage-hits heuristic.
// Kept local so this package doesn't reach into gameengine for it. A typical FaB card is
// worth ~3 points; blocking 1/4/7 with a card over-pays, so those amounts slip past.
// Dominate caps the defender to one block, so 5+ damage always lands at least 2.
func likelyDamageHits(n int, dominate bool) bool {
	if dominate && n >= 5 {
		return true
	}
	return n == 1 || n == 4 || n == 7
}
