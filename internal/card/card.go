// Package card defines the Card interface used by the simulator and basic/test implementations.
package card

// CardType is a card-type descriptor. Each constant corresponds to one keyword from a FaB
// card's type line (e.g. "Runeblade", "Action", "Attack").
type CardType uint64

const (
	TypeAction          CardType = 1 << iota // "Action"
	TypeAttack                               // "Attack"
	TypeAttackReaction                       // "Attack Reaction"
	TypeAura                                 // "Aura"
	TypeBlock                                // "Block"
	TypeDefenseReaction                      // "Defense Reaction"
	TypeGeneric                              // "Generic"
	TypeHero                                 // "Hero"
	TypeInstant                              // "Instant"
	TypeItem                                 // "Item"
	TypeOneHand                              // "1H"
	TypeRuneblade                            // "Runeblade"
	TypeScepter                              // "Scepter"
	TypeSword                                // "Sword"
	TypeTwoHand                              // "2H"
	TypeWeapon                               // "Weapon"
	TypeYoung                                // "Young"
)

// persistsInPlayMask is the set of types that keep a card in its zone after resolving rather
// than heading to the graveyard. Auras (e.g. Sigil of the Arknight: Runeblade, Action, Aura)
// and Items linger in the arena until a destroy condition fires; weapons stay equipped.
const persistsInPlayMask TypeSet = TypeSet(TypeAura) | TypeSet(TypeItem) | TypeSet(TypeWeapon)

// PersistsInPlay reports whether a card with this type set stays in its zone when it resolves
// instead of heading to the graveyard. Used by the solver to decide whether to append a
// just-played card to state.Graveyard.
func (s TypeSet) PersistsInPlay() bool {
	return s&persistsInPlayMask != 0
}

// TypeSet is a bitfield of CardType values — type checks become a single-word bitmask AND, no
// string hashing or map lookup on the hot path.
type TypeSet uint64

// NewTypeSet returns a TypeSet containing all of the given types.
func NewTypeSet(types ...CardType) TypeSet {
	var s TypeSet
	for _, t := range types {
		s |= TypeSet(t)
	}
	return s
}

// Has reports whether s contains the given type.
func (s TypeSet) Has(t CardType) bool { return s&TypeSet(t) != 0 }

// IsNonAttackAction reports whether s represents an Action that is not also an Attack —
// the bitmask check behind every "if a non-attack action card was played/pitched" rider.
func (s TypeSet) IsNonAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) == 0
}

// IsAttackAction reports whether s is an attack action card — both Action and Attack.
// Single-expression bitmask keeps the "next attack action" peek loops lean.
func (s TypeSet) IsAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) != 0
}

// IsRunebladeAttack reports whether s is a Runeblade attack — an attack action card OR a
// weapon swing. Used by "next Runeblade attack this turn" riders that peek CardsRemaining.
func (s TypeSet) IsRunebladeAttack() bool {
	return s&TypeSet(TypeRuneblade) != 0 && s&(TypeSet(TypeAttack)|TypeSet(TypeWeapon)) != 0
}

// IsDefenseReaction reports whether s has the Defense Reaction subtype.
func (s TypeSet) IsDefenseReaction() bool {
	return s&TypeSet(TypeDefenseReaction) != 0
}

// CardState wraps a Card with per-turn mutable flags that other cards' effects can toggle.
// Instances are created by the solver at the start of each attack chain and live only for that
// chain. Effects that grant keywords to "the next X" scan TurnState.CardsRemaining and flip
// flags on the matching entry; the card currently resolving receives its own CardState as
// the `self` parameter to Play.
type CardState struct {
	Card Card
	// GrantedGoAgain is set by a prior card's grant ("next X attack" riders) or by the card's
	// own Play flipping self.GrantedGoAgain = true. The solver's chain-legality check ORs
	// this with Card.GoAgain().
	GrantedGoAgain bool
	// GrantedDominate is the Dominate counterpart to GrantedGoAgain: set by a prior card's
	// grant or by this card's own Play flipping self.GrantedDominate = true when a
	// conditional "gains dominate" clause fires. LikelyToHit ORs this with the card's
	// Dominator marker (HasDominate) to decide whether to credit the "can't over-block" bump.
	GrantedDominate bool
	// FromArsenal flags the single CardState whose Card came from the arsenal slot at start of
	// turn. The solver sets it before the chain runs; CardStates for hand cards and mid-turn
	// extensions stay false. Cards gate "if this is played from arsenal" riders on
	// self.FromArsenal.
	FromArsenal bool
}

