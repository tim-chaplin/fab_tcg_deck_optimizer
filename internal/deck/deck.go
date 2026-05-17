// Package deck represents a candidate FaB deck — hero, weapons, cards, plus the user's
// sideboard / equipment lists. Search code creates many Decks via Random and edits them via
// ApplyDefaults. Simulation lives elsewhere; deck depends only on the narrow Hero / Weapon
// / Card / Registry contracts declared in this package, so a caller can swap in any concrete
// card / weapon / hero / registry implementation.
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Deck is a hero, equipped weapons, and a deck of cards. Sideboard is the reserve-card
// list the user manages for sideboarding between games; Equipment is the non-weapon arena
// loadout (head, chest, arms, legs). Both round-trip through the persistence layer; the
// simulator never reads either, so mutations leave them alone.
//
// Sideboard and Equipment are []string rather than []Card: equipment pieces and other items
// the user wants on their sideboard list (e.g. Nullrune cycle) aren't in the card registry,
// so a registry-backed field would force the user's data through a lossy lookup.
//
// cards doubles as the runtime deck: Shuffle / Draw / PeekTop / PutBottom / PutTop / Tutor
// mutate it directly. Callers running an evaluation trial should Copy() the master deck first
// so their deck mutations don't disturb the master — see deck.Copy. The field is unexported
// so external callers can't peek the runtime order; composition-level inspection goes through
// UniqueIDs / NameCounts / DisplayNames / PitchCounts.
type Deck struct {
	Hero      Hero
	Weapons   []Weapon
	cards     []Card
	Sideboard []string
	Equipment []string
}

// UniqueIDs returns the distinct card IDs in deck order of first appearance plus a
// position-lookup map keyed by ID. The simulator pre-builds this once per Evaluate run to
// drive its per-shuffle hand-presence accounting (the parallel marginal-stats buffers are
// indexed by position so shuffles can tally without growing maps in the inner loop).
func (d *Deck) UniqueIDs() ([]ids.CardID, map[ids.CardID]int) {
	if d == nil {
		return nil, nil
	}
	out := make([]ids.CardID, 0, len(d.cards))
	idx := make(map[ids.CardID]int, len(d.cards))
	for _, c := range d.cards {
		id := c.ID()
		if _, ok := idx[id]; ok {
			continue
		}
		idx[id] = len(out)
		out = append(out, id)
	}
	return out, idx
}

// New constructs a Deck. Panics if the weapon loadout violates the "0–2 weapons; if 2, both
// 1H" equipment rule. Sideboard and Equipment start empty; callers assign them directly when
// carrying them over.
func New(h Hero, weapons []Weapon, cards []Card) *Deck {
	validateWeapons(weapons)
	return &Deck{Hero: h, Weapons: weapons, cards: cards}
}

// Size reports the number of cards in the deck. Excludes Sideboard and Equipment, which are
// reserve / arena lists the simulator never reads. Nil-safe: a nil Deck reports 0 so
// callers holding a zero-value TurnState (no deck attached) read empty without a guard.
func (d *Deck) Size() int {
	if d == nil {
		return 0
	}
	return len(d.cards)
}

// Fingerprint returns a comparable summary of the deck for equality checks: the weapon
// loadout (names sorted) and a sorted card-count histogram. Two decks with the same cards
// in different orders, or with weapons listed in different orders, produce equal
// fingerprints — so callers can compare a mutation result against a baseline without
// caring about positional shuffles.
func (d *Deck) Fingerprint() string {
	var b strings.Builder
	weaponNames := sortedWeaponNames(d.Weapons)
	for i, name := range weaponNames {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(name)
	}
	b.WriteByte('|')
	counts := map[ids.CardID]int{}
	for _, c := range d.cards {
		counts[c.ID()]++
	}
	cardIDs := make([]int, 0, len(counts))
	for id := range counts {
		cardIDs = append(cardIDs, int(id))
	}
	sort.Ints(cardIDs)
	for _, id := range cardIDs {
		fmt.Fprintf(&b, "%d:%d,", id, counts[ids.CardID(id)])
	}
	return b.String()
}

// Copy returns a fresh Deck with independent backing slices for Weapons, Cards, Sideboard,
// and Equipment. Hero is shared (heroes are stateless concrete types). Used by the parallel
// evaluator: each worker calls master.Copy() before its trial so Shuffle / Draw / Reset /
// PutBottom mutations land on the worker's local copy and the master deck stays sharable.
func (d *Deck) Copy() *Deck {
	if d == nil {
		return &Deck{}
	}
	out := &Deck{Hero: d.Hero}
	if len(d.Weapons) > 0 {
		out.Weapons = append(make([]Weapon, 0, len(d.Weapons)), d.Weapons...)
	}
	if len(d.cards) > 0 {
		out.cards = append(make([]Card, 0, len(d.cards)), d.cards...)
	}
	if len(d.Sideboard) > 0 {
		out.Sideboard = append(make([]string, 0, len(d.Sideboard)), d.Sideboard...)
	}
	if len(d.Equipment) > 0 {
		out.Equipment = append(make([]string, 0, len(d.Equipment)), d.Equipment...)
	}
	return out
}

