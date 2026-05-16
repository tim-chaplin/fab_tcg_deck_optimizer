// Package token owns the factories for FaB's five built-in tokens: the item tokens
// Gold / Silver / Copper, and the aura-flavored tokens Runechant / Ponder. Each factory
// returns the concrete value the engine stores — *item.Item for the item tokens, *aura.Aura
// for the aura tokens — wiring in the per-token name, identifier, and (for auras) the
// trigger type + inlined fire closure.
//
// The activated-ability card types backing the item tokens (cards.GoldToken et al.) live
// in internal/cards alongside other card implementations. The narrow consumer interfaces
// the aura-token fire closures type-assert against (GameEngine, Aura) live in
// interfaces.go.
package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// NewGold returns a fresh Gold token item at count n.
func NewGold(n int) *item.Item {
	g := cards.GoldToken{}
	return item.NewFromToken(g.Name(), ids.GoldTokenID, g, n)
}

// NewSilver returns a fresh Silver token item at count n.
func NewSilver(n int) *item.Item {
	s := cards.SilverToken{}
	return item.NewFromToken(s.Name(), ids.SilverTokenID, s, n)
}

// NewCopper returns a fresh Copper token item at count n.
func NewCopper(n int) *item.Item {
	c := cards.CopperToken{}
	return item.NewFromToken(c.Name(), ids.CopperTokenID, c, n)
}

// NewRunechant returns a fresh Runechant token aura at count n. Fires per attack
// (triggertype.Attack): flips ArcaneDamageDealt when its count clears the damage-likely-
// to-hit window, then destroys. Damage is credited at creation time inside
// CreateRunechants; this handler is pure state cleanup.
func NewRunechant(n int) *aura.Aura {
	return aura.NewFromToken("Runechant", ids.RunechantTokenID, triggertype.Attack,
		func(engine, _, ctx any) {
			eng := engine.(GameEngine)
			a := ctx.(Aura)
			if eng.LikelyDamageHits(a.Count(), false) {
				eng.SetArcaneDamageDealt(true)
			}
			a.Destroy(false)
		}, n)
}

// NewPonder returns a fresh Ponder token aura at count n. Fires at end-of-turn
// (triggertype.EndOfTurn): for each token in play pops the top of the deck into the hand,
// letting the post-hoc arsenal-promotion step fill an otherwise-empty arsenal slot. Pops
// past deck-end are silently skipped.
func NewPonder(n int) *aura.Aura {
	return aura.NewFromToken("Ponder", ids.PonderTokenID, triggertype.EndOfTurn,
		func(engine, _, ctx any) {
			eng := engine.(GameEngine)
			a := ctx.(Aura)
			for i := 0; i < a.Count(); i++ {
				if !eng.PonderDrawOne() {
					break
				}
			}
			a.Destroy(false)
		}, n)
}