// EffectiveGoAgain reports whether this card has Go again this turn — from printed text or a
// grant by a prior card's effect.
func (p *CardState) EffectiveGoAgain() bool {
	return p.Card.GoAgain() || p.GrantedGoAgain
}

// EffectiveDominate reports whether this card attacks with Dominate this turn — from its
// printed Dominator marker or a grant flipping GrantedDominate (either by a prior card or by
// this card's own Play when a conditional "gains dominate" clause fires).
func (p *CardState) EffectiveDominate() bool {
	return p.GrantedDominate || HasDominate(p.Card)
}

// TurnState is the shared turn-level context passed to Card.Play alongside the per-card
// CardState wrapper. Cards read it to decide what effects to apply; the solver appends each
// played card to CardsPlayed after its Play returns so later cards this turn see what was
// played before them.
type TurnState struct {
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order. Populated by
	// the solver, not by Play itself.
	CardsPlayed []Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant
	// tokens, which are auras). Effects that check "if you've played or created an aura this
	// turn" should OR this with CardsPlayed containing an Aura-typed card.
	AuraCreated bool
	// CardsRemaining is the cards that will be played after the current one in turn order.
	// Populated by the solver before each Play so an effect can peek forward ("next X
	// attack") or grant keywords to a later card by flipping flags on its CardState entry.
	CardsRemaining []*CardState
	// Pitched is the cards pitched this turn for resources. Populated by the solver before any
	// Play. Effects that check "if an attack card was pitched" scan this list.
	Pitched []Card
	// Overpower is set when an attack with the Overpower keyword is being played. Not yet
	// consumed by the solver — blocked damage should eventually be forwarded to the hero when
	// Overpower is true.
	Overpower bool
	// Deck is the cards remaining in the deck (excluding the current hand), in top-of-deck
	// order. Effects that reveal or draw the top card (e.g. Sigil of the Arknight) inspect this.
	// Nil when unknown. Implementations must not mutate it.
	Deck []Card
	// Runechants is the live count of Runechant aura tokens in play. The solver seeds it with
	// the previous turn's carryover, CreateRunechants increments it, and the attack pipeline
	// consumes the running total on each attack / weapon swing (each token fires for 1 arcane
	// and is destroyed). Variable-cost cards (e.g. cost reduced per Runechant) read this in Cost.
	Runechants int
	// ArcaneDamageDealt sticks true once any source of arcane damage fires this turn: a
	// Runechant token consuming itself on an attack / weapon swing, or a card whose Play deals
	// arcane directly. Effects that read "if you've dealt arcane damage this turn" consult
	// this flag rather than Runechants (which only shows currently-alive tokens).
	//
	// playSequence sets the flag automatically for the Runechant-firing case by checking
	// Runechants > 0 before each attack/weapon's Play runs. Cards that deal arcane via their
	// Play text are responsible for flipping the flag themselves.
	ArcaneDamageDealt bool
	// NonAttackActionPlayed is set true once any non-attack action card has been appended to
	// CardsPlayed this turn. Maintained by playSequenceWithMeta when each card resolves so
	// hero triggers that ask "was a non-attack action played earlier?" can answer in O(1)
	// instead of rescanning CardsPlayed on every trigger.
	NonAttackActionPlayed bool
	// IncomingDamage is the opponent damage this turn (the value passed to hand.Best). Constant
	// across every partition the solver enumerates for this hand.
	IncomingDamage int
	// BlockTotal is the sum of Defense() across every Defend-role card in the current partition.
	// Uncapped: if the partition over-blocks, BlockTotal is the full sum, not clamped to
	// IncomingDamage. Cards that key on "will we block all incoming this turn?" read
	// BlockTotal >= IncomingDamage.
	BlockTotal int
	// Drawn records cards this turn has drawn mid-chain via DrawOne, in draw order. The solver
	// consumes it after the initial chain: leading entries may pitch to fund an under-budgeted
	// attack, next entries may play as free-cost chain extensions, and the rest carry as Held
	// (or compete for the empty arsenal slot) into the next hand.
	Drawn []Card
	// Graveyard is the cards that have entered the graveyard this turn — every card played or
	// blocked lands here after resolving. Pitched cards do not (they go back to the deck). In
	// the defense phase, the solver seeds Graveyard with every Defend-role card so effects
	// that scan the graveyard see plain blocks and other defenders. In the attack chain, each
	// card is appended after its Play returns so later attacks see what resolved before them.
	// Cards that read Graveyard must implement NoMemo since its contents aren't captured in
	// the hand's memo key.
	Graveyard []Card
	// Banish holds cards banished this turn — moved here by effects that pull a card out of
	// the graveyard (e.g. an aura-banish-for-arcane rider). Cards that key on "was a card
	// banished this turn" read this list.
	Banish []Card
	// AuraTriggers is the list of triggers from auras currently in play. Value-typed so the
	// sim can copy-restore it cheaply between permutations of the best-line search. Cards add
	// entries during Play via AddAuraTrigger; the sim fires matching entries on each
	// trigger-Type condition (start of turn for now), decrements Count in place, and drops
	// entries whose Count hits zero after sending Self to the graveyard.
	AuraTriggers []AuraTrigger
	// EphemeralAttackTriggers are same-turn, single-fire "next attack" triggers registered by
	// a card's Play (e.g. Mauvrion Skies's "if this hits, create Runechants" rider). They
	// differ from AuraTriggers on three axes: they fizzle at end of turn rather than carrying
	// across (no cross-turn seeding), they fire at most once (no Count), and they don't send
	// a source card to the graveyard on fire or fizzle — the registering card was already
	// graveyarded when its Play resolved; only the trigger "stays in play" for the turn. The
	// sim reset per permutation is an empty slice (no cross-turn carryover).
	EphemeralAttackTriggers []EphemeralAttackTrigger
	// Revealed is the side channel start-of-turn AuraTrigger handlers use to move a card
	// from the top of the post-draw deck into the hand (Sigil of the Arknight's reveal).
	// Handlers peek s.Deck[0], append to s.Revealed, and advance s.Deck past the popped
	// card; the deck loop consumes s.Revealed after firing every start-of-turn handler and
	// appends each entry to the dealt hand in order. Cascading reveals work because each
	// handler's pop shrinks the shared Deck view for the next handler.
	Revealed []Card
}

