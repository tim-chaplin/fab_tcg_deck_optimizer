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
	// (ephemeral attack triggers). The format layer attaches it to the previous chain
	// entry whose name matches Source.
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
	// ActionPoints is the chain runner's running Action Point pool. Seeded to 1 at the start
	// of each permutation, decremented by 1 before each non-Instant card resolves, and
	// incremented by 1 after a card with Go again (printed or granted) resolves. Action
	// cards (attack and non-attack) and weapon swings all cost 1 AP; only Instants
	// (card.TypeInstant) cost 0 AP. The chain is illegal when a non-Instant card would
	// resolve with no AP available.
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

	// --- Transient: reset by the sim per turn / chain step ---

	// Value is the running damage-equivalent total for this chain — damage dealt + damage
	// prevented + every aura-token / hero-trigger credit. The dispatcher records the chain
	// step's Play+BonusAttack contribution via AddLogEntry; trigger handlers (hero, aura,
	// ephemeral) credit themselves the same way. The solver compares permutations on this
	// field. Reset by the sim per permutation.
	Value int
	// Log is the per-event chain trace — one entry per chain step / hero / aura /
	// ephemeral / weapon swing. Chain-step producers (the sim) call AddLogEntry; pre-
	// trigger handlers (hero / aura attack-action) call AddPreTriggerLogEntry; post-
	// trigger handlers (ephemeral attack) call AddPostTriggerLogEntry. The format layer
	// uses the entry's Kind plus Source to cluster triggers under the right parent.
	// Reset per permutation.
	Log []LogEntry
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
	// BlockTotal is the sum of Defense() across every Defend-role card in the current
	// partition. Uncapped: if the partition over-blocks, BlockTotal is the full sum, not
	// clamped to IncomingDamage.
	BlockTotal int
	// EphemeralAttackTriggers are same-turn, single-fire "next attack" triggers registered
	// by a card's Play (e.g. Mauvrion Skies's "if this hits, create Runechants" rider).
	// Don't carry across turns; reset per chain.
	EphemeralAttackTriggers []EphemeralAttackTrigger
	// Revealed is the side channel start-of-turn AuraTrigger handlers use to move a card
	// from the top of the post-draw deck into the hand (Sigil of the Arknight's reveal).
	Revealed []Card
	// TriggeringCard is the card whose play caused the active aura attack-action trigger
	// to fire. The sim sets it before each AuraTrigger handler runs and clears it after;
	// the handler reads it to attribute its log line back to the triggering card. Hero
	// and ephemeral handlers receive the triggering card as a direct arg already and don't
	// need this field. Nil during direct chain-step resolution and start-of-turn fires.
	TriggeringCard Card
	// SkipLog short-circuits Log appends for chains the caller doesn't intend to display.
	// Value is still credited (the sim's running damage tally is correct) and triggers fire
	// normally (their Value contributions still flow through), but appendLog skips the slice
	// append so most chains pay zero Log cost. The eval loop runs every turn with SkipLog=true;
	// only the rare new-deck-best turn re-runs with SkipLog=false to materialise the Log for
	// the printout. Per-shuffle Log churn was the dominant allocation source — snapshotCarry's
	// Log slice copy was the biggest single field by bytes.
	SkipLog bool

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

	// SkipLog discards the entry; skip the descriptive Sprintf + three formatCardList allocs
	// when the caller doesn't intend to display the log.
	if s.SkipLog {
		return
	}
	s.AddLogEntry(fmt.Sprintf("Opted %s, put %s on top, put %s on bottom",
		formatCardList(cards), formatCardList(top), formatCardList(bottom)), 0)
}

// formatCardList renders cs as "[name1, name2, ...]" using DisplayName for each entry, or
// "[]" when cs is empty. Used by the Opt log entry so an empty top / bottom list shows up
// as a clear no-cards token rather than blank.
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

// AddLogEntry appends a freeform chain-step log line and credits n damage-equivalent to
// s.Value. text is the rendered display string. Returns the clamped n so callers can fold
// the call into a Play return. Trigger handlers call AddPreTriggerLogEntry or
// AddPostTriggerLogEntry instead so the format layer can group them under their parent.
func (s *TurnState) AddLogEntry(text string, n int) int {
	return s.appendLog(LogEntry{Text: text, Kind: LogEntryChainStep, N: n})
}

