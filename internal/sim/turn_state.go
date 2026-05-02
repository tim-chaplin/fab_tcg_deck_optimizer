package sim

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Per-turn shared context threaded through Card.Play. Cards mutate state directly — moving
// cards between Hand / Deck / Graveyard / Banish, registering triggers, creating runechants
// — and the sim copies the winning permutation's final state into next-turn state. There's
// no diff-signal indirection: a card that wants to draw appends to s.Hand and pops via
// PopDeckTop, full stop.
//
// Persistent fields (Hand, deck, Arsenal, graveyard, Banish, Runechants, AuraTriggers)
// carry across turns when the sim adopts the winner's snapshot. Transient fields
// (CardsPlayed, Pitched, IncomingDamage, etc.) are seeded by the sim per chain-step and
// reset at the turn boundary.
//
// deck and graveyard are unexported so card subpackages (internal/cards) can only reach them
// through the accessor methods below. Every accessor clears cacheable so a future hand-eval
// cache can key on the inputs and store only when IsCacheable() is true at chain end. The
// framework (this package) accesses the slices directly to seed and snapshot state without
// poisoning the bit; card code in a different package can't see the unexported field name,
// which is the language-level enforcement that the cacheable signal is sound.

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

// NextAttackActionHitTrigger is a one-shot rider queued by a card whose printed text reads
// "the next time an attack action card you control hits this turn, do X". The chain runner
// drains the queue inside finalizeActiveAttack on the first attack action that lands
// (IsAttackAction + LikelyToHit); every pending trigger fires together on that hit (the
// "next time" event resolves all listeners simultaneously) and the queue empties.
type NextAttackActionHitTrigger struct {
	Fire   func(s *TurnState, target *CardState, t *NextAttackActionHitTrigger)
	Source Card
}

