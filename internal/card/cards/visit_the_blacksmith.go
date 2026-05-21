// Visit the Blacksmith — Generic Action. Cost 0, Defense 2. Printed only in Blue (pitch 3).
// Text: "Your next sword attack this turn gains +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (VisitTheBlacksmithBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 1, card.IsSwordAttack)
}
