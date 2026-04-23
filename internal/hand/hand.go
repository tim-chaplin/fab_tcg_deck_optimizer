// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
	Held
	// Arsenal marks the card placed into the arsenal at end of turn. Contributes nothing to this
	// turn's Value (it scores on the future turn it's played) and carries across the turn boundary
	// as Play.ArsenalCard. At most one card in BestLine takes this role; when the partition
	// enumerator leaves arsenal empty, Best post-hoc promotes one Held card.
	Arsenal
)

// CardAssignment is a single card + the role it took this turn. Hand cards produce one per card;
// an arsenal-in card contributes one with FromArsenal set so a turn fits in one slice.
// Contribution is the per-card credit toward TurnSummary.Value (damage dealt, block share, or
// pitch resource depending on Role), filled by fillContributions once the winner is picked.
type CardAssignment struct {
	Card         card.Card
	Role         Role
	FromArsenal  bool
	Contribution float64
}

// TurnSummary is the result of running Best on a hand: the winning card-role assignments plus
// aggregate metadata about the turn.
type TurnSummary struct {
	// BestLine is the winning partition. Hand cards come first in canonical (post-sort) order;
	// the previous-turn arsenal card, if any, is the last entry with FromArsenal=true. Never
	// mutate — memoized results alias this slice.
	BestLine []CardAssignment
	// AttackChain is the winning chain in play order: attack-role cards from BestLine interleaved
	// with any swung weapons at the positions the solver picked. Each entry carries its Play-time
	// damage (plus hero-trigger damage for cards) so callers can attribute contribution for
	// weapons, which have no BestLine entry. Swung weapons are recoverable by type-asserting
	// AttackChainEntry.Card to weapon.Weapon. Empty when no attacks were played.
	AttackChain []AttackChainEntry
	// Value is the turn's total score (damage dealt + damage prevented). Equals the sum of
	// BestLine[].Contribution plus any weapon-swing damage (weapons aren't in BestLine).
	Value int
	// LeftoverRunechants is the Runechant token count at end of the chosen chain; the caller
	// feeds it back as the next turn's carryover. For partitions with no attacks, equals the
	// carryover the caller passed in.
	LeftoverRunechants int
	// ArsenalCard is the card occupying the arsenal slot at end of turn — either a hand card
	// just arsenaled (role=Arsenal) or a previous-turn arsenal card that stayed. Nil when empty.
	// The caller feeds this back as the next turn's arsenalCardIn.
	ArsenalCard card.Card
	// Drawn is the cards the winning chain drew mid-turn, in draw order, each paired with the
	// disposition the solver picked for it. Populated from state.Drawn during fillContributions's
	// tracked replay. Role is one of Pitch (consumed to fund the chain, Contribution = Pitch()),
	// Attack (played as a free-cost chain extension, Contribution = damage dealt), Arsenal
	// (promoted into an empty slot post-enumeration, Contribution 0), or Held (carries into
	// the next hand, Contribution 0). Nil when no draw rider fired.
	Drawn []CardAssignment
	// AuraTriggers is the surviving AuraTrigger set at end of this turn — triggers added by
	// this turn's winning Play chain. The deck loop feeds this into next turn's start-of-turn
	// trigger pass, closing the cross-turn loop. Nil when the turn played no trigger-creating
	// aura.
	AuraTriggers []card.AuraTrigger
	// TriggersFromLastTurn records the AuraTriggers whose start-of-turn handlers fired at the
	// top of this turn, each with the damage-equivalent it credited. Populated by the deck
	// loop before FormatBestTurn is called; Value already includes the sum.
	TriggersFromLastTurn []TriggerContribution
}

// TriggerContribution is one start-of-turn AuraTrigger fire: the aura that fired plus the
// Damage it credited (folded into Value). Surfaced in TurnSummary.TriggersFromLastTurn so
// FormatBestTurn can print a "(from previous turn)" line naming the outcome.
type TriggerContribution struct {
	Card   card.Card
	Damage int
}

// AttackChainEntry is a single played attack — a card with role=Attack or a swung weapon —
// carrying the damage it contributed when it resolved in the winning chain. Damage is the Play()
// return; TriggerDamage is the hero's OnCardPlayed contribution (e.g. Viserai creating a
// Runechant) so callers can surface hero attribution on its own line. For BestLine Attack entries
// Damage + TriggerDamage equals CardAssignment.Contribution; weapons live only here.
type AttackChainEntry struct {
	Card          card.Card
	Damage        float64
	TriggerDamage float64
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal, if any.
// Lets callers treat the arsenal-in card differently from hand cards without scanning BestLine.
func (t TurnSummary) ArsenalIn() (CardAssignment, bool) {
	for _, a := range t.BestLine {
		if a.FromArsenal {
			return a, true
		}
	}
	return CardAssignment{}, false
}

// String returns a human-readable role name ("PITCH", "ATTACK", "DEFEND", "HELD", "ARSENAL").
func (r Role) String() string {
	switch r {
	case Pitch:
		return "PITCH"
	case Attack:
		return "ATTACK"
	case Defend:
		return "DEFEND"
	case Held:
		return "HELD"
	case Arsenal:
		return "ARSENAL"
	}
	return "UNKNOWN"
}

// formatContribution renders a contribution/damage value for the best-turn printout. Integers
// render bare; fractional values (proportional defense share when multiple blockers split an
// incoming attack) show one decimal place.
func formatContribution(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%.1f", v)
}

// assignmentName returns the card name, suffixed with " (from arsenal)" when the assignment came
// from the arsenal slot — that tag tells readers why the card isn't in the dealt-hand list the
// optimiser reports alongside.
func assignmentName(a CardAssignment) string {
	if a.FromArsenal {
		return a.Card.Name() + " (from arsenal)"
	}
	return a.Card.Name()
}

// FormatBestLine pairs each card in BestLine with its assigned role for debug output, e.g.
// "Hocus Pocus (Blue): PITCH, Runic Reaping (Red): ATTACK". Compact one-line form; use
// FormatBestTurn for chronological play order.
func FormatBestLine(line []CardAssignment) string {
	parts := make([]string, len(line))
	for i, a := range line {
		parts[i] = assignmentName(a) + ": " + a.Role.String()
	}
	return strings.Join(parts, ", ")
}

// splitPitchesByPhase assigns each pitch card to the defense or attack phase, simulating the
// order FaB prompts them in. Smallest pitches fund the defense bucket until drCost is covered;
// the rest pay for this turn's attacks. Stable on ties so display order is deterministic.
func splitPitchesByPhase(pitched []CardAssignment, drCost int) (defensePitches, attackPitches []CardAssignment) {
	sorted := append([]CardAssignment(nil), pitched...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Card.Pitch() < sorted[j].Card.Pitch() })
	covered := 0
	for _, a := range sorted {
		if covered < drCost {
			defensePitches = append(defensePitches, a)
			covered += a.Card.Pitch()
		} else {
			attackPitches = append(attackPitches, a)
		}
	}
	return defensePitches, attackPitches
}

// appendAttackChainLines renders the Attack phase: one numbered line per AttackChain entry in
// solver-chosen play order, with the shared step counter advanced via stepPtr so later sections
// keep numbering contiguous. Non-weapon entries cross-reference BestLine by ID so arsenal-played
// cards get a "(from arsenal)" tag; weapons skip the match since they have no BestLine entry.
// Cards that aren't attacks (e.g. non-attack actions like Mauvrion Skies) use "PLAY" so the
// label matches what the card actually does on the chain. A non-zero TriggerDamage adds a
// trailing " (+M hero trigger)" so the attribution is visible instead of silently folded into
// the card's own damage number.
func appendAttackChainLines(lines []string, t TurnSummary, stepPtr *int) []string {
	used := make([]bool, len(t.BestLine))
	appendAttack := func(label, cardName string, e AttackChainEntry) {
		*stepPtr++
		line := fmt.Sprintf("  %d. %s: %s (+%s)", *stepPtr, cardName, label, formatContribution(e.Damage))
		if e.TriggerDamage > 0 {
			line += fmt.Sprintf(" (+%s hero trigger)", formatContribution(e.TriggerDamage))
		}
		lines = append(lines, line)
	}
	for _, e := range t.AttackChain {
		if _, isWeapon := e.Card.(weapon.Weapon); isWeapon {
			appendAttack("WEAPON ATTACK", e.Card.Name(), e)
			continue
		}
		// Match the first unused Attack-role BestLine entry by ID so we can detect FromArsenal.
		tag := ""
		for i := range t.BestLine {
			if used[i] || t.BestLine[i].Role != Attack || t.BestLine[i].Card.ID() != e.Card.ID() {
				continue
			}
			if t.BestLine[i].FromArsenal {
				tag = " (from arsenal)"
			}
			used[i] = true
			break
		}
		label := "ATTACK"
		if !e.Card.Types().Has(card.TypeAttack) {
			label = "PLAY"
		}
		appendAttack(label, e.Card.Name()+tag, e)
	}
	return lines
}

