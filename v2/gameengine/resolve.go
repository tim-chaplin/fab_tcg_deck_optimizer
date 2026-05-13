package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// ResolveChainStep runs card.Play on self and then applies the standard chain-step
// resolution: for an attack-action or weapon-attack, credit self.EffectiveAttack() to
// g.value; for a defense-reaction (or DefensiveInstant), credit the EffectiveDefense
// capped at IncomingDamage and decrement IncomingDamage; for everything else, log (+0).
// The "<DisplayName>: <VERB> (+N)" chain-step entry is appended after Play returns so
// any self-buffs Play applied (e.g. modal +2{p} riders flipping self.BonusAttack) are
// reflected in the displayed delta.
//
// Cards' Play body owns card-specific behaviour: riders that emit rider log lines,
// OnHit registration, conditional self-buffs, sub-card plays. The standard
// printed-attack-deals-damage / DR-blocks-incoming mechanic is the engine's job; cards
// don't reach for DealEffectiveAttack / DealEffectiveDefense or emit the chain step
// themselves.
func (g *GameEngine) ResolveChainStep(l card.Logger, self *card.CardState) {
	self.Card.Play(g, l, self)
	if self.Card.Types(nil).Has(card.TypeAura) {
		g.auraCreated = true
	}
	n := g.chainStepDelta(self)
	l.AppendChainStep(ChainStepText(self), n)
}

// PlayCard implements card.GameEngine.PlayCard. Cards reach this when they resolve another
// card mid-handler (Moon Wish tutoring Sun Kiss into play on go-again).
func (g *GameEngine) PlayCard(l card.Logger, self *card.CardState) {
	g.ResolveChainStep(l, self)
}

// chainStepDelta computes the chain step's display delta and applies the standard
// damage / block side effects. Returns the (+N) value for the log line.
func (g *GameEngine) chainStepDelta(self *card.CardState) int {
	types := self.Card.Types(nil)
	switch {
	case types.IsAttackAction() || types.IsWeaponAttack():
		n := self.EffectiveAttack()
		g.value += n
		return n
	case types.IsDefenseReaction() || isDefensiveInstant(self.Card):
		n := self.EffectiveDefense()
		if n > g.incomingDamage {
			n = g.incomingDamage
		}
		if n < 0 {
			n = 0
		}
		g.incomingDamage -= n
		g.value += n
		return n
	}
	return 0
}

// isDefensiveInstant reports whether c opts into the DR resolution path via the
// DefensiveInstant marker. Centralised here so ResolveChainStep doesn't repeat the
// type-assertion shape.
func isDefensiveInstant(c card.Card) bool {
	_, ok := c.(card.DefensiveInstant)
	return ok
}

// ChainStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix the chain-step
// log line is built from. VERB picks WEAPON ATTACK for weapon activated-ability cards
// (Weapon + Attack), ATTACK for attack-action cards, DEFENSE REACTION for Defense
// Reactions, and PLAY for everything else; the "from arsenal" suffix tags entries played
// out of the arsenal slot. Declared as a var so internal/optimizations can swap in a
// memoised per-(CardID, FromArsenal) implementation at init.
var ChainStepText = func(self *card.CardState) string {
	types := self.Card.Types(nil)
	var verb string
	switch {
	case types.IsWeaponAttack():
		verb = "WEAPON ATTACK"
	case types.IsAttackAction():
		verb = "ATTACK"
	case types.IsDefenseReaction():
		verb = "DEFENSE REACTION"
	default:
		verb = "PLAY"
	}
	if self.FromArsenal {
		verb += " from arsenal"
	}
	return self.Card.DisplayName() + ": " + verb
}