// AuraTriggerType categorizes when an AuraTrigger's Handler fires. The sim walks the
// TurnState's AuraTriggers list on each matching condition and invokes every applicable
// handler.
type AuraTriggerType int

const (
	// TriggerStartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search. The classic upkeep trigger for "at the beginning of your action phase
	// …" auras.
	TriggerStartOfTurn AuraTriggerType = iota
	// TriggerAttackAction fires each time an attack action card resolves during the attack
	// chain. Triggers that set OncePerTurn cap themselves at one fire per turn regardless of
	// how many attack actions resolve ("once per turn, when you play an attack action card
	// …" clauses).
	TriggerAttackAction
)

// OnAuraTrigger is the business-logic callback attached to an AuraTrigger. Called when the
// trigger's Type condition fires — it's where the printed "create a runechant", "gain 1{h}",
// "reveal top of deck" effect lives. Handlers mutate the passed TurnState directly
// (e.g. s.CreateRunechants, s.AddToGraveyard) and return the damage-equivalent that folds
// 1-to-1 into Value. The sim handles the counter bookkeeping (decrementing Count,
// graveyarding the aura when Count hits zero); the handler does not.
type OnAuraTrigger func(s *TurnState) int

// AuraTrigger is a counter-tracked handler attached to an aura in play. Each time Type's
// condition fires — and, when OncePerTurn is set, at most once per turn — the sim calls
// Handler and decrements Count. When Count reaches zero the sim sends Self to the graveyard
// and drops the trigger from TurnState.AuraTriggers. Self is the aura card itself so the
// sim can graveyard it without needing a back-reference.
type AuraTrigger struct {
	// Self is the aura card this trigger belongs to. Used by the sim to graveyard the aura
	// when Count reaches zero; also surfaced in per-turn summaries (e.g. the "(from previous
	// turn)" formatter line naming the aura that fired).
	Self Card
	// Type is the condition that fires this trigger.
	Type AuraTriggerType
	// Count is the number of times this trigger will still fire before the aura is destroyed.
	Count int
	// Handler runs when Type fires.
	Handler OnAuraTrigger
	// OncePerTurn caps the trigger at a single fire per turn regardless of how many matching
	// events occur. The sim sets FiredThisTurn the first time Handler runs each turn and
	// clears it at the next turn boundary.
	OncePerTurn bool
	// FiredThisTurn is sim-managed bookkeeping for OncePerTurn. Cards must not set it.
	FiredThisTurn bool
}