// FormatBestTurn renders a TurnSummary as a numbered play-order list, one card per line,
// matching the actual FaB turn sequence:
//
//  1. Defense-phase pitches (paying for Defense Reactions)
//  2. Plain blocks
//  3. Defense Reactions
//  4. Attack-phase pitches (paying for this turn's played cards)
//  5. Attack chain — played cards and swung weapons in the order the solver picked
//
// Held / Arsenal cards are summarized on trailing lines so the reader sees what's carrying over.
//
// Pitch-phase assignment uses a greedy split for display: smallest pitches first fund the defense
// pool until drCost is covered, the rest fund attack. The solver already validated some legal
// split exists; this picks one deterministically.
func FormatBestTurn(t TurnSummary) string {
	parts := partitionBestLineForDisplay(t.BestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	var lines []string
	step := 0
	nextStep := func() int { step++; return step }
	// Pitch lines don't show a damage contribution: pitches generate resource, not damage —
	// showing "+3" next to a pitch would double-count against the turn's Value.
	appendPitch := func(a CardAssignment, roleLabel string) {
		lines = append(lines, fmt.Sprintf("  %d. %s: %s", nextStep(), assignmentName(a), roleLabel))
	}
	// Defense lines contribute damage-prevention (block share + DR Play return), so the tag is
	// shown.
	appendDefense := func(a CardAssignment, roleLabel string) {
		lines = append(lines, fmt.Sprintf("  %d. %s: %s (+%s prevented)", nextStep(), assignmentName(a), roleLabel, formatContribution(a.Contribution)))
	}

	for _, a := range defensePitches {
		appendPitch(a, "PITCH (opponent's turn)")
	}
	for _, a := range parts.plainBlocks {
		appendDefense(a, "BLOCK")
	}
	for _, a := range parts.defenseReactions {
		appendDefense(a, "DEFENSE REACTION")
	}
	for _, d := range t.TriggersFromLastTurn {
		step++
		lines = append(lines, fmt.Sprintf("  %d. %s (from previous turn): START OF ACTION PHASE (+%d)",
			step, d.Card.Name(), d.Damage))
	}
	for _, a := range attackPitches {
		appendPitch(a, "PITCH (my turn)")
	}
	lines = appendAttackChainLines(lines, t, &step)
	lines = appendHeldArsenalFooter(lines, parts.held, parts.arsenal, t.Drawn)
	return strings.Join(lines, "\n")
}

// bestLineDisplayParts groups BestLine entries by the display section each belongs to. Pitches
// pool before being split into defense / attack phases; blocks split again by whether the card
// is a Defense Reaction (which has its own "DEFENSE REACTION" tag). drCost sums Defense-Reaction
// costs so splitPitchesByPhase can decide how much of the pitch pool funds the opponent's turn.
type bestLineDisplayParts struct {
	pitched          []CardAssignment
	plainBlocks      []CardAssignment
	defenseReactions []CardAssignment
	held             []CardAssignment
	arsenal          []CardAssignment
	drCost           int
}

// partitionBestLineForDisplay sorts the winning line into the buckets FormatBestTurn renders
// section-by-section. Defenders split on DR membership so DR-only lines get the right label
// and their cost contributes to the defense-phase pitch target.
func partitionBestLineForDisplay(line []CardAssignment) bestLineDisplayParts {
	var parts bestLineDisplayParts
	zeroState := &card.TurnState{}
	for _, a := range line {
		switch a.Role {
		case Pitch:
			parts.pitched = append(parts.pitched, a)
		case Attack:
			// Attack-phase cost sum is computed here to match the turn's modeling, but is not
			// surfaced in the rendered output — the attack chain's per-card lines already show
			// damage credit rather than cost.
			_ = a.Card.Cost(zeroState)
		case Defend:
			if a.Card.Types().IsDefenseReaction() {
				parts.drCost += a.Card.Cost(zeroState)
				parts.defenseReactions = append(parts.defenseReactions, a)
			} else {
				parts.plainBlocks = append(parts.plainBlocks, a)
			}
		case Held:
			parts.held = append(parts.held, a)
		case Arsenal:
			parts.arsenal = append(parts.arsenal, a)
		}
	}
	return parts
}

// appendHeldArsenalFooter appends the trailing "(held: ...)" / "(arsenal: ...)" lines that show
// unplayed cards outside the numbered sequence. Mid-turn-drawn Held / Arsenal cards render with
// a "(drawn)" suffix so the reader can tell them from starting-hand entries; the arsenal card
// itself shows "(stayed)" vs "(new)" so staying-in-place is distinguishable from being newly
// placed this turn.
func appendHeldArsenalFooter(lines []string, held, arsenal, drawn []CardAssignment) []string {
	var footers []string
	for _, a := range held {
		footers = append(footers, fmt.Sprintf("  (held: %s)", a.Card.Name()))
	}
	for _, d := range drawn {
		if d.Role == Held {
			footers = append(footers, fmt.Sprintf("  (held: %s (drawn))", d.Card.Name()))
		}
	}
	for _, a := range arsenal {
		label := a.Card.Name()
		if a.FromArsenal {
			label += " (stayed)"
		} else {
			label += " (new)"
		}
		footers = append(footers, fmt.Sprintf("  (arsenal: %s)", label))
	}
	for _, d := range drawn {
		if d.Role == Arsenal {
			footers = append(footers, fmt.Sprintf("  (arsenal: %s (drawn))", d.Card.Name()))
		}
	}
	return append(lines, footers...)
}

// Best returns the optimal TurnSummary for the given hand against an opponent that will attack
// for incomingDamage on their next turn. Equipped weapons may be swung for their Cost if
// resources allow.
//
// Cards partition into five roles:
//   - Pitch: contributes its Pitch value as resource paying for a played card this turn or a
//     Defense Reaction on the opponent's turn.
//   - Attack: consumes Cost on our turn; resolved by calling Card.Play in an order the optimizer
//     chooses. TurnState effects carry forward to later attacks in the same sequence.
//   - Defend: contributes Defense to damage prevented (capped at incomingDamage; excess wasted).
//     Plain blocking is free; Defense Reactions must pay Cost and contribute any Play() damage.
//   - Held: stays in hand for next turn. Contributes nothing this turn.
//   - Arsenal: moves into the arsenal slot at end of turn, or for an arsenal-in card, stays.
//     Contributes nothing this turn.
//
// Pitch resources split across two phases because resources don't carry between turns: attack
// pitches pay for this turn's played cards, defense pitches pay for Defense Reactions on the
// opponent's turn. A card can only be pitched while some unpaid card is on the stack in the
// matching phase, so a hand with no plays in a phase has no Pitch-role cards in that phase; any
// card that can't be legally pitched becomes Held.
//
// Results are memoized on (heroID, sorted weapon IDs, sorted card IDs, incomingDamage,
// runechantCarryover, arsenal-in ID) so repeat evaluations short-circuit. The hand is sorted in
// place into canonical order first; BestLine's hand entries align with that post-sort order.
// Every card in the hand must be registered in package cards or Best panics.
//
// runechantCarryover is the Runechant token count carrying in from the previous turn.
// TurnSummary.LeftoverRunechants is the count at end of the chosen chain; feed it back as the
// next turn's carryover.
//
// arsenalCardIn is the card sitting in the arsenal slot at start of turn (nil if empty). The
// enumerator pulls it in as an extra CardAssignment with restricted role options — Arsenal
// (stay), Attack (any non-DR card), or Defend (Defense Reactions only). Never Pitch or Held. A
// hand card may also take the Arsenal role so long as at most one BestLine entry ends up there.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return sharedEvaluator.Best(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn)
}

// BestWithTriggers is the package-level counterpart of Evaluator.BestWithTriggers, using the
// shared evaluator. Pass priorAuraTriggers to feed the cross-turn trigger carry into the
// search — those triggers may fire mid-chain (e.g. Malefic Incantation's TriggerAttackAction
// rune) and contribute damage to this turn's Value.
func BestWithTriggers(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	return sharedEvaluator.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
}

// Best is the method form of the package-level Best: same semantics, uses this Evaluator's
// scratch buffers so concurrent goroutines can each hold their own.
func (e *Evaluator) Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return e.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, nil)
}

// BestWithTriggers is Best plus an explicit priorAuraTriggers input — the AuraTriggers
// carrying in from the previous turn. Non-empty priorAuraTriggers disables memoization: the
// triggers contain Handler closures that aren't comparable, and the sim mutates trigger
// Count / FiredThisTurn mid-chain. With nil priorAuraTriggers, BestWithTriggers matches
// Best exactly (fully memoable).
func (e *Evaluator) BestWithTriggers(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	// IDs go into a fixed-size stack array to avoid a per-call slice alloc. Hand size is capped
	// at 8 (matches memoKey.cardIDs); larger hands panic out of the inner loops.
	n := len(hand)
	var ids [8]card.ID
	memoable := len(priorAuraTriggers) == 0
	for i, c := range hand {
		ids[i] = c.ID()
		if _, ok := c.(card.NoMemo); ok {
			memoable = false
		}
	}
	if arsenalCardIn != nil {
		if _, ok := arsenalCardIn.(card.NoMemo); ok {
			memoable = false
		}
	}

	sortHandByID(hand, ids[:], n)

	key := makeMemoKey(hero, weapons, &ids, n, incomingDamage, runechantCarryover, arsenalCardIn)
	if memoable {
		memoMu.RLock()
		cached, hit := memo[key]
		memoMu.RUnlock()
		if hit {
			return cached
		}
	}
	result := e.bestUncached(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
	if memoable {
		memoMu.Lock()
		memo[key] = result
		memoMu.Unlock()
	}
	return result
}

// sortHandByID sorts the first n entries of hand and ids in parallel by ascending id, in place.
// Insertion sort — for n ≤ 8 this beats sort.Sort and avoids boxing slices through sort.Interface.
// Canonicalizing the hand order is what lets the memo key collapse permutations onto one entry.
func sortHandByID(hand []card.Card, ids []card.ID, n int) {
	for i := 1; i < n; i++ {
		for j := i; j > 0 && ids[j-1] > ids[j]; j-- {
			ids[j-1], ids[j] = ids[j], ids[j-1]
			hand[j-1], hand[j] = hand[j], hand[j-1]
		}
	}
}

// memoKey is the comparable map key for the shared memo. Hand size is capped at 8 cards. Hero
// ID + weapon IDs live in the key so Evaluators with different (hero, weapons) tuples coexist in
// the memo without a scope wipe — distinct scopes just produce distinct keys. The uint16 hero.ID
// keeps the whole key a fixed-size integer struct; no string hashing per lookup.
type memoKey struct {
	heroID             hero.ID
	weaponIDs          [2]card.ID
	cardIDs            [8]card.ID
	cardCount          uint8
	incoming           int
	runechantCarryover int
	// arsenalInID is card.Invalid when the slot is empty, otherwise the ID of the starting
	// arsenal card — different arsenal-ins give distinct cache entries.
	arsenalInID card.ID
}

// Evaluator owns the per-goroutine mutable state hand.Best threads through an evaluation: a
// scratch-buffer cache keyed by (handSize, weapons). The memo cache is shared across all
// Evaluators so every worker benefits from cached hands — only the mutated scratch buffers
// must be per-goroutine. A long-lived Evaluator avoids reallocating ~20 scratch slices per call.
type Evaluator struct {
	// bufs holds the pre-allocated scratch slices for bestUncached / bestAttackWithWeapons /
	// bestSequence. Keyed by (handSize, weaponCount, weaponIDs); recreated when any differ so the
	// scratch sizing stays correct.
	bufs            *attackBufs
	bufsHandSize    int
	bufsWeaponIDs   [2]card.ID
	bufsWeaponCount int
	bufsValid       bool
}

// NewEvaluator returns a fresh Evaluator with empty scratch buffers. One Evaluator per goroutine
// is safe for concurrent use and still shares the global memo.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// sharedEvaluator backs the package-level Best — single-threaded callers don't need to construct
// their own. Parallel callers create one per goroutine for bufs but all share the global memo.
var sharedEvaluator = NewEvaluator()

// memo caches canonical-order TurnSummary results keyed by memoKey. Shared across all Evaluators
// and goroutines; protected by memoMu. Hero + weapon IDs live in memoKey so distinct (hero,
// weapons) scopes coexist without a wipe step.
var (
	memo   = map[memoKey]TurnSummary{}
	memoMu sync.RWMutex
)

// ClearMemo drops every cached TurnSummary. Callers (iterate-mode, benchmarks) use this to cap
// memo growth across unrelated runs; cross-run hit rate is near zero so nothing of value is lost.
func ClearMemo() {
	memoMu.Lock()
	clear(memo)
	memoMu.Unlock()
}

// MemoLen returns the current number of cached entries, for diagnostic logging.
func MemoLen() int {
	memoMu.RLock()
	n := len(memo)
	memoMu.RUnlock()
	return n
}

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID.
// sortedIDs is a pointer to the caller's stack [8]card.ID to avoid a slice-header escape. Weapon
// IDs are sorted into the two fixed slots so loadouts in any order hash to the same key.
func makeMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs *[8]card.ID, n int, incoming int, runechantCarryover int, arsenalCardIn card.Card) memoKey {
	k := memoKey{
		heroID:             hero.ID(),
		incoming:           incoming,
		runechantCarryover: runechantCarryover,
		cardCount:          uint8(n),
		cardIDs:            *sortedIDs,
	}
	if arsenalCardIn != nil {
		k.arsenalInID = arsenalCardIn.ID()
	}
	switch len(weapons) {
	case 1:
		k.weaponIDs[0] = weapons[0].ID()
	case 2:
		a, b := weapons[0].ID(), weapons[1].ID()
		if a > b {
			a, b = b, a
		}
		k.weaponIDs[0], k.weaponIDs[1] = a, b
	}
	return k
}

