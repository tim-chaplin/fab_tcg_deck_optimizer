// Shared helpers used across multiple card implementations.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// --- Aura helpers ---

// selfDestructAuraHandler is the start-of-action-phase trigger handler for an Aura whose
// only trigger clause is "destroy this". The leave payoff is the source card's
// OnLeavesArena hook, which the engine fires from DestroyAura.
func selfDestructAuraHandler(_ card.GameEngine, _ card.Logger, a card.Aura) {
	a.Destroy(true)
}

// installSigilSelfDestructAura registers the Sigil-cycle aura: a one-count card-sourced
// aura that destroys itself at the start of the owner's next action phase, deferring all
// payoff to the source card's OnLeavesArena clause. Every Sigil printed with the cycle's
// "At the beginning of your action phase, destroy this" wording installs through here.
func installSigilSelfDestructAura(ge card.GameEngine, source card.Card) {
	ge.CreateAura(source, triggertype.StartOfTurn, selfDestructAuraHandler, 1, false, nil)
}

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

// GrantNextCardInstant flags the first scheduled card in CardsRemaining for which match
// returns true to play as though it were an instant — paying no action point — then stops.
// Lands "play your next X as an instant" riders. Fizzles silently when no target follows.
func GrantNextCardInstant(ge card.GameEngine, match func(card.GameEngine, *card.CardState) bool) {
	for _, pc := range ge.CardsRemaining() {
		if match(ge, pc) {
			pc.GrantedInstant = true
			return
		}
	}
}
