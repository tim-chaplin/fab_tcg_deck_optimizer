// Fervent Forerunner — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Fervent Forerunner hits, **opt 2**. If Fervent Forerunner is played from arsenal, it
// gains **go again**."
//
// Simplification: On-hit Opt 2 and arsenal-only go-again aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var ferventForerunnerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FerventForerunnerRed struct{}

func (FerventForerunnerRed) ID() card.ID                 { return card.FerventForerunnerRed }
func (FerventForerunnerRed) Name() string                { return "Fervent Forerunner (Red)" }
func (FerventForerunnerRed) Cost() int                   { return 0 }
func (FerventForerunnerRed) Pitch() int                  { return 1 }
func (FerventForerunnerRed) Attack() int                 { return 3 }
func (FerventForerunnerRed) Defense() int                { return 2 }
func (FerventForerunnerRed) Types() card.TypeSet         { return ferventForerunnerTypes }
func (FerventForerunnerRed) GoAgain() bool               { return true }
func (c FerventForerunnerRed) Play(s *card.TurnState) int { return c.Attack() }

type FerventForerunnerYellow struct{}

func (FerventForerunnerYellow) ID() card.ID                 { return card.FerventForerunnerYellow }
func (FerventForerunnerYellow) Name() string                { return "Fervent Forerunner (Yellow)" }
func (FerventForerunnerYellow) Cost() int                   { return 0 }
func (FerventForerunnerYellow) Pitch() int                  { return 2 }
func (FerventForerunnerYellow) Attack() int                 { return 2 }
func (FerventForerunnerYellow) Defense() int                { return 2 }
func (FerventForerunnerYellow) Types() card.TypeSet         { return ferventForerunnerTypes }
func (FerventForerunnerYellow) GoAgain() bool               { return true }
func (c FerventForerunnerYellow) Play(s *card.TurnState) int { return c.Attack() }

type FerventForerunnerBlue struct{}

func (FerventForerunnerBlue) ID() card.ID                 { return card.FerventForerunnerBlue }
func (FerventForerunnerBlue) Name() string                { return "Fervent Forerunner (Blue)" }
func (FerventForerunnerBlue) Cost() int                   { return 0 }
func (FerventForerunnerBlue) Pitch() int                  { return 3 }
func (FerventForerunnerBlue) Attack() int                 { return 1 }
func (FerventForerunnerBlue) Defense() int                { return 2 }
func (FerventForerunnerBlue) Types() card.TypeSet         { return ferventForerunnerTypes }
func (FerventForerunnerBlue) GoAgain() bool               { return true }
func (c FerventForerunnerBlue) Play(s *card.TurnState) int { return c.Attack() }
