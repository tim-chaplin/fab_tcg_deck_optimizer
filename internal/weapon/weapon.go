// Package weapon defines the Weapon interface for Flesh and Blood weapons (equipment, not
// deck cards). A deck equips 0–2 weapons under the "2 × 1H or 1 × 2H" rule.
package weapon

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Weapon is an equipped weapon. Weapons are equipment rather than deck cards. The methods
// mirror the subset of card.Card the simulator pipeline calls on a swung weapon (Play /
// Cost / Types / Attack / GoAgain / Name / Pitch / Defense); Hands() returns 1 or 2 for the
// equipment-slot rule.
//
// Weapons satisfy card.Card structurally — ids.WeaponID is currently aliased to ids.CardID
// (see registry/ids/weapon_ids.go) so the chain runner can hold weapons alongside cards in
// a single attacker slice without an adapter.
type Weapon interface {
	ID() ids.WeaponID
	Name() string
	Cost(*card.TurnState) int
	Pitch() int
	Attack() int
	Defense() int
	Types() card.TypeSet
	GoAgain() bool
	Play(s *card.TurnState, self *card.CardState)
	Hands() int
}
