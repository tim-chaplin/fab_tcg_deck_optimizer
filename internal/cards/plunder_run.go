// Plunder Run — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time an attack action card you control hits this turn, draw a card. If
// Plunder Run is played from arsenal, the next attack action card you play this turn gains
// +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// The from-arsenal +N{p} grant gates on self.FromArsenal; played from hand, only the
// on-hit-draw trigger registers.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func plunderRunOnHitDraw(g card.GameEngine, l card.Logger, t card.Trigger) {
	g.DrawOne()
	l.AppendPostTriggerf(g.TriggeringCard().DisplayName(), 0,
		"%s drew a card on attack-action hit", t.CardName())
}

func plunderRunPlay(g card.GameEngine, l card.Logger, self *card.CardState, n int) {
	g.AddHitTrigger(self, plunderRunOnHitDraw, card.TypeSet.IsAttackAction)
	if self.FromArsenal {
		GrantNextCardBonusAttack(g, n, IsAttackAction)
	}
}

func (PlunderRunRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(g, l, self, 3)
}

func (PlunderRunYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(g, l, self, 2)
}

func (PlunderRunBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(g, l, self, 1)
}
