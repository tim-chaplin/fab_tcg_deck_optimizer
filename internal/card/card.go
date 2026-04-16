// Package card defines the Card interface used by the simulator and the basic / test card
// implementations.
package card

// CardType is an enumerated card-type descriptor. Each constant corresponds to a single keyword
// from a FaB card's type line (e.g. "Runeblade", "Action", "Attack").
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

// TypeSet is a bitfield of CardType values. It replaces map[string]bool for card type checks,
// eliminating string hashing on every lookup.
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

// PlayedCard wraps a Card with per-turn mutable flags that other cards' effects can toggle during
// the chain. Instances are created by the solver at the start of each attack chain and live only
// for that chain. Effects that grant keywords to "the next X" scan TurnState.CardsRemaining and
// flip flags on the matching entry — no special-cased fields on TurnState required.
type PlayedCard struct {
	Card Card
	// GrantedGoAgain is set by a prior card's effect to give this specific card Go again even if
	// its printed text doesn't (e.g. Mauvrion Skies targeting the next Runeblade attack action).
	// The solver's chain-legality check ORs this with Card.GoAgain().
	GrantedGoAgain bool
}

// EffectiveGoAgain reports whether this card has Go again for the current turn, whether from its
// printed text or a grant from a prior card's effect.
func (p *PlayedCard) EffectiveGoAgain() bool {
	return p.Card.GoAgain() || p.GrantedGoAgain
}

// TurnState is the context passed to Card.Play. Cards read it to decide what effects to apply;
// the solver appends each played card to CardsPlayed after its Play method returns, so later cards
// this turn can see what was played before them.
type TurnState struct {
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order. Populated by the
	// solver, not by Play itself.
	CardsPlayed []Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant tokens,
	// which are auras). Effects that check "if you've played or created an aura this turn" should
	// OR this with CardsPlayed containing an Aura-typed card.
	AuraCreated bool
	// CardsRemaining is the cards that will be played after the current one in the turn's ordering.
	// Populated by the solver before each Play so an effect can peek forward (e.g. Condemn to
	// Slaughter buffing the "next Runeblade attack") OR grant keywords to a later card by flipping
	// flags on its PlayedCard entry (e.g. Mauvrion Skies granting Go again).
	CardsRemaining []*PlayedCard
	// Pitched is the set of cards pitched this turn to generate resources. Populated by the solver
	// before any Play is called. Effects that check "if an attack card was pitched" scan this list.
	Pitched []Card
	// Self is the PlayedCard wrapper for the card currently being played. Effects that conditionally
	// grant the played card itself Go again (e.g. Runerager Swarm: "If you've played or created an
	// aura this turn, this gets go again") flip Self.GrantedGoAgain. The solver populates this
	// before each Play and consults EffectiveGoAgain after.
	Self *PlayedCard
	// Overpower is set when an attack with the Overpower keyword is being played. Not yet consumed by
	// the solver — blocked damage should eventually be forwarded to the hero when Overpower is true.
	Overpower bool
	// Deck is the cards remaining in the deck (excluding the current hand), in top-of-deck order.
	// Effects that reveal or draw the top card (e.g. Sigil of the Arknight) inspect this. Nil when
	// unknown / not provided. Implementations must not mutate it.
	Deck []Card
}

// Hero is the minimal hero profile card effects need. It's intentionally narrower than
// hero.Hero to avoid an import cycle. Package simstate holds the active hero for the run.
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

// CreateRunechants models the Runeblade mechanic of creating n Runechant token auras. Sets
// AuraCreated so effects that key on "aura created this turn" see it, and returns the damage the
// tokens will contribute this chain — currently 1 per token, since the solver assumes a following
// attack always arrives and triggers them. Callers add the returned value to whatever damage
// they're reporting back to the solver.
//
// Centralising creation in one call site makes it the only place to touch when we later track
// runechant counts for discount-per-token cards (e.g. Malefic Incantation's cost reduction).
func (s *TurnState) CreateRunechants(n int) int {
	if n > 0 {
		s.AuraCreated = true
	}
	return n
}

// CreateRunechant is shorthand for CreateRunechants(1) for the common single-token case.
func (s *TurnState) CreateRunechant() int {
	return s.CreateRunechants(1)
}

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static profile
// plus a Play hook for on-play logic.
type Card interface {
	// ID returns the card's canonical registry identifier. Stable within a build. Lets callers
	// key maps / slices on cards without string-hashing Name().
	ID() ID
	Name() string
	Cost() int
	Pitch() int
	// Attack is the card's base (printed) attack value. Conditional bonuses belong in Play, not here.
	Attack() int
	Defense() int
	// Types returns the card's type-line descriptors as a TypeSet bitfield, e.g.
	// NewTypeSet(TypeRuneblade, TypeAction, TypeAttack).
	Types() TypeSet
	// GoAgain reports whether playing this card grants an additional action point this turn. Cards
	// printed with "Go again" return true.
	GoAgain() bool
	// Play is called when the card is played — as an attack or as a defense reaction. It returns
	// damage dealt to the opposing hero (which may differ from Attack() after conditional bonuses)
	// and may read state to decide effects. When called on a defense reaction, the returned damage
	// is added to the turn's dealt total uncapped (the incoming-damage prevention cap applies only
	// to Defense()).
	Play(s *TurnState) int
}

// NoMemo is an optional marker. Cards that implement it signal that hands containing them must
// not be served from or written to the hand-evaluation memo — typically because the card's Play
// output depends on context (e.g. remaining deck composition) that the memo key doesn't capture.
type NoMemo interface {
	NoMemo()
}

