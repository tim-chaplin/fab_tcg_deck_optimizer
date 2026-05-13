package gameengine

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// GameEngine is the rules-engine wrapper around a *GameState. Cards play against this
// type via the v2/card.GameEngine interface; the engine's method surface mixes
// (a) cacheable-aware accessors that flip g.state.cacheable as a side effect of touching
// hidden state, and (b) rules-orchestration methods (Fire*, ResolveChainStep, Opt,
// Clash, DealArcaneDamage, token economy) that apply the FaB rule set on top of the raw
// state mutations the state itself exposes.
//
// GameState owns the data and the cheap pure accessors; GameEngine owns the rules. The
// split lets internal machinery (TurnSummary.State, sim's per-permutation copy, the
// find-best winner) pass around a *GameState pointer without dragging the whole engine
// API along when all it needs to do is read or copy raw state.
type GameEngine struct {
	state *GameState
}

// State returns the underlying *GameState. Sim uses this to pull the raw state out of an
// engine — e.g. capturing the winning permutation's state in TurnSummary.State.
func (g *GameEngine) State() *GameState { return g.state }

// === Hero + cacheable: read-only sugar; cards see the rules-aware view ===

func (g *GameEngine) Hero() Hero                 { return g.state.hero }
func (g *GameEngine) SetHero(h Hero)             { g.state.hero = h }
func (g *GameEngine) IsCacheable() bool          { return g.state.cacheable }
func (g *GameEngine) SetCacheable(v bool)        { g.state.cacheable = v }
func (g *GameEngine) HeroWantsLowerHealth() bool { return g.state.HeroWantsLowerHealth() }
func (g *GameEngine) CurrentHeroClass() card.CardType {
	return g.state.CurrentHeroClass()
}

// === Cards-facing zone accessors that flip cacheable ===

// Hand returns the live hand slice and flips IsCacheable to false. Cards must not mutate
// the returned slice; use AppendHand / PopHandAt for mutations.
func (g *GameEngine) Hand() []card.Card {
	g.state.cacheable = false
	return g.state.hand
}

// AppendHand appends c to the hand, flipping IsCacheable to false.
func (g *GameEngine) AppendHand(c card.Card) {
	g.state.cacheable = false
	g.state.hand = append(g.state.hand, c)
}

// PopHandAt removes and returns the card at index i, flipping IsCacheable to false.
func (g *GameEngine) PopHandAt(i int) card.Card {
	g.state.cacheable = false
	c := g.state.hand[i]
	g.state.hand = append(g.state.hand[:i], g.state.hand[i+1:]...)
	return c
}

// Graveyard returns the live graveyard slice and flips IsCacheable to false.
func (g *GameEngine) Graveyard() []card.Card {
	g.state.cacheable = false
	return g.state.graveyard
}

// Deck returns the chain-runner deck for read-only inspection and flips IsCacheable to
// false. Card handlers should not mutate the returned *deck.Deck directly; route
// mutations through PopDeckTop / PrependToDeck / Opt / TutorFromDeck /
// RecycleToDeckBottom.
func (g *GameEngine) Deck() *deck.Deck {
	g.state.cacheable = false
	return g.state.deck
}

// PeekTopN returns the top n cards of the deck (top first) without removing them and
// flips IsCacheable to false. Returns fewer cards when the deck has < n.
func (g *GameEngine) PeekTopN(n int) []card.Card {
	g.state.cacheable = false
	top := g.state.deck.PeekTopN(n)
	if len(top) == 0 {
		return nil
	}
	out := make([]card.Card, len(top))
	for i, c := range top {
		out[i] = c.(card.Card)
	}
	return out
}

// PopDeckTop removes the top card of the deck and returns it. Returns (nil, false) when
// the deck is empty. Flips IsCacheable to false.
func (g *GameEngine) PopDeckTop() (card.Card, bool) {
	g.state.cacheable = false
	if g.state.deck.Size() == 0 {
		return nil, false
	}
	return g.state.deck.Draw(1)[0].(card.Card), true
}

