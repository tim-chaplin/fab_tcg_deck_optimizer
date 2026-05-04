package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var goldTokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// GoldTokenAbility is the activated ability of a Gold token: cost 2, draw a card, destroy
// one Gold token. Token items don't head to the graveyard on destroy.
type GoldTokenAbility struct{}

func (GoldTokenAbility) ID() ids.CardID      { return ids.GoldTokenAbilityID }
func (GoldTokenAbility) Name() string        { return "Gold" }
func (GoldTokenAbility) Cost(*TurnState) int { return 2 }
func (GoldTokenAbility) Pitch() int          { return 0 }
func (GoldTokenAbility) Attack() int         { return 0 }
func (GoldTokenAbility) Defense() int        { return 0 }
func (GoldTokenAbility) Types() card.TypeSet { return goldTokenAbilityTypes }
func (GoldTokenAbility) GoAgain() bool       { return false }
func (GoldTokenAbility) Play(s *TurnState, self *CardState) {
	s.ConsumeItem(TokenTypeGold, 1)
	s.DrawOne()
	s.AddValue(DrawValue)
	s.Log(self, 0)
	s.LogRider(self, DrawValue, "Spent 1 gold to draw a card")
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
