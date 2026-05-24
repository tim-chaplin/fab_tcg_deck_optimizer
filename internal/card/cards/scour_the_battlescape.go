// Scour the Battlescape — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2,
// Blue 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a
// card. If Scour the Battlescape is played from arsenal, it gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func scourTheBattlescapePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	if len(ge.HeldHand()) == 0 {
		return
	}
	cycled := ge.PopHandAt(0)
	ge.AppendToDeck(cycled)
	ge.DrawOne()
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Cycled %s to deck bottom and drew", cycled.DisplayName())
}

func (ScourTheBattlescapeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	scourTheBattlescapePlay(ge, l, self)
}

func (ScourTheBattlescapeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	scourTheBattlescapePlay(ge, l, self)
}

func (ScourTheBattlescapeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	scourTheBattlescapePlay(ge, l, self)
}
