// Freewheeling Renegades — Generic Action - Attack. Cost 1. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If this is defended by an action card, this has -2{p}."
//
// Simplification: The 'defended by action card' -2{p} rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var freewheelingRenegadesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FreewheelingRenegadesRed struct{}

func (FreewheelingRenegadesRed) ID() card.ID                 { return card.FreewheelingRenegadesRed }
func (FreewheelingRenegadesRed) Name() string                { return "Freewheeling Renegades (Red)" }
func (FreewheelingRenegadesRed) Cost(*card.TurnState) int                   { return 1 }
func (FreewheelingRenegadesRed) Pitch() int                  { return 1 }
func (FreewheelingRenegadesRed) Attack() int                 { return 6 }
func (FreewheelingRenegadesRed) Defense() int                { return 2 }
func (FreewheelingRenegadesRed) Types() card.TypeSet         { return freewheelingRenegadesTypes }
func (FreewheelingRenegadesRed) GoAgain() bool               { return false }
func (c FreewheelingRenegadesRed) Play(s *card.TurnState) int { return c.Attack() }

type FreewheelingRenegadesYellow struct{}

func (FreewheelingRenegadesYellow) ID() card.ID                 { return card.FreewheelingRenegadesYellow }
func (FreewheelingRenegadesYellow) Name() string                { return "Freewheeling Renegades (Yellow)" }
func (FreewheelingRenegadesYellow) Cost(*card.TurnState) int                   { return 1 }
func (FreewheelingRenegadesYellow) Pitch() int                  { return 2 }
func (FreewheelingRenegadesYellow) Attack() int                 { return 5 }
func (FreewheelingRenegadesYellow) Defense() int                { return 2 }
func (FreewheelingRenegadesYellow) Types() card.TypeSet         { return freewheelingRenegadesTypes }
func (FreewheelingRenegadesYellow) GoAgain() bool               { return false }
func (c FreewheelingRenegadesYellow) Play(s *card.TurnState) int { return c.Attack() }

type FreewheelingRenegadesBlue struct{}

func (FreewheelingRenegadesBlue) ID() card.ID                 { return card.FreewheelingRenegadesBlue }
func (FreewheelingRenegadesBlue) Name() string                { return "Freewheeling Renegades (Blue)" }
func (FreewheelingRenegadesBlue) Cost(*card.TurnState) int                   { return 1 }
func (FreewheelingRenegadesBlue) Pitch() int                  { return 3 }
func (FreewheelingRenegadesBlue) Attack() int                 { return 4 }
func (FreewheelingRenegadesBlue) Defense() int                { return 2 }
func (FreewheelingRenegadesBlue) Types() card.TypeSet         { return freewheelingRenegadesTypes }
func (FreewheelingRenegadesBlue) GoAgain() bool               { return false }
func (c FreewheelingRenegadesBlue) Play(s *card.TurnState) int { return c.Attack() }
