// Fyendal's Fighting Spirit — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue
// 5. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks or defends, if you have less {h} than an opposing hero, gain 1{h}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// fyendalsFightingSpiritApplyRider emits the 1{h} gain as a sub-line under self when the
// current hero opts into LowerHealthWanter.
func fyendalsFightingSpiritApplyRider(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.HeroWantsLowerHealth() {
		return
	}
	ge.AddValue(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Gained 1 health (lower health than opposing hero)", 1)
}

func (FyendalsFightingSpiritRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fyendalsFightingSpiritApplyRider(ge, l, self)
}

func (FyendalsFightingSpiritYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fyendalsFightingSpiritApplyRider(ge, l, self)
}

func (FyendalsFightingSpiritBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fyendalsFightingSpiritApplyRider(ge, l, self)
}
