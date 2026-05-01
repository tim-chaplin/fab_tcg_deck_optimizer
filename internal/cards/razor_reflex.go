// Razor Reflex — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Choose 1; - Target dagger or sword weapon attack gets +3{p}. - Target attack action card
// with cost 1 or less gets +3{p} and "When this hits, it gets **go again**.""

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var razorReflexTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type RazorReflexRed struct{}

func (RazorReflexRed) ID() ids.CardID          { return ids.RazorReflexRed }
func (RazorReflexRed) Name() string            { return "Razor Reflex" }
func (RazorReflexRed) Cost(*sim.TurnState) int { return 1 }
func (RazorReflexRed) Pitch() int              { return 1 }
func (RazorReflexRed) Attack() int             { return 0 }
func (RazorReflexRed) Defense() int            { return 2 }
func (RazorReflexRed) Types() card.TypeSet     { return razorReflexTypes }
func (RazorReflexRed) GoAgain() bool           { return false }

// not implemented: modal AR +N{p}: dagger/sword weapon attack OR cost ≤1 attack action
// (on-hit go again)
func (RazorReflexRed) NotImplemented()                            {}
func (RazorReflexRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type RazorReflexYellow struct{}

func (RazorReflexYellow) ID() ids.CardID          { return ids.RazorReflexYellow }
func (RazorReflexYellow) Name() string            { return "Razor Reflex" }
func (RazorReflexYellow) Cost(*sim.TurnState) int { return 1 }
func (RazorReflexYellow) Pitch() int              { return 2 }
func (RazorReflexYellow) Attack() int             { return 0 }
func (RazorReflexYellow) Defense() int            { return 2 }
func (RazorReflexYellow) Types() card.TypeSet     { return razorReflexTypes }
func (RazorReflexYellow) GoAgain() bool           { return false }

// not implemented: modal AR +N{p}: dagger/sword weapon attack OR cost ≤1 attack action
// (on-hit go again)
func (RazorReflexYellow) NotImplemented()                            {}
func (RazorReflexYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type RazorReflexBlue struct{}

func (RazorReflexBlue) ID() ids.CardID          { return ids.RazorReflexBlue }
func (RazorReflexBlue) Name() string            { return "Razor Reflex" }
func (RazorReflexBlue) Cost(*sim.TurnState) int { return 1 }
func (RazorReflexBlue) Pitch() int              { return 3 }
func (RazorReflexBlue) Attack() int             { return 0 }
func (RazorReflexBlue) Defense() int            { return 2 }
func (RazorReflexBlue) Types() card.TypeSet     { return razorReflexTypes }
func (RazorReflexBlue) GoAgain() bool           { return false }

// not implemented: modal AR +N{p}: dagger/sword weapon attack OR cost ≤1 attack action
// (on-hit go again)
func (RazorReflexBlue) NotImplemented()                            {}
func (RazorReflexBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
