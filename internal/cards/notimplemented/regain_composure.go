// Regain Composure — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next attack this turn gets +1{p} and "When this hits, {u} your hero." **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: on-hit unfreeze rider (freeze/unfreeze state not tracked)

func (RegainComposureBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cards.GrantNextCardBonusAttack(ge, 1, cards.IsAttack)
}
