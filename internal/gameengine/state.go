package gameengine

import (
	"unsafe"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// auraPoolSize caps the number of simultaneously-live card-backed auras a single
// GameState supports. Real Viserai / Iyslander / Briar decks run 2-5 live card auras
// at a peak; 20 leaves 4x headroom. CreateAura panics on overflow so a regression that
// inflates aura count fails loud instead of silently slow-pathing onto the heap.
const auraPoolSize = 20

// GameState owns the raw per-turn data — every slice, scalar, and flag the engine reads or
// writes during an attack turn run. Internal machinery (per-permutation scratch, cross-turn
// snapshots, the find-best winning pointer) holds *GameState; GameEngine wraps it and adds
// the rules-engine API cards see. Unexported fields are touched only through methods.
//
// Fields split into cross-turn carryover (hero, weapons, hand, deck, arsenal, graveyard,
// banished, auras, items, opponentMarked, isMyTurn, incomingPhysicalDamage, incomingArcaneDamage)
// and the ephemeral struct of per-turn-resolution scratch. ephemeral.reset zeroes that
// group in one statement, so a new per-turn field can't leak across turn boundaries.
type GameState struct {
	hero Hero
	// heroTriggerType caches hero.TriggerType() so FireTriggers can gate the hero fire on a
	// field read instead of a per-event interface dispatch. SetHero keeps it in sync.
	heroTriggerType triggertype.Type
	weapons         []Weapon // currently-equipped weapon objects; persistent across turns

	hand      []card.CardState // role-tagged; each entry carries its partition role
	deck      *deck.Deck       // owned scratch; refilled on Reset via CopyFrom
	arsenal   card.Card
	graveyard []card.Card
	banished  []card.Card
	// auras holds interface boxes pointing into auraPool. Splicing on Destroy keeps the
	// slice live-only; the corresponding pool slot has its fields zeroed (SourceCard()
	// == nil) so CreateAura's next-free-slot scan picks it up for reuse.
	auras []Aura
	// auraPool is the per-GameState pre-allocated backing for card-backed auras.
	// CreateAura overwrites the next free slot in place; per-creation allocations are
	// gone. Pool exhaustion (>auraPoolSize simultaneously live) panics.
	auraPool [auraPoolSize]aura.Aura
	items    []Item // card-backed items only — tokens live in tokenItems

	// Token slots — one pre-allocated Aura / Item per kind, count=0 means "not in
	// play." Pointer identity stable for the GameState's lifetime; concrete pointer
	// types (not the Aura / Item interfaces) so Count / SetCount / Fire stay direct
	// calls instead of paying interface dispatch on every FireTriggers walk.
	tokenAuras [numTokenAuraKinds]*aura.Aura
	tokenItems [numTokenItemKinds]*item.Item

	// tokenAurasLiveBits / tokenItemsLiveBits mirror the per-slot Count() > 0
	// predicate as a bitmask: bit i set ⇔ tokenAuras[i].Count() > 0. Maintained by
	// bumpTokenAura / bumpTokenItem on 0→positive transitions and by DestroyAura /
	// DestroyItem (and Reset / Copy paths) on the reverse. Lets FireTriggers'
	// early-exit and per-pass snapshot reduce to a single byte load instead of an
	// N-slot interface walk.
	tokenAurasLiveBits uint8
	tokenItemsLiveBits uint8

	incomingPhysicalDamage int
	incomingArcaneDamage   int

	opponentMarked bool
	isMyTurn       bool

	ephemeral
}

// ephemeral groups per-turn-resolution scratch — every field an attack turn run accumulates or
// overwrites. reset zeroes the whole struct in one statement, with the four non-zero
// defaults (actionPoints, currentHookIdx, cacheable, logger) reseated in the same literal.
type ephemeral struct {
	cardsPlayed    []card.Card
	cardsRemaining []*card.CardState
	pitched        []*card.CardState
	defenders      []card.Card

	triggers             []EphemeralTrigger
	attackReactionTarget *card.CardState

	logger card.Logger

	actionPoints int
	value        int
	// damageDealt is the cumulative damage credited at every triggertype.Hit fire this
	// turn — gates "if you've dealt {N} this turn" riders. Resets at the turn boundary
	// via ephemeral.reset.
	damageDealt int
	// physicalDamageBlocked / arcaneDamageBlocked accumulate this turn's mitigation of each
	// incoming damage type — blocks and PreventPhysicalDamage on the physical side,
	// PreventArcaneDamage on the arcane side. Kept separate from the constant matchup
	// figures (incomingPhysicalDamage / incomingArcaneDamage) so RemainingPhysical/ArcaneDamage =
	// figure - blocked, and both reset per turn via ephemeral.reset rather than mutating
	// the carryover figure (which would leak prevention into every later turn).
	physicalDamageBlocked int
	arcaneDamageBlocked   int
	blockTotal            int
	currentHookIdx        int
	// currentFiringTokenAura / currentFiringTokenItem track which token slot is mid-Fire
	// so DestroyAura / DestroyItem route to "zero the count" instead of splicing s.auras
	// / s.items. -1 means no token is firing (the splice path is active).
	currentFiringTokenAura int
	currentFiringTokenItem int
	// cardsRemovedFromDeck counts deck → non-deck movements during this attack turn (draws,
	// tutors, peek-and-banish, etc.). The hand-eval cache stores it and refuses to replay
	// against a shallower deck — the cached attack turn consumed N cards, so replay needs ≥ N.
	cardsRemovedFromDeck int

	cardBanished          bool
	arcaneDamageDealt     bool
	auraCreated           bool
	nonAttackActionPlayed bool
	lastAttackHit         bool
	// destroyedOpponentArsenal flips true once the opposing arsenal is destroyed this turn;
	// the once-per-turn gate lives inside GameEngine.DestroyOpponentArsenal.
	destroyedOpponentArsenal bool
	// hitThisTurn flips true the first time a physical attack lands (the "if hit"
	// branch in the sequence runner). Distinct from damageDealt > 0 because arcane
	// damage contributes to damageDealt but does not count as a "hit" — same reason
	// arcane doesn't strip Mark.
	hitThisTurn          bool
	currentHookDestroyed bool
	currentStepRerouted  bool
	cacheable            bool
	heroTapped           bool
	hasCrowdCheered      bool
	hasCrowdBooed        bool
	// lightningFused records whether a Lightning card has been fused this turn — the
	// prerequisite Amulet of Lightning's instant ability gates on. Per-turn (no card sets it
	// across turns), so it lives here and clears on reset.
	lightningFused bool
}

// reset returns e to its start-of-turn baseline: scalars zero except actionPoints=1,
// currentHookIdx=-1, cacheable=true, logger=NoopLogger; slice fields keep their backing
// but truncate to zero length so the per-perm attack-turn runner's appends (notably
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
		cardsPlayed:            cardsPlayed,
		cardsRemaining:         cardsRemaining,
		pitched:                pitched,
		defenders:              defenders,
		triggers:               triggers,
		actionPoints:           1,
		currentHookIdx:         -1,
		currentFiringTokenAura: -1,
		currentFiringTokenItem: -1,
		cacheable:              true,
		logger:                 NoopLogger{},
	}
}

