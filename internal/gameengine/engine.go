package gameengine

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// GameEngine is the rules-engine wrapper around a *GameState (which is embedded so pure
// accessors and Copy / Reset promote automatically). GameState owns the data; GameEngine
// owns the rules. The methods below override the embedded ones to add either
// cacheable-flipping (for accessors that touch hidden state) or rules orchestration
// (Fire*, ResolveAttackStep, Opt, Clash, DealArcaneDamage, token economy).
type GameEngine struct {
	*GameState
	// pitchBonus accumulates extra resources a triggertype.Pitch handler grants for the
	// current pitch. FirePitchTriggers zeroes it before the fire and reads it after.
	// Lives on the wrapper (not GameState) so it doesn't get copied per permutation.
	pitchBonus int
}

// === Cards-facing zone accessors that flip cacheable. Shadow the *GameState methods;
//     ge.GameState.X reaches the non-flipping variant when engine internals need it.

// Hand returns the cards in hand — INCLUDING entries tagged Pitch / Attack by the
// partition. The hand→pitch-zone transition is rules-modelled at pay time, so "a card in
// your hand" reads must see scheduled-but-uncommitted entries. Flips IsCacheable; prefer
// HandSize / HandHasMatching / HeldHandSize for non-mutating gates.
func (ge *GameEngine) Hand() []card.Card {
	ge.cacheable = false
	return ge.GameState.Hand()
}

// HeldHand returns the Held subset of the hand — Pitch / Attack entries excluded. Drawn
// entries are included so length agrees with HeldHandSize. Flips IsCacheable.
func (ge *GameEngine) HeldHand() []card.Card {
	ge.cacheable = false
	return heldHandSlice(ge.GameState.HandStates())
}

// HandHasMatching reports whether any non-drawn hand entry satisfies pred. Pred receives
// the engine handle plus the entry's *CardState so it can dispatch to the existing
// card.IsAttackAction / card.IsNonAttackAction / EffectiveCost helpers without a closure.
// FromDraw entries are skipped: their identity is unknown to in-attack-turn attribute reads.
// Doesn't flip IsCacheable — the starting-hand multiset is already part of the cache key.
func (ge *GameEngine) HandHasMatching(pred func(card.GameEngine, *card.CardState) bool) bool {
	states := ge.handStatesForMatching()
	for i := range states {
		if states[i].FromDraw {
			continue
		}
		if pred(ge, &states[i]) {
			return true
		}
	}
	return false
}

