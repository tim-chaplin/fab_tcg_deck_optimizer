package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tokens are persistent-in-play entries with no originating card: destroying one just
// removes the entry (no graveyard append). Each token type has one fixed Aura handler or
// Item Ability defined here — behaviour is independent of the card that created it.
//
// Invariant: at most one entry per token type per TurnState (Auras and Items enforce this
// independently); helpers bump Count on the existing entry rather than appending. Keeps
// cache keys and the per-token loops compact.

// TokenType identifies a token kind. TokenTypeNone is the zero value used by card-based
// auras / items (which set Self.Card instead).
type TokenType int

const (
	// TokenTypeNone marks a non-token entry (Self.Card is set instead).
	TokenTypeNone TokenType = iota
	// TokenTypeRunechant is the runechant aura token. Consumed by the next attack or
	// weapon swing the controller resolves (see runechantAuraHandler).
	TokenTypeRunechant
	// TokenTypePonder is the ponder aura token. Destroys itself at end of the turn it
	// was created, drawing a card before the arsenal-promotion step (see
	// ponderAuraHandler).
	TokenTypePonder
	// TokenTypeGold is the Gold item token. Activated ability: pay {2}, destroy this
	// token, draw a card (see GoldTokenAbility).
	TokenTypeGold
	// TokenTypeSilver is the Silver item token. Same activation shape as Gold, but
	// costs {3} (see SilverTokenAbility).
	TokenTypeSilver
	// TokenTypeCopper is the Copper item token. Same activation shape as Gold, but
	// costs {4} (see CopperTokenAbility).
	TokenTypeCopper
)

// tokenDisplayName returns the printed name shown in logs and "(from previous turn)"
// summaries for the given token type. Mirrors Card.DisplayName() for card-based entries.
func tokenDisplayName(t TokenType) string {
	switch t {
	case TokenTypeRunechant:
		return "Runechant"
	case TokenTypePonder:
		return "Ponder"
	case TokenTypeGold:
		return "Gold"
	case TokenTypeSilver:
		return "Silver"
	case TokenTypeCopper:
		return "Copper"
	}
	return ""
}

// runechantAuraHandler is the TriggerAttack handler shared by every Runechant aura.
// Fires before each attack / weapon swing resolves: flips ArcaneDamageDealt when
// aura.Count clears the LikelyDamageHits window and destroys the aura. Damage was
// credited at creation time in CreateRunechants — this handler is pure state cleanup.
func runechantAuraHandler(s *TurnState, _ Logger, _ *Trigger, a *Aura) {
	if LikelyDamageHits(a.Count, false) {
		s.arcaneDamageDealt = true
	}
	s.DestroyAura(a, false)
}

// NewRunechantAura returns a runechant token aura at count n. Production code calls
// s.CreateRunechants instead — it bumps an existing aura and credits +n damage. This
// factory is for tests that need to seed a runechant aura without the damage credit.
func NewRunechantAura(n int) Aura {
	return Aura{
		Trigger: Trigger{TriggerType: TriggerAttack, Handler: runechantAuraHandler},
		Self:    CardOrTokenType{TokenType: TokenTypeRunechant},
		Count:   n,
	}
}

// ponderAuraHandler is the TriggerEndOfTurn handler shared by every Ponder aura. For
// each token in play it pops the deck top into the held cards (s.Hand), letting the
// post-hoc arsenal-promotion step fill an otherwise-empty arsenal slot. Pops past
// deck-end are silently skipped — empty deck just means no draw. Reading the deck top
// flips s.cacheable (PopDeckTop's contract).
func ponderAuraHandler(s *TurnState, _ Logger, _ *Trigger, a *Aura) {
	for i := 0; i < a.Count; i++ {
		c, ok := s.PopDeckTop()
		if !ok {
			break
		}
		s.hand = append(s.hand, c)
	}
	s.DestroyAura(a, false)
}

// NewPonderAura returns a ponder token aura at count n. Production code calls
// s.CreatePonder instead; this factory is for tests that need to seed the aura directly.
func NewPonderAura(n int) Aura {
	return Aura{
		Trigger: Trigger{TriggerType: TriggerEndOfTurn, Handler: ponderAuraHandler},
		Self:    CardOrTokenType{TokenType: TokenTypePonder},
		Count:   n,
	}
}

// tokenCountIn scans an aura slice for a token aura of the given type and returns its
// count. Single read site for the at-most-one-aura-per-token-type invariant.
func tokenCountIn(auras []Aura, t TokenType) int {
	for i := range auras {
		if auras[i].Self.TokenType == t {
			return auras[i].Count
		}
	}
	return 0
}

var goldTokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// GoldTokenAbility is the activated ability of a Gold token: cost {2}, draw a card,
// destroy one Gold token. Carries Go again so the action chain can fold gold-spending
// in alongside an attack. Token items don't head to the graveyard on destroy.
type GoldTokenAbility struct{}

