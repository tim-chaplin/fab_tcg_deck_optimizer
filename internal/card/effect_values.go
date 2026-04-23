// Core tuning assumptions of the model. Damage-equivalent values for non-damage riders ("draw
// a card", "create a Gold token", "opponent discards a card") and heuristics for when an attack
// is likely to land past the opponent's blocks. Any card that needs one of these should call
// through here rather than hardcoding — if we ever re-tune the model, the revision lands in
// one place.
//
// Core stats (damage, health, block) intentionally aren't factored out: they're 1-to-1 with
// damage by convention and threading them through named constants would just be noise.

package card

// DiscardValue is the damage-equivalent credited when the opponent is forced to discard one
// card — one card they won't get to play. A typical FaB card is worth ~3 points of tempo.
const DiscardValue = 3

// GoldTokenValue is the damage-equivalent credited when a card creates a Gold token. A Gold
// token is is one future card draw, at the expense of two resource.
const GoldTokenValue = 0

// LikelyToHit reports whether dealing n damage is likely to get through an opponent's blocks.
// A typical FaB card is worth ~3 points, so blocking 1/4/7 with a pitch or block card
// over-pays; the opponent would rather eat the damage. Multiples of 3 are the easy-to-block
// amounts. Used by riders like "if this hits, …" to decide whether the clause fires.
//
// dominate flips the math for cards printed (or granted) with the Dominate keyword: the
// defender is capped at one blocking card, so any attack of 5+ power slips at least 2 damage
// past that single block. The "if this hits" clause fires — we credit the rider.
func LikelyToHit(n int, dominate bool) bool {
	if dominate && n >= 5 {
		return true
	}
	return n == 1 || n == 4 || n == 7
}
