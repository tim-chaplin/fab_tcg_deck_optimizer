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
}

// === Cards-facing zone accessors that flip cacheable. These shadow the same-name
//     methods promoted from *GameState; the embedded versions stay reachable as
//     ge.GameState.X when the engine internals need the non-flipping variant.

// Hand returns the live hand slice and flips IsCacheable to false. Cards must not mutate
// the returned slice; use AppendHand / PopHandAt for mutations.
func (ge *GameEngine) Hand() []card.Card {
	ge.cacheable = false
	return ge.hand
}

// AppendHand inserts c into the hand, flipping IsCacheable to false. The hand is kept
// sorted by Card.ID() at all times — see insertHandSorted.
func (ge *GameEngine) AppendHand(c card.Card) {
	ge.insertHandSorted(c)
}

// insertHandSorted places c at the position that keeps the hand ordered by Card.ID(), and
// flips IsCacheable to false. Keeping the hand sorted by construction means every draw —
// start-of-turn aura, mid-chain DrawOne, end-of-turn refill — lands a canonical multiset,
// so the chain runner and the eval cache never need a separate normalising sort.
func (ge *GameEngine) insertHandSorted(c card.Card) {
	ge.cacheable = false
	i := sort.Search(len(ge.hand), func(j int) bool { return ge.hand[j].ID() > c.ID() })
	ge.hand = append(ge.hand, nil)
	copy(ge.hand[i+1:], ge.hand[i:])
	ge.hand[i] = c
}

