package gameengine

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// GameEngine is the per-turn shared state threaded through every Card.Play alongside the
// per-card CardState wrapper. Cards mutate state directly through the public methods —
// moving cards between hand / deck / graveyard / banish, registering triggers, creating
// runechants — and the sim copies the winning permutation's end-of-chain state into
// next-turn state via Snapshot.
//
// Persistent fields (hand, deck, arsenal, graveyard, banished, auras, items) carry across
// turns when the sim adopts the winner's snapshot. Transient fields (cardsPlayed, pitched,
// incomingDamage, etc.) are seeded by Reset at the top of each permutation.
//
// Every field is unexported; the public method surface is the only way to read or write
// state. v2/card.GameEngine is a narrow subset of these methods — what cards see — and is
// satisfied structurally. Sim imports this package directly and gets the rich API (Reset,
// Snapshot, Fire*, SetHero) on top.
type GameEngine struct {
	hero Hero

	hand      []card.Card
	deck      *deck.Deck // engine-owned scratch; refilled on Reset via CopyFrom
	arsenal   card.Card
	graveyard []card.Card
	banished  []card.Card
	auras     []Aura
	triggers  []Trigger
	items     []Item

	cardsPlayed    []card.Card
	cardsRemaining []*card.CardState
	pitched        []card.Card
	defenders      []card.Card

	logger               *turnlogger.TurnLogger
	triggeringCard       card.Card
	attackReactionTarget *card.CardState

	actionPoints         int
	value                int
	cardsDrawn           int
	incomingDamage       int
	arcaneIncomingDamage int
	blockTotal           int
	currentAuraIdx       int

	cardBanished          bool
	arcaneDamageDealt     bool
	opponentMarked        bool
	auraCreated           bool
	nonAttackActionPlayed bool
	currentAuraDestroyed  bool
	currentStepRerouted   bool
	cacheable             bool
}

// Hero returns the active hero — the value SetHero last installed (or the value Spec.Hero
// carried when the engine was built via NewFromSpec). Nil before any hero is set.
func (g *GameEngine) Hero() Hero { return g.hero }

// SetHero installs h as the active hero. Production code calls this at the top of each Best
// pass so card.Universal type-folding and hero ability hooks pick up the right hero; tests
// either supply Spec.Hero up-front or call this between cases.
func (g *GameEngine) SetHero(h Hero) { g.hero = h }

// === Card-facing methods (v2/card.GameEngine surface) ===

// IsCacheable reports whether the chain so far has not depended on hidden state — no card
// in this chain has read or mutated deck / graveyard via an accessor. A hand-eval cache
// stores results only when this is true at chain end.
func (g *GameEngine) IsCacheable() bool { return g.cacheable }

// AttackReactionTarget returns the buff target for the currently-resolving AR, or nil when
// no AR is resolving.
func (g *GameEngine) AttackReactionTarget() *card.CardState { return g.attackReactionTarget }

// SetAttackReactionTarget installs the buff target for the AR resolving next. The chain
// runner sets it around AR.Play and clears it after.
func (g *GameEngine) SetAttackReactionTarget(cs *card.CardState) { g.attackReactionTarget = cs }

// ActionPoints returns the chain runner's running AP pool.
func (g *GameEngine) ActionPoints() int { return g.actionPoints }

// AddActionPoints credits n AP to the running pool. Negative n is allowed.
func (g *GameEngine) AddActionPoints(n int) { g.actionPoints += n }

// ArcaneDamageDealt reports whether any arcane damage source has fired this turn.
func (g *GameEngine) ArcaneDamageDealt() bool { return g.arcaneDamageDealt }

// ArcaneIncomingDamage returns the opponent's arcane damage this turn.
func (g *GameEngine) ArcaneIncomingDamage() int { return g.arcaneIncomingDamage }

// HeroWantsLowerHealth reports whether the current hero opts into the LowerHealthWanter
// marker. Returns false when no hero is set.
func (g *GameEngine) HeroWantsLowerHealth() bool {
	if g.hero == nil {
		return false
	}
	_, ok := g.hero.(LowerHealthWanter)
	return ok
}

// CurrentHeroClass returns the active hero's primary class so "if you are a <class>"
// riders can gate on it. Zero when no hero is set.
func (g *GameEngine) CurrentHeroClass() card.CardType {
	if g.hero == nil {
		return 0
	}
	return g.hero.Class()
}

// AuraCreated reports whether a card or ability has created an aura this turn.
func (g *GameEngine) AuraCreated() bool { return g.auraCreated }

