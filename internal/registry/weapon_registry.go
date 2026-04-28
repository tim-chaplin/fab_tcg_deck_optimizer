package registry

import "github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"

// Weapon roster + name lookup. Weapons are equipment, not deck cards, so they aren't
// ID-indexed like cards; this registry is a small standalone lookup keyed by display name.

// AllWeapons lists every implemented weapon. Used by deck-search code to enumerate loadouts.
var AllWeapons = []weapon.Weapon{
	weapon.AnnalsOfSutcliffe{},
	weapon.NebulaBlade{},
	weapon.ReapingBlade{},
	weapon.RosettaThorn{},
	weapon.ScepterOfPain{},
	weapon.Talishar{},
}

// weaponsByName maps Weapon.Name() → Weapon for reverse lookup. Built once at init.
var weaponsByName = func() map[string]weapon.Weapon {
	m := make(map[string]weapon.Weapon, len(AllWeapons))
	for _, w := range AllWeapons {
		m[w.Name()] = w
	}
	return m
}()

// WeaponByName returns the registered Weapon for the given display name. Returns
// (nil, false) when no such weapon exists — serialization callers surface that to the user.
func WeaponByName(name string) (weapon.Weapon, bool) {
	w, ok := weaponsByName[name]
	return w, ok
}
