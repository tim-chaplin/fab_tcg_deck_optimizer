// Brandish — Generic Action - Attack. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If Brandish hits, your next weapon attack this turn gains +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func brandishPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(brandishOnHit)
}

// brandishOnHit grants +1{p} to the next weapon attack, fizzling if no weapon swing follows.
func brandishOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	GrantNextCardBonusAttack(ge, 1, card.IsWeaponAttack)
}

func (BrandishRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	brandishPlay(ge, l, self)
}

func (BrandishYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	brandishPlay(ge, l, self)
}

func (BrandishBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	brandishPlay(ge, l, self)
}
