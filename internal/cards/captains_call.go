// Captain's Call — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2. Go again.
//
// Text: "Choose 1; The next attack action card with cost N or less you play this turn gains
// +2{p}. The next attack action card with cost N or less you play this turn gains **go
// again**. **Go again**" (Red N=2, Yellow N=1, Blue N=0.)
//
// Mode 0 buffs the next cost-≤N attack action card +2{p}; mode 1 grants it go again.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var captainsCallTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// captainsCallPlay applies the modal grant to the next cost-≤maxCost attack action card in
// CardsRemaining. Fizzles silently if no follow-up attack action matches.
func captainsCallPlay(s *sim.TurnState, self *sim.CardState, maxCost int) {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Cost(s) > maxCost {
			continue
		}
		switch self.Mode {
		case 0:
			pc.BonusAttack += 2
		case 1:
			pc.GrantedGoAgain = true
		}
		break
	}
	s.Log(self, 0)
}

type CaptainsCallRed struct{}

func (CaptainsCallRed) ID() ids.CardID          { return ids.CaptainsCallRed }
func (CaptainsCallRed) Name() string            { return "Captain's Call" }
func (CaptainsCallRed) Cost(*sim.TurnState) int { return 0 }
func (CaptainsCallRed) Pitch() int              { return 1 }
func (CaptainsCallRed) Attack() int             { return 0 }
func (CaptainsCallRed) Defense() int            { return 2 }
func (CaptainsCallRed) Types() card.TypeSet     { return captainsCallTypes }
func (CaptainsCallRed) GoAgain() bool           { return true }
func (CaptainsCallRed) Modes() int              { return 2 }
func (CaptainsCallRed) Play(s *sim.TurnState, self *sim.CardState) {
	captainsCallPlay(s, self, 2)
}

type CaptainsCallYellow struct{}

func (CaptainsCallYellow) ID() ids.CardID          { return ids.CaptainsCallYellow }
func (CaptainsCallYellow) Name() string            { return "Captain's Call" }
func (CaptainsCallYellow) Cost(*sim.TurnState) int { return 0 }
func (CaptainsCallYellow) Pitch() int              { return 2 }
func (CaptainsCallYellow) Attack() int             { return 0 }
func (CaptainsCallYellow) Defense() int            { return 2 }
func (CaptainsCallYellow) Types() card.TypeSet     { return captainsCallTypes }
func (CaptainsCallYellow) GoAgain() bool           { return true }
func (CaptainsCallYellow) Modes() int              { return 2 }
func (CaptainsCallYellow) Play(s *sim.TurnState, self *sim.CardState) {
	captainsCallPlay(s, self, 1)
}

type CaptainsCallBlue struct{}

func (CaptainsCallBlue) ID() ids.CardID          { return ids.CaptainsCallBlue }
func (CaptainsCallBlue) Name() string            { return "Captain's Call" }
func (CaptainsCallBlue) Cost(*sim.TurnState) int { return 0 }
func (CaptainsCallBlue) Pitch() int              { return 3 }
func (CaptainsCallBlue) Attack() int             { return 0 }
func (CaptainsCallBlue) Defense() int            { return 2 }
func (CaptainsCallBlue) Types() card.TypeSet     { return captainsCallTypes }
func (CaptainsCallBlue) GoAgain() bool           { return true }
func (CaptainsCallBlue) Modes() int              { return 2 }
func (CaptainsCallBlue) Play(s *sim.TurnState, self *sim.CardState) {
	captainsCallPlay(s, self, 0)
}
