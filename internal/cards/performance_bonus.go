// Performance Bonus — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets **Go
// again**."
//
// Standard played-from-arsenal go-again (docs/dev-standards.md). Gold-token creation isn't
// modelled.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var performanceBonusTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func performanceBonusPlay(s *sim.TurnState, self *sim.CardState) {
	self.GrantGoAgainIfFromArsenal()
	s.ApplyAndLogEffectiveAttack(self)
	s.ApplyAndLogRiderOnHit(self, sim.GoldTokenValue, "On-hit created a gold token")
}

type PerformanceBonusRed struct{}

func (PerformanceBonusRed) ID() ids.CardID          { return ids.PerformanceBonusRed }
func (PerformanceBonusRed) Name() string            { return "Performance Bonus" }
func (PerformanceBonusRed) Cost(*sim.TurnState) int { return 0 }
func (PerformanceBonusRed) Pitch() int              { return 1 }
func (PerformanceBonusRed) Attack() int             { return 3 }
func (PerformanceBonusRed) Defense() int            { return 2 }
func (PerformanceBonusRed) Types() card.TypeSet     { return performanceBonusTypes }
func (PerformanceBonusRed) GoAgain() bool           { return false }

// not implemented: gold tokens
func (PerformanceBonusRed) NotImplemented() {}
func (PerformanceBonusRed) Play(s *sim.TurnState, self *sim.CardState) {
	performanceBonusPlay(s, self)
}

type PerformanceBonusYellow struct{}

func (PerformanceBonusYellow) ID() ids.CardID          { return ids.PerformanceBonusYellow }
func (PerformanceBonusYellow) Name() string            { return "Performance Bonus" }
func (PerformanceBonusYellow) Cost(*sim.TurnState) int { return 0 }
func (PerformanceBonusYellow) Pitch() int              { return 2 }
func (PerformanceBonusYellow) Attack() int             { return 2 }
func (PerformanceBonusYellow) Defense() int            { return 2 }
func (PerformanceBonusYellow) Types() card.TypeSet     { return performanceBonusTypes }
func (PerformanceBonusYellow) GoAgain() bool           { return false }

// not implemented: gold tokens
func (PerformanceBonusYellow) NotImplemented() {}
func (PerformanceBonusYellow) Play(s *sim.TurnState, self *sim.CardState) {
	performanceBonusPlay(s, self)
}

type PerformanceBonusBlue struct{}

func (PerformanceBonusBlue) ID() ids.CardID          { return ids.PerformanceBonusBlue }
func (PerformanceBonusBlue) Name() string            { return "Performance Bonus" }
func (PerformanceBonusBlue) Cost(*sim.TurnState) int { return 0 }
func (PerformanceBonusBlue) Pitch() int              { return 3 }
func (PerformanceBonusBlue) Attack() int             { return 1 }
func (PerformanceBonusBlue) Defense() int            { return 2 }
func (PerformanceBonusBlue) Types() card.TypeSet     { return performanceBonusTypes }
func (PerformanceBonusBlue) GoAgain() bool           { return false }

// not implemented: gold tokens
func (PerformanceBonusBlue) NotImplemented() {}
func (PerformanceBonusBlue) Play(s *sim.TurnState, self *sim.CardState) {
	performanceBonusPlay(s, self)
}

func (PerformanceBonusRed) ConditionalGoAgain()    {}
func (PerformanceBonusYellow) ConditionalGoAgain() {}
func (PerformanceBonusBlue) ConditionalGoAgain()   {}
