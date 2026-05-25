package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// GameState owns the raw per-turn data — every slice, scalar, and flag the engine reads or
// writes during a chain run. Internal machinery (per-permutation scratch, cross-turn
// snapshots, the find-best winning pointer) holds *GameState; GameEngine wraps it and adds
// the rules-engine API cards see. Unexported fields are touched only through methods.
//
// Fields split into cross-turn carryover (hero, weapons, hand, deck, arsenal, graveyard,
// banished, auras, items, opponentMarked, isMyTurn, incomingDamage, arcaneIncomingDamage)
// and the ephemeral struct of per-turn-resolution scratch. ephemeral.reset zeroes that
// group in one statement, so a new per-turn field can't leak across turn boundaries.
type GameState struct {
	hero Hero
	// heroTriggerType caches hero.TriggerType() so FireTriggers can gate the hero fire on a
	// field read instead of a per-event interface dispatch. SetHero keeps it in sync.
	heroTriggerType triggertype.Type
	weapons         []weapon.Weapon // currently-equipped weapons; persistent across turns

	hand      []card.CardState // role-tagged; each entry carries its partition role
	deck      *deck.Deck       // owned scratch; refilled on Reset via CopyFrom
	arsenal   card.Card
	graveyard []card.Card
	banished  []card.Card
	auras     []Aura
	items     []Item

	incomingDamage       int
	arcaneIncomingDamage int

	opponentMarked bool
	isMyTurn       bool

	ephemeral
}

// ephemeral groups per-turn-resolution scratch — every field a chain run accumulates or
// overwrites. reset zeroes the whole struct in one statement, with the four non-zero
// defaults (actionPoints, currentHookIdx, cacheable, logger) reseated in the same literal.
type ephemeral struct {
	cardsPlayed    []card.Card
	cardsRemaining []*card.CardState
	pitched        []card.Card
	defenders      []card.Card

	triggers             []EphemeralTrigger
	triggeringCard       card.Card
	attackReactionTarget *card.CardState

	logger card.Logger

	actionPoints   int
	value          int
	damageBlocked  int
	blockTotal     int
	currentHookIdx int
	// cardsRemovedFromDeck counts deck → non-deck movements during this chain (draws,
	// tutors, peek-and-banish, etc.). The hand-eval cache stores it and refuses to replay
	// against a shallower deck — the cached chain consumed N cards, so replay needs ≥ N.
	cardsRemovedFromDeck int

	cardBanished          bool
	arcaneDamageDealt     bool
	auraCreated           bool
	nonAttackActionPlayed bool
	lastAttackHit         bool
	currentHookDestroyed  bool
	currentStepRerouted   bool
	cacheable             bool
	heroTapped            bool
	hasCrowdCheered       bool
	hasCrowdBooed         bool
}

// reset returns e to its start-of-turn baseline: scalars zero except actionPoints=1,
// currentHookIdx=-1, cacheable=true, logger=NoopLogger; slice fields keep their backing
// but truncate to zero length so the per-perm chain runner's appends (notably
// AppendCardsPlayed) reuse the cap instead of allocating fresh each iteration. Aura /
// item FiredThisTurn flags live on the entries themselves and are rearmed by
// ResetEphemeralState's separate loop.
func (e *ephemeral) reset() {
	cardsPlayed := e.cardsPlayed[:0]
	cardsRemaining := e.cardsRemaining[:0]
	pitched := e.pitched[:0]
	defenders := e.defenders[:0]
	triggers := e.triggers[:0]
	*e = ephemeral{
		cardsPlayed:    cardsPlayed,
		cardsRemaining: cardsRemaining,
		pitched:        pitched,
		defenders:      defenders,
		triggers:       triggers,
		actionPoints:   1,
		currentHookIdx: -1,
		cacheable:      true,
		logger:         NoopLogger{},
	}
}

// Engine wraps s in a *GameEngine so the chain runner can drive Card.Play hooks against
// the rules-engine API while the underlying state remains the same pointer the caller
// holds. Cheap (single struct allocation); the engine doesn't copy state.
func (gs *GameState) Engine() *GameEngine { return &GameEngine{GameState: gs} }

// Copy returns a deep copy of gs. Slice and *deck.Deck fields get fresh backing storage;
// Aura / Item entries are deep-copied so per-permutation Count / FiredThisTurn mutations
// stay isolated. Triggers are effectively immutable, so only the slice header duplicates.
// Logger resets to nil — the caller installs a fresh per-clone logger when recording.
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
		out.triggers = append([]EphemeralTrigger(nil), gs.triggers...)
	} else {
		out.triggers = nil
	}
	out.items = copyItemsInto(nil, gs.items)
	out.logger = NoopLogger{}
	return &out
}