// attackerMeta caches the scalar card attributes playSequence reads on every permutation. With
// this hoisted to a per-attacker lookup, the hot inner loop skips Types / GoAgain interface
// dispatch; the one meta build amortises across all N! permutations.
//
// minCost / maxCost are static bounds on Card.Cost(s). For cards implementing card.VariableCost
// the solver uses them for O(1) partition pre-screens and falls through to card.Cost(state) in
// the chain inner loop. For non-VariableCost cards, minCost == maxCost == Cost(&TurnState{})
// and the cached value is used directly (no interface call per play).
type attackerMeta struct {
	types            card.TypeSet
	card             card.Card // held for variable-cost chain-time Cost(state) calls
	minCost          int
	maxCost          int
	isVariable       bool
	baseGoAgain      bool
	isAttackOrWeapon bool
	// isAttackAction is the "attack action card" test (Action+Attack, no Weapon) the sim uses
	// to pick which Play resolutions fire TriggerAttackAction AuraTriggers. Weapons carry
	// TypeAttack but aren't attack action CARDS; only the Action+Attack bitmask matches the
	// printed trigger text on cards like Malefic Incantation.
	isAttackAction bool
}

// costAt returns the card's effective cost given the current TurnState. Static cards return the
// cached value directly; variable-cost cards defer to card.Cost(s) so every game-state-dependent
// costing rule lives inside the card, not the solver.
func (m *attackerMeta) costAt(s *card.TurnState) int {
	if m.isVariable {
		return m.card.Cost(s)
	}
	return m.maxCost
}

// cardMetaCache / cardMetaReady are shared, read-only-after-init card metadata tables. Populated
// lazily via cardMetaSlowPath on first encounter, then read from all goroutines without sync.
// Sized for the full uint16 ID space so lookups are plain bounds-checked reads (~2 MB total).
const cardMetaCacheSize = 1 << 16

var (
	cardMetaCache [cardMetaCacheSize]attackerMeta
	cardMetaReady [cardMetaCacheSize]uint32 // written once (atomically) per ID; 0 = unready, 1 = ready
	cardMetaMu    sync.Mutex
)

// attackerMetaPtrFor returns a pointer to cached metadata for c, populating on first encounter.
// Hands back a direct pointer into the global cache so permutation swaps move 8 bytes instead of
// a full attackerMeta struct. The target is read-only after initialisation. Safe from multiple
// goroutines: the first writer per ID holds the mutex, later readers see the ready flag set with
// a release barrier and read the immutable meta entry directly.
func attackerMetaPtrFor(c card.Card) *attackerMeta {
	id := c.ID()
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return &cardMetaCache[id]
	}
	cardMetaSlowPath(c, id)
	return &cardMetaCache[id]
}

// cardMetaSlowPath populates the cache entry under cardMetaMu and returns the computed meta.
func cardMetaSlowPath(c card.Card, id card.ID) attackerMeta {
	cardMetaMu.Lock()
	defer cardMetaMu.Unlock()
	// Re-check under lock: another goroutine may have populated between the atomic load and here.
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return cardMetaCache[id]
	}
	t := c.Types()
	m := attackerMeta{
		types:            t,
		card:             c,
		baseGoAgain:      c.GoAgain(),
		isAttackOrWeapon: t.Has(card.TypeAttack) || t.Has(card.TypeWeapon),
		isAttackAction:   t.IsAttackAction(),
	}
	if vc, ok := c.(card.VariableCost); ok {
		m.minCost = vc.MinCost()
		m.maxCost = vc.MaxCost()
		m.isVariable = m.minCost != m.maxCost
	} else {
		// Static cost: any TurnState probe returns the same value. Cache once.
		fixed := c.Cost(&card.TurnState{})
		m.minCost = fixed
		m.maxCost = fixed
	}
	cardMetaCache[id] = m
	atomic.StoreUint32(&cardMetaReady[id], 1)
	return m
}

// attackBufs holds pre-allocated buffers for the attack-evaluation pipeline (bestSequence →
// playSequence) and the partition loop in bestUncached. Allocated once and cached on the
// Evaluator so a deck eval reuses them across every partition, mask, and permutation.
type attackBufs struct {
	pcBuf          []card.CardState
	ptrBuf         []*card.CardState
	cardsPlayedBuf []card.Card
	state          *card.TurnState
	// drScratch is a pooled TurnState for defense-reaction cost probing inside the
	// (pmask × wmask) loop; reusing its heap slot avoids a per-iteration alloc caused by
	// interface-call escape.
	drScratch card.TurnState
	// drCardStateScratch is a pooled *CardState handed to DR Card.Play calls. Each Play takes
	// a *CardState through an interface boundary so a literal &card.CardState{} would escape
	// and heap-alloc once per DR per partition — reusing this slot keeps the whole defense-phase
	// replay allocation-free. Reset per call by the caller.
	drCardStateScratch card.CardState
	attackerBuf    []card.Card // for bestAttackWithWeapons mask iteration
	// Pre-computed per-mask weapon data. Indexed by bitmask (0 to 2^len(weapons)-1):
	// weaponCosts[mask] is total Cost; weaponNames[mask] is the pre-built []string of names.
	weaponCosts []int
	weaponNames [][]string
	// permMeta parallels pcBuf: each entry points into the global cardMetaCache so playSequence's
	// inner loop skips interface dispatch on Types / GoAgain and reads cached cost bounds.
	// Pointer-valued so bestSequence's permutation swaps move 8 bytes instead of a full struct.
	permMeta []*attackerMeta
	// Partition-loop buffers, consumed by bestUncached. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n. isDRBuf caches TypeDefenseReaction
	// membership to skip Types().Has calls; addsFutureValueBuf caches
	// card.AddsFutureValue implementation so the beatsBest tiebreaker can count how many
	// hidden-future-value cards a partition queues.
	rolesBuf           []Role
	pitchVals          []int
	defenseVals        []int
	isDRBuf            []bool
	addsFutureValueBuf []bool
	// pitchedValsScratch backs the per-leaf "pitched values" slice consumed by phase-mask
	// enumeration. Re-sliced to [:0] at the start of every leaf to eliminate a per-leaf alloc.
	pitchedValsScratch []int
	pitchedBuf         []card.Card
	attackersBuf       []card.Card
	defendersBuf       []card.Card
	// defenseGravScratch / attackGravScratch back state.Graveyard during DR Plays and attack-
	// chain permutations respectively. Reset via [:0]+append per iteration so card effects can
	// freely mutate their view without leaking into the next one. Split so the two phases
	// never alias each other.
	defenseGravScratch []card.Card
	attackGravScratch  []card.Card
	// auraTriggersScratch backs state.AuraTriggers during attack-chain permutations. Reset
	// per permutation so AddAuraTrigger calls in one ordering don't leak into the next.
	auraTriggersScratch []card.AuraTrigger
	// ephemeralTriggersScratch backs state.EphemeralAttackTriggers during attack-chain
	// permutations. Reset per permutation (empty, no cross-turn carry) so one ordering's
	// registrations don't leak into the next.
	ephemeralTriggersScratch []card.EphemeralAttackTrigger
	// drawnWinnerScratch / auraTriggersWinnerScratch back sequenceContext.drawnWinner and
	// auraTriggersWinner. Each sequenceContext borrows them at construction so the eval
	// closure's winner snapshot (append-into-[:0] on every improved permutation) reuses one
	// backing array per Best call. Without this, each of the thousands of sequenceContexts a
	// Best call builds would start with a nil slice and allocate a fresh backing array the
	// first time a winner's AuraTriggers or Drawn is non-empty — the hottest alloc site in the
	// attack-enumeration inner loop. fillContributions clones the winners into summary.Drawn
	// / summary.AuraTriggers before returning so no downstream reader aliases this shared
	// storage.
	drawnWinnerScratch        []card.Card
	auraTriggersWinnerScratch []card.AuraTrigger
	// perCardScratch is sized maxAttackers (handSize + weaponCount). Written by playSequence only
	// when the caller passes a non-nil perCardOut; bestSequence snapshots the winning
	// permutation's per-card damage from here into the caller's output buffer. The partition-loop
	// hot path passes nil and never touches this slice.
	perCardScratch []float64
	// perCardTriggerScratch parallels perCardScratch for hero-trigger damage (OnCardPlayed
	// return). Only written when the caller tracks.
	perCardTriggerScratch []float64
	// fillContribWinnerOrder / fillContribPerCard are output buffers for bestSequence during
	// fillContributions's tracked replay. Kept on attackBufs so each Best call reuses the slab.
	fillContribWinnerOrder []card.Card
	fillContribPerCard     []float64
	fillContribTriggerDmg  []float64
	// fillContribUsed marks hand indices already assigned during chain→hand mapping. Sized
	// handSize; reset with clear before each fillContributions pass.
	fillContribUsed []bool
}

func newAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *attackBufs {
	// +1 reserves a slot for the arsenal-in card, which joins attackers or defenders when the
	// enumerator plays it from arsenal. +maxDrawnExtensions leaves headroom for mid-turn-drawn
	// cards that play as chain extensions — cheap cycling cards (cost 0, Go again, draws a
	// card) can extend a chain well past the starting hand size.
	const maxDrawnExtensions = 32
	maxAttackers := handSize + weaponCount + 1 + maxDrawnExtensions
	numMasks := 1 << weaponCount
	weaponCosts := make([]int, numMasks)
	weaponNames := make([][]string, numMasks)
	for mask := 0; mask < numMasks; mask++ {
		cost := 0
		var names []string
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				cost += w.Cost(&card.TurnState{})
				names = append(names, w.Name())
			}
		}
		weaponCosts[mask] = cost
		weaponNames[mask] = names
	}
	pcBuf := make([]card.CardState, maxAttackers)
	ptrBuf := make([]*card.CardState, maxAttackers)
	// Wire the ptrBuf entries to their pcBuf slots once — the mapping is stable across every
	// permutation so playSequenceWithMeta doesn't need to rewrite it per call.
	for i := range pcBuf {
		ptrBuf[i] = &pcBuf[i]
	}
	return &attackBufs{
		permMeta:                  make([]*attackerMeta, maxAttackers),
		pcBuf:                     pcBuf,
		ptrBuf:                    ptrBuf,
		cardsPlayedBuf:            make([]card.Card, 0, maxAttackers),
		state:                     &card.TurnState{},
		attackerBuf:               make([]card.Card, maxAttackers),
		weaponCosts:               weaponCosts,
		weaponNames:               weaponNames,
		rolesBuf:                  make([]Role, handSize+1),
		pitchVals:                 make([]int, handSize+1),
		defenseVals:               make([]int, handSize+1),
		isDRBuf:                   make([]bool, handSize+1),
		addsFutureValueBuf:        make([]bool, handSize+1),
		pitchedValsScratch:        make([]int, 0, handSize+1),
		pitchedBuf:                make([]card.Card, 0, handSize+1),
		attackersBuf:              make([]card.Card, 0, handSize+1),
		defendersBuf:              make([]card.Card, 0, handSize+1),
		defenseGravScratch:        make([]card.Card, 0, handSize+1),
		attackGravScratch:         make([]card.Card, 0, maxAttackers),
		auraTriggersScratch:       make([]card.AuraTrigger, 0, maxAttackers),
		ephemeralTriggersScratch:  make([]card.EphemeralAttackTrigger, 0, maxAttackers),
		drawnWinnerScratch:        make([]card.Card, 0, maxAttackers),
		auraTriggersWinnerScratch: make([]card.AuraTrigger, 0, maxAttackers),
		perCardScratch:            make([]float64, maxAttackers),
		perCardTriggerScratch:     make([]float64, maxAttackers),
		fillContribWinnerOrder:    make([]card.Card, maxAttackers),
		fillContribPerCard:        make([]float64, maxAttackers),
		fillContribTriggerDmg:     make([]float64, maxAttackers),
		fillContribUsed:           make([]bool, handSize),
	}
}

