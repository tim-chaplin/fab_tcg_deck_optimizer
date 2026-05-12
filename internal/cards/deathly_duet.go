// Deathly Duet — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Deathly Duet attacks, if an attack action card was pitched to play it, it gains
// +2{p}. If a 'non-attack' action card was pitched to play it, create 2 Runechant tokens."
//
// Both riders read self.PitchedToPlay independently — one fires per role in the pitched
// cards.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (DeathlyDuetRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	deathlyDuetApplyRiders(g, l, self)
}

func (DeathlyDuetYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	deathlyDuetApplyRiders(g, l, self)
}

func (DeathlyDuetBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	deathlyDuetApplyRiders(g, l, self)
}

// deathlyDuetApplyRiders folds the two pitch-conditional riders into self and state.
// Both riders can stack when self.PitchedToPlay contains both roles.
func deathlyDuetApplyRiders(g card.GameEngine, l card.Logger, self *card.CardState) {
	var attackPitched, nonAttackActionPitched bool
	for _, p := range self.PitchedToPlay {
		t := p.Types(nil)
		if t.Has(card.TypeAttack) {
			attackPitched = true
		}
		if t.IsNonAttackAction() {
			nonAttackActionPitched = true
		}
	}
	if attackPitched {
		self.BonusAttack += 2
	}
	if nonAttackActionPitched {
		g.CreateRunechants(2)
		l.AppendPostTrigger(self.Card.DisplayName(), "Created 2 runechants", 2)
	}
}
