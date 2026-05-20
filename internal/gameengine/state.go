package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// GameState owns the raw per-turn data — every slice, every scalar, every flag the
// engine reads or writes during a chain run. Internal machinery (sim's per-permutation
// scratch, snapshots that carry into the next turn, the find-best winning pointer) holds
// *GameState because it cares about the data, not the rules. GameEngine wraps a
// *GameState and adds the rules-engine API cards see; the unexported fields stay
// package-private so external callers can only touch them through methods.
type GameState struct {
	hero    Hero
	weapons []weapon.Weapon // currently-equipped weapons; persistent across turns

	hand      []card.Card
	deck      *deck.Deck // owned scratch; refilled on Reset via CopyFrom
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

	logger               card.Logger
	triggeringCard       card.Card
	attackReactionTarget *card.CardState

	actionPoints         int
	value                int
	incomingDamage       int
	damageBlocked        int
	arcaneIncomingDamage int
	blockTotal           int
	currentAuraIdx       int

	cardBanished          bool
	arcaneDamageDealt     bool
	opponentMarked        bool
	auraCreated           bool
	nonAttackActionPlayed bool
	lastAttackHit         bool
	currentAuraDestroyed  bool
	currentStepRerouted   bool
	cacheable             bool
	isMyTurn              bool
}

// Engine wraps s in a *GameEngine so the chain runner can drive Card.Play hooks against
// the rules-engine API while the underlying state remains the same pointer the caller
// holds. Cheap (single struct allocation); the engine doesn't copy state.
func (gs *GameState) Engine() *GameEngine { return &GameEngine{GameState: gs} }

// Copy returns a deep copy of gs. Slice and *deck.Deck fields get fresh backing storage;
// Aura / Item entries are deep-copied via their Copy() methods so per-permutation
// Count / FiredThisTurn mutations stay isolated. Triggers are effectively immutable after
// construction, so only the slice header is duplicated. Logger is reset to nil — the
// caller installs a fresh per-clone logger when recording, leaving find-best copies
// allocation-free.
func (gs *GameState) Copy() *GameState {
	out := *gs
	out.hand = appendCopy(nil, gs.hand)
	if gs.deck != nil {
		out.deck = gs.deck.Copy()
	}
	out.graveyard = appendCopy(nil, gs.graveyard)
	out.banished = appendCopy(nil, gs.banished)
	out.pitched = appendCopy(nil, gs.pitched)
	out.defenders = appendCopy(nil, gs.defenders)
	out.cardsPlayed = appendCopy(nil, gs.cardsPlayed)
	if len(gs.cardsRemaining) > 0 {
		out.cardsRemaining = append([]*card.CardState(nil), gs.cardsRemaining...)
	} else {
		out.cardsRemaining = nil
	}
	out.auras = copyAurasInto(nil, gs.auras)
	if len(gs.triggers) > 0 {
		out.triggers = append([]Trigger(nil), gs.triggers...)
	} else {
		out.triggers = nil
	}
	out.items = copyItemsInto(nil, gs.items)
	out.logger = NoopLogger{}
	return &out
}

// CopyFrom rewrites *gs in place to match what src.Copy() would produce. Reuses
// the receiver's slice and *deck.Deck backings when capacity permits so a pool slot can
// stand in for repeated per-Best masterState.Copy() allocations. Auras and items are
// deep-copied per entry via CopyInto on the pooled slot when present, so the pool
// amortises both the *GameState alloc and the per-entry *Aura / *Item allocs.
func (gs *GameState) CopyFrom(src *GameState) {
	pooledHand := gs.hand
	pooledGrav := gs.graveyard
	pooledBanished := gs.banished
	pooledPitched := gs.pitched
	pooledDefenders := gs.defenders
	pooledCardsPlayed := gs.cardsPlayed
	pooledCardsRemaining := gs.cardsRemaining
	pooledAuras := gs.auras
	pooledTriggers := gs.triggers
	pooledItems := gs.items
	pooledDeck := gs.deck
	*gs = *src
	gs.logger = NoopLogger{}
	gs.hand = resetCardSlice(pooledHand, src.hand)
	gs.graveyard = resetCardSlice(pooledGrav, src.graveyard)
	gs.banished = resetCardSlice(pooledBanished, src.banished)
	gs.pitched = resetCardSlice(pooledPitched, src.pitched)
	gs.defenders = resetCardSlice(pooledDefenders, src.defenders)
	gs.cardsPlayed = resetCardSlice(pooledCardsPlayed, src.cardsPlayed)
	if n := len(src.cardsRemaining); n > 0 {
		if cap(pooledCardsRemaining) >= n {
			gs.cardsRemaining = pooledCardsRemaining[:n]
		} else {
			gs.cardsRemaining = make([]*card.CardState, n)
		}
		copy(gs.cardsRemaining, src.cardsRemaining)
	} else {
		gs.cardsRemaining = nil
	}
	if n := len(src.triggers); n > 0 {
		if cap(pooledTriggers) >= n {
			gs.triggers = pooledTriggers[:n]
		} else {
			gs.triggers = make([]Trigger, n)
		}
		copy(gs.triggers, src.triggers)
	} else {
		gs.triggers = nil
	}
	gs.auras = copyAurasInto(pooledAuras, src.auras)
	gs.items = copyItemsInto(pooledItems, src.items)
	if src.deck != nil {
		if pooledDeck != nil {
			pooledDeck.CopyFrom(src.deck)
			gs.deck = pooledDeck
		} else {
			gs.deck = src.deck.Copy()
		}
	} else {
		gs.deck = nil
	}
}

