// Hand Behind the Pen — Generic Action - Attack. Cost 2, Pitch 1, Power 6, Defense 2. Only printed
// in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish a non-attack
// action card from their arsenal."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var handBehindThePenTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type HandBehindThePenRed struct{}

func (HandBehindThePenRed) ID() card.ID                 { return card.HandBehindThePenRed }
func (HandBehindThePenRed) Name() string                { return "Hand Behind the Pen" }
func (HandBehindThePenRed) Cost(*card.TurnState) int                   { return 2 }
func (HandBehindThePenRed) Pitch() int                  { return 1 }
func (HandBehindThePenRed) Attack() int                 { return 6 }
func (HandBehindThePenRed) Defense() int                { return 2 }
func (HandBehindThePenRed) Types() card.TypeSet         { return handBehindThePenTypes }
func (HandBehindThePenRed) GoAgain() bool               { return false }
// not implemented: on-hit opponent-arsenal manipulation rider
func (HandBehindThePenRed) NotImplemented()             {}
func (c HandBehindThePenRed) Play(s *card.TurnState, self *card.CardState) int { return handBehindThePenDamage(c.Attack(), self) }

// handBehindThePenDamage is a breadcrumb for the on-hit "arsenal face-up + banish non-attack
// action" rider — not modelled yet (see TODO.md).
func handBehindThePenDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(self) {
		// TODO: model on-hit arsenal manipulation rider.
	}
	return attack
}
