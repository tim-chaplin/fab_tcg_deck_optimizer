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
)

func moneyWhereYaMouthIsWagerOnHit(ge card.GameEngine, l card.Logger, target *card.CardState, h *card.OnHitHandler) {
	ge.CreateGold(1)
	l.AppendPostTriggerf(target.Card.DisplayName(), 0, "%s won wager", h.Source.DisplayName())
}

func moneyWhereYaMouthIsPlay(ge card.GameEngine, l card.Logger, self *card.CardState, source card.Card, n int) {
	GrantNextCardBonusAttack(ge, n, card.IsAttack)
	for _, pc := range ge.CardsRemaining() {
		if pc.EffectiveTypes(ge).IsAttack() {
			pc.OnHit = append(pc.OnHit, card.OnHitHandler{
				Fire:   moneyWhereYaMouthIsWagerOnHit,
				Source: source,
			})
			break
		}
	}
}

func (c MoneyWhereYaMouthIsRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moneyWhereYaMouthIsPlay(ge, l, self, c, 3)
}

func (c MoneyWhereYaMouthIsYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moneyWhereYaMouthIsPlay(ge, l, self, c, 2)
}

func (c MoneyWhereYaMouthIsBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moneyWhereYaMouthIsPlay(ge, l, self, c, 1)
}
