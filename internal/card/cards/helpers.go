// Shared helpers used across multiple card implementations.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// --- Aura helpers ---

// banishAuraFromGraveyard banishes the first aura-typed card in the graveyard and, on
// success, deals 1 arcane damage from source. Returns true when an aura was banished.
// Callers that also destroy the source card must run this BEFORE adding the source to
// the graveyard so the printed "another aura" restriction is satisfied naturally.
func banishAuraFromGraveyard(ge card.GameEngine, l card.Logger, source string) bool {
	if _, ok := ge.BanishFromGraveyard(func(c card.Card) bool {
		return c.Types(nil).Has(card.TypeAura)
	}); !ok {
		return false
	}
	ge.DealArcaneDamage(l, source, 1)
	return true
}

// --- Mark helpers ---

// markOpponentOnHit fires the printed "When this hits a hero, mark them" rider.
func markOpponentOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.MarkOpponent()
	l.AppendPostTrigger(self.Card.DisplayName(), "Marked the opposing hero", 0)
}

// --- "Next attack" rider helpers ---

// GrantNextCardBonusAttack adds n to the BonusAttack of the first scheduled card in
// CardsRemaining for which match returns true, then stops. Lands "your next attack" /
// "the next attack action card with X" riders on the buffed card itself so its
// EffectiveAttack folds the bonus into LikelyToHit. Fizzles silently when no qualifying
// target follows.
func GrantNextCardBonusAttack(ge card.GameEngine, n int, match func(card.GameEngine, *card.CardState) bool) {
	for _, pc := range ge.CardsRemaining() {
		if match(ge, pc) {
			pc.BonusAttack += n
			return
		}
	}
}
