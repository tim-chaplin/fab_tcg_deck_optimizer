package sim

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Per-turn shared context threaded through Card.Play. Cards mutate state directly — moving
// cards between hand / deck / graveyard / Banish, registering triggers, creating runechants
// — and the sim copies the winning permutation's final state into next-turn state. There's
// no diff-signal indirection: a card that wants to draw goes through AppendHand and pops
// via PopDeckTop, full stop.
//
// Persistent fields (hand, deck, Arsenal, graveyard, Banish, Runechants, Auras)
// carry across turns when the sim adopts the winner's snapshot. Transient fields
// (CardsPlayed, Pitched, IncomingDamage, etc.) are seeded by the sim per chain-step and
// reset at the turn boundary.
//
// hand, deck, and graveyard are unexported so card subpackages (internal/cards) can only
// reach them through the accessor methods below. Every accessor clears cacheable so the
// hand-eval cache can store results only when IsCacheable() is true at chain end. The
// framework (this package) accesses the slices directly to seed and snapshot state
// without poisoning the bit; card code in a different package can't see the unexported
// field name, which is the language-level enforcement that the cacheable signal is sound.

// LogEntryKind classifies a LogEntry. Triggers come in two flavours because they fire on
// opposite sides of their parent in the FaB stack — the format layer needs to know which
// side a given entry sits on so two cards with the same display name in the same chain
// don't steal each other's triggers during grouping.
type LogEntryKind int8

const (
	// LogEntryChainStep is the sim's "<Card>: <VERB>" line. Stands alone in the printout
	// and acts as the parent that triggers attach to.
	LogEntryChainStep LogEntryKind = iota
	// LogEntryPreTrigger is a trigger that fires before its parent chain entry resolves
	// (hero / aura attack-action triggers). The format layer attaches it to the next
	// chain entry whose name matches Source.
	LogEntryPreTrigger
	// LogEntryPostTrigger is a trigger that fires after its parent chain entry resolves
	// (OnHit riders, AR buffs). The format layer attaches it to the previous chain entry
	// whose name matches Source.
	LogEntryPostTrigger
)

// LogEntry is one chain-event entry in TurnState.Log. Text is the freeform display string
// the producer authored ("Viserai created a runechant", "Consuming Volition [R]: ATTACK")
// — the format layer renders it verbatim, no further opinions on phrasing. Kind tags the
// entry as a chain step or as a pre/post trigger so the grouping algorithm can attribute
// triggers correctly even when sibling chain entries share a name. Source names the card
// whose play caused the trigger and is matched against chain-entry names to pick the
// parent; only meaningful for triggers. N is the damage-equivalent credited to s.Value
// when the entry was added.
type LogEntry struct {
	Text   string
	Source string
	Kind   LogEntryKind
	N      int
}

