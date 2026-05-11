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
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func moneyWhereYaMouthIsWagerOnHit(s card.GameEngine, l card.Logger, target *card.CardState, h *card.OnHitHandler) {
	s.CreateGold(1)
	l.AppendPostTriggerf(target.Card.DisplayName(), 0, "%s won wager", h.Source.DisplayName())
}

func moneyWhereYaMouthIsPlay(s card.GameEngine, l card.Logger, self *card.CardState, source sim.Card, n int) {
	GrantNextCardBonusAttack(s, n, IsAttack)
	for _, pc := range s.CardsRemaining() {
		if pc.Card.Types().IsAttack() {
			pc.OnHit = append(pc.OnHit, card.OnHitHandler{
				Fire:   moneyWhereYaMouthIsWagerOnHit,
				Source: source,
			})
			break
		}
	}
}

func (c MoneyWhereYaMouthIsRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	moneyWhereYaMouthIsPlay(s, l, self, c, 3)
}

func (c MoneyWhereYaMouthIsYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	moneyWhereYaMouthIsPlay(s, l, self, c, 2)
}

func (c MoneyWhereYaMouthIsBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	moneyWhereYaMouthIsPlay(s, l, self, c, 1)
}
