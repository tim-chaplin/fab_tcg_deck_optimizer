// Ravenous Rabble — Generic Action - Attack. Cost 0. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, reveal the top card of your deck. This gets -X{p}, where X is the pitch
// value of the card revealed this way. **Go again**"
//
// Peek s.Deck[0].Pitch() and subtract from base power, floored at 0. If the deck is empty, no card
// is revealed so there's no penalty.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// ravenousRabbleApplyDebuff routes the -X{p} self-debuff (X = revealed deck-top pitch)
// through self.BonusAttack so EffectiveAttack and LikelyToHit see the debuffed power; the
// chain step's (+N) reflects the post-clamp result. No deck top means no penalty. Reads
// the deck top so the cacheable bit flips — the debuff size depends on hidden shuffle order.
func ravenousRabbleApplyDebuff(s card.GameEngine, l card.Logger, self *card.CardState) {
	top, ok := s.PeekDeck()
	if !ok {
		return
	}
	self.BonusAttack -= top.Pitch()
}

func (RavenousRabbleRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	ravenousRabbleApplyDebuff(s, l, self)
}

func (RavenousRabbleYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	ravenousRabbleApplyDebuff(s, l, self)
}

func (RavenousRabbleBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	ravenousRabbleApplyDebuff(s, l, self)
}
