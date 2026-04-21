// Shared helpers for fragile-aura cards (Arcane Cussing, Bloodspill Invocation). Both pay out
// N Runechants on destruction but die without value when we take damage, so their worth this
// game collapses to 0 in partitions where we don't block all incoming damage — unless we pop
// the aura ourselves the same turn by landing an attack.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// fragileAuraValue returns n when the aura is expected to pay out, 0 otherwise. Two paths to
// payment: (1) we land a same-turn attack that pops the aura now, or (2) we block all incoming
// so the aura survives opponent's turn into a future turn where it eventually pays.
//
// attackActionOnly gates the same-turn-pop check. Bloodspill's "attack action card you control
// hits" trigger ignores weapon swings, so it passes true. Cussing's looser "deal damage"
// trigger fires off any attacker, so it passes false.
//
// In the 0-bonus branch the aura was destroyed this turn by opponent damage, so the card moves
// to the graveyard immediately (otherwise auras stay in play until a destroy condition fires).
func fragileAuraValue(s *card.TurnState, n int, attackActionOnly bool) int {
	if popsThisTurn(s, attackActionOnly) {
		return n
	}
	if s.BlockTotal >= s.IncomingDamage {
		return n
	}
	if s.Self != nil {
		s.AddToGraveyard(s.Self.Card)
	}
	return 0
}

// popsThisTurn reports whether any subsequent attacker in the chain is likely to hit — via its
// own attack power, or via the Runechants that'll fire with the first attack after us. The
// first attack/weapon after our play consumes every live Runechant (playSequence zeroes them
// after it fires), so we credit the current Runechant count only to that first slot; later
// attackers see zero in this approximation.
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
		if card.LikelyToHit(pc.Card.Attack()) || card.LikelyToHit(runechants) {
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