// TurnState is the shared turn-level context passed to Card.Play alongside the per-card
// CardState wrapper.
type TurnState struct {
	// hand is the cards currently in hand. Starts as the dealt hand minus pitched /
	// attacker / defender cards (routed by the partition); cards drawn or tutored mid-
	// chain land here too. Whatever's left at end of chain becomes next turn's Held set.
	// Card subpackages reach hand only through Hand() / AppendHand / PopHandAt — see the
	// package-level docstring for the framework-vs-card boundary rationale.
	hand []Card
	// deck is the deck top-to-bottom. Unexported so card subpackages can only reach it via
	// the public Deck() / PopDeckTop / PrependToDeck / TutorFromDeck accessors, each of
	// which clears cacheable. Framework code in this package reads / writes deck directly
	// (resetStateForPermutation seed, snapshotCarry copy, applyTurnResult adoption) so the
	// non-card-driven path doesn't poison the cacheable bit.
	deck []Card
	// Arsenal is the arsenal slot's contents at this point in the chain — the arsenal-in
	// card at start of turn, nil after it plays / defends, refilled post-chain by the
	// arsenal-promotion step. Cards that read "from arsenal" use CardState.FromArsenal,
	// not this field.
	Arsenal Card
	// graveyard is cards that have entered the graveyard this turn — every card played or
	// blocked lands here after resolving. Pitched cards do not (they recycle to deck
	// bottom). Unexported for the same reason as deck: cards reach it only via Graveyard()
	// / BanishFromGraveyard / AddToGraveyard, all of which clear cacheable. Framework
	// code in this package writes graveyard directly (the dispatcher's "card resolved →
	// non-persistent goes to graveyard" rule, DestroyAura's aura-card append) so the
	// non-card-driven append doesn't poison cacheable.
	graveyard []Card
	// Banish holds cards moved into the banished zone this turn (e.g. an aura-banish-for-
	// arcane rider).
	Banish []Card
	// ActionPoints is the chain runner's running AP pool. Seeded to 1 per permutation,
	// decremented before each paying chain step, incremented after a Go-again card resolves.
	// Free chain steps (Instants, Attack Reactions) cost 0; Action cards and weapon swings
	// cost 1. A paying card resolving with no AP available makes the chain illegal.
	ActionPoints int
	// ArcaneDamageDealt sticks true once any source of arcane damage fires this turn:
	// a Runechant token consuming itself on an attack / weapon swing, or a card whose Play
	// deals arcane directly. Effects that read "if you've dealt arcane damage this turn"
	// consult this flag rather than Runechants. Reset at turn boundary.
	ArcaneDamageDealt bool
	// OpponentMarked tracks Mark on the opposing hero. Set by a Mark effect, cleared
	// when an attack action card or weapon swing resolves — modelled as the opponent's
	// first incoming physical damage stripping the mark. Arcane damage (Runechants,
	// arcane riders) doesn't clear the mark. Read by "if the defending hero is marked"
	// riders. Our own hero is never modelled as marked, so "if you are marked, can't
	// play" gates on opposing cards drop out. Carries across turns via CarryState so
	// next turn's chain seed sees a Mark applied late in the previous turn.
	OpponentMarked bool
	// Auras is the list of auras currently in play. Cards add entries during Play via
	// AddAura; the sim fires matching entries on each TriggerType condition and drops
	// entries whose handler called s.DestroyAura. Carries across turns.
	Auras []Aura
	// Triggers is the list of one-shot deferred handlers keyed to a TriggerType. Cards
	// add entries during Play via AddTrigger; the sim fires matching entries on each
	// TriggerType condition and removes them after firing. Reset per permutation; not
	// snapshotted into CarryState (TriggerEndOfTurn fires before the carry snapshot).
	Triggers []Trigger
	// Items is the list of items currently in play. Cards add entries via the
	// per-token-type Create helper (CreateGold); the chain runner enqueues each item's
	// Ability as a playable activated ability each turn. Carries across turns.
	Items []Item
	// CardsDrawn counts mid-chain card draws this turn — incremented by DrawOne and
	// any future tutor-into-hand helper. The partition tiebreaker prefers chains that
	// drew more cards: a draw is one extra play available next turn, comparable in
	// tempo to a future runechant. Reset per permutation; snapshotted into CarryState
	// so the partition recurse compares end-of-chain draw counts across leaves.
	CardsDrawn int
	// currentAuraIdx is the index in Auras of the handler currently running. DestroyAura
	// uses it as a fast-path hint to skip the linear scan in the common "handler destroys
	// its own aura" case; a Self comparison guards against stale hints.
	currentAuraIdx int
	// currentAuraDestroyed is set by DestroyAura when it splices the entry at
	// currentAuraIdx. The aura-loop reads it to decide whether to advance the cursor
	// (false) or stay (true — the next entry shifted into position i).
	currentAuraDestroyed bool

	// --- Transient: reset by the sim per turn / chain step ---

	// Value is the running damage-equivalent total for this chain — damage dealt + damage
	// prevented + every aura-token / hero-trigger credit. The dispatcher records the chain
	// step's Play+BonusAttack contribution via AddLogEntry; trigger handlers (hero, aura,
	// OnHit) credit themselves the same way. The solver compares permutations on this
	// field. Reset by the sim per permutation.
	Value int
	// turnLog is the per-event chain trace — one entry per chain step / hero / aura /
	// OnHit / weapon swing. Cards reach it through the Log / LogRider / LogPreTrigger /
	// LogPostTrigger family below; external readers (tests, format layer entry points)
	// use the LogEntries accessor. The format layer uses each entry's Kind plus Source to
	// cluster triggers under the right parent. Reset per permutation. Lowercase so callers
	// can't bypass the skipLog gate by appending directly.
	turnLog []LogEntry
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order.
	// Populated by the sim after each Play returns so later cards this turn see what was
	// played before them.
	CardsPlayed []Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant
	// tokens). Effects that check "if you've played or created an aura this turn" should
	// OR this with CardsPlayed containing an Aura-typed card.
	AuraCreated bool
	// CardsRemaining is the cards that will be played after the current one in chain order.
	// Populated by the sim before each Play so an effect can peek forward ("next X attack")
	// or grant keywords to a later card by flipping flags on its CardState entry.
	CardsRemaining []*CardState
	// Pitched is the cards pitched this turn for resources. Populated by the sim before any
	// Play. Effects that check "if an attack card was pitched" scan this list.
	Pitched []Card
	// Overpower is set when an attack with the Overpower keyword is being played. Not yet
	// consumed by the sim — blocked damage should eventually be forwarded to the hero when
	// Overpower is true.
	Overpower bool
	// NonAttackActionPlayed is set true once any non-attack action card has been appended to
	// CardsPlayed this turn. Maintained by the chain runner so hero triggers that ask "was a
	// non-attack action played earlier?" can answer in O(1).
	NonAttackActionPlayed bool
	// IncomingDamage is the opponent damage this turn — seeded by the sim from the value
	// passed to Best, and decremented by ApplyAndLogEffectiveDefense as defenders block.
	// Cards reading "did we block all incoming?" against the static partition aggregate use
	// BlockTotal instead.
	IncomingDamage int
	// ArcaneIncomingDamage is the opponent's arcane damage this turn, seeded from the
	// -arcane-incoming flag. Cards whose riders gate on "if you've been dealt arcane damage
	// this turn" read this directly; not decremented during the chain (defending arcane isn't
	// modelled).
	ArcaneIncomingDamage int
	// BlockTotal is the sum of Defense() across every Defend-role card in the current
	// partition. Uncapped: if the partition over-blocks, BlockTotal is the full sum, not
	// clamped to IncomingDamage.
	BlockTotal int
	// Defenders is the partition's defender slice (DRs + plain blocks together) populated
	// by the chain runner before invoking each defender's Play / Block hook. Cards reading
	// "how many other plain blockers / DRs are defending alongside" iterate this list — the
	// defender-side counterpart to CardsRemaining for attackers.
	Defenders []Card
	// attackReactionTarget is the buff target for the currently-resolving Attack Reaction.
	// Set by the chain runner around AR.Play; ARs read it via AttackReactionTarget().
	attackReactionTarget *CardState
	// TriggeringCard is the card whose play caused the active aura attack-action trigger
	// to fire. The sim sets it before each Aura handler runs and clears it after;
	// the handler reads it to attribute its log line back to the triggering card. Hero
	// and OnHit handlers receive the triggering card as a direct arg already and don't
	// need this field. Nil during direct chain-step resolution and start-of-turn fires.
	TriggeringCard Card
	// skipLog short-circuits Log appends and the per-entry text formatting for chains the
	// caller doesn't intend to display. The Log* helpers below own the gate end-to-end —
	// Value is still credited (so the sim's running damage tally stays correct) but the
	// helpers skip every fmt.Sprintf, DisplayName lookup, and slice append underneath. The
	// eval loop runs every turn silent; only the rare new-deck-best turn replays with
	// skipLog=false to materialise its Log for the printout. Cards must never read this
	// field — the helpers handle it. The lowercase name is the language-level enforcement.
	skipLog bool

	// cacheable is true while the chain hasn't read or mutated deck / graveyard through any
	// public accessor (Deck / Graveyard / PopDeckTop / PrependToDeck / TutorFromDeck /
	// BanishFromGraveyard / AddToGraveyard) or framework helper built on them (DrawOne,
	// Clash). Set to false by the accessor on first card-driven access; never restored
	// within a permutation. Constructors (NewTurnState, resetStateForPermutation,
	// defendersDamage's per-DR seed) explicitly set cacheable=true so a fresh state starts
	// cacheable; a zero-value `var s TurnState{}` defaults to false (uncacheable) — the more
	// conservative default that surfaces missing initialization rather than hiding it.
	cacheable bool
}

