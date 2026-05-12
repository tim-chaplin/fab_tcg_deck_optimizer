// Core tuning assumptions of the model. Damage-equivalent values for non-damage riders ("draw
// a card", "create a Gold token", "opponent discards a card") and heuristics for when an attack
// is likely to land past the opponent's blocks. The engine owns these — cards reach them
// through GameEngine accessors so the model's tuning lives in one place and stays out of
// card-side code.

package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// DiscardValue is the damage-equivalent credited when the opponent is forced to discard one
// card — one card they won't get to play. A typical FaB card is worth ~3 points of tempo.
const DiscardValue = 3

// GoldTokenValue is the placeholder credit for NotImplemented Gold-creating cards.
// Value 0 because Gold pays out only on activation (see GoldTokenAbility), not at
// creation.
const GoldTokenValue = 0

// LikelyToHit reports whether self's attack is likely to land past the opponent's blocks.
// Folds self.EffectiveAttack() (printed Card.Attack() + any granted BonusAttack, clamped at
// 0) and self.EffectiveDominate() (printed Dominator marker OR a granted Dominate flag) into
// the underlying threshold check. Cards reach this through GameEngine.LikelyToHit.
func LikelyToHit(self *card.CardState) bool {
	return LikelyDamageHits(self.EffectiveAttack(), self.EffectiveDominate())
}

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit. A typical FaB card
// is worth ~3 points, so blocking 1/4/7 with a pitch or block card over-pays; the opponent
// would rather eat the damage. Multiples of 3 are the easy-to-block amounts.
//
// dominate flips the math for cards printed (or granted) with the Dominate keyword: the
// defender is capped at one blocking card, so any attack of 5+ power slips at least 2 damage
// past that single block. The "if this hits" clause fires — we credit the rider.
//
// Most cards should use GameEngine.LikelyToHit instead. This raw form is for callers that
// probe a hypothetical damage value without a CardState (e.g. a fragile-aura helper asking
// "if N runechants fired at once, would they hit?").
func LikelyDamageHits(n int, dominate bool) bool {
	if dominate && n >= 5 {
		return true
	}
	return n == 1 || n == 4 || n == 7
}