// copyAurasInto returns a per-entry deep copy of src, reusing pool's backing slice when
// capacity permits and reaching for CopyInto on each prior-slot entry so concrete *Aura
// allocations are rewritten in place rather than replaced. Empty src returns nil to match
// the Copy() path's nil-on-empty semantics; a nil pool falls through to a fresh Copy()
// per entry (CopyInto's contract on a nil dst).
func copyAurasInto(pool, src []Aura) []Aura {
	n := len(src)
	if n == 0 {
		return nil
	}
	var out []Aura
	if cap(pool) >= n {
		out = pool[:n]
	} else {
		out = make([]Aura, n)
	}
	priorLen := len(pool)
	for i, a := range src {
		var prev any
		if i < priorLen {
			prev = pool[i]
		}
		out[i] = a.CopyInto(prev).(Aura)
	}
	return out
}

// copyItemsInto is the items counterpart of copyAurasInto. Item doesn't expose CopyInto
// today (no in-place reset surface), so the per-entry path always allocates via Copy().
// The pool's slice backing is still reused when capacity permits.
func copyItemsInto(pool, src []Item) []Item {
	n := len(src)
	if n == 0 {
		return nil
	}
	var out []Item
	if cap(pool) >= n {
		out = pool[:n]
	} else {
		out = make([]Item, n)
	}
	for i, it := range src {
		out[i] = it.Copy().(Item)
	}
	return out
}

// resetCardSlice returns a fresh slice header that aliases pooled when capacity permits,
// or a freshly allocated backing sized to src. Empty src returns nil to match the Copy()
// path's nil-on-empty semantics.
func resetCardSlice(pooled, src []card.Card) []card.Card {
	if len(src) == 0 {
		return nil
	}
	if cap(pooled) >= len(src) {
		out := pooled[:len(src)]
		copy(out, src)
		return out
	}
	out := make([]card.Card, len(src))
	copy(out, src)
	return out
}

// CopyPersistentState is a lighter variant of Copy that copies only the cross-turn
// persistent state — the inverse of ResetEphemeralState's reset set. Hand, pitched,
// defenders, cardsPlayed, cardsRemaining, triggers, deck, and logger are left nil
// (callers that want those populated set them after copying). Graveyard and banished
// get fresh backings so that source-side splice-style mutations (BanishFromGraveyard
// rewrites the source's backing in place during the same turn this snapshot is
// captured) don't bleed into the snapshot's view. Weapons are carried via the implicit
// `out := *gs` slice-header copy — no per-entry cloning (stateless structs). Auras
// and items get full per-entry copies because their fire-this-turn / count fields
// mutate independently of the source.
func (gs *GameState) CopyPersistentState() *GameState {
	out := *gs
	out.hand = nil
	out.pitched = nil
	out.defenders = nil
	out.cardsPlayed = nil
	out.cardsRemaining = nil
	out.triggers = nil
	out.deck = nil
	out.logger = NoopLogger{}
	if n := len(gs.graveyard); n > 0 {
		out.graveyard = append([]card.Card(nil), gs.graveyard...)
	}
	if n := len(gs.banished); n > 0 {
		out.banished = append([]card.Card(nil), gs.banished...)
	}
	out.auras = copyAurasInto(nil, gs.auras)
	out.items = copyItemsInto(nil, gs.items)
	return &out
}

