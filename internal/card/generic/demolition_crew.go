// Demolition Crew — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Demolition Crew, reveal a card in your hand with cost 2 or
// greater. **Dominate**"
//
// Simplification: Additional reveal cost isn't modelled; Dominate keyword is held but unused.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var demolitionCrewTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DemolitionCrewRed struct{}

func (DemolitionCrewRed) ID() card.ID                 { return card.DemolitionCrewRed }
func (DemolitionCrewRed) Name() string                { return "Demolition Crew (Red)" }
func (DemolitionCrewRed) Cost(*card.TurnState) int                   { return 2 }
func (DemolitionCrewRed) Pitch() int                  { return 1 }
func (DemolitionCrewRed) Attack() int                 { return 6 }
func (DemolitionCrewRed) Defense() int                { return 2 }
func (DemolitionCrewRed) Types() card.TypeSet         { return demolitionCrewTypes }
func (DemolitionCrewRed) GoAgain() bool               { return false }
func (c DemolitionCrewRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type DemolitionCrewYellow struct{}

func (DemolitionCrewYellow) ID() card.ID                 { return card.DemolitionCrewYellow }
func (DemolitionCrewYellow) Name() string                { return "Demolition Crew (Yellow)" }
func (DemolitionCrewYellow) Cost(*card.TurnState) int                   { return 2 }
func (DemolitionCrewYellow) Pitch() int                  { return 2 }
func (DemolitionCrewYellow) Attack() int                 { return 5 }
func (DemolitionCrewYellow) Defense() int                { return 2 }
func (DemolitionCrewYellow) Types() card.TypeSet         { return demolitionCrewTypes }
func (DemolitionCrewYellow) GoAgain() bool               { return false }
func (c DemolitionCrewYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type DemolitionCrewBlue struct{}

func (DemolitionCrewBlue) ID() card.ID                 { return card.DemolitionCrewBlue }
func (DemolitionCrewBlue) Name() string                { return "Demolition Crew (Blue)" }
func (DemolitionCrewBlue) Cost(*card.TurnState) int                   { return 2 }
func (DemolitionCrewBlue) Pitch() int                  { return 3 }
func (DemolitionCrewBlue) Attack() int                 { return 4 }
func (DemolitionCrewBlue) Defense() int                { return 2 }
func (DemolitionCrewBlue) Types() card.TypeSet         { return demolitionCrewTypes }
func (DemolitionCrewBlue) GoAgain() bool               { return false }
func (c DemolitionCrewBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
