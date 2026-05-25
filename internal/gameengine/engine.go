package gameengine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// GameEngine is the rules-engine wrapper around a *GameState. Cards play against this type
// via internal/card.GameEngine. The method surface mixes (a) cacheable-aware accessors that
// flip cacheable as a side effect of touching hidden state and (b) rules-orchestration
// methods (Fire*, ResolveChainStep, Opt, Clash, DealArcaneDamage, token economy).
//
// *GameState is embedded so every pure accessor and the Copy / Reset utilities promote
// automatically. Methods declared below override the embedded ones to add cacheable-flipping
// or rules logic. GameState owns the data; GameEngine owns the rules. The split lets
// internal machinery pass around a *GameState pointer when it just needs to read or copy
// raw state.
type GameEngine struct {
	*GameState
	// pitchBonus accumulates extra resources a triggertype.Pitch handler grants for the
	// card currently being pitched. FirePitchTriggers zeroes it before the fire and reads
	// it after; it lives on the engine wrapper (not GameState) since it never outlives a
	// single pitch and so must not be copied per permutation.
	pitchBonus int
}

// === Cards-facing zone accessors that flip cacheable. These shadow the same-name
//     methods promoted from *GameState; the embedded versions stay reachable as
//     ge.GameState.X when the engine internals need the non-flipping variant.

// Hand returns the cards in hand — INCLUDING entries tagged Pitch / Attack roles by the
// partition. The in-hand → pitch-zone transition is rules-modelled at the moment a card
// actually pays its cost, so "a card in your hand" reads must see scheduled-but-not-yet-
// committed entries. Flips IsCacheable; for non-mutating size / predicate gates prefer
// HandSize / HandHasMatching / HeldHandSize.
func (ge *GameEngine) Hand() []card.Card {
	ge.cacheable = false
	return ge.GameState.Hand()
}

// HeldHand returns the Held subset of the hand — Pitch / Attack entries are excluded
// (scheduled to commit downstream and not divertible). Drawn entries are included so
// the slice's length agrees with HeldHandSize. Flips IsCacheable since iterating the
// slice exposes each entry's attributes.
func (ge *GameEngine) HeldHand() []card.Card {
	ge.cacheable = false
	return heldHandSlice(ge.GameState.HandStates())
}

// heldHandSlice projects a role-tagged hand down to the Held-role cards (Pitch / Attack
// entries skipped; drawn entries included). Returns nil when there are none.
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

// insertHandSorted inserts c at the position that keeps the hand ordered by Card.ID(),
// flipping IsCacheable to false. Sorting on every insert keeps the hand a canonical
// multiset, which the chain runner and the eval-cache key both depend on.
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

// Deck returns the chain-runner deck for read-only inspection and flips IsCacheable to
// false. Card handlers should not mutate the returned *deck.Deck directly; route
// mutations through PopDeckTop / PrependToDeck / Opt / TutorFromDeck /
// RecycleToDeckBottom.
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

