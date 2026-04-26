// Humble — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, they lose all hero card abilities until the end of their next
// turn."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var humbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type HumbleRed struct{}

func (HumbleRed) ID() card.ID                 { return card.HumbleRed }
func (HumbleRed) Name() string                { return "Humble" }
func (HumbleRed) Cost(*card.TurnState) int                   { return 2 }
func (HumbleRed) Pitch() int                  { return 1 }
func (HumbleRed) Attack() int                 { return 6 }
func (HumbleRed) Defense() int                { return 2 }
func (HumbleRed) Types() card.TypeSet         { return humbleTypes }
func (HumbleRed) GoAgain() bool               { return false }
// not implemented: hero-ability suppression rider
func (HumbleRed) NotImplemented()             {}
func (c HumbleRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, humbleDamage(c.Attack(), self)-self.Card.Attack())
}
type HumbleYellow struct{}

func (HumbleYellow) ID() card.ID                 { return card.HumbleYellow }
func (HumbleYellow) Name() string                { return "Humble" }
func (HumbleYellow) Cost(*card.TurnState) int                   { return 2 }
func (HumbleYellow) Pitch() int                  { return 2 }
func (HumbleYellow) Attack() int                 { return 5 }
func (HumbleYellow) Defense() int                { return 2 }
func (HumbleYellow) Types() card.TypeSet         { return humbleTypes }
func (HumbleYellow) GoAgain() bool               { return false }
// not implemented: hero-ability suppression rider
func (HumbleYellow) NotImplemented()             {}
func (c HumbleYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, humbleDamage(c.Attack(), self)-self.Card.Attack())
}
type HumbleBlue struct{}

func (HumbleBlue) ID() card.ID                 { return card.HumbleBlue }
func (HumbleBlue) Name() string                { return "Humble" }
func (HumbleBlue) Cost(*card.TurnState) int                   { return 2 }
func (HumbleBlue) Pitch() int                  { return 3 }
func (HumbleBlue) Attack() int                 { return 4 }
func (HumbleBlue) Defense() int                { return 2 }
func (HumbleBlue) Types() card.TypeSet         { return humbleTypes }
func (HumbleBlue) GoAgain() bool               { return false }
// not implemented: hero-ability suppression rider
func (HumbleBlue) NotImplemented()             {}
func (c HumbleBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, humbleDamage(c.Attack(), self)-self.Card.Attack())
}
// humbleDamage is a breadcrumb for the on-hit "lose all hero card abilities" rider — not
// modelled yet (see TODO.md).
func humbleDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(self) {
		// TODO: model on-hit hero-ability suppression rider.
	}
	return attack
}
