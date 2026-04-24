package cards

// Weapon roster + name lookup. Weapons aren't ID-indexed like cards (decks don't hold them;
// they're equipment), so the roster lives alongside the card index as a small standalone
// registry keyed by display name.

import "github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"

// AllWeapons lists every implemented weapon. Used by deck-search code to enumerate loadouts.
var AllWeapons = []weapon.Weapon{
	weapon.NebulaBlade{},
	weapon.ReapingBlade{},
	weapon.ScepterOfPain{},
}

// weaponsByName maps Weapon.Name() → Weapon for reverse lookup. Built once at init.
var weaponsByName = func() map[string]weapon.Weapon {
	m := make(map[string]weapon.Weapon, len(AllWeapons))
	for _, w := range AllWeapons {
		m[w.Name()] = w
	}
	return m
}()

// WeaponByName returns the registered Weapon for the given display name. Returns (nil, false)
// when no such weapon exists — serialization callers surface that to the user.
func WeaponByName(name string) (weapon.Weapon, bool) {
	w, ok := weaponsByName[name]
	return w, ok
}
