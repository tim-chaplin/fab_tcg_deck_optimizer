// Rally the Coast Guard — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this only while this
// card is defending."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var rallyTheCoastGuardTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RallyTheCoastGuardRed struct{}

func (RallyTheCoastGuardRed) ID() ids.CardID          { return ids.RallyTheCoastGuardRed }
func (RallyTheCoastGuardRed) Name() string            { return "Rally the Coast Guard" }
func (RallyTheCoastGuardRed) Cost(*sim.TurnState) int { return 3 }
func (RallyTheCoastGuardRed) Pitch() int              { return 1 }
func (RallyTheCoastGuardRed) Attack() int             { return 7 }
func (RallyTheCoastGuardRed) Defense() int            { return 2 }
func (RallyTheCoastGuardRed) Types() card.TypeSet     { return rallyTheCoastGuardTypes }
func (RallyTheCoastGuardRed) GoAgain() bool           { return false }

// not implemented: defense-time instant activated ability
func (RallyTheCoastGuardRed) NotImplemented() {}
func (c RallyTheCoastGuardRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type RallyTheCoastGuardYellow struct{}

func (RallyTheCoastGuardYellow) ID() ids.CardID          { return ids.RallyTheCoastGuardYellow }
func (RallyTheCoastGuardYellow) Name() string            { return "Rally the Coast Guard" }
func (RallyTheCoastGuardYellow) Cost(*sim.TurnState) int { return 3 }
func (RallyTheCoastGuardYellow) Pitch() int              { return 2 }
func (RallyTheCoastGuardYellow) Attack() int             { return 6 }
func (RallyTheCoastGuardYellow) Defense() int            { return 2 }
func (RallyTheCoastGuardYellow) Types() card.TypeSet     { return rallyTheCoastGuardTypes }
func (RallyTheCoastGuardYellow) GoAgain() bool           { return false }

// not implemented: defense-time instant activated ability
func (RallyTheCoastGuardYellow) NotImplemented() {}
func (c RallyTheCoastGuardYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type RallyTheCoastGuardBlue struct{}

func (RallyTheCoastGuardBlue) ID() ids.CardID          { return ids.RallyTheCoastGuardBlue }
func (RallyTheCoastGuardBlue) Name() string            { return "Rally the Coast Guard" }
func (RallyTheCoastGuardBlue) Cost(*sim.TurnState) int { return 3 }
func (RallyTheCoastGuardBlue) Pitch() int              { return 3 }
func (RallyTheCoastGuardBlue) Attack() int             { return 5 }
func (RallyTheCoastGuardBlue) Defense() int            { return 2 }
func (RallyTheCoastGuardBlue) Types() card.TypeSet     { return rallyTheCoastGuardTypes }
func (RallyTheCoastGuardBlue) GoAgain() bool           { return false }

// not implemented: defense-time instant activated ability
func (RallyTheCoastGuardBlue) NotImplemented() {}
func (c RallyTheCoastGuardBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
