// Smashing Good Time — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time an attack action card hits a hero this turn, you may destroy an item they
// control with cost 2 or less. If Smashing Good Time is played from arsenal, the next attack action
// card you play this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: item-destruction rider isn't modelled. The +N{p} grant requires self.FromArsenal;
// when set, scan CardsRemaining for the next attack action card and credit the bonus.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var smashingGoodTimeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SmashingGoodTimeRed struct{}

func (SmashingGoodTimeRed) ID() ids.CardID          { return ids.SmashingGoodTimeRed }
func (SmashingGoodTimeRed) Name() string            { return "Smashing Good Time" }
func (SmashingGoodTimeRed) Cost(*sim.TurnState) int { return 0 }
func (SmashingGoodTimeRed) Pitch() int              { return 1 }
func (SmashingGoodTimeRed) Attack() int             { return 0 }
func (SmashingGoodTimeRed) Defense() int            { return 2 }
func (SmashingGoodTimeRed) Types() card.TypeSet     { return smashingGoodTimeTypes }
func (SmashingGoodTimeRed) GoAgain() bool           { return true }

// not implemented: on-hit item-destruction rider
func (SmashingGoodTimeRed) NotImplemented() {}
func (SmashingGoodTimeRed) Play(s *sim.TurnState, self *sim.CardState) {
	if self.FromArsenal {
		grantNextAttackActionBonus(s, 3)
	}
	s.LogPlay(self)
}

type SmashingGoodTimeYellow struct{}

func (SmashingGoodTimeYellow) ID() ids.CardID          { return ids.SmashingGoodTimeYellow }
func (SmashingGoodTimeYellow) Name() string            { return "Smashing Good Time" }
func (SmashingGoodTimeYellow) Cost(*sim.TurnState) int { return 0 }
func (SmashingGoodTimeYellow) Pitch() int              { return 2 }
func (SmashingGoodTimeYellow) Attack() int             { return 0 }
func (SmashingGoodTimeYellow) Defense() int            { return 2 }
func (SmashingGoodTimeYellow) Types() card.TypeSet     { return smashingGoodTimeTypes }
func (SmashingGoodTimeYellow) GoAgain() bool           { return true }

// not implemented: on-hit item-destruction rider
func (SmashingGoodTimeYellow) NotImplemented() {}
func (SmashingGoodTimeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	if self.FromArsenal {
		grantNextAttackActionBonus(s, 2)
	}
	s.LogPlay(self)
}

type SmashingGoodTimeBlue struct{}

func (SmashingGoodTimeBlue) ID() ids.CardID          { return ids.SmashingGoodTimeBlue }
func (SmashingGoodTimeBlue) Name() string            { return "Smashing Good Time" }
func (SmashingGoodTimeBlue) Cost(*sim.TurnState) int { return 0 }
func (SmashingGoodTimeBlue) Pitch() int              { return 3 }
func (SmashingGoodTimeBlue) Attack() int             { return 0 }
func (SmashingGoodTimeBlue) Defense() int            { return 2 }
func (SmashingGoodTimeBlue) Types() card.TypeSet     { return smashingGoodTimeTypes }
func (SmashingGoodTimeBlue) GoAgain() bool           { return true }

// not implemented: on-hit item-destruction rider
func (SmashingGoodTimeBlue) NotImplemented() {}
func (SmashingGoodTimeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	if self.FromArsenal {
		grantNextAttackActionBonus(s, 1)
	}
	s.LogPlay(self)
}
