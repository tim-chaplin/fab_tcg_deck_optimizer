// Smashing Good Time — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time an attack action card hits a hero this turn, you may destroy an item they
// control with cost 2 or less. If Smashing Good Time is played from arsenal, the next attack action
// card you play this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: the on-hit item-destruction rider isn't modelled — sim doesn't track opposing items
// and the value is matchup-dependent. The +N{p} grant requires self.FromArsenal; when set, scan
// CardsRemaining for the next attack action card and credit the bonus to its BonusAttack so
// EffectiveAttack folds it into LikelyToHit on the buffed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SmashingGoodTimeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.FromArsenal {
		GrantNextCardBonusAttack(ge, 3, card.IsAttackAction)
	}
}

func (SmashingGoodTimeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.FromArsenal {
		GrantNextCardBonusAttack(ge, 2, card.IsAttackAction)
	}
}

func (SmashingGoodTimeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.FromArsenal {
		GrantNextCardBonusAttack(ge, 1, card.IsAttackAction)
	}
}
