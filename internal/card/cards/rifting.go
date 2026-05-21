// Rifting — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Rifting hits, you may play your next 'non-attack' action card this turn as
// though it were an instant."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func riftingPlay(self *card.CardState) {
	self.RegisterOnHit(riftingOnHit)
}

// riftingOnHit lets the next non-attack action card play as an instant — paying no action
// point — fizzling if no non-attack action follows.
func riftingOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	GrantNextCardInstant(ge, card.IsNonAttackAction)
}

func (RiftingRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    { riftingPlay(self) }
func (RiftingYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) { riftingPlay(self) }
func (RiftingBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   { riftingPlay(self) }
