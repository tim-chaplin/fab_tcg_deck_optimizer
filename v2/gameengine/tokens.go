package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Token aura / item logic. The engine knows about FaB's five built-in token kinds
// (Runechant, Ponder, Gold, Silver, Copper) because v2/card.GameEngine exposes
// CreateXxx / XxxCount methods cards call against any engine. Handlers and ability cards
// are engine bookkeeping — bumping counters, destroying entries, drawing cards on activate.

// Token name constants — exact strings the engine and CardName comparisons key on.
const (
	tokenNameRunechant = "Runechant"
	tokenNamePonder    = "Ponder"
	tokenNameGold      = "Gold"
	tokenNameSilver    = "Silver"
	tokenNameCopper    = "Copper"
)

// NewRunechantAura returns a runechant token aura at count n. Production code calls
// g.CreateRunechants instead — it bumps an existing aura and credits +n damage. This
// factory is for tests / Spec seeding that need to add a runechant aura without the
// damage credit.
func NewRunechantAura(n int) Aura {
	return NewTokenAura(tokenNameRunechant, ids.RunechantTokenID, TriggerAttack, runechantAuraHandler, n)
}

// NewPonderAura returns a ponder token aura at count n. Production code calls
// g.CreatePonder instead; this factory is for tests / Spec seeding.
func NewPonderAura(n int) Aura {
	return NewTokenAura(tokenNamePonder, ids.PonderTokenID, TriggerEndOfTurn, ponderAuraHandler, n)
}

// runechantAuraHandler is the TriggerAttack handler shared by every Runechant aura.
// Fires before each attack / weapon swing: flips ArcaneDamageDealt when aura.Count clears
// the LikelyDamageHits window and destroys the aura. Damage was credited at creation time
// in CreateRunechants — this handler is pure state cleanup.
func runechantAuraHandler(g card.GameEngine, _ card.Logger, a card.Aura) {
	eng := g.(*GameEngine)
	if LikelyDamageHits(a.Count(), false) {
		eng.arcaneDamageDealt = true
	}
	a.Destroy(false)
}

// ponderAuraHandler is the TriggerEndOfTurn handler shared by every Ponder aura. For each
// token in play it pops the deck top into the hand, letting the post-hoc arsenal-promotion
// step fill an otherwise-empty arsenal slot. Pops past deck-end are silently skipped.
func ponderAuraHandler(g card.GameEngine, _ card.Logger, a card.Aura) {
	eng := g.(*GameEngine)
	for i := 0; i < a.Count(); i++ {
		c, ok := eng.PopDeckTop()
		if !ok {
			break
		}
		eng.hand = append(eng.hand, c)
	}
	a.Destroy(false)
}

// NewGoldItem / NewSilverItem / NewCopperItem return fresh token items at count n. Production
// code calls g.CreateGold / CreateSilver / CreateCopper instead; these factories are for
// tests / Spec seeding.
func NewGoldItem(n int) Item {
	return NewTokenItem(tokenNameGold, ids.GoldTokenID, GoldTokenAbility{}, n)
}
func NewSilverItem(n int) Item {
	return NewTokenItem(tokenNameSilver, ids.SilverTokenID, SilverTokenAbility{}, n)
}
func NewCopperItem(n int) Item {
	return NewTokenItem(tokenNameCopper, ids.CopperTokenID, CopperTokenAbility{}, n)
}

// Token activated abilities. Each is a Card the chain runner enqueues as a 1-AP playable
// activated ability whenever its token is in play. Play decrements the token count and
// draws a card; token items don't head to the graveyard on destroy.

var tokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// GoldTokenAbility: cost {2}, draw a card, destroy one Gold token. Carries Go again.
type GoldTokenAbility struct{}

func (GoldTokenAbility) ID() ids.CardID                     { return ids.GoldTokenAbilityID }
func (GoldTokenAbility) Name() string                       { return tokenNameGold }
func (GoldTokenAbility) DisplayName() string                { return tokenNameGold }
func (GoldTokenAbility) Cost(card.GameEngine) int           { return 2 }
func (GoldTokenAbility) Pitch() int                         { return 0 }
func (GoldTokenAbility) Attack() int                        { return 0 }
func (GoldTokenAbility) Defense() int                       { return 0 }
func (GoldTokenAbility) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (GoldTokenAbility) GoAgain(card.GameEngine) bool       { return true }

