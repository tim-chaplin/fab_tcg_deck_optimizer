package gameengine

import (
	"math/bits"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// Token kinds — sequential indices into GameState.tokenAuras / tokenItems. Every pool
// GameState pre-allocates one Aura / Item per kind (count=0); creating tokens of a
// kind increments the count, destruction zeroes it, and reset zeroes all counts. The
// instance pointer is stable for the GameState's lifetime, so trigger dispatch fires
// the same Fire closure across an entire run with no per-creation allocation.

type tokenAuraKind int

const (
	tokenAuraRunechant tokenAuraKind = iota
	tokenAuraPonder
	tokenAuraQuicken
	numTokenAuraKinds
)

type tokenItemKind int

const (
	tokenItemGold tokenItemKind = iota
	tokenItemSilver
	tokenItemCopper
	numTokenItemKinds
)

// initTokenSlots populates gs.tokenAuras / tokenItems with the canonical instance per
// kind at count=0. Called by NewPrewarmedState and GameStateBuilder so every constructed
// GameState carries a usable token slot for each kind.
func (gs *GameState) initTokenSlots() {
	gs.tokenAuras[tokenAuraRunechant] = token.NewRunechant(0)
	gs.tokenAuras[tokenAuraPonder] = token.NewPonder(0)
	gs.tokenAuras[tokenAuraQuicken] = token.NewQuicken(0)
	gs.tokenItems[tokenItemGold] = token.NewGold(0)
	gs.tokenItems[tokenItemSilver] = token.NewSilver(0)
	gs.tokenItems[tokenItemCopper] = token.NewCopper(0)
}

// resetTokenCounts zeroes Count + firedThisTurn on every token slot and clears both
// live bitmasks — the canonical "no tokens in play" reset.
func (gs *GameState) resetTokenCounts() {
	for i := range gs.tokenAuras {
		gs.tokenAuras[i].SetCount(0)
		gs.tokenAuras[i].SetFiredThisTurn(false)
	}
	for i := range gs.tokenItems {
		gs.tokenItems[i].SetCount(0)
		gs.tokenItems[i].SetFiredThisTurn(false)
	}
	gs.tokenAurasLiveBits = 0
	gs.tokenItemsLiveBits = 0
}

// bumpTokenAura increments a token aura's count by n; re-arms firedThisTurn when the
// transition takes it from zero to positive so a same-turn creation can fire its
// trigger this turn. Also flips auraCreated for "an aura was created this turn"
// readers, and sets the matching bit in tokenAurasLiveBits.
func (gs *GameState) bumpTokenAura(kind tokenAuraKind, n int) {
	a := gs.tokenAuras[kind]
	prev := a.Count()
	a.SetCount(prev + n)
	if prev == 0 {
		a.SetFiredThisTurn(false)
	}
	gs.tokenAurasLiveBits |= 1 << kind
	gs.auraCreated = true
}

// bumpTokenItem is the item counterpart of bumpTokenAura: re-arms firedThisTurn on the
// 0→positive transition so a same-turn creation can fire its trigger this turn, and
// sets the matching bit in tokenItemsLiveBits. No auraCreated flip — items don't gate
// any "aura created this turn" reader.
func (gs *GameState) bumpTokenItem(kind tokenItemKind, n int) {
	it := gs.tokenItems[kind]
	prev := it.Count()
	it.SetCount(prev + n)
	if prev == 0 {
		it.SetFiredThisTurn(false)
	}
	gs.tokenItemsLiveBits |= 1 << kind
}

// TotalAuraCount returns the count summed across card-backed auras + every active
// token-aura slot. Used by the optimizer's tiebreaker scoring and by "are any auras
// in play" gates. Skips dead token slots via tokenAurasLiveBits so live-empty pools
// (the common case) read one byte instead of walking all kinds.
func (gs *GameState) TotalAuraCount() int {
	total := 0
	for _, a := range gs.auras {
		total += a.Count()
	}
	mask := gs.tokenAurasLiveBits
	for mask != 0 {
		low := mask & -mask // lowest set bit
		total += gs.tokenAuras[bits.TrailingZeros8(low)].Count()
		mask ^= low
	}
	return total
}

// TotalItemCount mirrors TotalAuraCount for items.
func (gs *GameState) TotalItemCount() int {
	total := 0
	for _, it := range gs.items {
		total += it.Count()
	}
	mask := gs.tokenItemsLiveBits
	for mask != 0 {
		low := mask & -mask
		total += gs.tokenItems[bits.TrailingZeros8(low)].Count()
		mask ^= low
	}
	return total
}

// AnyAurasInPlay reports whether any aura entry (card-backed or token, with count > 0)
// would be visible to a FireTriggers walk.
func (gs *GameState) AnyAurasInPlay() bool {
	return len(gs.auras) > 0 || gs.tokenAurasLiveBits != 0
}

// AnyItemsInPlay reports whether any item entry (card-backed or token with count > 0)
// would be visible to a FireTriggers walk.
func (gs *GameState) AnyItemsInPlay() bool {
	return len(gs.items) > 0 || gs.tokenItemsLiveBits != 0
}

// ForEachTokenAura yields every token-aura slot with a positive count to fn. Mirrors
// ForEachTokenItem for token auras; used by the sim's start-of-turn aura-processing
// reporter so test fixtures can verify that token slots like Ponder weren't fired off
// the wrong event.
func (gs *GameState) ForEachTokenAura(fn func(Aura)) {
	mask := gs.tokenAurasLiveBits
	for mask != 0 {
		low := mask & -mask
		fn(gs.tokenAuras[bits.TrailingZeros8(low)])
		mask ^= low
	}
}

// ForEachTokenItem yields every token-item slot with a positive count to fn. Used by
// the sim's attack-turn ability enumeration so activated token abilities (Gold / Silver /
// Copper Spend) enqueue as playables alongside card-backed item abilities.
func (gs *GameState) ForEachTokenItem(fn func(Item)) {
	mask := gs.tokenItemsLiveBits
	for mask != 0 {
		low := mask & -mask
		fn(gs.tokenItems[bits.TrailingZeros8(low)])
		mask ^= low
	}
}
