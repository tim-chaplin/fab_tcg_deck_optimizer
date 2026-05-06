// Water the Seeds — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, your next attack this combat chain with 1 or less base {p} gets +1{p}.
// **Go again**"
//
// Scans TurnState.CardsRemaining for the first attack action card with base power 1 or less and
// credits the +1 assuming it will be played; if no matching attack follows, the rider fizzles.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var waterTheSeedsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// waterTheSeedsIsTarget gates the rider on attacks (action cards or weapon swings — "your
// next attack") with base power 1 or less.
func waterTheSeedsIsTarget(_ *sim.TurnState, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttack() && pc.Card.Attack() <= 1
}

type WaterTheSeedsRed struct{}

func (WaterTheSeedsRed) ID() ids.CardID          { return ids.WaterTheSeedsRed }
func (WaterTheSeedsRed) Name() string            { return "Water the Seeds" }
func (WaterTheSeedsRed) Cost(*sim.TurnState) int { return 1 }
func (WaterTheSeedsRed) Pitch() int              { return 1 }
func (WaterTheSeedsRed) Attack() int             { return 3 }
func (WaterTheSeedsRed) Defense() int            { return 2 }
func (WaterTheSeedsRed) Types() card.TypeSet     { return waterTheSeedsTypes }
func (WaterTheSeedsRed) GoAgain() bool           { return true }
func (WaterTheSeedsRed) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, waterTheSeedsIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type WaterTheSeedsYellow struct{}

func (WaterTheSeedsYellow) ID() ids.CardID          { return ids.WaterTheSeedsYellow }
func (WaterTheSeedsYellow) Name() string            { return "Water the Seeds" }
func (WaterTheSeedsYellow) Cost(*sim.TurnState) int { return 1 }
func (WaterTheSeedsYellow) Pitch() int              { return 2 }
func (WaterTheSeedsYellow) Attack() int             { return 2 }
func (WaterTheSeedsYellow) Defense() int            { return 2 }
func (WaterTheSeedsYellow) Types() card.TypeSet     { return waterTheSeedsTypes }
func (WaterTheSeedsYellow) GoAgain() bool           { return true }
func (WaterTheSeedsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, waterTheSeedsIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type WaterTheSeedsBlue struct{}

func (WaterTheSeedsBlue) ID() ids.CardID          { return ids.WaterTheSeedsBlue }
func (WaterTheSeedsBlue) Name() string            { return "Water the Seeds" }
func (WaterTheSeedsBlue) Cost(*sim.TurnState) int { return 1 }
func (WaterTheSeedsBlue) Pitch() int              { return 3 }
func (WaterTheSeedsBlue) Attack() int             { return 1 }
func (WaterTheSeedsBlue) Defense() int            { return 2 }
func (WaterTheSeedsBlue) Types() card.TypeSet     { return waterTheSeedsTypes }
func (WaterTheSeedsBlue) GoAgain() bool           { return true }
func (WaterTheSeedsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, waterTheSeedsIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
