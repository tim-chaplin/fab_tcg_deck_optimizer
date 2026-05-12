package gameengine

// Card-facing token creation / count methods on *GameEngine. v2/card.GameEngine requires
// these for cards to make / read FaB's five built-in tokens (Runechant, Ponder, Gold,
// Silver, Copper). The engine identifies live tokens by their CardName (the public
// canonical display name) and delegates aura / item construction to the registered
// Build*Aura / Build*Item factories so the concrete types live outside gameengine.
//
// Token tuning constants live with the heuristics; see DiscardValue / GoldTokenValue.

// Token display names — the engine matches by CardName when bumping an existing entry's
// Count or reading a count. Concrete Aura / Item impls report these strings via CardName.
const (
	tokenNameRunechant = "Runechant"
	tokenNamePonder    = "Ponder"
	tokenNameGold      = "Gold"
	tokenNameSilver    = "Silver"
	tokenNameCopper    = "Copper"
)

// CreateRunechants creates n Runechant tokens and credits +n damage at creation. Tokens
// are stored as a single Aura entry — bump an existing entry's Count or add a new one.
// Sets AuraCreated so same-turn "aura created this turn" effects see it.
func (g *GameEngine) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	g.AddValue(n)
	g.bumpOrAppendAura(tokenNameRunechant, BuildRunechantAura, n)
}

// CreatePonder creates n Ponder tokens. No Value credit — Ponder pays out at end of turn
// (see the runtime's Ponder aura handler).
func (g *GameEngine) CreatePonder(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendAura(tokenNamePonder, BuildPonderAura, n)
}

// CreateGold / CreateSilver / CreateCopper create the matching token items. No Value
// credit — items only pay out when the activated ability is spent.
func (g *GameEngine) CreateGold(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendItem(tokenNameGold, BuildGoldItem, n)
}
func (g *GameEngine) CreateSilver(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendItem(tokenNameSilver, BuildSilverItem, n)
}
func (g *GameEngine) CreateCopper(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendItem(tokenNameCopper, BuildCopperItem, n)
}

// RunechantCount / PonderCount / GoldCount / SilverCount / CopperCount return the live
// count of each token kind in play, or zero when none.
func (g *GameEngine) RunechantCount() int { return auraCountByName(g.auras, tokenNameRunechant) }
func (g *GameEngine) PonderCount() int    { return auraCountByName(g.auras, tokenNamePonder) }
func (g *GameEngine) GoldCount() int      { return itemCountByName(g.items, tokenNameGold) }
func (g *GameEngine) SilverCount() int    { return itemCountByName(g.items, tokenNameSilver) }
func (g *GameEngine) CopperCount() int    { return itemCountByName(g.items, tokenNameCopper) }

// bumpOrAppendAura increments an existing aura entry matching name, or appends a fresh
// one built by build(n). Flips auraCreated.
func (g *GameEngine) bumpOrAppendAura(name string, build func(int) Aura, n int) {
	g.auraCreated = true
	for i := range g.auras {
		if g.auras[i].CardName() == name {
			g.auras[i].SetCount(g.auras[i].Count() + n)
			return
		}
	}
	g.auras = append(g.auras, build(n))
}

// bumpOrAppendItem increments an existing item entry matching name, or appends a fresh
// one built by build(n). Items don't flip auraCreated.
func (g *GameEngine) bumpOrAppendItem(name string, build func(int) Item, n int) {
	for i := range g.items {
		if g.items[i].CardName() == name {
			g.items[i].SetCount(g.items[i].Count() + n)
			return
		}
	}
	g.items = append(g.items, build(n))
}

// ConsumeItemByName decrements the matching item's Count by n and removes the entry when
// Count reaches zero. Token items don't head to the graveyard on destroy. No-op when no
// item matches name. Called by token-ability Play implementations registered outside
// gameengine.
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