// IsCacheable reports whether the chain so far has not depended on hidden state — i.e. no
// card in this chain has read or mutated deck / graveyard via an accessor. A future
// hand-eval cache stores results only when this is true at chain end.
func (s *TurnState) IsCacheable() bool { return s.cacheable }

// AttackReactionTarget returns the buff target for the currently-resolving AR, or nil when
// no AR is resolving.
func (s *TurnState) AttackReactionTarget() *CardState { return s.attackReactionTarget }

// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. ARs use this to
// fold their +{p} buff into the buffed attack's display delta. No-op when skipLog elided
// log entries.
func (s *TurnState) AmendLastChainStepN(n int) {
	if s.skipLog || n == 0 {
		return
	}
	for i := len(s.turnLog) - 1; i >= 0; i-- {
		if s.turnLog[i].Kind == LogEntryChainStep {
			s.turnLog[i].N += n
			return
		}
	}
}

// Deck returns the live deck top-to-bottom and flips IsCacheable to false. Cards must not
// mutate the returned slice; use PopDeckTop / PrependToDeck / TutorFromDeck for mutations.
// Read-only callers that only inspect the slice still flip — the cache key can't depend on
// what the deck-order read produced.
func (s *TurnState) Deck() []Card {
	s.cacheable = false
	return s.deck
}

// Hand returns the live hand slice and flips IsCacheable to false. Cards must not mutate
// the returned slice; use AppendHand / PopHandAt for mutations. Read-only callers that
// only inspect length / contents still flip — the cache key can't depend on what the
// hand-content read produced.
func (s *TurnState) Hand() []Card {
	s.cacheable = false
	return s.hand
}

// AppendHand appends c to the hand, flipping IsCacheable to false. Cards that draw / tutor
// into hand mid-chain use this so the cache invalidation is automatic. Framework code in
// this package writes s.hand directly (resetStateForPermutation seeds, applyChainStep
// removes the playing card) and doesn't go through here — those mutations aren't
// card-driven and shouldn't poison the cacheable bit.
func (s *TurnState) AppendHand(c Card) {
	s.cacheable = false
	s.hand = append(s.hand, c)
}

// PopHandAt removes and returns the card at index i, flipping IsCacheable to false. Cards
// that pop hand cards (alt-cost effects, Moon Wish's "return a hand card to top of deck")
// use this so the cache invalidation is automatic. Panics on out-of-range i — callers
// should len-check via Hand() first.
func (s *TurnState) PopHandAt(i int) Card {
	s.cacheable = false
	c := s.hand[i]
	s.hand = append(s.hand[:i], s.hand[i+1:]...)
	return c
}

