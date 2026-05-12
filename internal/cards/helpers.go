// Shared helpers used across multiple card implementations. Each card's own .go file
// contains only that card's structs and its dedicated *Play / *IsTarget code; anything
// reused across cards lives here.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// --- Aura helpers ---

// banishAuraFromGraveyard banishes the first aura-typed card in the graveyard and, on
// success, deals 1 arcane damage from source via g.DealArcaneDamage — which credits
// value, flips ArcaneDamageDealt when the hit lands, and logs the "Dealt 1 arcane damage"
// rider under source. Returns true when an aura was banished, false when graveyard had
// no aura. Callers that also destroy the source card (e.g. Sigil of Silphidae's leave
// trigger) should run this scan BEFORE AddToGraveyard'ing the source so the printed
// "another aura" restriction is satisfied naturally.
func banishAuraFromGraveyard(g card.GameEngine, l card.Logger, source string) bool {
	if _, ok := g.BanishFromGraveyard(func(c card.Card) bool {
		return c.Types(nil).Has(card.TypeAura)
	}); !ok {
		return false
	}
	g.DealArcaneDamage(l, source, 1)
	return true
}

// --- Fragile-aura helpers ---
//
// A fragile aura pays N Runechants on destruction but dies without value if we take damage,
// so its worth collapses to 0 in partitions where we don't block all incoming damage —
// unless we pop the aura ourselves same turn by landing an attack.

// fragileAuraPlay emits the chain step for a fragile-aura card and writes the expected
// payoff as a sub-line under self when fragileAuraValue is non-zero. Auras have Attack=0,
// so LogPlay carries the chain entry; the rider line carries the predicted value.
func fragileAuraPlay(g card.GameEngine, l card.Logger, self *card.CardState, n int, attackActionOnly bool) {
	v := fragileAuraValue(g, n, attackActionOnly)
	if v <= 0 {
		return
	}
	g.AddValue(v)
	l.AppendPostTriggerf(self.Card.DisplayName(), v, "Aura expected to pay %d runechants", v)
}

// fragileAuraValue returns n when the aura is expected to pay out, 0 otherwise. Two paths
// to payment: (1) we land a same-turn attack that pops the aura now, or (2) we block all
// incoming so the aura survives into a future turn.
//
// attackActionOnly gates the same-turn-pop check. Triggers restricted to "attack action
// card" pass true (weapon swings don't qualify); triggers off any damage source pass false.
func fragileAuraValue(g card.GameEngine, n int, attackActionOnly bool) int {
	if popsThisTurn(g, attackActionOnly) {
		return n
	}
	if g.BlockTotal() >= g.IncomingDamage() {
		return n
	}
	return 0
}

// popsThisTurn reports whether any subsequent attacker in the chain is likely to hit — via
// its own attack power, or via the Runechants that'll fire with the first attack after us.
// The first attack/weapon after our play consumes every live Runechant (playSequence zeroes
// them after it fires), so we credit the current Runechant count only to that first slot;
// later attackers see zero in this approximation.
//
// The runechants check passes dominate=false because Runechant damage is ambient arcane
// damage from aura tokens, not a card attack — Dominate is an attack-keyword and doesn't
// apply. The attacker-power check threads pc.EffectiveDominate() so a target with printed
// (or granted) Dominate clears the 5+ bar.
func popsThisTurn(g card.GameEngine, attackActionOnly bool) bool {
	firstAttacker := true
	for _, pc := range g.CardsRemaining() {
		if !qualifiesAsAttacker(pc.Card, attackActionOnly) {
			continue
		}
		runechants := 0
		if firstAttacker {
			runechants = g.RunechantCount()
			firstAttacker = false
		}
		if g.LikelyToHit(pc) || g.LikelyDamageHits(runechants, false) {
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
func markOpponentOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.MarkOpponent()
	l.AppendPostTrigger(self.Card.DisplayName(), "Marked the opposing hero", 0)
}

// --- "Next attack" rider helpers ---

// GrantNextCardBonusAttack adds n to the BonusAttack of the first scheduled card in
// CardsRemaining for which match returns true, then stops. Lands "your next attack" /
// "the next attack action card with X" riders on the buffed card itself so its
// EffectiveAttack folds the bonus into LikelyToHit. Fizzles silently when no qualifying
// target follows. Pass match as a top-level function value to keep the call alloc-free —
// closures capturing s would escape to the heap.
func GrantNextCardBonusAttack(g card.GameEngine, n int, match func(card.GameEngine, *card.CardState) bool) {
	for _, pc := range g.CardsRemaining() {
		if match(g, pc) {
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
// fold the active hero's class into their Types — Wage Gold under Viserai then
// matches TypeRuneblade.
func IsRunebladeAttack(g card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(g).IsRunebladeAttack()
}
