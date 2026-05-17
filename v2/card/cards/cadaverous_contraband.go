// Cadaverous Contraband — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Cadaverous Contraband hits, you may put a 'non-attack' action card from your graveyard
// on top of your deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func cadaverousContrabandOnHitRecycle(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	if _, ok := ge.RecycleFromGraveyardToTop(isNonAttackAction); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled a non-attack action card to top of deck", 0)
	}
}

func isNonAttackAction(c card.Card) bool { return c.Types(nil).IsNonAttackAction() }

func cadaverousContrabandPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(cadaverousContrabandOnHitRecycle)
}

func (CadaverousContrabandRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cadaverousContrabandPlay(ge, l, self)
}

func (CadaverousContrabandYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cadaverousContrabandPlay(ge, l, self)
}

func (CadaverousContrabandBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cadaverousContrabandPlay(ge, l, self)
}
