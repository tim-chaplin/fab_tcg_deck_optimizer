// Sun Kiss — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
// Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}. If you have played a card named Moon Wish this turn, draw a card and Sun Kiss
// gains **go again**." (N is the printed variant value above.)
//
// Modelling: health is valued 1-to-1 with damage, so Play returns +N damage-equivalent per
// variant. The Moon Wish synergy (draw a card + go again on playing Moon Wish earlier this turn)
// isn't modelled, so go-again stays off — the printed keyword is conditional, not unconditional.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sunKissTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SunKissRed struct{}

func (SunKissRed) ID() card.ID                 { return card.SunKissRed }
func (SunKissRed) Name() string                { return "Sun Kiss (Red)" }
func (SunKissRed) Cost(*card.TurnState) int                   { return 0 }
func (SunKissRed) Pitch() int                  { return 1 }
func (SunKissRed) Attack() int                 { return 0 }
func (SunKissRed) Defense() int                { return 2 }
func (SunKissRed) Types() card.TypeSet         { return sunKissTypes }
func (SunKissRed) GoAgain() bool               { return false }
// not implemented: Moon Wish synergy (draw a card + go again on playing Moon Wish earlier this turn)
func (SunKissRed) NotImplemented()             {}
func (SunKissRed) Play(s *card.TurnState, _ *card.CardState) int { return 3 }

type SunKissYellow struct{}

func (SunKissYellow) ID() card.ID                 { return card.SunKissYellow }
func (SunKissYellow) Name() string                { return "Sun Kiss (Yellow)" }
func (SunKissYellow) Cost(*card.TurnState) int                   { return 0 }
func (SunKissYellow) Pitch() int                  { return 2 }
func (SunKissYellow) Attack() int                 { return 0 }
func (SunKissYellow) Defense() int                { return 2 }
func (SunKissYellow) Types() card.TypeSet         { return sunKissTypes }
func (SunKissYellow) GoAgain() bool               { return false }
// not implemented: Moon Wish synergy (draw a card + go again on playing Moon Wish earlier this turn)
func (SunKissYellow) NotImplemented()             {}
func (SunKissYellow) Play(s *card.TurnState, _ *card.CardState) int { return 2 }

type SunKissBlue struct{}

func (SunKissBlue) ID() card.ID                 { return card.SunKissBlue }
func (SunKissBlue) Name() string                { return "Sun Kiss (Blue)" }
func (SunKissBlue) Cost(*card.TurnState) int                   { return 0 }
func (SunKissBlue) Pitch() int                  { return 3 }
func (SunKissBlue) Attack() int                 { return 0 }
func (SunKissBlue) Defense() int                { return 2 }
func (SunKissBlue) Types() card.TypeSet         { return sunKissTypes }
func (SunKissBlue) GoAgain() bool               { return false }
// not implemented: Moon Wish synergy (draw a card + go again on playing Moon Wish earlier this turn)
func (SunKissBlue) NotImplemented()             {}
func (SunKissBlue) Play(s *card.TurnState, _ *card.CardState) int { return 1 }
