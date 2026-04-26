// Relentless Pursuit — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "**Mark** target opposing hero. If you've attacked them this turn, put this on the bottom
// of its owner's deck. **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var relentlessPursuitTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RelentlessPursuitBlue struct{}

func (RelentlessPursuitBlue) ID() card.ID                 { return card.RelentlessPursuitBlue }
func (RelentlessPursuitBlue) Name() string                { return "Relentless Pursuit" }
func (RelentlessPursuitBlue) Cost(*card.TurnState) int                   { return 0 }
func (RelentlessPursuitBlue) Pitch() int                  { return 3 }
func (RelentlessPursuitBlue) Attack() int                 { return 0 }
func (RelentlessPursuitBlue) Defense() int                { return 3 }
func (RelentlessPursuitBlue) Types() card.TypeSet         { return relentlessPursuitTypes }
func (RelentlessPursuitBlue) GoAgain() bool               { return true }
// not implemented: marked-target gate + 'attacked them this turn' chain rider
func (RelentlessPursuitBlue) NotImplemented()             {}
func (RelentlessPursuitBlue) Play(s *card.TurnState, _ *card.CardState) int { return 0 }
