package gameengine

import (
	"fmt"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/turnlogger"
)

// GameEngine is the rules-engine wrapper around a *GameState. Cards play against this
// type via the internal/card.GameEngine interface; the engine's method surface mixes
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
//     ge.GameState.X when the engine internals need the non-flipping variant.

// Hand returns the live hand slice and flips IsCacheable to false. Cards must not mutate
// the returned slice; use AppendHand / PopHandAt for mutations.
func (ge *GameEngine) Hand() []card.Card {
	ge.cacheable = false
	return ge.hand
}

// AppendHand appends c to the hand, flipping IsCacheable to false.
func (ge *GameEngine) AppendHand(c card.Card) {
	ge.cacheable = false
	ge.hand = append(ge.hand, c)
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

// RecycleToDeckBottom appends self.Card to the bottom of the deck and flags the chain
// dispatcher to skip the usual non-persistent → graveyard append. Models the FaB clause
// "put this on the bottom of its owner's deck" (Relentless Pursuit). Flips IsCacheable.
func (ge *GameEngine) RecycleToDeckBottom(self *card.CardState) {
	ge.cacheable = false
	ge.deck.PutBottom([]deck.Card{self.Card})
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

// AddToGraveyard appends c to graveyard so later-resolving cards see it. Used by cards
// running a mini-dispatcher inline (Moon Wish's go-again Sun Kiss play). Flips
// IsCacheable to false. The promoted AppendGraveyard does the same append without the
// flip — framework-internal callers (the chain dispatcher's "non-persistent →
// graveyard" rule, Aura.OnDestroy) reach for that one instead.
func (ge *GameEngine) AddToGraveyard(c card.Card) {
	ge.cacheable = false
	ge.graveyard = append(ge.graveyard, c)
}

// DrawOne models a mid-turn draw: pop the top of the deck and append it to Hand. No-op
// on an empty deck. Bumps CardsDrawn so the partition tiebreaker can prefer chains with
// more draws. Inherits the IsCacheable flip via PopDeckTop.
func (ge *GameEngine) DrawOne() {
	c, ok := ge.PopDeckTop()
	if !ok {
		return
	}
	ge.hand = append(ge.hand, c)
	ge.cardsDrawn++
}

// PonderDrawOne pops the deck top into the hand without bumping CardsDrawn. Returns false
// when the deck is empty. Used by the Ponder token aura at end of turn — ponder draws
// aren't "draws" in the partition-tiebreaker sense, so the counter stays put.
func (ge *GameEngine) PonderDrawOne() bool {
	c, ok := ge.PopDeckTop()
	if !ok {
		return false
	}
	ge.hand = append(ge.hand, c)
	return true
}

// === Rules-engine helpers cards reach through GameEngine ===

// LikelyToHit reports whether self's attack is likely to land past the opponent's blocks.
func (ge *GameEngine) LikelyToHit(self *card.CardState) bool { return LikelyToHit(self) }

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

// === Rules-orchestration methods. Each operates on the embedded *GameState's slices
// but applies game-rule semantics (cursor iteration for handler-side splices,
// OncePerTurn gating, FiredThisTurn accounting, post-fire trigger drainage).

// FireAttack walks the aura entries with TriggerType()==triggertype.Attack and invokes
// every one whose OncePerTurn gate is open. The triggering card is published on
// triggeringCard so handlers can attribute log lines back to the source. Cursor-based
// iteration so a handler-side splice (Destroy) advances only when the slice length
// didn't change.
func (ge *GameEngine) FireAttack(triggeringCard card.Card) {
	ge.fireMatching(triggeringCard, triggertype.Attack)
}

// FireAttackAction is the triggertype.AttackAction counterpart to FireAttack: walks
// the aura entries matching that type and fires those whose OncePerTurn gate is open.
func (ge *GameEngine) FireAttackAction(triggeringCard card.Card) {
	ge.fireMatching(triggeringCard, triggertype.AttackAction)
}

// fireMatching is the shared aura-fire walk for FireAttack / FireAttackAction /
// FireEndOfTurn. Iterates auras with a cursor so handler-side splicing (Destroy
// mutates the auras slice in place, shifting the next entry down to the cursor's
// index) advances only when the slice length didn't change. When a Fire body called
// Destroy, the spliced *Aura is released to the aura pool only after Fire fully
// returns — the post-fire activeEngine clear must complete before the *Aura can be
// handed off, otherwise a concurrent goroutine that pulled the recycled pointer would
// race the write.
func (ge *GameEngine) fireMatching(triggeringCard card.Card, trigger triggertype.Type) {
	for i := 0; i < len(ge.auras); {
		a := ge.auras[i]
		if a.TriggerType() != trigger || (a.OncePerTurn() && a.FiredThisTurn()) {
			i++
			continue
		}
		ge.triggeringCard = triggeringCard
		ge.currentAuraIdx = i
		ge.currentAuraDestroyed = false
		a.Fire(ge, ge.logger)
		ge.currentAuraIdx = -1
		ge.triggeringCard = nil
		if ge.currentAuraDestroyed {
			if r, ok := a.(auraReleaser); ok {
				r.Release()
			}
			continue
		}
		ge.auras[i].SetFiredThisTurn(true)
		i++
	}
}

// auraReleaser is the optional opt-in interface concrete Aura impls implement when
// their per-instance memory can be returned to a pool after destruction. The engine
// queries for it via type assertion so the gameengine.Aura interface stays minimal.
type auraReleaser interface{ Release() }

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

// FireEndOfTurn runs after the chain has finished resolving (and the legality gates
// have passed) but before the carry state is captured. Walks Auras and Triggers in one
// pass each:
//
//   - Aura entries respect OncePerTurn / FiredThisTurn semantics; the handler owns
//     destruction via the engine's destroyAura path.
//   - Trigger entries are one-shot; fired entries are removed afterward. Snapshotting
//     len(ge.triggers) before iterating keeps a handler that calls AddXxxTrigger from
//     firing its newcomer on the same pass — newcomers stay queued for the next
//     matching event.
func (ge *GameEngine) FireEndOfTurn() {
	ge.fireMatching(nil, triggertype.EndOfTurn)
	n := len(ge.triggers)
	for i := 0; i < n; i++ {
		tr := ge.triggers[i]
		if tr.TriggerType() != triggertype.EndOfTurn {
			continue
		}
		tr.Fire(ge, ge.logger)
	}
	kept := ge.triggers[:0]
	for i, tr := range ge.triggers {
		if i < n && tr.TriggerType() == triggertype.EndOfTurn {
			continue
		}
		kept = append(kept, tr)
	}
	ge.triggers = kept
}

// FireHit walks the one-shot trigger queue and invokes every triggertype.Hit entry
// whose type filter matches the attacking card's types. Surviving entries (filter
// mismatch) are kept; fired entries are removed.
func (ge *GameEngine) FireHit(attackerTypes card.TypeSet) {
	kept := ge.triggers[:0]
	for i := range ge.triggers {
		t := ge.triggers[i]
		if t.TriggerType() != triggertype.Hit || !t.Matches(attackerTypes) {
			kept = append(kept, t)
			continue
		}
		t.Fire(ge, ge.logger)
	}
	ge.triggers = kept
}

// FireStartOfTurn walks ge.auras and invokes every triggertype.StartOfTurn entry,
// calling onFire with each entry's pre-state snapshot so sim can attribute damage /
// draws / log lines back to the firing aura. Auras that destroy themselves splice
// out; FiredThisTurn flips reset on each fresh turn boundary.
//
// The onFire callback receives:
//   - pre is the index in ge.auras of the firing entry at the time of the call.
//   - damage is ge.value's delta during this handler — the partition tiebreaker uses
//     it.
//   - drawnCard is the first card the handler appended to hand, or nil. Used by
//     processAurasAtStartOfTurn to surface "revealed" entries.
//   - newLogEntries is the slice of LogEntries the handler appended (caller may copy
//     out).
func (ge *GameEngine) FireStartOfTurn(onFire func(idx int, damage int, drawnCard card.Card, newLogEntries []turnlogger.LogEntry)) {
	for i := 0; i < len(ge.auras); {
		a := ge.auras[i]
		ge.auras[i].SetFiredThisTurn(false)
		if a.TriggerType() != triggertype.StartOfTurn {
			i++
			continue
		}
		preHand := len(ge.hand)
		preLog := 0
		if ge.logger != nil {
			preLog = len(ge.logger.Entries())
		}
		preValue := ge.value
		ge.currentAuraIdx = i
		ge.currentAuraDestroyed = false
		a.Fire(ge, ge.logger)
		ge.currentAuraIdx = -1

		damage := ge.value - preValue
		var drawn card.Card
		if len(ge.hand) > preHand {
			drawn = ge.hand[preHand]
		}
		var newEntries []turnlogger.LogEntry
		if ge.logger != nil {
			if entries := ge.logger.Entries(); len(entries) > preLog {
				newEntries = entries[preLog:]
			}
		}
		if onFire != nil {
			onFire(i, damage, drawn, newEntries)
		}
		if ge.currentAuraDestroyed {
			if r, ok := a.(auraReleaser); ok {
				r.Release()
			}
			continue
		}
		i++
	}
}

// DestroyAura removes the aura currently being fired and, when addToGraveyard==true, pushes
// the aura's source card into the graveyard (token auras with no source no-op). Direct
// splice (no cacheable flip) — destruction is deterministic from the triggering event.
//
// Called by the card.Aura context the engine threads into each handler; cards do not call
// this directly.
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

// ResolveChainStep runs card.Play on self and then applies the standard chain-step
// resolution: for an attack-action or weapon-attack, credit self.EffectiveAttack() to
// ge.value; for a defense-reaction (or DefensiveInstant), credit the EffectiveDefense
// capped at the remaining unblocked damage and bank it via AddDamageBlocked; for
// everything else, log (+0). The "<DisplayName>: <VERB> (+N)" chain-step entry is
// appended after Play returns so any self-buffs Play applied (e.g. modal +2{p} riders
// flipping self.BonusAttack) are reflected in the displayed delta.
//
// Cards' Play body owns card-specific behaviour: riders that emit rider log lines,
// OnHit registration, conditional self-buffs, sub-card plays. The standard
// printed-attack-deals-damage / DR-blocks-incoming mechanic is the engine's job;
// cards don't reach for DealEffectiveAttack / DealEffectiveDefense or emit the chain
// step themselves.
func (ge *GameEngine) ResolveChainStep(l card.Logger, self *card.CardState) {
	self.Card.Play(ge, l, self)
	if self.Card.Types(nil).Has(card.TypeAura) {
		ge.auraCreated = true
	}
	n := ge.chainStepDelta(self)
	l.AppendChainStep(ChainStepText(self), n)
}

// PlayCard implements card.GameEngine.PlayCard. Cards reach this when they resolve
// another card mid-handler (Moon Wish tutoring Sun Kiss into play on go-again).
func (ge *GameEngine) PlayCard(l card.Logger, self *card.CardState) {
	ge.ResolveChainStep(l, self)
}

// chainStepDelta computes the chain step's display delta and applies the standard
// damage / block side effects. Returns the (+N) value for the log line.
func (ge *GameEngine) chainStepDelta(self *card.CardState) int {
	types := self.Card.Types(nil)
	switch {
	case types.IsAttackAction() || types.IsWeaponAttack():
		n := self.EffectiveAttack()
		ge.value += n
		return n
	case types.IsDefenseReaction() || isDefensiveInstant(self.Card):
		n := self.EffectiveDefense()
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
// qualifying hits to a card-type predicate; nil = any hit qualifies. The handler
// signature matches card.GameEngine's inline declaration.
func (ge *GameEngine) AddHitTrigger(self *card.CardState, handler func(card.GameEngine, card.Logger, card.Trigger), filter func(card.TypeSet) bool) {
	ge.CreateTrigger(trigger.NewFromCard(self.Card, triggertype.Hit, handler, filter))
}

// AddEndOfTurnTrigger registers a one-shot triggertype.EndOfTurn listener — fires
// after the chain finishes resolving but before the carry-state snapshot.
func (ge *GameEngine) AddEndOfTurnTrigger(self *card.CardState, handler func(card.GameEngine, card.Logger, card.Trigger)) {
	ge.CreateTrigger(trigger.NewFromCard(self.Card, triggertype.EndOfTurn, handler, nil))
}

// === Tokens ===

// Card-facing token creation / count methods on *GameEngine. internal/card.GameEngine
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
func (ge *GameEngine) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	ge.AddValue(n)
	bumpOrCreateAura(ge.GameState, tokenNameRunechant, func(n int) Aura { return token.NewRunechant(n) }, n)
}

// CreatePonders creates n Ponder tokens. No Value credit — Ponder pays out at end of
// turn (see the runtime's Ponder aura handler).
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

// ConsumeItemByName decrements the matching item's Count by n and removes the entry
// when Count reaches zero. Token items don't head to the graveyard on destroy. No-op
// when no item matches name. Called by token-ability Play implementations registered
// outside gameengine.
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
