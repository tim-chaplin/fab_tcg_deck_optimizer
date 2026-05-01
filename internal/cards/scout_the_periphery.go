// Scout the Periphery — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at the top card of target hero's deck. The next attack action card you play from
// arsenal this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: deck-peek rider isn't modelled. The +N{p} grant targets an attack action card
// that itself was played from arsenal — scan CardsRemaining for the first attack action with
// CardState.FromArsenal set. Since the arsenal holds at most one card, the grant only fires
// when the arsenal-in card is an attack action queued later in the chain.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var scoutThePeripheryTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// grantNextArsenalAttackActionBonus adds n to the first scheduled attack action played
// from arsenal via its BonusAttack so the buff lands on the buffed card. Fizzles silently
// when no qualifying target follows.
func grantNextArsenalAttackActionBonus(s *sim.TurnState, n int) {
	for _, pc := range s.CardsRemaining {
		if !pc.FromArsenal {
			continue
		}
		if pc.Card.Types().IsAttackAction() {
			pc.BonusAttack += n
			return
		}
	}
}

type ScoutThePeripheryRed struct{}

func (ScoutThePeripheryRed) ID() ids.CardID          { return ids.ScoutThePeripheryRed }
func (ScoutThePeripheryRed) Name() string            { return "Scout the Periphery" }
func (ScoutThePeripheryRed) Cost(*sim.TurnState) int { return 0 }
func (ScoutThePeripheryRed) Pitch() int              { return 1 }
func (ScoutThePeripheryRed) Attack() int             { return 0 }
func (ScoutThePeripheryRed) Defense() int            { return 2 }
func (ScoutThePeripheryRed) Types() card.TypeSet     { return scoutThePeripheryTypes }
func (ScoutThePeripheryRed) GoAgain() bool           { return true }
func (ScoutThePeripheryRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextArsenalAttackActionBonus(s, 3)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type ScoutThePeripheryYellow struct{}

func (ScoutThePeripheryYellow) ID() ids.CardID          { return ids.ScoutThePeripheryYellow }
func (ScoutThePeripheryYellow) Name() string            { return "Scout the Periphery" }
func (ScoutThePeripheryYellow) Cost(*sim.TurnState) int { return 0 }
func (ScoutThePeripheryYellow) Pitch() int              { return 2 }
func (ScoutThePeripheryYellow) Attack() int             { return 0 }
func (ScoutThePeripheryYellow) Defense() int            { return 2 }
func (ScoutThePeripheryYellow) Types() card.TypeSet     { return scoutThePeripheryTypes }
func (ScoutThePeripheryYellow) GoAgain() bool           { return true }
func (ScoutThePeripheryYellow) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextArsenalAttackActionBonus(s, 2)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type ScoutThePeripheryBlue struct{}

func (ScoutThePeripheryBlue) ID() ids.CardID          { return ids.ScoutThePeripheryBlue }
func (ScoutThePeripheryBlue) Name() string            { return "Scout the Periphery" }
func (ScoutThePeripheryBlue) Cost(*sim.TurnState) int { return 0 }
func (ScoutThePeripheryBlue) Pitch() int              { return 3 }
func (ScoutThePeripheryBlue) Attack() int             { return 0 }
func (ScoutThePeripheryBlue) Defense() int            { return 2 }
func (ScoutThePeripheryBlue) Types() card.TypeSet     { return scoutThePeripheryTypes }
func (ScoutThePeripheryBlue) GoAgain() bool           { return true }
func (ScoutThePeripheryBlue) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextArsenalAttackActionBonus(s, 1)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