// OnEphemeralAttackTrigger is the business-logic callback attached to an
// EphemeralAttackTrigger. target is the CardState of the attacker whose resolution triggered
// the fire; the handler may read target.Card.Attack(), target.EffectiveDominate(), etc. to
// decide whether a rider effect fires and what damage-equivalent to credit. Handlers mutate
// the passed TurnState directly (e.g. s.CreateRunechants) and return the damage-equivalent;
// the sim routes that damage back to the trigger's Source for per-card attribution.
type OnEphemeralAttackTrigger func(s *TurnState, target *CardState) int

// EphemeralAttackTrigger is a same-turn, fire-once "next attack" trigger registered by a
// card's Play (via TurnState.AddEphemeralAttackTrigger). Fires on the next attack action
// whose resolution matches its Matches predicate, AFTER the attacker's Play, hero
// OnCardPlayed, and AuraTriggers have all settled — so the Handler sees the fully-resolved
// attacker state (incl. any Dominate grants and hero-created auras).
//
// Distinct from AuraTrigger on three axes:
//   - Fire-once. No Count / OncePerTurn bookkeeping; the trigger resolves and drops out.
//   - Doesn't graveyard a source when it fires or fizzles — the registering card was
//     already graveyarded on its own resolution; only the trigger effect "stays in play."
//   - Doesn't persist across turns. Non-matching attack actions leave the trigger in place
//     for a later match, but anything unresolved at end of turn fizzles silently.
//
// The Source card keeps damage attribution clean: Handler's return is credited to Source's
// position in the chain (via SourceIndex), so a trigger fired by Mauvrion Skies during
// Drowning Dire's attack surfaces as damage on Mauvrion's BestLine entry rather than DD's.
type EphemeralAttackTrigger struct {
	// Source is the card that registered the trigger. Damage the handler returns accrues to
	// Source's per-card attribution; also surfaces in per-turn debug output.
	Source Card
	// Matches decides whether the trigger fires on a given attack action card's
	// resolution. Nil matches any attack action — non-matching resolutions leave the
	// trigger in place for a later attack.
	Matches func(target *CardState) bool
	// Handler runs when the trigger fires. Returns damage-equivalent credited to Source.
	Handler OnEphemeralAttackTrigger
	// SourceIndex is sim-managed bookkeeping: the position of Source in the played-chain
	// permutation, used to route Handler's damage back to Source's perCardOut slot. Cards
	// must not set it — the solver stamps it on registration.
	SourceIndex int
}

// DrawOne models a mid-turn draw: advance the deck by one card and append it to Drawn. No-op
// on an empty deck. Every draw-rider card routes through this helper.
func (s *TurnState) DrawOne() {
	if len(s.Deck) == 0 {
		return
	}
	s.Drawn = append(s.Drawn, s.Deck[0])
	s.Deck = s.Deck[1:]
}

// Hero is the minimal hero profile card effects need. Narrower than hero.Hero to avoid an
// import cycle; package simstate holds the active hero for the run.
type Hero interface {
	Name() string
	Intelligence() int
}

// HasPlayedType reports whether any card played this turn has the given type in its Types() set.
func (s *TurnState) HasPlayedType(t CardType) bool {
	for _, c := range s.CardsPlayed {
		if c.Types().Has(t) {
			return true
		}
	}
	return false
}

// HasAuraInPlay reports whether an aura was played or created this turn — the condition
// behind "if you've played or created an aura this turn" riders. Checks AuraCreated (set by
// CreateRunechants, Sigil plays, etc.) OR scans CardsPlayed for an Aura-typed card: the flag
// covers token creation, the scan covers explicit Aura cards.
func (s *TurnState) HasAuraInPlay() bool {
	return s.AuraCreated || s.HasPlayedType(TypeAura)
}

