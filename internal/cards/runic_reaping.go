// Runic Reaping — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains 'When this hits, create N
// Runechant tokens'. If an attack card was pitched to play Runic Reaping, the next Runeblade
// attack action card you play this turn gains +1{p}. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Both riders target Runeblade attack action cards only — weapon swings don't qualify.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// runicReapingTargetMatches accepts Runeblade attack action cards (weapons don't qualify).
func runicReapingTargetMatches(target *card.CardState) bool {
	t := target.Card.Types()
	return t.Has(card.TypeRuneblade) && t.IsAttackAction()
}

func (c RunicReapingRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	runicReapingPlay(s, l, self, c, 3)
}

func (c RunicReapingYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	runicReapingPlay(s, l, self, c, 2)
}

func (c RunicReapingBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	runicReapingPlay(s, l, self, c, 1)
}

// runicReapingPlay buffs the next matching attack +1{p} when an attack card was pitched
// and appends an on-hit n-runechant rider.
func runicReapingPlay(s card.GameEngine, l card.Logger, selfState *card.CardState, source card.Card, n int) {
	var target *card.CardState
	for _, pc := range s.CardsRemaining() {
		if runicReapingTargetMatches(pc) {
			target = pc
			break
		}
	}
	if target == nil {
		return
	}
	for _, p := range selfState.PitchedToPlay {
		if p.Types().Has(card.TypeAttack) {
			target.BonusAttack++
			break
		}
	}
	text := onHitRunechantText[source.ID()]
	target.OnHit = append(target.OnHit, card.OnHitHandler{
		Fire:    onHitCreateRunechants,
		Source:  source,
		LogText: text,
		N:       n,
	})
}