// Engine wraps s in a *GameEngine so the attack-turn runner can drive Card.Play hooks against
// the rules-engine API while the underlying state remains the same pointer the caller
// holds. Cheap (single struct allocation); the engine doesn't copy state.
func (gs *GameState) Engine() *GameEngine { return &GameEngine{GameState: gs} }

// Copy returns a deep copy of gs. Slice and *deck.Deck fields get fresh backing storage;
// Aura / Item entries are deep-copied so per-permutation Count / FiredThisTurn mutations
// stay isolated. Triggers are effectively immutable, so only the slice header duplicates.
// Logger resets to nil — the caller installs a fresh per-clone logger when recording.
func (gs *GameState) Copy() *GameState {
	out := *gs
	out.weapons = copyWeaponsInto(nil, gs.weapons)
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
	// out := *gs memcopied gs.auraPool byte-identically. Rebuild out.auras so its
	// interface values point into out.auraPool (not gs.auraPool) at the same pool
	// indices src.auras occupied.
	out.auras = rebuildPooledAuras(&out, gs.auras, &gs.auraPool[0], nil)
	if len(gs.triggers) > 0 {
		out.triggers = append([]EphemeralTrigger(nil), gs.triggers...)
	} else {
		out.triggers = nil
	}
	out.items = copyItemsInto(nil, gs.items)
	// out := *gs copied the token-slot interface values (pointer + type), aliasing src.
	// Make out's own per-slot instances so mutations don't cross-link the clones.
	out.initTokenSlots()
	for i := range out.tokenAuras {
		out.tokenAuras[i].SetCount(gs.tokenAuras[i].Count())
		out.tokenAuras[i].SetFiredThisTurn(gs.tokenAuras[i].FiredThisTurn())
	}
	for i := range out.tokenItems {
		out.tokenItems[i].SetCount(gs.tokenItems[i].Count())
		out.tokenItems[i].SetFiredThisTurn(gs.tokenItems[i].FiredThisTurn())
	}
	out.logger = NoopLogger{}
	return &out
}

