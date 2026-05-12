// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to target
// hero."
//
// Play routes through banishAuraFromGraveyard: if g.Graveyard has an aura, banish it for 1
// arcane and flip ArcaneDamageDealt. No aura means the banish clause fails and Play returns
// 0 — the printed 3 block still applies via Defense().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// weepingBattlegroundPlay emits the chain step then routes through banishAuraFromGraveyard
// for the banish-and-1-arcane rider (the helper handles value, log, and ArcaneDamageDealt).
// No-op when the graveyard has no aura.
func weepingBattlegroundPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	banishAuraFromGraveyard(g, l, self.Card.DisplayName())
}

func (WeepingBattlegroundRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(g, l, self)
}

func (WeepingBattlegroundYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(g, l, self)
}

func (WeepingBattlegroundBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	weepingBattlegroundPlay(g, l, self)
}
