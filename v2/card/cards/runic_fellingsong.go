// Runic Fellingsong — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks, you may banish an aura from your graveyard. If you do, deal 1 arcane
// damage to target hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// runicFellingsongPlay fires the banish-and-1-arcane rider via banishAuraFromGraveyard.
// No-op when the graveyard has no aura.
func runicFellingsongPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	banishAuraFromGraveyard(ge, l, self.Card.DisplayName())
}

func (RunicFellingsongRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	runicFellingsongPlay(ge, l, self)
}

func (RunicFellingsongYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	runicFellingsongPlay(ge, l, self)
}

func (RunicFellingsongBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	runicFellingsongPlay(ge, l, self)
}