// PeekDeck returns the top card of the deck without removing it. Returns (nil, false) on
// an empty deck. Flips IsCacheable to false.
func (g *GameEngine) PeekDeck() (card.Card, bool) {
	g.state.cacheable = false
	top := g.state.deck.PeekTop()
	if top == nil {
		return nil, false
	}
	return top.(card.Card), true
}

// PrependToDeck inserts c at the top of the deck. Flips IsCacheable to false.
func (g *GameEngine) PrependToDeck(c card.Card) {
	g.state.cacheable = false
	g.state.deck.PutTop([]deck.Card{c})
}

// RecycleToDeckBottom appends self.Card to the bottom of the deck and flags the chain
// dispatcher to skip the usual non-persistent → graveyard append. Models the FaB clause
// "put this on the bottom of its owner's deck" (Relentless Pursuit). Flips IsCacheable.
func (g *GameEngine) RecycleToDeckBottom(self *card.CardState) {
	g.state.cacheable = false
	g.state.deck.PutBottom([]deck.Card{self.Card})
	g.state.currentStepRerouted = true
}

// TutorFromDeck removes and returns the highest-scoring card per score. Returns (nil,
// false) when no card scores > 0 (or the deck is empty). Flips IsCacheable to false.
func (g *GameEngine) TutorFromDeck(score func(card.Card) int) (card.Card, bool) {
	g.state.cacheable = false
	got, ok := g.state.deck.Tutor(func(c deck.Card) int { return score(c.(card.Card)) })
	if !ok {
		return nil, false
	}
	return got.(card.Card), true
}

// BanishFromGraveyard removes the first graveyard card matching pred, appends it to the
// banished zone, and returns it. Returns (nil, false) when no card matches. Flips
// IsCacheable to false. Sets CardBanished so this-turn-banish riders fire correctly.
func (g *GameEngine) BanishFromGraveyard(pred func(card.Card) bool) (card.Card, bool) {
	g.state.cacheable = false
	for i, c := range g.state.graveyard {
		if !pred(c) {
			continue
		}
		g.state.banished = append(g.state.banished, c)
		g.state.cardBanished = true
		g.state.graveyard = append(g.state.graveyard[:i], g.state.graveyard[i+1:]...)
		return c, true
	}
	return nil, false
}

// Banished returns the slice of cards in the banished zone, top-to-bottom. Read-only —
// mutate via BanishFromGraveyard. Includes prior-turn entries; "did anything banish THIS
// turn" readers must use CardBanished instead.
func (g *GameEngine) Banished() []card.Card { return g.state.banished }

// RecycleFromGraveyardToTop / RecycleFromGraveyardToBottom remove the first graveyard
// card matching pred and put it on the top / bottom of the deck. Flip IsCacheable.
func (g *GameEngine) RecycleFromGraveyardToTop(pred func(card.Card) bool) (card.Card, bool) {
	return g.recycleFromGraveyard(pred, true)
}
func (g *GameEngine) RecycleFromGraveyardToBottom(pred func(card.Card) bool) (card.Card, bool) {
	return g.recycleFromGraveyard(pred, false)
}

func (g *GameEngine) recycleFromGraveyard(pred func(card.Card) bool, toTop bool) (card.Card, bool) {
	g.state.cacheable = false
	for i, c := range g.state.graveyard {
		if !pred(c) {
			continue
		}
		g.state.graveyard = append(g.state.graveyard[:i], g.state.graveyard[i+1:]...)
		if toTop {
			g.state.deck.PutTop([]deck.Card{c})
		} else {
			g.state.deck.PutBottom([]deck.Card{c})
		}
		return c, true
	}
	return nil, false
}

