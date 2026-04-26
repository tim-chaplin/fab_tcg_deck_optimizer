package card

// Per-turn shared context threaded through Card.Play. Cards mutate state directly — moving
// cards between Hand / Deck / Graveyard / Banish, registering triggers, creating runechants
// — and the sim copies the winning permutation's final state into next-turn state. There's
// no diff-signal indirection: a card that wants to draw appends to s.Hand and pops from
// s.Deck, full stop.
//
// Persistent fields (Hand, Deck, Arsenal, Graveyard, Banish, Runechants, AuraTriggers)
// carry across turns when the sim adopts the winner's snapshot. Transient fields
// (CardsPlayed, Pitched, IncomingDamage, etc.) are seeded by the sim per chain-step and
// reset at the turn boundary.

// TurnState is the shared turn-level context passed to Card.Play alongside the per-card
// CardState wrapper.
type TurnState struct {
	// Hand is the cards currently in hand. Starts as the dealt hand minus pitched / attacker
	// / defender cards (those have been routed by the partition). Cards that draw or tutor
	// append to Hand; alt-cost effects pop from Hand. Whatever's in Hand at end of chain
	// becomes next turn's Held cards.
	Hand []Card
	// Deck is the deck top-to-bottom. Cards mutate freely: DrawOne pops Deck[0]; tutor
	// removes a specific card; alt cost prepends to Deck. Whatever's in Deck at end of
	// chain becomes next turn's deck.
	Deck []Card
	// Arsenal is the arsenal slot's contents at this point in the chain — the arsenal-in
	// card at start of turn, nil after it plays / defends, refilled post-chain by the
	// arsenal-promotion step. Cards that read "from arsenal" use CardState.FromArsenal,
	// not this field.
	Arsenal Card
	// Graveyard is cards that have entered the graveyard this turn — every card played or
	// blocked lands here after resolving. Pitched cards do not (they recycle to deck
	// bottom). Cards that destroy themselves mid-chain route through AddToGraveyard.
	Graveyard []Card
	// Banish holds cards moved into the banished zone this turn (e.g. an aura-banish-for-
	// arcane rider).
	Banish []Card
	// Runechants is the live count of Runechant aura tokens in play. Carries across turns.
	// CreateRunechants increments it; the attack pipeline consumes the running total on each
	// attack / weapon swing.
	Runechants int
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
	// prevented + every aura-token / hero-trigger credit. The dispatcher calls RecordValue
	// after each Play / hero / aura / ephemeral / weapon return; the solver compares
	// permutations on this field. Reset by the sim per permutation.
	Value int
	// Log is reserved for an upcoming per-line trace of the chain — the dispatcher does not
	// yet write to it and FormatBestTurn does not yet read from it. Kept on the struct so the
	// follow-up wiring is a one-place change. Reset per permutation.
	Log []string
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
	// IncomingDamage is the opponent damage this turn (the value passed to hand.Best).
	// Constant across every partition the solver enumerates for this hand.
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
}

// DrawOne models a mid-turn draw: pop the top of Deck and append it to Hand. No-op on an
// empty deck. Every draw-rider card routes through this helper.
func (s *TurnState) DrawOne() {
	if len(s.Deck) == 0 {
		return
	}
	c := s.Deck[0]
	s.Deck = s.Deck[1:]
	s.Hand = append(s.Hand, c)
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
// behind "if you've played or created an aura this turn" riders.
func (s *TurnState) HasAuraInPlay() bool {
	return s.AuraCreated || s.HasPlayedType(TypeAura)
}

// ClashValue returns the net damage-equivalent of a clash (see comprehensive rules 8.5.45):
// we and the opponent reveal the top card of our decks and the higher {p} wins. We model
// from our side only — our deck's top card is read from s.Deck; the opponent's top is
// approximated as 5-power. So our {p} of 6-7 wins (credit +bonus), 5 ties (credit 0), and
// anything below 5 loses (credit -bonus). Returns 0 when s.Deck is empty.
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

// RecordValue bumps s.Value by n, clamping at 0 (FaB damage / prevention can't drive the
// running total negative). Negative n is a no-op. The dispatcher calls this after each
// Play / hero trigger / aura trigger / weapon swing / defense block so s.Value is the
// authoritative running total for the permutation. Cards don't call RecordValue themselves —
// they return the damage-equivalent from Play and let the dispatcher record it.
func (s *TurnState) RecordValue(n int) {
	if n <= 0 {
		return
	}
	s.Value += n
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

// DealArcaneDamage flips ArcaneDamageDealt so same-turn triggers reading "if you've dealt
// arcane damage this turn" fire, and returns n so callers can fold the arcane damage into
// their Play return in one expression.
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

// AddAuraTrigger is the Play-side combo every Action - Aura card reaches for: flip
// AuraCreated so same-turn "if you've played or created an aura" riders see the entry, and
// append t to s.AuraTriggers so the sim fires it on its matching Type condition.
func (s *TurnState) AddAuraTrigger(t AuraTrigger) {
	s.AuraCreated = true
	s.AuraTriggers = append(s.AuraTriggers, t)
}

// AddEphemeralAttackTrigger registers a same-turn, fire-once "next attack" trigger. The sim
// stamps t.SourceIndex after the registering card's Play returns. Fires on the next
// matching attack action's resolution; fizzles silently at end of turn if no match.
func (s *TurnState) AddEphemeralAttackTrigger(t EphemeralAttackTrigger) {
	s.EphemeralAttackTriggers = append(s.EphemeralAttackTriggers, t)
}
