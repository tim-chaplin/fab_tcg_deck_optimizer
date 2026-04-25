// Pummel — Generic Attack Reaction. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Choose 1; - Target club or hammer weapon attack gains +4{p}. - Target attack action card
// with cost 2 or more gets +4{p} and "When this hits a hero, they discard a card.""

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pummelTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type PummelRed struct{}

func (PummelRed) ID() card.ID                               { return card.PummelRed }
func (PummelRed) Name() string                              { return "Pummel (Red)" }
func (PummelRed) Cost(*card.TurnState) int                  { return 2 }
func (PummelRed) Pitch() int                                { return 1 }
func (PummelRed) Attack() int                               { return 0 }
func (PummelRed) Defense() int                              { return 2 }
func (PummelRed) Types() card.TypeSet                       { return pummelTypes }
func (PummelRed) GoAgain() bool                             { return false }
// not implemented: modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action (on-hit discard)
func (PummelRed) NotImplemented()                           {}
func (PummelRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type PummelYellow struct{}

func (PummelYellow) ID() card.ID                               { return card.PummelYellow }
func (PummelYellow) Name() string                              { return "Pummel (Yellow)" }
func (PummelYellow) Cost(*card.TurnState) int                  { return 2 }
func (PummelYellow) Pitch() int                                { return 2 }
func (PummelYellow) Attack() int                               { return 0 }
func (PummelYellow) Defense() int                              { return 2 }
func (PummelYellow) Types() card.TypeSet                       { return pummelTypes }
func (PummelYellow) GoAgain() bool                             { return false }
// not implemented: modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action (on-hit discard)
func (PummelYellow) NotImplemented()                           {}
func (PummelYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type PummelBlue struct{}

func (PummelBlue) ID() card.ID                               { return card.PummelBlue }
func (PummelBlue) Name() string                              { return "Pummel (Blue)" }
func (PummelBlue) Cost(*card.TurnState) int                  { return 2 }
func (PummelBlue) Pitch() int                                { return 3 }
func (PummelBlue) Attack() int                               { return 0 }
func (PummelBlue) Defense() int                              { return 2 }
func (PummelBlue) Types() card.TypeSet                       { return pummelTypes }
func (PummelBlue) GoAgain() bool                             { return false }
// not implemented: modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action (on-hit discard)
func (PummelBlue) NotImplemented()                           {}
func (PummelBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