// BlockTotal returns the partition's uncapped defense sum.
func (g *GameEngine) BlockTotal() int { return g.blockTotal }

// CardBanished reports whether any card has been banished this turn.
func (g *GameEngine) CardBanished() bool { return g.cardBanished }

// CardsPlayed returns the sequence of cards played (as attacks) this turn.
func (g *GameEngine) CardsPlayed() []card.Card { return g.cardsPlayed }

// SetCardsPlayed replaces the cards-played slice — used by Moon Wish's transient pre-append
// + pop around its go-again Sun Kiss invocation so the synergy fires.
func (g *GameEngine) SetCardsPlayed(cs []card.Card) { g.cardsPlayed = cs }

// CardsRemaining returns the cards scheduled after the current chain step.
func (g *GameEngine) CardsRemaining() []*card.CardState { return g.cardsRemaining }

// SetCardsRemaining replaces the look-ahead queue — used by tests that seed a partial
// chain for predicate evaluation.
func (g *GameEngine) SetCardsRemaining(cs []*card.CardState) { g.cardsRemaining = cs }

// Defenders returns the partition's defender slice (DRs + plain blocks).
func (g *GameEngine) Defenders() []card.Card { return g.defenders }

// SetDefenders replaces the defender slice — used by the chain runner's plain-block phase
// after the DR pass updates the running defender set.
func (g *GameEngine) SetDefenders(d []card.Card) { g.defenders = d }

// IncomingDamage returns the opponent damage left to allocate this turn.
func (g *GameEngine) IncomingDamage() int { return g.incomingDamage }

// SetIncomingDamage replaces the running incoming-damage tally — used by the per-DR
// reseed inside the defense phase.
func (g *GameEngine) SetIncomingDamage(n int) { g.incomingDamage = n }

// NonAttackActionPlayed reports whether any non-attack action has resolved this turn.
func (g *GameEngine) NonAttackActionPlayed() bool { return g.nonAttackActionPlayed }

// SetNonAttackActionPlayed flips the bookkeeping flag. The chain runner sets it after each
// non-attack action card resolves.
func (g *GameEngine) SetNonAttackActionPlayed(v bool) { g.nonAttackActionPlayed = v }

// OpponentMarked reports whether the opposing hero currently carries the Mark token.
func (g *GameEngine) OpponentMarked() bool { return g.opponentMarked }

// MarkOpponent puts the Mark token on the opposing hero. The next attack action card /
// weapon swing clears it.
func (g *GameEngine) MarkOpponent() { g.opponentMarked = true }

// ClearOpponentMarked strips the Mark token. The chain runner calls it after each attack
// action card / weapon swing resolves.
func (g *GameEngine) ClearOpponentMarked() { g.opponentMarked = false }

// OpponentDiscard credits n cards' worth of damage-equivalent value for forcing the
// opponent to discard. Returns the credited value for log attribution.
func (g *GameEngine) OpponentDiscard(n int) int {
	v := n * DiscardValue
	g.AddValue(v)
	return v
}

// LikelyToHit reports whether self's attack is likely to land past the opponent's blocks.
func (g *GameEngine) LikelyToHit(self *card.CardState) bool { return LikelyToHit(self) }

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit.
func (g *GameEngine) LikelyDamageHits(n int, dominate bool) bool {
	return LikelyDamageHits(n, dominate)
}

// Pitched returns the cards pitched this turn for resources.
func (g *GameEngine) Pitched() []card.Card { return g.pitched }

// TriggeringCard returns the card whose Play caused the currently-firing aura
// attack-action trigger, or nil outside of a trigger fire.
func (g *GameEngine) TriggeringCard() card.Card { return g.triggeringCard }

// SetTriggeringCard replaces the triggering-card slot. Used by tests that drive a trigger
// handler directly; production threads it through the fire loop.
func (g *GameEngine) SetTriggeringCard(c card.Card) { g.triggeringCard = c }

// Value returns the running damage-equivalent total for this chain.
func (g *GameEngine) Value() int { return g.value }

// AddValue credits n to the chain's value, clamped at 0. Negative n is allowed for "this
// card gives the opponent value" effects (Test of Strength's clash-loss).
func (g *GameEngine) AddValue(n int) { g.value += n }

// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. ARs use this to
// fold their +{p} buff into the buffed attack's display delta. No-op when the logger is
// nil (find-best pass) or when no chain-step entry exists yet.
func (g *GameEngine) AmendLastChainStepN(n int) {
	g.logger.AmendLastChainStepN(n)
}