// CopyFrom rewrites *gs in place to match what src.Copy() would produce. Reuses the
// receiver's slice / deck backings when capacity permits so a pool slot stands in for
// repeated masterState.Copy() allocations. Auras / items / weapons deep-copy per entry into
// the pool, amortising both the *GameState alloc and the per-entry *Aura / *Item / *Weapon
// allocs.
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
	pooledWeapons := gs.weapons
	pooledDeck := gs.deck
	tokenAuras := gs.tokenAuras
	tokenItems := gs.tokenItems
	*gs = *src
	// Restore the receiver's own token-slot pointers, then transcribe counts +
	// firedThisTurn from src's slots so the receiver mirrors src's token state.
	gs.tokenAuras = tokenAuras
	gs.tokenItems = tokenItems
	for i := range gs.tokenAuras {
		gs.tokenAuras[i].SetCount(src.tokenAuras[i].Count())
		gs.tokenAuras[i].SetFiredThisTurn(src.tokenAuras[i].FiredThisTurn())
	}
	for i := range gs.tokenItems {
		gs.tokenItems[i].SetCount(src.tokenItems[i].Count())
		gs.tokenItems[i].SetFiredThisTurn(src.tokenItems[i].FiredThisTurn())
	}
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
	// *gs = *src memcopied gs.auraPool byte-identically from src. Rebuild gs.auras so
	// its interface values point into gs.auraPool (not src.auraPool) at the same pool
	// indices src.auras occupied.
	gs.auras = rebuildPooledAuras(gs, src.auras, &src.auraPool[0], pooledAuras)
	gs.items = copyItemsInto(pooledItems, src.items)
	gs.weapons = copyWeaponsInto(pooledWeapons, src.weapons)
	if pooledDeck != nil {
		// Receiver owns a wrapper (pool slot prewarm); rebind its slice headers to src's
		// content rather than deep-copying. ShallowCopyFrom handles src==nil too.
		pooledDeck.ShallowCopyFrom(src.deck)
		gs.deck = pooledDeck
	} else if src.deck != nil {
		// Non-pool receivers (no prewarmed wrapper): deep copy keeps them independent of src.
		gs.deck = src.deck.Copy()
	}
}

