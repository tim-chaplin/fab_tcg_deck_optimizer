package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/internal/token"

// Token kinds — sequential indices into GameState.tokenAuras / tokenItems. Every pool
// GameState pre-allocates one Aura / Item per kind (count=0); creating tokens of a
// kind increments the count, destruction zeroes it, and reset zeroes all counts. The
// instance pointer is stable for the GameState's lifetime, so trigger dispatch fires
// the same Fire closure across an entire run with no per-creation allocation.

type tokenAuraKind int

const (
	tokenAuraRunechant tokenAuraKind = iota
	tokenAuraPonder
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
	gs.tokenItems[tokenItemGold] = token.NewGold(0)
	gs.tokenItems[tokenItemSilver] = token.NewSilver(0)
	gs.tokenItems[tokenItemCopper] = token.NewCopper(0)
}

// resetTokenCounts zeroes the Count on every token slot and clears their firedThisTurn
// flags. Called by ResetEphemeralState (firedThisTurn) and by Reset (counts + fired).
func (gs *GameState) resetTokenCounts() {
	for i := range gs.tokenAuras {
		gs.tokenAuras[i].SetCount(0)
		gs.tokenAuras[i].SetFiredThisTurn(false)
	}
	for i := range gs.tokenItems {
		gs.tokenItems[i].SetCount(0)
		gs.tokenItems[i].SetFiredThisTurn(false)
	}
}

// bumpTokenAura increments a token aura's count by n; re-arms firedThisTurn when the
// transition takes it from zero to positive so a same-turn creation can fire its
// trigger this turn. Also flips auraCreated for "an aura was created this turn"
// readers.
func (gs *GameState) bumpTokenAura(kind tokenAuraKind, n int) {
	a := gs.tokenAuras[kind]
	prev := a.Count()
	a.SetCount(prev + n)
	if prev == 0 {
		a.SetFiredThisTurn(false)
	}
	gs.auraCreated = true
}

// bumpTokenItem is the item counterpart of bumpTokenAura: re-arms firedThisTurn on the
// 0→positive transition so a same-turn creation can fire its trigger this turn. No
// auraCreated flip — items don't gate any "aura created this turn" reader.
func (gs *GameState) bumpTokenItem(kind tokenItemKind, n int) {
	it := gs.tokenItems[kind]
	prev := it.Count()
	it.SetCount(prev + n)
	if prev == 0 {
		it.SetFiredThisTurn(false)
	}
}

// TotalAuraCount returns the count summed across card-backed auras + every active
// token-aura slot. Used by the optimizer's tiebreaker scoring and by "are any auras
// in play" gates.
func (gs *GameState) TotalAuraCount() int {
	total := 0
	for _, a := range gs.auras {
		total += a.Count()
	}
	for _, a := range gs.tokenAuras {
		if a != nil {
			total += a.Count()
		}
	}
	return total
}

// TotalItemCount mirrors TotalAuraCount for items.
func (gs *GameState) TotalItemCount() int {
	total := 0
	for _, it := range gs.items {
		total += it.Count()
	}
	for _, it := range gs.tokenItems {
		if it != nil {
			total += it.Count()
		}
	}
	return total
}

// AnyAurasInPlay reports whether any aura entry (card-backed or token, with count > 0)
// would be visible to a FireTriggers walk.
func (gs *GameState) AnyAurasInPlay() bool {
	if len(gs.auras) > 0 {
		return true
	}
	for _, a := range gs.tokenAuras {
		if a != nil && a.Count() > 0 {
			return true
		}
	}
	return false
}

// AnyItemsInPlay reports whether any item entry (card-backed or token with count > 0)
// would be visible to a FireTriggers walk.
func (gs *GameState) AnyItemsInPlay() bool {
	if len(gs.items) > 0 {
		return true
	}
	for _, it := range gs.tokenItems {
		if it != nil && it.Count() > 0 {
			return true
		}
	}
	return false
}

// ForEachTokenAura yields every token-aura slot with a positive count to fn. Mirrors
// ForEachTokenItem for token auras; used by the sim's start-of-turn aura-processing
// reporter so test fixtures can verify that token slots like Ponder weren't fired off
// the wrong event.
func (gs *GameState) ForEachTokenAura(fn func(Aura)) {
	for _, a := range gs.tokenAuras {
		if a != nil && a.Count() > 0 {
			fn(a)
		}
	}
}

// ForEachTokenItem yields every token-item slot with a positive count to fn. Used by
// the sim's attack-turn ability enumeration so activated token abilities (Gold / Silver /
// Copper Spend) enqueue as playables alongside card-backed item abilities.
func (gs *GameState) ForEachTokenItem(fn func(Item)) {
	for _, it := range gs.tokenItems {
		if it != nil && it.Count() > 0 {
			fn(it)
		}
	}
}

