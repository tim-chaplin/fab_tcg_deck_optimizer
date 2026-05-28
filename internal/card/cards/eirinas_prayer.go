// Eirina's Prayer — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Reveal the top card of your deck. Prevent the next X arcane damage that would be dealt to
// your hero this turn, where X is N minus the pitch value of the card revealed this way."
// Red N=6, Yellow N=5, Blue N=4.
//
// DefensiveInstant routes the card through the defense partition slot so the attack-turn runner
// runs it at defense time (when IncomingArcaneDamage is still on the matchup figure). Play
// peeks the deck, computes x = N - top.Pitch (or N when the deck is empty), and asks the
// engine to prevent that much arcane damage; PreventArcaneDamage credits the prevented
// amount as Value and returns it for the rider log.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func eirinasPrayerPlay(ge card.GameEngine, l card.Logger, self *card.CardState, n int) {
	x := n
	if top, ok := ge.PeekDeck(); ok {
		x -= top.Pitch()
	}
	prevented := ge.PreventArcaneDamage(x)
	if prevented <= 0 {
		return
	}
	l.AppendPostTriggerf(self.Card.DisplayName(), prevented, "Prevented %d arcane damage", prevented)
}

func (EirinasPrayerRed) DefensiveInstant() {}
func (EirinasPrayerRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	eirinasPrayerPlay(ge, l, self, 6)
}

func (EirinasPrayerYellow) DefensiveInstant() {}
func (EirinasPrayerYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	eirinasPrayerPlay(ge, l, self, 5)
}

func (EirinasPrayerBlue) DefensiveInstant() {}
func (EirinasPrayerBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	eirinasPrayerPlay(ge, l, self, 4)
}
