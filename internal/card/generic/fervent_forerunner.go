// Fervent Forerunner — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Fervent Forerunner hits, **opt 2**. If Fervent Forerunner is played from arsenal, it
// gains **go again**."
//
// Simplification: on-hit Opt 2 and the played-from-arsenal go-again aren't modelled. The latter
// means GoAgain() returns false unconditionally (matching how the rest of the unmodelled riders
// here default off) — a true return would let Fervent Forerunner always chain, over-crediting
// the vast majority of sequences where it wasn't played from arsenal.
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
func (FerventForerunnerRed) GoAgain() bool               { return false }
func (c FerventForerunnerRed) Play(s *card.TurnState) int { return c.Attack() }

type FerventForerunnerYellow struct{}

func (FerventForerunnerYellow) ID() card.ID                 { return card.FerventForerunnerYellow }
func (FerventForerunnerYellow) Name() string                { return "Fervent Forerunner (Yellow)" }
func (FerventForerunnerYellow) Cost() int                   { return 0 }
func (FerventForerunnerYellow) Pitch() int                  { return 2 }
func (FerventForerunnerYellow) Attack() int                 { return 2 }
func (FerventForerunnerYellow) Defense() int                { return 2 }
func (FerventForerunnerYellow) Types() card.TypeSet         { return ferventForerunnerTypes }
func (FerventForerunnerYellow) GoAgain() bool               { return false }
func (c FerventForerunnerYellow) Play(s *card.TurnState) int { return c.Attack() }

type FerventForerunnerBlue struct{}

func (FerventForerunnerBlue) ID() card.ID                 { return card.FerventForerunnerBlue }
func (FerventForerunnerBlue) Name() string                { return "Fervent Forerunner (Blue)" }
func (FerventForerunnerBlue) Cost() int                   { return 0 }
func (FerventForerunnerBlue) Pitch() int                  { return 3 }
func (FerventForerunnerBlue) Attack() int                 { return 1 }
func (FerventForerunnerBlue) Defense() int                { return 2 }
func (FerventForerunnerBlue) Types() card.TypeSet         { return ferventForerunnerTypes }
func (FerventForerunnerBlue) GoAgain() bool               { return false }
func (c FerventForerunnerBlue) Play(s *card.TurnState) int { return c.Attack() }
