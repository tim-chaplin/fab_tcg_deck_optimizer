// Pick a Card, Any Card — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at target opponent's hand then name a card. Choose a random card from their hand and
// reveal it. If it's the named card, create a Silver token. Repeat this process thrice. **Go
// again**"
//
// Simplification: Opponent hand inspection and Silver-token economy aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pickACardAnyCardTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type PickACardAnyCardRed struct{}

func (PickACardAnyCardRed) ID() card.ID                 { return card.PickACardAnyCardRed }
func (PickACardAnyCardRed) Name() string                { return "Pick a Card, Any Card (Red)" }
func (PickACardAnyCardRed) Cost(*card.TurnState) int                   { return 0 }
func (PickACardAnyCardRed) Pitch() int                  { return 1 }
func (PickACardAnyCardRed) Attack() int                 { return 0 }
func (PickACardAnyCardRed) Defense() int                { return 2 }
func (PickACardAnyCardRed) Types() card.TypeSet         { return pickACardAnyCardTypes }
func (PickACardAnyCardRed) GoAgain() bool               { return true }
func (PickACardAnyCardRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return 0 }

type PickACardAnyCardYellow struct{}

func (PickACardAnyCardYellow) ID() card.ID                 { return card.PickACardAnyCardYellow }
func (PickACardAnyCardYellow) Name() string                { return "Pick a Card, Any Card (Yellow)" }
func (PickACardAnyCardYellow) Cost(*card.TurnState) int                   { return 0 }
func (PickACardAnyCardYellow) Pitch() int                  { return 2 }
func (PickACardAnyCardYellow) Attack() int                 { return 0 }
func (PickACardAnyCardYellow) Defense() int                { return 2 }
func (PickACardAnyCardYellow) Types() card.TypeSet         { return pickACardAnyCardTypes }
func (PickACardAnyCardYellow) GoAgain() bool               { return true }
func (PickACardAnyCardYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return 0 }

type PickACardAnyCardBlue struct{}

func (PickACardAnyCardBlue) ID() card.ID                 { return card.PickACardAnyCardBlue }
func (PickACardAnyCardBlue) Name() string                { return "Pick a Card, Any Card (Blue)" }
func (PickACardAnyCardBlue) Cost(*card.TurnState) int                   { return 0 }
func (PickACardAnyCardBlue) Pitch() int                  { return 3 }
func (PickACardAnyCardBlue) Attack() int                 { return 0 }
func (PickACardAnyCardBlue) Defense() int                { return 2 }
func (PickACardAnyCardBlue) Types() card.TypeSet         { return pickACardAnyCardTypes }
func (PickACardAnyCardBlue) GoAgain() bool               { return true }
func (PickACardAnyCardBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return 0 }
