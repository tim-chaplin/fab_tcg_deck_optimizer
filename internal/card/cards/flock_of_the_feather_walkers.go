// Flock of the Feather Walkers — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4,
// Blue 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Flock of the Feather Walkers, reveal a card in your hand
// with cost 1 or less. When you attack with Flock of the Feather Walkers, create a Quicken token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// flockOfTheFeatherWalkersPrecondition gates Flock on "some card in hand has effective
// cost ≤ 1". Reads EffectiveCost so a discount-style VariableCost card (Drawn to the Dark
// Dimension at Runechant 2, Amplify the Arknight at Runechant 3, etc.) qualifies when its
// discounted cost actually fits.
func flockOfTheFeatherWalkersPrecondition(ge card.GameEngine) bool {
	return ge.HandHasMatching(func(ge card.GameEngine, pc *card.CardState) bool {
		return pc.EffectiveCost(ge) <= 1
	})
}

func flockOfTheFeatherWalkersPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateQuicken(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created Quicken", 0)
}

func (FlockOfTheFeatherWalkersRed) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return flockOfTheFeatherWalkersPrecondition(ge)
}

func (FlockOfTheFeatherWalkersYellow) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return flockOfTheFeatherWalkersPrecondition(ge)
}

func (FlockOfTheFeatherWalkersBlue) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return flockOfTheFeatherWalkersPrecondition(ge)
}

func (FlockOfTheFeatherWalkersRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	flockOfTheFeatherWalkersPlay(ge, l, self)
}

func (FlockOfTheFeatherWalkersYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	flockOfTheFeatherWalkersPlay(ge, l, self)
}

func (FlockOfTheFeatherWalkersBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	flockOfTheFeatherWalkersPlay(ge, l, self)
}