// getAttackBufs returns this Evaluator's scratch-buffer set, rebuilding when (handSize, weapons)
// changes. Single-slot per Evaluator: iterate runs tens of thousands of same-shape hands, so a
// slot outperforms a keyed pool for this workload.
func (e *Evaluator) getAttackBufs(handSize int, weapons []weapon.Weapon) *attackBufs {
	var wids [2]card.ID
	for i, w := range weapons {
		if i >= len(wids) {
			break
		}
		wids[i] = w.ID()
	}
	if e.bufsValid &&
		e.bufsHandSize == handSize &&
		e.bufsWeaponCount == len(weapons) &&
		e.bufsWeaponIDs == wids {
		return e.bufs
	}
	e.bufs = newAttackBufs(handSize, len(weapons), weapons)
	e.bufsHandSize = handSize
	e.bufsWeaponCount = len(weapons)
	e.bufsWeaponIDs = wids
	e.bufsValid = true
	return e.bufs
}

// fillPartitionPerCardBufs writes the per-card values the partition recurse reads at each leaf:
// Pitch / Defense magnitudes, Defense-Reaction membership, and AddsFutureValue interface
// satisfaction. Computing them up front keeps the recurse's inner body free of card-method /
// type-assert calls, which would otherwise repeat on every leaf. totalN covers the optional
// arsenal-in slot at index n; when present, its Defense picks up ArsenalDefenseBonus so the
// partition / capping pipeline sees the effective value. Returns whether any card is a
// Defense Reaction so the leaf branch can pick between the full three-bucket grouper and the
// faster reaction-free grouper.
func fillPartitionPerCardBufs(hand []card.Card, n, totalN int, arsenalCardIn card.Card, pvals, dvals []int, isDR, addsFutureValue []bool) bool {
	hasReactions := false
	for i := 0; i < totalN; i++ {
		var c card.Card
		if i < n {
			c = hand[i]
		} else {
			c = arsenalCardIn
		}
		pvals[i] = c.Pitch()
		dvals[i] = c.Defense()
		// Arsenal slot (i == n) lives at the end. Defense Reactions whose +N{d} rider only fires
		// when played from arsenal (Unmovable, Springboard Somersault) opt in via
		// card.ArsenalDefenseBonus; bump the static Defense() up here so the partition / capping
		// pipeline sees the effective value.
		if i == n {
			if ab, ok := c.(card.ArsenalDefenseBonus); ok {
				dvals[i] += ab.ArsenalDefenseBonus()
			}
		}
		isDR[i] = c.Types().IsDefenseReaction()
		if isDR[i] {
			hasReactions = true
		}
		_, addsFutureValue[i] = c.(card.AddsFutureValue)
	}
	return hasReactions
}

func (e *Evaluator) bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	n := len(hand)
	// The partition recurse treats the arsenal-in card as an extra entry at index n with a
	// restricted role menu (Arsenal / Attack / Defend), so everything about it is decided inside
	// the enumeration. totalN is the effective size of BestLine.
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	// Seed best.LeftoverRunechants with the carryover: partitions with no attacks don't reduce
	// it, so carryover is the baseline to beat. BestLine starts with every hand card Held and
	// the arsenal-in card (if any) staying in the slot, so a hand with no Value-adding partition
	// still reports sensible "nothing played, nothing pitched" assignments.
	best := TurnSummary{BestLine: make([]CardAssignment, totalN), LeftoverRunechants: runechantCarryover}
	// bestSwung holds the winning partition's swung weapon names so fillContributions can rebuild
	// the chain it runs bestSequence over. Lives outside TurnSummary since weapons are
	// recoverable from AttackChain once fillContributions finishes. bestBudget captures the
	// winning phase-split's chain-resource state; the replay re-seeds ctx with it so
	// bestSequence finds the exact permutation that won during enumeration.
	var bestSwung []string
	var bestBudget chainBudget
	// bestHasHeld tracks whether the current best has at least one Held hand card — lets
	// beatsBest distinguish "arsenal will be occupied post-hoc" from "arsenal will be empty."
	// Seeded true when the hand is non-empty: the initial best puts every hand card into Held,
	// so a post-hoc promotion would fill arsenal. Candidates need both a Value/leftover tie and
	// some way to end with arsenal occupied to displace it.
	bestHasHeld := n > 0
	// bestFutureValuePlayed tracks how many card.AddsFutureValue cards the current best is
	// playing (Role=Attack). Seeded 0 because the initial best assigns every card Held. The
	// beatsBest tiebreaker prefers partitions that play MORE future-value cards at equal
	// Value/leftover — their hidden-later-turn payoff is invisible to the current-turn
	// score, so without this bias a lone sigil loses to Held → arsenal promotion on the
	// arsenal-occupancy tiebreak.
	bestFutureValuePlayed := 0
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
		best.ArsenalCard = arsenalCardIn
	}

	// bufs is the pooled scratch space for this deck evaluation. Partition scratch is sized
	// handSize+1, big enough for totalN when an arsenal-in card inflates the effective hand.
	// Each field is re-sliced and rewritten below, so carry-over from prior calls can't leak.
	bufs := e.getAttackBufs(n, weapons)
	rolesBuf := bufs.rolesBuf[:totalN]
	pvals := bufs.pitchVals[:totalN]
	dvals := bufs.defenseVals[:totalN]
	isDR := bufs.isDRBuf[:totalN]
	addsFutureValue := bufs.addsFutureValueBuf[:totalN]

	hasReactions := fillPartitionPerCardBufs(hand, n, totalN, arsenalCardIn, pvals, dvals, isDR, addsFutureValue)
	pitched := bufs.pitchedBuf
	attackers := bufs.attackersBuf
	defenders := bufs.defendersBuf

	var recurse func(i, pitchSum, defenseSum int)
	recurse = func(i, pitchSum, defenseSum int) {
		if i == totalN {
			prevented := defenseSum
			if prevented > incomingDamage {
				prevented = incomingDamage
			}
			// Group roles into played / pitched / defending buckets. Iterates the hand (size n),
			// then layers in the arsenal slot (index n) based on its assigned role. Arsenal-role
			// cards contribute nothing this turn whether they came from hand or the slot.
			var p, a, d []card.Card
			hasAnyDefender := hasReactions
			if !hasAnyDefender && arsenalCardIn != nil && rolesBuf[n] == Defend {
				hasAnyDefender = true
			}
			if hasAnyDefender {
				p, a, d = groupByRoleInto(hand, rolesBuf[:n], pitched[:0], attackers[:0], defenders[:0])
				if arsenalCardIn != nil {
					switch rolesBuf[n] {
					case Attack:
						a = append(a, arsenalCardIn)
					case Defend:
						d = append(d, arsenalCardIn)
					}
				}
			} else {
				p, a = groupPitchAttack(hand, rolesBuf[:n], pitched[:0], attackers[:0])
				if arsenalCardIn != nil && rolesBuf[n] == Attack {
					a = append(a, arsenalCardIn)
				}
			}
			// Arsenal-in is appended last to a / d above, so its index in the attackers slice is
			// len(a)-1 when present in the chain. -1 means no arsenal-in card in the attackers
			// (either no arsenal-in card at all, or it took a different role).
			arsenalInIdx := -1
			if arsenalCardIn != nil && rolesBuf[n] == Attack {
				arsenalInIdx = len(a) - 1
			}
			attackDealt, defenseDealt, leftoverRunechants, budget, swung, ok := bestAttackWithWeapons(hero, weapons, a, d, p, deck, bufs, runechantCarryover, incomingDamage, defenseSum, arsenalInIdx, priorAuraTriggers)
			if !ok {
				return
			}

			v := attackDealt + defenseDealt + prevented
			arsenalCard := findArsenalCard(rolesBuf, arsenalCardIn, n)
			// Hand cards never take Arsenal role during enumeration, so arsenalCard is only set
			// when arsenal-in stayed; post-hoc promotion potential is tracked via hasHeld.
			hasHeld := false
			futureValuePlayed := 0
			for j := 0; j < n; j++ {
				if rolesBuf[j] == Held {
					hasHeld = true
				}
				if rolesBuf[j] == Attack && addsFutureValue[j] {
					futureValuePlayed++
				}
			}
			if arsenalCardIn != nil && rolesBuf[n] == Attack && addsFutureValue[n] {
				futureValuePlayed++
			}
			willOccupy := arsenalCard != nil || hasHeld
			bestWillOccupy := best.ArsenalCard != nil || bestHasHeld
			if !beatsBest(v, leftoverRunechants, futureValuePlayed, willOccupy, best, bestFutureValuePlayed, bestWillOccupy) {
				return
			}
			best.Value = v
			bestSwung = swung
			bestBudget = budget
			best.LeftoverRunechants = leftoverRunechants
			best.ArsenalCard = arsenalCard
			bestHasHeld = hasHeld
			bestFutureValuePlayed = futureValuePlayed
			// Write the winning roles into BestLine. Cards and FromArsenal flags were populated
			// at construction; only Role varies. Contribution is cleared here and filled by
			// fillContributions below for the winning line.
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = rolesBuf[j]
				best.BestLine[j].Contribution = 0
			}
			return
		}
		isArsenalSlot := i == n && arsenalCardIn != nil
		// Hand cards can't take Arsenal role (post-hoc promotion handles that). Cap the range at
		// Held for hand slots to skip the roleAllowed-rejection work for Arsenal.
		maxRole := Held
		if isArsenalSlot {
			maxRole = Arsenal
		}
		for r := Role(0); r <= maxRole; r++ {
			if !roleAllowed(r, isArsenalSlot, isDR[i]) {
				continue
			}
			// With no damage coming in and no Defense Reactions in the hand, a non-DR card's
			// Defend contribution is 0 — same as Held — and nothing scans the defender set,
			// so the two partitions produce the same Value / leftover / futureValuePlayed and
			// Held wins the arsenal-occupancy tiebreaker. Skip the dominated Defend branch.
			// DR-present hands keep Defend because DR Play effects scan defenders as a
			// graveyard seed (e.g. Weeping Battleground banishing an aura a non-DR blocker
			// put there).
			if r == Defend && incomingDamage == 0 && !isDR[i] && !hasReactions {
				continue
			}
			rolesBuf[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pvals[i], defenseSum)
			case Defend:
				recurse(i+1, pitchSum, defenseSum+dvals[i])
			case Attack, Held, Arsenal:
				recurse(i+1, pitchSum, defenseSum)
			}
		}
	}
	recurse(0, 0, 0)
	// Once per Best call, on the winning line only, attribute per-card contribution.
	if len(best.BestLine) > 0 {
		fillContributions(&best, hero, weapons, bestSwung, bestBudget, deck, bufs, incomingDamage, runechantCarryover, priorAuraTriggers)
	}
	// If the arsenal slot is empty after enumeration, promote one Held card into it. Held hand
	// cards and Held mid-turn-drawn cards are treated as one pool — neither source is preferred,
	// because both end the turn as a single card of equivalent future-turn value. The pick is
	// deterministic per-hand (hashed from sorted card IDs + drawn card IDs + arsenal-in ID) so
	// the memo stays consistent, but spreads across candidates to avoid a lowest-ID bias.
	if best.ArsenalCard == nil {
		promoteRandomHeldToArsenal(&best, hand, n, arsenalCardIn)
	}
	return best
}

