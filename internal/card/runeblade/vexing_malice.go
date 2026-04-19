// Vexing Malice — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 2.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 2 arcane damage to target hero."
//
// Simplification: the printed 2 arcane is added to base damage unconditionally. Play sets
// ArcaneDamageDealt so later-this-turn triggers see the flag.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var vexingMaliceTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// vexingMalicePlay adds the 2 arcane to base damage and marks ArcaneDamageDealt.
func vexingMalicePlay(attack int, s *card.TurnState) int {
	s.ArcaneDamageDealt = true
	return attack + 2
}

type VexingMaliceRed struct{}

func (VexingMaliceRed) ID() card.ID                   { return card.VexingMaliceRed }
func (VexingMaliceRed) Name() string                  { return "Vexing Malice (Red)" }
func (VexingMaliceRed) Cost(*card.TurnState) int                     { return 1 }
func (VexingMaliceRed) Pitch() int                    { return 1 }
func (VexingMaliceRed) Attack() int                   { return 3 }
func (VexingMaliceRed) Defense() int                  { return 3 }
func (VexingMaliceRed) Types() card.TypeSet           { return vexingMaliceTypes }
func (VexingMaliceRed) GoAgain() bool                 { return false }
func (c VexingMaliceRed) Play(s *card.TurnState) int  { return vexingMalicePlay(c.Attack(), s) }

type VexingMaliceYellow struct{}

func (VexingMaliceYellow) ID() card.ID                   { return card.VexingMaliceYellow }
func (VexingMaliceYellow) Name() string                  { return "Vexing Malice (Yellow)" }
func (VexingMaliceYellow) Cost(*card.TurnState) int                     { return 1 }
func (VexingMaliceYellow) Pitch() int                    { return 2 }
func (VexingMaliceYellow) Attack() int                   { return 2 }
func (VexingMaliceYellow) Defense() int                  { return 3 }
func (VexingMaliceYellow) Types() card.TypeSet           { return vexingMaliceTypes }
func (VexingMaliceYellow) GoAgain() bool                 { return false }
func (c VexingMaliceYellow) Play(s *card.TurnState) int  { return vexingMalicePlay(c.Attack(), s) }

type VexingMaliceBlue struct{}

func (VexingMaliceBlue) ID() card.ID                   { return card.VexingMaliceBlue }
func (VexingMaliceBlue) Name() string                  { return "Vexing Malice (Blue)" }
func (VexingMaliceBlue) Cost(*card.TurnState) int                     { return 1 }
func (VexingMaliceBlue) Pitch() int                    { return 3 }
func (VexingMaliceBlue) Attack() int                   { return 1 }
func (VexingMaliceBlue) Defense() int                  { return 3 }
func (VexingMaliceBlue) Types() card.TypeSet           { return vexingMaliceTypes }
func (VexingMaliceBlue) GoAgain() bool                 { return false }
func (c VexingMaliceBlue) Play(s *card.TurnState) int  { return vexingMalicePlay(c.Attack(), s) }