// ClashValue returns the net damage-equivalent of a clash (see comprehensive rules 8.5.45): we
// and the opponent reveal the top card of our decks and the higher {p} wins. We model from our
// side only — our deck's top card is read from s.Deck; the opponent's top is approximated as
// 5-power (the median of an aggressive FaB deck). So our {p} of 6-7 wins (credit +bonus), 5
// ties (credit 0), and anything below 5 loses (credit -bonus: the bonus accrues to the
// opponent in those cases).
//
// bonus is the damage-equivalent of whatever the clash winner receives. Returns 0 when
// s.Deck is empty: no card to reveal means the clash effect fails per rule 8.5.45d.
func ClashValue(s *TurnState, bonus int) int {
	if len(s.Deck) == 0 {
		return 0
	}
	switch top := s.Deck[0].Attack(); {
	case top >= 6:
		return bonus
	case top == 5:
		return 0
	default:
		return -bonus
	}
}

// CreateRunechants adds n Runechant token auras to the count, sets AuraCreated so effects that
// key on "aura created this turn" see it, and returns n — each token is credited as +1 damage
// at creation time (it'll fire on some future attack, possibly via carryover). The attack
// pipeline consumes state.Runechants without re-crediting damage so every token counts once.
// Tokens that never fire (end-of-sim leftovers) are slightly over-credited — accepted.
func (s *TurnState) CreateRunechants(n int) int {
	if n > 0 {
		s.AuraCreated = true
		s.Runechants += n
	}
	return n
}

// CreateRunechant is shorthand for CreateRunechants(1) for the common single-token case.
func (s *TurnState) CreateRunechant() int {
	return s.CreateRunechants(1)
}

// DealArcaneDamage flips ArcaneDamageDealt so same-turn triggers reading "if you've dealt
// arcane damage this turn" fire, and returns n so callers can fold the arcane damage into
// their Play return in one expression (e.g. `return attack + s.DealArcaneDamage(1)`).
func (s *TurnState) DealArcaneDamage(n int) int {
	s.ArcaneDamageDealt = true
	return n
}

// AddToGraveyard appends c to s.Graveyard so later-resolving cards see it. Persistent-type
// cards (Auras, Items) don't enter the graveyard on play, so effects that destroy or banish
// themselves mid-chain route through here to make the move visible to downstream readers.
func (s *TurnState) AddToGraveyard(c Card) {
	s.Graveyard = append(s.Graveyard, c)
}

// AddAuraTrigger is the Play-side combo every Action - Aura card reaches for: flip AuraCreated
// so same-turn "if you've played or created an aura" riders see the entry, and append t to
// s.AuraTriggers so the sim fires it on its matching Type condition. Pairing the two in one
// method keeps a card from accidentally advertising the aura without the trigger or vice
// versa. The sim owns the trigger's lifecycle from here on: ticking Count and graveyarding
// Self when Count hits zero.
func (s *TurnState) AddAuraTrigger(t AuraTrigger) {
	s.AuraCreated = true
	s.AuraTriggers = append(s.AuraTriggers, t)
}

// AddEphemeralAttackTrigger registers a same-turn, fire-once "next attack" trigger. The sim
// stamps t.SourceIndex after the registering card's Play returns, so cards don't need to
// (and must not) set it. Fires on the next matching attack action's resolution; fizzles
// silently at end of turn if no match occurs.
func (s *TurnState) AddEphemeralAttackTrigger(t EphemeralAttackTrigger) {
	s.EphemeralAttackTriggers = append(s.EphemeralAttackTriggers, t)
}

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static
// profile plus a Play hook for on-play logic.
type Card interface {
	// ID returns the card's canonical registry identifier. Stable within a build. Lets callers
	// key maps / slices on cards without string-hashing Name().
	ID() ID
	Name() string
	// Cost returns the card's current resource cost given the turn state. Cards with a static
	// printed cost ignore s and return a constant; cards that read s (e.g. discount-per-token
	// effects) additionally implement VariableCost so the solver can pre-screen with cheap
	// MinCost / MaxCost bounds before enumerating chain permutations.
	Cost(s *TurnState) int
	Pitch() int
	// Attack is the printed attack value. Conditional bonuses belong in Play, not here.
	Attack() int
	Defense() int
	// Types returns the card's type-line descriptors as a TypeSet bitfield, e.g.
	// NewTypeSet(TypeRuneblade, TypeAction, TypeAttack).
	Types() TypeSet
	// GoAgain reports whether playing this card grants an additional action point. Cards
	// printed with "Go again" return true.
	GoAgain() bool
	// Play is called when the card resolves — as an attack or as a defense reaction. Returns
	// damage dealt to the opposing hero (may differ from Attack() after conditional bonuses) and
	// may read state to decide effects. self is the CardState wrapper for this resolution:
	// cards read self.FromArsenal for arsenal-gated riders and write self.GrantedGoAgain = true
	// to grant themselves Go again. When called on a defense reaction, the returned damage is
	// added uncapped to the turn's dealt total (the incoming-damage cap applies only to
	// Defense()).
	Play(s *TurnState, self *CardState) int
}