// promoteRandomHeldToArsenal picks one Held card — a hand card in best.BestLine or a mid-turn-
// drawn card in best.Drawn — and flips its role to Arsenal. Both sources share a single
// candidate pool so the draw isn't preferred over hand Helds (nor the other way around). No-op
// when nothing is Held.
func promoteRandomHeldToArsenal(best *TurnSummary, hand []card.Card, n int, arsenalCardIn card.Card) {
	handHeldCount := countHeldInBestLine(best.BestLine, n)
	drawnHeldCount := countHeldInDrawn(best.Drawn)
	total := handHeldCount + drawnHeldCount
	if total == 0 {
		return
	}
	pick := int(arsenalPromotionHash(hand, best.Drawn, arsenalCardIn) % uint64(total))
	// Walk hand Helds first (in BestLine order), then drawn Helds (in draw order), mapping pick
	// to the matching slot.
	if pick < handHeldCount {
		promoteNthHeldInBestLine(best, n, pick)
		return
	}
	promoteNthHeldInDrawn(best, pick-handHeldCount)
}

// countHeldInBestLine returns how many of the first n BestLine entries are still Role=Held after
// enumeration. The first-n restriction excludes any arsenal-in entry (which lives at index n and
// is never Held).
func countHeldInBestLine(line []CardAssignment, n int) int {
	c := 0
	for i := 0; i < n; i++ {
		if line[i].Role == Held {
			c++
		}
	}
	return c
}

// countHeldInDrawn returns how many mid-turn-drawn cards are still Role=Held after the winning
// chain resolved.
func countHeldInDrawn(drawn []CardAssignment) int {
	c := 0
	for i := range drawn {
		if drawn[i].Role == Held {
			c++
		}
	}
	return c
}

// arsenalPromotionHash computes the deterministic bucket seed that picks which Held card fills
// an empty arsenal slot. Uses FNV-1a over the sorted hand IDs + drawn card IDs + arsenal-in ID —
// the only requirement is a uniform spread across bucket counts 1..total so no lowest-ID bias
// creeps in while the memo stays consistent per hand.
func arsenalPromotionHash(hand []card.Card, drawn []CardAssignment, arsenalCardIn card.Card) uint64 {
	const (
		fnvOffsetBasis uint64 = 1469598103934665603
		fnvPrime       uint64 = 1099511628211
	)
	h := fnvOffsetBasis
	for _, c := range hand {
		h ^= uint64(c.ID())
		h *= fnvPrime
	}
	for _, d := range drawn {
		h ^= uint64(d.Card.ID())
		h *= fnvPrime
	}
	if arsenalCardIn != nil {
		h ^= uint64(arsenalCardIn.ID())
		h *= fnvPrime
	}
	return h
}

// promoteNthHeldInBestLine flips the pick-th Held hand card (in BestLine order) to Arsenal, and
// records it on best.ArsenalCard. Caller guarantees pick < count of Held entries.
func promoteNthHeldInBestLine(best *TurnSummary, n, pick int) {
	idx := 0
	for i := 0; i < n; i++ {
		if best.BestLine[i].Role != Held {
			continue
		}
		if idx == pick {
			best.BestLine[i].Role = Arsenal
			best.ArsenalCard = best.BestLine[i].Card
			return
		}
		idx++
	}
}

// promoteNthHeldInDrawn flips the pick-th Held mid-turn-drawn card (in draw order) to Arsenal,
// and records it on best.ArsenalCard. Caller guarantees pick < count of Held drawn entries.
func promoteNthHeldInDrawn(best *TurnSummary, pick int) {
	idx := 0
	for i := range best.Drawn {
		if best.Drawn[i].Role != Held {
			continue
		}
		if idx == pick {
			best.Drawn[i].Role = Arsenal
			best.ArsenalCard = best.Drawn[i].Card
			return
		}
		idx++
	}
}

// groupByRoleInto appends hand cards into caller-provided pitched/attackers/defenders slices
// (passed pre-reset to length 0) to avoid per-partition heap allocation.
func groupByRoleInto(hand []card.Card, roles []Role, pitched, attackers, defenders []card.Card) ([]card.Card, []card.Card, []card.Card) {
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		case Defend:
			defenders = append(defenders, c)
		}
	}
	return pitched, attackers, defenders
}

// groupPitchAttack is the reaction-free leaf's grouping step: skips the defenders bucket (only
// needed for Defense-Reaction-Play dispatch, which this path doesn't run).
func groupPitchAttack(hand []card.Card, roles []Role, pitched, attackers []card.Card) ([]card.Card, []card.Card) {
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		}
	}
	return pitched, attackers
}

// findArsenalCard returns the arsenal-in card when it stays in the arsenal slot, nil otherwise.
// Hand cards never take Arsenal role during enumeration (post-hoc promotion handles that), so
// the only slot that can be Arsenal is the arsenal-in slot at index n.
func findArsenalCard(rolesBuf []Role, arsenalCardIn card.Card, n int) card.Card {
	if arsenalCardIn != nil && rolesBuf[n] == Arsenal {
		return arsenalCardIn
	}
	return nil
}

// beatsBest decides whether a candidate partition displaces the current best. Tiebreak
// order: Value → leftover runechants (future arcane) → more AddsFutureValue cards played
// (hidden later-turn payoff the current-turn Value misses) → arsenal slot ending occupied
// (saves a hand slot next refill; covers both arsenal-in-stayed and Held-for-promotion).
func beatsBest(v, leftoverRunechants, futureValuePlayed int, willOccupyArsenal bool, best TurnSummary, bestFutureValuePlayed int, bestWillOccupyArsenal bool) bool {
	if v > best.Value {
		return true
	}
	if v < best.Value {
		return false
	}
	if leftoverRunechants > best.LeftoverRunechants {
		return true
	}
	if leftoverRunechants < best.LeftoverRunechants {
		return false
	}
	if futureValuePlayed > bestFutureValuePlayed {
		return true
	}
	if futureValuePlayed < bestFutureValuePlayed {
		return false
	}
	return willOccupyArsenal && !bestWillOccupyArsenal
}

// roleAllowed decides whether the partition enumerator may assign role r to the current card.
// The arsenal-in slot may only take Arsenal (stay), Attack (any non-DR card — non-attack actions
// play fine from arsenal on your turn), or Defend (Defense Reactions only — plain-blocking from
// arsenal isn't legal). Hand cards take any role except Attack for Defense Reactions (DRs only
// fire on the opponent's turn); their role loop caps at Held, so the "which Held card gets
// arsenaled" choice happens post-hoc and doesn't bias toward low-ID slots.
func roleAllowed(r Role, isArsenalSlot, isDefenseReaction bool) bool {
	if isArsenalSlot {
		switch r {
		case Pitch, Held:
			return false
		case Attack:
			return !isDefenseReaction
		case Defend:
			return isDefenseReaction
		}
		return true // Arsenal is always allowed on the arsenal-in slot.
	}
	return !(r == Attack && isDefenseReaction)
}

