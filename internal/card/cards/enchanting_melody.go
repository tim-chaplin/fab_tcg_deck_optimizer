// Enchanting Melody — Generic Action - Aura. Cost 2. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Go again** If your hero would be dealt damage, instead destroy Enchanting Melody and
// prevent 4 damage that source would deal. At the beginning of your end phase, destroy Enchanting
// Melody unless you have played a 'non-attack' action card this turn."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func enchantingMelodyHandler(ge card.GameEngine, l card.Logger, a card.Aura, firingType triggertype.Type) {
	switch firingType {
	case triggertype.DamageAboutToBeTaken:
		prevented := ge.PreventIncomingDamage(4)
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

func enchantingMelodyPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self.Card, triggertype.DamageAboutToBeTaken|triggertype.EndOfTurn,
		enchantingMelodyHandler, 1, false, nil)
}

func (EnchantingMelodyRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	enchantingMelodyPlay(ge, l, self)
}

func (EnchantingMelodyYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	enchantingMelodyPlay(ge, l, self)
}

func (EnchantingMelodyBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	enchantingMelodyPlay(ge, l, self)
}
