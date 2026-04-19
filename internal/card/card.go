// Package card defines the Card interface used by the simulator and basic/test implementations.
package card

// CardType is a card-type descriptor. Each constant corresponds to one keyword from a FaB
// card's type line (e.g. "Runeblade", "Action", "Attack").
type CardType uint64

const (
	TypeAction          CardType = 1 << iota // "Action"
	TypeAttack                               // "Attack"
	TypeAura                                 // "Aura"
	TypeDefenseReaction                      // "Defense Reaction"
	TypeGeneric                              // "Generic"
	TypeHero                                 // "Hero"
	TypeOneHand                              // "1H"
	TypeRuneblade                            // "Runeblade"
	TypeScepter                              // "Scepter"
	TypeSword                                // "Sword"
	TypeTwoHand                              // "2H"
	TypeWeapon                               // "Weapon"
	TypeYoung                                // "Young"
)

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
	// Self is the PlayedCard wrapper for the card currently being played. Effects that
	// conditionally grant the played card itself Go again (e.g. Runerager Swarm) flip
	// Self.GrantedGoAgain. The solver populates Self before each Play and consults
	// EffectiveGoAgain after.
	Self *PlayedCard
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
	// and is destroyed). DiscountPerRunechant cards read this to compute effective cost.
	Runechants int
	// DelayedRunechants are tokens that skip this turn entirely and go to next turn's carryover.
	// DelayRunechants adds here; same-turn attacks don't consume them and discount checks don't
	// see them. playSequence folds Runechants + DelayedRunechants into LeftoverRunechants.
	DelayedRunechants int
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

// LikelyToHit reports whether dealing n damage is likely to get through an opponent's blocks.
// A typical FaB card is worth ~3 points, so blocking 1/4/7 with a pitch or block card over-pays;
// the opponent would rather eat the damage. Multiples of 3 are the easy-to-block amounts.
// Used by fragile-aura cards (Arcane Cussing, Bloodspill Invocation) to decide whether a
// same-turn attack will actually land and pop the aura for its pay-off.
func LikelyToHit(n int) bool {
	return n == 1 || n == 4 || n == 7
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

// DelayRunechants adds n Runechant tokens that skip this turn entirely — they go to next turn's
// carryover without being available to same-turn attacks or DiscountPerRunechant checks. Used
// by cards whose text fires at the start of a future turn (e.g. Blessing of Occult's "at start
// of your turn, create N Runechant tokens"). Returns n; each token is credited as +1 damage at
// creation.
func (s *TurnState) DelayRunechants(n int) int {
	if n > 0 {
		s.AuraCreated = true
		s.DelayedRunechants += n
	}
	return n
}

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static
// profile plus a Play hook for on-play logic.
type Card interface {
	// ID returns the card's canonical registry identifier. Stable within a build. Lets callers
	// key maps / slices on cards without string-hashing Name().
	ID() ID
	Name() string
	Cost() int
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

// DiscountPerRunechant is optionally implemented by cards whose printed cost is reduced by 1
// per Runechant in play (e.g. Amplify the Arknight, Reduce to Runechant, Rune Flash).
// PrintedCost returns the undiscounted cost; the solver computes the effective per-play cost as
// max(0, PrintedCost() - TurnState.Runechants) at play time.
//
// Cost() on these cards returns 0 so the partition-level affordability check treats them as
// their fully-discounted minimum. The permutation pipeline enforces the actual per-play cost
// against the running resource pool.
type DiscountPerRunechant interface {
	PrintedCost() int
}

// NotSilverAgeLegal is an optional marker. Cards that implement it signal they're banned in the
// Silver Age format and must be excluded from format-restricted deck pools. Source of truth is
// data_sources/silver_age_banlist.txt — keep the two in sync.
type NotSilverAgeLegal interface {
	NotSilverAgeLegal()
}

