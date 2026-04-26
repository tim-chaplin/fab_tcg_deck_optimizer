// Rosetta Thorn — Runeblade Weapon - Sword (2H). Power 2, Arcane 2.
//
// Text: "**Once per Turn Action** - {r}: **Attack**
//
// Whenever you attack with Rosetta Thorn, if you've played an attack action card and a
// 'non-attack' action card this turn, deal 2 arcane damage to target hero."

package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var rosettaThornTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type RosettaThorn struct{}

func (RosettaThorn) ID() card.ID              { return card.RosettaThornID }
func (RosettaThorn) Name() string             { return "Rosetta Thorn" }
func (RosettaThorn) Cost(*card.TurnState) int { return 1 }
func (RosettaThorn) Pitch() int               { return 0 }
func (RosettaThorn) Attack() int              { return 2 }
func (RosettaThorn) Defense() int             { return 0 }
func (RosettaThorn) Types() card.TypeSet      { return rosettaThornTypes }
func (RosettaThorn) GoAgain() bool            { return false }
func (RosettaThorn) Hands() int               { return 2 }
func (RosettaThorn) NotSilverAgeLegal()       {}

// not implemented: on-attack 2 arcane damage rider gated on having played an attack action AND
// a non-attack action this turn
func (RosettaThorn) NotImplemented() {}
func (RosettaThorn) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
