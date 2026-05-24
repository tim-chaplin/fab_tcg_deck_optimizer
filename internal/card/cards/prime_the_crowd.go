// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2. Go again.
//
// Text: "The next attack action card you play this turn gets +N{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**" (Red N=4, Yellow N=3, Blue
// N=2.)
//
// Crowd reactions land only on the active hero in this solver — opposing heroes aren't
// modelled, so cheers/boos targeting them have no observable effect and are skipped.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func primeTheCrowdPlay(ge card.GameEngine, bonus int) {
	if ge.HeroHasType(card.TypeRevered) {
		ge.CrowdCheer()
	}
	if ge.HeroHasType(card.TypeReviled) {
		ge.CrowdBoo()
	}
	GrantNextCardBonusAttack(ge, bonus, card.IsAttackAction)
}

func (PrimeTheCrowdRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	primeTheCrowdPlay(ge, 4)
}

func (PrimeTheCrowdYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	primeTheCrowdPlay(ge, 3)
}

func (PrimeTheCrowdBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	primeTheCrowdPlay(ge, 2)
}