// SetHandForTesting replaces the hand with the supplied cards. Test-only — production
// hand seeding goes through resetStateForPermutation (framework) or AppendHand (cards).
// Doesn't flip cacheable: tests that assert on hand state typically don't care about the
// cache, and a NewTurnState-then-SetHandForTesting flow shouldn't poison the bit before
// the test has run a single Play.
func (s *TurnState) SetHandForTesting(cards []Card) {
	s.hand = cards
}

// SetCurrentAuraIdxForTesting sets currentAuraIdx so AuraFor / DestroyAura resolve
// inside a manually-fired aura handler. Test-only — production fire loops maintain
// currentAuraIdx themselves.
func (s *TurnState) SetCurrentAuraIdxForTesting(i int) {
	s.currentAuraIdx = i
}

// Graveyard returns the live graveyard slice and flips IsCacheable to false. Cards must
// not mutate the returned slice; use BanishFromGraveyard for mutations or AddToGraveyard
// for the deterministic append-only path.
func (s *TurnState) Graveyard() []Card {
	s.cacheable = false
	return s.graveyard
}

// PopDeckTop removes the top card of the deck and returns it. Returns (nil, false) when
// the deck is empty. Flips IsCacheable to false.
func (s *TurnState) PopDeckTop() (Card, bool) {
	s.cacheable = false
	if len(s.deck) == 0 {
		return nil, false
	}
	top := s.deck[0]
	s.deck = s.deck[1:]
	return top, true
}

// PeekDeck returns the top card of the deck without removing it. Returns (nil, false) on
// an empty deck. Flips IsCacheable to false — observing the deck top makes the chain's
// output depend on hidden shuffle order, same as PopDeckTop.
func (s *TurnState) PeekDeck() (Card, bool) {
	s.cacheable = false
	if len(s.deck) == 0 {
		return nil, false
	}
	return s.deck[0], true
}

// PrependToDeck inserts c at the top of the deck. Flips IsCacheable to false. Allocates a
// fresh backing slice so subsequent mid-chain mutations don't poison sibling permutations
// that share the per-leaf deck reference.
func (s *TurnState) PrependToDeck(c Card) {
	s.cacheable = false
	newDeck := make([]Card, 0, len(s.deck)+1)
	newDeck = append(newDeck, c)
	newDeck = append(newDeck, s.deck...)
	s.deck = newDeck
}

// RecycleToDeckBottom appends self.Card to the bottom of the deck and flips
// self.SkipGraveyard so the chain dispatcher skips the usual non-persistent → graveyard
// append. Models the FaB clause "put this on the bottom of its owner's deck"
// (Relentless Pursuit). Flips IsCacheable to false. Allocates a fresh deck backing slice
// so the per-leaf deck reference shared across permutations stays untouched.
func (s *TurnState) RecycleToDeckBottom(self *CardState) {
	s.cacheable = false
	newDeck := make([]Card, 0, len(s.deck)+1)
	newDeck = append(newDeck, s.deck...)
	newDeck = append(newDeck, self.Card)
	s.deck = newDeck
	self.SkipGraveyard = true
}

// Opt resolves the FaB "Opt N" keyword: pops up to n cards from the top of the deck and
// hands them to the current hero's Opt heuristic. The handler returns a (top, bottom)
// split; the top list is placed back on top of the deck (in returned order) and the
// bottom list appends to the bottom of the deck (in returned order). n is clamped to the
// current deck length, so an Opt N call against a shorter deck reshapes whatever's there
// without error. Always flips IsCacheable to false — Opt always reads the deck, so the
// chain becomes uncacheable regardless of whether n is positive or whether any cards
// were available.
//
// Emits a log entry "Opted X, put Y on top, put Z on bottom" naming the revealed cards
// and the chosen split when the handler ran (no-op paths skip the log to keep the trace
// quiet on degenerate cases).
//
// Panics if the handler's combined output isn't exactly the input multiset. The contract
// is that Opt only re-orders cards; adding, dropping, or substituting any card is a bug.
//
// Allocates a fresh deck backing slice so the per-leaf deck reference shared across
// permutations stays untouched (same convention as PrependToDeck / TutorFromDeck).
func (s *TurnState) Opt(n int) {
	s.cacheable = false
	if n <= 0 || len(s.deck) == 0 {
		return
	}
	if n > len(s.deck) {
		n = len(s.deck)
	}
	// Copy off the popped slice so the handler can't mutate s.deck through aliasing.
	cards := append([]Card(nil), s.deck[:n]...)
	rest := s.deck[n:]

	top, bottom := CurrentHero.Opt(cards)
	panicIfOptViolatesMultiset(cards, top, bottom)

	newDeck := make([]Card, 0, len(top)+len(rest)+len(bottom))
	newDeck = append(newDeck, top...)
	newDeck = append(newDeck, rest...)
	newDeck = append(newDeck, bottom...)
	s.deck = newDeck

	if OptDebug {
		fmt.Printf("Opt(%d): cards=%s -> top=%s bottom=%s\n",
			n, formatCardList(cards), formatCardList(top), formatCardList(bottom))
	}
	if s.skipLog {
		return
	}
	s.Logf(0, "Opted %s, put %s on top, put %s on bottom",
		formatCardList(cards), formatCardList(top), formatCardList(bottom))
}

