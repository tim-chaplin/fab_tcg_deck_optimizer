// Tip-Off — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Instant** - Discard this: **Mark** target opposing hero."
//
// Modal modelling of the Instant activation: mode 0 is the printed attack (cost 1, full
// power); mode 1 is the Instant-discard activation, modelled as cost 0, zeroes the attack
// damage (via BonusAttack), marks the opponent, and grants Go Again so the action point is
// refunded (so the chain slot is effectively free). The mark survives the chain-runner's
// per-attack OpponentMarked clear because that clear is gated on the attack having positive
// EffectiveAttack — mode 1's 0-damage swing leaves the mark intact for downstream readers.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func tipOffPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 0 {
		return
	}
	self.BonusAttack -= self.Card.Attack()
	ge.MarkOpponent()
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
