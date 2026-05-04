package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// GrantNextCardBonusAttack adds n to the first scheduled card in CardsRemaining whose
// TypeSet matches filter. Caller passes the predicate matching the printed wording
// (card.TypeSet.IsAttack for "your next attack", card.TypeSet.IsAttackAction for "the next
// attack action card") so the wording stays at the call site.
func GrantNextCardBonusAttack(s *sim.TurnState, n int, filter func(card.TypeSet) bool) {
	for _, pc := range s.CardsRemaining {
		if filter(pc.Card.Types()) {
			pc.BonusAttack += n
			return
		}
	}
}
