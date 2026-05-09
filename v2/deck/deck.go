// Package deck represents a candidate FaB deck — hero, weapons, cards, plus the user's
// sideboard / equipment lists. Search code creates many Decks via Random and edits them via
// ApplyDefaults / SanitizeNotImplemented. Simulation lives elsewhere; deck depends only on
// the narrow Hero / Weapon / Card / Registry contracts declared in this package, so a
// caller can swap in any concrete card / weapon / hero / registry implementation.
package deck

import (
	"fmt"
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Deck is a hero, equipped weapons, and a deck of cards. Sideboard is the reserve-card
// list the user manages for sideboarding between games; Equipment is the non-weapon arena
// loadout (head, chest, arms, legs). Both round-trip through the persistence layer; the
// simulator never reads either, so mutations leave them alone.
//
// Sideboard and Equipment are []string rather than []Card: equipment pieces and other items
// the user wants on their sideboard list (e.g. Nullrune cycle) aren't in the card registry,
// so a registry-backed field would force the user's data through a lossy lookup.
type Deck struct {
	Hero      Hero
	Weapons   []Weapon
	Cards     []Card
	Sideboard []string
	Equipment []string
}

// New constructs a Deck. Panics if the weapon loadout violates the "0–2 weapons; if 2, both
// 1H" equipment rule. Sideboard and Equipment start empty; callers assign them directly when
// carrying them over.
func New(h Hero, weapons []Weapon, cards []Card) *Deck {
	validateWeapons(weapons)
	return &Deck{Hero: h, Weapons: weapons, Cards: cards}
}

// SideboardDefault is one "always include in the sideboard" entry the caller passes to
// ApplyDefaults: a card display name plus the target copy count to top the sideboard up
// toward. Count must be in [1, SideboardCopyCap]; a larger target silently clamps when the
// merge respects the main-deck + sideboard copy cap.
type SideboardDefault struct {
	Name  string
	Count int
}

// Defaults is the loadout ApplyDefaults tops a deck up toward. Equipment names are added
// at one copy each; Sideboard entries are added up to their per-entry Count, capped by
// SideboardCopyCap against main-deck + sideboard copies.
type Defaults struct {
	Equipment []string
	Sideboard []SideboardDefault
}

// SideboardCopyCap is the per-card copy limit across main deck + sideboard combined. The
// default-sideboard merger respects it so a default addition never pushes a card past the
// normal deck-construction max.
const SideboardCopyCap = 2

// ApplyDefaults tops d.Equipment and d.Sideboard up toward the supplied defaults so
// persisted decks always carry the common "every deck runs these" slots the caller picks
// for this hero / format. Idempotent: running it twice is a no-op because each entry is
// only added when the current count falls below its target. Equipment targets 1 copy per
// entry; sideboard targets each entry's Count, clamped by SideboardCopyCap against
// main-deck + sideboard copies so the merge never pushes a card past the deck-construction
// limit.
func (d *Deck) ApplyDefaults(defaults Defaults) {
	equipCounts := map[string]int{}
	for _, name := range d.Equipment {
		equipCounts[name]++
	}
	for _, name := range defaults.Equipment {
		if equipCounts[name] < 1 {
			d.Equipment = append(d.Equipment, name)
			equipCounts[name]++
		}
	}

	mainCounts := map[string]int{}
	for _, c := range d.Cards {
		mainCounts[c.DisplayName()]++
	}
	sideCounts := map[string]int{}
	for _, name := range d.Sideboard {
		sideCounts[name]++
	}
	for _, entry := range defaults.Sideboard {
		room := SideboardCopyCap - mainCounts[entry.Name] - sideCounts[entry.Name]
		if room <= 0 {
			continue
		}
		want := entry.Count - sideCounts[entry.Name]
		if want <= 0 {
			continue
		}
		if want > room {
			want = room
		}
		for i := 0; i < want; i++ {
			d.Sideboard = append(d.Sideboard, entry.Name)
			sideCounts[entry.Name]++
		}
	}
}

// Random generates a random legal deck for h: a random weapon loadout from
// reg.LegalWeapons (one 2H or two 1H; dual-wielding the same weapon allowed) and size
// cards drawn uniformly from reg.LegalCards one at a time, skipping any roll that would
// exceed maxCopies for the picked card's ID.
//
// legal further filters the card pool: only cards for which legal(c) returns true are
// candidates. Pass nil for no filtering. Callers typically wire a constructed format's
// banlist predicate through here.
func Random(h Hero, size, maxCopies int, rng *rand.Rand, legal func(Card) bool, reg Registry) *Deck {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 1 (got %d)", maxCopies))
	}
	loadouts := weaponLoadouts(reg.LegalWeapons())
	if len(loadouts) == 0 {
		panic("deck: Random has no legal weapon loadouts — registry rejected every weapon")
	}
	weapons := loadouts[rng.Intn(len(loadouts))]

	pool := legalPool(reg, legal)
	if len(pool) == 0 {
		panic("deck: Random's legal filter rejected every card — cannot build a deck")
	}
	counts := map[ids.CardID]int{}
	picks := make([]Card, 0, size)
	for len(picks) < size {
		c := pool[rng.Intn(len(pool))]
		if counts[c.ID()]+1 > maxCopies {
			continue
		}
		counts[c.ID()]++
		picks = append(picks, c)
	}
	return New(h, weapons, picks)
}

