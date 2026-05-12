// Runic Fellingsong — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks, you may banish an aura from your graveyard. If you do, deal 1 arcane
// damage to target hero."
//
// Play credits Attack() plus 1 arcane when banishAuraFromGraveyard lands an aura in
// the banished zone.
// No aura in the graveyard → the banish rider fizzles and Play returns just Attack().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// runicFellingsongPlay emits the chain step at printed power, then routes through
// banishAuraFromGraveyard for the banish-and-1-arcane rider (the helper handles value,
// log, and ArcaneDamageDealt). No-op when the graveyard has no aura.
func runicFellingsongPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	banishAuraFromGraveyard(g, l, self.Card.DisplayName())
}

func (RunicFellingsongRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	runicFellingsongPlay(g, l, self)
}

func (RunicFellingsongYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	runicFellingsongPlay(g, l, self)
}

func (RunicFellingsongBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	runicFellingsongPlay(g, l, self)
}