// CopyPersistentStateFrom overwrites *gs in place to match what CopyPersistentState(src)
// would produce. Reuses gs's auras / items slice backing when capacity permits, avoiding
// the per-permutation slice allocation in the chain runner's hot loop. The chain runner
// follows with ResetEphemeralState which wipes the same ephemeral fields CopyPersistentState
// left nil, so this path mirrors CopyPersistentState exactly — only the persistent
// carryover fields and the aura / item per-entry copies are load-bearing here.
func (gs *GameState) CopyPersistentStateFrom(src *GameState) {
	// Stash the pool's existing slice backings before *gs = *src aliases them onto src.
	// copyAurasInto / copyItemsInto then write per-entry copies into the pooled backings
	// without trampling the source the leafState pool keeps.
	pooledAuras := gs.auras
	pooledItems := gs.items
	*gs = *src
	gs.hand = nil
	gs.pitched = nil
	gs.defenders = nil
	gs.cardsPlayed = nil
	gs.cardsRemaining = nil
	gs.triggers = nil
	gs.deck = nil
	gs.logger = NoopLogger{}
	if n := len(src.graveyard); n > 0 {
		gs.graveyard = src.graveyard[:n:n]
	} else {
		gs.graveyard = nil
	}
	if n := len(src.banished); n > 0 {
		gs.banished = src.banished[:n:n]
	} else {
		gs.banished = nil
	}
	gs.auras = copyAurasInto(pooledAuras, src.auras)
	gs.items = copyItemsInto(pooledItems, src.items)
}

// ResetEphemeralState returns gs to its start-of-turn baseline: it discards every field
// that playing out a turn accumulates, keeping only the cross-turn carryover (hero, deck,
// hand, arsenal, graveyard, banished, the aura / item lists, opponentMarked, and the
// matchup's incoming-damage figures).
//
// What it resets, by category:
//   - per-turn zones — pitched, defenders (hand persists; it was refilled at the end of
//     the previous turn and is the dealt hand visible to start-of-turn aura handlers)
//   - resolution scratch — the played / remaining lists, the one-shot trigger queue, the
//     value accumulator and draw counter, action points, the block total, the current-step
//     machinery, the "happened this resolution" flags, the cacheable bit, the logger
//   - aura gates — every aura's FiredThisTurn flag rearms, so OncePerTurn auras can fire
//     again
//
// auraCreated resets to false — it means "an aura was played or created THIS turn", so
// auras carried over from a previous turn must not satisfy it.
//
// incomingDamage stays put — it's the constant matchup figure, carried over untouched.
// damageBlocked (how much of it defense has absorbed so far) resets to zero.
func (gs *GameState) ResetEphemeralState() {
	gs.pitched = nil
	gs.defenders = nil
	gs.cardsPlayed = nil
	gs.cardsRemaining = nil
	gs.triggers = nil
	gs.triggeringCard = nil
	gs.attackReactionTarget = nil
	gs.actionPoints = 1
	gs.value = 0
	gs.damageBlocked = 0
	gs.blockTotal = 0
	gs.currentAuraDestroyed = false
	gs.currentStepRerouted = false
	gs.currentAuraIdx = -1
	gs.cardBanished = false
	gs.arcaneDamageDealt = false
	gs.nonAttackActionPlayed = false
	gs.lastAttackHit = false
	gs.cacheable = true
	gs.logger = NoopLogger{}
	gs.auraCreated = false
	for _, a := range gs.auras {
		a.SetFiredThisTurn(false)
	}
}

// === Pure state accessors. No cacheable flips; sim uses these to drive the chain runner. ===

func (gs *GameState) Hero() Hero                   { return gs.hero }
func (gs *GameState) SetHero(h Hero)               { gs.hero = h }
func (gs *GameState) Weapons() []weapon.Weapon     { return gs.weapons }
func (gs *GameState) SetWeapons(w []weapon.Weapon) { gs.weapons = w }
func (gs *GameState) IsCacheable() bool            { return gs.cacheable }
func (gs *GameState) SetCacheable(v bool)          { gs.cacheable = v }

func (gs *GameState) Hand() []card.Card     { return gs.hand }
func (gs *GameState) SetHand(h []card.Card) { gs.hand = h }
func (gs *GameState) AppendHandRaw(c card.Card) {
	gs.hand = append(gs.hand, c)
}

// RemoveFromHand removes the first matching card from the hand without flipping
// IsCacheable. Returns true if a card was removed. Does not preserve order — the
// removed slot is filled by the last element (swap-with-last). The chain runner
// reads hand by membership / length, not by index, so order doesn't matter.
func (gs *GameState) RemoveFromHand(c card.Card) bool {
	for i := range gs.hand {
		if gs.hand[i] == c {
			last := len(gs.hand) - 1
			gs.hand[i] = gs.hand[last]
			gs.hand = gs.hand[:last]
			return true
		}
	}
	return false
}