// defenseReactionDamage runs Play() for every Defense Reaction in defenders and sums the damage
// they deal back to the attacker (e.g. a banish-an-aura-for-arcane rider). Played in isolation
// — no attack ordering; TurnState carries Pitched / Deck plus a per-DR fresh copy of the
// defenders list in Graveyard so effects that scan the graveyard see plain blocks and other
// defenders. Uncapped: this damage is dealt, not prevented.
//
// state is caller-provided (from attackBufs) and reset per call. gravBuf is the caller-owned
// scratch backing state.Graveyard; the returned slice is the (possibly grown) buffer for reuse.
func defenseReactionDamage(defenders, pitched, deck []card.Card, state *card.TurnState, gravBuf []card.Card, cs *card.CardState) (int, []card.Card) {
	total := 0
	for _, d := range defenders {
		if !d.Types().IsDefenseReaction() {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		*state = card.TurnState{Pitched: pitched, Deck: deck, Graveyard: gravBuf}
		*cs = card.CardState{Card: d}
		total += d.Play(state, cs)
	}
	return total, gravBuf
}

// chainBudget captures the winning phase-split's attack-chain resource state. Reusing it to seed
// the replay ctx in fillContributions ensures playSequenceWithMeta finds the exact permutation
// that won during partition enumeration — critical for per-card attribution since different
// permutations can deal different per-card damage.
type chainBudget struct {
	resource         int
	maxPitch         int
	hasAttackPitches bool
}

// bestAttackWithWeapons evaluates one partition leaf end-to-end: it enumerates every split of
// pitched cards between the attack phase (funding played cards + swung weapons) and the defense
// phase (funding Defense Reactions), every subset of weapons to swing, and — via bestSequence —
// every ordering of the attack chain. Returns the winning (attackDamage, defenseDamage,
// leftoverRunechants, chainBudget, swungWeaponNames) with ok=true, or ok=false when no split is
// legal (no pitching arrangement covers the chain and DR costs without violating the pitch-timing
// rule).
//
// Pitch-timing rule: every pitched card must pay for something on the stack. For each phase that
// has any pitch, the residual budget after paying all costs must be less than the max pitch in
// that phase — otherwise one pitch could have been Held. playSequenceWithMeta enforces the
// attack-phase check per permutation; this function applies the defense-phase check after
// computing the DR cost at the chain's final runechant count.
//
// phaseBudgets is one (pmask) configuration's split of pitched-resource totals across the
// attack and defense phases. Each side tracks both its running total and the largest single
// pitch assigned to it — the "largest pitch" feeds the pitch-timing waste check (if the
// residual budget after paying all costs is at least that value, one pitch could have been
// Held, and the partition is illegal).
type phaseBudgets struct {
	attackBudget, defendBudget       int
	maxAttackPitch, maxDefendPitch   int
	hasAttackPitches, hasDefendPitches bool
}

// splitPitchesAcrossPhases assigns each pitch to the attack or defense phase based on the
// bitmask and computes the per-phase resource summary. Bit i set → pitchedVals[i] funds
// defense; bit i clear → it funds attack. phaseCount==1 forces every pitch to the attack
// phase (no DRs present or no pitches to split) regardless of pmask.
func splitPitchesAcrossPhases(pitchedVals []int, pmask, phaseCount int) phaseBudgets {
	var p phaseBudgets
	for i, v := range pitchedVals {
		if phaseCount > 1 && pmask&(1<<i) != 0 {
			p.defendBudget += v
			if v > p.maxDefendPitch {
				p.maxDefendPitch = v
			}
			p.hasDefendPitches = true
		} else {
			p.attackBudget += v
			if v > p.maxAttackPitch {
				p.maxAttackPitch = v
			}
			p.hasAttackPitches = true
		}
	}
	return p
}

// containsDefenseReaction reports whether any card in cards is a Defense Reaction. The
// partition-leaf precompute uses this to decide whether the defense-phase pitch enumeration
// needs to split budgets at all (no DRs means every pitch funds the attack phase).
func containsDefenseReaction(cards []card.Card) bool {
	for _, c := range cards {
		if c.Types().IsDefenseReaction() {
			return true
		}
	}
	return false
}

// Phase masks: when no Defense Reactions are present (or no pitches exist), all pitches go to
// the attack phase, so we visit one configuration. Otherwise we enumerate 2^|pitched| splits.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, defenders, pitched, deck []card.Card, bufs *attackBufs, runechantCarryover, incomingDamage, blockTotal, arsenalInIdx int, priorAuraTriggers []card.AuraTrigger) (int, int, int, chainBudget, []string, bool) {
	ctx := &sequenceContext{
		hero:               hero,
		pitched:            pitched,
		deck:               deck,
		bufs:               bufs,
		runechantCarryover: runechantCarryover,
		incomingDamage:     incomingDamage,
		blockTotal:         blockTotal,
		arsenalInIdx:       arsenalInIdx,
		priorAuraTriggers:  priorAuraTriggers,
		// Borrow bufs' pre-sized winner scratch so the eval closure's append-winner step reuses
		// one backing array per Best call instead of allocating per sequenceContext.
		drawnWinner:        bufs.drawnWinnerScratch[:0],
		auraTriggersWinner: bufs.auraTriggersWinnerScratch[:0],
	}
	// Hoist leaf-constant TurnState fields out of the per-permutation reset in
	// playSequenceWithMeta.
	ctx.seedState()

	// Defense Reactions fire independently of ordering and attack chain (each sees a fresh
	// TurnState with only Pitched + Deck), so their Play-return damage is constant across phase /
	// weapon masks. Compute it once; reseed ctx state for the attack chain afterwards.
	hasDRs := containsDefenseReaction(defenders)
	var defenseDealt int
	if hasDRs {
		defenseDealt, bufs.defenseGravScratch = defenseReactionDamage(defenders, pitched, deck, bufs.state, bufs.defenseGravScratch, &bufs.drCardStateScratch)
		ctx.seedState()
	}

	pitchedVals := bufs.pitchedValsScratch[:0]
	for _, c := range pitched {
		pitchedVals = append(pitchedVals, c.Pitch())
	}

	// Phase splits only matter when there is actually a defense phase to fund (a DR exists) AND
	// there are pitches to split. Otherwise every pitch goes to the attack phase and we visit a
	// single configuration.
	phaseCount := 1
	if hasDRs && len(pitched) > 0 {
		phaseCount = 1 << len(pitched)
	}

	// Pre-screen precomputation: printed-cost sums let us reject doomed (pmask, wmask) pairs in
	// O(1) before spinning up bestSequence's N! permutation loop. attackersMinCost sums the
	// floor-cost of each attacker (non-discount: printed Cost; discount: 0), a safe under-estimate
	// of chain cost. attackersPrinted is the no-discount upper bound, used for the pitch-waste
	// upper bound check.
	attackersMinCost := 0
	attackersMaxCost := 0
	for _, a := range attackers {
		m := attackerMetaPtrFor(a)
		attackersMinCost += m.minCost
		attackersMaxCost += m.maxCost
	}

	copy(bufs.attackerBuf, attackers)

	bestDealt := 0
	bestLeftoverRunechants := runechantCarryover
	var bestSwung []string
	var bestBudget chainBudget
	foundFeasible := false

	for pmask := 0; pmask < phaseCount; pmask++ {
		phase := splitPitchesAcrossPhases(pitchedVals, pmask, phaseCount)

		ctx.resourceBudget = phase.attackBudget
		ctx.hasAttackPitches = phase.hasAttackPitches
		ctx.maxAttackPitch = phase.maxAttackPitch

		for wmask := 0; wmask < 1<<len(weapons); wmask++ {
			weaponCost := bufs.weaponCosts[wmask] // weapons are static-cost
			// Lower bound on total chain cost (sum of MinCost across attackers + weapons). If the
			// attack budget can't cover even this floor, no permutation is feasible. Mid-turn
			// draws can pitch on top of the committed hand pitch ("hopeful" partitions) but
			// can't reduce the base cost, so this MinCost prune is safe. No matching pitch-timing
			// pre-screen here: drawn cards play as chain extensions and consume the residual, so
			// playSequenceWithMeta enforces pitch-timing post-extension instead.
			if attackersMinCost+weaponCost > phase.attackBudget {
				continue
			}
			allAttackers := bufs.attackerBuf[:len(attackers)]
			for i, w := range weapons {
				if wmask&(1<<i) != 0 {
					allAttackers = append(allAttackers, w)
				}
			}
			dealt, leftoverRunechants, legal := ctx.bestSequence(allAttackers, nil, nil, nil)
			if !legal {
				continue
			}
			// Cost the DRs against the chain's final runechant count. DRs with variable cost
			// read state.Runechants inside their Cost; static DRs return a constant. Reuse
			// bufs.drScratch instead of allocating a fresh TurnState per mask iteration — the
			// interface call boxes the pointer, so a stack allocation would escape and heap-alloc
			// every loop.
			bufs.drScratch = card.TurnState{Runechants: leftoverRunechants}
			drCost := 0
			for _, d := range defenders {
				if !d.Types().IsDefenseReaction() {
					continue
				}
				drCost += d.Cost(&bufs.drScratch)
			}
			if drCost > phase.defendBudget {
				continue
			}
			if phase.hasDefendPitches && phase.defendBudget-drCost >= phase.maxDefendPitch {
				continue
			}
			if !foundFeasible || dealt > bestDealt ||
				(dealt == bestDealt && leftoverRunechants > bestLeftoverRunechants) {
				bestDealt = dealt
				bestLeftoverRunechants = leftoverRunechants
				bestSwung = bufs.weaponNames[wmask]
				bestBudget = chainBudget{resource: phase.attackBudget, maxPitch: phase.maxAttackPitch, hasAttackPitches: phase.hasAttackPitches}
				foundFeasible = true
			}
		}
	}

	if !foundFeasible {
		return 0, 0, 0, chainBudget{}, nil, false
	}
	return bestDealt, defenseDealt, bestLeftoverRunechants, bestBudget, bestSwung, true
}

// sequenceContext carries the stable per-partition-leaf environment: hero (for OnCardPlayed
// triggers), pitched / deck refs for Card.Play, shared scratch buffers, and the numeric budgets
// that persist across permutation and mask iterations. Built once per leaf so the hot inner
// calls (playSequence, bestSequence) shrink to their varying inputs and tracking outputs.
//
// resourceBudget / hasAttackPitches / maxAttackPitch are rewritten by bestAttackWithWeapons on
// each phase-mask iteration: they fund the attack chain and let playSequenceWithMeta reject
// permutations whose final residual breaks FaB's pitch-timing rule (excess >= max pitch means
// one pitch could have been Held instead).
type sequenceContext struct {
	hero               hero.Hero
	pitched, deck      []card.Card
	bufs               *attackBufs
	resourceBudget     int
	runechantCarryover int
	incomingDamage     int
	blockTotal         int
	hasAttackPitches   bool
	maxAttackPitch     int
	// arsenalInIdx is the index in the attackers slice (the slice passed to bestSequence) of
	// the card that came from the arsenal slot at start of turn, or -1 when no arsenal-in card
	// is in the chain. Lets bestSequence flag the matching pcBuf entry's FromArsenal as the
	// permutation moves it around.
	arsenalInIdx int
	// priorAuraTriggers are the AuraTriggers carried in from the previous turn (e.g. an
	// AttackAction trigger from a Malefic Incantation played a turn ago). Each permutation
	// seeds state.AuraTriggers with a fresh copy of this slice so mid-chain firing can
	// decrement Count / set FiredThisTurn without leaking those mutations across permutations.
	priorAuraTriggers []card.AuraTrigger
	// drawnWinner snapshots the winning permutation's drawn cards so fillContributions can
	// surface them on summary.Drawn. Populated from state.Drawn during the winner path because
	// Heap's algorithm keeps iterating after the winner is chosen and state.Drawn reflects the
	// last permutation's draws. Every entry in drawnWinner is assigned Role=Held by the
	// caller (post-Best, promoteRandomHeldToArsenal may flip one to Arsenal).
	drawnWinner []card.Card
	// auraTriggersWinner snapshots the winning permutation's final state.AuraTriggers so the
	// deck loop can carry them into next turn. Includes both inherited triggers from
	// priorAuraTriggers (with mutated Count / FiredThisTurn) and ones added by Play.
	auraTriggersWinner []card.AuraTrigger
}

// seedState writes the TurnState fields that are constant across a partition's permutations
// fireAttackActionTriggers walks state.AuraTriggers after an attack action card resolves and
// invokes every TriggerAttackAction entry whose OncePerTurn gate is open. Each fire
// decrements the trigger's Count; when Count hits zero the aura drops out of the list and
// Self lands in the graveyard so downstream same-turn effects see the destroy. Returns the
// summed Damage from all fires, folded into chain damage by the caller.
//
// Slice mutation: a survivors prefix is built in place over the existing slice; entries
// kept after firing are written back at increasing indices, exhausted ones are skipped.
func fireAttackActionTriggers(state *card.TurnState) int {
	total := 0
	triggers := state.AuraTriggers
	dst := triggers[:0]
	for i := range triggers {
		t := triggers[i]
		if t.Type != card.TriggerAttackAction || (t.OncePerTurn && t.FiredThisTurn) {
			dst = append(dst, t)
			continue
		}
		total += t.Handler(state)
		t.FiredThisTurn = true
		t.Count--
		if t.Count <= 0 {
			state.AddToGraveyard(t.Self)
			continue
		}
		dst = append(dst, t)
	}
	state.AuraTriggers = dst
	return total
}

// fireEphemeralAttackTriggers walks state.EphemeralAttackTriggers after an attack action
// card resolves and invokes every entry whose Matches predicate accepts the attacker. Each
// fire consumes the trigger (fire-once semantics) and routes its damage to the source's
// perCardOut slot via SourceIndex — Mauvrion Skies's "if hits" Runechants, for instance,
// surface on Mauvrion's BestLine entry rather than the attacker's. Non-matching entries stay
// in the slice for a later attack action; anything still in the list at end of chain
// fizzles silently (no graveyard bookkeeping — the source was already graveyarded when its
// own Play resolved).
//
// Slice mutation parallels fireAttackActionTriggers: a survivors prefix is built in place
// over the existing slice, with fired entries skipped.
func fireEphemeralAttackTriggers(state *card.TurnState, target *card.CardState, perCardOut []float64) int {
	total := 0
	triggers := state.EphemeralAttackTriggers
	dst := triggers[:0]
	for i := range triggers {
		t := triggers[i]
		if t.Matches != nil && !t.Matches(target) {
			dst = append(dst, t)
			continue
		}
		dmg := t.Handler(state, target)
		total += dmg
		if perCardOut != nil && t.SourceIndex >= 0 && t.SourceIndex < len(perCardOut) {
			perCardOut[t.SourceIndex] += float64(dmg)
		}
	}
	state.EphemeralAttackTriggers = dst
	return total
}