// heldHandSlice projects a role-tagged hand to its Held cards. Returns nil when empty.
func heldHandSlice(states []card.CardState) []card.Card {
	n := 0
	for i := range states {
		if states[i].Role == card.Held {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	out := make([]card.Card, 0, n)
	for i := range states {
		if states[i].Role == card.Held {
			out = append(out, states[i].Card)
		}
	}
	return out
}

// AppendHand inserts c into the hand at its Card.ID()-sorted position, flipping
// IsCacheable to false.
func (ge *GameEngine) AppendHand(c card.Card) {
	ge.insertHandSorted(c)
}

// insertHandSorted inserts c at the Card.ID()-sorted position, flipping IsCacheable. The
// canonical multiset ordering is required by the attack-turn runner and the eval-cache key.
func (ge *GameEngine) insertHandSorted(c card.Card) {
	ge.cacheable = false
	i := sort.Search(len(ge.hand), func(j int) bool { return ge.hand[j].Card.ID() > c.ID() })
	ge.hand = append(ge.hand, card.CardState{})
	copy(ge.hand[i+1:], ge.hand[i:])
	ge.hand[i] = card.CardState{Card: c, Role: card.Held}
}

// Graveyard returns the live graveyard slice and flips IsCacheable to false.
func (ge *GameEngine) Graveyard() []card.Card {
	ge.cacheable = false
	return ge.graveyard
}

// Deck returns the attack-turn runner deck for read-only inspection and flips IsCacheable.
// Card handlers must not mutate it directly — route through PopDeckTop / PrependToDeck /
// Opt / TutorFromDeck / RecycleToDeckBottom.
func (ge *GameEngine) Deck() *deck.Deck {
	ge.cacheable = false
	return ge.deck
}

// PeekTopN returns the top n cards of the deck (top first) without removing them and
// flips IsCacheable to false. Returns fewer cards when the deck has < n.
func (ge *GameEngine) PeekTopN(n int) []card.Card {
	ge.cacheable = false
	top := ge.deck.PeekTopN(n)
	if len(top) == 0 {
		return nil
	}
	out := make([]card.Card, len(top))
	for i, c := range top {
		out[i] = c.(card.Card)
	}
	return out
}

// PopDeckTop removes and returns the top card, (nil, false) when empty. Flips IsCacheable
// (caller observes identity) and notes the removal so the eval cache's depth check stays
// accurate even on uncacheable paths.
func (ge *GameEngine) PopDeckTop() (card.Card, bool) {
	ge.cacheable = false
	if ge.deck.Size() == 0 {
		return nil, false
	}
	ge.noteDeckRemoval(1)
	return ge.deck.Draw(1)[0].(card.Card), true
}

// MoveFromTopOfDeckToArsenal pops the top of the deck into the (currently-empty) arsenal
// slot without revealing it. Returns true when it ran (deck non-empty AND arsenal empty),
// false otherwise. Does NOT flip IsCacheable: the caller never reads the card's identity,
// so the operation is reproducible from the cache key (the popped card is determined by
// start-of-turn deck order, and the arsenal's identity at NEXT turn start is captured by
// that turn's own cache key). Notes the deck removal so the eval cache's depth check
// stays accurate.
func (ge *GameEngine) MoveFromTopOfDeckToArsenal() bool {
	if ge.arsenal != nil || ge.deck.Size() == 0 {
		return false
	}
	ge.noteDeckRemoval(1)
	ge.arsenal = ge.deck.Draw(1)[0].(card.Card)
	return true
}

// PeekDeck returns the top card of the deck without removing it. Returns (nil, false) on
// an empty deck. Flips IsCacheable to false.
func (ge *GameEngine) PeekDeck() (card.Card, bool) {
	ge.cacheable = false
	top := ge.deck.PeekTop()
	if top == nil {
		return nil, false
	}
	return top.(card.Card), true
}

// PrependToDeck inserts c at the top of the deck. Doesn't flip IsCacheable — caller
// supplied c, so the write is reproducible from cache key + attack-turn order.
func (ge *GameEngine) PrependToDeck(c card.Card) {
	ge.deck.PutTop([]deck.Card{c})
}

// AppendToDeck inserts c at the bottom of the deck. Cache-friendly like PrependToDeck.
func (ge *GameEngine) AppendToDeck(c card.Card) {
	ge.deck.PutBottom([]deck.Card{c})
}

// Discard pops the first Held-role hand card to the graveyard and logs under source. Cache-
// safe: discarded card's identity never escapes the engine, so cache can't diverge. After
// the move, fires the OnDiscardHook on the popped card when implemented (Fool's Gold:
// "when this is discarded, create a Gold token"). MoveFromHandToTopOfDeck /
// MoveFromHandToBottomOfDeck don't fire the hook — those model "put on top/bottom of deck"
// effects, not actual discards.
func (ge *GameEngine) Discard(source string) bool {
	c, ok := ge.popFirstHeldCard()
	if !ok {
		return false
	}
	ge.graveyard = append(ge.graveyard, c)
	ge.logger.AppendPostTriggerf(source, 0, "Discarded a card")
	if hook, ok := c.(card.OnDiscardHook); ok {
		hook.OnDiscard(ge, ge.logger)
	}
	return true
}

// MoveFromHandToTopOfDeck pops the first Held-role hand card to the top of the deck and
// logs under source. Models FaB's "put a card from your hand on top of your deck" effects
// (Moon Wish / Rise Above alt cost, Seek Horizon). Not a discard event — the OnDiscardHook
// (Fool's Gold) doesn't fire. Cache-safe.
func (ge *GameEngine) MoveFromHandToTopOfDeck(source string) bool {
	c, ok := ge.popFirstHeldCard()
	if !ok {
		return false
	}
	ge.deck.PutTop([]deck.Card{c})
	ge.logger.AppendPostTriggerf(source, 0, "Cycled a card to top of deck")
	return true
}

// MoveFromHandToBottomOfDeck pops the first Held-role hand card to the bottom of the deck
// and logs under source. Models FaB's "put a card from your hand on the bottom of your
// deck" effects (Emissary cycle, Scour the Battlescape, Sift). Not a discard event — the
// OnDiscardHook (Fool's Gold) doesn't fire. Cache-safe.
func (ge *GameEngine) MoveFromHandToBottomOfDeck(source string) bool {
	c, ok := ge.popFirstHeldCard()
	if !ok {
		return false
	}
	ge.deck.PutBottom([]deck.Card{c})
	ge.logger.AppendPostTriggerf(source, 0, "Cycled a card to bottom of deck")
	return true
}

// popFirstHeldCard removes the first Held-role hand entry and returns its Card. Drawn
// entries are eligible — indistinguishable to no-return-card callers.
func (ge *GameEngine) popFirstHeldCard() (card.Card, bool) {
	for rawIdx := range ge.hand {
		if ge.hand[rawIdx].Role != card.Held {
			continue
		}
		c := ge.hand[rawIdx].Card
		ge.hand = append(ge.hand[:rawIdx], ge.hand[rawIdx+1:]...)
		return c, true
	}
	return nil, false
}

// RecycleToDeckBottom appends pc.Card to the bottom of the deck and flags the attack turn
// dispatcher to skip the usual non-persistent → graveyard append. Models the FaB clause
// "put this on the bottom of its owner's deck". Doesn't flip IsCacheable — the caller
// supplies pc, so the write is reproducible from the cache key + attack-turn order.
func (ge *GameEngine) RecycleToDeckBottom(pc *card.CardState) {
	ge.deck.PutBottom([]deck.Card{pc.Card})
	ge.currentStepRerouted = true
}

// TutorFromDeck removes and returns the highest-scoring card per score. Returns (nil,
// false) when no card scores > 0 (or the deck is empty). Flips IsCacheable to false.
func (ge *GameEngine) TutorFromDeck(score func(card.Card) int) (card.Card, bool) {
	ge.cacheable = false
	got, ok := ge.deck.Tutor(func(c deck.Card) int { return score(c.(card.Card)) })
	if !ok {
		return nil, false
	}
	ge.noteDeckRemoval(1)
	return got.(card.Card), true
}

// BanishFromGraveyard removes the first graveyard card matching pred, appends it to the
// banished zone, and returns it. Returns (nil, false) when no card matches. Flips
// IsCacheable to false. Sets CardBanished so this-turn-banish riders fire correctly.
func (ge *GameEngine) BanishFromGraveyard(pred func(card.GameEngine, *card.CardState) bool) (card.Card, bool) {
	ge.cacheable = false
	var cs card.CardState
	for i, c := range ge.graveyard {
		cs = card.CardState{Card: c}
		if !pred(ge, &cs) {
			continue
		}
		ge.banished = append(ge.banished, c)
		ge.cardBanished = true
		ge.graveyard = append(ge.graveyard[:i], ge.graveyard[i+1:]...)
		return c, true
	}
	return nil, false
}

// RecycleFromGraveyardToTop / RecycleFromGraveyardToBottom remove the first graveyard
// card matching pred and put it on the top / bottom of the deck. Flip IsCacheable.
// pred receives the engine handle plus a transient *card.CardState wrapping each
// graveyard entry so it can read EffectiveCost / EffectiveTypes uniformly with the hand
// and in-play sites, or dispatch directly to card.IsAttackAction / card.IsNonAttackAction.
func (ge *GameEngine) RecycleFromGraveyardToTop(pred func(card.GameEngine, *card.CardState) bool) (card.Card, bool) {
	return ge.recycleFromGraveyard(pred, true)
}
func (ge *GameEngine) RecycleFromGraveyardToBottom(pred func(card.GameEngine, *card.CardState) bool) (card.Card, bool) {
	return ge.recycleFromGraveyard(pred, false)
}

func (ge *GameEngine) recycleFromGraveyard(pred func(card.GameEngine, *card.CardState) bool, toTop bool) (card.Card, bool) {
	ge.cacheable = false
	var cs card.CardState
	for i, c := range ge.graveyard {
		cs = card.CardState{Card: c}
		if !pred(ge, &cs) {
			continue
		}
		ge.graveyard = append(ge.graveyard[:i], ge.graveyard[i+1:]...)
		if toTop {
			ge.deck.PutTop([]deck.Card{c})
		} else {
			ge.deck.PutBottom([]deck.Card{c})
		}
		return c, true
	}
	return nil, false
}

// AddToGraveyard appends c to graveyard so later-resolving cards see it. Doesn't flip
// IsCacheable — the caller supplies c (a card whose identity is already known to the
// caller), so the write is reproducible from the cache key + attack-turn order.
func (ge *GameEngine) AddToGraveyard(c card.Card) {
	ge.graveyard = append(ge.graveyard, c)
}

// DrawOne models a mid-turn draw: pop the top of the deck into the hand at its sorted
// position. Reports whether a card was drawn; false on an empty deck. A successful
// draw doesn't flip IsCacheable — the cached attack turn replays DrawOne against the
// caller's current deck, and any downstream read of the drawn card's attributes
// (Hand / HeldHand / PeekTopN) flips IsCacheable through its own accessor. An
// empty-deck failure DOES flip IsCacheable: the card's no-draw branch took here was
// forced by an exhausted deck, and a future Best call with the same cache key but a
// non-empty deck would explore the successful-draw branch and produce a different
// optimum, so the entry must not be stored.
func (ge *GameEngine) DrawOne() bool {
	if ge.deck == nil || ge.deck.Size() == 0 {
		ge.cacheable = false
		return false
	}
	c := ge.deck.Draw(1)[0].(card.Card)
	ge.noteDeckRemoval(1)
	i := sort.Search(len(ge.hand), func(j int) bool { return ge.hand[j].Card.ID() > c.ID() })
	ge.hand = append(ge.hand, card.CardState{})
	copy(ge.hand[i+1:], ge.hand[i:])
	ge.hand[i] = card.CardState{Card: c, Role: card.Held, FromDraw: true}
	return true
}

// === Rules-engine helpers cards reach through GameEngine ===

// LikelyToHit reports whether pc's attack is likely to land past the opponent's blocks.
func (ge *GameEngine) LikelyToHit(pc *card.CardState) bool { return LikelyToHit(pc) }

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit.
func (ge *GameEngine) LikelyDamageHits(n int, dominate bool) bool {
	return LikelyDamageHits(n, dominate)
}

// LikelyDamageDealt is the "how much" sibling — see the package-level function.
func (ge *GameEngine) LikelyDamageDealt(n int, dominate bool) int {
	return LikelyDamageDealt(n, dominate)
}

// OpponentDiscard credits n cards' worth of damage-equivalent value for forcing the
// opponent to discard. Returns the credited value for log attribution.
func (ge *GameEngine) OpponentDiscard(n int) int {
	v := n * DiscardValue
	ge.AddValue(v)
	return v
}

// DestroyOpponentArsenal destroys the opposing arsenal, crediting one card's worth of the
// OpponentDiscard heuristic (DiscardValue). Idempotent per turn: a second call returns 0
// and credits nothing. Returns the value credited for log attribution.
func (ge *GameEngine) DestroyOpponentArsenal() int {
	if ge.destroyedOpponentArsenal {
		return 0
	}
	ge.destroyedOpponentArsenal = true
	return ge.OpponentDiscard(1)
}

// PreventArcaneDamage caps remaining arcane damage by up to n, crediting the prevented
// amount to Value. Returns the amount actually prevented — the lesser of n and
// RemainingArcaneDamage, clamped at 0 — for log attribution. Banks the prevention into
// arcaneDamageBlocked (which resets per turn) rather than mutating the constant matchup
// figure, so prevention doesn't leak into later turns.
func (ge *GameEngine) PreventArcaneDamage(n int) int {
	if n <= 0 {
		return 0
	}
	rem := ge.incomingArcaneDamage - ge.arcaneDamageBlocked
	if rem <= 0 {
		return 0
	}
	if n > rem {
		n = rem
	}
	ge.arcaneDamageBlocked += n
	ge.AddValue(n)
	return n
}

// PreventPhysicalDamage caps remaining physical damage by up to n, crediting the prevented
// amount to Value. Returns the amount actually prevented (lesser of n and the current
// RemainingPhysicalDamage, clamped at 0) for log attribution. Banks the prevention into
// physicalDamageBlocked so RemainingPhysicalDamage reads the reduced figure for downstream
// triggers (a DamageAboutToBeTaken handler that absorbs the swing chains into the
// DamageTaken gate this way).
func (ge *GameEngine) PreventPhysicalDamage(n int) int {
	if n <= 0 {
		return 0
	}
	rem := ge.incomingPhysicalDamage - ge.physicalDamageBlocked
	if rem <= 0 {
		return 0
	}
	if n > rem {
		n = rem
	}
	ge.physicalDamageBlocked += n
	ge.AddValue(n)
	return n
}

// PreventGenericDamage absorbs up to n damage of whichever type is currently active —
// physical when there's remaining unblocked physical, otherwise arcane. Mirrors the
// official-rules wording "prevent N damage": the card asks for N prevented without naming
// a type, the engine picks. The chosen typed method credits Value; this returns the amount
// actually prevented (0 when both types are empty) for log attribution.
func (ge *GameEngine) PreventGenericDamage(n int) int {
	if prevented := ge.PreventPhysicalDamage(n); prevented > 0 {
		return prevented
	}
	return ge.PreventArcaneDamage(n)
}

// TurnFaceUp flips pc.FaceUp = true on the specific CardState the caller passes — found
// by scanning CardsRemaining for an in-attack-turn target, or held directly when the target is
// the arsenal-in or another known CardState pointer — then fires pc.Card.OnFaceUp if
// pc.Card implements card.FaceUpHook. Touches a single instance rather than every
// scheduled copy with the same identity; the caller picks the target.
func (ge *GameEngine) TurnFaceUp(pc *card.CardState) {
	pc.FaceUp = true
	if hook, ok := pc.Card.(card.FaceUpHook); ok {
		hook.OnFaceUp(ge, ge.logger)
	}
}

// Clash models a clash (rule 8.5.45): we and the opponent reveal the top card of our
// decks and the higher {p} wins. We model from our side only — our deck's top is read
// via PeekDeck; the opponent's top is approximated as 5-power. On a win (our top ≥ 6),
// win fires; on a loss (our top ≤ 4), lose fires; ties (top == 5) and empty deck fire
// neither. PeekDeck flips IsCacheable to false.
func (ge *GameEngine) Clash(win, lose func()) {
	top, ok := ge.PeekDeck()
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

// Opt resolves the FaB "Opt N" keyword via the current hero's splitter. Wrapper around
// OptWith; see that method for the contract.
func (ge *GameEngine) Opt(l card.Logger, n int) {
	split := func(cs []card.Card) (top, bottom []card.Card) { return cs, nil }
	if ge.hero != nil {
		split = ge.hero.Opt
	}
	ge.OptWith(l, n, split)
}

// OptWith is Opt with a caller-supplied splitter — for cards whose printed text overrides
// the hero's Opt heuristic with a card-local rule. Pops up to n cards from the top of the
// deck and hands them to split; the returned top list goes back on top of the deck (in
// returned order) and the bottom list appends to the bottom (in returned order). n is
// clamped to the current deck length. Always flips IsCacheable to false.
//
// Emits a log entry naming the revealed cards and the split when the handler ran.
//
// Panics if the handler's combined output isn't exactly the input multiset.
func (ge *GameEngine) OptWith(l card.Logger, n int, split func([]card.Card) (top, bottom []card.Card)) {
	ge.cacheable = false
	if n <= 0 || ge.deck.Size() == 0 {
		return
	}
	if n > ge.deck.Size() {
		n = ge.deck.Size()
	}
	drawn := ge.deck.Draw(n)
	ge.noteDeckRemoval(n)
	cards := make([]card.Card, len(drawn))
	for i, c := range drawn {
		cards[i] = c.(card.Card)
	}

	top, bottom := split(cards)
	panicIfOptViolatesMultiset(cards, top, bottom)

	deckTop := make([]deck.Card, len(top))
	for i, c := range top {
		deckTop[i] = c
	}
	ge.deck.PutTop(deckTop)
	deckBottom := make([]deck.Card, len(bottom))
	for i, c := range bottom {
		deckBottom[i] = c
	}
	ge.deck.PutBottom(deckBottom)

	if OptDebug {
		fmt.Printf("Opt(%d): cards=%s -> top=%s bottom=%s\n",
			n, formatCardList(cards), formatCardList(top), formatCardList(bottom))
	}
	if l == nil {
		return
	}
	l.AppendAttackStepf(0, "Opted %s, put %s on top, put %s on bottom",
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

// === Trigger and aura dispatch ===

// FireTriggers fires the hero plus every Aura, EphemeralTrigger, and Item registered for
// trigger type t. It is the single dispatch point for every triggertype.Type lifecycle
// event.
//
// triggeringCard is the physical *CardState whose resolution raised the event, or nil for
// turn-boundary events. It is published on ge.triggeringCard so handlers can read its
// in-attack-turn ephemeral state and write back to it. The hero fires first, then auras,
// ephemeral triggers, and items.
func (ge *GameEngine) FireTriggers(t triggertype.Type, triggeringCard *card.CardState) {
	heroFires := ge.heroTriggerType&t != 0
	if !heroFires && !ge.AnyAurasInPlay() && len(ge.triggers) == 0 && !ge.AnyItemsInPlay() {
		return
	}
	ge.triggeringCard = triggeringCard

	var triggeringTypes card.TypeSet
	if triggeringCard != nil {
		triggeringTypes = triggeringCard.Card.Types(ge)
	}

	// Snapshot which token slots are live at fire-time. fireHooks uses
	// `n := len(*hooks)` for the same reason: a token created by an aura/trigger fired
	// during this FireTriggers call must not fire on the same event (FaB rule: triggers
	// only see events that happened *before* the trigger source entered the game). The
	// liveBits caches mean the snapshot is a single byte load instead of an N-slot
	// interface walk.
	liveAuraBits := ge.tokenAurasLiveBits
	liveItemBits := ge.tokenItemsLiveBits

	if heroFires {
		fireHero(ge, t, triggeringCard, triggeringTypes)
	}
	fireHooks(ge, &ge.auras, t, triggeringCard, triggeringTypes, false)
	if liveAuraBits != 0 {
		fireTokenAuras(ge, t, triggeringCard, triggeringTypes, liveAuraBits)
	}
	fireHooks(ge, &ge.triggers, t, triggeringCard, triggeringTypes, true)
	fireHooks(ge, &ge.items, t, triggeringCard, triggeringTypes, false)
	if liveItemBits != 0 {
		fireTokenItems(ge, t, triggeringCard, triggeringTypes, liveItemBits)
	}

	ge.triggeringCard = nil
}

// fireHero applies the OncePerTurn / Matches gates and invokes the hero handler. The
// hero is singular (no slice splicing, no removeAfterFire) so it bypasses fireHooks's
// cursor walk. The TriggerType bit-and check is done at the caller; t is the firing
// event so a multi-subscription hero can dispatch.
func fireHero(ge *GameEngine, t triggertype.Type, triggeringCard *card.CardState, triggeringTypes card.TypeSet) {
	h := ge.hero
	if h.OncePerTurn() && h.FiredThisTurn() {
		return
	}
	if triggeringCard != nil && !h.Matches(triggeringTypes) {
		return
	}
	h.Fire(ge, ge.logger, t)
	if h.OncePerTurn() {
		h.SetFiredThisTurn(true)
	}
}

// fireHooks fires every entry of *hooks subscribed to event t: an open once-per-turn gate
// is required, and card-raised events additionally need the entry's type filter to accept
// triggeringCard. The snapshot length is taken up front so an entry a handler creates
// lands past it and isn't fired this pass. A cursor walk keeps a handler-side destroy from
// skipping the next entry — currentHookIdx is published before each Fire so the entry's
// Destroy splices the right slot and sets currentHookDestroyed, which shortens the walk.
// removeAfterFire splices every fired entry unconditionally: the one-shot semantics of
// ephemeral triggers, which the engine drops once they fire.
func fireHooks[H trigger.Hook](ge *GameEngine, hooks *[]H, t triggertype.Type, triggeringCard *card.CardState, triggeringTypes card.TypeSet, removeAfterFire bool) {
	n := len(*hooks)
	for i := 0; i < n; {
		h := (*hooks)[i]
		if h.TriggerType()&t == 0 || (h.OncePerTurn() && h.FiredThisTurn()) ||
			(triggeringCard != nil && !h.Matches(triggeringTypes)) {
			i++
			continue
		}
		ge.currentHookIdx = i
		ge.currentHookDestroyed = false
		h.Fire(ge, ge.logger, t)
		switch {
		case ge.currentHookDestroyed:
			n--
		case removeAfterFire:
			*hooks = append((*hooks)[:i], (*hooks)[i+1:]...)
			n--
		default:
			(*hooks)[i].SetFiredThisTurn(true)
			i++
		}
	}
	ge.currentHookIdx = -1
}

// fireTokenAuras walks the live token-aura slots indicated by liveAtStart (a snapshot
// of tokenAurasLiveBits taken at FireTriggers entry), firing each whose hook gate /
// type-filter accepts the event. The liveAtStart guard mirrors fireHooks's snapshot
// length trick: a token created mid-FireTriggers won't have its bit set here, so it
// can't fire on the same event. Publishes currentFiringTokenAura before each Fire so
// DestroyAura routes to the slot's SetCount(0) path.
func fireTokenAuras(ge *GameEngine, t triggertype.Type, triggeringCard *card.CardState, triggeringTypes card.TypeSet, liveAtStart uint8) {
	for mask := liveAtStart; mask != 0; {
		low := mask & -mask
		i := bits.TrailingZeros8(low)
		mask ^= low
		a := ge.tokenAuras[i]
		// Re-check count: a prior token's Fire in this same pass may have routed
		// through DestroyAura and zeroed this slot via the currentFiringTokenAura
		// path. liveAtStart only guarantees the slot was live at entry.
		if a.Count() == 0 {
			continue
		}
		if a.TriggerType()&t == 0 || (a.OncePerTurn() && a.FiredThisTurn()) ||
			(triggeringCard != nil && !a.Matches(triggeringTypes)) {
			continue
		}
		ge.currentFiringTokenAura = i
		ge.currentHookDestroyed = false
		a.Fire(ge, ge.logger, t)
		if !ge.currentHookDestroyed {
			a.SetFiredThisTurn(true)
		}
	}
	ge.currentFiringTokenAura = -1
}

// fireTokenItems is the item-side counterpart of fireTokenAuras.
func fireTokenItems(ge *GameEngine, t triggertype.Type, triggeringCard *card.CardState, triggeringTypes card.TypeSet, liveAtStart uint8) {
	for mask := liveAtStart; mask != 0; {
		low := mask & -mask
		i := bits.TrailingZeros8(low)
		mask ^= low
		it := ge.tokenItems[i]
		if it.Count() == 0 {
			continue
		}
		if it.TriggerType()&t == 0 || (it.OncePerTurn() && it.FiredThisTurn()) ||
			(triggeringCard != nil && !it.Matches(triggeringTypes)) {
			continue
		}
		ge.currentFiringTokenItem = i
		ge.currentHookDestroyed = false
		it.Fire(ge, ge.logger, t)
		if !ge.currentHookDestroyed {
			it.SetFiredThisTurn(true)
		}
	}
	ge.currentFiringTokenItem = -1
}

// DestroyAura removes the aura currently being fired from the arena. It then fires the
// source card's OnLeavesArena hook — the printed "when this leaves the arena" clause, when
// the card implements card.LeavesArenaAura — and, when addToGraveyard==true, pushes the
// source card into the graveyard (no-op for token auras with no source). OnLeavesArena runs
// before the graveyard append so a "banish another aura from your graveyard" leave clause
// doesn't see the just-left card. Direct splice with no cacheable flip — destruction is
// deterministic from the triggering event.
func (ge *GameEngine) DestroyAura(addToGraveyard bool) {
	// Token-aura path: when a token slot is mid-Fire, zero its count instead of
	// splicing ge.auras. Tokens have no source card, so the addToGraveyard /
	// OnLeavesArena dance doesn't apply to them.
	if ge.currentFiringTokenAura >= 0 {
		ge.tokenAuras[ge.currentFiringTokenAura].SetCount(0)
		ge.tokenAurasLiveBits &^= 1 << ge.currentFiringTokenAura
		ge.currentHookDestroyed = true
		return
	}
	i := ge.currentHookIdx
	if i < 0 || i >= len(ge.auras) {
		return
	}
	dying := ge.auras[i]
	src := dying.SourceCard()
	// Clear the underlying pool slot (SourceCard()==nil, Count()==0) so the next
	// CreateAura's free-slot scan can reuse it. Splicing ge.auras alone leaves the
	// pool slot occupied with stale fields.
	dying.Clear()
	ge.auras = append(ge.auras[:i], ge.auras[i+1:]...)
	ge.currentHookDestroyed = true
	if la, ok := src.(card.LeavesArenaAura); ok {
		la.OnLeavesArena(ge, ge.logger)
	}
	if src != nil && addToGraveyard {
		ge.AppendGraveyard(src.(card.Card))
	}
}

// DestroyItem removes the item currently being fired from the arena and, when
// addToGraveyard is true, pushes its source card into the graveyard (no-op for token
// items with no source). The item counterpart of DestroyAura: direct splice with no
// cacheable flip — destruction is deterministic from the triggering event.
func (ge *GameEngine) DestroyItem(addToGraveyard bool) {
	// Token-item path: see DestroyAura.
	if ge.currentFiringTokenItem >= 0 {
		ge.tokenItems[ge.currentFiringTokenItem].SetCount(0)
		ge.tokenItemsLiveBits &^= 1 << ge.currentFiringTokenItem
		ge.currentHookDestroyed = true
		return
	}
	i := ge.currentHookIdx
	if i < 0 || i >= len(ge.items) {
		return
	}
	src := ge.items[i].SourceCard()
	ge.items = append(ge.items[:i], ge.items[i+1:]...)
	ge.currentHookDestroyed = true
	if src != nil && addToGraveyard {
		ge.AppendGraveyard(src.(card.Card))
	}
}

// AddResourcePoints adds n resources to the card currently being pitched. A triggertype.Pitch
// handler calls it to boost what the pitched card yields beyond its printed Pitch value;
// pay folds the total into the pitch pool. No effect outside a pitch fire.
func (ge *GameEngine) AddResourcePoints(n int) { ge.pitchBonus += n }

// FirePitchTriggers fires the triggertype.Pitch event for a just-pitched card and returns
// the resource bonus its handlers granted via AddResourcePoints — the amount pay adds to the
// pitched card's contribution on top of its printed Pitch value. pitched is the physical
// pitched-zone copy.
func (ge *GameEngine) FirePitchTriggers(pitched *card.CardState) int {
	ge.pitchBonus = 0
	ge.FireTriggers(triggertype.Pitch, pitched)
	return ge.pitchBonus
}

// SacrificePayoffAura destroys one aura the player controls and reports whether it
// destroyed one. It targets the first aura whose source card carries a leave-the-arena
// payoff (card.LeavesArenaAura), fires that OnLeavesArena clause, and graveyards the card.
// Auras with no leave payoff are skipped — the method exists to cash a payoff on demand,
// so there is nothing to gain from destroying one without it.
func (ge *GameEngine) SacrificePayoffAura() bool {
	for i, a := range ge.auras {
		la, ok := a.SourceCard().(card.LeavesArenaAura)
		if !ok {
			continue
		}
		src := a.SourceCard().(card.Card)
		a.Clear()
		ge.auras = append(ge.auras[:i], ge.auras[i+1:]...)
		la.OnLeavesArena(ge, ge.logger)
		ge.AppendGraveyard(src)
		return true
	}
	return false
}

// === Attack-step resolution ===

// ResolveAttackStep runs card.Play on pc and then applies the standard attack-step
// resolution: attack-action / weapon-attack credit pc.EffectiveAttack() to ge.value;
// defense-reaction (or DefensiveInstant) credits EffectiveDefense capped at the remaining
// unblocked damage; everything else logs (+0). The "<DisplayName>: <VERB> (+N)" attack-step
// entry is appended after Play returns so self-buffs Play applied are reflected in the
// displayed delta.
func (ge *GameEngine) ResolveAttackStep(l card.Logger, pc *card.CardState) {
	pc.Card.Play(ge, l, pc)
	// EffectiveTypes routes through ModalTypes for cards whose type-line shifts per mode
	// (Tip-Off mode 1 reads as Generic Instant, not Generic Action - Attack), so the
	// downstream aura-create flip and attack-step delta land on the mode-correct TypeSet.
	types := pc.EffectiveTypes(ge)
	if types.Has(card.TypeAura) {
		ge.auraCreated = true
	}
	n := ge.attackStepDelta(pc, types)
	// NoopLogger discards the text, so skip the cached-string lookup and the AppendAttackStep
	// dispatch entirely on the eval hot path — every attack-step resolution hits this.
	if _, isNoop := l.(NoopLogger); isNoop {
		return
	}
	l.AppendAttackStep(AttackStepText(pc), n)
}

// PlayCard implements card.GameEngine.PlayCard — resolves another card mid-handler.
func (ge *GameEngine) PlayCard(l card.Logger, pc *card.CardState) {
	ge.ResolveAttackStep(l, pc)
}

// attackStepDelta computes the attack step's display delta and applies the standard damage /
// block side effects. Returns the (+N) value for the log line. types is the caller's
// already-resolved Types(nil) so we skip a second interface dispatch.
func (ge *GameEngine) attackStepDelta(pc *card.CardState, types card.TypeSet) int {
	switch {
	case types.IsAttackAction() || types.IsWeaponAttack():
		n := pc.EffectiveAttack()
		ge.value += n
		return n
	case types.IsDefenseReaction() || isDefensiveInstant(pc.Card):
		n := pc.EffectiveDefense()
		if rem := ge.incomingPhysicalDamage - ge.physicalDamageBlocked; n > rem {
			n = rem
		}
		if n < 0 {
			n = 0
		}
		ge.physicalDamageBlocked += n
		ge.value += n
		return n
	}
	return 0
}

// isDefensiveInstant reports whether c opts into the DR resolution path via the
// DefensiveInstant marker.
func isDefensiveInstant(c card.Card) bool {
	_, ok := c.(card.DefensiveInstant)
	return ok
}

// AttackStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix for the attack turn-
// step log line. VERB picks WEAPON ATTACK for Weapon+Attack, ATTACK for attack-actions,
// DEFENSE REACTION for DRs, and PLAY otherwise. EffectiveTypes dispatches on mode so
// ModalTypes cards (Tip-Off mode 1) log under the mode-correct verb.
//
// Output is memoised on (Card.ID, FromArsenal) — DisplayName / type-line / verb are all
// static per ID — so the per-Play DisplayName concat happens at most twice per card kind
// for the whole process. ModalTypes cards (their type-line shifts with self.Mode) bypass
// the cache and rebuild every call. Tests using a Card with ids.InvalidCard hit the same
// no-cache path.
func AttackStepText(pc *card.CardState) string {
	if _, ok := pc.Card.(card.ModalTypes); ok {
		return buildAttackStepText(pc)
	}
	id := pc.Card.ID()
	if id == ids.InvalidCard {
		return buildAttackStepText(pc)
	}
	idx := attackStepCacheIndex(id, pc.FromArsenal)
	if s := attackStepCache[idx].Load(); s != nil {
		return *s
	}
	out := buildAttackStepText(pc)
	attackStepCache[idx].Store(&out)
	return out
}

// buildAttackStepText is the uncached renderer; the cached path falls through to it on
// miss and ModalTypes / InvalidCard inputs route here every call.
func buildAttackStepText(pc *card.CardState) string {
	types := pc.EffectiveTypes(nil)
	var verb string
	switch {
	case types.IsWeaponAttack():
		verb = "WEAPON ATTACK"
	case types.IsAttackAction():
		verb = "ATTACK"
	case types.IsDefenseReaction():
		verb = "DEFENSE REACTION"
	default:
		verb = "PLAY"
	}
	if pc.FromArsenal {
		verb += " from arsenal"
	}
	return pc.Card.DisplayName() + ": " + verb
}

// attackStepCache memoises AttackStepText results keyed by (Card.ID, FromArsenal). Two rows
// per card cover the in-hand and from-arsenal verb suffixes. Sized for the full uint16 ID
// space (~1 MB) so lookups are direct bounds-checked array reads. Multiple goroutines
// computing the same entry produce the same string, so racing writers converge.
const attackStepCacheSize = 1 << 17 // 2 entries per ID × 65536 IDs

var attackStepCache [attackStepCacheSize]atomic.Pointer[string]

// attackStepCacheIndex packs (id, fromArsenal) into a single uint32 cache index. Bit 16 is
// the FromArsenal flag, bits 0-15 are the card ID — keeps the in-hand and from-arsenal
// variants in adjacent halves so the hot path is a plain array read.
func attackStepCacheIndex(id ids.CardID, fromArsenal bool) uint32 {
	idx := uint32(id)
	if fromArsenal {
		idx |= 1 << 16
	}
	return idx
}

// === Arcane damage ===

// DealArcaneDamage credits n arcane damage to Value, runs the arcane side effects
// (RegisterArcaneDamage), and writes the "Dealt n arcane damage" rider line under source.
// Routes through dealtArcaneText to avoid per-call fmt.Sprintf and variadic-int boxing.
func (ge *GameEngine) DealArcaneDamage(l card.Logger, source string, n int) {
	ge.AddValue(n)
	ge.RegisterArcaneDamage(n)
	if n >= 0 && n < len(dealtArcaneText) {
		l.AppendPostTrigger(source, dealtArcaneText[n], n)
		return
	}
	l.AppendPostTriggerf(source, n, "Dealt %d arcane damage", n)
}

// === Crowd reactions ===

// CrowdCheer flips hasCrowdCheered and fires the CrowdCheer trigger so "if you've been
// cheered this turn" gates and "whenever the crowd cheers you" handlers both see this
// turn's cheer. Source-side gating (which heroes are Revered / Reviled) belongs to the
// caller; this method only records the cheer landing on your hero.
func (ge *GameEngine) CrowdCheer() {
	ge.hasCrowdCheered = true
	ge.FireTriggers(triggertype.CrowdCheer, nil)
}

// CrowdBoo is the boo-side counterpart to CrowdCheer.
func (ge *GameEngine) CrowdBoo() {
	ge.hasCrowdBooed = true
	ge.FireTriggers(triggertype.CrowdBoo, nil)
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

// === Tokens ===

// Card-facing token creation / count methods on *GameEngine. Each method bumps the
// pre-allocated slot for its kind via bumpTokenAura / bumpTokenItem.

// CreateRunechants creates n Runechant tokens and credits +n damage at creation. Bumps
// the count on the pre-allocated Runechant slot; the slot's *Aura is reused across
// creations / destructions for the GameState's lifetime so there's no allocation here.
func (ge *GameEngine) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	ge.AddValue(n)
	ge.GameState.bumpTokenAura(tokenAuraRunechant, n)
}

// CreatePonders creates n Ponder tokens. No Value credit — Ponder pays out at end of turn.
func (ge *GameEngine) CreatePonders(n int) {
	if n <= 0 {
		return
	}
	ge.GameState.bumpTokenAura(tokenAuraPonder, n)
}

// CreateQuicken creates n Quicken tokens. Each charge grants the next attack-action card
// played Go again, then consumes. No Value credit at creation — the value lands when the
// chain extends from the granted go-again.
func (ge *GameEngine) CreateQuicken(n int) {
	if n <= 0 {
		return
	}
	ge.GameState.bumpTokenAura(tokenAuraQuicken, n)
}

// CreateGold / CreateSilver / CreateCopper create the matching token items. No Value
// credit — items only pay out when the activated ability is spent.
func (ge *GameEngine) CreateGold(n int) {
	if n <= 0 {
		return
	}
	ge.GameState.bumpTokenItem(tokenItemGold, n)
}
func (ge *GameEngine) CreateSilver(n int) {
	if n <= 0 {
		return
	}
	ge.GameState.bumpTokenItem(tokenItemSilver, n)
}
func (ge *GameEngine) CreateCopper(n int) {
	if n <= 0 {
		return
	}
	ge.GameState.bumpTokenItem(tokenItemCopper, n)
}

// CreateFrailtyForOpponent / CreateInertiaForOpponent / CreateBloodrotPoxForOpponent credit
// the matching damage-equivalent heuristic when a status token is created under the
// opponent's control. We don't track opposing status-token state, so these are flat value
// credits — see heuristics.go for the per-token rationale.
func (ge *GameEngine) CreateFrailtyForOpponent()     { ge.AddValue(FrailtyValue) }
func (ge *GameEngine) CreateInertiaForOpponent()     { ge.AddValue(InertiaValue) }
func (ge *GameEngine) CreateBloodrotPoxForOpponent() { ge.AddValue(BloodrotPoxValue) }

// SpendGold / SpendSilver / SpendCopper decrement the matching token slot by n, floored
// at zero. The token cards' Play handlers call these to consume the activated ability's
// payment.
func (ge *GameEngine) SpendGold(n int)   { ge.consumeTokenItem(tokenItemGold, n) }
func (ge *GameEngine) SpendSilver(n int) { ge.consumeTokenItem(tokenItemSilver, n) }
func (ge *GameEngine) SpendCopper(n int) { ge.consumeTokenItem(tokenItemCopper, n) }

func (ge *GameEngine) consumeTokenItem(kind tokenItemKind, n int) {
	it := ge.tokenItems[kind]
	c := it.Count() - n
	if c < 0 {
		c = 0
	}
	it.SetCount(c)
	if c == 0 {
		ge.tokenItemsLiveBits &^= 1 << kind
	}
}