// PopDeckTop removes the top card of the deck and returns it. Returns (nil, false) when
// the deck is empty. Flips IsCacheable to false (the caller observes the popped card's
// identity), and notes the deck removal so the eval cache's depth check stays accurate
// even on uncacheable paths.
func (ge *GameEngine) PopDeckTop() (card.Card, bool) {
	ge.cacheable = false
	if ge.deck.Size() == 0 {
		return nil, false
	}
	ge.noteDeckRemoval(1)
	return ge.deck.Draw(1)[0].(card.Card), true
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

// PrependToDeck inserts c at the top of the deck. The caller supplies c, so the write
// is reproducible from the cache key + chain order; doesn't flip IsCacheable.
func (ge *GameEngine) PrependToDeck(c card.Card) {
	ge.deck.PutTop([]deck.Card{c})
}

// AppendToDeck inserts c at the bottom of the deck. Cache-friendly for the same reason
// as PrependToDeck.
func (ge *GameEngine) AppendToDeck(c card.Card) {
	ge.deck.PutBottom([]deck.Card{c})
}

// Discard pops the first Held-role hand card to the graveyard and logs the action
// under source. Returns true on success; false when no Held card exists. Cache-safe:
// the discarded card's identity never escapes the engine, so the caller can't branch
// on it and the cache can't diverge across replays with different drawn cards.
func (ge *GameEngine) Discard(source string) bool {
	c, ok := ge.popFirstHeldCard()
	if !ok {
		return false
	}
	ge.graveyard = append(ge.graveyard, c)
	ge.logger.AppendPostTriggerf(source, 0, "Discarded a card")
	return true
}

// DiscardToTopOfDeck pops the first Held-role hand card to the top of the deck and
// logs under source. Returns true on success. Cache-safe for the same reason as
// Discard — identity never escapes.
func (ge *GameEngine) DiscardToTopOfDeck(source string) bool {
	c, ok := ge.popFirstHeldCard()
	if !ok {
		return false
	}
	ge.deck.PutTop([]deck.Card{c})
	ge.logger.AppendPostTriggerf(source, 0, "Cycled a card to top of deck")
	return true
}

// DiscardToBottomOfDeck pops the first Held-role hand card to the bottom of the deck
// and logs under source. Returns true on success. Cache-safe.
func (ge *GameEngine) DiscardToBottomOfDeck(source string) bool {
	c, ok := ge.popFirstHeldCard()
	if !ok {
		return false
	}
	ge.deck.PutBottom([]deck.Card{c})
	ge.logger.AppendPostTriggerf(source, 0, "Cycled a card to bottom of deck")
	return true
}

// popFirstHeldCard removes the first Held-role hand entry and returns its Card. Drawn
// entries are eligible — they're indistinguishable to the caller of the no-return
// accessors above, so identity leakage isn't a concern at this layer.
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

// RecycleToDeckBottom appends pc.Card to the bottom of the deck and flags the chain
// dispatcher to skip the usual non-persistent → graveyard append. Models the FaB clause
// "put this on the bottom of its owner's deck". Doesn't flip IsCacheable — the caller
// supplies pc, so the write is reproducible from the cache key + chain order.
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
func (ge *GameEngine) BanishFromGraveyard(pred func(card.Card) bool) (card.Card, bool) {
	ge.cacheable = false
	for i, c := range ge.graveyard {
		if !pred(c) {
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
func (ge *GameEngine) RecycleFromGraveyardToTop(pred func(card.Card) bool) (card.Card, bool) {
	return ge.recycleFromGraveyard(pred, true)
}
func (ge *GameEngine) RecycleFromGraveyardToBottom(pred func(card.Card) bool) (card.Card, bool) {
	return ge.recycleFromGraveyard(pred, false)
}

func (ge *GameEngine) recycleFromGraveyard(pred func(card.Card) bool, toTop bool) (card.Card, bool) {
	ge.cacheable = false
	for i, c := range ge.graveyard {
		if !pred(c) {
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
// caller), so the write is reproducible from the cache key + chain order.
func (ge *GameEngine) AddToGraveyard(c card.Card) {
	ge.graveyard = append(ge.graveyard, c)
}

// DrawOne models a mid-turn draw: pop the top of the deck into the hand at its sorted
// position. Reports whether a card was drawn; false on an empty deck. Doesn't flip
// IsCacheable — the cached chain runs DrawOne again on replay against the caller's
// current deck, so the drawn card's identity is naturally re-resolved per call. Chain
// steps that read the drawn card's attributes via Hand / HeldHand / PeekTopN still
// flip IsCacheable through those accessors, so a Hand-reader downstream of DrawOne
// remains uncacheable.
func (ge *GameEngine) DrawOne() bool {
	if ge.deck == nil || ge.deck.Size() == 0 {
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

// OpponentDiscard credits n cards' worth of damage-equivalent value for forcing the
// opponent to discard. Returns the credited value for log attribution.
func (ge *GameEngine) OpponentDiscard(n int) int {
	v := n * DiscardValue
	ge.AddValue(v)
	return v
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

// Opt resolves the FaB "Opt N" keyword: pops up to n cards from the top of the deck and
// hands them to the current hero's Opt heuristic. The handler returns a (top, bottom)
// split; the top list goes back on top of the deck (in returned order) and the bottom
// list appends to the bottom (in returned order). n is clamped to the current deck
// length. Always flips IsCacheable to false.
//
// Emits a log entry naming the revealed cards and the split when the handler ran.
//
// Panics if the handler's combined output isn't exactly the input multiset.
func (ge *GameEngine) Opt(l card.Logger, n int) {
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

	var top, bottom []card.Card
	if ge.hero == nil {
		top = cards
	} else {
		top, bottom = ge.hero.Opt(cards)
	}
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

// === Trigger and aura dispatch ===

// FireTriggers fires the hero plus every Aura, EphemeralTrigger, and Item registered for
// trigger type t. It is the single dispatch point for every triggertype.Type lifecycle
// event.
//
// triggeringCard is the card whose resolution raised the event, or nil for turn-boundary
// events. It is published on ge.triggeringCard so handlers can attribute log lines, and
// its type set is what each entry's type filter matches against. The hero fires first,
// then auras, ephemeral triggers, and items.
//
// triggeringCard.Types is resolved lazily — only when an entry of the right trigger type
// is actually found — so the common per-card fire with no subscribers stays off the Types
// interface-dispatch path.
func (ge *GameEngine) FireTriggers(t triggertype.Type, triggeringCard card.Card) {
	heroFires := ge.hero != nil && ge.hero.TriggerType()&t != 0
	if !heroFires && len(ge.auras) == 0 && len(ge.triggers) == 0 && len(ge.items) == 0 {
		return
	}
	ge.triggeringCard = triggeringCard

	var triggeringTypes card.TypeSet
	typesResolved := false
	matchTypes := func() card.TypeSet {
		if !typesResolved {
			if triggeringCard != nil {
				triggeringTypes = triggeringCard.Types(ge)
			}
			typesResolved = true
		}
		return triggeringTypes
	}

	if heroFires {
		fireHero(ge, triggeringCard, matchTypes)
	}
	fireHooks(ge, &ge.auras, t, triggeringCard, matchTypes, false)
	fireHooks(ge, &ge.triggers, t, triggeringCard, matchTypes, true)
	fireHooks(ge, &ge.items, t, triggeringCard, matchTypes, false)

	ge.triggeringCard = nil
}

// fireHero applies the OncePerTurn / Matches gates and invokes the hero handler. The
// hero is singular (no slice splicing, no removeAfterFire) so it bypasses fireHooks's
// cursor walk. The TriggerType bit-and check is done at the caller.
func fireHero(ge *GameEngine, triggeringCard card.Card, matchTypes func() card.TypeSet) {
	h := ge.hero
	if h.OncePerTurn() && h.FiredThisTurn() {
		return
	}
	if triggeringCard != nil && !h.Matches(matchTypes()) {
		return
	}
	h.Fire(ge, ge.logger)
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
func fireHooks[H trigger.Hook](ge *GameEngine, hooks *[]H, t triggertype.Type, triggeringCard card.Card, matchTypes func() card.TypeSet, removeAfterFire bool) {
	n := len(*hooks)
	for i := 0; i < n; {
		h := (*hooks)[i]
		if h.TriggerType()&t == 0 || (h.OncePerTurn() && h.FiredThisTurn()) ||
			(triggeringCard != nil && !h.Matches(matchTypes())) {
			i++
			continue
		}
		ge.currentHookIdx = i
		ge.currentHookDestroyed = false
		h.Fire(ge, ge.logger)
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

// DestroyAura removes the aura currently being fired from the arena. It then fires the
// source card's OnLeavesArena hook — the printed "when this leaves the arena" clause, when
// the card implements card.LeavesArenaAura — and, when addToGraveyard==true, pushes the
// source card into the graveyard (no-op for token auras with no source). OnLeavesArena runs
// before the graveyard append so a "banish another aura from your graveyard" leave clause
// doesn't see the just-left card. Direct splice with no cacheable flip — destruction is
// deterministic from the triggering event.
func (ge *GameEngine) DestroyAura(addToGraveyard bool) {
	i := ge.currentHookIdx
	if i < 0 || i >= len(ge.auras) {
		return
	}
	src := ge.auras[i].SourceCard()
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
// pitched card's contribution on top of its printed Pitch value.
func (ge *GameEngine) FirePitchTriggers(pitched card.Card) int {
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
		ge.auras = append(ge.auras[:i], ge.auras[i+1:]...)
		la.OnLeavesArena(ge, ge.logger)
		ge.AppendGraveyard(src)
		return true
	}
	return false
}

// === Chain-step resolution ===

// ResolveChainStep runs card.Play on pc and then applies the standard chain-step
// resolution: attack-action / weapon-attack credit pc.EffectiveAttack() to ge.value;
// defense-reaction (or DefensiveInstant) credits EffectiveDefense capped at the remaining
// unblocked damage; everything else logs (+0). The "<DisplayName>: <VERB> (+N)" chain-step
// entry is appended after Play returns so self-buffs Play applied are reflected in the
// displayed delta.
func (ge *GameEngine) ResolveChainStep(l card.Logger, pc *card.CardState) {
	pc.Card.Play(ge, l, pc)
	types := pc.Card.Types(nil)
	if types.Has(card.TypeAura) {
		ge.auraCreated = true
	}
	n := ge.chainStepDelta(pc, types)
	l.AppendChainStep(ChainStepText(pc), n)
}

// PlayCard implements card.GameEngine.PlayCard — resolves another card mid-handler.
func (ge *GameEngine) PlayCard(l card.Logger, pc *card.CardState) {
	ge.ResolveChainStep(l, pc)
}

// chainStepDelta computes the chain step's display delta and applies the standard damage /
// block side effects. Returns the (+N) value for the log line. types is the caller's
// already-resolved Types(nil) so we skip a second interface dispatch.
func (ge *GameEngine) chainStepDelta(pc *card.CardState, types card.TypeSet) int {
	switch {
	case types.IsAttackAction() || types.IsWeaponAttack():
		n := pc.EffectiveAttack()
		ge.value += n
		return n
	case types.IsDefenseReaction() || isDefensiveInstant(pc.Card):
		n := pc.EffectiveDefense()
		if rem := ge.incomingDamage - ge.damageBlocked; n > rem {
			n = rem
		}
		if n < 0 {
			n = 0
		}
		ge.damageBlocked += n
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

// ChainStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix for the chain-
// step log line. VERB picks WEAPON ATTACK for Weapon+Attack, ATTACK for attack-actions,
// DEFENSE REACTION for DRs, and PLAY otherwise. Declared as a var so a memoised
// implementation can be swapped in at init.
var ChainStepText = func(pc *card.CardState) string {
	types := pc.Card.Types(nil)
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

// === Arcane damage ===

// DealArcaneDamage credits n arcane damage to Value, writes a "Dealt n arcane damage" rider
// line under source, and flips ArcaneDamageDealt when LikelyDamageHits(n, false) approves
// so same-turn triggers reading "if you've dealt arcane damage this turn" fire. Routes
// through dealtArcaneText to avoid per-call fmt.Sprintf and variadic-int boxing. Any
// positive arcane damage also clears OpponentMarked — mark is consumed by any damage the
// marked hero takes, gated the same way as the physical-attack clear (positive damage,
// not LikelyDamageHits).
func (ge *GameEngine) DealArcaneDamage(l card.Logger, source string, n int) {
	ge.AddValue(n)
	if ge.LikelyDamageHits(n, false) {
		ge.arcaneDamageDealt = true
	}
	if n > 0 {
		ge.opponentMarked = false
	}
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

// Card-facing token creation / count methods on *GameEngine. Live tokens are identified by
// CardName (the canonical display name); concrete Aura / Item types live outside gameengine
// and are produced by the token-package factories used below.

// Token display names — the engine matches by CardName when bumping an existing entry's
// Count or reading a count.
const (
	tokenNameRunechant = "Runechant"
	tokenNamePonder    = "Ponder"
	tokenNameGold      = "Gold"
	tokenNameSilver    = "Silver"
	tokenNameCopper    = "Copper"
)

// CreateRunechants creates n Runechant tokens and credits +n damage at creation.
// Tokens are stored as a single Aura entry — bump an existing entry's Count or add a
// new one. Sets AuraCreated so same-turn "aura created this turn" effects see it.
func (ge *GameEngine) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	ge.AddValue(n)
	bumpOrCreateAura(ge.GameState, tokenNameRunechant, func(n int) Aura { return token.NewRunechant(n) }, n)
}

// CreatePonders creates n Ponder tokens. No Value credit — Ponder pays out at end of turn.
func (ge *GameEngine) CreatePonders(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateAura(ge.GameState, tokenNamePonder, func(n int) Aura { return token.NewPonder(n) }, n)
}

// CreateGold / CreateSilver / CreateCopper create the matching token items. No Value
// credit — items only pay out when the activated ability is spent.
func (ge *GameEngine) CreateGold(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateItem(ge.GameState, tokenNameGold, func(n int) Item { return token.NewGold(n) }, n)
}
func (ge *GameEngine) CreateSilver(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateItem(ge.GameState, tokenNameSilver, func(n int) Item { return token.NewSilver(n) }, n)
}
func (ge *GameEngine) CreateCopper(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateItem(ge.GameState, tokenNameCopper, func(n int) Item { return token.NewCopper(n) }, n)
}

// RunechantCount / PonderCount / GoldCount / SilverCount / CopperCount return the
// live count of each token kind in play, or zero when none.
func (ge *GameEngine) RunechantCount() int { return auraCountByName(ge.auras, tokenNameRunechant) }
func (ge *GameEngine) PonderCount() int    { return auraCountByName(ge.auras, tokenNamePonder) }
func (ge *GameEngine) GoldCount() int      { return itemCountByName(ge.items, tokenNameGold) }
func (ge *GameEngine) SilverCount() int    { return itemCountByName(ge.items, tokenNameSilver) }
func (ge *GameEngine) CopperCount() int    { return itemCountByName(ge.items, tokenNameCopper) }

// bumpOrCreateAura increments an existing aura entry matching name on s, or appends
// a new one built by build(n). Flips gs.auraCreated.
func bumpOrCreateAura(s *GameState, name string, build func(int) Aura, n int) {
	s.auraCreated = true
	for i := range s.auras {
		if s.auras[i].CardName() == name {
			s.auras[i].SetCount(s.auras[i].Count() + n)
			return
		}
	}
	s.auras = append(s.auras, build(n))
}

// bumpOrCreateItem increments an existing item entry matching name on s, or appends
// a fresh one built by build(n). Items don't flip auraCreated.
func bumpOrCreateItem(s *GameState, name string, build func(int) Item, n int) {
	for i := range s.items {
		if s.items[i].CardName() == name {
			s.items[i].SetCount(s.items[i].Count() + n)
			return
		}
	}
	s.items = append(s.items, build(n))
}

// ConsumeItemByName decrements the matching item's Count by n and removes the entry when
// Count reaches zero. Token items don't head to the graveyard on destroy. No-op when no
// item matches name.
func (ge *GameEngine) ConsumeItemByName(name string, n int) {
	for i := range ge.items {
		if ge.items[i].CardName() != name {
			continue
		}
		newCount := ge.items[i].Count() - n
		if newCount <= 0 {
			ge.items = append(ge.items[:i], ge.items[i+1:]...)
		} else {
			ge.items[i].SetCount(newCount)
		}
		return
	}
}

// auraCountByName scans auras for a token aura by display name.
func auraCountByName(auras []Aura, name string) int {
	for _, a := range auras {
		if a.CardName() == name {
			return a.Count()
		}
	}
	return 0
}

// itemCountByName scans items for a token item by display name.
func itemCountByName(items []Item, name string) int {
	for _, i := range items {
		if i.CardName() == name {
			return i.Count()
		}
	}
	return 0
}
