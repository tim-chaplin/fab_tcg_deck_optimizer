// Package deck represents a candidate FaB deck and the hand-value stats accumulated from
// simulating it. Search code creates many Decks, evaluates each, and compares their Stats.
//
// The Deck type and its construction live in this file. Cohesive concern groups are split
// across sibling files in this package: weapon_loadouts.go (loadout helpers + validation),
// stats.go (Stats / BestTurn / CardPlayStats / CycleStats), mutations.go (iterate-mode
// candidate generation), evaluate.go (hand-by-hand simulation), iterate.go (parallel
// simulated-annealing round runner).
package deck

import (
	"fmt"
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Deck is a hero, equipped weapons, a deck of cards, and the simulated hand-value stats.
// Sideboard is the reserve-card list the user manages for sideboarding between games;
// Equipment is the non-weapon arena loadout (head, chest, arms, legs). Both round-trip
// through deckio and fabrary; the simulator never reads either, so mutations and Evaluate
// leave them alone.
//
// Both are []string rather than []card.Card: equipment pieces and other items the user
// wants on their sideboard list (e.g. Nullrune cycle) aren't in the card registry, so a
// registry-backed field would force the user's data through a lossy lookup.
type Deck struct {
	Hero      hero.Hero
	Weapons   []weapon.Weapon
	Cards     []card.Card
	Sideboard []string
	Equipment []string
	Stats     Stats
}

// New constructs a Deck. Panics if the weapon loadout violates the "0–2 weapons; if 2, both 1H"
// equipment rule. Sideboard and Equipment start empty; callers assign them directly when
// carrying them over.
func New(h hero.Hero, weapons []weapon.Weapon, cards []card.Card) *Deck {
	validateWeapons(weapons)
	return &Deck{Hero: h, Weapons: weapons, Cards: cards}
}

// defaultEquipment lists equipment names that every persisted deck should carry in its
// Equipment section. ApplyDefaults tops Equipment up to include each of these at least once;
// the user can add more copies but never drops below one.
var defaultEquipment = []string{
	"Beckoning Haunt",
	"Blade Beckoner Boots",
	"Blade Beckoner Helm",
	"Blossom of Spring",
}

// defaultSideboardEntry is one "always include in the sideboard" default: a name plus the
// target copy count ApplyDefaults tops the sideboard up toward. Invariant: count must be in
// [1, sideboardCopyCap] — a larger target would silently clamp when the merge respects the
// main-deck + sideboard copy cap.
type defaultSideboardEntry struct {
	name  string
	count int
}

// defaultSideboard items are appended to Sideboard by ApplyDefaults. For each entry, the
// merger tops the sideboard count up toward `count`, but never past sideboardCopyCap (2 per
// card across main + sideboard). Equipment-slot items (Crown of Dichotomy, Nullrune
// boots/gloves, Runebleed Robe) target 1 copy; deck cards target 2.
//
// Card names must match card.DisplayName format ("Read the Runes [R]") since ApplyDefaults
// dedupes via DisplayName-keyed counts. Equipment slots stay bare since they're not pitch-
// varying cards.
var defaultSideboard = []defaultSideboardEntry{
	{"Crown of Dichotomy", 1},
	{"Nullrune Boots", 1},
	{"Nullrune Gloves", 1},
	{"Runebleed Robe", 1},
	{"Read the Runes [R]", 2},
	{"Reduce to Runechant [R]", 2},
	{"Sigil of Suffering [R]", 2},
}

// sideboardCopyCap is the per-card copy limit across main deck + sideboard combined. The
// default-sideboard merger respects this so a default addition never pushes a card past the
// normal deck-construction max.
const sideboardCopyCap = 2

// ApplyDefaults tops d.Equipment and d.Sideboard up toward the hardcoded default loadout so
// persisted decks always carry the common "every Viserai deck runs these" slots. Idempotent:
// running it twice is a no-op because each entry is only added when the current count falls
// below its target. Equipment targets 1 copy per entry; sideboard targets each entry's
// count, but is clamped by sideboardCopyCap against main-deck + sideboard copies so the
// merge never pushes a card past the deck-construction limit.
func (d *Deck) ApplyDefaults() {
	equipCounts := map[string]int{}
	for _, name := range d.Equipment {
		equipCounts[name]++
	}
	for _, name := range defaultEquipment {
		if equipCounts[name] < 1 {
			d.Equipment = append(d.Equipment, name)
			equipCounts[name]++
		}
	}

	mainCounts := map[string]int{}
	for _, c := range d.Cards {
		mainCounts[card.DisplayName(c)]++
	}
	sideCounts := map[string]int{}
	for _, name := range d.Sideboard {
		sideCounts[name]++
	}
	for _, entry := range defaultSideboard {
		room := sideboardCopyCap - mainCounts[entry.name] - sideCounts[entry.name]
		if room <= 0 {
			continue
		}
		want := entry.count - sideCounts[entry.name]
		if want <= 0 {
			continue
		}
		if want > room {
			want = room
		}
		for i := 0; i < want; i++ {
			d.Sideboard = append(d.Sideboard, entry.name)
			sideCounts[entry.name]++
		}
	}
}

// Random generates a random legal deck for h: a random weapon loadout from registry.AllWeapons (one 2H
// or two 1H; dual-wielding the same weapon allowed) and size cards drawn uniformly from
// registry.DeckableCards() one at a time, skipping any roll that would exceed maxCopies for the picked
// ID. Matches the single-slot granularity of deck.AllMutations so the hill-climb can explore
// the space the generator actually produces.
//
// legal filters the card pool: only IDs for which legal(registry.GetCard(id)) returns true are
// candidates. Pass nil for no filtering. Callers typically wire deckformat.Format.IsLegal
// through here to restrict generation to a constructed format's banlist.
func Random(h hero.Hero, size, maxCopies int, rng *rand.Rand, legal func(card.Card) bool) *Deck {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 1 (got %d)", maxCopies))
	}
	loadouts := weaponLoadouts(registry.AllWeapons)
	weapons := loadouts[rng.Intn(len(loadouts))]

	pool := legalPool(legal)
	if len(pool) == 0 {
		panic("deck: Random's legal filter rejected every card — cannot build a deck")
	}
	counts := map[ids.CardID]int{}
	picks := make([]card.Card, 0, size)
	for len(picks) < size {
		id := pool[rng.Intn(len(pool))]
		if counts[id]+1 > maxCopies {
			continue
		}
		counts[id]++
		picks = append(picks, registry.GetCard(id))
	}
	return New(h, weapons, picks)
}

