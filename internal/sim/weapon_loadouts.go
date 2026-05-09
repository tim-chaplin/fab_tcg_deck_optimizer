package sim

// Weapon-loadout helpers used by sim's mutation enumeration. The enum-every-combo logic
// also lives in v2/deck (private), but the mutation generator works at the deck.Weapon
// level too — duplicating the small helpers here keeps mutations.go from depending on
// v2/deck internals.

import (
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// sortedWeaponNames returns the weapon names in ascending order. The canonical form both
// loadoutLabel and weaponKey build on so two loadouts with the same weapons in different
// orders compare equal.
func sortedWeaponNames(ws []deck.Weapon) []string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	sort.Strings(names)
	return names
}

// loadoutLabel formats a weapon loadout for mutation descriptions, e.g. "[Nebula Blade]"
// or "[Reaping Blade, Scepter of Pain]".
func loadoutLabel(ws []deck.Weapon) string {
	if len(ws) == 0 {
		return "[]"
	}
	return "[" + strings.Join(sortedWeaponNames(ws), ", ") + "]"
}

// weaponKey returns a comparable string for a weapon loadout so we can check equality.
func weaponKey(ws []deck.Weapon) string {
	return strings.Join(sortedWeaponNames(ws), ",")
}

// weaponLoadouts enumerates every legal equip combination from ws: each 2H weapon as a solo
// loadout, plus every unordered pair of 1H weapons (including dual-wielding the same weapon).
func weaponLoadouts(ws []deck.Weapon) [][]deck.Weapon {
	var oneHand, twoHand []deck.Weapon
	for _, w := range ws {
		if w.Hands() == 1 {
			oneHand = append(oneHand, w)
		} else {
			twoHand = append(twoHand, w)
		}
	}
	var out [][]deck.Weapon
	for _, w := range twoHand {
		out = append(out, []deck.Weapon{w})
	}
	for i := 0; i < len(oneHand); i++ {
		for j := i; j < len(oneHand); j++ {
			out = append(out, []deck.Weapon{oneHand[i], oneHand[j]})
		}
	}
	return out
}
