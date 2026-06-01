// Package deck represents a candidate FaB deck — hero, weapons, cards, plus sideboard /
// equipment lists. Depends only on the narrow Hero / Weapon / Card / Registry contracts
// declared here, so callers can swap in any concrete implementation.
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Deck is a hero, equipped weapons, and a deck of cards. Sideboard is the user-managed
// reserve-card list; Equipment is the non-weapon arena loadout (head, chest, arms, legs).
// Both round-trip through persistence; the simulator never reads them.
//
// Sideboard and Equipment are []string rather than []Card so non-registry entries (equipment
// pieces, Nullrune cycle, ...) round-trip without a lossy registry lookup.
//
// cards doubles as the runtime deck: Shuffle / Draw / PeekTop / PutBottom / PutTop / Tutor
// mutate it directly. Evaluation trials should Copy() first to keep the master shareable.
// The field is unexported; composition inspection goes through UniqueIDs / NameCounts /
// DisplayNames / PitchCounts.
type Deck struct {
	Hero Hero
	// Format is the constructed format the deck is built and optimized for. Together with
	// Hero it scopes the legal card pool (Registry.LegalCardsFor); both are durable deck
	// attributes that persist with the deck. New defaults it to format.SilverAge; Random and
	// the textio loaders set it explicitly. Empty-wrapper constructors (NewShallowSafe, the
	// nil-receiver Copy) leave it nil until rebound from a source deck.
	Format    format.Format
	Weapons   []Weapon
	cards     []Card
	Sideboard []string
	Equipment []string
	// mustNotShuffle marks wrappers that share slice backing with another *Deck (produced by
	// ShallowCopy). Shuffle panics on these to prevent silently corrupting peer wrappers.
	mustNotShuffle bool
}

// UniqueIDs returns the distinct card IDs in first-appearance order, plus a lookup map giving
// each ID's index in that slice.
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

// Without returns a copy of the deck with every copy of id removed (a smaller deck; Weapons,
// Sideboard, and Equipment carried over verbatim, surviving cards' order preserved).
func (d *Deck) Without(id ids.CardID) *Deck {
	out := d.Copy()
	kept := out.cards[:0]
	for _, c := range out.cards {
		if c.ID() != id {
			kept = append(kept, c)
		}
	}
	out.cards = kept
	return out
}

// Count returns how many copies of id the deck holds.
func (d *Deck) Count(id ids.CardID) int {
	n := 0
	for _, c := range d.cards {
		if c.ID() == id {
			n++
		}
	}
	return n
}

// New constructs a Deck. Panics if the weapon loadout violates the "0–2 weapons; if 2, both
// 1H" equipment rule. Format defaults to format.SilverAge (the only format today) so a
// New-built deck always has a non-nil format for the mutation pool; callers building for a
// specific format set Format afterwards (Random and the textio loaders do). Sideboard and
// Equipment start empty; callers assign them directly when carrying them over.
func New(h Hero, weapons []Weapon, cards []Card) *Deck {
	validateWeapons(weapons)
	return &Deck{Hero: h, Format: format.SilverAge, Weapons: weapons, cards: cards}
}

// NewShallowSafe returns an empty *Deck wrapper marked mustNotShuffle, ready to be
// rebound via ShallowCopyFrom against a source deck.
func NewShallowSafe() *Deck {
	return &Deck{mustNotShuffle: true}
}

// Size reports the number of cards in the deck (Sideboard / Equipment excluded). Nil-safe:
// a nil Deck reports 0.
func (d *Deck) Size() int {
	if d == nil {
		return 0
	}
	return len(d.cards)
}