// seedState writes the leaf-constant TurnState fields (pitched / deck refs, incoming damage,
// block total) so the per-permutation reset in playSequenceWithMeta can skip them.
func (ctx *sequenceContext) seedState() {
	s := ctx.bufs.state
	s.Pitched = ctx.pitched
	s.Deck = ctx.deck
	s.IncomingDamage = ctx.incomingDamage
	s.BlockTotal = ctx.blockTotal
}

// bestSequence tries every ordering of attackers and returns the max total damage plus the
// runechant count at the end of the winning permutation. Between each card's Play() and its
// append to CardsPlayed, the hero's OnCardPlayed hook fires so triggered abilities contribute.
// legal=true when at least one ordering is playable; false when every permutation is rejected
// by playSequenceWithMeta's resource / go-again / pitch-waste checks.
//
// Uses Heap's algorithm (iterative) — no closure/callback alloc, no recursive call per perm.
//
// When winnerOrderOut is non-nil (len >= len(attackers)) the winning permutation is copied into
// it. perCardOut / perCardTriggerOut (same size rule) receive the winning line's per-card Play
// damage and hero-trigger damage. fillContributions uses these; the partition-loop caller
// passes nil for all three so the permutation search stays allocation-free.
func (ctx *sequenceContext) bestSequence(attackers, winnerOrderOut []card.Card, perCardOut, perCardTriggerOut []float64) (int, int, bool) {
	n := len(attackers)
	if n == 0 {
		// No attackers means no chain costs are deducted — the attack phase spends zero from the
		// budget. If any attack-phase pitches exist they over-pay (residual == budget >= maxPitch
		// since the budget is the sum of those pitches); pitch-timing fails.
		if ctx.hasAttackPitches && ctx.resourceBudget >= ctx.maxAttackPitch {
			return 0, 0, false
		}
		return 0, ctx.runechantCarryover, true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		pcBuf[idx] = card.CardState{Card: c, FromArsenal: idx == ctx.arsenalInIdx}
	}

	// Scratch buffers are playSequence's per-card outputs, overwritten every permutation. On a
	// new winner we copy them into the caller's perCardOut / perCardTriggerOut. Only populated
	// when the caller asked to track.
	var scratch, triggerScratch []float64
	if perCardOut != nil {
		scratch = ctx.bufs.perCardScratch[:n]
	}
	if perCardTriggerOut != nil {
		if cap(ctx.bufs.perCardTriggerScratch) < n {
			ctx.bufs.perCardTriggerScratch = make([]float64, n)
		}
		triggerScratch = ctx.bufs.perCardTriggerScratch[:n]
	}

	best := 0
	bestLeftoverRunechants := ctx.runechantCarryover
	foundLegal := false
	ctx.drawnWinner = ctx.drawnWinner[:0]
	ctx.auraTriggersWinner = ctx.auraTriggersWinner[:0]
	eval := func() {
		dmg, leftoverRunechants, _, legal := ctx.playSequenceWithMeta(n, scratch, triggerScratch)
		if !legal {
			return
		}
		if !foundLegal || dmg > best ||
			(dmg == best && leftoverRunechants > bestLeftoverRunechants) {
			best = dmg
			bestLeftoverRunechants = leftoverRunechants
			foundLegal = true
			ctx.drawnWinner = append(ctx.drawnWinner[:0], ctx.bufs.state.Drawn...)
			ctx.auraTriggersWinner = append(ctx.auraTriggersWinner[:0], ctx.bufs.state.AuraTriggers...)
			if winnerOrderOut != nil {
				for i := 0; i < n; i++ {
					winnerOrderOut[i] = pcBuf[i].Card
				}
			}
			if perCardOut != nil {
				copy(perCardOut[:n], scratch)
			}
			if perCardTriggerOut != nil {
				copy(perCardTriggerOut[:n], triggerScratch)
			}
		}
	}
	eval()
	// Heap's algorithm, iterative: c[] counts how many times each stack frame has iterated.
	// pcBuf and permMeta swap together so playSequenceWithMeta sees meta aligned with the
	// current permutation. FromArsenal rides inside pcBuf (one byte), so it permutes for free;
	// no separate permFromArsenal slice to maintain.
	var c [8]int
	i := 0
	for i < n {
		if c[i] < i {
			if i&1 == 0 {
				pcBuf[0], pcBuf[i] = pcBuf[i], pcBuf[0]
				permMeta[0], permMeta[i] = permMeta[i], permMeta[0]
			} else {
				pcBuf[c[i]], pcBuf[i] = pcBuf[i], pcBuf[c[i]]
				permMeta[c[i]], permMeta[i] = permMeta[i], permMeta[c[i]]
			}
			eval()
			c[i]++
			i = 0
		} else {
			c[i] = 0
			i++
		}
	}
	return best, bestLeftoverRunechants, foundLegal
}

// playSequence plays `order` as a sequence of cards, reusing ctx.bufs' pooled buffers.
// Buffers are mutated in place; the caller must not read them concurrently.
//
// When perCardOut is non-nil (len >= n) each entry is the card's Play return for that
// position; perCardTriggerOut (same size rule) receives the hero's OnCardPlayed return.
// The hot partition-loop callers pass nil for both.
//
// Runechant flow:
//   - state.Runechants starts at ctx.runechantCarryover.
//   - Play / OnCardPlayed calling CreateRunechants increments the count AND returns n damage
//     — tokens are credited exactly once, at creation.
//   - After each Attack / Weapon card resolves, all current tokens fire and are destroyed;
//     state.Runechants is zeroed but damage is NOT re-added (tokens were credited at
//     creation).
//   - At end of the sequence, state.Runechants is the leftover count carrying into next turn.
//
// Resource flow: ctx.resourceBudget is the starting pool; each card deducts
// attackerMeta.costAt(state). Negative remaining budget returns legal=false.
//
// Populates permMeta from order and then calls playSequenceWithMeta. The hot path
// (bestSequence) builds meta once and calls playSequenceWithMeta directly to amortise
// interface dispatch across the N! permutations.
func (ctx *sequenceContext) playSequence(order []card.Card, perCardOut, perCardTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	ctx.seedState()
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		pcBuf[i] = card.CardState{Card: c, FromArsenal: i == ctx.arsenalInIdx}
	}
	return ctx.playSequenceWithMeta(n, perCardOut, perCardTriggerOut)
}

// playSequenceWithMeta runs the permutation currently held in ctx.bufs.pcBuf[:n] with
// aligned permMeta[:n]. CardState (Card + FromArsenal) persists across permutations, so only
// GrantedGoAgain needs a per-permutation reset — Play is the only thing that flips it.
func (ctx *sequenceContext) playSequenceWithMeta(n int, perCardOut, perCardTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		pcBuf[i].GrantedGoAgain = false
		if perCardOut != nil {
			perCardOut[i] = 0
		}
		if perCardTriggerOut != nil {
			perCardTriggerOut[i] = 0
		}
	}
	played := ptrBuf[:n]
	state := ctx.bufs.state
	// Reset only the fields playSequence mutates. Pitched and Deck are stable across permutations
	// in a partition leaf; CardsRemaining and Self are rewritten per-card in the loop. A
	// Per-permutation reset. Only touch fields the cards mutate; leaf-stable fields (Pitched,
	// Deck, IncomingDamage, BlockTotal) come from seedState. A full-struct replace here
	// memcpies big slice headers on every permutation and profiles dramatically slower.
	state.CardsPlayed = ctx.bufs.cardsPlayedBuf[:0]
	state.Runechants = ctx.runechantCarryover
	state.ArcaneDamageDealt = false
	state.AuraCreated = false
	state.Overpower = false
	state.NonAttackActionPlayed = false
	// Deck and Drawn reset per permutation: DrawOne mutates them, so a prior permutation's
	// consumption would poison the next.
	state.Deck = ctx.deck
	state.Drawn = nil
	// Graveyard and Banish reset per permutation: cards append themselves to Graveyard as
	// they resolve, and graveyard-banish effects shift cards into Banish. Reusing the scratch
	// backing array keeps the reset allocation-free.
	state.Graveyard = ctx.bufs.attackGravScratch[:0]
	state.Banish = nil
	// AuraTriggers reset per permutation: seeded with a copy of priorAuraTriggers so
	// mid-chain attack-action triggers can fire without their Count / FiredThisTurn
	// mutations leaking across permutations. Cards adding triggers via AddAuraTrigger extend
	// the same scratch slice.
	state.AuraTriggers = append(ctx.bufs.auraTriggersScratch[:0], ctx.priorAuraTriggers...)
	// EphemeralAttackTriggers reset per permutation as empty — fire-once triggers never
	// carry across turns, so there's nothing to seed from prior state.
	state.EphemeralAttackTriggers = ctx.bufs.ephemeralTriggersScratch[:0]
	resources := ctx.resourceBudget
	for i, pc := range played {
		m := meta[i]
		cost := m.costAt(state)
		resources -= cost
		if resources < 0 {
			return 0, 0, 0, false
		}

		state.CardsRemaining = played[i+1:]

		// If this card is an attack or weapon and any Runechant is live, those tokens fire on
		// its damage step. Set ArcaneDamageDealt now — before Play and OnCardPlayed — so Play
		// effects that read "if you've dealt arcane damage this turn" see the flag for same-hand
		// triggers. Cards that deal arcane damage via their Play text flip the flag themselves.
		isAttackOrWeapon := m.isAttackOrWeapon
		if isAttackOrWeapon && state.Runechants > 0 {
			state.ArcaneDamageDealt = true
		}

		// Hero ability fires BEFORE the card's own Play so "aura created this turn" checks
		// inside the card's Play see the runechant (or other aura) the hero just made.
		// Viserai's "another non-attack action" gate still excludes the current card because
		// NonAttackActionPlayed isn't flipped until the end of the iteration.
		triggerDmg := ctx.hero.OnCardPlayed(pc.Card, state)
		ephemeralsBefore := len(state.EphemeralAttackTriggers)
		playDmg := pc.Card.Play(state, pc)
		// Stamp SourceIndex on any EphemeralAttackTriggers the card registered during Play
		// so fireEphemeralAttackTriggers can route their damage back to this card's
		// perCardOut slot.
		for k := ephemeralsBefore; k < len(state.EphemeralAttackTriggers); k++ {
			state.EphemeralAttackTriggers[k].SourceIndex = i
		}
		auraTriggerDmg := 0
		ephemeralDmg := 0
		if m.isAttackAction {
			auraTriggerDmg = fireAttackActionTriggers(state)
			// Fire ephemeral triggers AFTER hero and aura triggers so the handler sees the
			// fully-resolved attacker state (Dominate grants, hero-created auras, fresh
			// Runechants from aura triggers). Damage is routed back to each trigger's
			// source via SourceIndex, so perCardOut is updated in place inside the helper.
			ephemeralDmg = fireEphemeralAttackTriggers(state, pc, perCardOut)
		}
		damage += playDmg + triggerDmg + auraTriggerDmg + ephemeralDmg
		if perCardOut != nil {
			perCardOut[i] = float64(playDmg)
		}
		if perCardTriggerOut != nil {
			// Fold AuraTrigger damage into the hero-trigger slot for per-card attribution —
			// the damage is driven by the attack action card that resolved, not by the aura
			// (which has no BestLine entry by mid-chain). Ephemeral trigger damage is NOT
			// folded in here — fireEphemeralAttackTriggers already credited it to the
			// trigger's source via perCardOut[SourceIndex], which is the semantically
			// correct attribution (the effect belongs to the card that registered the
			// trigger, not the attacker that happened to consume it).
			perCardTriggerOut[i] = float64(triggerDmg + auraTriggerDmg)
		}
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)
		if m.types.IsNonAttackAction() {
			state.NonAttackActionPlayed = true
		}
		// Weapons and persistent card types (Auras, Items) stay in their zone when they
		// resolve; any destroy event that should send them to the graveyard is a separate
		// trigger. Everything else — Actions, Attack Reactions, Defense Reactions, Blocks,
		// Instants — heads to the graveyard immediately.
		if !m.types.PersistsInPlay() {
			state.Graveyard = append(state.Graveyard, pc.Card)
		}

		// Attacks and weapon swings consume all runechants in play. Damage isn't re-added: each
		// token was credited +1 at creation time, so this is pure state cleanup.
		if isAttackOrWeapon {
			state.Runechants = 0
		}

		if i < n-1 && !(m.baseGoAgain || pc.GrantedGoAgain) {
			return 0, 0, 0, false
		}
	}

	// Mid-turn-drawn cards always carry to the next hand as Held or compete for the empty
	// arsenal slot; they never pitch or extend the chain. If they could, the solver's best line
	// would depend on Deck[0], which the player commits before the draw reveals.

	// Pitch-timing rule: every Pitch-role card must have paid for something on the stack. If the
	// chain's leftover budget is at least the max attack-phase pitch, one pitch could have been
	// Held instead — this permutation violates FaB's rules.
	if ctx.hasAttackPitches && resources >= ctx.maxAttackPitch {
		return 0, 0, 0, false
	}
	return damage, state.Runechants, resources, true
}

