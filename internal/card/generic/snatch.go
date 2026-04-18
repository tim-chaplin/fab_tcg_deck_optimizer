// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."
//
// Simplification: Draw on hit isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var snatchTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SnatchRed struct{}

func (SnatchRed) ID() card.ID                 { return card.SnatchRed }
func (SnatchRed) Name() string                { return "Snatch (Red)" }
func (SnatchRed) Cost() int                   { return 0 }
func (SnatchRed) Pitch() int                  { return 1 }
func (SnatchRed) Attack() int                 { return 4 }
func (SnatchRed) Defense() int                { return 2 }
func (SnatchRed) Types() card.TypeSet         { return snatchTypes }
func (SnatchRed) GoAgain() bool               { return false }
func (c SnatchRed) Play(s *card.TurnState) int { return c.Attack() }

type SnatchYellow struct{}

func (SnatchYellow) ID() card.ID                 { return card.SnatchYellow }
func (SnatchYellow) Name() string                { return "Snatch (Yellow)" }
func (SnatchYellow) Cost() int                   { return 0 }
func (SnatchYellow) Pitch() int                  { return 2 }
func (SnatchYellow) Attack() int                 { return 3 }
func (SnatchYellow) Defense() int                { return 2 }
func (SnatchYellow) Types() card.TypeSet         { return snatchTypes }
func (SnatchYellow) GoAgain() bool               { return false }
func (c SnatchYellow) Play(s *card.TurnState) int { return c.Attack() }

type SnatchBlue struct{}

func (SnatchBlue) ID() card.ID                 { return card.SnatchBlue }
func (SnatchBlue) Name() string                { return "Snatch (Blue)" }
func (SnatchBlue) Cost() int                   { return 0 }
func (SnatchBlue) Pitch() int                  { return 3 }
func (SnatchBlue) Attack() int                 { return 2 }
func (SnatchBlue) Defense() int                { return 2 }
func (SnatchBlue) Types() card.TypeSet         { return snatchTypes }
func (SnatchBlue) GoAgain() bool               { return false }
func (c SnatchBlue) Play(s *card.TurnState) int { return c.Attack() }
