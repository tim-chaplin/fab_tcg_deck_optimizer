package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Equip-time weapon object construction. Weapons aren't played from hand — they're equipped
// at game start, so there is no card-facing CreateWeapon (the Aura / Item counterpart). The
// sim and the StateBuilder build the mutable weapon objects from the deck's platonic weapon
// cards through EquipFromCards.

// EquipFromCards builds the engine-side weapon object for each platonic weapon card via
// weapon.NewFromCard, registering the card's WeaponTrigger (the firing event + handler; (0,
// nil) for an untriggered weapon) and caching the card's Hands / Ability. The object lives in
// GameState.weapons across turns. Returns nil for an empty loadout so the no-weapon case
// carries no allocation.
func EquipFromCards(cards []weapon.Card) []Weapon {
	if len(cards) == 0 {
		return nil
	}
	out := make([]Weapon, len(cards))
	for i, c := range cards {
		tt, h := c.WeaponTrigger()
		out[i] = weapon.NewFromCard(c, tt, h, false, nil)
	}
	return out
}

// DestroyWeapon removes the weapon currently being fired from the arena and, when
// addToGraveyard is true, pushes its source weapon card into the graveyard. The weapon
// counterpart of DestroyAura / DestroyItem: direct splice with no cacheable flip —
// destruction is deterministic from the triggering event. Not used by any weapon yet
// (Talishar's end-phase self-destruct will be the first); wired so a future triggered
// weapon's handler-side Destroy routes back here.
func (ge *GameEngine) DestroyWeapon(addToGraveyard bool) {
	i := ge.currentHookIdx
	if i < 0 || i >= len(ge.weapons) {
		return
	}
	src := ge.weapons[i].SourceCard()
	ge.weapons = append(ge.weapons[:i], ge.weapons[i+1:]...)
	ge.currentHookDestroyed = true
	if src != nil && addToGraveyard {
		ge.AppendGraveyard(src.(card.Card))
	}
}
