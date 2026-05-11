// Vigor Rush — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If you have played a 'non-attack' action card this turn, Vigor Rush gains **go again**."
//
// Conditional go-again gated on s.NonAttackActionPlayed() (an O(1) flag the sim maintains as it
// walks the chain).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func vigorRushPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if s.NonAttackActionPlayed() {
		self.GrantedGoAgain = true
	}
}

func (VigorRushRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	vigorRushPlay(s, l, self)
}

func (VigorRushYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	vigorRushPlay(s, l, self)
}

func (VigorRushBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	vigorRushPlay(s, l, self)
}
