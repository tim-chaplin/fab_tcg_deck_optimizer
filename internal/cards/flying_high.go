// Flying High — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Your next attack this turn gets **go again**. If it's <matching color>, it gets +1{p}.
// **Go again**" (Red checks for a red attack, Yellow for a yellow attack, Blue for a blue attack.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var flyingHighTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// flyingHighApplySideEffect grants go again to the next attack scheduled later this turn
// — attack action card OR weapon swing per the "your next attack" wording. If the target
// is an attack action card whose pitch matches matchPitch (this card's own pitch), we
// also add +1 to its BonusAttack — the "+1{p} if it's <matching color>" rider — so
// EffectiveAttack picks the buff up in any LikelyToHit check on the buffed attack. The
// +1 attributes to the target's slot, not Flying High's.
func flyingHighApplySideEffect(s *sim.TurnState, matchPitch int) {
	for _, pc := range s.CardsRemaining {
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

type FlyingHighRed struct{}

func (FlyingHighRed) ID() ids.CardID          { return ids.FlyingHighRed }
func (FlyingHighRed) Name() string            { return "Flying High" }
func (FlyingHighRed) Cost(*sim.TurnState) int { return 0 }
func (FlyingHighRed) Pitch() int              { return 1 }
func (FlyingHighRed) Attack() int             { return 0 }
func (FlyingHighRed) Defense() int            { return 2 }
func (FlyingHighRed) Types() card.TypeSet     { return flyingHighTypes }
func (FlyingHighRed) GoAgain() bool           { return true }
func (FlyingHighRed) Play(s *sim.TurnState, self *sim.CardState) {
	flyingHighApplySideEffect(s, 1)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type FlyingHighYellow struct{}

func (FlyingHighYellow) ID() ids.CardID          { return ids.FlyingHighYellow }
func (FlyingHighYellow) Name() string            { return "Flying High" }
func (FlyingHighYellow) Cost(*sim.TurnState) int { return 0 }
func (FlyingHighYellow) Pitch() int              { return 2 }
func (FlyingHighYellow) Attack() int             { return 0 }
func (FlyingHighYellow) Defense() int            { return 2 }
func (FlyingHighYellow) Types() card.TypeSet     { return flyingHighTypes }
func (FlyingHighYellow) GoAgain() bool           { return true }
func (FlyingHighYellow) Play(s *sim.TurnState, self *sim.CardState) {
	flyingHighApplySideEffect(s, 2)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type FlyingHighBlue struct{}

func (FlyingHighBlue) ID() ids.CardID          { return ids.FlyingHighBlue }
func (FlyingHighBlue) Name() string            { return "Flying High" }
func (FlyingHighBlue) Cost(*sim.TurnState) int { return 0 }
func (FlyingHighBlue) Pitch() int              { return 3 }
func (FlyingHighBlue) Attack() int             { return 0 }
func (FlyingHighBlue) Defense() int            { return 2 }
func (FlyingHighBlue) Types() card.TypeSet     { return flyingHighTypes }
func (FlyingHighBlue) GoAgain() bool           { return true }
func (FlyingHighBlue) Play(s *sim.TurnState, self *sim.CardState) {
	flyingHighApplySideEffect(s, 3)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
