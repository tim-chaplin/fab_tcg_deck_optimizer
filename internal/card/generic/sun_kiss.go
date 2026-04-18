// Sun Kiss — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Gain 3{h}. If you have played a card named Moon Wish this turn, draw a card and Sun Kiss
// gains **go again**."
//
// Simplification: Health gain and the Moon Wish synergy aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sunKissTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SunKissRed struct{}

func (SunKissRed) ID() card.ID                 { return card.SunKissRed }
func (SunKissRed) Name() string                { return "Sun Kiss (Red)" }
func (SunKissRed) Cost() int                   { return 0 }
func (SunKissRed) Pitch() int                  { return 1 }
func (SunKissRed) Attack() int                 { return 0 }
func (SunKissRed) Defense() int                { return 2 }
func (SunKissRed) Types() card.TypeSet         { return sunKissTypes }
func (SunKissRed) GoAgain() bool               { return true }
func (SunKissRed) Play(s *card.TurnState) int { return 0 }

type SunKissYellow struct{}

func (SunKissYellow) ID() card.ID                 { return card.SunKissYellow }
func (SunKissYellow) Name() string                { return "Sun Kiss (Yellow)" }
func (SunKissYellow) Cost() int                   { return 0 }
func (SunKissYellow) Pitch() int                  { return 2 }
func (SunKissYellow) Attack() int                 { return 0 }
func (SunKissYellow) Defense() int                { return 2 }
func (SunKissYellow) Types() card.TypeSet         { return sunKissTypes }
func (SunKissYellow) GoAgain() bool               { return true }
func (SunKissYellow) Play(s *card.TurnState) int { return 0 }

type SunKissBlue struct{}

func (SunKissBlue) ID() card.ID                 { return card.SunKissBlue }
func (SunKissBlue) Name() string                { return "Sun Kiss (Blue)" }
func (SunKissBlue) Cost() int                   { return 0 }
func (SunKissBlue) Pitch() int                  { return 3 }
func (SunKissBlue) Attack() int                 { return 0 }
func (SunKissBlue) Defense() int                { return 2 }
func (SunKissBlue) Types() card.TypeSet         { return sunKissTypes }
func (SunKissBlue) GoAgain() bool               { return true }
func (SunKissBlue) Play(s *card.TurnState) int { return 0 }
