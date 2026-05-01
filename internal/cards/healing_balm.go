// Healing Balm — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2. Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}." (N is the printed variant value above.)
//
// Modelling: health is valued 1-to-1 with damage, so Play credits +N damage-equivalent per
// variant.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var healingBalmTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// healingBalmPlay emits the chain step then writes the printed N{h} as a "Gained N health"
// sub-line under self. Health is valued 1-to-1 with damage.
func healingBalmPlay(s *sim.TurnState, self *sim.CardState, heal int) {
	s.LogChain(self, 0)
	s.LogRiderf(self, s.AddValue(heal), "Gained %d health", heal)
}

type HealingBalmRed struct{}

func (HealingBalmRed) ID() ids.CardID          { return ids.HealingBalmRed }
func (HealingBalmRed) Name() string            { return "Healing Balm" }
func (HealingBalmRed) Cost(*sim.TurnState) int { return 0 }
func (HealingBalmRed) Pitch() int              { return 1 }
func (HealingBalmRed) Attack() int             { return 0 }
func (HealingBalmRed) Defense() int            { return 2 }
func (HealingBalmRed) Types() card.TypeSet     { return healingBalmTypes }
func (HealingBalmRed) GoAgain() bool           { return false }
func (HealingBalmRed) Play(s *sim.TurnState, self *sim.CardState) {
	healingBalmPlay(s, self, 3)
}

type HealingBalmYellow struct{}

func (HealingBalmYellow) ID() ids.CardID          { return ids.HealingBalmYellow }
func (HealingBalmYellow) Name() string            { return "Healing Balm" }
func (HealingBalmYellow) Cost(*sim.TurnState) int { return 0 }
func (HealingBalmYellow) Pitch() int              { return 2 }
func (HealingBalmYellow) Attack() int             { return 0 }
func (HealingBalmYellow) Defense() int            { return 2 }
func (HealingBalmYellow) Types() card.TypeSet     { return healingBalmTypes }
func (HealingBalmYellow) GoAgain() bool           { return false }
func (HealingBalmYellow) Play(s *sim.TurnState, self *sim.CardState) {
	healingBalmPlay(s, self, 2)
}

type HealingBalmBlue struct{}

func (HealingBalmBlue) ID() ids.CardID          { return ids.HealingBalmBlue }
func (HealingBalmBlue) Name() string            { return "Healing Balm" }
func (HealingBalmBlue) Cost(*sim.TurnState) int { return 0 }
func (HealingBalmBlue) Pitch() int              { return 3 }
func (HealingBalmBlue) Attack() int             { return 0 }
func (HealingBalmBlue) Defense() int            { return 2 }
func (HealingBalmBlue) Types() card.TypeSet     { return healingBalmTypes }
func (HealingBalmBlue) GoAgain() bool           { return false }
func (HealingBalmBlue) Play(s *sim.TurnState, self *sim.CardState) {
	healingBalmPlay(s, self, 1)
}
