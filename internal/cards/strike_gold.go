// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var strikeGoldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func strikeGoldOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	s.LogRider(self, 0, "On-hit created a gold token")
}

func strikeGoldPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, sim.OnHitHandler{Fire: strikeGoldOnHit})
}

func (StrikeGoldRed) CreatesItem() sim.TokenType    { return sim.TokenTypeGold }
func (StrikeGoldYellow) CreatesItem() sim.TokenType { return sim.TokenTypeGold }
func (StrikeGoldBlue) CreatesItem() sim.TokenType   { return sim.TokenTypeGold }

type StrikeGoldRed struct{}

func (StrikeGoldRed) ID() ids.CardID                             { return ids.StrikeGoldRed }
func (StrikeGoldRed) Name() string                               { return "Strike Gold" }
func (StrikeGoldRed) Cost(*sim.TurnState) int                    { return 0 }
func (StrikeGoldRed) Pitch() int                                 { return 1 }
func (StrikeGoldRed) Attack() int                                { return 4 }
func (StrikeGoldRed) Defense() int                               { return 2 }
func (StrikeGoldRed) Types() card.TypeSet                        { return strikeGoldTypes }
func (StrikeGoldRed) GoAgain() bool                              { return false }
func (StrikeGoldRed) Play(s *sim.TurnState, self *sim.CardState) { strikeGoldPlay(s, self) }

type StrikeGoldYellow struct{}

func (StrikeGoldYellow) ID() ids.CardID                             { return ids.StrikeGoldYellow }
func (StrikeGoldYellow) Name() string                               { return "Strike Gold" }
func (StrikeGoldYellow) Cost(*sim.TurnState) int                    { return 0 }
func (StrikeGoldYellow) Pitch() int                                 { return 2 }
func (StrikeGoldYellow) Attack() int                                { return 3 }
func (StrikeGoldYellow) Defense() int                               { return 2 }
func (StrikeGoldYellow) Types() card.TypeSet                        { return strikeGoldTypes }
func (StrikeGoldYellow) GoAgain() bool                              { return false }
func (StrikeGoldYellow) Play(s *sim.TurnState, self *sim.CardState) { strikeGoldPlay(s, self) }

type StrikeGoldBlue struct{}

func (StrikeGoldBlue) ID() ids.CardID                             { return ids.StrikeGoldBlue }
func (StrikeGoldBlue) Name() string                               { return "Strike Gold" }
func (StrikeGoldBlue) Cost(*sim.TurnState) int                    { return 0 }
func (StrikeGoldBlue) Pitch() int                                 { return 3 }
func (StrikeGoldBlue) Attack() int                                { return 2 }
func (StrikeGoldBlue) Defense() int                               { return 2 }
func (StrikeGoldBlue) Types() card.TypeSet                        { return strikeGoldTypes }
func (StrikeGoldBlue) GoAgain() bool                              { return false }
func (StrikeGoldBlue) Play(s *sim.TurnState, self *sim.CardState) { strikeGoldPlay(s, self) }