// rebuildPooledAuras translates each entry in srcAuras from its src-pool index to the
// equivalent slot in dst.auraPool via pointer arithmetic, then CopyInto-copies the
// 96-byte aura.Aura from the src slot into the matching dst slot. Pool-index identity
// is preserved across the copy so existing src-relative invariants (currentHookIdx,
// pool-slot pointer equality between paired states) still hold.
//
// Precondition: dst's auraPool slots not covered by srcAuras must read as free
// (SourceCard nil) before the call. Stale dst slots that aren't refreshed here would
// otherwise steer the next CreateAura's free-slot scan to the wrong slot. `pooled` is
// the receiver's prior auras-slice backing (reused when cap permits).
//
// Panics if an aura's pointer falls outside the source pool — the auras list is
// pool-only by invariant; a stray heap-allocated entry would silently corrupt the
// adjacent aura field on indexed write.
func rebuildPooledAuras(dst *GameState, srcAuras []Aura, srcPoolBase *aura.Aura, pooled []Aura) []Aura {
	if len(srcAuras) == 0 {
		if pooled == nil {
			return nil
		}
		return pooled[:0]
	}
	out := pooled[:0]
	if cap(pooled) < len(srcAuras) {
		out = make([]Aura, 0, auraPoolSize)
	}
	auraSize := unsafe.Sizeof(aura.Aura{})
	base := unsafe.Pointer(srcPoolBase)
	for _, srcA := range srcAuras {
		srcPtr := unsafe.Pointer(srcA.(*aura.Aura))
		idx := (uintptr(srcPtr) - uintptr(base)) / auraSize
		if idx >= auraPoolSize {
			panic("rebuildPooledAuras: aura outside the source pool — only pool slots may appear in gs.auras")
		}
		srcA.CopyInto(&dst.auraPool[idx])
		out = append(out, &dst.auraPool[idx])
	}
	return out
}

// copyItemsInto returns a per-entry copy of src, reusing pool's slice backing when
// capacity permits. Each entry is duplicated through Item.Copy (always allocates a
// fresh *item.Item — Item has no in-place reset surface today).
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

// copyWeaponsInto returns a per-entry copy of src, reusing pool's slice backing AND its
// per-entry *weapon.Weapon allocations when present. Weapons are mutable once equipped
// (per-turn counters), so each entry is deep-copied to keep one perm's counter mutations
// from leaking into a sibling perm — but the equip loadout is identical across perms, so
// CopyInto rewrites the prior perm's slot in place (no per-perm allocation) on the hot
// CopyPersistentStateFrom path. A grown / fresh backing falls back to Copy for new slots.
func copyWeaponsInto(pool, src []Weapon) []Weapon {
	n := len(src)
	if n == 0 {
		return nil
	}
	reuse := cap(pool) >= n
	var out []Weapon
	if reuse {
		out = pool[:n]
	} else {
		out = make([]Weapon, n)
	}
	for i, w := range src {
		if reuse && out[i] != nil {
			out[i] = w.CopyInto(out[i]).(Weapon)
		} else {
			out[i] = w.Copy().(Weapon)
		}
	}
	return out
}

