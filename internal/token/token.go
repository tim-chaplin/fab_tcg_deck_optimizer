// Package token owns the factories for FaB's five built-in tokens: the item tokens
// Gold / Silver / Copper, and the aura-flavored tokens Runechant / Ponder. Each factory
// returns the concrete value the engine stores — *item.Item for item tokens, *aura.Aura
// for aura tokens — wiring in the name, identifier, and (for auras) the trigger type +
// fire closure.
package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
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

// NewRunechant returns a fresh Runechant token aura at count n. Fires when an attack is
// played (triggertype.CardOrAbility, filtered to attacks): flips ArcaneDamageDealt when its
// count clears the damage-likely-to-hit window, then destroys. The Runechant resolves
// before the attack's own effect, so it can turn on that attack's "dealt arcane damage
// this turn" rider. Damage is credited at creation time inside CreateRunechants; this
// handler is pure state cleanup.
func NewRunechant(n int) *aura.Aura {
	return aura.NewFromToken("Runechant", ids.RunechantTokenID, triggertype.CardOrAbility,
		func(engine card.GameEngine, _ card.Logger, ctx card.Aura, _ triggertype.Type) {
			if engine.LikelyDamageHits(ctx.Count(), false) {
				engine.(GameEngine).SetArcaneDamageDealt(true)
			}
			ctx.Destroy(false)
		}, n, card.TypeSet.IsAttack)
}

// NewPonder returns a fresh Ponder token aura at count n. Fires at end-of-turn
// (triggertype.EndOfTurn): for each token in play pops the top of the deck into the hand,
// letting the post-hoc arsenal-promotion step fill an otherwise-empty arsenal slot. Pops
// past deck-end are silently skipped.
func NewPonder(n int) *aura.Aura {
	return aura.NewFromToken("Ponder", ids.PonderTokenID, triggertype.EndOfTurn,
		func(engine card.GameEngine, _ card.Logger, ctx card.Aura, _ triggertype.Type) {
			for i := 0; i < ctx.Count(); i++ {
				if !engine.DrawOne() {
					break
				}
			}
			ctx.Destroy(false)
		}, n, nil)
}

// NewQuicken returns a fresh Quicken token aura at count n. Fires when an attack is played
// (triggertype.CardOrAbility, filtered to attacks, which fires before the attack resolves):
// sets GrantedGoAgain on the triggering attack's CardState and consumes one charge.
func NewQuicken(n int) *aura.Aura {
	return aura.NewFromToken("Quicken", ids.QuickenTokenID, triggertype.CardOrAbility,
		func(engine card.GameEngine, _ card.Logger, ctx card.Aura, _ triggertype.Type) {
			pc := engine.TriggeringCard()
			if pc == nil {
				return
			}
			pc.GrantedGoAgain = true
			if ctx.DecrementCount() <= 0 {
				ctx.Destroy(false)
			}
		}, n, card.TypeSet.IsAttack)
}
