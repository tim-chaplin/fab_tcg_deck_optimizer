package sim

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// LogEntry, LogEntryKind, and the LogEntry* constants live in v2/turnlogger; the
// aliases below let sibling files in this package keep referencing sim.LogEntry /
// sim.LogEntryChainStep without churn while the refactor is in flight.
type (
	LogEntry     = turnlogger.LogEntry
	LogEntryKind = turnlogger.LogEntryKind
)

const (
	LogEntryChainStep   = turnlogger.LogEntryChainStep
	LogEntryPreTrigger  = turnlogger.LogEntryPreTrigger
	LogEntryPostTrigger = turnlogger.LogEntryPostTrigger
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

// TurnState is the shared turn-level context passed to Card.Play alongside the per-card
// CardState wrapper.
type TurnState struct {
	// hand is the cards currently in hand. Starts as the dealt hand minus pitched /
	// attacker / defender cards (routed by the partition); cards drawn or tutored mid-
	// chain land here too. Whatever's left at end of chain becomes next turn's Held set.
	// Card subpackages reach hand only through Hand() / AppendHand / PopHandAt — see the
	// package-level docstring for the framework-vs-card boundary rationale.
	hand []card.Card
	// deck is the chain's working deck. Unexported so card subpackages can only reach
	// it via the public Deck() / PopDeckTop / PrependToDeck / TutorFromDeck accessors, each
	// of which clears cacheable. Framework code in this package reads / writes deck directly
	// (resetStateForPermutation seed, snapshotCarry copy) so the non-card-driven path
	// doesn't poison the cacheable bit. Always non-nil during a chain — pointed at a
	// per-permutation Copy of the leaf's starting deck.
	deck *deck.Deck
	// Arsenal is the arsenal slot's contents at this point in the chain — the arsenal-in
	// card at start of turn, nil after it plays / defends, refilled post-chain by the
	// arsenal-promotion step. Cards that read "from arsenal" use CardState.FromArsenal,
	// not this field.
	Arsenal card.Card
	// graveyard is cards that have entered the graveyard this turn — every card played or
	// blocked lands here after resolving. Pitched cards do not (they recycle to deck
	// bottom). Unexported for the same reason as deck: cards reach it only via Graveyard()
	// / BanishFromGraveyard / AddToGraveyard, all of which clear cacheable. Framework
	// code in this package writes graveyard directly (the dispatcher's "card resolved →
	// non-persistent goes to graveyard" rule, DestroyAura's aura-card append) so the
	// non-card-driven append doesn't poison cacheable.
	graveyard []card.Card
	// banished holds every card in the banished zone — both prior-turn carryover and
	// this-turn appends. Per FaB rules cards stay banished forever by default, so the
	// per-permutation reset re-seeds from the priorBanish snapshot rather than starting
	// empty. Unexported so the only mid-chain writer is BanishFromGraveyard, which
	// flips CardBanished alongside the append; cross-turn carry assigns the field
	// directly via the same-package next-turn construction. External readers consult
	// Banished() / CardBanished; external constructors go through TurnStateSpec.
	banished []card.Card
	// CardBanished is the per-turn flag BanishFromGraveyard sets the first time it
	// appends a card. Reset (alongside banished's per-permutation re-seed) in
	// resetStateForPermutation so "any card banished this turn" reads true only for
	// the current chain. Tremor of íArathael's +2{p} rider is the canary.
	cardBanished bool
	// ActionPoints is the chain runner's running AP pool. Seeded to 1 per permutation,
	// decremented before each paying chain step, incremented after a Go-again card resolves.
	// Free chain steps (Instants, Attack Reactions) cost 0; Action cards and weapon swings
	// cost 1. A paying card resolving with no AP available makes the chain illegal.
	actionPoints int
	// ArcaneDamageDealt sticks true once any source of arcane damage fires this turn:
	// a Runechant token consuming itself on an attack / weapon swing, or a card whose Play
	// deals arcane directly. Effects that read "if you've dealt arcane damage this turn"
	// consult this flag rather than Runechants. Reset at turn boundary.
	arcaneDamageDealt bool
	// OpponentMarked tracks Mark on the opposing hero. Set by a Mark effect, cleared
	// when an attack action card or weapon swing resolves — modelled as the opponent's
	// first incoming physical damage stripping the mark. Arcane damage (Runechants,
	// arcane riders) doesn't clear the mark. Read by "if the defending hero is marked"
	// riders. Our own hero is never modelled as marked, so "if you are marked, can't
	// play" gates on opposing cards drop out. Carries across turns via CarryState so
	// next turn's chain seed sees a Mark applied late in the previous turn.
	opponentMarked bool
	// Auras is the list of auras currently in play. Cards add entries during Play via
	// AddAura; the sim fires matching entries on each TriggerType condition and drops
	// entries whose handler called s.DestroyAura. Carries across turns.
	auras []Aura
	// Triggers is the list of one-shot deferred handlers keyed to a TriggerType. Cards
	// add entries during Play via AddTrigger; the sim fires matching entries on each
	// TriggerType condition and removes them after firing. Reset per permutation; not
	// snapshotted into CarryState (TriggerEndOfTurn fires before the carry snapshot).
	triggers []Trigger
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
	value int
	// logger is the per-turn log sink. Nil during the eval-loop's "find best" pass so
	// every Append* helper short-circuits via the nil-receiver guard on TurnLogger;
	// points at bufs.logger during the rare "replay best turn" pass that materialises
	// the printout. The chain runner threads this same value into Card.Play / Block /
	// OnHitHandler.Fire / TriggerHandler / Hero.OnCardPlayed as the explicit Logger
	// arg cards write through; framework methods that need direct access
	// (defendersDamage's per-DR seed, the LogEntries accessor) read it here.
	logger *turnlogger.TurnLogger
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order.
	// Populated by the sim after each Play returns so later cards this turn see what was
	// played before them.
	cardsPlayed []card.Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant
	// tokens). Effects that check "if you've played or created an aura this turn" should
	// OR this with CardsPlayed containing an Aura-typed card.
	auraCreated bool
	// CardsRemaining is the cards that will be played after the current one in chain order.
	// Populated by the sim before each Play so an effect can peek forward ("next X attack")
	// or grant keywords to a later card by flipping flags on its CardState entry.
	cardsRemaining []*card.CardState
	// Pitched is the cards pitched this turn for resources. Populated by the sim before any
	// Play. Effects that check "if an attack card was pitched" scan this list.
	pitched []card.Card
	// Overpower is set when an attack with the Overpower keyword is being played. Not yet
	// consumed by the sim — blocked damage should eventually be forwarded to the hero when
	// Overpower is true.
	overpower bool
	// NonAttackActionPlayed is set true once any non-attack action card has been appended to
	// CardsPlayed this turn. Maintained by the chain runner so hero triggers that ask "was a
	// non-attack action played earlier?" can answer in O(1).
	nonAttackActionPlayed bool
	// IncomingDamage is the opponent damage this turn — seeded by the sim from the value
	// passed to Best, and decremented by ApplyAndLogEffectiveDefense as defenders block.
	// Cards reading "did we block all incoming?" against the static partition aggregate use
	// BlockTotal instead.
	incomingDamage int
	// ArcaneIncomingDamage is the opponent's arcane damage this turn, seeded from the
	// -arcane-incoming flag. Cards whose riders gate on "if you've been dealt arcane damage
	// this turn" read this directly; not decremented during the chain (defending arcane isn't
	// modelled).
	arcaneIncomingDamage int
	// BlockTotal is the sum of Defense() across every Defend-role card in the current
	// partition. Uncapped: if the partition over-blocks, BlockTotal is the full sum, not
	// clamped to IncomingDamage.
	blockTotal int
	// Defenders is the partition's defender slice (DRs + plain blocks together) populated
	// by the chain runner before invoking each defender's Play / Block hook. Cards reading
	// "how many other plain blockers / DRs are defending alongside" iterate this list — the
	// defender-side counterpart to CardsRemaining for attackers.
	defenders []card.Card
	// attackReactionTarget is the buff target for the currently-resolving Attack Reaction.
	// Set by the chain runner around AR.Play; ARs read it via AttackReactionTarget().
	attackReactionTarget *card.CardState
	// TriggeringCard is the card whose play caused the active aura attack-action trigger
	// to fire. The sim sets it before each Aura handler runs and clears it after;
	// the handler reads it to attribute its log line back to the triggering card. Hero
	// and OnHit handlers receive the triggering card as a direct arg already and don't
	// need this field. Nil during direct chain-step resolution and start-of-turn fires.
	triggeringCard card.Card
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
func (s *TurnState) AttackReactionTarget() *card.CardState { return s.attackReactionTarget }

// Field accessors below: each returns or mutates one privatized struct field. Cards
// reach them through these methods so the cards-side surface matches v2/card.GameEngine.
// Setters exist only for fields cards actually mutate; other fields are seeded via
// TurnStateSpec or via dedicated mutation helpers (AddValue, AddAura, AddTrigger, …).

// ActionPoints returns the chain runner's running AP pool.
func (s *TurnState) ActionPoints() int { return s.actionPoints }

// AddActionPoints credits n AP to the running pool. Negative n is allowed (some cards
// spend AP) but callers should ensure the result stays non-negative.
func (s *TurnState) AddActionPoints(n int) { s.actionPoints += n }

// ArcaneDamageDealt reports whether any arcane damage source has fired this turn.
func (s *TurnState) ArcaneDamageDealt() bool { return s.arcaneDamageDealt }

// SetArcaneDamageDealt flips the sticky arcane-damage-dealt flag.
func (s *TurnState) SetArcaneDamageDealt(v bool) { s.arcaneDamageDealt = v }

// ArcaneIncomingDamage returns the opponent's arcane damage this turn.
func (s *TurnState) ArcaneIncomingDamage() int { return s.arcaneIncomingDamage }

// AuraCreated reports whether a card or ability has created an aura this turn.
func (s *TurnState) AuraCreated() bool { return s.auraCreated }

// SetAuraCreated flips the aura-created-this-turn flag.
func (s *TurnState) SetAuraCreated(v bool) { s.auraCreated = v }

// Auras returns the live aura set.
func (s *TurnState) Auras() []Aura { return s.auras }

// SetAuras replaces the aura set wholesale. Used by tests that seed a prior-turn aura
// carryover; production code uses AddAura, which also flips AuraCreated.
func (s *TurnState) SetAuras(a []Aura) { s.auras = a }

// BlockTotal returns the partition's uncapped defense sum.
func (s *TurnState) BlockTotal() int { return s.blockTotal }

// CardBanished reports whether any card has been banished this turn.
func (s *TurnState) CardBanished() bool { return s.cardBanished }

// CardsPlayed returns the sequence of cards played (as attacks) this turn.
func (s *TurnState) CardsPlayed() []card.Card { return s.cardsPlayed }

// SetCardsPlayed replaces the cards-played slice — used by Moon Wish's transient
// pre-append + pop around its go-again Sun Kiss invocation so the synergy fires.
func (s *TurnState) SetCardsPlayed(cs []card.Card) { s.cardsPlayed = cs }

// CardsRemaining returns the cards scheduled after the current chain step.
func (s *TurnState) CardsRemaining() []*card.CardState { return s.cardsRemaining }

// SetCardsRemaining replaces the look-ahead queue — currently used only by tests that
// seed a partial chain for predicate evaluation.
func (s *TurnState) SetCardsRemaining(cs []*card.CardState) { s.cardsRemaining = cs }

// Defenders returns the partition's defender slice (DRs + plain blocks).
func (s *TurnState) Defenders() []card.Card { return s.defenders }

// IncomingDamage returns the opponent damage left to allocate this turn.
func (s *TurnState) IncomingDamage() int { return s.incomingDamage }

// SetIncomingDamage replaces the running incoming-damage tally — used by the per-DR
// seed inside defendersDamage.
func (s *TurnState) SetIncomingDamage(n int) { s.incomingDamage = n }

// NonAttackActionPlayed reports whether any non-attack action has resolved this turn.
func (s *TurnState) NonAttackActionPlayed() bool { return s.nonAttackActionPlayed }

// OpponentMarked reports whether the opposing hero currently carries the Mark token.
func (s *TurnState) OpponentMarked() bool { return s.opponentMarked }

// SetOpponentMarked flips the opposing hero's Mark.
func (s *TurnState) SetOpponentMarked(v bool) { s.opponentMarked = v }

// Overpower reports whether the chain step currently in flight has Overpower.
func (s *TurnState) Overpower() bool { return s.overpower }

// GrantOverpower flips the Overpower flag on the current chain step and credits the
// engine's per-Overpower-grant value. Returns the credited value so the calling card
// can attribute the rider in its own log line.
func (s *TurnState) GrantOverpower(self *card.CardState) int {
	s.overpower = true
	s.AddValue(OverpowerValue)
	return OverpowerValue
}

// OpponentDiscard credits n cards' worth of damage-equivalent value for forcing the
// opponent to discard. Returns the credited value for log attribution.
func (s *TurnState) OpponentDiscard(n int) int {
	v := n * DiscardValue
	s.AddValue(v)
	return v
}

// LikelyToHit reports whether self's attack is likely to land past the opponent's
// blocks — engine-side delegation to the package heuristic.
func (s *TurnState) LikelyToHit(self *card.CardState) bool { return LikelyToHit(self) }

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit.
func (s *TurnState) LikelyDamageHits(n int, dominate bool) bool {
	return LikelyDamageHits(n, dominate)
}

// Pitched returns the cards pitched this turn for resources.
func (s *TurnState) Pitched() []card.Card { return s.pitched }

// TriggeringCard returns the card whose Play caused the currently-firing aura
// attack-action trigger, or nil outside of a trigger fire.
func (s *TurnState) TriggeringCard() card.Card { return s.triggeringCard }

// SetTriggeringCard replaces the triggering-card slot. Used by tests that drive a
// trigger handler directly; production code threads it through the trigger fire loop.
func (s *TurnState) SetTriggeringCard(c card.Card) { s.triggeringCard = c }

// Triggers returns the one-shot trigger queue.
func (s *TurnState) Triggers() []Trigger { return s.triggers }

// Value returns the running damage-equivalent total for this chain.
func (s *TurnState) Value() int { return s.value }

// SetValue replaces the running damage-equivalent total. Prefer AddValue for the usual
// "credit n" case; this exists for the rare decrement (e.g. Test of Strength's clash
// loss concedes a Gold token).
func (s *TurnState) SetValue(n int) { s.value = n }

// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. ARs use this to
// fold their +{p} buff into the buffed attack's display delta. No-op when the logger is
// nil (find-best pass) or when no chain-step entry exists yet.
func (s *TurnState) AmendLastChainStepN(n int) {
	s.logger.AmendLastChainStepN(n)
}

// Deck returns the chain-runner deck for read-only inspection and flips IsCacheable to
// false. Card handlers should not mutate the returned *deck.Deck directly; route mutations
// through PopDeckTop / PrependToDeck / Opt / TutorFromDeck / RecycleToDeckBottom so the
// cacheable bookkeeping stays consistent.
func (s *TurnState) Deck() *deck.Deck {
	s.cacheable = false
	return s.deck
}

// PeekTopN returns the top n cards of the deck (top first) without removing them and
// flips IsCacheable to false. Returns fewer cards when the deck has < n. Used by cards
// that scan or count the top N for a rider.
func (s *TurnState) PeekTopN(n int) []card.Card {
	s.cacheable = false
	top := s.deck.PeekTopN(n)
	if len(top) == 0 {
		return nil
	}
	out := make([]card.Card, len(top))
	for i, c := range top {
		out[i] = c.(card.Card)
	}
	return out
}

// Hand returns the live hand slice and flips IsCacheable to false. Cards must not mutate
// the returned slice; use AppendHand / PopHandAt for mutations. Read-only callers that
// only inspect length / contents still flip — the cache key can't depend on what the
// hand-content read produced.
func (s *TurnState) Hand() []card.Card {
	s.cacheable = false
	return s.hand
}

// AppendHand appends c to the hand, flipping IsCacheable to false. Cards that draw / tutor
// into hand mid-chain use this so the cache invalidation is automatic. Framework code in
// this package writes s.hand directly (resetStateForPermutation seeds, applyChainStep
// removes the playing card) and doesn't go through here — those mutations aren't
// card-driven and shouldn't poison the cacheable bit.
func (s *TurnState) AppendHand(c card.Card) {
	s.cacheable = false
	s.hand = append(s.hand, c)
}

// PopHandAt removes and returns the card at index i, flipping IsCacheable to false. Cards
// that pop hand cards (alt-cost effects, Moon Wish's "return a hand card to top of deck")
// use this so the cache invalidation is automatic. Panics on out-of-range i — callers
// should len-check via Hand() first.
func (s *TurnState) PopHandAt(i int) card.Card {
	s.cacheable = false
	c := s.hand[i]
	s.hand = append(s.hand[:i], s.hand[i+1:]...)
	return c
}

// TurnStateSpec is the cross-turn / test-construction input shape: every field
// outside-package callers (turntests, internal/cards tests, anything assembling a
// prior state) need to seed on a fresh TurnState. Mirrors TurnState's private fields
// in exported form so external code doesn't have to go through one Set* call per
// field. Construct via NewTurnStateFromSpec.
type TurnStateSpec struct {
	Arsenal               card.Card
	Auras                 []Aura
	Triggers              []Trigger
	Items                 []Item
	Banished              []deck.Card
	Graveyard             []deck.Card
	Pitched               []card.Card
	Defenders             []card.Card
	CardsPlayed           []card.Card
	CardsRemaining        []*card.CardState
	OpponentMarked        bool
	ArcaneDamageDealt     bool
	AuraCreated           bool
	CardBanished          bool
	Overpower             bool
	NonAttackActionPlayed bool
	ActionPoints          int
	IncomingDamage        int
	ArcaneIncomingDamage  int
	BlockTotal            int
	Value                 int
	TriggeringCard        card.Card
	AttackReactionTarget  *card.CardState
}

// NewTurnStatePtr is the pointer variant of NewTurnStateFromSpec — convenience for callers
// that pass the result by pointer in the same expression (e.g. card.Cost(s) tests). The
// returned pointer is independent of the caller's TurnStateSpec value.
func NewTurnStatePtr(spec TurnStateSpec) *TurnState {
	s := NewTurnStateFromSpec(spec)
	return &s
}

// NewTurnStateFromSpec builds a TurnState from a TurnStateSpec, sealing the unexported
// fields (banish, graveyard) inside the package. Use this when an external caller needs
// to construct a prior-turn state for EvalOneTurnForTesting; production code in package
// sim that already has the cross-turn buffers in hand can construct TurnState directly.
func NewTurnStateFromSpec(spec TurnStateSpec) TurnState {
	banished := make([]card.Card, len(spec.Banished))
	for i, c := range spec.Banished {
		banished[i] = c.(card.Card)
	}
	graveyard := make([]card.Card, len(spec.Graveyard))
	for i, c := range spec.Graveyard {
		graveyard[i] = c.(card.Card)
	}
	return TurnState{
		Arsenal:               spec.Arsenal,
		auras:                 spec.Auras,
		triggers:              spec.Triggers,
		Items:                 spec.Items,
		banished:              banished,
		graveyard:             graveyard,
		pitched:               spec.Pitched,
		defenders:             spec.Defenders,
		cardsPlayed:           spec.CardsPlayed,
		cardsRemaining:        spec.CardsRemaining,
		opponentMarked:        spec.OpponentMarked,
		arcaneDamageDealt:     spec.ArcaneDamageDealt,
		auraCreated:           spec.AuraCreated,
		cardBanished:          spec.CardBanished,
		overpower:             spec.Overpower,
		nonAttackActionPlayed: spec.NonAttackActionPlayed,
		actionPoints:          spec.ActionPoints,
		incomingDamage:        spec.IncomingDamage,
		arcaneIncomingDamage:  spec.ArcaneIncomingDamage,
		blockTotal:            spec.BlockTotal,
		value:                 spec.Value,
		triggeringCard:        spec.TriggeringCard,
		attackReactionTarget:  spec.AttackReactionTarget,
		cacheable:             true,
		currentAuraIdx:        -1,
		logger:                turnlogger.New(),
	}
}

// SetHandForTesting replaces the hand with the supplied cards. Test-only — production
// hand seeding goes through resetStateForPermutation (framework) or AppendHand (cards).
// Doesn't flip cacheable: tests that assert on hand state typically don't care about the
// cache, and a NewTurnState-then-SetHandForTesting flow shouldn't poison the bit before
// the test has run a single Play.
func (s *TurnState) SetHandForTesting(cards []card.Card) {
	s.hand = cards
}

// SetCurrentAuraIdxForTesting sets currentAuraIdx so DestroyAura resolves to the
// expected slot when a test invokes an aura handler directly. Test-only — production
// fire loops maintain currentAuraIdx themselves.
func (s *TurnState) SetCurrentAuraIdxForTesting(i int) {
	s.currentAuraIdx = i
}

// Graveyard returns the live graveyard slice and flips IsCacheable to false. Cards must
// not mutate the returned slice; use BanishFromGraveyard for mutations or AddToGraveyard
// for the deterministic append-only path.
func (s *TurnState) Graveyard() []card.Card {
	s.cacheable = false
	return s.graveyard
}

// PopDeckTop removes the top card of the deck and returns it. Returns (nil, false) when
// the deck is empty. Flips IsCacheable to false.
func (s *TurnState) PopDeckTop() (card.Card, bool) {
	s.cacheable = false
	if s.deck.Size() == 0 {
		return nil, false
	}
	return s.deck.Draw(1)[0].(card.Card), true
}

// PeekDeck returns the top card of the deck without removing it. Returns (nil, false) on
// an empty deck. Flips IsCacheable to false — observing the deck top makes the chain's
// output depend on hidden shuffle order, same as PopDeckTop.
func (s *TurnState) PeekDeck() (card.Card, bool) {
	s.cacheable = false
	top := s.deck.PeekTop()
	if top == nil {
		return nil, false
	}
	return top.(card.Card), true
}

// PrependToDeck inserts c at the top of the deck. Flips IsCacheable to false.
func (s *TurnState) PrependToDeck(c card.Card) {
	s.cacheable = false
	s.deck.PutTop([]deck.Card{c})
}

// RecycleToDeckBottom appends self.Card to the bottom of the deck and flips
// self.SkipGraveyard so the chain dispatcher skips the usual non-persistent → graveyard
// append. Models the FaB clause "put this on the bottom of its owner's deck"
// (Relentless Pursuit). Flips IsCacheable to false.
func (s *TurnState) RecycleToDeckBottom(self *card.CardState) {
	s.cacheable = false
	s.deck.PutBottom([]deck.Card{self.Card})
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
func (s *TurnState) Opt(l card.Logger, n int) {
	s.cacheable = false
	if n <= 0 || s.deck.Size() == 0 {
		return
	}
	if n > s.deck.Size() {
		n = s.deck.Size()
	}
	// Copy off the popped slice so the handler can't mutate s.deck through aliasing.
	drawn := s.deck.Draw(n)
	cards := make([]card.Card, len(drawn))
	for i, c := range drawn {
		cards[i] = c.(card.Card)
	}

	top, bottom := CurrentHero.Opt(cards)
	panicIfOptViolatesMultiset(cards, top, bottom)

	deckTop := make([]deck.Card, len(top))
	for i, c := range top {
		deckTop[i] = c
	}
	s.deck.PutTop(deckTop)
	deckBottom := make([]deck.Card, len(bottom))
	for i, c := range bottom {
		deckBottom[i] = c
	}
	s.deck.PutBottom(deckBottom)

	if OptDebug {
		fmt.Printf("Opt(%d): cards=%s -> top=%s bottom=%s\n",
			n, formatCardList(cards), formatCardList(top), formatCardList(bottom))
	}
	if l == nil {
		return
	}
	l.AppendChainStepf(0, "Opted %s, put %s on top, put %s on bottom",
		formatCardList(cards), formatCardList(top), formatCardList(bottom))
}

// formatCardList renders cs as "[name1, name2, ...]" using DisplayName for each entry, or
// "[]" when cs is empty. Used by the Opt log entry; the caller gates the call on whether
// the logger is recording so this only runs when the chain materialises its log.
func formatCardList(cs []card.Card) string {
	if len(cs) == 0 {
		return "[]"
	}
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = c.DisplayName()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// panicIfOptViolatesMultiset enforces TurnState.Opt's contract that the hero handler's
// combined (top, bottom) output is exactly the input multiset — a permutation of the
// input cards, no additions or removals. Panics with a descriptive message naming the
// failure mode (size mismatch, foreign card, or dropped card). Cards are zero-sized
// structs in production and small POD structs in tests; both flavours are usable as
// map keys for the multiset count.
func panicIfOptViolatesMultiset(in, top, bottom []card.Card) {
	if len(top)+len(bottom) != len(in) {
		panic(fmt.Sprintf("Opt: handler returned %d+%d cards, want %d (input multiset)",
			len(top), len(bottom), len(in)))
	}
	counts := make(map[card.Card]int, len(in))
	for _, c := range in {
		counts[c]++
	}
	check := func(out []card.Card, label string) {
		for _, c := range out {
			counts[c]--
			if counts[c] < 0 {
				panic(fmt.Sprintf("Opt: %s list returned card %s not in input",
					label, c.DisplayName()))
			}
		}
	}
	check(top, "top")
	check(bottom, "bottom")
	for c, n := range counts {
		if n != 0 {
			panic(fmt.Sprintf("Opt: handler dropped %d copy of %s from input", n, c.DisplayName()))
		}
	}
}

// TutorFromDeck removes and returns the highest-scoring card per score. Returns (nil,
// false) when no card scores > 0 (or the deck is empty). Flips IsCacheable to false.
func (s *TurnState) TutorFromDeck(score func(card.Card) int) (card.Card, bool) {
	s.cacheable = false
	got, ok := s.deck.Tutor(func(c deck.Card) int { return score(c.(card.Card)) })
	if !ok {
		return nil, false
	}
	return got.(card.Card), true
}

// BanishFromGraveyard removes the first graveyard card matching pred, appends it to
// s.banished, and returns it. Returns (nil, false) when no card matches. Flips
// IsCacheable to false. Reads graveyard contents from a previous turn (or from this
// turn's plain blocks the partition put there) — the chain output depends on hidden
// prior-turn state. Sets CardBanished so this-turn-banish riders fire correctly without
// scanning a slice that may contain prior-turn entries.
func (s *TurnState) BanishFromGraveyard(pred func(card.Card) bool) (card.Card, bool) {
	s.cacheable = false
	for i, c := range s.graveyard {
		if !pred(c) {
			continue
		}
		s.banished = append(s.banished, c)
		s.cardBanished = true
		s.graveyard = append(s.graveyard[:i], s.graveyard[i+1:]...)
		return c, true
	}
	return nil, false
}

// Banished returns the slice of cards in the banished zone, top-to-bottom (in landing
// order). Read-only — mutate via BanishFromGraveyard. Includes prior-turn entries since
// banished cards stay banished by default; "did anything banish THIS turn" readers
// must use CardBanished instead.
func (s *TurnState) Banished() []card.Card {
	return s.banished
}

// RecycleFromGraveyardToTop removes the first graveyard card matching pred, prepends it to
// the deck, and returns it. Returns (nil, false) when no card matches. Flips IsCacheable to
// false (reading graveyard + mutating deck). The deck mutation IS the model — the recycled
// card resurfaces in next turn's deal naturally, so callers don't credit Value here. Pair
// with a rider log line so the trace records the recycle.
func (s *TurnState) RecycleFromGraveyardToTop(pred func(card.Card) bool) (card.Card, bool) {
	return s.recycleFromGraveyard(pred, true)
}

// RecycleFromGraveyardToBottom is RecycleFromGraveyardToTop's bottom-of-deck variant: same
// IsCacheable flip, same no-Value-credit contract; the recycled card lands at the bottom of
// the deck instead of the top.
func (s *TurnState) RecycleFromGraveyardToBottom(pred func(card.Card) bool) (card.Card, bool) {
	return s.recycleFromGraveyard(pred, false)
}

func (s *TurnState) recycleFromGraveyard(pred func(card.Card) bool, toTop bool) (card.Card, bool) {
	s.cacheable = false
	for i, c := range s.graveyard {
		if !pred(c) {
			continue
		}
		s.graveyard = append(s.graveyard[:i], s.graveyard[i+1:]...)
		if toTop {
			s.deck.PutTop([]deck.Card{c})
		} else {
			s.deck.PutBottom([]deck.Card{c})
		}
		return c, true
	}
	return nil, false
}

// NewTurnState constructs a *TurnState with the given deck and graveyard seed. The
// unexported deck / graveyard fields aren't reachable via a composite literal from outside
// this package, so callers route through this constructor. The returned state has
// IsCacheable()==true (cacheable has to be set explicitly because the field's zero value
// is false — see the field doc) and a fresh recording *TurnLogger; the framework's
// resetStateForPermutation overrides the logger to nil for the eval-loop's silent
// find-best pass.
func NewTurnState(d *deck.Deck, graveyard []card.Card) *TurnState {
	if d == nil {
		d = &deck.Deck{}
	}
	return &TurnState{deck: d, graveyard: graveyard, cacheable: true, currentAuraIdx: -1, logger: turnlogger.New()}
}

// NewTurnStateFromCards is a test-only constructor that wraps a Card slice in a fresh
// *deck.Deck and forwards to NewTurnState. Lets card tests build a TurnState seeded with
// a hand-rolled deck order without each test plumbing the deck construction inline.
func NewTurnStateFromCards(deckCards, graveyard []card.Card) *TurnState {
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	return NewTurnState(deck.New(nil, nil, dc), graveyard)
}

// AddValue credits n to s.value, clamped at 0. Pair with a Log helper when you also want a
// log line; call alone for silent value (an aura that pays out without surfacing in the
// printout). Negative n is a no-op (FaB damage / prevention can't drive the running total
// negative). The convention is to put AddValue on its own line, separate from any Log call,
// so a line beginning with Log( has no side effects.
func (s *TurnState) AddValue(n int) {
	if n > 0 {
		s.value += n
	}
}

// LogEntries returns the per-event chain trace accumulated by the Log family. External
// readers (tests, format layer) use this; package-internal code reads the underlying
// logger field directly.
func (s *TurnState) LogEntries() []LogEntry { return s.logger.Entries() }

// Logger returns the chain runner's currently-active log sink as a Logger interface
// value. Returns the same value the framework threads into Card.Play so test harnesses
// (and any other framework caller that needs to invoke Card-shaped hooks directly) can
// pass it through. Nil during the find-best pass.
func (s *TurnState) Logger() card.Logger { return s.logger }

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
	for _, c := range s.cardsPlayed {
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
	return s.auraCreated || s.HasPlayedType(card.TypeAura)
}

// Clash models a clash (rule 8.5.45): we and the opponent reveal the top card of our decks
// and the higher {p} wins. We model from our side only — our deck's top is read via
// PeekDeck; the opponent's top is approximated as 5-power. On a win (our top ≥ 6), win
// fires; on a loss (our top ≤ 4), lose fires; ties (top == 5) and an empty deck fire
// neither. Either callback may be nil. PeekDeck flips IsCacheable to false — a clash
// result depends on hidden shuffle order.
func (s *TurnState) Clash(win, lose func()) {
	top, ok := s.PeekDeck()
	if !ok {
		return
	}
	power := top.Attack()
	switch {
	case power >= 6:
		if win != nil {
			win()
		}
	case power <= 4:
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
	s.auraCreated = true
	s.AddValue(n)
	for i := range s.auras {
		if s.auras[i].Self.TokenType == TokenTypeRunechant {
			s.auras[i].Count += n
			return
		}
	}
	s.auras = append(s.auras, NewRunechantAura(n))
}

// CreatePonder creates n Ponder tokens, bumping the existing aura entry's Count or
// adding a new one, and flips AuraCreated. No Value credit — see ponderAuraHandler.
func (s *TurnState) CreatePonder(n int) {
	if n <= 0 {
		return
	}
	s.auraCreated = true
	for i := range s.auras {
		if s.auras[i].Self.TokenType == TokenTypePonder {
			s.auras[i].Count += n
			return
		}
	}
	s.auras = append(s.auras, NewPonderAura(n))
}

// Runechants returns the current Runechant token count, or zero when none are in play.
func (s *TurnState) Runechants() int { return tokenCountIn(s.auras, TokenTypeRunechant) }

// Ponders returns the current Ponder token count, or zero when none are in play.
func (s *TurnState) Ponders() int { return tokenCountIn(s.auras, TokenTypePonder) }

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
// line under self's chain entry, and flips ArcaneDamageDealt when s.LikelyDamageHits(n, false)
// approves so same-turn triggers reading "if you've dealt arcane damage this turn" fire.
// Routes through dealtArcaneText[n] so the hot path avoids per-call fmt.Sprintf and
// variadic-int boxing.
func (s *TurnState) DealArcaneDamage(l card.Logger, self *card.CardState, n int) {
	s.AddValue(n)
	if s.LikelyDamageHits(n, false) {
		s.arcaneDamageDealt = true
	}
	if n >= 0 && n < len(dealtArcaneText) {
		l.AppendPostTrigger(self.Card.DisplayName(), dealtArcaneText[n], n)
		return
	}
	l.AppendPostTriggerf(self.Card.DisplayName(), n, "Dealt %d arcane damage", n)
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
func (s *TurnState) AddToGraveyard(c card.Card) {
	s.cacheable = false
	s.graveyard = append(s.graveyard, c)
}

// AddStartOfTurnAura registers a TriggerStartOfTurn aura: the handler fires at the
// start of each subsequent turn. Source and Self.Card both come from self.Card.
func (s *TurnState) AddStartOfTurnAura(self *card.CardState, handler card.AuraHandler, count int) {
	s.appendCardAura(self, TriggerStartOfTurn, handler, count, false)
}

// AddAttackActionAura registers a TriggerAttackAction aura: the handler fires once
// per attack action card that resolves. See AddOncePerTurnAttackActionAura for the
// OncePerTurn variant.
func (s *TurnState) AddAttackActionAura(self *card.CardState, handler card.AuraHandler, count int) {
	s.appendCardAura(self, TriggerAttackAction, handler, count, false)
}

// AddOncePerTurnAttackActionAura is the OncePerTurn variant of AddAttackActionAura —
// fires at most once per turn regardless of how many attack actions resolve.
func (s *TurnState) AddOncePerTurnAttackActionAura(self *card.CardState, handler card.AuraHandler, count int) {
	s.appendCardAura(self, TriggerAttackAction, handler, count, true)
}

// AddAttackAura registers a TriggerAttack aura: the handler fires each time any
// attack (attack action card or weapon swing) resolves.
func (s *TurnState) AddAttackAura(self *card.CardState, handler card.AuraHandler, count int) {
	s.appendCardAura(self, TriggerAttack, handler, count, false)
}

// AddEndOfTurnAura registers a TriggerEndOfTurn aura: the handler fires after the
// chain finishes resolving, before the carry-state snapshot.
func (s *TurnState) AddEndOfTurnAura(self *card.CardState, handler card.AuraHandler, count int) {
	s.appendCardAura(self, TriggerEndOfTurn, handler, count, false)
}

// appendCardAura is the shared body for the AddXxxAura family. Self / Source both
// derive from self.Card; AuraCreated flips so same-turn "if you've played or created
// an aura" riders see the entry.
func (s *TurnState) appendCardAura(self *card.CardState, tt TriggerType, handler card.AuraHandler, count int, oncePerTurn bool) {
	s.auraCreated = true
	s.auras = append(s.auras, Aura{
		Source:      self.Card,
		TriggerType: tt,
		Handler:     handler,
		Self:        CardOrTokenType{Card: self.Card},
		Count:       count,
		OncePerTurn: oncePerTurn,
	})
}

// AddHitTrigger registers a one-shot TriggerHit listener. filter narrows the
// qualifying hits to a card-type predicate (typically TypeSet.IsAttack or
// TypeSet.IsAttackAction); nil = any hit qualifies.
func (s *TurnState) AddHitTrigger(self *card.CardState, handler card.TriggerHandler, filter func(card.TypeSet) bool) {
	s.triggers = append(s.triggers, Trigger{
		Source:      self.Card,
		TriggerType: TriggerHit,
		TypeFilter:  filter,
		Handler:     handler,
	})
}

// AddEndOfTurnTrigger registers a one-shot TriggerEndOfTurn listener — fires after
// the chain finishes resolving but before the carry-state snapshot.
func (s *TurnState) AddEndOfTurnTrigger(self *card.CardState, handler card.TriggerHandler) {
	s.triggers = append(s.triggers, Trigger{
		Source:      self.Card,
		TriggerType: TriggerEndOfTurn,
		Handler:     handler,
	})
}

// AuraCount returns the count of live auras. Used by gates like Yinti Yanti's
// "while you control an aura" rider.
func (s *TurnState) AuraCount() int { return len(s.auras) }

// FireAuraForTesting invokes the aura at index i's handler, threading the same
// auraCtx adapter the production fire walks use. Test-only — production firing
// (fireStartOfTurn, fireAttackActionAuras, …) drives the cursor and the
// OncePerTurn / FiredThisTurn bookkeeping in addition to invoking the handler.
func (s *TurnState) FireAuraForTesting(i int) {
	s.currentAuraIdx = i
	s.currentAuraDestroyed = false
	a := &s.auras[i]
	ctx := auraCtx{a: a, s: s}
	a.Handler(s, s.logger, &ctx)
	s.currentAuraIdx = -1
}

// FireTriggerForTesting invokes the trigger at index i's handler with the trigger
// adapter the production trigger walks use. Test-only.
func (s *TurnState) FireTriggerForTesting(i int) {
	t := &s.triggers[i]
	ctx := triggerCtx{t: t}
	t.Handler(s, s.logger, &ctx)
}

// destroyAura removes a from s.auras and, when addToGraveyard, appends a.Self.Card to
// s.graveyard. Token-style auras (Self.Card == nil) skip the graveyard append
// unconditionally. The splice uses currentAuraIdx (the fire loop maintains it through
// any mid-handler s.auras realloc), so it resolves to the correct live slot even when
// CreateRunechants has appended a new aura since the handler started; a is read for
// Self only and may legitimately point at the pre-realloc backing.
//
// Direct graveyard append (no cacheable flip): destruction is deterministic from the
// triggering event the sim already accounts for, not from hidden state.
func (s *TurnState) destroyAura(a *Aura, addToGraveyard bool) {
	if addToGraveyard && a.Self.Card != nil {
		s.graveyard = append(s.graveyard, a.Self.Card)
	}
	i := s.currentAuraIdx
	if i < 0 || i >= len(s.auras) {
		return
	}
	s.auras = append(s.auras[:i], s.auras[i+1:]...)
	s.currentAuraDestroyed = true
}
