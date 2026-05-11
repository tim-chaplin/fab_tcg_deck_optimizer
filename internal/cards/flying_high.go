// Flying High — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Your next attack this turn gets **go again**. If it's <matching color>, it gets +1{p}.
// **Go again**" (Red checks for a red attack, Yellow for a yellow attack, Blue for a blue attack.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// flyingHighApplySideEffect grants go again to the next attack scheduled later this turn
// — attack action card OR weapon swing per the "your next attack" wording. If the target
// is an attack action card whose pitch matches matchPitch (this card's own pitch), we
// also add +1 to its BonusAttack — the "+1{p} if it's <matching color>" rider — so
// EffectiveAttack picks the buff up in any LikelyToHit check on the buffed attack. The
// +1 attributes to the target's slot, not Flying High's.
func flyingHighApplySideEffect(s *sim.TurnState, matchPitch int) {
	for _, pc := range s.CardsRemaining() {
		if !pc.Card.Types().IsAttack() {
			continue
		}
		pc.GrantedGoAgain = true
		if pc.Card.Types().IsAttackAction() && pc.Card.Pitch() == matchPitch {
			pc.BonusAttack++
		}
		return
	}
}

func (FlyingHighRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	flyingHighApplySideEffect(s, 1)
}

func (FlyingHighYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	flyingHighApplySideEffect(s, 2)
}

func (FlyingHighBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	flyingHighApplySideEffect(s, 3)
}
