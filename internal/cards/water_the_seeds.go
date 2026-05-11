// Water the Seeds — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, your next attack this combat chain with 1 or less base {p} gets +1{p}.
// **Go again**"
//
// Scans TurnState.CardsRemaining for the first attack action card with base power 1 or less and
// credits the +1 assuming it will be played; if no matching attack follows, the rider fizzles.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// waterTheSeedsIsTarget gates the rider on attacks (action cards or weapon swings — "your
// next attack") with base power 1 or less.
func waterTheSeedsIsTarget(_ card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types().IsAttack() && pc.Card.Attack() <= 1
}

func (WaterTheSeedsRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(s, 1, waterTheSeedsIsTarget)
}

func (WaterTheSeedsYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(s, 1, waterTheSeedsIsTarget)
}

func (WaterTheSeedsBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(s, 1, waterTheSeedsIsTarget)
}
