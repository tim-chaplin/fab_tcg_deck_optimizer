// Money Where Ya Mouth Is — Generic Action. Cost 1.
// Text: "Your next attack this turn gets +N{p} and "When this attacks a hero, you may
// **wager** a Gold token with them."" (Red N=3, Yellow N=2, Blue N=1.)
//
// "Your next attack" includes weapon swings, so the filter is IsAttack. The granted "may"
// wager opts in only when the buffed attack is likely to hit; the win (a Gold token)
// resolves on hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var moneyWhereYaMouthIsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

func moneyWhereYaMouthIsWagerOnHit(s *sim.TurnState, target *sim.CardState, h *sim.OnHitHandler) {
	s.CreateGold(1)
	s.LogPostTriggerf(sim.DisplayName(target.Card), 0, "%s won wager", sim.DisplayName(h.Source))
}

func moneyWhereYaMouthIsPlay(s *sim.TurnState, self *sim.CardState, source sim.Card, n int) {
	GrantNextCardBonusAttack(s, n, IsAttack)
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsAttack() {
			pc.OnHit = append(pc.OnHit, sim.OnHitHandler{
				Fire:   moneyWhereYaMouthIsWagerOnHit,
				Source: source,
			})
			break
		}
	}
	s.Log(self, 0)
}

type MoneyWhereYaMouthIsRed struct{}

func (MoneyWhereYaMouthIsRed) ID() ids.CardID          { return ids.MoneyWhereYaMouthIsRed }
func (MoneyWhereYaMouthIsRed) Name() string            { return "Money Where Ya Mouth Is" }
func (MoneyWhereYaMouthIsRed) Cost(*sim.TurnState) int { return 1 }
func (MoneyWhereYaMouthIsRed) Pitch() int              { return 1 }
func (MoneyWhereYaMouthIsRed) Attack() int             { return 0 }
func (MoneyWhereYaMouthIsRed) Defense() int            { return 2 }
func (MoneyWhereYaMouthIsRed) Types() card.TypeSet     { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsRed) GoAgain() bool           { return true }
func (c MoneyWhereYaMouthIsRed) Play(s *sim.TurnState, self *sim.CardState) {
	moneyWhereYaMouthIsPlay(s, self, c, 3)
}

type MoneyWhereYaMouthIsYellow struct{}

func (MoneyWhereYaMouthIsYellow) ID() ids.CardID          { return ids.MoneyWhereYaMouthIsYellow }
func (MoneyWhereYaMouthIsYellow) Name() string            { return "Money Where Ya Mouth Is" }
func (MoneyWhereYaMouthIsYellow) Cost(*sim.TurnState) int { return 1 }
func (MoneyWhereYaMouthIsYellow) Pitch() int              { return 2 }
func (MoneyWhereYaMouthIsYellow) Attack() int             { return 0 }
func (MoneyWhereYaMouthIsYellow) Defense() int            { return 2 }
func (MoneyWhereYaMouthIsYellow) Types() card.TypeSet     { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsYellow) GoAgain() bool           { return true }
func (c MoneyWhereYaMouthIsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	moneyWhereYaMouthIsPlay(s, self, c, 2)
}

type MoneyWhereYaMouthIsBlue struct{}

func (MoneyWhereYaMouthIsBlue) ID() ids.CardID          { return ids.MoneyWhereYaMouthIsBlue }
func (MoneyWhereYaMouthIsBlue) Name() string            { return "Money Where Ya Mouth Is" }
func (MoneyWhereYaMouthIsBlue) Cost(*sim.TurnState) int { return 1 }
func (MoneyWhereYaMouthIsBlue) Pitch() int              { return 3 }
func (MoneyWhereYaMouthIsBlue) Attack() int             { return 0 }
func (MoneyWhereYaMouthIsBlue) Defense() int            { return 2 }
func (MoneyWhereYaMouthIsBlue) Types() card.TypeSet     { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsBlue) GoAgain() bool           { return true }
func (c MoneyWhereYaMouthIsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	moneyWhereYaMouthIsPlay(s, self, c, 1)
}
