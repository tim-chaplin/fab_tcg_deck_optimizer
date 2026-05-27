// Fool's Gold — Generic Resource. Cost 0. Printed pitch variants: Yellow 2. Attack 0, Defense 0.
//
// Text: "When this is discarded, create a Gold token."
//
// Pure resource card — only meaningfully assigned the Pitch or Held role (no attack power, no
// defense). The OnDiscard rider fires only on hand→graveyard discards (Punish the Weak, Lay
// Low's hand mode, etc.), not on cycle-to-top/bottom or pitch.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (FoolsGoldYellow) OnDiscard(ge card.GameEngine, l card.Logger) {
	ge.CreateGold(1)
	l.AppendPostTrigger("Fool's Gold [Y]", "Discarded → created a Gold", 0)
}

func (FoolsGoldYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
