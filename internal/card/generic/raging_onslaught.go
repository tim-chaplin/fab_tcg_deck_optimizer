// Raging Onslaught — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var ragingOnslaughtTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RagingOnslaughtRed struct{}

func (RagingOnslaughtRed) ID() card.ID                 { return card.RagingOnslaughtRed }
func (RagingOnslaughtRed) Name() string                { return "Raging Onslaught (Red)" }
func (RagingOnslaughtRed) Cost(*card.TurnState) int                   { return 3 }
func (RagingOnslaughtRed) Pitch() int                  { return 1 }
func (RagingOnslaughtRed) Attack() int                 { return 7 }
func (RagingOnslaughtRed) Defense() int                { return 3 }
func (RagingOnslaughtRed) Types() card.TypeSet         { return ragingOnslaughtTypes }
func (RagingOnslaughtRed) GoAgain() bool               { return false }
func (c RagingOnslaughtRed) Play(s *card.TurnState) int { return c.Attack() }

type RagingOnslaughtYellow struct{}

func (RagingOnslaughtYellow) ID() card.ID                 { return card.RagingOnslaughtYellow }
func (RagingOnslaughtYellow) Name() string                { return "Raging Onslaught (Yellow)" }
func (RagingOnslaughtYellow) Cost(*card.TurnState) int                   { return 3 }
func (RagingOnslaughtYellow) Pitch() int                  { return 2 }
func (RagingOnslaughtYellow) Attack() int                 { return 6 }
func (RagingOnslaughtYellow) Defense() int                { return 3 }
func (RagingOnslaughtYellow) Types() card.TypeSet         { return ragingOnslaughtTypes }
func (RagingOnslaughtYellow) GoAgain() bool               { return false }
func (c RagingOnslaughtYellow) Play(s *card.TurnState) int { return c.Attack() }

type RagingOnslaughtBlue struct{}

func (RagingOnslaughtBlue) ID() card.ID                 { return card.RagingOnslaughtBlue }
func (RagingOnslaughtBlue) Name() string                { return "Raging Onslaught (Blue)" }
func (RagingOnslaughtBlue) Cost(*card.TurnState) int                   { return 3 }
func (RagingOnslaughtBlue) Pitch() int                  { return 3 }
func (RagingOnslaughtBlue) Attack() int                 { return 5 }
func (RagingOnslaughtBlue) Defense() int                { return 3 }
func (RagingOnslaughtBlue) Types() card.TypeSet         { return ragingOnslaughtTypes }
func (RagingOnslaughtBlue) GoAgain() bool               { return false }
func (c RagingOnslaughtBlue) Play(s *card.TurnState) int { return c.Attack() }