// AddToGraveyard appends c to graveyard so later-resolving cards see it. Used by cards
// running a mini-dispatcher inline (Moon Wish's go-again Sun Kiss play). Flips
// IsCacheable to false.
func (g *GameEngine) AddToGraveyard(c card.Card) {
	g.state.cacheable = false
	g.state.graveyard = append(g.state.graveyard, c)
}

// AppendGraveyard appends c to graveyard without flipping IsCacheable. Framework-
// internal: the chain dispatcher's "non-persistent → graveyard" rule and Aura.OnDestroy
// use it so engine bookkeeping doesn't poison the cacheable bit.
func (g *GameEngine) AppendGraveyard(c card.Card) {
	g.state.graveyard = append(g.state.graveyard, c)
}

// === Framework-internal forwarders to *GameState. These exist on the engine so the
// chain runner can mutate state without alternating between an engine pointer (for
// rules-orchestration calls) and a state pointer (for raw access) on every line. Cards
// shouldn't reach for these — they aren't on card.GameEngine.

// HandRaw / SetHand / AppendHandRaw / RemoveFromHand expose hand mutation without the
// IsCacheable flip the cards-facing Hand() / AppendHand() / PopHandAt() impose.
func (g *GameEngine) HandRaw() []card.Card      { return g.state.hand }
func (g *GameEngine) SetHand(h []card.Card)     { g.state.hand = h }
func (g *GameEngine) AppendHandRaw(c card.Card) { g.state.hand = append(g.state.hand, c) }
func (g *GameEngine) RemoveFromHand(c card.Card) bool {
	for i := range g.state.hand {
		if g.state.hand[i] == c {
			g.state.hand = append(g.state.hand[:i], g.state.hand[i+1:]...)
			return true
		}
	}
	return false
}

// DeckRaw / SetDeck / GraveyardRaw / SetGraveyard: raw read / write of the deck and
// graveyard slices without flipping IsCacheable.
func (g *GameEngine) DeckRaw() *deck.Deck         { return g.state.deck }
func (g *GameEngine) SetDeck(d *deck.Deck)        { g.state.deck = d }
func (g *GameEngine) GraveyardRaw() []card.Card   { return g.state.graveyard }
func (g *GameEngine) SetGraveyard(gv []card.Card) { g.state.graveyard = gv }

// SetArsenal / Arsenal: chain runner uses these post-chain to promote a leftover hand
// card into next turn's arsenal.
func (g *GameEngine) Arsenal() card.Card     { return g.state.arsenal }
func (g *GameEngine) SetArsenal(c card.Card) { g.state.arsenal = c }

// SetBanished / SetPitched: per-permutation seeding.
func (g *GameEngine) SetBanished(b []card.Card) { g.state.banished = b }
func (g *GameEngine) SetPitched(p []card.Card)  { g.state.pitched = p }

// Auras / Items / Triggers + their Clear* counterparts: chain runner reads / clears the
// engine's persistent-in-play lists when seeding a fresh permutation.
func (g *GameEngine) Auras() []Aura       { return g.state.auras }
func (g *GameEngine) Items() []Item       { return g.state.items }
func (g *GameEngine) Triggers() []Trigger { return g.state.triggers }
func (g *GameEngine) ClearAuras()         { g.state.auras = nil }
func (g *GameEngine) ClearTriggers()      { g.state.triggers = nil }
func (g *GameEngine) ClearItems()         { g.state.items = nil }

// SetValue / SetActionPoints / SetCardsDrawn / SetArcaneIncomingDamage / SetBlockTotal /
// SetOpponentMarked / SetAuraCreated / SetCardBanished: per-permutation chain-locals
// reset paths.
func (g *GameEngine) SetValue(v int)                { g.state.value = v }
func (g *GameEngine) SetActionPoints(n int)         { g.state.actionPoints = n }
func (g *GameEngine) SetCardsDrawn(n int)           { g.state.cardsDrawn = n }
func (g *GameEngine) SetArcaneIncomingDamage(n int) { g.state.arcaneIncomingDamage = n }
func (g *GameEngine) SetBlockTotal(n int)           { g.state.blockTotal = n }
func (g *GameEngine) SetOpponentMarked(v bool)      { g.state.opponentMarked = v }
func (g *GameEngine) SetAuraCreated(v bool)         { g.state.auraCreated = v }
func (g *GameEngine) SetCardBanished(v bool)        { g.state.cardBanished = v }

// Copy returns a fresh *GameEngine wrapping a deep copy of g.state. The chain runner
// uses this to branch a per-permutation engine off the leaf master without mutating
// shared state.
func (g *GameEngine) Copy() *GameEngine {
	return &GameEngine{state: g.state.Copy()}
}

// BeginPermutation forwards to *GameState.BeginPermutation. Resets chain-locals on the
// underlying state in preparation for a fresh permutation's chain run.
func (g *GameEngine) BeginPermutation(hand []card.Card, incomingDamage int, logger *turnlogger.TurnLogger) {
	g.state.BeginPermutation(hand, incomingDamage, logger)
}

// DrawOne models a mid-turn draw: pop the top of the deck and append it to Hand. No-op
// on an empty deck. Bumps CardsDrawn so the partition tiebreaker can prefer chains with
// more draws. Inherits the IsCacheable flip via PopDeckTop.
func (g *GameEngine) DrawOne() {
	c, ok := g.PopDeckTop()
	if !ok {
		return
	}
	g.state.hand = append(g.state.hand, c)
	g.state.cardsDrawn++
}

// === Pure forwarders for cards-facing scalar accessors ===

func (g *GameEngine) ActionPoints() int                          { return g.state.actionPoints }
func (g *GameEngine) AddActionPoints(n int)                      { g.state.actionPoints += n }
func (g *GameEngine) ArcaneDamageDealt() bool                    { return g.state.arcaneDamageDealt }
func (g *GameEngine) SetArcaneDamageDealt(v bool)                { g.state.arcaneDamageDealt = v }
func (g *GameEngine) ArcaneIncomingDamage() int                  { return g.state.arcaneIncomingDamage }
func (g *GameEngine) AuraCreated() bool                          { return g.state.auraCreated }
func (g *GameEngine) BlockTotal() int                            { return g.state.blockTotal }
func (g *GameEngine) CardBanished() bool                         { return g.state.cardBanished }
func (g *GameEngine) CardsPlayed() []card.Card                   { return g.state.cardsPlayed }
func (g *GameEngine) SetCardsPlayed(cs []card.Card)              { g.state.cardsPlayed = cs }
func (g *GameEngine) CardsRemaining() []*card.CardState          { return g.state.cardsRemaining }
func (g *GameEngine) SetCardsRemaining(cs []*card.CardState)     { g.state.cardsRemaining = cs }
func (g *GameEngine) Defenders() []card.Card                     { return g.state.defenders }
func (g *GameEngine) SetDefenders(d []card.Card)                 { g.state.defenders = d }
func (g *GameEngine) IncomingDamage() int                        { return g.state.incomingDamage }
func (g *GameEngine) SetIncomingDamage(n int)                    { g.state.incomingDamage = n }
func (g *GameEngine) NonAttackActionPlayed() bool                { return g.state.nonAttackActionPlayed }
func (g *GameEngine) SetNonAttackActionPlayed(v bool)            { g.state.nonAttackActionPlayed = v }
func (g *GameEngine) OpponentMarked() bool                       { return g.state.opponentMarked }
func (g *GameEngine) MarkOpponent()                              { g.state.opponentMarked = true }
func (g *GameEngine) ClearOpponentMarked()                       { g.state.opponentMarked = false }
func (g *GameEngine) Pitched() []card.Card                       { return g.state.pitched }
func (g *GameEngine) TriggeringCard() card.Card                  { return g.state.triggeringCard }
func (g *GameEngine) SetTriggeringCard(c card.Card)              { g.state.triggeringCard = c }
func (g *GameEngine) AttackReactionTarget() *card.CardState      { return g.state.attackReactionTarget }
func (g *GameEngine) SetAttackReactionTarget(cs *card.CardState) { g.state.attackReactionTarget = cs }
func (g *GameEngine) Value() int                                 { return g.state.value }
func (g *GameEngine) AddValue(n int)                             { g.state.value += n }
func (g *GameEngine) CardsDrawn() int                            { return g.state.cardsDrawn }
func (g *GameEngine) AuraCount() int                             { return len(g.state.auras) }
func (g *GameEngine) CurrentStepRerouted() bool                  { return g.state.currentStepRerouted }
func (g *GameEngine) SetCurrentStepRerouted(v bool)              { g.state.currentStepRerouted = v }
func (g *GameEngine) AmendLastChainStepN(n int)                  { g.state.logger.AmendLastChainStepN(n) }
func (g *GameEngine) Logger() *turnlogger.TurnLogger             { return g.state.logger }
func (g *GameEngine) SetLogger(l *turnlogger.TurnLogger)         { g.state.logger = l }
func (g *GameEngine) LogEntries() []turnlogger.LogEntry          { return g.state.logger.Entries() }

// HasPlayedType reports whether any card played this turn has the given type. Universal
// cards' Types() folds the active hero's class through g so class-gated triggers see the
// right type-line.
func (g *GameEngine) HasPlayedType(t card.CardType) bool {
	for _, c := range g.state.cardsPlayed {
		if c.Types(g).Has(t) {
			return true
		}
	}
	return false
}

// === Rules-engine helpers cards reach through GameEngine ===

// LikelyToHit reports whether self's attack is likely to land past the opponent's blocks.
func (g *GameEngine) LikelyToHit(self *card.CardState) bool { return LikelyToHit(self) }

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit.
func (g *GameEngine) LikelyDamageHits(n int, dominate bool) bool {
	return LikelyDamageHits(n, dominate)
}

// OpponentDiscard credits n cards' worth of damage-equivalent value for forcing the
// opponent to discard. Returns the credited value for log attribution.
func (g *GameEngine) OpponentDiscard(n int) int {
	v := n * DiscardValue
	g.AddValue(v)
	return v
}

// Clash models a clash (rule 8.5.45): we and the opponent reveal the top card of our
// decks and the higher {p} wins. We model from our side only — our deck's top is read
// via PeekDeck; the opponent's top is approximated as 5-power. On a win (our top ≥ 6),
// win fires; on a loss (our top ≤ 4), lose fires; ties (top == 5) and empty deck fire
// neither. PeekDeck flips IsCacheable to false.
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
	g.state.cacheable = false
	if n <= 0 || g.state.deck.Size() == 0 {
		return
	}
	if n > g.state.deck.Size() {
		n = g.state.deck.Size()
	}
	drawn := g.state.deck.Draw(n)
	cards := make([]card.Card, len(drawn))
	for i, c := range drawn {
		cards[i] = c.(card.Card)
	}

	var top, bottom []card.Card
	if g.state.hero == nil {
		top = cards
	} else {
		top, bottom = g.state.hero.Opt(cards)
	}
	panicIfOptViolatesMultiset(cards, top, bottom)

	deckTop := make([]deck.Card, len(top))
	for i, c := range top {
		deckTop[i] = c
	}
	g.state.deck.PutTop(deckTop)
	deckBottom := make([]deck.Card, len(bottom))
	for i, c := range bottom {
		deckBottom[i] = c
	}
	g.state.deck.PutBottom(deckBottom)

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

// formatCardList renders cs as "[name1, name2, ...]" using DisplayName, or "[]" when
// empty.
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
// combined (top, bottom) output is exactly the input multiset.
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