func (gs *GameState) Deck() *deck.Deck     { return gs.deck }
func (gs *GameState) SetDeck(d *deck.Deck) { gs.deck = d }

func (gs *GameState) Graveyard() []card.Card      { return gs.graveyard }
func (gs *GameState) SetGraveyard(gv []card.Card) { gs.graveyard = gv }
func (gs *GameState) AppendGraveyard(c card.Card) { gs.graveyard = append(gs.graveyard, c) }

func (gs *GameState) Arsenal() card.Card     { return gs.arsenal }
func (gs *GameState) SetArsenal(c card.Card) { gs.arsenal = c }

func (gs *GameState) Banished() []card.Card     { return gs.banished }
func (gs *GameState) SetBanished(b []card.Card) { gs.banished = b }

func (gs *GameState) Auras() []Aura       { return gs.auras }
func (gs *GameState) ClearAuras()         { gs.auras = nil }
func (gs *GameState) Triggers() []Trigger { return gs.triggers }
func (gs *GameState) ClearTriggers()      { gs.triggers = nil }
func (gs *GameState) Items() []Item       { return gs.items }
func (gs *GameState) ClearItems()         { gs.items = nil }

// CreateAura appends a to the aura list. Flips AuraCreated so same-turn "if you've
// played or created an aura" riders see the entry.
func (gs *GameState) CreateAura(a Aura) {
	gs.auras = append(gs.auras, a)
	gs.auraCreated = true
}

// CreateTrigger appends t to the one-shot trigger queue.
func (gs *GameState) CreateTrigger(t Trigger) { gs.triggers = append(gs.triggers, t) }

// CreateItem appends i to the item list.
func (gs *GameState) CreateItem(i Item) { gs.items = append(gs.items, i) }

// AuraCount returns the count of live auras. Used by gates like Yinti Yanti's "while you
// control an aura" rider.
func (gs *GameState) AuraCount() int { return len(gs.auras) }

// RunechantCount / PonderCount / GoldCount / SilverCount / CopperCount return the live
// token-aura / token-item count by display name. Both GameState and GameEngine expose
// these so end-of-turn callers (TurnSummary readers) can read counts off the state
// pointer directly without needing an engine wrapper.
func (gs *GameState) RunechantCount() int { return auraCountByName(gs.auras, tokenNameRunechant) }
func (gs *GameState) PonderCount() int    { return auraCountByName(gs.auras, tokenNamePonder) }
func (gs *GameState) GoldCount() int      { return itemCountByName(gs.items, tokenNameGold) }
func (gs *GameState) SilverCount() int    { return itemCountByName(gs.items, tokenNameSilver) }
func (gs *GameState) CopperCount() int    { return itemCountByName(gs.items, tokenNameCopper) }

func (gs *GameState) Pitched() []card.Card     { return gs.pitched }
func (gs *GameState) SetPitched(p []card.Card) { gs.pitched = p }

func (gs *GameState) Defenders() []card.Card     { return gs.defenders }
func (gs *GameState) SetDefenders(d []card.Card) { gs.defenders = d }

func (gs *GameState) CardsPlayed() []card.Card               { return gs.cardsPlayed }
func (gs *GameState) SetCardsPlayed(cs []card.Card)          { gs.cardsPlayed = cs }
func (gs *GameState) AppendCardsPlayed(c card.Card)          { gs.cardsPlayed = append(gs.cardsPlayed, c) }
func (gs *GameState) CardsRemaining() []*card.CardState      { return gs.cardsRemaining }
func (gs *GameState) SetCardsRemaining(cs []*card.CardState) { gs.cardsRemaining = cs }

func (gs *GameState) Logger() card.Logger     { return gs.logger }
func (gs *GameState) SetLogger(l card.Logger) { gs.logger = l }

func (gs *GameState) TriggeringCard() card.Card                  { return gs.triggeringCard }
func (gs *GameState) AttackReactionTarget() *card.CardState      { return gs.attackReactionTarget }
func (gs *GameState) SetAttackReactionTarget(cs *card.CardState) { gs.attackReactionTarget = cs }

func (gs *GameState) ActionPoints() int     { return gs.actionPoints }
func (gs *GameState) SetActionPoints(n int) { gs.actionPoints = n }
func (gs *GameState) AddActionPoints(n int) { gs.actionPoints += n }

func (gs *GameState) Value() int     { return gs.value }
func (gs *GameState) SetValue(v int) { gs.value = v }
func (gs *GameState) AddValue(n int) { gs.value += n }

