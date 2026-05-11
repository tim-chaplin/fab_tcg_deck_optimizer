// Shared helpers used across multiple card implementations. Each card's own .go file
// contains only that card's structs and its dedicated *Play / *IsTarget code; anything
// reused across cards lives here.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// --- Aura helpers ---

// banishAuraFromGraveyard banishes the first aura-typed card in the graveyard, flips
// ArcaneDamageDealt, and returns 1. Returns 0 when no aura is found. Callers that also
// destroy the source card (e.g. Sigil of Silphidae's leave trigger) should run this scan
// BEFORE AddToGraveyard'ing the source so the printed "another aura" restriction is
// satisfied naturally. Routes through s.BanishFromGraveyard so the cacheable bit flips —
// reading the graveyard contents (which may include cards from prior turns) makes the
// chain output depend on hidden state. Each caller owns its own AddValue + LogRider
// because the printed wording differs ("Banished an aura, dealt 1 arcane damage").
func banishAuraFromGraveyard(s sim.GameEngine) int {
	if _, ok := s.BanishFromGraveyard(func(c sim.Card) bool {
		return c.Types().Has(card.TypeAura)
	}); !ok {
		return 0
	}
	if sim.LikelyDamageHits(1, false) {
		s.SetArcaneDamageDealt(true)
	}
	return 1
}

// --- Fragile-aura helpers ---
//
// A fragile aura pays N Runechants on destruction but dies without value if we take damage,
// so its worth collapses to 0 in partitions where we don't block all incoming damage —
// unless we pop the aura ourselves same turn by landing an attack.

// fragileAuraPlay emits the chain step for a fragile-aura card and writes the expected
// payoff as a sub-line under self when fragileAuraValue is non-zero. Auras have Attack=0,
// so LogPlay carries the chain entry; the rider line carries the predicted value.
func fragileAuraPlay(s sim.GameEngine, l sim.Logger, self *sim.CardState, n int, attackActionOnly bool) {
	v := fragileAuraValue(s, n, attackActionOnly)
	if v <= 0 {
		return
	}
	s.AddValue(v)
	l.AppendPostTriggerf(self.Card.DisplayName(), v, "Aura expected to pay %d runechants", v)
}

// fragileAuraValue returns n when the aura is expected to pay out, 0 otherwise. Two paths
// to payment: (1) we land a same-turn attack that pops the aura now, or (2) we block all
// incoming so the aura survives into a future turn.
//
// attackActionOnly gates the same-turn-pop check. Triggers restricted to "attack action
// card" pass true (weapon swings don't qualify); triggers off any damage source pass false.
func fragileAuraValue(s sim.GameEngine, n int, attackActionOnly bool) int {
	if popsThisTurn(s, attackActionOnly) {
		return n
	}
	if s.BlockTotal() >= s.IncomingDamage() {
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
func popsThisTurn(s sim.GameEngine, attackActionOnly bool) bool {
	firstAttacker := true
	for _, pc := range s.CardsRemaining() {
		if !qualifiesAsAttacker(pc.Card, attackActionOnly) {
			continue
		}
		runechants := 0
		if firstAttacker {
			runechants = s.Runechants()
			firstAttacker = false
		}
		if sim.LikelyToHit(pc) || sim.LikelyDamageHits(runechants, false) {
			return true
		}
	}
	return false
}

func qualifiesAsAttacker(c sim.Card, attackActionOnly bool) bool {
	ts := c.Types()
	if attackActionOnly {
		return ts.IsAttackAction()
	}
	return ts.IsAttack()
}

// --- Mark helpers ---

// markOpponentOnHit fires the printed "When this hits a hero, mark them" rider.
func markOpponentOnHit(s sim.GameEngine, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.SetOpponentMarked(true)
	l.AppendPostTrigger(self.Card.DisplayName(), "Marked the opposing hero", 0)
}

// --- "Next attack" rider helpers ---

// GrantNextCardBonusAttack adds n to the BonusAttack of the first scheduled card in
// CardsRemaining for which match returns true, then stops. Lands "your next attack" /
// "the next attack action card with X" riders on the buffed card itself so its
// EffectiveAttack folds the bonus into LikelyToHit. Fizzles silently when no qualifying
// target follows. Pass match as a top-level function value to keep the call alloc-free —
// closures capturing s would escape to the heap.
func GrantNextCardBonusAttack(s sim.GameEngine, n int, match func(sim.GameEngine, *sim.CardState) bool) {
	for _, pc := range s.CardsRemaining() {
		if match(s, pc) {
			pc.BonusAttack += n
			return
		}
	}
}

// IsAttack matches any scheduled attack — action card OR weapon swing — for
// "your next attack" wording.
func IsAttack(_ sim.GameEngine, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttack()
}

// IsAttackAction matches scheduled attack action cards (excludes weapon swings)
// for "the next attack action card" wording.
func IsAttackAction(_ sim.GameEngine, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttackAction()
}

// IsRunebladeAttack matches scheduled Runeblade attacks (action or weapon) for
// "your next Runeblade attack" wording.
func IsRunebladeAttack(_ sim.GameEngine, pc *sim.CardState) bool {
	return pc.Card.Types().IsRunebladeAttack()
}
