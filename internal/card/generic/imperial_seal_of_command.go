// Imperial Seal of Command — Generic Action - Item. Cost 0. Printed pitch variants: Red 1.
//
// Text: "**Legendary** **Action** - Destroy this: Defense reaction cards can't be played this turn.
// If you are Royal, the next time you hit a hero this turn, destroy all cards in their arsenal.
// **Go again**"
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var imperialSealOfCommandTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ImperialSealOfCommandRed struct{}

func (ImperialSealOfCommandRed) ID() card.ID                               { return card.ImperialSealOfCommandRed }
func (ImperialSealOfCommandRed) Name() string                              { return "Imperial Seal of Command (Red)" }
func (ImperialSealOfCommandRed) Cost(*card.TurnState) int                  { return 0 }
func (ImperialSealOfCommandRed) Pitch() int                                { return 1 }
func (ImperialSealOfCommandRed) Attack() int                               { return 0 }
func (ImperialSealOfCommandRed) Defense() int                              { return 0 }
func (ImperialSealOfCommandRed) Types() card.TypeSet                       { return imperialSealOfCommandTypes }
func (ImperialSealOfCommandRed) GoAgain() bool                             { return false }
func (ImperialSealOfCommandRed) NotImplemented()                           {}
func (ImperialSealOfCommandRed) Play(*card.TurnState, *card.CardState) int { return 0 }
