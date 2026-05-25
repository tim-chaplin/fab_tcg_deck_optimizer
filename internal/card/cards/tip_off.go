// Tip-Off — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Instant** - Discard this: **Mark** target opposing hero."
//
// Modal modelling of the Instant activation: mode 0 is the printed attack (cost 1, full
// power); mode 1 is the Instant-discard activation, modelled as cost 0, zeroes the attack
// damage (via BonusAttack), schedules the mark via an end-of-turn trigger (the chain
// runner's per-attack ClearOpponentMarked would strip an immediate mark before any rider
// could read it; the deferred mark lands after the chain and carries to next turn), and
// grants Go Again so the action point is refunded (so the chain slot is effectively free).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func tipOffMarkAtEndOfTurn(ge card.GameEngine, _ card.Logger, _ card.EphemeralTrigger) {
	ge.MarkOpponent()
}

func tipOffPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 0 {
		return
	}
	self.BonusAttack -= self.Card.Attack()
	ge.CreateTrigger(self.Card, triggertype.EndOfTurn, tipOffMarkAtEndOfTurn, nil)
	self.GrantedGoAgain = true
}

func tipOffModalCost(mode int8) int {
	if mode == 0 {
		return 1
	}
	return 0
}

func (TipOffRed) Modes() int              { return 2 }
func (TipOffRed) ModalCost(mode int8) int { return tipOffModalCost(mode) }
func (TipOffRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tipOffPlay(ge, l, self)
}

func (TipOffYellow) Modes() int              { return 2 }
func (TipOffYellow) ModalCost(mode int8) int { return tipOffModalCost(mode) }
func (TipOffYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tipOffPlay(ge, l, self)
}

func (TipOffBlue) Modes() int              { return 2 }
func (TipOffBlue) ModalCost(mode int8) int { return tipOffModalCost(mode) }
func (TipOffBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tipOffPlay(ge, l, self)
}
