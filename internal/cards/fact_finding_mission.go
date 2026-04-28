// Fact-Finding Mission — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, you may look at a face-down card in their arsenal or equipment
// zones."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var factFindingMissionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FactFindingMissionRed struct{}

func (FactFindingMissionRed) ID() ids.CardID          { return ids.FactFindingMissionRed }
func (FactFindingMissionRed) Name() string            { return "Fact-Finding Mission" }
func (FactFindingMissionRed) Cost(*sim.TurnState) int { return 2 }
func (FactFindingMissionRed) Pitch() int              { return 1 }
func (FactFindingMissionRed) Attack() int             { return 6 }
func (FactFindingMissionRed) Defense() int            { return 2 }
func (FactFindingMissionRed) Types() card.TypeSet     { return factFindingMissionTypes }
func (FactFindingMissionRed) GoAgain() bool           { return false }

// not implemented: on-hit opponent-arsenal/equipment peek
func (FactFindingMissionRed) NotImplemented() {}
func (FactFindingMissionRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type FactFindingMissionYellow struct{}

func (FactFindingMissionYellow) ID() ids.CardID          { return ids.FactFindingMissionYellow }
func (FactFindingMissionYellow) Name() string            { return "Fact-Finding Mission" }
func (FactFindingMissionYellow) Cost(*sim.TurnState) int { return 2 }
func (FactFindingMissionYellow) Pitch() int              { return 2 }
func (FactFindingMissionYellow) Attack() int             { return 5 }
func (FactFindingMissionYellow) Defense() int            { return 2 }
func (FactFindingMissionYellow) Types() card.TypeSet     { return factFindingMissionTypes }
func (FactFindingMissionYellow) GoAgain() bool           { return false }

// not implemented: on-hit opponent-arsenal/equipment peek
func (FactFindingMissionYellow) NotImplemented() {}
func (FactFindingMissionYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type FactFindingMissionBlue struct{}

func (FactFindingMissionBlue) ID() ids.CardID          { return ids.FactFindingMissionBlue }
func (FactFindingMissionBlue) Name() string            { return "Fact-Finding Mission" }
func (FactFindingMissionBlue) Cost(*sim.TurnState) int { return 2 }
func (FactFindingMissionBlue) Pitch() int              { return 3 }
func (FactFindingMissionBlue) Attack() int             { return 4 }
func (FactFindingMissionBlue) Defense() int            { return 2 }
func (FactFindingMissionBlue) Types() card.TypeSet     { return factFindingMissionTypes }
func (FactFindingMissionBlue) GoAgain() bool           { return false }

// not implemented: on-hit opponent-arsenal/equipment peek
func (FactFindingMissionBlue) NotImplemented() {}
func (FactFindingMissionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