// CopyFrom provides a memory-efficient way to copy a Deck by reusing an already-allocated
// Deck (the receiver) that's no longer needed: d's cards / Weapons / Hero are overwritten
// from src, reusing d's existing slice backing arrays when they have enough capacity.
// Sideboard and Equipment are skipped. Callers that need a full deep copy should use Copy.
//
// d must be a non-nil receiver. A nil src is treated as an empty deck — d's slice lengths
// drop to 0 (backings retained) and Hero zeroes.
func (d *Deck) CopyFrom(src *Deck) {
	if src == nil {
		d.Hero = nil
		d.cards = d.cards[:0]
		d.Weapons = d.Weapons[:0]
		return
	}
	d.Hero = src.Hero
	if cap(d.cards) >= len(src.cards) {
		d.cards = append(d.cards[:0], src.cards...)
	} else {
		d.cards = append(make([]Card, 0, len(src.cards)), src.cards...)
	}
	if cap(d.Weapons) >= len(src.Weapons) {
		d.Weapons = append(d.Weapons[:0], src.Weapons...)
	} else {
		d.Weapons = append(make([]Weapon, 0, len(src.Weapons)), src.Weapons...)
	}
}

// Shuffle randomises the deck in place via Fisher-Yates. Mutates the receiver — callers
// running independent trials should Copy() the master deck first.
func (d *Deck) Shuffle(rng *rand.Rand) {
	for i := len(d.cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

// PeekTop returns the top card of the deck without removing it, or nil when the deck is
// empty. Used by cards that get a buff based on the top card's type / cost / etc.
// (e.g. "if the top card of your deck is an attack action, this gets +1{p}"). Tests
// reading the whole deck should drive Draw / PutBottom instead of reaching for a Peek-all
// API — the deck is meant to be a black box past its top.
func (d *Deck) PeekTop() Card {
	if d == nil || len(d.cards) == 0 {
		return nil
	}
	return d.cards[0]
}

// PeekTopN returns the top n cards of the deck (top first) without removing them, or
// fewer when the deck has < n cards. Reserved for cards whose printed effect literally
// reads "reveal the top N cards" (Sutcliffe's Research Notes and similar) — not for tests
// that want a back-door view of the deck's full contents. The returned slice aliases the
// deck's backing storage; mutating it would corrupt the deck.
func (d *Deck) PeekTopN(n int) []Card {
	if d == nil {
		return nil
	}
	if n > len(d.cards) {
		n = len(d.cards)
	}
	return d.cards[:n]
}

// Draw removes the top n cards from the deck and returns them (top to bottom).
// The returned slice aliases the deck's backing storage; same retention caveat as Peek.
// Panics when n exceeds Size.
func (d *Deck) Draw(n int) []Card {
	if n > len(d.cards) {
		panic(fmt.Sprintf("deck: Draw(%d) exceeds remaining size %d", n, len(d.cards)))
	}
	out := d.cards[:n]
	d.cards = d.cards[n:]
	return out
}

// PutBottom appends cards to the bottom of the deck, preserving the relative order
// passed in. Used by the per-turn loop to recycle pitched cards onto the deck bottom per
// FaB's end-of-turn pitch-zone-to-deck rule.
func (d *Deck) PutBottom(cards []Card) {
	d.cards = append(d.cards, cards...)
}

// PutTop prepends cards to the top of the deck, preserving the relative order passed
// in (cards[0] becomes the new top). Used by mid-chain effects that put a card back on
// top (PrependToDeck) or that re-order the top N (Opt).
func (d *Deck) PutTop(cards []Card) {
	if len(cards) == 0 {
		return
	}
	combined := make([]Card, 0, len(cards)+len(d.cards))
	combined = append(combined, cards...)
	combined = append(combined, d.cards...)
	d.cards = combined
}

// Tutor scans the entire deck, removes the highest-scoring card per score, and returns it.
// Returns (nil, false) when no card scores > 0 or the deck is empty. Used by mid-chain
// tutor effects ("search your deck for a … with X").
func (d *Deck) Tutor(score func(Card) int) (Card, bool) {
	bestIdx := -1
	bestScore := 0
	for i, c := range d.cards {
		sc := score(c)
		if sc > bestScore {
			bestScore = sc
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return nil, false
	}
	found := d.cards[bestIdx]
	out := make([]Card, 0, len(d.cards)-1)
	out = append(out, d.cards[:bestIdx]...)
	out = append(out, d.cards[bestIdx+1:]...)
	d.cards = out
	return found, true
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
	for _, c := range d.cards {
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
func Random(h Hero, size, maxCopies int, rng *rand.Rand, reg Registry) *Deck {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 1 (got %d)", maxCopies))
	}
	loadouts := weaponLoadouts(reg.LegalWeapons())
	if len(loadouts) == 0 {
		panic("deck: Random has no legal weapon loadouts — registry rejected every weapon")
	}
	weapons := loadouts[rng.Intn(len(loadouts))]

	pool := reg.LegalCards()
	if len(pool) == 0 {
		panic("deck: Random has no legal cards — cannot build a deck")
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
