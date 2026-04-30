// Demolition Crew — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Demolition Crew, reveal a card in your hand with cost 2 or
// greater. **Dominate**"
//
// The reveal is non-consuming — the revealed card stays in hand. With no eligible reveal the
// additional cost can't be paid and the card resolves as a no-op so the optimizer can compare
// cleanly against the Held alternative.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var demolitionCrewTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// demolitionCrewPlay applies the additional-cost gate before crediting the attack. The
// chain runner removes self from s.Hand at iteration start, so the scan can't pick the
// Demolition Crew copy resolving right now as its own reveal target.
func demolitionCrewPlay(s *sim.TurnState, self *sim.CardState) {
	for _, c := range s.Hand {
		if c.Cost(s) >= 2 {
			s.ApplyAndLogEffectiveAttack(self)
			return
		}
	}
	s.LogPlay(self)
}

type DemolitionCrewRed struct{}

func (DemolitionCrewRed) ID() ids.CardID          { return ids.DemolitionCrewRed }
func (DemolitionCrewRed) Name() string            { return "Demolition Crew" }
func (DemolitionCrewRed) Cost(*sim.TurnState) int { return 2 }
func (DemolitionCrewRed) Pitch() int              { return 1 }
func (DemolitionCrewRed) Attack() int             { return 6 }
func (DemolitionCrewRed) Defense() int            { return 2 }
func (DemolitionCrewRed) Types() card.TypeSet     { return demolitionCrewTypes }
func (DemolitionCrewRed) GoAgain() bool           { return false }
func (DemolitionCrewRed) Dominate()               {}
func (c DemolitionCrewRed) Play(s *sim.TurnState, self *sim.CardState) {
	demolitionCrewPlay(s, self)
}

type DemolitionCrewYellow struct{}

func (DemolitionCrewYellow) ID() ids.CardID          { return ids.DemolitionCrewYellow }
func (DemolitionCrewYellow) Name() string            { return "Demolition Crew" }
func (DemolitionCrewYellow) Cost(*sim.TurnState) int { return 2 }
func (DemolitionCrewYellow) Pitch() int              { return 2 }
func (DemolitionCrewYellow) Attack() int             { return 5 }
func (DemolitionCrewYellow) Defense() int            { return 2 }
func (DemolitionCrewYellow) Types() card.TypeSet     { return demolitionCrewTypes }
func (DemolitionCrewYellow) GoAgain() bool           { return false }
func (DemolitionCrewYellow) Dominate()               {}
func (c DemolitionCrewYellow) Play(s *sim.TurnState, self *sim.CardState) {
	demolitionCrewPlay(s, self)
}

type DemolitionCrewBlue struct{}

func (DemolitionCrewBlue) ID() ids.CardID          { return ids.DemolitionCrewBlue }
func (DemolitionCrewBlue) Name() string            { return "Demolition Crew" }
func (DemolitionCrewBlue) Cost(*sim.TurnState) int { return 2 }
func (DemolitionCrewBlue) Pitch() int              { return 3 }
func (DemolitionCrewBlue) Attack() int             { return 4 }
func (DemolitionCrewBlue) Defense() int            { return 2 }
func (DemolitionCrewBlue) Types() card.TypeSet     { return demolitionCrewTypes }
func (DemolitionCrewBlue) GoAgain() bool           { return false }
func (DemolitionCrewBlue) Dominate()               {}
func (c DemolitionCrewBlue) Play(s *sim.TurnState, self *sim.CardState) {
	demolitionCrewPlay(s, self)
}
