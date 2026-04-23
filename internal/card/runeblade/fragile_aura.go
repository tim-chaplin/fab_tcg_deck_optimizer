// Shared helpers for fragile-aura cards. A fragile aura pays N Runechants on destruction
// but dies without value if we take damage, so its worth collapses to 0 in partitions where
// we don't block all incoming damage — unless we pop the aura ourselves same turn by landing
// an attack.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// fragileAuraValue returns n when the aura is expected to pay out, 0 otherwise. Two paths
// to payment: (1) we land a same-turn attack that pops the aura now, or (2) we block all
// incoming so the aura survives into a future turn.
//
// attackActionOnly gates the same-turn-pop check. Triggers restricted to "attack action
// card" pass true (weapon swings don't qualify); triggers off any damage source pass false.
func fragileAuraValue(s *card.TurnState, n int, attackActionOnly bool) int {
	if popsThisTurn(s, attackActionOnly) {
		return n
	}
	if s.BlockTotal >= s.IncomingDamage {
		return n
	}
	return 0
}

// popsThisTurn reports whether any subsequent attacker in the chain is likely to hit — via its
// own attack power, or via the Runechants that'll fire with the first attack after us. The
// first attack/weapon after our play consumes every live Runechant (playSequence zeroes them
// after it fires), so we credit the current Runechant count only to that first slot; later
// attackers see zero in this approximation.
//
// The runechants check passes dominate=false because Runechant damage is ambient arcane
// damage from aura tokens, not a card attack — Dominate is an attack-keyword and doesn't
// apply. The attacker-power check threads pc.EffectiveDominate() so a target with printed
// (or granted) Dominate clears the 5+ bar.
func popsThisTurn(s *card.TurnState, attackActionOnly bool) bool {
	firstAttacker := true
	for _, pc := range s.CardsRemaining {
		if !qualifiesAsAttacker(pc.Card, attackActionOnly) {
			continue
		}
		runechants := 0
		if firstAttacker {
			runechants = s.Runechants
			firstAttacker = false
		}
		if card.LikelyToHit(pc.Card.Attack(), pc.EffectiveDominate()) || card.LikelyToHit(runechants, false) {
			return true
		}
	}
	return false
}

func qualifiesAsAttacker(c card.Card, attackActionOnly bool) bool {
	ts := c.Types()
	if attackActionOnly {
		return ts.Has(card.TypeAttack)
	}
	return ts.Has(card.TypeAttack) || ts.Has(card.TypeWeapon)
}
