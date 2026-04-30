// Scour the Battlescape — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card. If
// Scour the Battlescape is played from arsenal, it gains **go again**."
//
// Modelling: hand-cycle isn't modelled. Standard played-from-arsenal go-again
// (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var scourTheBattlescapeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// scourTheBattlescapePlay grants self Go again when this copy was played from arsenal,
// then emits the chain step.
func scourTheBattlescapePlay(s *sim.TurnState, self *sim.CardState) {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type ScourTheBattlescapeRed struct{}

func (ScourTheBattlescapeRed) ID() ids.CardID          { return ids.ScourTheBattlescapeRed }
func (ScourTheBattlescapeRed) Name() string            { return "Scour the Battlescape" }
func (ScourTheBattlescapeRed) Cost(*sim.TurnState) int { return 0 }
func (ScourTheBattlescapeRed) Pitch() int              { return 1 }
func (ScourTheBattlescapeRed) Attack() int             { return 3 }
func (ScourTheBattlescapeRed) Defense() int            { return 2 }
func (ScourTheBattlescapeRed) Types() card.TypeSet     { return scourTheBattlescapeTypes }
func (ScourTheBattlescapeRed) GoAgain() bool           { return false }

// not implemented: hand-cycle rider (put a card on bottom of deck, draw)
func (ScourTheBattlescapeRed) NotImplemented() {}
func (ScourTheBattlescapeRed) Play(s *sim.TurnState, self *sim.CardState) {
	scourTheBattlescapePlay(s, self)
}

type ScourTheBattlescapeYellow struct{}

func (ScourTheBattlescapeYellow) ID() ids.CardID          { return ids.ScourTheBattlescapeYellow }
func (ScourTheBattlescapeYellow) Name() string            { return "Scour the Battlescape" }
func (ScourTheBattlescapeYellow) Cost(*sim.TurnState) int { return 0 }
func (ScourTheBattlescapeYellow) Pitch() int              { return 2 }
func (ScourTheBattlescapeYellow) Attack() int             { return 2 }
func (ScourTheBattlescapeYellow) Defense() int            { return 2 }
func (ScourTheBattlescapeYellow) Types() card.TypeSet     { return scourTheBattlescapeTypes }
func (ScourTheBattlescapeYellow) GoAgain() bool           { return false }

// not implemented: hand-cycle rider (put a card on bottom of deck, draw)
func (ScourTheBattlescapeYellow) NotImplemented() {}
func (ScourTheBattlescapeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	scourTheBattlescapePlay(s, self)
}

type ScourTheBattlescapeBlue struct{}

func (ScourTheBattlescapeBlue) ID() ids.CardID          { return ids.ScourTheBattlescapeBlue }
func (ScourTheBattlescapeBlue) Name() string            { return "Scour the Battlescape" }
func (ScourTheBattlescapeBlue) Cost(*sim.TurnState) int { return 0 }
func (ScourTheBattlescapeBlue) Pitch() int              { return 3 }
func (ScourTheBattlescapeBlue) Attack() int             { return 1 }
func (ScourTheBattlescapeBlue) Defense() int            { return 2 }
func (ScourTheBattlescapeBlue) Types() card.TypeSet     { return scourTheBattlescapeTypes }
func (ScourTheBattlescapeBlue) GoAgain() bool           { return false }

// not implemented: hand-cycle rider (put a card on bottom of deck, draw)
func (ScourTheBattlescapeBlue) NotImplemented() {}
func (ScourTheBattlescapeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	scourTheBattlescapePlay(s, self)
}

func (ScourTheBattlescapeRed) ConditionalGoAgain()    {}
func (ScourTheBattlescapeYellow) ConditionalGoAgain() {}
func (ScourTheBattlescapeBlue) ConditionalGoAgain()   {}
