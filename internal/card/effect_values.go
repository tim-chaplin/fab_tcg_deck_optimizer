// Damage-equivalent values for non-damage card effects. Any card that models a rider like
// "draw a card", "create a Gold token", "opponent discards a card" should credit one of these
// constants rather than hardcoding a number — so if we ever re-evaluate what a card's worth of
// future tempo actually costs the opponent, the revision lands in one place.
//
// Core stats (damage, health, block) intentionally aren't factored out: they're 1-to-1 with
// damage by convention and threading them through named constants would just be noise.

package card

// DrawValue is the damage-equivalent credited for drawing one extra card. A typical FaB card
// is worth ~3 points of tempo when played on a future turn.
const DrawValue = 3

// DiscardValue is the damage-equivalent credited when the opponent is forced to discard one
// card. Symmetric to DrawValue — a discarded card is one the opponent won't get to play.
const DiscardValue = 3

// GoldTokenValue is the damage-equivalent credited when a card creates a Gold token. A Gold
// token is one free resource (one {r}) worth of future tempo — about 1 damage's worth once
// spent.
const GoldTokenValue = 1
