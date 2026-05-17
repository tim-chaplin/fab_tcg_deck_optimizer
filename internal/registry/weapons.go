package registry

import "github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"

// Weapons are equipment, not deck cards; this roster is a small standalone lookup keyed by
// display name rather than the ID-indexed card slice.

// AllWeapons lists every implemented weapon. Used by deck-search code to enumerate loadouts.
var AllWeapons = []Weapon{
	weapons.AnnalsOfSutcliffe{},
	weapons.NebulaBlade{},
	weapons.ReapingBlade{},
	weapons.RosettaThorn{},
	weapons.ScepterOfPain{},
	weapons.Talishar{},
}

// weaponsByName maps Weapon.Name() → Weapon for reverse lookup.
var weaponsByName = func() map[string]Weapon {
	m := make(map[string]Weapon, len(AllWeapons))
	for _, w := range AllWeapons {
		m[w.Name()] = w
	}
	return m
}()

// WeaponByName returns the registered Weapon for the given display name. Returns
// (nil, false) when no such weapon exists; callers surface that to the user.
func WeaponByName(name string) (Weapon, bool) {
	w, ok := weaponsByName[name]
	return w, ok
}
