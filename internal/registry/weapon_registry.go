package registry

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Weapon roster + name lookup. Weapons are equipment, not deck cards, so they aren't
// ID-indexed like cards; this registry is a small standalone lookup keyed by display name.

// AllWeapons lists every implemented weapon. Used by deck-search code to enumerate loadouts.
var AllWeapons = []sim.Weapon{
	weapons.AnnalsOfSutcliffe{},
	weapons.NebulaBlade{},
	weapons.ReapingBlade{},
	weapons.RosettaThorn{},
	weapons.ScepterOfPain{},
	weapons.Talishar{},
}

// weaponsByName maps Weapon.Name() → Weapon for reverse lookup. Built once at init.
var weaponsByName = func() map[string]sim.Weapon {
	m := make(map[string]sim.Weapon, len(AllWeapons))
	for _, w := range AllWeapons {
		m[w.Name()] = w
	}
	return m
}()

// WeaponByName returns the registered Weapon for the given display name. Returns
// (nil, false) when no such weapon exists — serialization callers surface that to the user.
func WeaponByName(name string) (sim.Weapon, bool) {
	w, ok := weaponsByName[name]
	return w, ok
}