// PopHandAt removes and returns the card at index i, flipping IsCacheable to false.
func (ge *GameEngine) PopHandAt(i int) card.Card {
	ge.cacheable = false
	c := ge.hand[i]
	ge.hand = append(ge.hand[:i], ge.hand[i+1:]...)
	return c
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
// the deck is empty. Flips IsCacheable to false.
func (ge *GameEngine) PopDeckTop() (card.Card, bool) {
	ge.cacheable = false
	if ge.deck.Size() == 0 {
		return nil, false
	}
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

// PrependToDeck inserts c at the top of the deck. Flips IsCacheable to false.
func (ge *GameEngine) PrependToDeck(c card.Card) {
	ge.cacheable = false
	ge.deck.PutTop([]deck.Card{c})
}

// AppendToDeck inserts c at the bottom of the deck. Flips IsCacheable to false.
func (ge *GameEngine) AppendToDeck(c card.Card) {
	ge.cacheable = false
	ge.deck.PutBottom([]deck.Card{c})
}

// Discard pops the first hand card and appends it to the graveyard. Returns the discarded
// card and true; returns (nil, false) when the hand is empty. Flips IsCacheable via
// PopHandAt.
func (ge *GameEngine) Discard() (card.Card, bool) {
	if len(ge.hand) == 0 {
		return nil, false
	}
	c := ge.PopHandAt(0)
	ge.graveyard = append(ge.graveyard, c)
	return c, true
}

// RecycleToDeckBottom appends pc.Card to the bottom of the deck and flags the chain
// dispatcher to skip the usual non-persistent → graveyard append. Models the FaB clause
// "put this on the bottom of its owner's deck". Flips IsCacheable.
func (ge *GameEngine) RecycleToDeckBottom(pc *card.CardState) {
	ge.cacheable = false
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

// AddToGraveyard appends c to graveyard so later-resolving cards see it. Card-facing entry
// point — flips IsCacheable to false. Framework-internal callers use the promoted
// AppendGraveyard, which appends without the cacheable flip.
func (ge *GameEngine) AddToGraveyard(c card.Card) {
	ge.cacheable = false
	ge.graveyard = append(ge.graveyard, c)
}

// DrawOne models a mid-turn draw: pop the top of the deck into the hand at its sorted
// position. No-op on an empty deck. Inherits the IsCacheable flip via PopDeckTop.
func (ge *GameEngine) DrawOne() {
	c, ok := ge.PopDeckTop()
	if !ok {
		return
	}
	ge.insertHandSorted(c)
}

// PonderDrawOne pops the deck top into the hand at its sorted position and returns false
// on an empty deck.
func (ge *GameEngine) PonderDrawOne() bool {
	c, ok := ge.PopDeckTop()
	if !ok {
		return false
	}
	ge.insertHandSorted(c)
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

// HasEndOfTurnFire reports whether either Auras or Triggers carries a
// triggertype.EndOfTurn entry. Lets the chain runner skip the end-of-turn walk when
// nothing would fire.
func (ge *GameEngine) HasEndOfTurnFire() bool {
	for _, a := range ge.auras {
		if a.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	for _, t := range ge.triggers {
		if t.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	return false
}

// FireTriggers fires every Aura and one-shot Trigger registered for trigger type t. It is
// the single dispatch point for every triggertype.Type lifecycle event.
//
// triggeringCard is the card whose resolution raised the event, or nil for turn-boundary
// events. It is published on ge.triggeringCard so handlers can attribute log lines, and
// its type set is what a Trigger's type filter matches against.
//
// Auras fire in a cursor walk so a handler-side Destroy splice doesn't skip the next
// entry; an open OncePerTurn gate is required and FiredThisTurn is set after a fire. The
// gate is re-armed at the turn boundary by ResetEphemeralState, not here.
//
// Triggers are one-shot. The queue length is snapshotted before firing so a handler that
// queues a new trigger doesn't fire it on the same pass; fired entries are dropped after.
func (ge *GameEngine) FireTriggers(t triggertype.Type, triggeringCard card.Card) {
	ge.triggeringCard = triggeringCard

	for i := 0; i < len(ge.auras); {
		a := ge.auras[i]
		if a.TriggerType() != t || (a.OncePerTurn() && a.FiredThisTurn()) {
			i++
			continue
		}
		ge.currentAuraIdx = i
		ge.currentAuraDestroyed = false
		a.Fire(ge, ge.logger)
		if !ge.currentAuraDestroyed {
			ge.auras[i].SetFiredThisTurn(true)
			i++
		}
	}
	ge.currentAuraIdx = -1

	if n := len(ge.triggers); n > 0 {
		var triggeringTypes card.TypeSet
		if triggeringCard != nil {
			triggeringTypes = triggeringCard.Types(ge)
		}
		firedAny := false
		for i := 0; i < n; i++ {
			tr := ge.triggers[i]
			if tr.TriggerType() != t || !tr.Matches(triggeringTypes) {
				continue
			}
			tr.Fire(ge, ge.logger)
			firedAny = true
		}
		if firedAny {
			kept := ge.triggers[:0]
			for i, tr := range ge.triggers {
				if i < n && tr.TriggerType() == t && tr.Matches(triggeringTypes) {
					continue
				}
				kept = append(kept, tr)
			}
			ge.triggers = kept
		}
	}

	ge.triggeringCard = nil
}

// DestroyAura removes the aura currently being fired and, when addToGraveyard==true, pushes
// the aura's source card into the graveyard (no-op for token auras with no source). Direct
// splice with no cacheable flip — destruction is deterministic from the triggering event.
func (ge *GameEngine) DestroyAura(addToGraveyard bool) {
	i := ge.currentAuraIdx
	if i < 0 || i >= len(ge.auras) {
		return
	}
	if addToGraveyard {
		if src := ge.auras[i].SourceCard(); src != nil {
			ge.AppendGraveyard(src.(card.Card))
		}
	}
	ge.auras = append(ge.auras[:i], ge.auras[i+1:]...)
	ge.currentAuraDestroyed = true
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
// through dealtArcaneText to avoid per-call fmt.Sprintf and variadic-int boxing.
func (ge *GameEngine) DealArcaneDamage(l card.Logger, source string, n int) {
	ge.AddValue(n)
	if ge.LikelyDamageHits(n, false) {
		ge.arcaneDamageDealt = true
	}
	if n >= 0 && n < len(dealtArcaneText) {
		l.AppendPostTrigger(source, dealtArcaneText[n], n)
		return
	}
	l.AppendPostTriggerf(source, n, "Dealt %d arcane damage", n)
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

// === Trigger registration ===

// AddHitTrigger registers a one-shot triggertype.Hit listener. filter narrows the qualifying
// hits to a card-type predicate; nil = any hit qualifies.
func (ge *GameEngine) AddHitTrigger(pc *card.CardState, handler func(card.GameEngine, card.Logger, card.Trigger), filter func(card.TypeSet) bool) {
	ge.CreateTrigger(trigger.NewFromCard(pc.Card, triggertype.Hit, handler, filter))
}

// AddEndOfTurnTrigger registers a one-shot triggertype.EndOfTurn listener — fires
// after the chain finishes resolving but before the carry-state snapshot.
func (ge *GameEngine) AddEndOfTurnTrigger(pc *card.CardState, handler func(card.GameEngine, card.Logger, card.Trigger)) {
	ge.CreateTrigger(trigger.NewFromCard(pc.Card, triggertype.EndOfTurn, handler, nil))
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
// a fresh one built by build(n). Flips gs.auraCreated.
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