// Fingerprint returns a comparable summary of the deck — sorted weapon names and sorted
// card-count histogram — so order-insensitive equality checks compare mutated decks against
// a baseline without caring about positional shuffles.
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
// and Equipment. Hero is shared (heroes are stateless). Per-worker trials Copy() first so
// the master stays sharable across goroutines.
func (d *Deck) Copy() *Deck {
	if d == nil {
		return &Deck{}
	}
	out := &Deck{Hero: d.Hero, Format: d.Format}
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

// ShallowCopy returns a fresh *Deck wrapper that shares slice backing with the receiver.
// cards / Weapons are sliced with cap=len so any future append allocates fresh rather than
// writing past the shared region. Safe with read paths and copy-on-write mutators (PutTop /
// PutBottom / Tutor); unsafe with Shuffle, which mutates in place and panics here.
func (d *Deck) ShallowCopy() *Deck {
	if d == nil {
		return &Deck{}
	}
	out := &Deck{Hero: d.Hero, Format: d.Format, mustNotShuffle: true}
	if len(d.cards) > 0 {
		out.cards = d.cards[:len(d.cards):len(d.cards)]
	}
	if len(d.Weapons) > 0 {
		out.Weapons = d.Weapons[:len(d.Weapons):len(d.Weapons)]
	}
	out.Sideboard = d.Sideboard
	out.Equipment = d.Equipment
	return out
}

// ShallowCopyFrom resets d to mirror src with shared slice backings, writing into the
// receiver instead of allocating. cards / Weapons are sliced with cap=len so future appends
// allocate fresh. Same Shuffle-safety contract as ShallowCopy.
func (d *Deck) ShallowCopyFrom(src *Deck) {
	if src == nil {
		d.Hero = nil
		d.Format = nil
		d.cards = nil
		d.Weapons = nil
		d.Sideboard = nil
		d.Equipment = nil
		return
	}
	d.Hero = src.Hero
	d.Format = src.Format
	if len(src.cards) > 0 {
		d.cards = src.cards[:len(src.cards):len(src.cards)]
	} else {
		d.cards = nil
	}
	if len(src.Weapons) > 0 {
		d.Weapons = src.Weapons[:len(src.Weapons):len(src.Weapons)]
	} else {
		d.Weapons = nil
	}
	d.Sideboard = src.Sideboard
	d.Equipment = src.Equipment
}

// CopyFrom provides a memory-efficient way to copy a Deck by reusing an already-allocated
// Deck (the receiver) that's no longer needed: d's cards / Weapons / Hero / Format are
// overwritten from src, reusing d's existing slice backing arrays when they have enough
// capacity. Sideboard and Equipment are skipped. Callers that need a full deep copy should
// use Copy.
//
// d must be a non-nil receiver. A nil src is treated as an empty deck — d's slice lengths
// drop to 0 (backings retained) and Hero / Format zero.
func (d *Deck) CopyFrom(src *Deck) {
	if src == nil {
		d.Hero = nil
		d.Format = nil
		d.cards = d.cards[:0]
		d.Weapons = d.Weapons[:0]
		return
	}
	d.Hero = src.Hero
	d.Format = src.Format
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

// Shuffle randomises the deck in place via Fisher-Yates. Mutates the receiver. Panics on
// wrappers produced by ShallowCopy (mustNotShuffle set): in-place shuffle would silently
// corrupt peer wrappers sharing the same slice backing.
func (d *Deck) Shuffle(rng *rand.Rand) {
	if d.mustNotShuffle {
		panic("deck: Shuffle called on a ShallowCopy-produced wrapper — a card mutated the per-permutation deck via Shuffle, which would corrupt sibling permutations sharing the same slice backing")
	}
	for i := len(d.cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

// PeekTop returns the top card of the deck without removing it, or nil when empty. Used by
// cards that read the top card to compute a buff. The deck is otherwise a black box past
// its top — tests reading the full deck should drive Draw / PutBottom instead.
func (d *Deck) PeekTop() Card {
	if d == nil || len(d.cards) == 0 {
		return nil
	}
	return d.cards[0]
}

// PeekTopN returns the top n cards (top first) without removing them; fewer when the deck
// has < n. Reserved for cards whose printed effect reveals the top N — not a back-door
// inspection API. The returned slice aliases the deck's backing storage; do not mutate.
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

// PutBottom appends cards to the bottom of the deck, preserving the input order. Used by
// the end-of-turn loop to recycle pitched cards per FaB's pitch-zone-to-deck rule.
func (d *Deck) PutBottom(cards []Card) {
	d.cards = append(d.cards, cards...)
}

// PutTop prepends cards to the top of the deck, preserving the relative order passed
// in (cards[0] becomes the new top). Used by mid-attack-turn effects that put a card back on
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
// Returns (nil, false) when no card scores > 0 or the deck is empty. Used by mid-attack-turn
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

// ApplyDefaults reconciles d.Equipment and d.Sideboard with the supplied defaults. Equipment is
// topped up to one copy per entry. The sideboard is first trimmed so main-deck + sideboard copies
// of any card stay within SideboardCopyCap, then topped up toward each default entry's Count within
// the same cap. Idempotent.
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

	// Trim the existing sideboard so main-deck + sideboard copies of any card never exceed the cap
	// (a card now in the main shouldn't also sit in the sideboard past the cap); keeps the first
	// (cap - main) copies of each, preserving order.
	sideCounts := map[string]int{}
	kept := make([]string, 0, len(d.Sideboard))
	for _, name := range d.Sideboard {
		if mainCounts[name]+sideCounts[name] < SideboardCopyCap {
			kept = append(kept, name)
			sideCounts[name]++
		}
	}
	d.Sideboard = kept

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

// Random generates a random legal deck for hero h in format f: a random weapon loadout from
// reg.LegalWeaponsFor (one 2H or two 1H; dual-wielding the same weapon allowed) and size
// cards drawn uniformly from reg.LegalCardsFor one at a time, skipping any roll that would
// exceed maxCopies for the picked card's ID. The deck's Format is set to f.
func Random(h Hero, f format.Format, size, maxCopies int, rng *rand.Rand, reg Registry) *Deck {
	if maxCopies < 1 {
		panic(fmt.Sprintf("deck: Random requires maxCopies >= 1 (got %d)", maxCopies))
	}
	loadouts := weaponLoadouts(reg.LegalWeaponsFor(f, h))
	if len(loadouts) == 0 {
		panic("deck: Random has no legal weapon loadouts — registry rejected every weapon")
	}
	weapons := loadouts[rng.Intn(len(loadouts))]

	pool := reg.LegalCardsFor(f, h)
	if len(pool) == 0 {
		panic("deck: Random has no legal cards — cannot build a deck")
	}
	counts := map[ids.CardID]int{}
	picks := make([]Card, 0, size)
	for len(picks) < size {
		c := pool[rng.Intn(len(pool))]
		if counts[c.ID()]+1 > effectiveMaxCopies(c, maxCopies) {
			continue
		}
		counts[c.ID()]++
		picks = append(picks, c)
	}
	d := New(h, weapons, picks)
	d.Format = f
	return d
}

// effectiveMaxCopies returns 1 for Legendary cards and maxCopies otherwise.
func effectiveMaxCopies(c Card, maxCopies int) int {
	if _, ok := c.(interface{ Legendary() }); ok {
		return 1
	}
	return maxCopies
}
