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
	TypeDefenseReaction                      // "Defense Reaction"
	TypeGeneric                              // "Generic"
	TypeHero                                 // "Hero"
	TypeInstant                              // "Instant"
	TypeOneHand                              // "1H"
	TypeRuneblade                            // "Runeblade"
	TypeScepter                              // "Scepter"
	TypeSword                                // "Sword"
	TypeTwoHand                              // "2H"
	TypeWeapon                               // "Weapon"
	TypeYoung                                // "Young"
)

// graveyardOnResolveMask is the set of types that hit the graveyard the moment they resolve:
// Action, Attack Reaction, Defense Reaction, Instant. Persistent types (Aura today; Item,
// Weapon, Hero, etc. for future card roll-outs) stay in their zone until a card-specific
// destroy condition fires even when they also carry one of the four resolve-bound subtypes.
const graveyardOnResolveMask TypeSet = TypeSet(TypeAction) | TypeSet(TypeAttackReaction) |
	TypeSet(TypeDefenseReaction) | TypeSet(TypeInstant)

// persistsInPlayMask is the set of types that keep a card in the arena (or a dedicated zone)
// after resolving. When present the card doesn't hit the graveyard on resolve even if the
// card's type line also includes Action / Attack / etc. — aura-actions like Sigil of the
// Arknight are the canonical case: typed Runeblade/Action/Aura yet they linger until PlayNextTurn
// destroys them. Keep this mask in sync with the set of implemented persistent types.
const persistsInPlayMask TypeSet = TypeSet(TypeAura)

// GraveyardOnResolve reports whether a card with this type set goes to the graveyard the moment
// it resolves. Used by the solver to decide whether to append a just-played card to
// state.Graveyard (true) or leave it on the battlefield (false — a separate destroy event will
// move it, typically via card.DelayedPlay).
func (s TypeSet) GraveyardOnResolve() bool {
	if s&persistsInPlayMask != 0 {
		return false
	}
	return s&graveyardOnResolveMask != 0
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

// IsNonAttackAction reports whether s represents an Action that is not also an Attack. Used by
// effects keyed on "if a non-attack action card was played/pitched" (Viserai's trigger, Vigor
// Rush's go-again rider, Aether Slash's arcane rider, Deathly Duet's runechant rider, Nebula
// Blade's +3 power rider). A single bitmask check avoids duplicating Has(Action) && !Has(Attack)
// in every caller.
func (s TypeSet) IsNonAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) == 0
}

// IsRunebladeAttack reports whether s is a Runeblade attack — an attack action card OR a weapon
// swing. Used by "next Runeblade attack this turn" riders (Mauvrion Skies, Runic Reaping, Oath of
// the Arknight, Condemn to Slaughter) that peek CardsRemaining.
func (s TypeSet) IsRunebladeAttack() bool {
	return s&TypeSet(TypeRuneblade) != 0 && s&(TypeSet(TypeAttack)|TypeSet(TypeWeapon)) != 0
}

// IsDefenseReaction reports whether s has the Defense Reaction subtype. Named because five
// sites in the solver (partition-scratch isDR precompute, per-turn summary grouping,
// defenseReactionDamage filter, fillContributions per-card Play re-play) all reach for the same
// bit.
func (s TypeSet) IsDefenseReaction() bool {
	return s&TypeSet(TypeDefenseReaction) != 0
}

// PlayedCard wraps a Card with per-turn mutable flags that other cards' effects can toggle.
// Instances are created by the solver at the start of each attack chain and live only for that
// chain. Effects that grant keywords to "the next X" scan TurnState.CardsRemaining and flip
// flags on the matching entry.
type PlayedCard struct {
	Card Card
	// GrantedGoAgain is set by a prior card's effect to give this specific card Go again even if
	// its printed text doesn't (e.g. Mauvrion Skies targeting the next Runeblade attack). The
	// solver's chain-legality check ORs this with Card.GoAgain().
	GrantedGoAgain bool
	// FromArsenal flags the single PlayedCard whose Card came from the arsenal slot at start of
	// turn. The solver sets it before the chain runs; PlayedCards for hand cards and mid-turn
	// extensions stay false.
	FromArsenal bool
}

// EffectiveGoAgain reports whether this card has Go again this turn — from printed text or a
// grant by a prior card's effect.
func (p *PlayedCard) EffectiveGoAgain() bool {
	return p.Card.GoAgain() || p.GrantedGoAgain
}

