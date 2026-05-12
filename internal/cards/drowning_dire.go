// Drowning Dire — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
//
// Text: "If you have played or created an aura this turn, Drowning Dire gains **dominate**.
//
// When Drowning Dire hits, you may put a 'non-attack' action card from your graveyard on the
// bottom of your deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func drowningDireOnHitRecycle(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	if _, ok := g.RecycleFromGraveyardToBottom(isNonAttackAction); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled a non-attack action card to bottom of deck", 0)
	}
}

func drowningDirePlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.AuraCreated() {
		self.GrantedDominate = true
	}
	self.RegisterOnHit(drowningDireOnHitRecycle)
}

func (DrowningDireRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	drowningDirePlay(g, l, self)
}

func (DrowningDireYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	drowningDirePlay(g, l, self)
}

func (DrowningDireBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	drowningDirePlay(g, l, self)
}
