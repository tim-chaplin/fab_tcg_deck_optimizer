// Package weapon defines the Weapon interface for Flesh and Blood weapons (equipment, not deck
// cards). A deck equips 0–2 weapons under the "2 × 1H or 1 × 2H" rule.
package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Weapon is an equipped weapon. Implements card.Card so swings slot into the same attacker /
// ordering pipeline as hand cards. Pitch() and Defense() return 0 (weapons don't fill those
// roles), Attack() is the printed Power, and Hands() returns 1 or 2 for equipment-slot rules.
type Weapon interface {
	card.Card
	Hands() int
}
