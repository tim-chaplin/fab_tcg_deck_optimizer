// Strike Gold — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, create a Gold token."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var strikeGoldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type StrikeGoldRed struct{}

func (StrikeGoldRed) ID() ids.CardID          { return ids.StrikeGoldRed }
func (StrikeGoldRed) Name() string            { return "Strike Gold" }
func (StrikeGoldRed) Cost(*sim.TurnState) int { return 0 }
func (StrikeGoldRed) Pitch() int              { return 1 }
func (StrikeGoldRed) Attack() int             { return 4 }
func (StrikeGoldRed) Defense() int            { return 2 }
func (StrikeGoldRed) Types() card.TypeSet     { return strikeGoldTypes }
func (StrikeGoldRed) GoAgain() bool           { return false }

// not implemented: gold tokens
func (StrikeGoldRed) NotImplemented() {}
func (StrikeGoldRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, func(state *sim.TurnState) {
		state.AddValue(sim.GoldTokenValue)
		state.LogRider(self, sim.GoldTokenValue, "On-hit created a gold token")
	})
}

type StrikeGoldYellow struct{}

func (StrikeGoldYellow) ID() ids.CardID          { return ids.StrikeGoldYellow }
func (StrikeGoldYellow) Name() string            { return "Strike Gold" }
func (StrikeGoldYellow) Cost(*sim.TurnState) int { return 0 }
func (StrikeGoldYellow) Pitch() int              { return 2 }
func (StrikeGoldYellow) Attack() int             { return 3 }
func (StrikeGoldYellow) Defense() int            { return 2 }
func (StrikeGoldYellow) Types() card.TypeSet     { return strikeGoldTypes }
func (StrikeGoldYellow) GoAgain() bool           { return false }

// not implemented: gold tokens
func (StrikeGoldYellow) NotImplemented() {}
func (StrikeGoldYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, func(state *sim.TurnState) {
		state.AddValue(sim.GoldTokenValue)
		state.LogRider(self, sim.GoldTokenValue, "On-hit created a gold token")
	})
}

type StrikeGoldBlue struct{}

func (StrikeGoldBlue) ID() ids.CardID          { return ids.StrikeGoldBlue }
func (StrikeGoldBlue) Name() string            { return "Strike Gold" }
func (StrikeGoldBlue) Cost(*sim.TurnState) int { return 0 }
func (StrikeGoldBlue) Pitch() int              { return 3 }
func (StrikeGoldBlue) Attack() int             { return 2 }
func (StrikeGoldBlue) Defense() int            { return 2 }
func (StrikeGoldBlue) Types() card.TypeSet     { return strikeGoldTypes }
func (StrikeGoldBlue) GoAgain() bool           { return false }

// not implemented: gold tokens
func (StrikeGoldBlue) NotImplemented() {}
func (StrikeGoldBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, func(state *sim.TurnState) {
		state.AddValue(sim.GoldTokenValue)
		state.LogRider(self, sim.GoldTokenValue, "On-hit created a gold token")
	})
}