func (GoldTokenAbility) ID() ids.CardID      { return ids.GoldTokenAbilityID }
func (GoldTokenAbility) Name() string        { return "Gold" }
func (GoldTokenAbility) DisplayName() string { return "Gold" }
func (GoldTokenAbility) Cost(GameEngine) int { return 2 }
func (GoldTokenAbility) Pitch() int          { return 0 }
func (GoldTokenAbility) Attack() int         { return 0 }
func (GoldTokenAbility) Defense() int        { return 0 }
func (GoldTokenAbility) Types() card.TypeSet { return goldTokenAbilityTypes }
func (GoldTokenAbility) GoAgain() bool       { return true }

// PlayPrecondition gates the activated ability on having a Gold token to spend. Rejects
// permutations that order the ability before the card / OnHit that creates the token —
// the chain runner finds the legal ordering (token created first) via Heap's algorithm.
func (GoldTokenAbility) PlayPrecondition(s *TurnState, self *CardState) bool {
	return s.Gold() > 0
}

func (GoldTokenAbility) Play(g GameEngine, l Logger, self *CardState) {
	s := g.(*TurnState)
	s.ConsumeItem(TokenTypeGold, 1)
	s.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 gold to draw a card", 0)
}

// NewGoldItem returns a fresh Gold token Item with the given count. Production code calls
// s.CreateGold instead — it bumps an existing entry. Test seeding only.
func NewGoldItem(n int) Item {
	return Item{
		Self:    CardOrTokenType{TokenType: TokenTypeGold},
		Count:   n,
		Ability: GoldTokenAbility{},
	}
}

var silverTokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// SilverTokenAbility is the activated ability of a Silver token: pay {3}, destroy one
// Silver token, draw a card. Same shape as GoldTokenAbility — token items don't head
// to the graveyard on destroy.
type SilverTokenAbility struct{}

func (SilverTokenAbility) ID() ids.CardID      { return ids.SilverTokenAbilityID }
func (SilverTokenAbility) Name() string        { return "Silver" }
func (SilverTokenAbility) DisplayName() string { return "Silver" }
func (SilverTokenAbility) Cost(GameEngine) int { return 3 }
func (SilverTokenAbility) Pitch() int          { return 0 }
func (SilverTokenAbility) Attack() int         { return 0 }
func (SilverTokenAbility) Defense() int        { return 0 }
func (SilverTokenAbility) Types() card.TypeSet { return silverTokenAbilityTypes }
func (SilverTokenAbility) GoAgain() bool       { return true }

// PlayPrecondition gates the activated ability on having a Silver token to spend.
func (SilverTokenAbility) PlayPrecondition(s *TurnState, self *CardState) bool {
	return s.Silver() > 0
}

func (SilverTokenAbility) Play(g GameEngine, l Logger, self *CardState) {
	s := g.(*TurnState)
	s.ConsumeItem(TokenTypeSilver, 1)
	s.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 silver to draw a card", 0)
}

// NewSilverItem returns a fresh Silver token Item with the given count. Production code
// calls s.CreateSilver instead — it bumps an existing entry. Test seeding only.
func NewSilverItem(n int) Item {
	return Item{
		Self:    CardOrTokenType{TokenType: TokenTypeSilver},
		Count:   n,
		Ability: SilverTokenAbility{},
	}
}

var copperTokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// CopperTokenAbility is the activated ability of a Copper token: pay {4}, destroy one
// Copper token, draw a card. Same shape as GoldTokenAbility — token items don't head
// to the graveyard on destroy.
type CopperTokenAbility struct{}

func (CopperTokenAbility) ID() ids.CardID      { return ids.CopperTokenAbilityID }
func (CopperTokenAbility) Name() string        { return "Copper" }
func (CopperTokenAbility) DisplayName() string { return "Copper" }
func (CopperTokenAbility) Cost(GameEngine) int { return 4 }
func (CopperTokenAbility) Pitch() int          { return 0 }
func (CopperTokenAbility) Attack() int         { return 0 }
func (CopperTokenAbility) Defense() int        { return 0 }
func (CopperTokenAbility) Types() card.TypeSet { return copperTokenAbilityTypes }
func (CopperTokenAbility) GoAgain() bool       { return true }

// PlayPrecondition gates the activated ability on having a Copper token to spend.
func (CopperTokenAbility) PlayPrecondition(s *TurnState, self *CardState) bool {
	return s.Copper() > 0
}

func (CopperTokenAbility) Play(g GameEngine, l Logger, self *CardState) {
	s := g.(*TurnState)
	s.ConsumeItem(TokenTypeCopper, 1)
	s.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}

// NewCopperItem returns a fresh Copper token Item with the given count. Production code
// calls s.CreateCopper instead — it bumps an existing entry. Test seeding only.
func NewCopperItem(n int) Item {
	return Item{
		Self:    CardOrTokenType{TokenType: TokenTypeCopper},
		Count:   n,
		Ability: CopperTokenAbility{},
	}
}
