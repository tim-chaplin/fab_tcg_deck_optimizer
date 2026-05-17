package deck

import (
	"fmt"
	"sort"
	"strings"
)

// validateWeapons enforces FaB's "0–2 weapons; if 2, both 1H" equipment rule. Called from
// New so a panic surfaces at construction rather than mid-simulation.
func validateWeapons(weapons []Weapon) {
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

// weaponLoadouts enumerates every legal equip combination from ws: each 2H weapon as a
// solo loadout, plus every unordered pair of 1H weapons (including dual-wielding the same
// weapon). Used by Random to pick a starting loadout uniformly across all legal shapes,
// and by the mutation enumerator to surface every alternative loadout for an existing
// deck.
func weaponLoadouts(ws []Weapon) [][]Weapon {
	var oneHand, twoHand []Weapon
	for _, w := range ws {
		if w.Hands() == 1 {
			oneHand = append(oneHand, w)
		} else {
			twoHand = append(twoHand, w)
		}
	}
	var out [][]Weapon
	for _, w := range twoHand {
		out = append(out, []Weapon{w})
	}
	for i := 0; i < len(oneHand); i++ {
		for j := i; j < len(oneHand); j++ {
			out = append(out, []Weapon{oneHand[i], oneHand[j]})
		}
	}
	return out
}

// sortedWeaponNames returns the weapon names in ascending order. The canonical form both
// loadoutLabel and weaponKey build on so two loadouts with the same weapons in different
// orders compare equal.
func sortedWeaponNames(ws []Weapon) []string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	sort.Strings(names)
	return names
}

// loadoutLabel formats a weapon loadout for mutation descriptions, e.g. "[Nebula Blade]"
// or "[Reaping Blade, Scepter of Pain]".
func loadoutLabel(ws []Weapon) string {
	if len(ws) == 0 {
		return "[]"
	}
	return "[" + strings.Join(sortedWeaponNames(ws), ", ") + "]"
}

// weaponKey returns a comparable string for a weapon loadout so callers can check
// equality (e.g. tests asserting that two mutations produced identical loadouts).
func weaponKey(ws []Weapon) string {
	return strings.Join(sortedWeaponNames(ws), ",")
}