// formatCardList renders cs as "[name1, name2, ...]" using DisplayName for each entry, or
// "[]" when cs is empty. Used by the Opt log entry; the caller gates the call on s.skipLog
// so this only runs when the chain materialises its log.
func formatCardList(cs []Card) string {
	if len(cs) == 0 {
		return "[]"
	}
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = DisplayName(c)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// panicIfOptViolatesMultiset enforces TurnState.Opt's contract that the hero handler's
// combined (top, bottom) output is exactly the input multiset — a permutation of the
// input cards, no additions or removals. Panics with a descriptive message naming the
// failure mode (size mismatch, foreign card, or dropped card). Cards are zero-sized
// structs in production and small POD structs in tests; both flavours are usable as
// map keys for the multiset count.
func panicIfOptViolatesMultiset(in, top, bottom []Card) {
	if len(top)+len(bottom) != len(in) {
		panic(fmt.Sprintf("Opt: handler returned %d+%d cards, want %d (input multiset)",
			len(top), len(bottom), len(in)))
	}
	counts := make(map[Card]int, len(in))
	for _, c := range in {
		counts[c]++
	}
	check := func(out []Card, label string) {
		for _, c := range out {
			counts[c]--
			if counts[c] < 0 {
				panic(fmt.Sprintf("Opt: %s list returned card %s not in input",
					label, DisplayName(c)))
			}
		}
	}
	check(top, "top")
	check(bottom, "bottom")
	for c, n := range counts {
		if n != 0 {
			panic(fmt.Sprintf("Opt: handler dropped %d copy of %s from input", n, DisplayName(c)))
		}
	}
}

// TutorFromDeck removes and returns the highest-scoring card per score. Returns (nil,
// false) when no card scores > 0 (or the deck is empty). Flips IsCacheable to false.
// Allocates a fresh backing slice so the per-leaf deck reference shared across
// permutations stays untouched.
func (s *TurnState) TutorFromDeck(score func(Card) int) (Card, bool) {
	s.cacheable = false
	bestIdx := -1
	bestScore := 0
	for i, c := range s.deck {
		sc := score(c)
		if sc > bestScore {
			bestScore = sc
			bestIdx = i
		}
	}
	if bestIdx < 0 {
		return nil, false
	}
	found := s.deck[bestIdx]
	out := make([]Card, 0, len(s.deck)-1)
	out = append(out, s.deck[:bestIdx]...)
	out = append(out, s.deck[bestIdx+1:]...)
	s.deck = out
	return found, true
}

// BanishFromGraveyard removes the first graveyard card matching pred, appends it to
// s.Banish, and returns it. Returns (nil, false) when no card matches. Flips IsCacheable
// to false. Reads graveyard contents from a previous turn (or from this turn's plain
// blocks the partition put there) — the chain output depends on hidden prior-turn state.
func (s *TurnState) BanishFromGraveyard(pred func(Card) bool) (Card, bool) {
	s.cacheable = false
	for i, c := range s.graveyard {
		if !pred(c) {
			continue
		}
		s.Banish = append(s.Banish, c)
		s.graveyard = append(s.graveyard[:i], s.graveyard[i+1:]...)
		return c, true
	}
	return nil, false
}

// RecycleFromGraveyardToTop removes the first graveyard card matching pred, prepends it to
// the deck, and returns it. Returns (nil, false) when no card matches. Flips IsCacheable to
// false (reading graveyard + mutating deck). The deck mutation IS the model — the recycled
// card resurfaces in next turn's deal naturally, so callers don't credit Value here. Pair
// with a rider log line so the trace records the recycle.
func (s *TurnState) RecycleFromGraveyardToTop(pred func(Card) bool) (Card, bool) {
	return s.recycleFromGraveyard(pred, true)
}

// RecycleFromGraveyardToBottom is RecycleFromGraveyardToTop's bottom-of-deck variant: same
// IsCacheable flip, same no-Value-credit contract; the recycled card lands at the bottom of
// the deck instead of the top.
func (s *TurnState) RecycleFromGraveyardToBottom(pred func(Card) bool) (Card, bool) {
	return s.recycleFromGraveyard(pred, false)
}

func (s *TurnState) recycleFromGraveyard(pred func(Card) bool, toTop bool) (Card, bool) {
	s.cacheable = false
	for i, c := range s.graveyard {
		if !pred(c) {
			continue
		}
		s.graveyard = append(s.graveyard[:i], s.graveyard[i+1:]...)
		newDeck := make([]Card, 0, len(s.deck)+1)
		if toTop {
			newDeck = append(newDeck, c)
			newDeck = append(newDeck, s.deck...)
		} else {
			newDeck = append(newDeck, s.deck...)
			newDeck = append(newDeck, c)
		}
		s.deck = newDeck
		return c, true
	}
	return nil, false
}

// NewTurnState constructs a *TurnState with the given deck and graveyard seed. Test /
// utility constructor: the unexported deck / graveyard fields aren't reachable via a
// composite literal from outside this package, so callers in card subpackages and other
// non-sim packages route through this constructor (or set the slices via the accessor
// methods after construction). The returned state has IsCacheable()==true; cacheable
// has to be set explicitly because the field's zero value is false (see the field doc).
func NewTurnState(deck, graveyard []Card) *TurnState {
	return &TurnState{deck: deck, graveyard: graveyard, cacheable: true, currentAuraIdx: -1}
}

// AddValue credits n to s.Value, clamped at 0. Pair with a Log helper when you also want a
// log line; call alone for silent value (an aura that pays out without surfacing in the
// printout). Negative n is a no-op (FaB damage / prevention can't drive the running total
// negative). The convention is to put AddValue on its own line, separate from any Log call,
// so a line beginning with Log( has no side effects.
func (s *TurnState) AddValue(n int) {
	if n > 0 {
		s.Value += n
	}
}

// LogEntries returns the per-event chain trace accumulated by the Log family. External
// readers (tests, format layer) use this; package-internal code reads the underlying field.
func (s *TurnState) LogEntries() []LogEntry { return s.turnLog }

// log is the single skipLog gate. When not running silent, appends a LogEntry of the given
// kind, source, and pre-built text. Every public Log helper funnels through here or its
// variadic sibling logf, so the gate lives in exactly one place and cards never check
// skipLog themselves. log does NOT credit s.Value — pair the Log helper with AddValue when
// you also want to record damage.
func (s *TurnState) log(kind LogEntryKind, source, text string, n int) {
	if s.skipLog {
		return
	}
	if n < 0 {
		n = 0
	}
	s.turnLog = append(s.turnLog, LogEntry{
		Kind:   kind,
		Text:   text,
		Source: source,
		N:      n,
	})
}

// logf is the format variant: same gate as log, but fmt.Sprintf only runs on the !skipLog
// branch. Callers pay variadic-arg boxing at the call site regardless, so prefer the
// non-format Log helpers when text is constant or pre-built.
func (s *TurnState) logf(kind LogEntryKind, source string, n int, format string, args ...any) {
	if s.skipLog {
		return
	}
	if n < 0 {
		n = 0
	}
	s.turnLog = append(s.turnLog, LogEntry{
		Kind:   kind,
		Text:   fmt.Sprintf(format, args...),
		Source: source,
		N:      n,
	})
}

// Log appends the canonical "<DisplayName>: <VERB>[ from arsenal]" main-line chain-step
// entry for self, with display suffix "(+n)". Use for both attacks (n = effective attack)
// and non-attack chain steps (n = 0). Pair with AddValue or self.DealEffectiveAttack /
// self.DealEffectiveDefense on a separate line so the Log call itself has no side effects.
// ChainStepText is deferred into the !skipLog branch.
func (s *TurnState) Log(self *CardState, n int) {
	if s.skipLog {
		return
	}
	if n < 0 {
		n = 0
	}
	s.turnLog = append(s.turnLog, LogEntry{
		Kind: LogEntryChainStep,
		Text: ChainStepText(self),
		N:    n,
	})
}

// Logf appends a free-form main-line chain-step entry with formatted text. Use when no
// CardState applies (Opt's "Opted X, put Y on top, put Z on bottom").
func (s *TurnState) Logf(n int, format string, args ...any) {
	s.logf(LogEntryChainStep, "", n, format, args...)
}

// LogRider appends an indented post-trigger sub-line under self's chain entry. Use for
// "Created a runechant", "Gained 3 health (graveyard trigger)", "On-hit discarded a card",
// etc. Pair with AddValue on a separate preceding line when n > 0.
func (s *TurnState) LogRider(self *CardState, n int, text string) {
	if s.skipLog {
		return
	}
	s.log(LogEntryPostTrigger, DisplayName(self.Card), text, n)
}

// LogRiderf is the format variant of LogRider — defers fmt.Sprintf and DisplayName into the
// !skipLog branch.
func (s *TurnState) LogRiderf(self *CardState, n int, format string, args ...any) {
	if s.skipLog {
		return
	}
	s.logf(LogEntryPostTrigger, DisplayName(self.Card), n, format, args...)
}

// LogPreTrigger appends an indented pre-trigger sub-line attributed to source — a hero or
// aura-attack-action trigger that fires before its parent chain entry. The format layer
// attaches this entry to the next chain entry whose name matches source.
func (s *TurnState) LogPreTrigger(source, text string, n int) {
	s.log(LogEntryPreTrigger, source, text, n)
}

// LogPreTriggerf is the format variant of LogPreTrigger.
func (s *TurnState) LogPreTriggerf(source string, n int, format string, args ...any) {
	s.logf(LogEntryPreTrigger, source, n, format, args...)
}

// LogPostTrigger appends an indented post-trigger sub-line attributed to source — for
// rider lines whose host differs from self (e.g. an OnHit attached to a target card). Use
// LogRider when self's CardState is the host.
func (s *TurnState) LogPostTrigger(source, text string, n int) {
	s.log(LogEntryPostTrigger, source, text, n)
}

// LogPostTriggerf is the format variant of LogPostTrigger.
func (s *TurnState) LogPostTriggerf(source string, n int, format string, args ...any) {
	s.logf(LogEntryPostTrigger, source, n, format, args...)
}

// DrawOne models a mid-turn draw: pop the top of the deck and append it to Hand. No-op on
// an empty deck. Every draw-rider card routes through this helper. Inherits the flip via
// PopDeckTop — a card that draws makes the chain's output depend on hidden shuffle order,
// same as a card that reads the deck top. Bumps CardsDrawn so the partition tiebreaker
// can prefer chains with more draws.
func (s *TurnState) DrawOne() {
	c, ok := s.PopDeckTop()
	if !ok {
		return
	}
	s.hand = append(s.hand, c)
	s.CardsDrawn++
}

// HasPlayedType reports whether any card played this turn has the given type in its Types() set.
func (s *TurnState) HasPlayedType(t card.CardType) bool {
	for _, c := range s.CardsPlayed {
		if c.Types().Has(t) {
			return true
		}
	}
	return false
}

// HasPlayedOrCreatedAura reports whether an aura was played or created this turn — the
// condition behind "if you've played or created an aura this turn" riders. The aura need
// not still be on the battlefield; the flag is sticky for the rest of the turn.
func (s *TurnState) HasPlayedOrCreatedAura() bool {
	return s.AuraCreated || s.HasPlayedType(card.TypeAura)
}

// Clash models a clash (rule 8.5.45): we and the opponent reveal the top card of our decks
// and the higher {p} wins. We model from our side only — our deck's top is read via
// s.Deck(); the opponent's top is approximated as 5-power. On a win (our top ≥ 6), win
// fires; on a loss (our top ≤ 4), lose fires; ties (top == 5) and an empty deck fire
// neither. Either callback may be nil. Reading the deck top through Deck() flips
// IsCacheable to false — a clash result depends on hidden shuffle order.
func (s *TurnState) Clash(win, lose func()) {
	deck := s.Deck()
	if len(deck) == 0 {
		return
	}
	top := deck[0].Attack()
	switch {
	case top >= 6:
		if win != nil {
			win()
		}
	case top <= 4:
		if lose != nil {
			lose()
		}
	}
}

// CreateRunechants creates n Runechant tokens and credits +n damage at creation time.
// Tokens are stored as a single Aura entry — bump an existing entry's Count or add a
// new one. Sets AuraCreated so same-turn "aura created this turn" effects see it.
// Tokens that never fire (end-of-sim leftovers) are slightly over-credited — accepted.
func (s *TurnState) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	s.AuraCreated = true
	s.AddValue(n)
	for i := range s.Auras {
		if s.Auras[i].Self.TokenType == TokenTypeRunechant {
			s.Auras[i].Count += n
			return
		}
	}
	s.AddAura(NewRunechantAura(n))
}

