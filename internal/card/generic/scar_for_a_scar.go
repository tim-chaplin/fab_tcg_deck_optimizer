// Scar for a Scar — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**."

package generic

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

var scarForAScarTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type ScarForAScarRed struct{}

func (ScarForAScarRed) ID() card.ID                 { return card.ScarForAScarRed }
func (ScarForAScarRed) Name() string                { return "Scar for a Scar" }
func (ScarForAScarRed) Cost(*card.TurnState) int                   { return 0 }
func (ScarForAScarRed) Pitch() int                  { return 1 }
func (ScarForAScarRed) Attack() int                 { return 4 }
func (ScarForAScarRed) Defense() int                { return 2 }
func (ScarForAScarRed) Types() card.TypeSet         { return scarForAScarTypes }
func (ScarForAScarRed) GoAgain() bool               { return simstate.HeroWantsLowerHealth() }
func (c ScarForAScarRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type ScarForAScarYellow struct{}

func (ScarForAScarYellow) ID() card.ID                 { return card.ScarForAScarYellow }
func (ScarForAScarYellow) Name() string                { return "Scar for a Scar" }
func (ScarForAScarYellow) Cost(*card.TurnState) int                   { return 0 }
func (ScarForAScarYellow) Pitch() int                  { return 2 }
func (ScarForAScarYellow) Attack() int                 { return 3 }
func (ScarForAScarYellow) Defense() int                { return 2 }
func (ScarForAScarYellow) Types() card.TypeSet         { return scarForAScarTypes }
func (ScarForAScarYellow) GoAgain() bool               { return simstate.HeroWantsLowerHealth() }
func (c ScarForAScarYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type ScarForAScarBlue struct{}

func (ScarForAScarBlue) ID() card.ID                 { return card.ScarForAScarBlue }
func (ScarForAScarBlue) Name() string                { return "Scar for a Scar" }
func (ScarForAScarBlue) Cost(*card.TurnState) int                   { return 0 }
func (ScarForAScarBlue) Pitch() int                  { return 3 }
func (ScarForAScarBlue) Attack() int                 { return 2 }
func (ScarForAScarBlue) Defense() int                { return 2 }
func (ScarForAScarBlue) Types() card.TypeSet         { return scarForAScarTypes }
func (ScarForAScarBlue) GoAgain() bool               { return simstate.HeroWantsLowerHealth() }
func (c ScarForAScarBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
