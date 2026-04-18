// Frontline Scout — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may look at the defending hero's hand. If Frontline Scout is played from arsenal, it
// gains **go again**."
//
// Simplification: Hand-peek and arsenal-only go-again aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var frontlineScoutTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FrontlineScoutRed struct{}

func (FrontlineScoutRed) ID() card.ID                 { return card.FrontlineScoutRed }
func (FrontlineScoutRed) Name() string                { return "Frontline Scout (Red)" }
func (FrontlineScoutRed) Cost() int                   { return 0 }
func (FrontlineScoutRed) Pitch() int                  { return 1 }
func (FrontlineScoutRed) Attack() int                 { return 3 }
func (FrontlineScoutRed) Defense() int                { return 2 }
func (FrontlineScoutRed) Types() card.TypeSet         { return frontlineScoutTypes }
func (FrontlineScoutRed) GoAgain() bool               { return false }
func (c FrontlineScoutRed) Play(s *card.TurnState) int { return c.Attack() }

type FrontlineScoutYellow struct{}

func (FrontlineScoutYellow) ID() card.ID                 { return card.FrontlineScoutYellow }
func (FrontlineScoutYellow) Name() string                { return "Frontline Scout (Yellow)" }
func (FrontlineScoutYellow) Cost() int                   { return 0 }
func (FrontlineScoutYellow) Pitch() int                  { return 2 }
func (FrontlineScoutYellow) Attack() int                 { return 2 }
func (FrontlineScoutYellow) Defense() int                { return 2 }
func (FrontlineScoutYellow) Types() card.TypeSet         { return frontlineScoutTypes }
func (FrontlineScoutYellow) GoAgain() bool               { return false }
func (c FrontlineScoutYellow) Play(s *card.TurnState) int { return c.Attack() }

type FrontlineScoutBlue struct{}

func (FrontlineScoutBlue) ID() card.ID                 { return card.FrontlineScoutBlue }
func (FrontlineScoutBlue) Name() string                { return "Frontline Scout (Blue)" }
func (FrontlineScoutBlue) Cost() int                   { return 0 }
func (FrontlineScoutBlue) Pitch() int                  { return 3 }
func (FrontlineScoutBlue) Attack() int                 { return 1 }
func (FrontlineScoutBlue) Defense() int                { return 2 }
func (FrontlineScoutBlue) Types() card.TypeSet         { return frontlineScoutTypes }
func (FrontlineScoutBlue) GoAgain() bool               { return false }
func (c FrontlineScoutBlue) Play(s *card.TurnState) int { return c.Attack() }
