package card

// Per-turn shared context threaded through Card.Play: the TurnState type and its helper
// methods for mid-turn mutations (draws, graveyard / banish moves, aura creation, trigger
// registration) that more than one card reaches for. Cards read the fields to decide effects;
// the sim owns the top-level bookkeeping (appending to CardsPlayed, firing triggers, etc.).

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
	// Held is the partition's Held-role cards at start of the chain — the hand cards the
	// solver assigned no Pitch / Attack / Defend role. Read-only by Play unless the card
	// implements an alt-cost "use a Held card" effect, in which case Play pops the consumed
	// card off Held and appends it to HeldConsumed so the post-chain accounting
	// (recycleCardStates, arsenal-promotion candidate counts) skips it.
	Held []Card
	// HeldConsumed records cards moved out of Held by alt-cost effects mid-chain. The
	// deck-loop accounting compares BestLine's Held-role cards against this list and
	// suppresses the nextHeld carry for any match. Cards listed here are also inserted at
	// the top of the next-turn deck buffer (the "rather than pay" rule), so the same card
	// commonly reappears in the next turn's hand or feeds a same-turn DrawOne when a tutor
	// fires.
	HeldConsumed []Card
	// DeckRemoved records cards taken out of the deck this turn by any means — DrawOne,
	// tutor effects (Moon Wish's Sun Kiss search), or future deck-search riders. The
	// deck-loop's applyTurnResult patches the underlying deck buffer to actually remove
	// each listed card so it can't be drawn again on a later turn. Without this list the
	// buf would still hold the tutored card at its original position, and a duplicate would
	// surface once the head pointer reached that slot.
	DeckRemoved []Card
}

// DrawOne models a mid-turn draw: advance the deck by one card and append it to Drawn. No-op
// on an empty deck. Every draw-rider card routes through this helper. Also appends to
// DeckRemoved so applyTurnResult patches the same card out of the underlying deck buffer.
func (s *TurnState) DrawOne() {
	if len(s.Deck) == 0 {
		return
	}
	c := s.Deck[0]
	s.Drawn = append(s.Drawn, c)
	s.DeckRemoved = append(s.DeckRemoved, c)
	s.Deck = s.Deck[1:]
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
