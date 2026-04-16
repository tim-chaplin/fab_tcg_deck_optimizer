// Rise Above — Generic Defense Reaction. Cost 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on top of your deck rather than pay Rise Above's {r}
// cost." Simplification: the alternative hand-as-cost option isn't modelled — we pay the printed
// 2{r} or the card isn't legal (the partition's canAfford check handles that).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type RiseAboveRed struct{}

func (RiseAboveRed) ID() card.ID                 { return card.RiseAboveRed }
func (RiseAboveRed) Name() string             { return "Rise Above (Red)" }
func (RiseAboveRed) Cost() int                { return 2 }
func (RiseAboveRed) Pitch() int               { return 1 }
func (RiseAboveRed) Attack() int              { return 0 }
func (RiseAboveRed) Defense() int             { return 4 }
func (RiseAboveRed) Types() card.TypeSet      { return defenseReactionTypes }
func (RiseAboveRed) GoAgain() bool            { return false }
func (RiseAboveRed) Play(*card.TurnState) int { return 0 }

type RiseAboveYellow struct{}

func (RiseAboveYellow) ID() card.ID                 { return card.RiseAboveYellow }
func (RiseAboveYellow) Name() string             { return "Rise Above (Yellow)" }
func (RiseAboveYellow) Cost() int                { return 2 }
func (RiseAboveYellow) Pitch() int               { return 2 }
func (RiseAboveYellow) Attack() int              { return 0 }
func (RiseAboveYellow) Defense() int             { return 3 }
func (RiseAboveYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (RiseAboveYellow) GoAgain() bool            { return false }
func (RiseAboveYellow) Play(*card.TurnState) int { return 0 }

type RiseAboveBlue struct{}

func (RiseAboveBlue) ID() card.ID                 { return card.RiseAboveBlue }
func (RiseAboveBlue) Name() string             { return "Rise Above (Blue)" }
func (RiseAboveBlue) Cost() int                { return 2 }
func (RiseAboveBlue) Pitch() int               { return 3 }
func (RiseAboveBlue) Attack() int              { return 0 }
func (RiseAboveBlue) Defense() int             { return 2 }
func (RiseAboveBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (RiseAboveBlue) GoAgain() bool            { return false }
func (RiseAboveBlue) Play(*card.TurnState) int { return 0 }
