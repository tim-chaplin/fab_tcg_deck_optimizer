// Shared helper for Generic Action cards whose rider is "the next attack action card you play
// this turn gets +N{p}". Peeks TurnState.CardsRemaining; if any follow-up is an attack action,
// the bonus is credited assuming it will be played.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// nextAttackActionBonus returns n when some attack action card is scheduled later this turn,
// otherwise 0. Callers with extra gating (cost/power caps, pitch-color matching) should scan
// CardsRemaining themselves rather than trying to parameterise this helper further.
func nextAttackActionBonus(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if t.Has(card.TypeAttack) && t.Has(card.TypeAction) {
			return n
		}
	}
	return 0
}
