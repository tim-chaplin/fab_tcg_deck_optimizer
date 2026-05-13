// Shared helpers used across multiple card implementations.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
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

// --- Fragile-aura helpers ---
//
// A fragile aura pays N Runechants on destruction but dies without value if we take damage,
// so its worth collapses to 0 in partitions where we don't block all incoming damage —
// unless we pop the aura ourselves same turn by landing an attack.

// fragileAuraPlay writes the expected payoff as a sub-line under self when
// fragileAuraValue is non-zero.
func fragileAuraPlay(ge card.GameEngine, l card.Logger, self *card.CardState, n int, attackActionOnly bool) {
	v := fragileAuraValue(ge, n, attackActionOnly)
	if v <= 0 {
		return
	}
	ge.AddValue(v)
	l.AppendPostTriggerf(self.Card.DisplayName(), v, "Aura expected to pay %d runechants", v)
}

// fragileAuraValue returns n when the aura is expected to pay out, 0 otherwise. Two paths
// to payment: (1) we land a same-turn attack that pops the aura now, or (2) we block all
// incoming so the aura survives into a future turn.
//
// attackActionOnly gates the same-turn-pop check. Triggers restricted to "attack action
// card" pass true (weapon swings don't qualify); triggers off any damage source pass false.
func fragileAuraValue(ge card.GameEngine, n int, attackActionOnly bool) int {
	if popsThisTurn(ge, attackActionOnly) {
		return n
	}
	if ge.BlockTotal() >= ge.IncomingDamage() {
		return n
	}
	return 0
}

// popsThisTurn reports whether any subsequent attacker in the chain is likely to hit — via
// its own attack power, or via Runechants that'll fire with the first attack after us. Only
// the first attack/weapon after our play sees live Runechants (it consumes them); later
// attackers see zero in this approximation.
//
// The runechants check passes dominate=false because Runechant damage is ambient arcane,
// not a card attack.
func popsThisTurn(ge card.GameEngine, attackActionOnly bool) bool {
	firstAttacker := true
	for _, pc := range ge.CardsRemaining() {
		if !qualifiesAsAttacker(pc.Card, attackActionOnly) {
			continue
		}
		runechants := 0
		if firstAttacker {
			runechants = ge.RunechantCount()
			firstAttacker = false
		}
		if ge.LikelyToHit(pc) || ge.LikelyDamageHits(runechants, false) {
			return true
		}
	}
	return false
}

func qualifiesAsAttacker(c card.Card, attackActionOnly bool) bool {
	ts := c.Types(nil)
	if attackActionOnly {
		return ts.IsAttackAction()
	}
	return ts.IsAttack()
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

// IsAttack matches any scheduled attack — action card OR weapon swing — for
// "your next attack" wording.
func IsAttack(_ card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(nil).IsAttack()
}

// IsAttackAction matches scheduled attack action cards (excludes weapon swings)
// for "the next attack action card" wording.
func IsAttackAction(_ card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(nil).IsAttackAction()
}

// IsRunebladeAttack matches scheduled Runeblade attacks (action or weapon) for
// "your next Runeblade attack" wording. Engine threaded through so Universal cards
// fold the active hero's class into their Types.
func IsRunebladeAttack(ge card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(ge).IsRunebladeAttack()
}
