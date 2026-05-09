// Money Where Ya Mouth Is — Generic Action. Cost 1.
// Text: "Your next attack this turn gets +N{p} and "When this attacks a hero, you may
// **wager** a Gold token with them."" (Red N=3, Yellow N=2, Blue N=1.)
//
// "Your next attack" includes weapon swings, so the filter is IsAttack. The granted "may"
// wager opts in only when the buffed attack is likely to hit; the win (a Gold token)
// resolves on hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

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

func (c MoneyWhereYaMouthIsRed) Play(s *sim.TurnState, self *sim.CardState) {
	moneyWhereYaMouthIsPlay(s, self, c, 3)
}

func (c MoneyWhereYaMouthIsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	moneyWhereYaMouthIsPlay(s, self, c, 2)
}

func (c MoneyWhereYaMouthIsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	moneyWhereYaMouthIsPlay(s, self, c, 1)
}
