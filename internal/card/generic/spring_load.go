// Spring Load — Generic Action - Attack. Cost 1. Printed power: Red 2, Yellow 2, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, if you have no cards in hand, it gains +3{p}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var springLoadTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SpringLoadRed struct{}

func (SpringLoadRed) ID() card.ID                 { return card.SpringLoadRed }
func (SpringLoadRed) Name() string                { return "Spring Load (Red)" }
func (SpringLoadRed) Cost(*card.TurnState) int                   { return 1 }
func (SpringLoadRed) Pitch() int                  { return 1 }
func (SpringLoadRed) Attack() int                 { return 2 }
func (SpringLoadRed) Defense() int                { return 2 }
func (SpringLoadRed) Types() card.TypeSet         { return springLoadTypes }
func (SpringLoadRed) GoAgain() bool               { return false }
// not implemented: +3{p} 'no cards in hand' rider never fires
func (SpringLoadRed) NotImplemented()             {}
func (c SpringLoadRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type SpringLoadYellow struct{}

func (SpringLoadYellow) ID() card.ID                 { return card.SpringLoadYellow }
func (SpringLoadYellow) Name() string                { return "Spring Load (Yellow)" }
func (SpringLoadYellow) Cost(*card.TurnState) int                   { return 1 }
func (SpringLoadYellow) Pitch() int                  { return 2 }
func (SpringLoadYellow) Attack() int                 { return 2 }
func (SpringLoadYellow) Defense() int                { return 2 }
func (SpringLoadYellow) Types() card.TypeSet         { return springLoadTypes }
func (SpringLoadYellow) GoAgain() bool               { return false }
// not implemented: +3{p} 'no cards in hand' rider never fires
func (SpringLoadYellow) NotImplemented()             {}
func (c SpringLoadYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type SpringLoadBlue struct{}

func (SpringLoadBlue) ID() card.ID                 { return card.SpringLoadBlue }
func (SpringLoadBlue) Name() string                { return "Spring Load (Blue)" }
func (SpringLoadBlue) Cost(*card.TurnState) int                   { return 1 }
func (SpringLoadBlue) Pitch() int                  { return 3 }
func (SpringLoadBlue) Attack() int                 { return 2 }
func (SpringLoadBlue) Defense() int                { return 2 }
func (SpringLoadBlue) Types() card.TypeSet         { return springLoadTypes }
func (SpringLoadBlue) GoAgain() bool               { return false }
// not implemented: +3{p} 'no cards in hand' rider never fires
func (SpringLoadBlue) NotImplemented()             {}
func (c SpringLoadBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
