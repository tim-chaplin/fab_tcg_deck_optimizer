// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per turn,
// when you play an attack action card, remove a verse counter from this. If you do, create a
// Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: front-load all N verse counters' worth of Runechants on play rather than one
// per attack action card over multiple turns. The once-per-turn and counter-depletion timing
// isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var maleficTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type MaleficIncantationRed struct{}

func (MaleficIncantationRed) ID() card.ID                 { return card.MaleficIncantationRed }
func (MaleficIncantationRed) Name() string              { return "Malefic Incantation (Red)" }
func (MaleficIncantationRed) Cost() int                 { return 0 }
func (MaleficIncantationRed) Pitch() int                { return 1 }
func (MaleficIncantationRed) Attack() int               { return 0 }
func (MaleficIncantationRed) Defense() int              { return 2 }
func (MaleficIncantationRed) Types() card.TypeSet        { return maleficTypes }
func (MaleficIncantationRed) GoAgain() bool             { return true }
func (MaleficIncantationRed) Play(s *card.TurnState) int { return s.CreateRunechants(3) }

type MaleficIncantationYellow struct{}

func (MaleficIncantationYellow) ID() card.ID                 { return card.MaleficIncantationYellow }
func (MaleficIncantationYellow) Name() string              { return "Malefic Incantation (Yellow)" }
func (MaleficIncantationYellow) Cost() int                 { return 0 }
func (MaleficIncantationYellow) Pitch() int                { return 2 }
func (MaleficIncantationYellow) Attack() int               { return 0 }
func (MaleficIncantationYellow) Defense() int              { return 2 }
func (MaleficIncantationYellow) Types() card.TypeSet        { return maleficTypes }
func (MaleficIncantationYellow) GoAgain() bool             { return true }
func (MaleficIncantationYellow) Play(s *card.TurnState) int { return s.CreateRunechants(2) }

type MaleficIncantationBlue struct{}

func (MaleficIncantationBlue) ID() card.ID                 { return card.MaleficIncantationBlue }
func (MaleficIncantationBlue) Name() string              { return "Malefic Incantation (Blue)" }
func (MaleficIncantationBlue) Cost() int                 { return 0 }
func (MaleficIncantationBlue) Pitch() int                { return 3 }
func (MaleficIncantationBlue) Attack() int               { return 0 }
func (MaleficIncantationBlue) Defense() int              { return 2 }
func (MaleficIncantationBlue) Types() card.TypeSet        { return maleficTypes }
func (MaleficIncantationBlue) GoAgain() bool             { return true }
func (MaleficIncantationBlue) Play(s *card.TurnState) int { return s.CreateRunechants(1) }
