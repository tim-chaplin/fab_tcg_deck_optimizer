// Scepter of Pain — Runeblade Weapon - Scepter (1H). Cost 2, Arcane 1.
// Text: "Once per Turn Action - {r}{r}: Deal 1 arcane damage to any opposing target. Create a
// Runechant token for each damage dealt this way."
//
// Simplification: modelled as an attack source dealing 1 arcane + 1 Runechant (+1 future damage,
// per the Malefic convention) = 2 damage total. The ability is not strictly an attack in FaB terms
// (the card has no "Attack" type), but the simulator treats any weapon swing as the turn's
// damage-dealing action.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var scepterOfPainTypes = map[string]bool{
	"Runeblade": true,
	"Weapon":    true,
	"Scepter":   true,
	"1H":        true,
}

type ScepterOfPain struct{}

func (ScepterOfPain) Name() string                 { return "Scepter of Pain" }
func (ScepterOfPain) Cost() int                    { return 2 }
func (ScepterOfPain) Pitch() int                   { return 0 }
func (ScepterOfPain) Attack() int                  { return 2 }
func (ScepterOfPain) Defense() int                 { return 0 }
func (ScepterOfPain) Types() map[string]bool       { return scepterOfPainTypes }
func (ScepterOfPain) GoAgain() bool                { return false }
func (ScepterOfPain) Hands() int                   { return 1 }
func (c ScepterOfPain) Play(*card.TurnState) int   { return c.Attack() }