// AddPreTriggerLogEntry appends a pre-trigger log line — a hero or aura-attack-action
// trigger that fires before its parent chain entry. text is the rendered display string
// ("Viserai created a runechant"); source is the DisplayName of the card whose play
// caused the trigger to fire. The format layer attaches this entry to the next chain
// entry whose name matches source. Returns the clamped n so handlers can fold the call
// into a single return:
//
//	return s.AddPreTriggerLogEntry("Viserai created a runechant",
//	    DisplayName(played), s.CreateRunechant())
func (s *TurnState) AddPreTriggerLogEntry(text, source string, n int) int {
	return s.appendLog(LogEntry{Text: text, Source: source, Kind: LogEntryPreTrigger, N: n})
}

// AddPostTriggerLogEntry appends a post-trigger log line — an ephemeral attack trigger
// that fires after its parent chain entry resolves. The format layer attaches this entry
// to the previous chain entry whose name matches source. Same return contract as
// AddPreTriggerLogEntry.
func (s *TurnState) AddPostTriggerLogEntry(text, source string, n int) int {
	return s.appendLog(LogEntry{Text: text, Source: source, Kind: LogEntryPostTrigger, N: n})
}

// appendLog credits the entry's N to s.Value (clamped at 0) and, when SkipLog is false,
// appends it to s.Log. SkipLog=true keeps the Value tally accurate but elides every Log
// append in the chain — used for chains that won't be displayed (every turn except the
// rare new-deck-best, which is replayed with SkipLog=false to materialise its Log).
func (s *TurnState) appendLog(e LogEntry) int {
	if e.N < 0 {
		e.N = 0
	}
	s.Value += e.N
	if !s.SkipLog {
		s.Log = append(s.Log, e)
	}
	return e.N
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

// RecordValue bumps s.Value by n, clamping at 0 (FaB damage / prevention can't drive the
// running total negative). Negative n is a no-op. Cards rarely call this directly — the
// AddLogEntry / AddPreTriggerLogEntry / AddPostTriggerLogEntry helpers credit Value while
// also appending a log entry; ApplyAndLogEffectiveAttack does the same for the chain step.
func (s *TurnState) RecordValue(n int) {
	if n <= 0 {
		return
	}
	s.Value += n
}

// ApplyAndLogEffectiveAttack is the canonical chain-step finisher every Card.Play invokes:
// appends the chain-step log entry "<DisplayName>: <VERB>[ from arsenal]" (where VERB is
// ATTACK for attack actions, WEAPON ATTACK for weapons, PLAY for everything else) and
// credits Card.Attack() + self.BonusAttack to s.Value, clamped at 0. Cards with separable
// rider effects (conditional arcane bonuses, runechant creation, on-hit credits) emit
// each rider as its own post-trigger child line via ApplyAndLogRiderOnPlay /
// CreateAndLogRunechantsOnPlay / DealAndLogArcaneDamage so the rider's contribution is
// visible in the printout instead of bundled into the chain step's (+N).
func (s *TurnState) ApplyAndLogEffectiveAttack(self *CardState) {
	n := self.EffectiveAttack()
	if n < 0 {
		n = 0
	}
	s.AddLogEntry(ChainStepText(self), n)
}

// LogPlay is the chain-step finisher for non-attack cards (auras, non-attack actions, items)
// — emits "<DisplayName>: PLAY[ from arsenal]" with no value contribution. The "(+0)" suffix
// is dropped because these cards never deal printed damage; any value they contribute lands
// via separate AddPostTriggerLogEntry / aura trigger paths.
func (s *TurnState) LogPlay(self *CardState) {
	s.AddLogEntry(ChainStepText(self), 0)
}

// ApplyAndLogEffectiveDefense is the Defense Reaction counterpart to ApplyAndLogEffectiveAttack:
// emits the "<DisplayName>: DEFENSE REACTION[ from arsenal]" chain step and credits the
// effective Defense (printed Defense + BonusDefense + ArsenalDefenseBonus when from arsenal)
// to s.Value, clamped at the remaining IncomingDamage so an over-blocked DR doesn't credit
// past what was actually prevented. The credited amount is decremented from s.IncomingDamage
// so a later defender sees the reduced pool. Cards with separable rider effects (arcane
// pings, runechant creation, on-hit credits) emit each rider as its own post-trigger child
// line via DealAndLogArcaneDamage / CreateAndLogRunechantsOnPlay / ApplyAndLogRiderOnPlay after the
// chain step; conditional "+N{d}" bonuses fold into BonusDefense before the chain step so
// they roll into the same (+N), mirroring how BonusAttack feeds ApplyAndLogEffectiveAttack.
func (s *TurnState) ApplyAndLogEffectiveDefense(self *CardState) {
	n := self.EffectiveDefense()
	if n > s.IncomingDamage {
		n = s.IncomingDamage
	}
	if n < 0 {
		n = 0
	}
	s.IncomingDamage -= n
	s.AddLogEntry(ChainStepText(self), n)
}

// CreateRunechants adds n Runechant token auras to the count, sets AuraCreated so effects
// that key on "aura created this turn" see it, and returns n — each token is credited as
// +1 damage at creation time. Tokens that never fire (end-of-sim leftovers) are slightly
// over-credited — accepted.
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

// CreateAndLogRunechants creates n Runechant tokens, writes the canonical pre-trigger log
// line ("<selfName> created a runechant" for n==1, "<selfName> created N runechants" for
// n>1) sourced under sourceName, and returns the damage-equivalent credited. Trigger
// handlers that fire before their parent (Viserai's hero ability, Malefic Incantation's
// aura) call this in a single return statement.
func (s *TurnState) CreateAndLogRunechants(selfName, sourceName string, n int) int {
	return s.AddPreTriggerLogEntry(selfName+" "+runechantsCreatedPhrase(n), sourceName, s.CreateRunechants(n))
}

// CreateAndLogRunechantsOnHit is the post-trigger variant of CreateAndLogRunechants —
// the trigger log line reads "<selfName> created N runechants on hit" so the conditional
// gate on the ephemeral attack trigger (Mauvrion Skies, Runic Reaping) is visible in the
// printout. Same return contract as CreateAndLogRunechants.
func (s *TurnState) CreateAndLogRunechantsOnHit(selfName, sourceName string, n int) int {
	return s.AddPostTriggerLogEntry(selfName+" "+runechantsCreatedPhrase(n)+" on hit", sourceName, s.CreateRunechants(n))
}

// CreateAndLogRunechantsOnPlay is the on-play self-rider variant: the chain step's own
// "Created N runechants" sub-line, sourced under self so the format layer attaches it as
// a child of self's chain entry. The line uses indentation to convey source (no card-name
// prefix, sentence-cap leading verb) since the format layer renders it indented under
// self's chain entry. n>0 only — n=0 returns 0 without writing a line.
func (s *TurnState) CreateAndLogRunechantsOnPlay(self *CardState, n int) int {
	if n <= 0 {
		return 0
	}
	var text string
	if n == 1 {
		text = "Created a runechant"
	} else {
		text = fmt.Sprintf("Created %d runechants", n)
	}
	return s.AddPostTriggerLogEntry(text, DisplayName(self.Card), s.CreateRunechants(n))
}

// runechantsCreatedPhrase returns "created a runechant" / "created N runechants" — the
// canonical verb phrase for runechant-creation log lines.
func runechantsCreatedPhrase(n int) string {
	if n == 1 {
		return "created a runechant"
	}
	return fmt.Sprintf("created %d runechants", n)
}

// DealArcaneDamage credits n arcane damage and, when LikelyDamageHits(n, false) approves,
// flips ArcaneDamageDealt so same-turn triggers reading "if you've dealt arcane damage this
// turn" fire. The value is credited unconditionally — even if the opponent expends a card or
// resource to negate it, that's still net tempo gained — so only the trigger flag is gated by
// the hit heuristic. Returns n so callers can fold the arcane damage into their Play return.
func (s *TurnState) DealArcaneDamage(n int) int {
	if LikelyDamageHits(n, false) {
		s.ArcaneDamageDealt = true
	}
	return n
}

// DealAndLogArcaneDamage is the rider-line variant: credits n arcane damage (routed through
// DealArcaneDamage so the same hit-gated ArcaneDamageDealt flip applies) and writes a
// "Dealt N arcane damage" sub-line sourced under self so the format layer attaches it as a
// child of self's chain entry. n>0 only — n=0 returns 0 without writing a line.
func (s *TurnState) DealAndLogArcaneDamage(self *CardState, n int) int {
	if n <= 0 {
		return 0
	}
	var text string
	if n == 1 {
		text = "Dealt 1 arcane damage"
	} else {
		text = fmt.Sprintf("Dealt %d arcane damage", n)
	}
	return s.AddPostTriggerLogEntry(text, DisplayName(self.Card), s.DealArcaneDamage(n))
}

// ApplyAndLogRiderOnPlay writes a freeform rider sub-line under self's chain entry. text is a
// terse description of what the rider did (e.g. "On-hit discarded a card", "Gained 3
// health"); the format layer renders it indented under self's chain step, so the line
// should read as a complete utterance without a card-name prefix. n is the damage-
// equivalent credit. Returns the credited n (clamped at 0 by the underlying log helper).
func (s *TurnState) ApplyAndLogRiderOnPlay(self *CardState, text string, n int) int {
	return s.AddPostTriggerLogEntry(text, DisplayName(self.Card), n)
}

// ApplyAndLogRiderOnHit is the on-hit-gated variant of ApplyAndLogRiderOnPlay: writes the rider sub-line
// only when self's effective attack is likely to land (per LikelyToHit), otherwise no log
// entry is written and no value is credited. Use this for "When this hits, …" riders whose
// only gate is the hit check; riders with extra preconditions (ArcaneDamageDealt,
// HasPlayedOrCreatedAura, …) should test the precondition first and call this for the hit gate.
// Don't use this when n is computed via a side-effecting expression that must fire only on
// hit (e.g. s.CreateRunechant()) — Go evaluates the argument before the call, so the side
// effect would always fire. Returns 0 on miss, the credited n on hit.
func (s *TurnState) ApplyAndLogRiderOnHit(self *CardState, text string, n int) int {
	if !LikelyToHit(self) {
		return 0
	}
	return s.ApplyAndLogRiderOnPlay(self, text, n)
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

// RegisterStartOfTurn registers a TriggerStartOfTurn AuraTrigger as the canonical shape
// for "at the beginning of your action phase ..." aura clauses. self is the aura card
// (used by the sim to graveyard the source after the final fire); count is how many
// start-of-turn fires the aura survives before the sim destroys it (1 for one-shot
// destroy-on-fire auras, N for verse-counter / charge-counter auras); text is the
// per-fire effect description ("Gained 1 health", "Created a runechant", …) auto-logged
// alongside the trigger so the printout names what happened — pass "" when the handler
// authors its own log line (dynamic wording, e.g. Sigil of the Arknight's "drew X into
// hand"); handler runs each fire and returns the damage-equivalent the trigger credits.
//
// When text is non-empty and the handler returns n > 0, the wrapper writes a post-trigger
// log entry "<DisplayName>: text (+n)" attributed to self so the sim's per-trigger
// contribution renders as a descriptive line instead of the generic "<DisplayName>:
// START OF ACTION PHASE (+N)" fallback. n == 0 fires log nothing (matches the no-banish
// edge case of Sigil of Silphidae's leave trigger). The damage accumulator in
// processTriggersAtStartOfTurn still folds n into the turn's value via the handler's
// return; the log entry's N is purely cosmetic for the rendered line.
func (s *TurnState) RegisterStartOfTurn(self Card, count int, text string, handler OnAuraTrigger) {
	finalHandler := handler
	if text != "" {
		source := DisplayName(self)
		prefix := source + ": " + text
		finalHandler = func(s *TurnState) int {
			n := handler(s)
			if n > 0 {
				s.AddPostTriggerLogEntry(fmt.Sprintf("%s (+%d)", prefix, n), source, n)
			}
			return n
		}
	}
	s.AddAuraTrigger(AuraTrigger{
		Self:    self,
		Type:    TriggerStartOfTurn,
		Count:   count,
		Handler: finalHandler,
	})
}

// AddEphemeralAttackTrigger registers a same-turn, fire-once "next attack" trigger. The sim
// stamps t.SourceIndex after the registering card's Play returns. Fires on the next
// matching attack action's resolution; fizzles silently at end of turn if no match.
func (s *TurnState) AddEphemeralAttackTrigger(t EphemeralAttackTrigger) {
	s.EphemeralAttackTriggers = append(s.EphemeralAttackTriggers, t)
}
