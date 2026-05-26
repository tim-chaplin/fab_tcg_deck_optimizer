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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func plunderRunOnHitDraw(ge card.GameEngine, l card.Logger, t card.EphemeralTrigger, _ triggertype.Type) {
	ge.DrawOne()
	l.AppendPostTriggerf(ge.TriggeringCard().DisplayName(), 0,
		"%s drew a card on attack-action hit", t.CardName())
}

func plunderRunPlay(ge card.GameEngine, l card.Logger, self *card.CardState, n int) {
	ge.CreateTrigger(self.Card, triggertype.Hit, plunderRunOnHitDraw, card.TypeSet.IsAttackAction)
	if self.FromArsenal {
		GrantNextCardBonusAttack(ge, n, card.IsAttackAction)
	}
}

func (PlunderRunRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(ge, l, self, 3)
}

func (PlunderRunYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(ge, l, self, 2)
}

func (PlunderRunBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(ge, l, self, 1)
}
