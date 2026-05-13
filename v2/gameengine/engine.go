package gameengine

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// GameEngine is the rules-engine wrapper around a *GameState. Cards play against this
// type via the v2/card.GameEngine interface; the engine's method surface mixes
// (a) cacheable-aware accessors that flip cacheable as a side effect of touching hidden
// state and (b) rules-orchestration methods (Fire*, ResolveChainStep, Opt, Clash,
// DealArcaneDamage, token economy) that apply the FaB rule set on top of the raw state
// mutations the state itself exposes.
//
// *GameState is embedded so every pure state accessor (Hand, Deck, Auras, SetX, …) and
// the Copy / Reset utilities promote automatically. Methods listed below
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

// === Rules-orchestration methods. Each operates on the embedded *GameState's slices
// but applies game-rule semantics (cursor iteration for handler-side splices,
// OncePerTurn gating, FiredThisTurn accounting, post-fire trigger drainage).

// FireAttack walks the aura entries with TriggerType()==triggertype.Attack and invokes
// every one whose OncePerTurn gate is open. The triggering card is published on
// triggeringCard so handlers can attribute log lines back to the source. Cursor-based
// iteration so a handler-side splice (Destroy) advances only when the slice length
// didn't change.
func (g *GameEngine) FireAttack(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, triggertype.Attack)
}

// FireAttackAction is the triggertype.AttackAction counterpart to FireAttack: walks
// the aura entries matching that type and fires those whose OncePerTurn gate is open.
func (g *GameEngine) FireAttackAction(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, triggertype.AttackAction)
}

// fireMatching is the shared aura-fire walk for FireAttack / FireAttackAction /
// FireEndOfTurn. Iterates auras with a cursor so handler-side splicing (Destroy
// mutates the auras slice in place, shifting the next entry down to the cursor's
// index) advances only when the slice length didn't change.
func (g *GameEngine) fireMatching(triggeringCard card.Card, trigger triggertype.Type) {
	for i := 0; i < len(g.auras); {
		a := g.auras[i]
		if a.TriggerType() != trigger || (a.OncePerTurn() && a.FiredThisTurn()) {
			i++
			continue
		}
		g.triggeringCard = triggeringCard
		g.currentAuraIdx = i
		g.currentAuraDestroyed = false
		a.Fire(g, g.logger)
		g.currentAuraIdx = -1
		g.triggeringCard = nil
		if !g.currentAuraDestroyed {
			g.auras[i].SetFiredThisTurn(true)
			i++
		}
	}
}

