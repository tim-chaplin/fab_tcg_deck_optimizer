// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to
// target hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// weepingBattlegroundPlay fires the banish-and-1-arcane rider via banishAuraFromGraveyard.
// No-op when the graveyard has no aura.
func weepingBattlegroundPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	banishAuraFromGraveyard(ge, l, self.Card.DisplayName())
}

func (WeepingBattlegroundRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(ge, l, self)
}

func (WeepingBattlegroundYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(ge, l, self)
}

func (WeepingBattlegroundBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(ge, l, self)
}
