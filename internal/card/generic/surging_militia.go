// Surging Militia — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Surging Militia has +1{p} for each non-equipment card defending it."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var surgingMilitiaTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SurgingMilitiaRed struct{}

func (SurgingMilitiaRed) ID() card.ID                 { return card.SurgingMilitiaRed }
func (SurgingMilitiaRed) Name() string                { return "Surging Militia (Red)" }
func (SurgingMilitiaRed) Cost(*card.TurnState) int                   { return 2 }
func (SurgingMilitiaRed) Pitch() int                  { return 1 }
func (SurgingMilitiaRed) Attack() int                 { return 5 }
func (SurgingMilitiaRed) Defense() int                { return 2 }
func (SurgingMilitiaRed) Types() card.TypeSet         { return surgingMilitiaTypes }
func (SurgingMilitiaRed) GoAgain() bool               { return false }
// not implemented: defended-by +N{p} rider (defender's hand contents not exposed)
func (SurgingMilitiaRed) NotImplemented()             {}
func (c SurgingMilitiaRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type SurgingMilitiaYellow struct{}

func (SurgingMilitiaYellow) ID() card.ID                 { return card.SurgingMilitiaYellow }
func (SurgingMilitiaYellow) Name() string                { return "Surging Militia (Yellow)" }
func (SurgingMilitiaYellow) Cost(*card.TurnState) int                   { return 2 }
func (SurgingMilitiaYellow) Pitch() int                  { return 2 }
func (SurgingMilitiaYellow) Attack() int                 { return 4 }
func (SurgingMilitiaYellow) Defense() int                { return 2 }
func (SurgingMilitiaYellow) Types() card.TypeSet         { return surgingMilitiaTypes }
func (SurgingMilitiaYellow) GoAgain() bool               { return false }
// not implemented: defended-by +N{p} rider (defender's hand contents not exposed)
func (SurgingMilitiaYellow) NotImplemented()             {}
func (c SurgingMilitiaYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type SurgingMilitiaBlue struct{}

func (SurgingMilitiaBlue) ID() card.ID                 { return card.SurgingMilitiaBlue }
func (SurgingMilitiaBlue) Name() string                { return "Surging Militia (Blue)" }
func (SurgingMilitiaBlue) Cost(*card.TurnState) int                   { return 2 }
func (SurgingMilitiaBlue) Pitch() int                  { return 3 }
func (SurgingMilitiaBlue) Attack() int                 { return 3 }
func (SurgingMilitiaBlue) Defense() int                { return 2 }
func (SurgingMilitiaBlue) Types() card.TypeSet         { return surgingMilitiaTypes }
func (SurgingMilitiaBlue) GoAgain() bool               { return false }
// not implemented: defended-by +N{p} rider (defender's hand contents not exposed)
func (SurgingMilitiaBlue) NotImplemented()             {}
func (c SurgingMilitiaBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