// CopyFrom rewrites *gs in place to match what src.Copy() would produce. Reuses the
// receiver's slice / deck backings when capacity permits so a pool slot stands in for
// repeated masterState.Copy() allocations. Auras / items deep-copy per entry into the
// pool, amortising both the *GameState alloc and the per-entry *Aura / *Item allocs.
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
			gs.triggers = make([]EphemeralTrigger, n)
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

// copyAurasInto returns a per-entry deep copy of src, reusing pool's backing when capacity
// permits and calling CopyInto on each prior-slot entry to rewrite *Aura allocations in
// place. Empty src returns nil (matches Copy()'s nil-on-empty); a nil pool falls through
// to a fresh Copy() per entry.
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
func resetCardSlice[T any](pooled, src []T) []T {
	if len(src) == 0 {
		return nil
	}
	if cap(pooled) >= len(src) {
		out := pooled[:len(src)]
		copy(out, src)
		return out
	}
	out := make([]T, len(src))
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
// follows with ResetEphemeralState which wipes ephemeral fields, so this only needs to
// touch the persistent carryover; the chain-step pool slot's ephemeral content gets
// thrown out either way.
func (gs *GameState) CopyPersistentStateFrom(src *GameState) {
	gs.hero = src.hero
	gs.heroTriggerType = src.heroTriggerType
	gs.weapons = src.weapons
	gs.arsenal = src.arsenal
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
	gs.auras = copyAurasInto(gs.auras, src.auras)
	gs.items = copyItemsInto(gs.items, src.items)
	gs.incomingDamage = src.incomingDamage
	gs.arcaneIncomingDamage = src.arcaneIncomingDamage
	gs.opponentMarked = src.opponentMarked
	gs.isMyTurn = src.isMyTurn
}

// ResetEphemeralState returns gs to its start-of-turn baseline: it discards every field
// that playing out a turn accumulates, keeping only the cross-turn carryover (hero, deck,
// hand, arsenal, graveyard, banished, the aura / item lists, opponentMarked, isMyTurn,
// and the matchup's incoming-damage figures).
//
// gs.ephemeral.reset wipes every per-turn scratch field in one struct assignment. The
// aura / item / hero FiredThisTurn flags live on those entries themselves, not in
// ephemeral, so they get their own rearm loop here so OncePerTurn auras can fire again.
func (gs *GameState) ResetEphemeralState() {
	gs.ephemeral.reset()
	for _, a := range gs.auras {
		a.SetFiredThisTurn(false)
	}
	for _, it := range gs.items {
		it.SetFiredThisTurn(false)
	}
	if gs.hero != nil && gs.hero.OncePerTurn() {
		gs.hero.SetFiredThisTurn(false)
	}
}

// Reset re-seeds gs with the given hero / weapons / incoming-damage values, preserving
// pre-allocated slice backings (hand, graveyard, banished, auras, items, ephemeral).
// Lets a pooled GameState start a fresh shuffle without losing prewarmed backings.
func (gs *GameState) Reset(h Hero, weapons []weapon.Weapon, incoming, arcaneIncoming int) {
	hand := gs.hand[:0]
	graveyard := gs.graveyard[:0]
	banished := gs.banished[:0]
	auras := gs.auras[:0]
	items := gs.items[:0]
	eph := gs.ephemeral
	*gs = GameState{
		hand:                 hand,
		graveyard:            graveyard,
		banished:             banished,
		auras:                auras,
		items:                items,
		weapons:              weapons,
		incomingDamage:       incoming,
		arcaneIncomingDamage: arcaneIncoming,
		ephemeral:            eph,
	}
	gs.ephemeral.reset()
	gs.SetHero(h)
}

// === Pure state accessors. No cacheable flips; sim uses these to drive the chain runner. ===

func (gs *GameState) Hero() Hero { return gs.hero }
func (gs *GameState) SetHero(h Hero) {
	gs.hero = h
	if h != nil {
		gs.heroTriggerType = h.TriggerType()
	} else {
		gs.heroTriggerType = 0
	}
}
func (gs *GameState) Weapons() []weapon.Weapon     { return gs.weapons }
func (gs *GameState) SetWeapons(w []weapon.Weapon) { gs.weapons = w }
func (gs *GameState) IsCacheable() bool            { return gs.cacheable }
func (gs *GameState) SetCacheable(v bool)          { gs.cacheable = v }

// CardsRemovedFromDeck reports how many cards have been moved out of the deck during
// this chain resolution (mid-chain draws, tutors, peek-and-banish, etc.).
func (gs *GameState) CardsRemovedFromDeck() int { return gs.cardsRemovedFromDeck }

// noteDeckRemoval increments the deck-removal counter; called by every helper that
// pops a card off the deck. Package-private — only the engine knows when a removal
// has happened.
func (gs *GameState) noteDeckRemoval(n int) { gs.cardsRemovedFromDeck += n }

// Hand projects the role-tagged hand down to the bare cards. Callers needing the roles
// (Discard, the chain runner) read HandStates.
func (gs *GameState) Hand() []card.Card {
	if len(gs.hand) == 0 {
		return nil
	}
	out := make([]card.Card, len(gs.hand))
	for i := range gs.hand {
		out[i] = gs.hand[i].Card
	}
	return out
}

// HandStates returns the live role-tagged hand.
func (gs *GameState) HandStates() []card.CardState { return gs.hand }

// HandSize reports the total number of cards in the hand zone, including entries added
// by mid-chain DrawOne. The value is determined by partition + chain progress alone, so
// this accessor doesn't flip IsCacheable.
func (gs *GameState) HandSize() int { return len(gs.hand) }

// HandHasMatching reports whether any non-drawn hand entry satisfies pred. FromDraw
// entries are skipped: their identity is unknown to in-chain attribute reads. Doesn't
// flip IsCacheable — the starting-hand multiset is already part of the cache key.
func (gs *GameState) HandHasMatching(pred func(card.Card) bool) bool {
	for i := range gs.hand {
		if gs.hand[i].FromDraw {
			continue
		}
		if pred(gs.hand[i].Card) {
			return true
		}
	}
	return false
}

// HeldHandSize reports the total Held-role hand-entry count, including drawn entries.
// Counting alone doesn't reveal drawn-card attributes, so this accessor stays
// cacheable; an actual pop that lands on a drawn entry flips IsCacheable in PopHandAt
// because the popped card's identity then leaks out.
func (gs *GameState) HeldHandSize() int {
	n := 0
	for i := range gs.hand {
		if gs.hand[i].Role == card.Held {
			n++
		}
	}
	return n
}

// SetHand installs h as the hand with every card defaulting to the Held role. Role-aware
// callers (the defense pass, the chain runner) use SetHandStates.
func (gs *GameState) SetHand(h []card.Card) {
	gs.hand = gs.hand[:0]
	for _, c := range h {
		gs.hand = append(gs.hand, card.CardState{Card: c, Role: card.Held})
	}
}

// SetHandStates installs an already role-tagged hand.
func (gs *GameState) SetHandStates(h []card.CardState) { gs.hand = h }

func (gs *GameState) AppendHandRaw(c card.Card) {
	gs.hand = append(gs.hand, card.CardState{Card: c, Role: card.Held})
}

// RemoveFromHand removes the first matching card from the hand without flipping
// IsCacheable. Returns true if a card was removed. Does not preserve order — the
// removed slot is filled by the last element (swap-with-last). The chain runner
// reads hand by membership / length, not by index, so order doesn't matter.
func (gs *GameState) RemoveFromHand(c card.Card) bool {
	for i := range gs.hand {
		if gs.hand[i].Card == c {
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

func (gs *GameState) Auras() []Aura                { return gs.auras }
func (gs *GameState) SetAuras(a []Aura)            { gs.auras = a }
func (gs *GameState) ClearAuras()                  { gs.auras = nil }
func (gs *GameState) Triggers() []EphemeralTrigger { return gs.triggers }
func (gs *GameState) ClearTriggers()               { gs.triggers = nil }
func (gs *GameState) Items() []Item                { return gs.items }
func (gs *GameState) ClearItems()                  { gs.items = nil }

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

// HasCrowdCheered / HasCrowdBooed report whether the crowd has cheered / booed your hero
// at any point during the current turn. Used by "if you've been cheered this turn" gates;
// "whenever the crowd cheers you" handlers subscribe to triggertype.CrowdCheer / CrowdBoo
// instead. Both flip together with the trigger fire via GameState.CrowdCheer / CrowdBoo
// methods on the engine wrapper.
func (gs *GameState) HasCrowdCheered() bool { return gs.hasCrowdCheered }
func (gs *GameState) HasCrowdBooed() bool   { return gs.hasCrowdBooed }

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

// HeroTapped reports whether the hero is tapped.
func (gs *GameState) HeroTapped() bool { return gs.heroTapped }

// UntapHero untaps the hero — the printed "untap your hero" effect.
func (gs *GameState) UntapHero() { gs.heroTapped = false }

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

// HeroHasType reports whether the active hero's type line contains t. Returns false when
// no hero is set. Used by cards that gate on hero-only type keywords (TypeRevered /
// TypeReviled, etc.).
func (gs *GameState) HeroHasType(t card.CardType) bool {
	if gs.hero == nil {
		return false
	}
	return gs.hero.Types().Has(t)
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
func appendCopy[T any](dst []T, src []T) []T {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}