// RemainingUnblockedDamage returns the opponent damage still unblocked this turn — the
// constant matchup figure minus everything defense has absorbed so far.
func (gs *GameState) RemainingUnblockedDamage() int { return gs.incomingDamage - gs.damageBlocked }

// SetIncomingDamage installs the turn's incoming-damage figure and zeroes the
// damage-blocked accumulator — "n incoming, none blocked yet". Defense reactions and
// blocks then chip away at it via AddDamageBlocked (and the engine's DR resolution)
// rather than mutating the figure itself, so the matchup number stays constant and
// carries across turns untouched.
func (gs *GameState) SetIncomingDamage(n int) {
	gs.incomingDamage = n
	gs.damageBlocked = 0
}

// AddDamageBlocked credits n damage as absorbed by defense, shrinking
// RemainingUnblockedDamage by n. The engine's DR resolution accumulates through here; the
// chain runner's plain-block pass calls it directly.
func (gs *GameState) AddDamageBlocked(n int) { gs.damageBlocked += n }

func (gs *GameState) IncomingDamage() int           { return gs.incomingDamage }
func (gs *GameState) ArcaneIncomingDamage() int     { return gs.arcaneIncomingDamage }
func (gs *GameState) SetArcaneIncomingDamage(n int) { gs.arcaneIncomingDamage = n }

func (gs *GameState) BlockTotal() int     { return gs.blockTotal }
func (gs *GameState) SetBlockTotal(n int) { gs.blockTotal = n }

func (gs *GameState) ArcaneDamageDealt() bool     { return gs.arcaneDamageDealt }
func (gs *GameState) SetArcaneDamageDealt(v bool) { gs.arcaneDamageDealt = v }

func (gs *GameState) OpponentMarked() bool     { return gs.opponentMarked }
func (gs *GameState) SetOpponentMarked(v bool) { gs.opponentMarked = v }
func (gs *GameState) MarkOpponent()            { gs.opponentMarked = true }
func (gs *GameState) ClearOpponentMarked()     { gs.opponentMarked = false }

func (gs *GameState) AuraCreated() bool     { return gs.auraCreated }
func (gs *GameState) SetAuraCreated(v bool) { gs.auraCreated = v }

func (gs *GameState) CardBanished() bool     { return gs.cardBanished }
func (gs *GameState) SetCardBanished(v bool) { gs.cardBanished = v }

func (gs *GameState) NonAttackActionPlayed() bool     { return gs.nonAttackActionPlayed }
func (gs *GameState) SetNonAttackActionPlayed(v bool) { gs.nonAttackActionPlayed = v }

func (gs *GameState) LastAttackHit() bool     { return gs.lastAttackHit }
func (gs *GameState) SetLastAttackHit(v bool) { gs.lastAttackHit = v }

// IsMyTurn reports whether the active phase is the owning player's action phase (true) or
// the defense phase (false). The chain runner sets it; cards read it for "during your
// turn" riders.
func (gs *GameState) IsMyTurn() bool     { return gs.isMyTurn }
func (gs *GameState) SetIsMyTurn(v bool) { gs.isMyTurn = v }

func (gs *GameState) CurrentStepRerouted() bool     { return gs.currentStepRerouted }
func (gs *GameState) SetCurrentStepRerouted(v bool) { gs.currentStepRerouted = v }

// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. No-op when
// the logger is nil or when no chain-step entry exists yet.
func (gs *GameState) AmendLastChainStepN(n int) { gs.logger.AmendLastChainStepN(n) }

// HeroWantsLowerHealth reports whether the current hero opts into the LowerHealthWanter
// marker. Returns false when no hero is set.
func (gs *GameState) HeroWantsLowerHealth() bool {
	if gs.hero == nil {
		return false
	}
	_, ok := gs.hero.(LowerHealthWanter)
	return ok
}

// CurrentHeroClass returns the active hero's primary class. Zero when no hero is set.
func (gs *GameState) CurrentHeroClass() card.CardType {
	if gs.hero == nil {
		return 0
	}
	return gs.hero.Class()
}

// HasPlayedType reports whether any card played this turn has the given type. Universal
// cards' Types() folds the active hero's class through the engine, but at the GameState
// level we just check the recorded types directly.
func (gs *GameState) HasPlayedType(t card.CardType) bool {
	for _, c := range gs.cardsPlayed {
		if c.Types(nil).Has(t) {
			return true
		}
	}
	return false
}

// appendCopy is the small slice-clone helper Copy uses. Inlined into its own helper so
// each persistent slice's branch is a one-liner.
func appendCopy(dst []card.Card, src []card.Card) []card.Card {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}