// fillDefenseContributions writes Contribution on each Defend-role entry. The block-prevention
// share is proportional to the card's effective defense out of sumDef, capped by incomingDamage
// so over-blocking doesn't inflate attribution past what actually stopped. Effective defense is
// Defense() plus the arsenal bonus when FromArsenal is set on a card.ArsenalDefenseBonus
// implementer. Defense Reactions add their own Play return on top, evaluated against a fresh
// TurnState seeded with the turn's pitched pool and remaining deck so card effects see the same
// context the solver scored them in.
func fillDefenseContributions(line []CardAssignment, pitched []card.Card, deck []card.Card, bufs *attackBufs, sumDef, incomingDamage int) {
	prevented := sumDef
	if prevented > incomingDamage {
		prevented = incomingDamage
	}
	// Collect defenders first so each DR's Play sees the full set in state.Graveyard — mirroring
	// the seeding defenseReactionDamage does during partition enumeration.
	defenders := bufs.defendersBuf[:0]
	for i := range line {
		if line[i].Role == Defend {
			defenders = append(defenders, line[i].Card)
		}
	}
	for i := range line {
		if line[i].Role != Defend {
			continue
		}
		c := line[i].Card
		def := c.Defense()
		if line[i].FromArsenal {
			if ab, ok := c.(card.ArsenalDefenseBonus); ok {
				def += ab.ArsenalDefenseBonus()
			}
		}
		if sumDef > 0 {
			line[i].Contribution = float64(def) * float64(prevented) / float64(sumDef)
		}
		if c.Types().IsDefenseReaction() {
			bufs.defenseGravScratch = append(bufs.defenseGravScratch[:0], defenders...)
			*bufs.state = card.TurnState{Pitched: pitched, Deck: deck, Graveyard: bufs.defenseGravScratch}
			bufs.drCardStateScratch = card.CardState{Card: c, FromArsenal: line[i].FromArsenal}
			line[i].Contribution += float64(c.Play(bufs.state, &bufs.drCardStateScratch))
		}
	}
}

// fillContributions populates each BestLine entry's Contribution from the winning line:
//   - Pitch:  Card.Pitch() as resource value.
//   - Defend: proportional share of Prevented plus own Play return if a Defense Reaction.
//   - Attack: per-card damage from the winning attack-chain replay.
//   - Held / Arsenal: zero (contributed nothing this turn).
//
// Called once per Best call after the partition loop picks the winner. All transient slices
// (pitched/attackers/chain/winnerOrder/perCard/used) borrow attackBufs slots so nothing
// allocates here.
func fillContributions(summary *TurnSummary, hero hero.Hero, weapons []weapon.Weapon, swungNames []string, budget chainBudget, deck []card.Card, bufs *attackBufs, incomingDamage, runechantCarryover int, priorAuraTriggers []card.AuraTrigger) {
	line := summary.BestLine

	// Reconstruct pitched and attackers from the winning line. The arsenal-in entry
	// (FromArsenal=true, last slot) participates in attackers / defenders identically to hand
	// entries when its role is Attack / Defend.
	pitched := bufs.pitchedBuf[:0]
	attackers := bufs.attackersBuf[:0]
	arsenalInIdx := -1
	var sumDef int
	for _, a := range line {
		switch a.Role {
		case Pitch:
			pitched = append(pitched, a.Card)
		case Attack:
			if a.FromArsenal {
				arsenalInIdx = len(attackers)
			}
			attackers = append(attackers, a.Card)
		case Defend:
			def := a.Card.Defense()
			if a.FromArsenal {
				if ab, ok := a.Card.(card.ArsenalDefenseBonus); ok {
					def += ab.ArsenalDefenseBonus()
				}
			}
			sumDef += def
		}
	}

	// Pitch contributions.
	for i := range line {
		if line[i].Role == Pitch {
			line[i].Contribution = float64(line[i].Card.Pitch())
		}
	}

	fillDefenseContributions(line, pitched, deck, bufs, sumDef, incomingDamage)

	chain := buildAttackChain(bufs.attackerBuf[:0], attackers, weapons, swungNames)
	if len(chain) > 0 {
		// Re-seed ctx with the winning phase split's chain-resource state so bestSequence
		// reproduces the exact permutation that won during enumeration; per-card damage
		// depends on order.
		ctx := &sequenceContext{
			hero:               hero,
			pitched:            pitched,
			deck:               deck,
			bufs:               bufs,
			resourceBudget:     budget.resource,
			runechantCarryover: runechantCarryover,
			incomingDamage:     incomingDamage,
			blockTotal:         sumDef,
			hasAttackPitches:   budget.hasAttackPitches,
			maxAttackPitch:     budget.maxPitch,
			arsenalInIdx:       arsenalInIdx,
			priorAuraTriggers:  priorAuraTriggers,
			// Same borrow as bestAttackWithWeapons above — fillContributions clones the
			// winners into summary before returning, so sharing bufs-backed storage is safe.
			drawnWinner:        bufs.drawnWinnerScratch[:0],
			auraTriggersWinner: bufs.auraTriggersWinnerScratch[:0],
		}
		ctx.seedState()
		fillAttackChainContributions(summary, chain, ctx)
		// Copy the winning permutation's drawn cards out as CardAssignments. Read from
		// ctx.drawnWinner (bestSequence's winner snapshot), not bufs.state.Drawn — state.Drawn
		// reflects whichever permutation Heap's algorithm iterated last, which can diverge
		// from the winner when different permutations trigger different draws. Drawn cards
		// start Held with zero contribution; promoteRandomHeldToArsenal may flip one to
		// Arsenal post-enumeration.
		if drawn := ctx.drawnWinner; len(drawn) > 0 {
			summary.Drawn = make([]CardAssignment, len(drawn))
			for i, c := range drawn {
				summary.Drawn[i] = CardAssignment{Card: c, Role: Held}
			}
		}
		// Fresh slice so the returned TurnSummary doesn't alias ctx's buf-backed scratch — the
		// memo keeps TurnSummaries around and a later Best call would otherwise overwrite the
		// cached entry's triggers on its next permutation sweep.
		if n := len(ctx.auraTriggersWinner); n > 0 {
			summary.AuraTriggers = make([]card.AuraTrigger, n)
			copy(summary.AuraTriggers, ctx.auraTriggersWinner)
		}
	}
}

// buildAttackChain appends attackers first, then the weapons named in swungNames in that order,
// so the sequence search sees the same chain composition the partition loop priced. Uses the
// passed-in slice's backing array (typically bufs.attackerBuf) to stay allocation-free.
func buildAttackChain(dst []card.Card, attackers []card.Card, weapons []weapon.Weapon, swungNames []string) []card.Card {
	dst = append(dst, attackers...)
	for _, name := range swungNames {
		for _, w := range weapons {
			if w.Name() == name {
				dst = append(dst, w)
				break
			}
		}
	}
	return dst
}

// fillAttackChainContributions re-runs the sequence search with tracking enabled to recover
// the winning permutation, snapshots it into summary.AttackChain (fresh slice to avoid
// aliasing the buf-backed winnerOrder), and maps each position's damage back to BestLine's
// Attack-role entries. Weapons have no BestLine entry; their damage is already in
// summary.Value. Duplicate printings disambiguate by scan order. Contribution bundles Play
// return + hero-trigger so per-card stats reflect total this-turn impact.
func fillAttackChainContributions(summary *TurnSummary, chain []card.Card, ctx *sequenceContext) {
	line := summary.BestLine
	total := len(line)
	bufs := ctx.bufs
	winnerOrder := bufs.fillContribWinnerOrder[:len(chain)]
	perCardDmg := bufs.fillContribPerCard[:len(chain)]
	if cap(bufs.fillContribTriggerDmg) < len(chain) {
		bufs.fillContribTriggerDmg = make([]float64, len(chain))
	}
	perCardTrigger := bufs.fillContribTriggerDmg[:len(chain)]
	ctx.bestSequence(chain, winnerOrder, perCardDmg, perCardTrigger)
	summary.AttackChain = make([]AttackChainEntry, len(winnerOrder))
	for i := range winnerOrder {
		summary.AttackChain[i] = AttackChainEntry{
			Card:          winnerOrder[i],
			Damage:        perCardDmg[i],
			TriggerDamage: perCardTrigger[i],
		}
	}
	if cap(bufs.fillContribUsed) < total {
		bufs.fillContribUsed = make([]bool, total)
	}
	used := bufs.fillContribUsed[:total]
	for i := range used {
		used[i] = false
	}
	for k, c := range winnerOrder {
		if _, isWeapon := c.(weapon.Weapon); isWeapon {
			continue
		}
		for i := range line {
			if used[i] || line[i].Role != Attack || line[i].Card.ID() != c.ID() {
				continue
			}
			line[i].Contribution = perCardDmg[k] + perCardTrigger[k]
			used[i] = true
			break
		}
	}
}
