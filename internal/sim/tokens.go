package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Token aura / item logic. The engine knows about FaB's five built-in token kinds by name
// (Runechant, Ponder, Gold, Silver, Copper) — cards call g.CreateRunechants /
// g.RunechantCount and the engine routes through sim's builders registered at init.

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
func NewRunechantAura(n int) *Aura {
	return NewTokenAura(tokenNameRunechant, ids.RunechantTokenID, gameengine.TriggerAttack, runechantAuraHandler, n)
}

// NewPonderAura returns a ponder token aura at count n. Production code calls
// g.CreatePonder instead; this factory is for tests / Spec seeding.
func NewPonderAura(n int) *Aura {
	return NewTokenAura(tokenNamePonder, ids.PonderTokenID, gameengine.TriggerEndOfTurn, ponderAuraHandler, n)
}

// runechantAuraHandler is the TriggerAttack handler shared by every Runechant aura.
// Fires before each attack / weapon swing: flips ArcaneDamageDealt when aura.Count clears
// the LikelyDamageHits window and destroys the aura. Damage was credited at creation time
// in CreateRunechants — this handler is pure state cleanup.
func runechantAuraHandler(g card.GameEngine, _ card.Logger, a card.Aura) {
	eng := g.(*gameengine.GameEngine)
	if gameengine.LikelyDamageHits(a.Count(), false) {
		eng.SetArcaneDamageDealt(true)
	}
	a.Destroy(false)
}

// ponderAuraHandler is the TriggerEndOfTurn handler shared by every Ponder aura. For each
// token in play it pops the deck top into the hand, letting the post-hoc arsenal-promotion
// step fill an otherwise-empty arsenal slot. Pops past deck-end are silently skipped.
func ponderAuraHandler(g card.GameEngine, _ card.Logger, a card.Aura) {
	eng := g.(*gameengine.GameEngine)
	for i := 0; i < a.Count(); i++ {
		c, ok := eng.PopDeckTop()
		if !ok {
			break
		}
		eng.AppendHandRaw(c)
	}
	a.Destroy(false)
}

// NewGoldItem / NewSilverItem / NewCopperItem return fresh token items at count n.
// Production code calls g.CreateGold / CreateSilver / CreateCopper instead.
func NewGoldItem(n int) *Item {
	return NewTokenItem(tokenNameGold, ids.GoldTokenID, GoldTokenAbility{}, n)
}
func NewSilverItem(n int) *Item {
	return NewTokenItem(tokenNameSilver, ids.SilverTokenID, SilverTokenAbility{}, n)
}
func NewCopperItem(n int) *Item {
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

// PlayPrecondition gates the ability on having a Gold token to spend.
func (GoldTokenAbility) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return g.GoldCount() > 0
}

func (GoldTokenAbility) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	eng := g.(*gameengine.GameEngine)
	eng.ConsumeItemByName(tokenNameGold, 1)
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
	eng := g.(*gameengine.GameEngine)
	eng.ConsumeItemByName(tokenNameSilver, 1)
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
	eng := g.(*gameengine.GameEngine)
	eng.ConsumeItemByName(tokenNameCopper, 1)
	eng.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}
