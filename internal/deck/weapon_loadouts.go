package deck

// Weapon-loadout helpers: canonical name sorting, display labels, a comparable key for equality
// checks, enumeration of every legal equip combination, and the 0–2 / 1H-pair validation applied
// at construction.

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// sortedWeaponNames returns the weapon names in ascending order. The canonical form both
// loadoutLabel and weaponKey build on so two loadouts with the same weapons in different orders
// compare equal.
func sortedWeaponNames(ws []weapons.Weapon) []string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	sort.Strings(names)
	return names
}

// loadoutLabel formats a weapon loadout for mutation descriptions, e.g. "[Nebula Blade]" or
// "[Reaping Blade, Scepter of Pain]".
func loadoutLabel(ws []weapons.Weapon) string {
	if len(ws) == 0 {
		return "[]"
	}
	return "[" + strings.Join(sortedWeaponNames(ws), ", ") + "]"
}

// weaponKey returns a comparable string for a weapon loadout so we can check equality.
func weaponKey(ws []weapons.Weapon) string {
	return strings.Join(sortedWeaponNames(ws), ",")
}

// weaponLoadouts enumerates every legal equip combination from ws: each 2H weapon as a solo
// loadout, plus every unordered pair of 1H weapons (including dual-wielding the same weapon).
func weaponLoadouts(ws []weapons.Weapon) [][]weapons.Weapon {
	var oneHand, twoHand []weapons.Weapon
	for _, w := range ws {
		if w.Hands() == 1 {
			oneHand = append(oneHand, w)
		} else {
			twoHand = append(twoHand, w)
		}
	}
	var out [][]weapons.Weapon
	for _, w := range twoHand {
		out = append(out, []weapons.Weapon{w})
	}
	for i := 0; i < len(oneHand); i++ {
		for j := i; j < len(oneHand); j++ {
			out = append(out, []weapons.Weapon{oneHand[i], oneHand[j]})
		}
	}
	return out
}

func validateWeapons(weapons []weapons.Weapon) {
	switch len(weapons) {
	case 0, 1:
		return
	case 2:
		if weapons[0].Hands() != 1 || weapons[1].Hands() != 1 {
			panic("deck: two-weapon loadout requires both weapons to be 1H")
		}
	default:
		panic(fmt.Sprintf("deck: invalid weapon count %d (max 2)", len(weapons)))
	}
}
