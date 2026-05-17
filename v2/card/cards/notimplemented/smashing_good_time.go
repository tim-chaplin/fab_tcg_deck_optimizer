// Smashing Good Time — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time an attack action card hits a hero this turn, you may destroy an item they
// control with cost 2 or less. If Smashing Good Time is played from arsenal, the next attack action
// card you play this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: item-destruction rider isn't modelled. The +N{p} grant requires self.FromArsenal;
// when set, scan CardsRemaining for the next attack action card and credit the bonus.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
)

// not implemented: on-hit item-destruction rider

func (SmashingGoodTimeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.FromArsenal {
		cards.GrantNextCardBonusAttack(ge, 3, card.IsAttackAction)
	}
}

// not implemented: on-hit item-destruction rider

func (SmashingGoodTimeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.FromArsenal {
		cards.GrantNextCardBonusAttack(ge, 2, card.IsAttackAction)
	}
}

// not implemented: on-hit item-destruction rider

func (SmashingGoodTimeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.FromArsenal {
		cards.GrantNextCardBonusAttack(ge, 1, card.IsAttackAction)
	}
}
