// Shared helper for Generic Action cards whose rider is "the next attack action card you play
// this turn gets +N{p}". Peeks TurnState.CardsRemaining; if any follow-up is an attack action,
// the bonus is added to that card's BonusAttack so the +N is attributed to the buffed attack
// (not the granter) and EffectiveAttack picks it up in hit-likelihood checks.

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

// grantNextAttackActionBonus adds n to the first scheduled attack action's BonusAttack and
// returns 0 — the granter's own contribution is zero; the +N rides on the target. If no
// attack action follows in CardsRemaining, the grant fizzles silently (no card to land on).
// Callers with extra gating (cost/power caps, pitch-color matching) should scan
// CardsRemaining themselves rather than trying to parameterise this helper further.
func grantNextAttackActionBonus(s *sim.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsAttackAction() {
			pc.BonusAttack += n
			return 0
		}
	}
	return 0
}