// resetCardSlice returns a slice of length len(src) holding src's contents, preferring to
// alias pooled's backing when capacity permits. For empty src, returns pooled[:0] — keeping
// the prewarmed backing alive so subsequent appends within the prewarmed cap don't realloc.
// (Returning nil on empty src would silently throw the backing away every time a per-perm
// state was reset from a zero-graveyard start, which dominated AppendGraveyard's allocs.)
func resetCardSlice[T any](pooled, src []T) []T {
	if len(src) == 0 {
		return pooled[:0]
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
// captured) don't bleed into the snapshot's view. Weapons, auras, and items get full
// per-entry copies because their counter / fire-this-turn fields mutate independently of
// the source.
func (gs *GameState) CopyPersistentState() *GameState {
	out := *gs
	out.weapons = copyWeaponsInto(nil, gs.weapons)
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
	// out := *gs memcopied out.auraPool from gs. Rebuild out.auras so its interface
	// values point into out.auraPool at the same pool indices gs.auras occupied.
	out.auras = rebuildPooledAuras(&out, gs.auras, &gs.auraPool[0], nil)
	out.items = copyItemsInto(nil, gs.items)
	// Token slots: out := *gs copied the slot pointers, aliasing src. Give out its own
	// instances so the snapshot's counts can't drift if src's slots mutate.
	out.initTokenSlots()
	for i := range out.tokenAuras {
		out.tokenAuras[i].SetCount(gs.tokenAuras[i].Count())
		out.tokenAuras[i].SetFiredThisTurn(gs.tokenAuras[i].FiredThisTurn())
	}
	for i := range out.tokenItems {
		out.tokenItems[i].SetCount(gs.tokenItems[i].Count())
		out.tokenItems[i].SetFiredThisTurn(gs.tokenItems[i].FiredThisTurn())
	}
	return &out
}

// CopyPersistentStateFrom overwrites *gs in place to match what CopyPersistentState(src)
// would produce. graveyard / banished copy into gs's prewarmed backing (always fits at
// the worst-case cap); auras / items reuse gs's backing when cap permits. The attack-turn
// runner's per-permutation hot path avoids slice allocation entirely.
func (gs *GameState) CopyPersistentStateFrom(src *GameState) {
	gs.hero = src.hero
	gs.heroTriggerType = src.heroTriggerType
	// Weapons deep-copy into gs's reused backing: per-perm counter mutations must isolate
	// from the master, mirroring auras / items.
	gs.weapons = copyWeaponsInto(gs.weapons, src.weapons)
	gs.arsenal = src.arsenal
	gs.graveyard = resetCardSlice(gs.graveyard, src.graveyard)
	gs.banished = resetCardSlice(gs.banished, src.banished)
	// Clear only the auraPool slots that gs.auras currently references; the rest already
	// read as free (SourceCard nil), which is all rebuildPooledAuras needs before writing
	// src's live entries into matching pool slots.
	for _, a := range gs.auras {
		a.Clear()
	}
	gs.auras = rebuildPooledAuras(gs, src.auras, &src.auraPool[0], gs.auras)
	gs.items = copyItemsInto(gs.items, src.items)
	for i := range gs.tokenAuras {
		gs.tokenAuras[i].SetCount(src.tokenAuras[i].Count())
		gs.tokenAuras[i].SetFiredThisTurn(src.tokenAuras[i].FiredThisTurn())
	}
	for i := range gs.tokenItems {
		gs.tokenItems[i].SetCount(src.tokenItems[i].Count())
		gs.tokenItems[i].SetFiredThisTurn(src.tokenItems[i].FiredThisTurn())
	}
	gs.tokenAurasLiveBits = src.tokenAurasLiveBits
	gs.tokenItemsLiveBits = src.tokenItemsLiveBits
	gs.incomingPhysicalDamage = src.incomingPhysicalDamage
	gs.incomingArcaneDamage = src.incomingArcaneDamage
	gs.opponentMarked = src.opponentMarked
	gs.isMyTurn = src.isMyTurn
}

// ResetEphemeralState returns gs to its start-of-turn baseline: it discards every field
// that playing out a turn accumulates, keeping only the cross-turn carryover (hero, deck,
// hand, arsenal, graveyard, banished, the aura / item lists, opponentMarked, isMyTurn,
// and the matchup's incoming-damage figures).
//
// gs.ephemeral.reset wipes every per-turn scratch field in one struct assignment. The
// aura / item / weapon / hero FiredThisTurn flags live on those entries themselves, not in
// ephemeral, so they get their own rearm loop here so OncePerTurn handlers can fire again.
func (gs *GameState) ResetEphemeralState() {
	gs.ephemeral.reset()
	for _, a := range gs.auras {
		a.SetFiredThisTurn(false)
	}
	for _, it := range gs.items {
		it.SetFiredThisTurn(false)
	}
	for _, w := range gs.weapons {
		w.SetFiredThisTurn(false)
	}
	for i := range gs.tokenAuras {
		if gs.tokenAuras[i] != nil {
			gs.tokenAuras[i].SetFiredThisTurn(false)
		}
	}
	for i := range gs.tokenItems {
		if gs.tokenItems[i] != nil {
			gs.tokenItems[i].SetFiredThisTurn(false)
		}
	}
	if gs.hero != nil && gs.hero.OncePerTurn() {
		gs.hero.SetFiredThisTurn(false)
	}
}

// Reset re-seeds gs with the given hero / incoming-damage values, preserving pre-allocated
// slice backings (hand, graveyard, banished, auras, items, weapons, ephemeral) so a pooled
// GameState starts a fresh shuffle without losing prewarmed backings. Weapons are emptied
// (backing kept); the caller re-equips via EquipFromCards after Reset.
func (gs *GameState) Reset(h Hero, incoming, arcaneIncoming int) {
	hand := gs.hand[:0]
	graveyard := gs.graveyard[:0]
	banished := gs.banished[:0]
	auras := gs.auras[:0]
	items := gs.items[:0]
	weapons := gs.weapons[:0]
	deckWrapper := gs.deck
	if deckWrapper != nil {
		deckWrapper.ShallowCopyFrom(nil)
	}
	tokenAuras := gs.tokenAuras
	tokenItems := gs.tokenItems
	eph := gs.ephemeral
	*gs = GameState{
		hand:                   hand,
		graveyard:              graveyard,
		banished:               banished,
		auras:                  auras,
		items:                  items,
		tokenAuras:             tokenAuras,
		tokenItems:             tokenItems,
		weapons:                weapons,
		incomingPhysicalDamage: incoming,
		incomingArcaneDamage:   arcaneIncoming,
		deck:                   deckWrapper,
		ephemeral:              eph,
	}
	gs.ephemeral.reset()
	gs.resetTokenCounts()
	gs.SetHero(h)
}

// === Pure state accessors. No cacheable flips; sim uses these to drive the attack-turn runner. ===

func (gs *GameState) Hero() Hero { return gs.hero }
func (gs *GameState) SetHero(h Hero) {
	gs.hero = h
	if h != nil {
		gs.heroTriggerType = h.TriggerType()
	} else {
		gs.heroTriggerType = 0
	}
}
func (gs *GameState) Weapons() []Weapon     { return gs.weapons }
func (gs *GameState) SetWeapons(w []Weapon) { gs.weapons = w }

// anyTriggeredWeapon reports whether any equipped weapon subscribes to a trigger event
// (TriggerType != 0). The FireTriggers early-exit gates on this rather than "any weapon
// equipped": a handler-less weapon (every current weapon) can never match a firing event,
// so it must not keep FireTriggers from short-circuiting on the common no-subscriber path.
func (gs *GameState) anyTriggeredWeapon() bool {
	for _, w := range gs.weapons {
		if w.TriggerType() != 0 {
			return true
		}
	}
	return false
}

func (gs *GameState) IsCacheable() bool   { return gs.cacheable }
func (gs *GameState) SetCacheable(v bool) { gs.cacheable = v }

// CardsRemovedFromDeck reports how many cards have been moved out of the deck during
// this attack turn resolution (mid-attack-turn draws, tutors, peek-and-banish, etc.).
func (gs *GameState) CardsRemovedFromDeck() int { return gs.cardsRemovedFromDeck }

// noteDeckRemoval increments the deck-removal counter; called by every helper that
// pops a card off the deck. Package-private — only the engine knows when a removal
// has happened.
func (gs *GameState) noteDeckRemoval(n int) { gs.cardsRemovedFromDeck += n }

// Hand projects the role-tagged hand down to the bare cards. Callers needing the roles
// (Discard, the attack-turn runner) read HandStates.
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
// by mid-attack-turn DrawOne. The value is determined by partition + attack-turn progress alone, so
// this accessor doesn't flip IsCacheable.
func (gs *GameState) HandSize() int { return len(gs.hand) }

// handStates returns the underlying CardState slice for HandHasMatching's iteration.
// Kept on GameState; the cards-facing HandHasMatching is on *GameEngine so it can pass
// itself as the engine handle into the predicate.
func (gs *GameState) handStatesForMatching() []card.CardState { return gs.hand }

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
// callers (the defense pass, the attack-turn runner) use SetHandStates.
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
// removed slot is filled by the last element (swap-with-last). The attack-turn runner
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

// AuraCount returns the number of aura entries currently in play (card-backed + token
// slots with count > 0). Used by gates like Yinti Yanti's "while you control an aura"
// rider — a Runechant or Ponder slot counts as one entry regardless of its stack count.
func (gs *GameState) AuraCount() int {
	n := len(gs.auras)
	for _, a := range gs.tokenAuras {
		if a != nil && a.Count() > 0 {
			n++
		}
	}
	return n
}

// RunechantCount / PonderCount / GoldCount / SilverCount / CopperCount return the live
// token count from the pre-allocated slot for each kind. Promoted onto *GameEngine via
// the embedded *GameState, so card-facing engine callers and TurnSummary readers share
// the same accessor.
func (gs *GameState) RunechantCount() int { return gs.tokenAuras[tokenAuraRunechant].Count() }
func (gs *GameState) PonderCount() int    { return gs.tokenAuras[tokenAuraPonder].Count() }
func (gs *GameState) QuickenCount() int   { return gs.tokenAuras[tokenAuraQuicken].Count() }
func (gs *GameState) GoldCount() int      { return gs.tokenItems[tokenItemGold].Count() }
func (gs *GameState) SilverCount() int    { return gs.tokenItems[tokenItemSilver].Count() }
func (gs *GameState) CopperCount() int    { return gs.tokenItems[tokenItemCopper].Count() }

// SetRunechantCount rewrites the Runechant slot's count directly and reconciles
// tokenAurasLiveBits with the new count. Used by the sim's pooled DR-cost probe,
// which churns the count across probes without going through the CreateRunechants
// damage-credit path.
func (gs *GameState) SetRunechantCount(n int) {
	gs.tokenAuras[tokenAuraRunechant].SetCount(n)
	if n > 0 {
		gs.tokenAurasLiveBits |= 1 << tokenAuraRunechant
	} else {
		gs.tokenAurasLiveBits &^= 1 << tokenAuraRunechant
	}
}

func (gs *GameState) Pitched() []*card.CardState     { return gs.pitched }
func (gs *GameState) SetPitched(p []*card.CardState) { gs.pitched = p }

func (gs *GameState) Defenders() []card.Card     { return gs.defenders }
func (gs *GameState) SetDefenders(d []card.Card) { gs.defenders = d }

func (gs *GameState) CardsPlayed() []card.Card               { return gs.cardsPlayed }
func (gs *GameState) SetCardsPlayed(cs []card.Card)          { gs.cardsPlayed = cs }
func (gs *GameState) AppendCardsPlayed(c card.Card)          { gs.cardsPlayed = append(gs.cardsPlayed, c) }
func (gs *GameState) CardsRemaining() []*card.CardState      { return gs.cardsRemaining }
func (gs *GameState) SetCardsRemaining(cs []*card.CardState) { gs.cardsRemaining = cs }

func (gs *GameState) Logger() card.Logger     { return gs.logger }
func (gs *GameState) SetLogger(l card.Logger) { gs.logger = l }

func (gs *GameState) AttackReactionTarget() *card.CardState      { return gs.attackReactionTarget }
func (gs *GameState) SetAttackReactionTarget(cs *card.CardState) { gs.attackReactionTarget = cs }

func (gs *GameState) ActionPoints() int     { return gs.actionPoints }
func (gs *GameState) SetActionPoints(n int) { gs.actionPoints = n }
func (gs *GameState) AddActionPoints(n int) { gs.actionPoints += n }

func (gs *GameState) Value() int     { return gs.value }
func (gs *GameState) SetValue(v int) { gs.value = v }
func (gs *GameState) AddValue(n int) { gs.value += n }

// DamageDealt returns the cumulative damage credited at every Hit / arcane-damage event
// this turn — "if you've dealt {N} this turn" gates read this. Zeroed at the turn boundary.
// Updated only by RegisterPhysicalDamage / RegisterArcaneDamage so the bookkeeping stays
// centralized.
func (gs *GameState) DamageDealt() int { return gs.damageDealt }

// HitThisTurn reports whether at least one physical attack has landed this turn. Arcane
// damage does not count — same reason it doesn't strip Mark.
func (gs *GameState) HitThisTurn() bool     { return gs.hitThisTurn }
func (gs *GameState) SetHitThisTurn(v bool) { gs.hitThisTurn = v }

func (gs *GameState) HasLightningFused() bool  { return gs.lightningFused }
func (gs *GameState) SetLightningFused(v bool) { gs.lightningFused = v }

// RegisterPhysicalDamage runs the side effects for an n-power physical attack: when the
// attack clears the likely-to-hit gate it flips HitThisTurn and credits the unblocked
// amount (per LikelyDamageDealt) to DamageDealt. Self-gating — a miss is a no-op, so
// callers don't need an outer "did this hit" guard around the call. The single
// bookkeeping point on the physical side; counterpart to RegisterArcaneDamage.
func (gs *GameState) RegisterPhysicalDamage(n int, dominate bool) {
	d := LikelyDamageDealt(n, dominate)
	if d == 0 {
		return
	}
	gs.hitThisTurn = true
	gs.damageDealt += d
}

// RegisterArcaneDamage runs the side effects for an n-arcane-damage event: flips
// arcaneDamageDealt and credits the unblocked amount to DamageDealt when n clears the
// LikelyDamageHits gate. The single bookkeeping point on the arcane side —
// DealArcaneDamage routes through here, and tokens whose Value is pre-credited at
// creation (Runechant) call it directly at fire time. Doesn't touch OpponentMarked or
// HitThisTurn — neither responds to arcane.
func (gs *GameState) RegisterArcaneDamage(n int) {
	if d := LikelyDamageDealt(n, false); d > 0 {
		gs.arcaneDamageDealt = true
		gs.damageDealt += d
	}
}

// RemainingPhysicalDamage returns the physical damage still unmitigated this turn — the
// constant matchup figure minus everything defense has blocked / prevented so far.
func (gs *GameState) RemainingPhysicalDamage() int {
	return gs.incomingPhysicalDamage - gs.physicalDamageBlocked
}

// RemainingArcaneDamage is the arcane counterpart of RemainingPhysicalDamage — the
// constant matchup arcane figure minus everything prevented so far this turn.
func (gs *GameState) RemainingArcaneDamage() int {
	return gs.incomingArcaneDamage - gs.arcaneDamageBlocked
}

// SetIncomingPhysicalDamage installs the turn's physical incoming-damage figure and zeroes BOTH
// damage-blocked accumulators — "fresh defense computation: n physical incoming, nothing
// mitigated yet on either side". Defense reactions, blocks, and prevention chip away at
// the accumulators rather than mutating the matchup figures, so the constants carry
// across turns untouched. Called at the start of every defense computation, so it also
// re-arms the arcane accumulator between per-leaf re-runs.
func (gs *GameState) SetIncomingPhysicalDamage(n int) {
	gs.incomingPhysicalDamage = n
	gs.physicalDamageBlocked = 0
	gs.arcaneDamageBlocked = 0
}

// AddPhysicalDamageBlocked credits n physical damage as absorbed by defense, shrinking
// RemainingPhysicalDamage by n. The engine's DR resolution accumulates through here; the
// attack-turn runner's plain-block pass calls it directly.
func (gs *GameState) AddPhysicalDamageBlocked(n int) { gs.physicalDamageBlocked += n }

func (gs *GameState) IncomingPhysicalDamage() int   { return gs.incomingPhysicalDamage }
func (gs *GameState) IncomingArcaneDamage() int     { return gs.incomingArcaneDamage }
func (gs *GameState) SetIncomingArcaneDamage(n int) { gs.incomingArcaneDamage = n }

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

// DestroyedOpponentArsenal reports whether the opposing arsenal has already been
// destroyed this turn. Cards mutate it via GameEngine.DestroyOpponentArsenal.
func (gs *GameState) DestroyedOpponentArsenal() bool { return gs.destroyedOpponentArsenal }

// IsMyTurn reports whether the active phase is the owning player's action phase (true) or
// the defense phase (false). The attack-turn runner sets it; cards read it for "during your
// turn" riders.
func (gs *GameState) IsMyTurn() bool     { return gs.isMyTurn }
func (gs *GameState) SetIsMyTurn(v bool) { gs.isMyTurn = v }

// HeroTapped reports whether the hero is tapped.
func (gs *GameState) HeroTapped() bool { return gs.heroTapped }

// UntapHero untaps the hero — the printed "untap your hero" effect.
func (gs *GameState) UntapHero() { gs.heroTapped = false }

func (gs *GameState) CurrentStepRerouted() bool     { return gs.currentStepRerouted }
func (gs *GameState) SetCurrentStepRerouted(v bool) { gs.currentStepRerouted = v }

// AmendLastAttackStepN adds n to the most recent AttackStep entry's N field. No-op when
// the logger is nil or when no attack-step entry exists yet.
func (gs *GameState) AmendLastAttackStepN(n int) { gs.logger.AmendLastAttackStepN(n) }

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
