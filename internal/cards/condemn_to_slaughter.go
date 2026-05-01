// Condemn to Slaughter — Runeblade Action. Cost 1, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gets +N{p}. You may destroy an aura you control. If
// you do, each opponent destroys an aura permanent they control. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var condemnToSlaughterTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type CondemnToSlaughterRed struct{}

func (CondemnToSlaughterRed) ID() ids.CardID          { return ids.CondemnToSlaughterRed }
func (CondemnToSlaughterRed) Name() string            { return "Condemn to Slaughter" }
func (CondemnToSlaughterRed) Cost(*sim.TurnState) int { return 1 }
func (CondemnToSlaughterRed) Pitch() int              { return 1 }
func (CondemnToSlaughterRed) Attack() int             { return 0 }
func (CondemnToSlaughterRed) Defense() int            { return 3 }
func (CondemnToSlaughterRed) Types() card.TypeSet     { return condemnToSlaughterTypes }
func (CondemnToSlaughterRed) GoAgain() bool           { return true }

// not implemented: aura-trade rider and opponent-aura destruction clause; only same-turn
// Runeblade-attack +N{p} is modelled
func (CondemnToSlaughterRed) NotImplemented() {}
func (CondemnToSlaughterRed) Play(s *sim.TurnState, self *sim.CardState) {
	condemnToSlaughterApplySideEffect(s, 3)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type CondemnToSlaughterYellow struct{}

func (CondemnToSlaughterYellow) ID() ids.CardID          { return ids.CondemnToSlaughterYellow }
func (CondemnToSlaughterYellow) Name() string            { return "Condemn to Slaughter" }
func (CondemnToSlaughterYellow) Cost(*sim.TurnState) int { return 1 }
func (CondemnToSlaughterYellow) Pitch() int              { return 2 }
func (CondemnToSlaughterYellow) Attack() int             { return 0 }
func (CondemnToSlaughterYellow) Defense() int            { return 3 }
func (CondemnToSlaughterYellow) Types() card.TypeSet     { return condemnToSlaughterTypes }
func (CondemnToSlaughterYellow) GoAgain() bool           { return true }

// not implemented: aura-trade rider and opponent-aura destruction clause; only same-turn
// Runeblade-attack +N{p} is modelled
func (CondemnToSlaughterYellow) NotImplemented() {}
func (CondemnToSlaughterYellow) Play(s *sim.TurnState, self *sim.CardState) {
	condemnToSlaughterApplySideEffect(s, 2)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type CondemnToSlaughterBlue struct{}

func (CondemnToSlaughterBlue) ID() ids.CardID          { return ids.CondemnToSlaughterBlue }
func (CondemnToSlaughterBlue) Name() string            { return "Condemn to Slaughter" }
func (CondemnToSlaughterBlue) Cost(*sim.TurnState) int { return 1 }
func (CondemnToSlaughterBlue) Pitch() int              { return 3 }
func (CondemnToSlaughterBlue) Attack() int             { return 0 }
func (CondemnToSlaughterBlue) Defense() int            { return 3 }
func (CondemnToSlaughterBlue) Types() card.TypeSet     { return condemnToSlaughterTypes }
func (CondemnToSlaughterBlue) GoAgain() bool           { return true }

// not implemented: aura-trade rider and opponent-aura destruction clause; only same-turn
// Runeblade-attack +N{p} is modelled
func (CondemnToSlaughterBlue) NotImplemented() {}
func (CondemnToSlaughterBlue) Play(s *sim.TurnState, self *sim.CardState) {
	condemnToSlaughterApplySideEffect(s, 1)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// condemnToSlaughterApplySideEffect grants +n to the first scheduled Runeblade attack (attack
// action card or weapon swing) via pc.BonusAttack so the buffed attack's EffectiveAttack folds
// the bonus into LikelyToHit and the chain credit lands on the target's slot. Condemn's own
// contribution is zero.
func condemnToSlaughterApplySideEffect(s *sim.TurnState, n int) {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsRunebladeAttack() {
			pc.BonusAttack += n
			return
		}
	}
}
