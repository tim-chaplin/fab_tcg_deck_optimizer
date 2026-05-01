// On the Horizon — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Printed
// defense: Red 4, Yellow 3, Blue 2.
//
// Text: "When this defends, look at the top card of your deck."
//
// The deck-peek defend trigger isn't modelled — it surfaces information for the player,
// not a state change the solver can credit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type OnTheHorizonRed struct{}

func (OnTheHorizonRed) ID() ids.CardID          { return ids.OnTheHorizonRed }
func (OnTheHorizonRed) Name() string            { return "On the Horizon" }
func (OnTheHorizonRed) Cost(*sim.TurnState) int { return 0 }
func (OnTheHorizonRed) Pitch() int              { return 1 }
func (OnTheHorizonRed) Attack() int             { return 0 }
func (OnTheHorizonRed) Defense() int            { return 4 }
func (OnTheHorizonRed) Types() card.TypeSet     { return defenseReactionTypes }
func (OnTheHorizonRed) GoAgain() bool           { return false }
func (OnTheHorizonRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.ApplyDefenseValue(self.EffectiveDefense()))
}

type OnTheHorizonYellow struct{}

func (OnTheHorizonYellow) ID() ids.CardID          { return ids.OnTheHorizonYellow }
func (OnTheHorizonYellow) Name() string            { return "On the Horizon" }
func (OnTheHorizonYellow) Cost(*sim.TurnState) int { return 0 }
func (OnTheHorizonYellow) Pitch() int              { return 2 }
func (OnTheHorizonYellow) Attack() int             { return 0 }
func (OnTheHorizonYellow) Defense() int            { return 3 }
func (OnTheHorizonYellow) Types() card.TypeSet     { return defenseReactionTypes }
func (OnTheHorizonYellow) GoAgain() bool           { return false }
func (OnTheHorizonYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.ApplyDefenseValue(self.EffectiveDefense()))
}

type OnTheHorizonBlue struct{}

func (OnTheHorizonBlue) ID() ids.CardID          { return ids.OnTheHorizonBlue }
func (OnTheHorizonBlue) Name() string            { return "On the Horizon" }
func (OnTheHorizonBlue) Cost(*sim.TurnState) int { return 0 }
func (OnTheHorizonBlue) Pitch() int              { return 3 }
func (OnTheHorizonBlue) Attack() int             { return 0 }
func (OnTheHorizonBlue) Defense() int            { return 2 }
func (OnTheHorizonBlue) Types() card.TypeSet     { return defenseReactionTypes }
func (OnTheHorizonBlue) GoAgain() bool           { return false }
func (OnTheHorizonBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.ApplyDefenseValue(self.EffectiveDefense()))
}
