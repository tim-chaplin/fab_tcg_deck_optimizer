// Oath of the Arknight — Runeblade Action. Cost 2, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gains +N{p}. Create a Runechant token. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var oathOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type OathOfTheArknightRed struct{}

func (OathOfTheArknightRed) ID() card.ID              { return card.OathOfTheArknightRed }
func (OathOfTheArknightRed) Name() string             { return "Oath of the Arknight" }
func (OathOfTheArknightRed) Cost(*card.TurnState) int { return 2 }
func (OathOfTheArknightRed) Pitch() int               { return 1 }
func (OathOfTheArknightRed) Attack() int              { return 0 }
func (OathOfTheArknightRed) Defense() int             { return 3 }
func (OathOfTheArknightRed) Types() card.TypeSet      { return oathOfTheArknightTypes }
func (OathOfTheArknightRed) GoAgain() bool            { return true }
func (OathOfTheArknightRed) Play(s *card.TurnState, self *card.CardState) {
	oathPlay(s, self, 3)
}

type OathOfTheArknightYellow struct{}

func (OathOfTheArknightYellow) ID() card.ID              { return card.OathOfTheArknightYellow }
func (OathOfTheArknightYellow) Name() string             { return "Oath of the Arknight" }
func (OathOfTheArknightYellow) Cost(*card.TurnState) int { return 2 }
func (OathOfTheArknightYellow) Pitch() int               { return 2 }
func (OathOfTheArknightYellow) Attack() int              { return 0 }
func (OathOfTheArknightYellow) Defense() int             { return 3 }
func (OathOfTheArknightYellow) Types() card.TypeSet      { return oathOfTheArknightTypes }
func (OathOfTheArknightYellow) GoAgain() bool            { return true }
func (OathOfTheArknightYellow) Play(s *card.TurnState, self *card.CardState) {
	oathPlay(s, self, 2)
}

type OathOfTheArknightBlue struct{}

func (OathOfTheArknightBlue) ID() card.ID              { return card.OathOfTheArknightBlue }
func (OathOfTheArknightBlue) Name() string             { return "Oath of the Arknight" }
func (OathOfTheArknightBlue) Cost(*card.TurnState) int { return 2 }
func (OathOfTheArknightBlue) Pitch() int               { return 3 }
func (OathOfTheArknightBlue) Attack() int              { return 0 }
func (OathOfTheArknightBlue) Defense() int             { return 3 }
func (OathOfTheArknightBlue) Types() card.TypeSet      { return oathOfTheArknightTypes }
func (OathOfTheArknightBlue) GoAgain() bool            { return true }
func (OathOfTheArknightBlue) Play(s *card.TurnState, self *card.CardState) {
	oathPlay(s, self, 1)
}

// oathPlay grants +n to the first scheduled Runeblade attack via pc.BonusAttack so the
// buffed attack's EffectiveAttack folds the bonus into LikelyToHit and the chain credit
// lands on the target's slot, not Oath's. Always creates a Runechant token, which IS
// Oath's own contribution and lands as a sub-line under self's chain entry.
func oathPlay(s *card.TurnState, self *card.CardState, n int) {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsRunebladeAttack() {
			pc.BonusAttack += n
			break
		}
	}
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}
