// Fate Foreseen — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Opt 1"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type FateForeseenRed struct{}

func (FateForeseenRed) ID() card.ID                 { return card.FateForeseenRed }
func (FateForeseenRed) Name() string             { return "Fate Foreseen (Red)" }
func (FateForeseenRed) Cost(*card.TurnState) int                { return 0 }
func (FateForeseenRed) Pitch() int               { return 1 }
func (FateForeseenRed) Attack() int              { return 0 }
func (FateForeseenRed) Defense() int             { return 4 }
func (FateForeseenRed) Types() card.TypeSet      { return defenseReactionTypes }
func (FateForeseenRed) GoAgain() bool            { return false }
func (FateForeseenRed) NotSilverAgeLegal()       {}
// not implemented: Opt 1 rider; block value is printed defence only
func (FateForeseenRed) NotImplemented()           {}
func (FateForeseenRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type FateForeseenYellow struct{}

func (FateForeseenYellow) ID() card.ID                 { return card.FateForeseenYellow }
func (FateForeseenYellow) Name() string             { return "Fate Foreseen (Yellow)" }
func (FateForeseenYellow) Cost(*card.TurnState) int                { return 0 }
func (FateForeseenYellow) Pitch() int               { return 2 }
func (FateForeseenYellow) Attack() int              { return 0 }
func (FateForeseenYellow) Defense() int             { return 3 }
func (FateForeseenYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (FateForeseenYellow) GoAgain() bool            { return false }
func (FateForeseenYellow) NotSilverAgeLegal()       {}
// not implemented: Opt 1 rider; block value is printed defence only
func (FateForeseenYellow) NotImplemented()           {}
func (FateForeseenYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type FateForeseenBlue struct{}

func (FateForeseenBlue) ID() card.ID                 { return card.FateForeseenBlue }
func (FateForeseenBlue) Name() string             { return "Fate Foreseen (Blue)" }
func (FateForeseenBlue) Cost(*card.TurnState) int                { return 0 }
func (FateForeseenBlue) Pitch() int               { return 3 }
func (FateForeseenBlue) Attack() int              { return 0 }
func (FateForeseenBlue) Defense() int             { return 2 }
func (FateForeseenBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (FateForeseenBlue) GoAgain() bool            { return false }
func (FateForeseenBlue) NotSilverAgeLegal()       {}
// not implemented: Opt 1 rider; block value is printed defence only
func (FateForeseenBlue) NotImplemented()           {}
func (FateForeseenBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