// NoMemo is an optional marker. Cards that implement it opt out of the hand-evaluation memo —
// typically because the card's Play output depends on context (e.g. remaining deck composition)
// that the memo key doesn't capture.
type NoMemo interface {
	NoMemo()
}

// VariableCost is optionally implemented by cards whose Cost(s) varies with TurnState (e.g.
// discount-per-token effects). MinCost and MaxCost are static bounds on the Cost output across
// any state; the solver uses them for cheap O(1) pre-screens before enumerating chain
// permutations. Non-implementers must return the same value for Cost(s) regardless of s.
type VariableCost interface {
	MinCost() int
	MaxCost() int
}

// NotSilverAgeLegal is an optional marker. Cards that implement it signal they're banned in the
// Silver Age format and must be excluded from format-restricted deck pools. Source of truth is
// data_sources/silver_age_banlist.txt — keep the two in sync.
type NotSilverAgeLegal interface {
	NotSilverAgeLegal()
}

// NotImplemented is an optional marker. Cards whose printed effect references mechanics the
// simulator doesn't model (e.g. Gold / Silver / Copper token economies, Landmarks) opt in so
// random deck generation and mutation pools skip them. A NotImplemented card is still a valid
// Card — hands that already contain one still evaluate using its static Attack / Pitch /
// Defense — but the optimizer won't introduce it into a new deck or swap one in via mutation.
// Orthogonal to NotSilverAgeLegal: a card can be format-legal yet unimplemented, or banned yet
// fully implemented, or both.
type NotImplemented interface {
	NotImplemented()
}

// Dominator is an optional marker. Attack action cards printed with the Dominate keyword
// implement it; the defender is capped at one blocking card, so LikelyToHit credits the
// "slips past one block" bump at 5+ power. Conditional grants ("if X, it gains dominate")
// stay off this marker and flow through CardState.GrantedDominate instead.
type Dominator interface {
	Dominate()
}

// HasDominate reports whether c is printed with the Dominate keyword — a type assertion to
// the Dominator marker. Used by CardState.EffectiveDominate and any future scanner that
// needs the static printed-keyword check without going through a CardState.
func HasDominate(c Card) bool {
	_, ok := c.(Dominator)
	return ok
}

// LowerHealthWanter is an optional Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than their opponent (deck building, sandbagging, self-damage) opt in. Cards with a
// "less {h} than an opposing hero" rider assume the clause always fires for these heroes and never
// fires for anyone else — a coarse proxy that skips per-turn life tracking.
type LowerHealthWanter interface {
	WantsLowerHealth()
}

// AddsFutureValue is an optional marker for cards whose printed effect delivers value on a
// LATER turn rather than the one they're played — next-turn triggers, cross-turn counters,
// and the like. The solver uses it as a beatsBest tiebreaker: at equal current-turn Value
// and equal leftover-runechants, a partition that plays more AddsFutureValue cards wins,
// because their hidden future payoff isn't reflected in this turn's score. Without the
// bias, a lone future-value aura loses to Held → arsenal promotion on the arsenal-occupancy
// tiebreak.
//
// The marker is intentionally decoupled from AuraTrigger so future hidden-value mechanisms
// can opt in without piggybacking on the trigger system.
type AddsFutureValue interface {
	AddsFutureValue()
}

// ArsenalDefenseBonus is an optional marker for Defense Reactions whose printed text grants
// extra defense only when the card is played from arsenal. Implementers return the
// additional defense added to Defense() when this copy came from the arsenal slot at start
// of turn. Defense() itself stays the printed value so the hand-played path is unaffected.
type ArsenalDefenseBonus interface {
	ArsenalDefenseBonus() int
}