// NotImplementedReplacement records one swap made by Deck.SanitizeNotImplemented: the
// card.NotImplemented-tagged card that was removed and the card.NotImplemented-free card
// that took its slot.
type NotImplementedReplacement struct {
	From card.Card
	To   card.Card
}

// SanitizeNotImplemented scans d.Cards in order and replaces every card carrying
// card.NotImplemented with a random legal replacement drawn from legalPool(legal), rolling
// again if the pick would exceed maxCopies for that printing. Weapons and hero are
// untouched. The sanitized deck stays size-stable and copy-cap-legal so the caller can
// re-evaluate directly.
//
// Returns the ordered list of swaps made (one entry per replaced slot; duplicates of the
// same tagged ID produce one entry each since each slot is picked independently). Returns
// an empty slice when nothing needed replacement.
//
// Panics when maxCopies < 1 or when the legal/NotImplemented pool is smaller than the
// per-printing maxCopies budget d already uses — both indicate a config so degenerate that
// there's no sensible recovery.
func (d *Deck) SanitizeNotImplemented(maxCopies int, rng *rand.Rand, legal func(card.Card) bool) []NotImplementedReplacement {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: SanitizeNotImplemented requires maxCopies >= 1 (got %d)", maxCopies))
	}
	pool := legalPool(legal)
	if len(pool) == 0 {
		panic("deck: SanitizeNotImplemented's legal filter rejected every implemented card — cannot build a replacement")
	}
	// Seed counts with the implemented-keeper cards already in the deck so replacements
	// respect maxCopies against the surviving slots. The tagged slots we're about to
	// overwrite don't count.
	counts := map[ids.CardID]int{}
	var slots []int
	for i, c := range d.Cards {
		if _, unimplemented := c.(card.NotImplemented); unimplemented {
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
		var pick ids.CardID
		for {
			pick = pool[rng.Intn(len(pool))]
			if counts[pick]+1 <= maxCopies {
				break
			}
		}
		counts[pick]++
		from := d.Cards[idx]
		to := registry.GetCard(pick)
		d.Cards[idx] = to
		replacements = append(replacements, NotImplementedReplacement{From: from, To: to})
	}
	return replacements
}

// legalPool returns registry.DeckableCards() filtered by legal, with any card carrying the
// card.NotImplemented marker removed. The NotImplemented filter is always applied — a card
// whose printed effect the sim can't faithfully reproduce shouldn't land in a random deck or
// become a mutation candidate regardless of format legality. Pass nil for legal to apply only
// the NotImplemented filter. Shared by Random and AllMutations so both agree on the pool.
func legalPool(legal func(card.Card) bool) []ids.CardID {
	pool := registry.DeckableCards()
	filtered := pool[:0]
	for _, id := range pool {
		c := registry.GetCard(id)
		if _, unimplemented := c.(card.NotImplemented); unimplemented {
			continue
		}
		if legal != nil && !legal(c) {
			continue
		}
		filtered = append(filtered, id)
	}
	return filtered
}
