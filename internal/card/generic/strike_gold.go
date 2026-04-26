// Strike Gold — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, create a Gold token."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var strikeGoldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// strikeGoldDamage returns the base attack plus the Gold-token rider when the attack is likely
// to land.
func strikeGoldDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(self) {
		return attack + card.GoldTokenValue
	}
	return attack
}

type StrikeGoldRed struct{}

func (StrikeGoldRed) ID() card.ID              { return card.StrikeGoldRed }
func (StrikeGoldRed) Name() string             { return "Strike Gold" }
func (StrikeGoldRed) Cost(*card.TurnState) int { return 0 }
func (StrikeGoldRed) Pitch() int               { return 1 }
func (StrikeGoldRed) Attack() int              { return 4 }
func (StrikeGoldRed) Defense() int             { return 2 }
func (StrikeGoldRed) Types() card.TypeSet      { return strikeGoldTypes }
func (StrikeGoldRed) GoAgain() bool            { return false }

// not implemented: gold tokens
func (StrikeGoldRed) NotImplemented() {}
func (c StrikeGoldRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, strikeGoldDamage(c.Attack(), self)-self.Card.Attack())
}

type StrikeGoldYellow struct{}

func (StrikeGoldYellow) ID() card.ID              { return card.StrikeGoldYellow }
func (StrikeGoldYellow) Name() string             { return "Strike Gold" }
func (StrikeGoldYellow) Cost(*card.TurnState) int { return 0 }
func (StrikeGoldYellow) Pitch() int               { return 2 }
func (StrikeGoldYellow) Attack() int              { return 3 }
func (StrikeGoldYellow) Defense() int             { return 2 }
func (StrikeGoldYellow) Types() card.TypeSet      { return strikeGoldTypes }
func (StrikeGoldYellow) GoAgain() bool            { return false }

// not implemented: gold tokens
func (StrikeGoldYellow) NotImplemented() {}
func (c StrikeGoldYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, strikeGoldDamage(c.Attack(), self)-self.Card.Attack())
}

type StrikeGoldBlue struct{}

func (StrikeGoldBlue) ID() card.ID              { return card.StrikeGoldBlue }
func (StrikeGoldBlue) Name() string             { return "Strike Gold" }
func (StrikeGoldBlue) Cost(*card.TurnState) int { return 0 }
func (StrikeGoldBlue) Pitch() int               { return 3 }
func (StrikeGoldBlue) Attack() int              { return 2 }
func (StrikeGoldBlue) Defense() int             { return 2 }
func (StrikeGoldBlue) Types() card.TypeSet      { return strikeGoldTypes }
func (StrikeGoldBlue) GoAgain() bool            { return false }

// not implemented: gold tokens
func (StrikeGoldBlue) NotImplemented() {}
func (c StrikeGoldBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, strikeGoldDamage(c.Attack(), self)-self.Card.Attack())
}
