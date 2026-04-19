// Arcanic Crackle — Runeblade Action - Attack. Cost 0, Defense 3, Arcane 1.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to target hero."
//
// Simplification: the printed 1 arcane is added to base damage unconditionally. Play sets
// ArcaneDamageDealt so later-this-turn triggers reading "if you've dealt arcane damage this
// turn" see it.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcanicCrackleTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// arcanicCracklePlay adds the printed +1 arcane to base damage and marks ArcaneDamageDealt.
func arcanicCracklePlay(attack int, s *card.TurnState) int {
	s.ArcaneDamageDealt = true
	return attack + 1
}

type ArcanicCrackleRed struct{}

func (ArcanicCrackleRed) ID() card.ID                   { return card.ArcanicCrackleRed }
func (ArcanicCrackleRed) Name() string                  { return "Arcanic Crackle (Red)" }
func (ArcanicCrackleRed) Cost(*card.TurnState) int                     { return 0 }
func (ArcanicCrackleRed) Pitch() int                    { return 1 }
func (ArcanicCrackleRed) Attack() int                   { return 3 }
func (ArcanicCrackleRed) Defense() int                  { return 3 }
func (ArcanicCrackleRed) Types() card.TypeSet           { return arcanicCrackleTypes }
func (ArcanicCrackleRed) GoAgain() bool                 { return false }
func (c ArcanicCrackleRed) Play(s *card.TurnState) int  { return arcanicCracklePlay(c.Attack(), s) }

type ArcanicCrackleYellow struct{}

func (ArcanicCrackleYellow) ID() card.ID                   { return card.ArcanicCrackleYellow }
func (ArcanicCrackleYellow) Name() string                  { return "Arcanic Crackle (Yellow)" }
func (ArcanicCrackleYellow) Cost(*card.TurnState) int                     { return 0 }
func (ArcanicCrackleYellow) Pitch() int                    { return 2 }
func (ArcanicCrackleYellow) Attack() int                   { return 2 }
func (ArcanicCrackleYellow) Defense() int                  { return 3 }
func (ArcanicCrackleYellow) Types() card.TypeSet           { return arcanicCrackleTypes }
func (ArcanicCrackleYellow) GoAgain() bool                 { return false }
func (c ArcanicCrackleYellow) Play(s *card.TurnState) int  { return arcanicCracklePlay(c.Attack(), s) }

type ArcanicCrackleBlue struct{}

func (ArcanicCrackleBlue) ID() card.ID                   { return card.ArcanicCrackleBlue }
func (ArcanicCrackleBlue) Name() string                  { return "Arcanic Crackle (Blue)" }
func (ArcanicCrackleBlue) Cost(*card.TurnState) int                     { return 0 }
func (ArcanicCrackleBlue) Pitch() int                    { return 3 }
func (ArcanicCrackleBlue) Attack() int                    { return 1 }
func (ArcanicCrackleBlue) Defense() int                   { return 3 }
func (ArcanicCrackleBlue) Types() card.TypeSet            { return arcanicCrackleTypes }
func (ArcanicCrackleBlue) GoAgain() bool                  { return false }
func (c ArcanicCrackleBlue) Play(s *card.TurnState) int   { return arcanicCracklePlay(c.Attack(), s) }
