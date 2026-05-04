// Shared helper for cards whose rider is "the next [filter] you play this turn gets +N{p}".
// Peeks TurnState.CardsRemaining; if any follow-up matches filter, the bonus is added to
// that card's BonusAttack so the +N is attributed to the buffed attack (not the granter)
// and EffectiveAttack picks it up in hit-likelihood checks.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// GrantNextCardBonusAttack adds n to the first scheduled card in CardsRemaining whose
// TypeSet matches filter. Caller passes a card.TypeSet predicate (typically
// card.TypeSet.IsAttack for "your next attack" or card.TypeSet.IsAttackAction for "the
// next attack action card") so the call site states the card-text wording explicitly.
func GrantNextCardBonusAttack(s *sim.TurnState, n int, filter func(card.TypeSet) bool) {
	for _, pc := range s.CardsRemaining {
		if filter(pc.Card.Types()) {
			pc.BonusAttack += n
			return
		}
	}
}