// PlayPrecondition gates the ability on having a Gold token to spend. Rejects permutations
// that order the ability before the card / OnHit that creates the token.
func (GoldTokenAbility) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return g.GoldCount() > 0
}

func (GoldTokenAbility) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	eng := g.(*GameEngine)
	eng.consumeItem(tokenNameGold, 1)
	eng.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 gold to draw a card", 0)
}

// SilverTokenAbility: cost {3}, draw a card, destroy one Silver token.
type SilverTokenAbility struct{}

func (SilverTokenAbility) ID() ids.CardID                     { return ids.SilverTokenAbilityID }
func (SilverTokenAbility) Name() string                       { return tokenNameSilver }
func (SilverTokenAbility) DisplayName() string                { return tokenNameSilver }
func (SilverTokenAbility) Cost(card.GameEngine) int           { return 3 }
func (SilverTokenAbility) Pitch() int                         { return 0 }
func (SilverTokenAbility) Attack() int                        { return 0 }
func (SilverTokenAbility) Defense() int                       { return 0 }
func (SilverTokenAbility) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (SilverTokenAbility) GoAgain(card.GameEngine) bool       { return true }

func (SilverTokenAbility) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return g.SilverCount() > 0
}

func (SilverTokenAbility) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	eng := g.(*GameEngine)
	eng.consumeItem(tokenNameSilver, 1)
	eng.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 silver to draw a card", 0)
}

// CopperTokenAbility: cost {4}, draw a card, destroy one Copper token.
type CopperTokenAbility struct{}

func (CopperTokenAbility) ID() ids.CardID                     { return ids.CopperTokenAbilityID }
func (CopperTokenAbility) Name() string                       { return tokenNameCopper }
func (CopperTokenAbility) DisplayName() string                { return tokenNameCopper }
func (CopperTokenAbility) Cost(card.GameEngine) int           { return 4 }
func (CopperTokenAbility) Pitch() int                         { return 0 }
func (CopperTokenAbility) Attack() int                        { return 0 }
func (CopperTokenAbility) Defense() int                       { return 0 }
func (CopperTokenAbility) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (CopperTokenAbility) GoAgain(card.GameEngine) bool       { return true }

func (CopperTokenAbility) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return g.CopperCount() > 0
}

func (CopperTokenAbility) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	eng := g.(*GameEngine)
	eng.consumeItem(tokenNameCopper, 1)
	eng.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}

// === Card-facing token creation / consumption methods on GameEngine ===

// CreateRunechants creates n Runechant tokens and credits +n damage at creation. Tokens
// are stored as a single Aura entry — bump an existing entry's Count or add a new one.
// Sets AuraCreated so same-turn "aura created this turn" effects see it.
func (g *GameEngine) CreateRunechants(n int) {
	if n <= 0 {
		return
	}
	g.AddValue(n)
	g.bumpOrAppendAura(tokenNameRunechant, NewRunechantAura, n)
}

// CreatePonder creates n Ponder tokens. No Value credit — see ponderAuraHandler.
func (g *GameEngine) CreatePonder(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendAura(tokenNamePonder, NewPonderAura, n)
}

// CreateGold / CreateSilver / CreateCopper create the matching token items. No Value
// credit — these only pay out when the player activates the ability (which draws a card).
func (g *GameEngine) CreateGold(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendItem(tokenNameGold, NewGoldItem, n)
}
func (g *GameEngine) CreateSilver(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendItem(tokenNameSilver, NewSilverItem, n)
}
func (g *GameEngine) CreateCopper(n int) {
	if n <= 0 {
		return
	}
	g.bumpOrAppendItem(tokenNameCopper, NewCopperItem, n)
}

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

// consumeItem decrements the matching item's Count by n and removes the entry when Count
// reaches zero. Token items don't head to the graveyard on destroy. No-op when no item
// matches name.
func (g *GameEngine) consumeItem(name string, n int) {
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
