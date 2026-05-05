// Wage Gold — Generic Action - Attack. Cost 3.
// Text: "**Universal** When this attacks a hero, you may **wager** a Gold token with them."
//
// Types() ORs in sim.Universal() so the keyword applies. The "may" wager opts in only when
// likely-to-hit; the win (a Gold token) resolves on hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var wageGoldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func wageGoldOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	s.LogRider(self, 0, "On-hit won wager")
}

func wageGoldPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, sim.OnHitHandler{Fire: wageGoldOnHit})
}

func (WageGoldRed) CreatesItem() sim.TokenType    { return sim.TokenTypeGold }
func (WageGoldYellow) CreatesItem() sim.TokenType { return sim.TokenTypeGold }
func (WageGoldBlue) CreatesItem() sim.TokenType   { return sim.TokenTypeGold }

type WageGoldRed struct{}

func (WageGoldRed) ID() ids.CardID                             { return ids.WageGoldRed }
func (WageGoldRed) Name() string                               { return "Wage Gold" }
func (WageGoldRed) Cost(*sim.TurnState) int                    { return 3 }
func (WageGoldRed) Pitch() int                                 { return 1 }
func (WageGoldRed) Attack() int                                { return 7 }
func (WageGoldRed) Defense() int                               { return 2 }
func (WageGoldRed) Types() card.TypeSet                        { return wageGoldTypes | sim.Universal() }
func (WageGoldRed) GoAgain() bool                              { return false }
func (WageGoldRed) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }

type WageGoldYellow struct{}

func (WageGoldYellow) ID() ids.CardID                             { return ids.WageGoldYellow }
func (WageGoldYellow) Name() string                               { return "Wage Gold" }
func (WageGoldYellow) Cost(*sim.TurnState) int                    { return 3 }
func (WageGoldYellow) Pitch() int                                 { return 2 }
func (WageGoldYellow) Attack() int                                { return 6 }
func (WageGoldYellow) Defense() int                               { return 2 }
func (WageGoldYellow) Types() card.TypeSet                        { return wageGoldTypes | sim.Universal() }
func (WageGoldYellow) GoAgain() bool                              { return false }
func (WageGoldYellow) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }

type WageGoldBlue struct{}

func (WageGoldBlue) ID() ids.CardID                             { return ids.WageGoldBlue }
func (WageGoldBlue) Name() string                               { return "Wage Gold" }
func (WageGoldBlue) Cost(*sim.TurnState) int                    { return 3 }
func (WageGoldBlue) Pitch() int                                 { return 3 }
func (WageGoldBlue) Attack() int                                { return 5 }
func (WageGoldBlue) Defense() int                               { return 2 }
func (WageGoldBlue) Types() card.TypeSet                        { return wageGoldTypes | sim.Universal() }
func (WageGoldBlue) GoAgain() bool                              { return false }
func (WageGoldBlue) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }
