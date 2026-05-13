package gameengine

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// GameEngine is the rules-engine wrapper around a *GameState. Cards play against this
// type via the v2/card.GameEngine interface; the engine's method surface mixes
// (a) cacheable-aware accessors that flip cacheable as a side effect of touching hidden
// state and (b) rules-orchestration methods (Fire*, ResolveChainStep, Opt, Clash,
// DealArcaneDamage, token economy) that apply the FaB rule set on top of the raw state
// mutations the state itself exposes.
//
// *GameState is embedded so every pure state accessor (Hand, Deck, Auras, SetX, …) and
// the Copy / BeginPermutation utilities promote automatically. Methods listed below
// override the embedded ones to add cacheable-flipping or rules logic; everything else
// falls through to GameState verbatim.
//
// GameState owns the data and the cheap pure accessors; GameEngine owns the rules. The
// split lets internal machinery (TurnSummary.State, sim's per-permutation copy, the
// find-best winner) pass around a *GameState pointer without dragging the whole engine
// API along when all it needs to do is read or copy raw state.
type GameEngine struct {
	*GameState
}

// === Cards-facing zone accessors that flip cacheable. These shadow the same-name
//     methods promoted from *GameState; the embedded versions stay reachable as
//     g.GameState.X when the engine internals need the non-flipping variant.

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

// Deck returns the chain-runner deck for read-only inspection and flips IsCacheable to
// false. Card handlers should not mutate the returned *deck.Deck directly; route
// mutations through PopDeckTop / PrependToDeck / Opt / TutorFromDeck /
// RecycleToDeckBottom.
func (g *GameEngine) Deck() *deck.Deck {
	g.cacheable = false
	return g.deck
}

// PeekTopN returns the top n cards of the deck (top first) without removing them and
// flips IsCacheable to false. Returns fewer cards when the deck has < n.
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

// RecycleFromGraveyardToTop / RecycleFromGraveyardToBottom remove the first graveyard
// card matching pred and put it on the top / bottom of the deck. Flip IsCacheable.
func (g *GameEngine) RecycleFromGraveyardToTop(pred func(card.Card) bool) (card.Card, bool) {
	return g.recycleFromGraveyard(pred, true)
}
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
// IsCacheable to false. The promoted AppendGraveyard does the same append without the
// flip — framework-internal callers (the chain dispatcher's "non-persistent →
// graveyard" rule, Aura.OnDestroy) reach for that one instead.
func (g *GameEngine) AddToGraveyard(c card.Card) {
	g.cacheable = false
	g.graveyard = append(g.graveyard, c)
}

// DrawOne models a mid-turn draw: pop the top of the deck and append it to Hand. No-op
// on an empty deck. Bumps CardsDrawn so the partition tiebreaker can prefer chains with
// more draws. Inherits the IsCacheable flip via PopDeckTop.
func (g *GameEngine) DrawOne() {
	c, ok := g.PopDeckTop()
	if !ok {
		return
	}
	g.hand = append(g.hand, c)
	g.cardsDrawn++
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
