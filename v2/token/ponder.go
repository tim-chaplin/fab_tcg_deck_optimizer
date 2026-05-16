// Ponder token: aura-flavored token created by ge.CreatePonders. Fires at end-of-turn
// (triggertype.EndOfTurn) — for each token in play pops the top of the deck into the hand,
// letting the post-hoc arsenal-promotion step fill an otherwise-empty arsenal slot. Pops
// past deck-end are silently skipped.

package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// ponderTokenName is the canonical display name. The engine matches by CardName when
// bumping an existing entry's Count or reading a count.
const ponderTokenName = "Ponder"

// NewPonder returns a fresh Ponder token aura at count n.
func NewPonder(n int) *aura.Aura {
	return aura.NewFromToken(ponderTokenName, ids.PonderTokenID, triggertype.EndOfTurn, ponderFire, n)
}

// ponderEngine is the slice of engine surface ponderFire reaches for — PonderDrawOne pops
// the deck top into the hand without bumping CardsDrawn. Defined narrowly so this file
// doesn't reference the engine's card-typed PopDeckTop / AppendHandRaw.
type ponderEngine interface {
	PonderDrawOne() bool
}

// ponderAura is the slice of the firing aura's surface ponderFire needs.
type ponderAura interface {
	Count() int
	Destroy(addToGraveyard bool)
}

// ponderFire is the triggertype.EndOfTurn handler shared by every Ponder aura.
func ponderFire(engine, _, ctx any) {
	eng := engine.(ponderEngine)
	a := ctx.(ponderAura)
	for i := 0; i < a.Count(); i++ {
		if !eng.PonderDrawOne() {
			break
		}
	}
	a.Destroy(false)
}
