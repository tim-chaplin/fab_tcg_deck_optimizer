// Enchanting Melody — Generic Action - Aura. Cost 2. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Go again** If your hero would be dealt damage, instead destroy Enchanting Melody and
// prevent N damage that source would deal. At the beginning of your end phase, destroy Enchanting
// Melody unless you have played a 'non-attack' action card this turn."
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Handler reads a.Count() for the per-variant prevention amount. The aura subscribes to
// DamageAboutToBeTaken | EndOfTurn; firingType dispatches each event to its clause.
//
// Modelling: each DamageAboutToBeTaken fire is a single damage moment for one damage type.
// Physical and arcane fire separately, so EM picks whichever type is active at this fire,
// absorbs up to a.Count() from it, and self-destructs. This enforces the "OR, not AND"
// rule — even when both damage types arrive the same turn, EM only eats the first fire
// and is gone by the time the second fires.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func enchantingMelodyHandler(ge card.GameEngine, l card.Logger, a card.Aura, firingType triggertype.Type) {
	switch firingType {
	case triggertype.DamageAboutToBeTaken:
		var prevented int
		if ge.RemainingUnblockedDamage() > 0 {
			prevented = ge.PreventIncomingDamage(a.Count())
		} else if ge.ArcaneIncomingDamage() > 0 {
			prevented = ge.PreventArcaneDamage(a.Count())
		}
		if prevented > 0 {
			ge.AddValue(prevented)
			l.AppendPostTriggerf(a.CardName(), prevented, "Prevented %d damage", prevented)
		}
		a.Destroy(true)
	case triggertype.EndOfTurn:
		if !ge.NonAttackActionPlayed() {
			a.Destroy(true)
		}
	}
}

func enchantingMelodyPlay(ge card.GameEngine, self *card.CardState, prevent int) {
	ge.CreateAura(self.Card, triggertype.DamageAboutToBeTaken|triggertype.EndOfTurn,
		enchantingMelodyHandler, prevent, false, nil)
}

func (EnchantingMelodyRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	enchantingMelodyPlay(ge, self, 4)
}

func (EnchantingMelodyYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	enchantingMelodyPlay(ge, self, 3)
}

func (EnchantingMelodyBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	enchantingMelodyPlay(ge, self, 2)
}