// CreatePonder creates n Ponder tokens, bumping the existing aura entry's Count or
// adding a new one, and flips AuraCreated. No Value credit — see ponderAuraHandler.
func (s *TurnState) CreatePonder(n int) {
	if n <= 0 {
		return
	}
	s.AuraCreated = true
	for i := range s.Auras {
		if s.Auras[i].Self.TokenType == TokenTypePonder {
			s.Auras[i].Count += n
			return
		}
	}
	s.AddAura(NewPonderAura(n))
}

// Runechants returns the current Runechant token count, or zero when none are in play.
func (s *TurnState) Runechants() int { return tokenCountIn(s.Auras, TokenTypeRunechant) }

// Ponders returns the current Ponder token count, or zero when none are in play.
func (s *TurnState) Ponders() int { return tokenCountIn(s.Auras, TokenTypePonder) }

// CreateGold creates n Gold tokens, bumping the existing item entry's Count or adding a
// new one. No Value credit — Gold only pays out when the player spends one via
// GoldTokenAbility (which decrements Count and draws a card).
func (s *TurnState) CreateGold(n int) {
	if n <= 0 {
		return
	}
	for i := range s.Items {
		if s.Items[i].Self.TokenType == TokenTypeGold {
			s.Items[i].Count += n
			return
		}
	}
	s.Items = append(s.Items, NewGoldItem(n))
}

