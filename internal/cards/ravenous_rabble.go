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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var ravenousRabbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// ravenousRabbleApplyDebuff routes the -X{p} self-debuff (X = revealed deck-top pitch)
// through self.BonusAttack so EffectiveAttack and LikelyToHit see the debuffed power; the
// chain step's (+N) reflects the post-clamp result. No deck top means no penalty.
func ravenousRabbleApplyDebuff(s *sim.TurnState, self *sim.CardState) {
	if len(s.Deck) == 0 {
		return
	}
	self.BonusAttack -= s.Deck[0].Pitch()
}

type RavenousRabbleRed struct{}

func (RavenousRabbleRed) ID() ids.CardID          { return ids.RavenousRabbleRed }
func (RavenousRabbleRed) Name() string            { return "Ravenous Rabble" }
func (RavenousRabbleRed) Cost(*sim.TurnState) int { return 0 }
func (RavenousRabbleRed) Pitch() int              { return 1 }
func (RavenousRabbleRed) Attack() int             { return 5 }
func (RavenousRabbleRed) Defense() int            { return 2 }
func (RavenousRabbleRed) Types() card.TypeSet     { return ravenousRabbleTypes }
func (RavenousRabbleRed) GoAgain() bool           { return true }
func (RavenousRabbleRed) Play(s *sim.TurnState, self *sim.CardState) {
	ravenousRabbleApplyDebuff(s, self)
	s.ApplyAndLogEffectiveAttack(self)
}

type RavenousRabbleYellow struct{}

func (RavenousRabbleYellow) ID() ids.CardID          { return ids.RavenousRabbleYellow }
func (RavenousRabbleYellow) Name() string            { return "Ravenous Rabble" }
func (RavenousRabbleYellow) Cost(*sim.TurnState) int { return 0 }
func (RavenousRabbleYellow) Pitch() int              { return 2 }
func (RavenousRabbleYellow) Attack() int             { return 4 }
func (RavenousRabbleYellow) Defense() int            { return 2 }
func (RavenousRabbleYellow) Types() card.TypeSet     { return ravenousRabbleTypes }
func (RavenousRabbleYellow) GoAgain() bool           { return true }
func (RavenousRabbleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	ravenousRabbleApplyDebuff(s, self)
	s.ApplyAndLogEffectiveAttack(self)
}

type RavenousRabbleBlue struct{}

func (RavenousRabbleBlue) ID() ids.CardID          { return ids.RavenousRabbleBlue }
func (RavenousRabbleBlue) Name() string            { return "Ravenous Rabble" }
func (RavenousRabbleBlue) Cost(*sim.TurnState) int { return 0 }
func (RavenousRabbleBlue) Pitch() int              { return 3 }
func (RavenousRabbleBlue) Attack() int             { return 3 }
func (RavenousRabbleBlue) Defense() int            { return 2 }
func (RavenousRabbleBlue) Types() card.TypeSet     { return ravenousRabbleTypes }
func (RavenousRabbleBlue) GoAgain() bool           { return true }
func (RavenousRabbleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	ravenousRabbleApplyDebuff(s, self)
	s.ApplyAndLogEffectiveAttack(self)
}
