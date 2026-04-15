// Package weapon defines the Weapon interface for Flesh and Blood weapons (equipment, not deck
// cards). A deck equips 0–2 weapons subject to the "2 × 1H or 1 × 2H" rule.
package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Weapon is an equipped weapon. It implements card.Card so that swinging it can slot into the same
// attacker/ordering pipeline used for cards from hand; Pitch() and Defense() should return 0 since
// weapons don't fill those roles, and Attack() is the weapon's printed Power. Hands() returns 1 or
// 2 to enforce equipment-slot rules.
type Weapon interface {
	card.Card
	Hands() int
}
