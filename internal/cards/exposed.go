// Exposed — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "If you are **marked**, you can't play this. Target attack gets +1{p}. **Mark** the
// defending hero."
//
// Mark plumbing isn't modelled: the optimizer only tracks our own state, not whether the
// opponent has marked us, and the "Mark the defending hero" rider has no in-model effect
// either since no implemented card reads opponent-mark state. Both clauses drop out, which
// makes Exposed effectively "Target attack gets +1{p}" — a strict overcount of the card's
// real value (a marked attacker can't legally play it). Acceptable upper bound; flag for
// retune when marks land.
//
// "Target attack" includes weapon attacks (per the printed wording, no "action card"
// qualifier).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var exposedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type ExposedBlue struct{}

func (ExposedBlue) ID() ids.CardID          { return ids.ExposedBlue }
func (ExposedBlue) Name() string            { return "Exposed" }
func (ExposedBlue) Cost(*sim.TurnState) int { return 0 }
func (ExposedBlue) Pitch() int              { return 3 }
func (ExposedBlue) Attack() int             { return 0 }
func (ExposedBlue) Defense() int            { return 0 }
func (ExposedBlue) Types() card.TypeSet     { return exposedTypes }
func (ExposedBlue) GoAgain() bool           { return false }
func (ExposedBlue) ARTargetAllowed(c sim.Card) bool {
	return c.Types().IsAttack()
}
func (ExposedBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, ExposedBlue{}.ARTargetAllowed, 1)
	s.Log(self, 0)
}
