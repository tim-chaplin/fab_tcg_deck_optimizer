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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func plunderRunOnHitDraw(s card.GameEngine, l card.Logger, t *sim.Trigger, _ *sim.Aura) {
	s.DrawOne()
	l.AppendPostTriggerf(s.TriggeringCard().DisplayName(), 0,
		"%s drew a card on attack-action hit", t.Source.DisplayName())
}

func plunderRunPlay(s card.GameEngine, l card.Logger, self *card.CardState, source card.Card, n int) {
	s.AddTrigger(sim.Trigger{
		Source:      source,
		TriggerType: sim.TriggerHit,
		TypeFilter:  card.TypeSet.IsAttackAction,
		Handler:     plunderRunOnHitDraw,
	})
	if self.FromArsenal {
		GrantNextCardBonusAttack(s, n, IsAttackAction)
	}
}

func (c PlunderRunRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(s, l, self, c, 3)
}

func (c PlunderRunYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(s, l, self, c, 2)
}

func (c PlunderRunBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	plunderRunPlay(s, l, self, c, 1)
}