// TurnState is the shared turn-level context passed to Card.Play alongside the per-card
// CardState wrapper.
type TurnState struct {
	// Hand is the cards currently in hand. Starts as the dealt hand minus pitched / attacker
	// / defender cards (those have been routed by the partition). Cards that draw or tutor
	// append to Hand; alt-cost effects pop from Hand. Whatever's in Hand at end of chain
	// becomes next turn's Held cards.
	Hand []Card
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
	// non-persistent goes to graveyard" rule, fireAttackActionTriggers's aura-destroy on
	// count zero, processTriggersAtStartOfTurn's start-of-turn trigger destroy) so the
	// non-card-driven append doesn't poison cacheable.
	graveyard []Card
	// Banish holds cards moved into the banished zone this turn (e.g. an aura-banish-for-
	// arcane rider).
	Banish []Card
	// Runechants is the live count of Runechant aura tokens in play. Carries across turns.
	// CreateRunechants increments it; the attack pipeline consumes the running total on each
	// attack / weapon swing.
	Runechants int
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
	// AuraTriggers is the list of triggers from auras currently in play. Cards add entries
	// during Play via AddAuraTrigger; the sim fires matching entries on each trigger-Type
	// condition, decrements Count in place, and drops entries whose Count hits zero after
	// sending Self to the graveyard. Carries across turns.
	AuraTriggers []AuraTrigger
	// pendingNextAttackActionHit queues NextAttackActionHitTriggers; reset per permutation.
	// Lowercase so cards register through RegisterNextAttackActionHit instead of appending.
	pendingNextAttackActionHit []NextAttackActionHitTrigger

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
	// attackReactionTarget is the buff target for the currently-resolving Attack Reaction.
	// Set by the chain runner around AR.Play; ARs read it via AttackReactionTarget().
	attackReactionTarget *CardState
	// Revealed is the side channel start-of-turn AuraTrigger handlers use to move a card
	// from the top of the post-draw deck into the hand (Sigil of the Arknight's reveal).
	Revealed []Card
	// TriggeringCard is the card whose play caused the active aura attack-action trigger
	// to fire. The sim sets it before each AuraTrigger handler runs and clears it after;
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
	// ClashValue). Set to false by the accessor on first card-driven access; never restored
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

// RegisterNextAttackActionHit queues t. See NextAttackActionHitTrigger for resolution.
func (s *TurnState) RegisterNextAttackActionHit(t NextAttackActionHitTrigger) {
	s.pendingNextAttackActionHit = append(s.pendingNextAttackActionHit, t)
}

// PendingNextAttackActionHits returns the number of currently queued triggers. For tests.
func (s *TurnState) PendingNextAttackActionHits() int {
	return len(s.pendingNextAttackActionHit)
}

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

// NewTurnState constructs a *TurnState with the given deck and graveyard seed. Test /
// utility constructor: the unexported deck / graveyard fields aren't reachable via a
// composite literal from outside this package, so callers in card subpackages and other
// non-sim packages route through this constructor (or set the slices via the accessor
// methods after construction). The returned state has IsCacheable()==true; cacheable
// has to be set explicitly because the field's zero value is false (see the field doc).
func NewTurnState(deck, graveyard []Card) *TurnState {
	return &TurnState{deck: deck, graveyard: graveyard, cacheable: true}
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
// same as a card that reads the deck top.
func (s *TurnState) DrawOne() {
	c, ok := s.PopDeckTop()
	if !ok {
		return
	}
	s.Hand = append(s.Hand, c)
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

// ClashValue returns the net damage-equivalent of a clash (see comprehensive rules 8.5.45):
// we and the opponent reveal the top card of our decks and the higher {p} wins. We model
// from our side only — our deck's top card is read via s.Deck(); the opponent's top is
// approximated as 5-power. So our {p} of 6-7 wins (credit +bonus), 5 ties (credit 0), and
// anything below 5 loses (credit -bonus). Returns 0 when the deck is empty. Reading the
// deck top through Deck() flips IsCacheable to false — a clash result depends on hidden
// shuffle order.
func (s *TurnState) ClashValue(bonus int) int {
	deck := s.Deck()
	if len(deck) == 0 {
		return 0
	}
	switch top := deck[0].Attack(); {
	case top >= 6:
		return bonus
	case top == 5:
		return 0
	default:
		return -bonus
	}
}

// CreateRunechants adds n Runechant token auras to the count, sets AuraCreated so effects
// that key on "aura created this turn" see it, and returns n — each token is credited as
// +1 damage at creation time when callers pair this with AddValue. Tokens that never fire
// (end-of-sim leftovers) are slightly over-credited — accepted.
func (s *TurnState) CreateRunechants(n int) int {
	if n > 0 {
		s.AuraCreated = true
		s.Runechants += n
	}
	return n
}

// CreateRunechant is shorthand for CreateRunechants(1).
func (s *TurnState) CreateRunechant() int {
	return s.CreateRunechants(1)
}

// DealArcaneDamage credits n arcane damage and, when LikelyDamageHits(n, false) approves,
// flips ArcaneDamageDealt so same-turn triggers reading "if you've dealt arcane damage this
// turn" fire. Pair with AddValue to credit the damage-equivalent. Returns n so callers can
// fold the call into a single AddValue argument.
func (s *TurnState) DealArcaneDamage(n int) int {
	if LikelyDamageHits(n, false) {
		s.ArcaneDamageDealt = true
	}
	return n
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

// AddAuraTrigger is the Play-side combo every Action - Aura card reaches for: flip
// AuraCreated so same-turn "if you've played or created an aura" riders see the entry, and
// append t to s.AuraTriggers so the sim fires it on its matching Type condition.
func (s *TurnState) AddAuraTrigger(t AuraTrigger) {
	s.AuraCreated = true
	s.AuraTriggers = append(s.AuraTriggers, t)
}

// RegisterStartOfTurn registers a TriggerStartOfTurn AuraTrigger as the canonical shape for
// "at the beginning of your action phase ..." aura clauses. self is the aura card (used by
// the sim to graveyard the source after the final fire); count is how many start-of-turn
// fires the aura survives before the sim destroys it (1 for one-shot destroy-on-fire auras,
// N for verse-counter / charge-counter auras); text is the per-fire effect description
// ("Gained 1 health", "Created a runechant", …) auto-logged alongside the trigger so the
// printout names what happened — pass "" when the handler authors its own log line (dynamic
// wording, e.g. Sigil of the Arknight's "drew X into hand"); handler runs each fire and
// returns the damage-equivalent the trigger credits.
//
// When text is non-empty, the framework's start-of-turn fire path writes a post-trigger
// log entry "<DisplayName>: text" attributed to self after handler returns and only when
// the handler returned n > 0. The pre-built LogText is stored on the AuraTrigger so the
// per-fire path runs zero string allocations (no per-Play closure either — handler stays a
// top-level function).
func (s *TurnState) RegisterStartOfTurn(self Card, count int, text string, handler OnAuraTrigger) {
	var logText string
	if text != "" {
		logText = DisplayName(self) + ": " + text
	}
	s.AddAuraTrigger(AuraTrigger{
		Self:    self,
		Type:    TriggerStartOfTurn,
		Count:   count,
		Handler: handler,
		LogText: logText,
	})
}
