// Consuming Volition — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've dealt arcane damage this turn, this gets 'When this hits a hero, they discard
// a card.'"
//
// Rider reads TurnState.ArcaneDamageDealt. When set, the "when this hits a hero" discard rider
// fires only if this card's own printed attack is likely to land (1/4/7 per card.LikelyToHit).
// Runechants firing alongside don't count — "this hits" is strictly about this card's damage
// reaching the hero, not separate arcane tokens.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var consumingVolitionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// consumingVolitionApplyRider emits the on-hit discard rider as a sub-line under self's
// chain step when ArcaneDamageDealt is set AND this card's printed attack is likely to
// land on its own. The credit is card.DiscardValue; the line reads "On-hit discarded a
// card (+N)" indented under the parent.
func consumingVolitionApplyRider(s *card.TurnState, self *card.CardState) {
	if !s.ArcaneDamageDealt {
		return
	}
	s.ApplyAndLogRiderOnHit(self, "On-hit discarded a card", card.DiscardValue)
}

type ConsumingVolitionRed struct{}

func (ConsumingVolitionRed) ID() ids.CardID           { return ids.ConsumingVolitionRed }
func (ConsumingVolitionRed) Name() string             { return "Consuming Volition" }
func (ConsumingVolitionRed) Cost(*card.TurnState) int { return 1 }
func (ConsumingVolitionRed) Pitch() int               { return 1 }
func (ConsumingVolitionRed) Attack() int              { return 4 }
func (ConsumingVolitionRed) Defense() int             { return 3 }
func (ConsumingVolitionRed) Types() card.TypeSet      { return consumingVolitionTypes }
func (ConsumingVolitionRed) GoAgain() bool            { return false }
func (ConsumingVolitionRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	consumingVolitionApplyRider(s, self)
}

type ConsumingVolitionYellow struct{}

func (ConsumingVolitionYellow) ID() ids.CardID           { return ids.ConsumingVolitionYellow }
func (ConsumingVolitionYellow) Name() string             { return "Consuming Volition" }
func (ConsumingVolitionYellow) Cost(*card.TurnState) int { return 1 }
func (ConsumingVolitionYellow) Pitch() int               { return 2 }
func (ConsumingVolitionYellow) Attack() int              { return 3 }
func (ConsumingVolitionYellow) Defense() int             { return 3 }
func (ConsumingVolitionYellow) Types() card.TypeSet      { return consumingVolitionTypes }
func (ConsumingVolitionYellow) GoAgain() bool            { return false }
func (ConsumingVolitionYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	consumingVolitionApplyRider(s, self)
}

type ConsumingVolitionBlue struct{}

func (ConsumingVolitionBlue) ID() ids.CardID           { return ids.ConsumingVolitionBlue }
func (ConsumingVolitionBlue) Name() string             { return "Consuming Volition" }
func (ConsumingVolitionBlue) Cost(*card.TurnState) int { return 1 }
func (ConsumingVolitionBlue) Pitch() int               { return 3 }
func (ConsumingVolitionBlue) Attack() int              { return 2 }
func (ConsumingVolitionBlue) Defense() int             { return 3 }
func (ConsumingVolitionBlue) Types() card.TypeSet      { return consumingVolitionTypes }
func (ConsumingVolitionBlue) GoAgain() bool            { return false }
func (ConsumingVolitionBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	consumingVolitionApplyRider(s, self)
}
