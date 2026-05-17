// Package weapon defines the Weapon interface for Flesh and Blood weapons (equipment, not
// deck cards). A deck equips 0–2 weapons under the "2 × 1H or 1 × 2H" rule.
package weapon

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Weapon is an equipped permanent. Each turn the chain runner enqueues the weapon's
// Ability() Card as a playable option (1 AP, pays the ability's printed activation cost);
// the weapon itself never enters the chain. Hands() is 1 or 2 for the "2 × 1H or 1 × 2H"
// equipment-slot rule. See docs/dev-standards.md → "Weapon activated abilities".
type Weapon interface {
	ID() ids.WeaponID
	Name() string
	Types() card.TypeSet
	Hands() int
	Ability() card.Card
}