// HasEndOfTurnFire reports whether either Auras or Triggers carries a
// triggertype.EndOfTurn entry. Lets the chain runner skip the end-of-turn walk when
// nothing would fire.
func (g *GameEngine) HasEndOfTurnFire() bool {
	for _, a := range g.auras {
		if a.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	for _, t := range g.triggers {
		if t.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	return false
}

// FireEndOfTurn runs after the chain has finished resolving (and the legality gates
// have passed) but before the carry state is captured. Walks Auras and Triggers in one
// pass each:
//
//   - Aura entries respect OncePerTurn / FiredThisTurn semantics; the handler owns
//     destruction via the engine's destroyAura path.
//   - Trigger entries are one-shot; fired entries are removed afterward. Snapshotting
//     len(g.triggers) before iterating keeps a handler that calls AddXxxTrigger from
//     firing its newcomer on the same pass — newcomers stay queued for the next
//     matching event.
func (g *GameEngine) FireEndOfTurn() {
	g.fireMatching(nil, triggertype.EndOfTurn)
	n := len(g.triggers)
	for i := 0; i < n; i++ {
		tr := g.triggers[i]
		if tr.TriggerType() != triggertype.EndOfTurn {
			continue
		}
		tr.Fire(g, g.logger)
	}
	kept := g.triggers[:0]
	for i, tr := range g.triggers {
		if i < n && tr.TriggerType() == triggertype.EndOfTurn {
			continue
		}
		kept = append(kept, tr)
	}
	g.triggers = kept
}

// FireHit walks the one-shot trigger queue and invokes every triggertype.Hit entry
// whose type filter matches the attacking card's types. Surviving entries (filter
// mismatch) are kept; fired entries are removed.
func (g *GameEngine) FireHit(attackerTypes card.TypeSet) {
	kept := g.triggers[:0]
	for i := range g.triggers {
		t := g.triggers[i]
		if t.TriggerType() != triggertype.Hit || !t.Matches(attackerTypes) {
			kept = append(kept, t)
			continue
		}
		t.Fire(g, g.logger)
	}
	g.triggers = kept
}

// FireStartOfTurn walks g.auras and invokes every triggertype.StartOfTurn entry,
// calling onFire with each entry's pre-state snapshot so sim can attribute damage /
// draws / log lines back to the firing aura. Auras that destroy themselves splice
// out; FiredThisTurn flips reset on each fresh turn boundary.
//
// The onFire callback receives:
//   - pre is the index in g.auras of the firing entry at the time of the call.
//   - damage is g.value's delta during this handler — the partition tiebreaker uses
//     it.
//   - drawnCard is the first card the handler appended to hand, or nil. Used by
//     processAurasAtStartOfTurn to surface "revealed" entries.
//   - newLogEntries is the slice of LogEntries the handler appended (caller may copy
//     out).
func (g *GameEngine) FireStartOfTurn(onFire func(idx int, damage int, drawnCard card.Card, newLogEntries []turnlogger.LogEntry)) {
	for i := 0; i < len(g.auras); {
		a := g.auras[i]
		g.auras[i].SetFiredThisTurn(false)
		if a.TriggerType() != triggertype.StartOfTurn {
			i++
			continue
		}
		preHand := len(g.hand)
		preLog := 0
		if g.logger != nil {
			preLog = len(g.logger.Entries())
		}
		preValue := g.value
		g.currentAuraIdx = i
		g.currentAuraDestroyed = false
		a.Fire(g, g.logger)
		g.currentAuraIdx = -1

		damage := g.value - preValue
		var drawn card.Card
		if len(g.hand) > preHand {
			drawn = g.hand[preHand]
		}
		var newEntries []turnlogger.LogEntry
		if g.logger != nil {
			if entries := g.logger.Entries(); len(entries) > preLog {
				newEntries = entries[preLog:]
			}
		}
		if onFire != nil {
			onFire(i, damage, drawn, newEntries)
		}
		if !g.currentAuraDestroyed {
			i++
		}
	}
}

// AdvanceTurnBoundary clears the per-turn FiredThisTurn flag on every persisted aura.
// The chain runner calls this when advancing across the turn boundary so the
// OncePerTurn gate rearms.
func (g *GameEngine) AdvanceTurnBoundary() {
	for i := range g.auras {
		g.auras[i].SetFiredThisTurn(false)
	}
}

// DestroyAura removes the aura currently being fired and, when addToGraveyard==true,
// invokes the aura's OnDestroy hook to push the aura's source card into the graveyard
// (token auras no-op). Direct splice (no cacheable flip) — destruction is
// deterministic from the triggering event, not hidden state.
//
// Called by the card.Aura context the engine threads into each handler; cards do not
// call this directly.
func (g *GameEngine) DestroyAura(addToGraveyard bool) {
	i := g.currentAuraIdx
	if i < 0 || i >= len(g.auras) {
		return
	}
	if addToGraveyard {
		g.auras[i].OnDestroy(g)
	}
	g.auras = append(g.auras[:i], g.auras[i+1:]...)
	g.currentAuraDestroyed = true
}

// === Chain-step resolution ===

// ResolveChainStep runs card.Play on self and then applies the standard chain-step
// resolution: for an attack-action or weapon-attack, credit self.EffectiveAttack() to
// g.value; for a defense-reaction (or DefensiveInstant), credit the EffectiveDefense
// capped at IncomingDamage and decrement IncomingDamage; for everything else, log
// (+0). The "<DisplayName>: <VERB> (+N)" chain-step entry is appended after Play
// returns so any self-buffs Play applied (e.g. modal +2{p} riders flipping
// self.BonusAttack) are reflected in the displayed delta.
//
// Cards' Play body owns card-specific behaviour: riders that emit rider log lines,
// OnHit registration, conditional self-buffs, sub-card plays. The standard
// printed-attack-deals-damage / DR-blocks-incoming mechanic is the engine's job;
// cards don't reach for DealEffectiveAttack / DealEffectiveDefense or emit the chain
// step themselves.
func (g *GameEngine) ResolveChainStep(l card.Logger, self *card.CardState) {
	self.Card.Play(g, l, self)
	if self.Card.Types(nil).Has(card.TypeAura) {
		g.auraCreated = true
	}
	n := g.chainStepDelta(self)
	l.AppendChainStep(ChainStepText(self), n)
}

// PlayCard implements card.GameEngine.PlayCard. Cards reach this when they resolve
// another card mid-handler (Moon Wish tutoring Sun Kiss into play on go-again).
func (g *GameEngine) PlayCard(l card.Logger, self *card.CardState) {
	g.ResolveChainStep(l, self)
}

// chainStepDelta computes the chain step's display delta and applies the standard
// damage / block side effects. Returns the (+N) value for the log line.
func (g *GameEngine) chainStepDelta(self *card.CardState) int {
	types := self.Card.Types(nil)
	switch {
	case types.IsAttackAction() || types.IsWeaponAttack():
		n := self.EffectiveAttack()
		g.value += n
		return n
	case types.IsDefenseReaction() || isDefensiveInstant(self.Card):
		n := self.EffectiveDefense()
		if n > g.incomingDamage {
			n = g.incomingDamage
		}
		if n < 0 {
			n = 0
		}
		g.incomingDamage -= n
		g.value += n
		return n
	}
	return 0
}

// isDefensiveInstant reports whether c opts into the DR resolution path via the
// DefensiveInstant marker. Centralised here so ResolveChainStep doesn't repeat the
// type-assertion shape.
func isDefensiveInstant(c card.Card) bool {
	_, ok := c.(card.DefensiveInstant)
	return ok
}

// ChainStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix the chain-
// step log line is built from. VERB picks WEAPON ATTACK for weapon activated-ability
// cards (Weapon + Attack), ATTACK for attack-action cards, DEFENSE REACTION for
// Defense Reactions, and PLAY for everything else; the "from arsenal" suffix tags
// entries played out of the arsenal slot. Declared as a var so internal/optimizations
// can swap in a memoised per-(CardID, FromArsenal) implementation at init.
var ChainStepText = func(self *card.CardState) string {
	types := self.Card.Types(nil)
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
	if self.FromArsenal {
		verb += " from arsenal"
	}
	return self.Card.DisplayName() + ": " + verb
}

// === Arcane damage ===

// DealArcaneDamage credits n arcane damage to Value, writes a "Dealt n arcane damage"
// rider line under source, and flips ArcaneDamageDealt when LikelyDamageHits(n, false)
// approves so same-turn triggers reading "if you've dealt arcane damage this turn"
// fire. Routes through dealtArcaneText[n] so the hot path avoids per-call fmt.Sprintf
// and variadic-int boxing.
func (g *GameEngine) DealArcaneDamage(l card.Logger, source string, n int) {
	g.AddValue(n)
	if g.LikelyDamageHits(n, false) {
		g.arcaneDamageDealt = true
	}
	if n >= 0 && n < len(dealtArcaneText) {
		l.AppendPostTrigger(source, dealtArcaneText[n], n)
		return
	}
	l.AppendPostTriggerf(source, n, "Dealt %d arcane damage", n)
}

// dealtArcaneText is the pre-built rider-line cache indexed by arcane-damage count,
// keeping DealArcaneDamage alloc-free on the hot path. Extend if a new card prints
// higher arcane.
var dealtArcaneText = [...]string{
	0: "Dealt 0 arcane damage",
	1: "Dealt 1 arcane damage",
	2: "Dealt 2 arcane damage",
	3: "Dealt 3 arcane damage",
	4: "Dealt 4 arcane damage",
}

// === Trigger registration ===

// AddHitTrigger registers a one-shot triggertype.Hit listener. filter narrows the
// qualifying hits to a card-type predicate; nil = any hit qualifies.
func (g *GameEngine) AddHitTrigger(self *card.CardState, handler card.TriggerHandler, filter func(card.TypeSet) bool) {
	g.CreateTrigger(BuildCardTrigger(self, triggertype.Hit, handler, filter))
}

// AddEndOfTurnTrigger registers a one-shot triggertype.EndOfTurn listener — fires
// after the chain finishes resolving but before the carry-state snapshot.
func (g *GameEngine) AddEndOfTurnTrigger(self *card.CardState, handler card.TriggerHandler) {
	g.CreateTrigger(BuildCardTrigger(self, triggertype.EndOfTurn, handler, nil))
}

// === Tokens ===

// Card-facing token creation / count methods on *GameEngine. v2/card.GameEngine
// requires these for cards to make / read FaB's five built-in tokens (Runechant,
// Ponder, Gold, Silver, Copper). The engine identifies live tokens by their CardName
// (the public canonical display name) and delegates aura / item construction to the
// registered Build*Aura / Build*Item factories so the concrete types live outside
// gameengine.

// Token display names — the engine matches by CardName when bumping an existing
// entry's Count or reading a count. Concrete Aura / Item impls report these strings
// via CardName.
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
func (g *GameEngine) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	g.AddValue(n)
	bumpOrCreateAura(g.GameState, tokenNameRunechant, BuildRunechantAura, n)
}

// CreatePonders creates n Ponder tokens. No Value credit — Ponder pays out at end of
// turn (see the runtime's Ponder aura handler).
func (g *GameEngine) CreatePonders(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateAura(g.GameState, tokenNamePonder, BuildPonderAura, n)
}

// CreateGold / CreateSilver / CreateCopper create the matching token items. No Value
// credit — items only pay out when the activated ability is spent.
func (g *GameEngine) CreateGold(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateItem(g.GameState, tokenNameGold, BuildGoldItem, n)
}
func (g *GameEngine) CreateSilver(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateItem(g.GameState, tokenNameSilver, BuildSilverItem, n)
}
func (g *GameEngine) CreateCopper(n int) {
	if n <= 0 {
		return
	}
	bumpOrCreateItem(g.GameState, tokenNameCopper, BuildCopperItem, n)
}

// RunechantCount / PonderCount / GoldCount / SilverCount / CopperCount return the
// live count of each token kind in play, or zero when none.
func (g *GameEngine) RunechantCount() int { return auraCountByName(g.auras, tokenNameRunechant) }
func (g *GameEngine) PonderCount() int    { return auraCountByName(g.auras, tokenNamePonder) }
func (g *GameEngine) GoldCount() int      { return itemCountByName(g.items, tokenNameGold) }
func (g *GameEngine) SilverCount() int    { return itemCountByName(g.items, tokenNameSilver) }
func (g *GameEngine) CopperCount() int    { return itemCountByName(g.items, tokenNameCopper) }

// bumpOrCreateAura increments an existing aura entry matching name on s, or appends
// a fresh one built by build(n). Flips s.auraCreated.
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

// ConsumeItemByName decrements the matching item's Count by n and removes the entry
// when Count reaches zero. Token items don't head to the graveyard on destroy. No-op
// when no item matches name. Called by token-ability Play implementations registered
// outside gameengine.
func (g *GameEngine) ConsumeItemByName(name string, n int) {
	for i := range g.items {
		if g.items[i].CardName() != name {
			continue
		}
		newCount := g.items[i].Count() - n
		if newCount <= 0 {
			g.items = append(g.items[:i], g.items[i+1:]...)
		} else {
			g.items[i].SetCount(newCount)
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
