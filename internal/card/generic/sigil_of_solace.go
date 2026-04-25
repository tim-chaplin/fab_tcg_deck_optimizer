// Sigil of Solace — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 3{h}"
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSolaceTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type SigilOfSolaceRed struct{}

func (SigilOfSolaceRed) ID() card.ID                               { return card.SigilOfSolaceRed }
func (SigilOfSolaceRed) Name() string                              { return "Sigil of Solace (Red)" }
func (SigilOfSolaceRed) Cost(*card.TurnState) int                  { return 0 }
func (SigilOfSolaceRed) Pitch() int                                { return 1 }
func (SigilOfSolaceRed) Attack() int                               { return 0 }
func (SigilOfSolaceRed) Defense() int                              { return 0 }
func (SigilOfSolaceRed) Types() card.TypeSet                       { return sigilOfSolaceTypes }
func (SigilOfSolaceRed) GoAgain() bool                             { return false }
func (SigilOfSolaceRed) NotSilverAgeLegal()                        {}
func (SigilOfSolaceRed) NotImplemented()                           {}
func (SigilOfSolaceRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type SigilOfSolaceYellow struct{}

func (SigilOfSolaceYellow) ID() card.ID                               { return card.SigilOfSolaceYellow }
func (SigilOfSolaceYellow) Name() string                              { return "Sigil of Solace (Yellow)" }
func (SigilOfSolaceYellow) Cost(*card.TurnState) int                  { return 0 }
func (SigilOfSolaceYellow) Pitch() int                                { return 2 }
func (SigilOfSolaceYellow) Attack() int                               { return 0 }
func (SigilOfSolaceYellow) Defense() int                              { return 0 }
func (SigilOfSolaceYellow) Types() card.TypeSet                       { return sigilOfSolaceTypes }
func (SigilOfSolaceYellow) GoAgain() bool                             { return false }
func (SigilOfSolaceYellow) NotSilverAgeLegal()                        {}
func (SigilOfSolaceYellow) NotImplemented()                           {}
func (SigilOfSolaceYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type SigilOfSolaceBlue struct{}

func (SigilOfSolaceBlue) ID() card.ID                               { return card.SigilOfSolaceBlue }
func (SigilOfSolaceBlue) Name() string                              { return "Sigil of Solace (Blue)" }
func (SigilOfSolaceBlue) Cost(*card.TurnState) int                  { return 0 }
func (SigilOfSolaceBlue) Pitch() int                                { return 3 }
func (SigilOfSolaceBlue) Attack() int                               { return 0 }
func (SigilOfSolaceBlue) Defense() int                              { return 0 }
func (SigilOfSolaceBlue) Types() card.TypeSet                       { return sigilOfSolaceTypes }
func (SigilOfSolaceBlue) GoAgain() bool                             { return false }
func (SigilOfSolaceBlue) NotSilverAgeLegal()                        {}
func (SigilOfSolaceBlue) NotImplemented()                           {}
func (SigilOfSolaceBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
