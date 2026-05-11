// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to target
// hero."
//
// Play routes through banishAuraFromGraveyard: if s.Graveyard has an aura, banish it for 1
// arcane and flip ArcaneDamageDealt. No aura means the banish clause fails and Play returns
// 0 — the printed 3 block still applies via Defense().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// weepingBattlegroundPlay emits the chain step then writes the banish-for-arcane rider as
// a sub-line under self when an aura was successfully banished from the graveyard.
func weepingBattlegroundPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if n := banishAuraFromGraveyard(s); n > 0 {
		s.AddValue(n)
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished an aura, dealt 1 arcane damage", n)
	}
}

func (WeepingBattlegroundRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(s, l, self)
}

func (WeepingBattlegroundYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(s, l, self)
}

func (WeepingBattlegroundBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(s, l, self)
}
