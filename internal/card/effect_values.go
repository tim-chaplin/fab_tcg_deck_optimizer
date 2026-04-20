// Core tuning assumptions of the model. Damage-equivalent values for non-damage riders ("draw
// a card", "create a Gold token", "opponent discards a card") and heuristics for when an attack
// is likely to land past the opponent's blocks. Any card that needs one of these should call
// through here rather than hardcoding — if we ever re-tune the model, the revision lands in
// one place.
//
// Core stats (damage, health, block) intentionally aren't factored out: they're 1-to-1 with
// damage by convention and threading them through named constants would just be noise.

package card

// DrawValue is the damage-equivalent credited for a cross-turn additive card draw — one the
// sim doesn't currently route through TurnState.DrawOne because the card enters hand outside
// the normal refill step (e.g. Sigil of the Arknight's start-of-action-phase reveal). A typical
// FaB card is worth ~3 points of tempo when played on a future turn.
const DrawValue = 3

// DiscardValue is the damage-equivalent credited when the opponent is forced to discard one
// card. Symmetric to DrawValue — a discarded card is one the opponent won't get to play.
const DiscardValue = 3

// GoldTokenValue is the damage-equivalent credited when a card creates a Gold token. A Gold
// token is one free resource (one {r}) worth of future tempo — about 1 damage's worth once
// spent.
const GoldTokenValue = 1

// LikelyToHit reports whether dealing n damage is likely to get through an opponent's blocks.
// A typical FaB card is worth ~3 points, so blocking 1/4/7 with a pitch or block card over-pays;
// the opponent would rather eat the damage. Multiples of 3 are the easy-to-block amounts.
// Used by fragile-aura cards (Arcane Cussing, Bloodspill Invocation) to decide whether a
// same-turn attack will actually land and pop the aura for its pay-off.
func LikelyToHit(n int) bool {
	return n == 1 || n == 4 || n == 7
}