// === Zone accessors — cards see these; each flips cacheable to false ===

// Deck returns the chain-runner deck for read-only inspection and flips IsCacheable to
// false. Card handlers should not mutate the returned *deck.Deck directly; route mutations
// through PopDeckTop / PrependToDeck / Opt / TutorFromDeck / RecycleToDeckBottom.
func (g *GameEngine) Deck() *deck.Deck {
	g.cacheable = false
	return g.deck
}

// PeekTopN returns the top n cards of the deck (top first) without removing them and flips
// IsCacheable to false. Returns fewer cards when the deck has < n.
func (g *GameEngine) PeekTopN(n int) []card.Card {
	g.cacheable = false
	top := g.deck.PeekTopN(n)
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
// the returned slice; use AppendHand / PopHandAt for mutations.
func (g *GameEngine) Hand() []card.Card {
	g.cacheable = false
	return g.hand
}

// AppendHand appends c to the hand, flipping IsCacheable to false.
func (g *GameEngine) AppendHand(c card.Card) {
	g.cacheable = false
	g.hand = append(g.hand, c)
}

// PopHandAt removes and returns the card at index i, flipping IsCacheable to false.
func (g *GameEngine) PopHandAt(i int) card.Card {
	g.cacheable = false
	c := g.hand[i]
	g.hand = append(g.hand[:i], g.hand[i+1:]...)
	return c
}

// Graveyard returns the live graveyard slice and flips IsCacheable to false.
func (g *GameEngine) Graveyard() []card.Card {
	g.cacheable = false
	return g.graveyard
}

// PopDeckTop removes the top card of the deck and returns it. Returns (nil, false) when
// the deck is empty. Flips IsCacheable to false.
func (g *GameEngine) PopDeckTop() (card.Card, bool) {
	g.cacheable = false
	if g.deck.Size() == 0 {
		return nil, false
	}
	return g.deck.Draw(1)[0].(card.Card), true
}

// PeekDeck returns the top card of the deck without removing it. Returns (nil, false) on
// an empty deck. Flips IsCacheable to false.
func (g *GameEngine) PeekDeck() (card.Card, bool) {
	g.cacheable = false
	top := g.deck.PeekTop()
	if top == nil {
		return nil, false
	}
	return top.(card.Card), true
}

// PrependToDeck inserts c at the top of the deck. Flips IsCacheable to false.
func (g *GameEngine) PrependToDeck(c card.Card) {
	g.cacheable = false
	g.deck.PutTop([]deck.Card{c})
}

// RecycleToDeckBottom appends self.Card to the bottom of the deck and flags the chain
// dispatcher to skip the usual non-persistent → graveyard append. Models the FaB clause
// "put this on the bottom of its owner's deck" (Relentless Pursuit). Flips IsCacheable.
func (g *GameEngine) RecycleToDeckBottom(self *card.CardState) {
	g.cacheable = false
	g.deck.PutBottom([]deck.Card{self.Card})
	g.currentStepRerouted = true
}

// CurrentStepRerouted is the per-step flag the chain dispatcher reads after Play to decide
// whether to skip the "non-persistent → graveyard" append.
func (g *GameEngine) CurrentStepRerouted() bool { return g.currentStepRerouted }

// SetCurrentStepRerouted resets the per-step flag. The chain dispatcher clears it before
// each step.
func (g *GameEngine) SetCurrentStepRerouted(v bool) { g.currentStepRerouted = v }

// Opt resolves the FaB "Opt N" keyword: pops up to n cards from the top of the deck and
// hands them to the current hero's Opt heuristic. The handler returns a (top, bottom)
// split; the top list goes back on top of the deck (in returned order) and the bottom
// list appends to the bottom (in returned order). n is clamped to the current deck
// length. Always flips IsCacheable to false.
//
// Emits a log entry naming the revealed cards and the split when the handler ran.
//
// Panics if the handler's combined output isn't exactly the input multiset.
func (g *GameEngine) Opt(l card.Logger, n int) {
	g.cacheable = false
	if n <= 0 || g.deck.Size() == 0 {
		return
	}
	if n > g.deck.Size() {
		n = g.deck.Size()
	}
	drawn := g.deck.Draw(n)
	cards := make([]card.Card, len(drawn))
	for i, c := range drawn {
		cards[i] = c.(card.Card)
	}

	var top, bottom []card.Card
	if g.hero == nil {
		// Default passthrough: every revealed card goes back on top in input order.
		top = cards
	} else {
		top, bottom = g.hero.Opt(cards)
	}
	panicIfOptViolatesMultiset(cards, top, bottom)

	deckTop := make([]deck.Card, len(top))
	for i, c := range top {
		deckTop[i] = c
	}
	g.deck.PutTop(deckTop)
	deckBottom := make([]deck.Card, len(bottom))
	for i, c := range bottom {
		deckBottom[i] = c
	}
	g.deck.PutBottom(deckBottom)

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

// formatCardList renders cs as "[name1, name2, ...]" using DisplayName, or "[]" when empty.
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

// panicIfOptViolatesMultiset enforces GameEngine.Opt's contract that the hero handler's
// combined (top, bottom) output is exactly the input multiset — a permutation of the
// input cards, no additions or removals.
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
func (g *GameEngine) TutorFromDeck(score func(card.Card) int) (card.Card, bool) {
	g.cacheable = false
	got, ok := g.deck.Tutor(func(c deck.Card) int { return score(c.(card.Card)) })
	if !ok {
		return nil, false
	}
	return got.(card.Card), true
}

// BanishFromGraveyard removes the first graveyard card matching pred, appends it to the
// banished zone, and returns it. Returns (nil, false) when no card matches. Flips
// IsCacheable to false. Sets CardBanished so this-turn-banish riders fire correctly.
func (g *GameEngine) BanishFromGraveyard(pred func(card.Card) bool) (card.Card, bool) {
	g.cacheable = false
	for i, c := range g.graveyard {
		if !pred(c) {
			continue
		}
		g.banished = append(g.banished, c)
		g.cardBanished = true
		g.graveyard = append(g.graveyard[:i], g.graveyard[i+1:]...)
		return c, true
	}
	return nil, false
}

// Banished returns the slice of cards in the banished zone, top-to-bottom. Read-only —
// mutate via BanishFromGraveyard. Includes prior-turn entries; "did anything banish THIS
// turn" readers must use CardBanished instead.
func (g *GameEngine) Banished() []card.Card {
	return g.banished
}

// RecycleFromGraveyardToTop removes the first graveyard card matching pred, prepends it to
// the deck, and returns it. Returns (nil, false) when no card matches. Flips IsCacheable
// to false. The deck mutation IS the model — the recycled card resurfaces in next turn's
// deal naturally, so callers don't credit Value here.
func (g *GameEngine) RecycleFromGraveyardToTop(pred func(card.Card) bool) (card.Card, bool) {
	return g.recycleFromGraveyard(pred, true)
}

// RecycleFromGraveyardToBottom is the bottom-of-deck variant of
// RecycleFromGraveyardToTop.
func (g *GameEngine) RecycleFromGraveyardToBottom(pred func(card.Card) bool) (card.Card, bool) {
	return g.recycleFromGraveyard(pred, false)
}

func (g *GameEngine) recycleFromGraveyard(pred func(card.Card) bool, toTop bool) (card.Card, bool) {
	g.cacheable = false
	for i, c := range g.graveyard {
		if !pred(c) {
			continue
		}
		g.graveyard = append(g.graveyard[:i], g.graveyard[i+1:]...)
		if toTop {
			g.deck.PutTop([]deck.Card{c})
		} else {
			g.deck.PutBottom([]deck.Card{c})
		}
		return c, true
	}
	return nil, false
}

// AddToGraveyard appends c to graveyard so later-resolving cards see it. Used by cards
// running a mini-dispatcher inline (Moon Wish's go-again Sun Kiss play). Flips
// IsCacheable to false so the convention "every public accessor that touches deck /
// graveyard flips cacheable" stays universal.
func (g *GameEngine) AddToGraveyard(c card.Card) {
	g.cacheable = false
	g.graveyard = append(g.graveyard, c)
}

// AppendGraveyard appends c to graveyard without flipping IsCacheable. Framework-internal:
// the chain dispatcher's "non-persistent → graveyard" rule and Aura.OnDestroy use it so
// engine bookkeeping doesn't poison the cacheable bit.
func (g *GameEngine) AppendGraveyard(c card.Card) {
	g.graveyard = append(g.graveyard, c)
}

// AppendHandRaw appends c to the hand without flipping IsCacheable. Framework-internal:
// the chain runner uses it to seed chain attackers + pitched cards into hand at the start
// of each permutation so cards' Hand() reads see the upcoming bag.
func (g *GameEngine) AppendHandRaw(c card.Card) {
	g.hand = append(g.hand, c)
}

// RemoveFromHand removes the first matching card from the hand without flipping
// IsCacheable. Returns true if a card was removed. Framework-internal: the chain runner
// pops the playing card and freshly-popped pitch cards out of hand as each chain step
// resolves.
func (g *GameEngine) RemoveFromHand(c card.Card) bool {
	for i := range g.hand {
		if g.hand[i] == c {
			g.hand = append(g.hand[:i], g.hand[i+1:]...)
			return true
		}
	}
	return false
}

// HandRaw returns the hand slice without flipping IsCacheable. Framework-internal:
// the chain runner reads it for permutation bookkeeping. Cards must continue to use
// Hand().
func (g *GameEngine) HandRaw() []card.Card { return g.hand }

// SetHand replaces the hand slice without flipping IsCacheable. Framework / test seeding
// — tests that need to seed a hand without going through Spec / Reset use this.
func (g *GameEngine) SetHand(h []card.Card) { g.hand = h }

// GraveyardRaw returns the graveyard slice without flipping IsCacheable.
// Framework-internal: snapshot / display code reads it for end-of-chain assembly.
func (g *GameEngine) GraveyardRaw() []card.Card { return g.graveyard }

// DeckRaw returns the engine's scratch *deck.Deck without flipping IsCacheable.
// Framework-internal.
func (g *GameEngine) DeckRaw() *deck.Deck { return g.deck }

// DrawOne models a mid-turn draw: pop the top of the deck and append it to Hand. No-op on
// an empty deck. Bumps CardsDrawn so the partition tiebreaker can prefer chains with more
// draws. Inherits the IsCacheable flip via PopDeckTop.
func (g *GameEngine) DrawOne() {
	c, ok := g.PopDeckTop()
	if !ok {
		return
	}
	g.hand = append(g.hand, c)
	g.cardsDrawn++
}

// CardsDrawn returns the count of mid-chain card draws this turn.
func (g *GameEngine) CardsDrawn() int { return g.cardsDrawn }

// HasPlayedType reports whether any card played this turn has the given type. Universal
// cards' Types() folds the active hero's class through g.
func (g *GameEngine) HasPlayedType(t card.CardType) bool {
	for _, c := range g.cardsPlayed {
		if c.Types(g).Has(t) {
			return true
		}
	}
	return false
}

// Clash models a clash (rule 8.5.45): we and the opponent reveal the top card of our decks
// and the higher {p} wins. We model from our side only — our deck's top is read via
// PeekDeck; the opponent's top is approximated as 5-power. On a win (our top ≥ 6), win
// fires; on a loss (our top ≤ 4), lose fires; ties (top == 5) and empty deck fire neither.
// PeekDeck flips IsCacheable to false.
func (g *GameEngine) Clash(win, lose func()) {
	top, ok := g.PeekDeck()
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

// === Aura / Trigger / token accessors used by cards ===

// Auras returns the live aura set. Read-only.
func (g *GameEngine) Auras() []Aura { return g.auras }

// SetAuras replaces the aura set wholesale. Used by tests that seed a prior-turn aura
// carryover; production code uses CreateXxxAura, which also flips AuraCreated.
func (g *GameEngine) SetAuras(a []Aura) { g.auras = a }

// AuraCount returns the count of live auras. Used by gates like Yinti Yanti's "while you
// control an aura" rider.
func (g *GameEngine) AuraCount() int { return len(g.auras) }

// Triggers returns the one-shot trigger queue. Read-only.
func (g *GameEngine) Triggers() []Trigger { return g.triggers }

// Items returns the live item set. Read-only.
func (g *GameEngine) Items() []Item { return g.items }

// LogEntries returns the per-event chain trace accumulated by the Log family.
func (g *GameEngine) LogEntries() []turnlogger.LogEntry { return g.logger.Entries() }

// Logger returns the chain runner's currently-active log sink. Returns the underlying
// *turnlogger.TurnLogger so sim can rebind buffers / pass it through Reset; the type
// satisfies card.Logger structurally so cards can still treat it as the cards-facing
// interface. Nil during the find-best pass.
func (g *GameEngine) Logger() *turnlogger.TurnLogger { return g.logger }

// SetLogger replaces the per-permutation log sink. Used by the per-DR seed inside the
// defense phase, which constructs a fresh recording / find-best logger to capture the
// DR's chain step. Production wires this through Reset.
func (g *GameEngine) SetLogger(l *turnlogger.TurnLogger) { g.logger = l }