// NotImplementedReplacement records one swap made by Deck.SanitizeNotImplemented: the
// pool-excluded card that was removed and the pool-eligible card that took its slot.
type NotImplementedReplacement struct {
	From Card
	To   Card
}

// SanitizeNotImplemented replaces every card no longer present in reg.LegalCards with a
// random legal replacement, respecting maxCopies. Returns the ordered list of swaps (empty
// when nothing needed replacement). Panics when maxCopies < 1, when the legal pool is
// empty, or when every pool entry is already at maxCopies in the surviving slots (no
// legal replacement can be picked).
func (d *Deck) SanitizeNotImplemented(maxCopies int, rng *rand.Rand, legal func(Card) bool, reg Registry) []NotImplementedReplacement {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: SanitizeNotImplemented requires maxCopies >= 1 (got %d)", maxCopies))
	}
	pool := legalPool(reg, legal)
	if len(pool) == 0 {
		panic("deck: SanitizeNotImplemented's legal filter rejected every pool-eligible card — cannot build a replacement")
	}
	// Build the legal-id set from the pool so the per-slot check is a map lookup.
	legalSet := make(map[ids.CardID]bool, len(pool))
	for _, c := range pool {
		legalSet[c.ID()] = true
	}
	// Seed counts with the keeper cards already in the deck so replacements respect
	// maxCopies against the surviving slots. The excluded slots we're about to overwrite
	// don't count.
	counts := map[ids.CardID]int{}
	var slots []int
	for i, c := range d.Cards {
		if !legalSet[c.ID()] {
			slots = append(slots, i)
			continue
		}
		counts[c.ID()]++
	}
	if len(slots) == 0 {
		return nil
	}
	replacements := make([]NotImplementedReplacement, 0, len(slots))
	for _, idx := range slots {
		// Reject-resample the pool until a non-saturated card lands. The bounded loop
		// guards against the pathological "every pool entry is already at maxCopies"
		// case the pre-rewrite version would loop forever on; once attempts exceeds the
		// pool-size × maxCopies upper bound on legal picks, panic with a clear error.
		var pick Card
		maxAttempts := len(pool) * (maxCopies + 1)
		picked := false
		for attempt := 0; attempt < maxAttempts; attempt++ {
			pick = pool[rng.Intn(len(pool))]
			if counts[pick.ID()]+1 <= maxCopies {
				picked = true
				break
			}
		}
		if !picked {
			panic("deck: SanitizeNotImplemented's pool can't satisfy the per-printing budget the surviving slots already use up")
		}
		counts[pick.ID()]++
		from := d.Cards[idx]
		d.Cards[idx] = pick
		replacements = append(replacements, NotImplementedReplacement{From: from, To: pick})
	}
	return replacements
}

// legalPool returns reg.LegalCards filtered by legal. legal == nil disables the extra
// filter (still excludes pool-marked cards via reg's own definition). Shared by Random
// and SanitizeNotImplemented.
func legalPool(reg Registry, legal func(Card) bool) []Card {
	pool := reg.LegalCards()
	if legal == nil {
		return pool
	}
	filtered := pool[:0]
	for _, c := range pool {
		if legal(c) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
