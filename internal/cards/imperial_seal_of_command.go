// Imperial Seal of Command — Generic Action - Item. Cost 0. Printed pitch variants: Red 1.
//
// Text: "**Legendary** **Action** - Destroy this: Defense reaction cards can't be played this turn.
// If you are Royal, the next time you hit a hero this turn, destroy all cards in their arsenal.
// **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var imperialSealOfCommandTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ImperialSealOfCommandRed struct{}

func (ImperialSealOfCommandRed) ID() ids.CardID          { return ids.ImperialSealOfCommandRed }
func (ImperialSealOfCommandRed) Name() string            { return "Imperial Seal of Command" }
func (ImperialSealOfCommandRed) Cost(*sim.TurnState) int { return 0 }
func (ImperialSealOfCommandRed) Pitch() int              { return 1 }
func (ImperialSealOfCommandRed) Attack() int             { return 0 }
func (ImperialSealOfCommandRed) Defense() int            { return 0 }
func (ImperialSealOfCommandRed) Types() card.TypeSet     { return imperialSealOfCommandTypes }
func (ImperialSealOfCommandRed) GoAgain() bool           { return false }

// not implemented: activated 'no DR this turn' + Royal-only arsenal-wipe on hit
func (ImperialSealOfCommandRed) NotImplemented()                            {}
func (ImperialSealOfCommandRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