// Gold returns the current Gold token count, or zero when none are in play.
func (s *TurnState) Gold() int { return itemCountIn(s.Items, TokenTypeGold) }

// CreateSilver creates n Silver tokens, same shape as CreateGold.
func (s *TurnState) CreateSilver(n int) {
	if n <= 0 {
		return
	}
	for i := range s.Items {
		if s.Items[i].Self.TokenType == TokenTypeSilver {
			s.Items[i].Count += n
			return
		}
	}
	s.Items = append(s.Items, NewSilverItem(n))
}

// Silver returns the current Silver token count, or zero when none are in play.
func (s *TurnState) Silver() int { return itemCountIn(s.Items, TokenTypeSilver) }

// CreateCopper creates n Copper tokens, same shape as CreateGold.
func (s *TurnState) CreateCopper(n int) {
	if n <= 0 {
		return
	}
	for i := range s.Items {
		if s.Items[i].Self.TokenType == TokenTypeCopper {
			s.Items[i].Count += n
			return
		}
	}
	s.Items = append(s.Items, NewCopperItem(n))
}

// Copper returns the current Copper token count, or zero when none are in play.
func (s *TurnState) Copper() int { return itemCountIn(s.Items, TokenTypeCopper) }

// ConsumeItem decrements the matching item's Count by n and removes the entry when
// Count reaches zero. Token items don't head to the graveyard on destroy. No-op when
// no item matches t.
func (s *TurnState) ConsumeItem(t TokenType, n int) {
	for i := range s.Items {
		if s.Items[i].Self.TokenType != t {
			continue
		}
		s.Items[i].Count -= n
		if s.Items[i].Count <= 0 {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
		}
		return
	}
}