// TurnState is the context passed to Card.Play. Cards read it to decide what effects to apply;
// the solver appends each played card to CardsPlayed after its Play returns so later cards this
// turn see what was played before them.
type TurnState struct {
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order. Populated by
	// the solver, not by Play itself.
	CardsPlayed []Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant
	// tokens, which are auras). Effects that check "if you've played or created an aura this
	// turn" should OR this with CardsPlayed containing an Aura-typed card.
	AuraCreated bool
	// CardsRemaining is the cards that will be played after the current one in turn order.
	// Populated by the solver before each Play so an effect can peek forward (e.g. Condemn to
	// Slaughter buffing the "next Runeblade attack") or grant keywords to a later card by
	// flipping flags on its PlayedCard entry (e.g. Mauvrion Skies granting Go again).
	CardsRemaining []*PlayedCard
	// Pitched is the cards pitched this turn for resources. Populated by the solver before any
	// Play. Effects that check "if an attack card was pitched" scan this list.
	Pitched []Card
	// SelfFromArsenal is true for the single Play call whose card came from the arsenal slot at
	// start of turn. The solver flips it on before calling Play and clears it afterwards.
	// Effects gated on "if this is played from arsenal" read PlayedFromArsenal(s), which checks
	// this flag.
	SelfFromArsenal bool
	// SelfGoAgain is set by a card's Play to grant itself Go again for this chain (e.g. when
	// Runerager Swarm's aura-played-this-turn rider fires). The solver reads it after Play
	// returns and, if true, marks the card's PlayedCard.GrantedGoAgain.
	SelfGoAgain bool
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
	// arcane directly (e.g. Arcanic Crackle, Vexing Malice, Sigil of Suffering). Effects that
	// read "if you've dealt arcane damage this turn" consult this flag rather than Runechants
	// (which only shows currently-alive tokens).
	//
	// playSequence sets the flag automatically for the Runechant-firing case by checking
	// Runechants > 0 before each attack/weapon's Play runs. Cards that deal arcane via their
	// Play text are responsible for flipping the flag themselves.
	ArcaneDamageDealt bool
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
	// the defense phase, the solver seeds Graveyard with every Defend-role card so effects like
	// Weeping Battleground's aura banish see plain-blocked auras. In the attack chain, each
	// card is appended after its Play returns so later attacks see what resolved before them.
	// Cards that read Graveyard must implement NoMemo since its contents aren't captured in
	// the hand's memo key.
	Graveyard []Card
	// Banish holds cards banished this turn — moved here by effects that pull a card out of
	// the graveyard (e.g. Weeping Battleground banishing an aura). Cards that key on "was a
	// card banished this turn" read this list.
	Banish []Card
}

// PlayedFromArsenal reports whether the card currently being played came from the arsenal
// slot. Reads s.SelfFromArsenal — the solver flips that flag on before calling Play for the
// arsenal-in attacker.
func PlayedFromArsenal(s *TurnState) bool {
	return s != nil && s.SelfFromArsenal
}

// AddToGraveyard moves c into the graveyard — the single entry point every card implementation
// uses when an aura / item / other persistent card leaves the arena, whether that's during a
// mid-turn self-destroy (e.g. a fragile aura taking unblocked damage) or from a
// DelayedPlay.PlayNextTurn callback.
func (s *TurnState) AddToGraveyard(c Card) {
	s.Graveyard = append(s.Graveyard, c)
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

// HasAuraInPlay reports whether an aura was played or created this turn — the condition six
// "if you've played or created an aura this turn" riders check (Reek of Corruption, Hit the High
// Notes, Shrill of Skullform, Vantage Point, Runerager Swarm, Yinti Yanti). Checks the
// AuraCreated flag (set by CreateRunechants, Sigil plays, etc.) OR scans CardsPlayed for an
// Aura-typed card — the flag covers token creation, the scan covers explicit Aura cards.
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
	// may read state to decide effects. When called on a defense reaction, the returned damage
	// is added uncapped to the turn's dealt total (the incoming-damage cap applies only to
	// Defense()).
	Play(s *TurnState) int
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

// LowerHealthWanter is an optional Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than their opponent (deck building, sandbagging, self-damage) opt in. Cards with a
// "less {h} than an opposing hero" rider assume the clause always fires for these heroes and never
// fires for anyone else — a coarse proxy that skips per-turn life tracking.
type LowerHealthWanter interface {
	WantsLowerHealth()
}

// DelayedPlay is an optional marker for cards whose effect fires at the START of the owner's
// NEXT action phase rather than the turn they're played. A card that implements this is still
// played normally — Play runs this turn, typically to flip AuraCreated so same-turn aura-readers
// see it — and is then queued for a PlayNextTurn callback that fires at the top of the next
// turn, after the hand is drawn but before the best-line search.
//
// Cross-turn auras whose printed text reads "At the beginning of your action phase, destroy
// this. When this leaves the arena, <effect>" belong here — Sigil of the Arknight's next-turn
// reveal, Sigil of Fyendal's next-turn 1{h} gain.
//
// The TurnState passed to PlayNextTurn has Deck populated with the remaining deck after the
// next hand has been drawn (so Deck[0] is the card about to be revealed by a top-of-deck
// effect); every other field is zero.
//
// PlayNextTurn fires exactly once, at the top of the turn after the card was played. Cards
// that leave the arena at that point call s.AddToGraveyard(self) to move themselves to the
// graveyard; cards that return something to the hand set ToHand on the result. Effects that
// should carry across additional turns have to be modelled separately — there's no automatic
// re-queue.
type DelayedPlay interface {
	PlayNextTurn(s *TurnState) DelayedPlayResult
}

// DelayedPlayResult is what a DelayedPlay callback returns. Damage is credited 1-to-1 toward
// the next turn's Value. ToHand, when non-nil, is popped off the top of the post-draw deck
// and appended to that turn's hand — modelling "reveal top of deck; if <condition>, put it
// into your hand" without collapsing the effect into a flat damage-equivalent. Callbacks that
// don't reveal leave ToHand nil; callbacks that don't credit damage leave Damage 0.
type DelayedPlayResult struct {
	Damage int
	ToHand Card
}

// ArsenalDefenseBonus is an optional marker for Defense Reactions whose printed text grants
// extra defense only when the card is played from arsenal (e.g. Unmovable, Springboard
// Somersault). Implementers return the additional defense added to Defense() when this copy
// came from the arsenal slot at start of turn. Defense() itself stays the printed value so the
// hand-played path is unaffected.
type ArsenalDefenseBonus interface {
	ArsenalDefenseBonus() int
}

