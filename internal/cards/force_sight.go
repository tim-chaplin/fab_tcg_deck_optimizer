// Force Sight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p}. If Force Sight is played from
// arsenal, **opt 2**. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// The Opt 2 fires only when this copy was played from arsenal (self.FromArsenal).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// forceSightPlay grants the next attack action +bonus{p}, logs the chain step (Force
// Sight is a non-attack action — no Attack() to apply), and resolves the arsenal-gated
// Opt 2.
func forceSightPlay(s card.GameEngine, l card.Logger, self *card.CardState, bonus int) {
	GrantNextCardBonusAttack(s, bonus, IsAttackAction)
	if self.FromArsenal {
		s.Opt(l, 2)
	}
}

func (ForceSightRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	forceSightPlay(s, l, self, 3)
}

func (ForceSightYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	forceSightPlay(s, l, self, 2)
}

func (ForceSightBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	forceSightPlay(s, l, self, 1)
}
