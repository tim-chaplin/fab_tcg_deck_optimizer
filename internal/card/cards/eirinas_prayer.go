// Eirina's Prayer — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Reveal the top card of your deck. Prevent the next X arcane damage that would be dealt to
// your hero this turn, where X is 6 minus the pitch value of the card revealed this way."
//
// DefensiveInstant routes the card through the defense partition slot so the chain runner
// runs it at defense time (when ArcaneIncomingDamage is still on the matchup figure). Play
// peeks the deck, computes x = 6 - top.Pitch (or 6 when the deck is empty), and asks the
// engine to prevent that much arcane damage; the engine returns the amount actually
// prevented (clamped at the remaining ArcaneIncomingDamage) and we credit it as Value.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func eirinasPrayerPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	x := 6
	if top, ok := ge.PeekDeck(); ok {
		x -= top.Pitch()
	}
	prevented := ge.PreventArcaneDamage(x)
	if prevented <= 0 {
		return
	}
	ge.AddValue(prevented)
	l.AppendPostTriggerf(self.Card.DisplayName(), prevented, "Prevented %d arcane damage", prevented)
}

func (EirinasPrayerRed) DefensiveInstant() {}
func (EirinasPrayerRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	eirinasPrayerPlay(ge, l, self)
}

func (EirinasPrayerYellow) DefensiveInstant() {}
func (EirinasPrayerYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	eirinasPrayerPlay(ge, l, self)
}

func (EirinasPrayerBlue) DefensiveInstant() {}
func (EirinasPrayerBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	eirinasPrayerPlay(ge, l, self)
}