// DealArcaneDamage credits n arcane damage to Value, writes a "Dealt n arcane damage" rider
// line under self's chain entry, and flips ArcaneDamageDealt when LikelyDamageHits(n, false)
// approves so same-turn triggers reading "if you've dealt arcane damage this turn" fire.
// Routes through dealtArcaneText[n] so the hot path avoids per-call fmt.Sprintf and
// variadic-int boxing.
func (s *TurnState) DealArcaneDamage(self *CardState, n int) {
	s.AddValue(n)
	if LikelyDamageHits(n, false) {
		s.ArcaneDamageDealt = true
	}
	if n >= 0 && n < len(dealtArcaneText) {
		s.LogRider(self, n, dealtArcaneText[n])
		return
	}
	s.LogRiderf(self, n, "Dealt %d arcane damage", n)
}

// dealtArcaneText is the pre-built rider-line cache indexed by arcane-damage count, keeping
// DealArcaneDamage alloc-free on the hot path. Extend if a new card prints higher arcane.
var dealtArcaneText = [...]string{
	0: "Dealt 0 arcane damage",
	1: "Dealt 1 arcane damage",
	2: "Dealt 2 arcane damage",
	3: "Dealt 3 arcane damage",
	4: "Dealt 4 arcane damage",
}

// AddToGraveyard appends c to graveyard so later-resolving cards see it. Used by cards
// running a mini-dispatcher inline (Moon Wish's go-again Sun Kiss play) that need to route
// the inline-played card through the same "non-persistent → graveyard" rule the framework
// dispatcher applies. Flips IsCacheable to false so the convention "every public accessor
// that touches deck / graveyard flips cacheable" stays universal — framework code that
// graveyards a played card writes s.graveyard directly (same package, no flip) and only
// card-driven calls reach this method, so the flip is sound and conservative.
func (s *TurnState) AddToGraveyard(c Card) {
	s.cacheable = false
	s.graveyard = append(s.graveyard, c)
}

// AddAura is the Play-side combo every Action - Aura card reaches for: flip
// AuraCreated so same-turn "if you've played or created an aura" riders see the entry, and
// append t to s.Auras so the sim fires it on its matching TriggerType condition.
func (s *TurnState) AddAura(t Aura) {
	s.AuraCreated = true
	s.Auras = append(s.Auras, t)
}

// AddTrigger appends t to s.Triggers. The sim fires t once on its matching TriggerType
// condition then removes it.
func (s *TurnState) AddTrigger(t Trigger) {
	s.Triggers = append(s.Triggers, t)
}

// DestroyAura splices t out of s.Auras and, when addToGraveyard, appends t.Self.Card to
// s.graveyard. Token-style auras (Self.Card == nil) skip the graveyard append
// unconditionally. During the aura fire walk the destroy uses currentAuraIdx directly
// (the fire loop maintains it through any mid-handler s.Auras realloc); off-fire callers
// fall back to a pointer-equality scan.
//
// Direct graveyard append (no cacheable flip): destruction is deterministic from the
// triggering event the sim already accounts for, not from hidden state.
func (s *TurnState) DestroyAura(t *Trigger, addToGraveyard bool) {
	a := s.AuraFor(t)
	if a == nil {
		return
	}
	if addToGraveyard && a.Self.Card != nil {
		s.graveyard = append(s.graveyard, a.Self.Card)
	}
	if i := s.currentAuraIdx; i >= 0 && i < len(s.Auras) && &s.Auras[i] == a {
		s.Auras = append(s.Auras[:i], s.Auras[i+1:]...)
		s.currentAuraDestroyed = true
		return
	}
	for i := range s.Auras {
		if &s.Auras[i] == a {
			s.Auras = append(s.Auras[:i], s.Auras[i+1:]...)
			return
		}
	}
}

// AuraFor returns the Aura whose embedded Trigger is t — used by aura handlers to
// recover Self / Count given the *Trigger they were called with. During an aura fire
// walk currentAuraIdx points at the firing aura (the fire loop maintains it through
// any mid-handler s.Auras realloc); the index path always wins inside a handler. The
// pointer-equality scan is the off-fire-path fallback for callers without
// currentAuraIdx context. Returns nil for standalone Triggers (no parent Aura).
func (s *TurnState) AuraFor(t *Trigger) *Aura {
	if i := s.currentAuraIdx; i >= 0 && i < len(s.Auras) {
		return &s.Auras[i]
	}
	for i := range s.Auras {
		if &s.Auras[i].Trigger == t {
			return &s.Auras[i]
		}
	}
	return nil
}
