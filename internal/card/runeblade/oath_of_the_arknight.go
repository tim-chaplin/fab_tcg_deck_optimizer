// Oath of the Arknight — Runeblade Action. Cost 2, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gains +N{p}. Create a Runechant token. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: the +N{p} rider contributes N damage only if a Runeblade attack (an attack
// action card OR a weapon swing) follows later in this turn's ordering (peeking
// TurnState.CardsRemaining). The Runechant is always created (+1 damage) regardless.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var oathOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type OathOfTheArknightRed struct{}

func (OathOfTheArknightRed) ID() card.ID                 { return card.OathOfTheArknightRed }
func (OathOfTheArknightRed) Name() string               { return "Oath of the Arknight (Red)" }
func (OathOfTheArknightRed) Cost() int                  { return 2 }
func (OathOfTheArknightRed) Pitch() int                 { return 1 }
func (OathOfTheArknightRed) Attack() int                { return 0 }
func (OathOfTheArknightRed) Defense() int               { return 3 }
func (OathOfTheArknightRed) Types() card.TypeSet        { return oathOfTheArknightTypes }
func (OathOfTheArknightRed) GoAgain() bool              { return true }
func (OathOfTheArknightRed) Play(s *card.TurnState) int { return oathPlay(s, 3) }

type OathOfTheArknightYellow struct{}

func (OathOfTheArknightYellow) ID() card.ID                 { return card.OathOfTheArknightYellow }
func (OathOfTheArknightYellow) Name() string               { return "Oath of the Arknight (Yellow)" }
func (OathOfTheArknightYellow) Cost() int                  { return 2 }
func (OathOfTheArknightYellow) Pitch() int                 { return 2 }
func (OathOfTheArknightYellow) Attack() int                { return 0 }
func (OathOfTheArknightYellow) Defense() int               { return 3 }
func (OathOfTheArknightYellow) Types() card.TypeSet        { return oathOfTheArknightTypes }
func (OathOfTheArknightYellow) GoAgain() bool              { return true }
func (OathOfTheArknightYellow) Play(s *card.TurnState) int { return oathPlay(s, 2) }

type OathOfTheArknightBlue struct{}

func (OathOfTheArknightBlue) ID() card.ID                 { return card.OathOfTheArknightBlue }
func (OathOfTheArknightBlue) Name() string               { return "Oath of the Arknight (Blue)" }
func (OathOfTheArknightBlue) Cost() int                  { return 2 }
func (OathOfTheArknightBlue) Pitch() int                 { return 3 }
func (OathOfTheArknightBlue) Attack() int                { return 0 }
func (OathOfTheArknightBlue) Defense() int               { return 3 }
func (OathOfTheArknightBlue) Types() card.TypeSet        { return oathOfTheArknightTypes }
func (OathOfTheArknightBlue) GoAgain() bool              { return true }
func (OathOfTheArknightBlue) Play(s *card.TurnState) int { return oathPlay(s, 1) }

func oathPlay(s *card.TurnState, n int) int {
	s.AuraCreated = true
	bonus := 0
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeRuneblade) {
			continue
		}
		if t.Has(card.TypeAttack) || t.Has(card.TypeWeapon) {
			bonus = n
			break
		}
	}
	return 1 + bonus
}
