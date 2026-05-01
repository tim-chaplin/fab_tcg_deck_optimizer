// Drowning Dire — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you have played or created an aura this turn, Drowning Dire gains **dominate**.
//
// When Drowning Dire hits, you may put a 'non-attack' action card from your graveyard on the
// bottom of your deck."
//
// Modelling: the Dominate grant is conditional on s.HasPlayedOrCreatedAura(). Standard
// self.GrantedDominate wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var drowningDireTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// drowningDirePlay grants self Dominate when an aura has been played or created this turn,
// then emits the chain step.
func drowningDirePlay(s *sim.TurnState, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.GrantedDominate = true
	}
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type DrowningDireRed struct{}

func (DrowningDireRed) ID() ids.CardID          { return ids.DrowningDireRed }
func (DrowningDireRed) Name() string            { return "Drowning Dire" }
func (DrowningDireRed) Cost(*sim.TurnState) int { return 2 }
func (DrowningDireRed) Pitch() int              { return 1 }
func (DrowningDireRed) Attack() int             { return 5 }
func (DrowningDireRed) Defense() int            { return 3 }
func (DrowningDireRed) Types() card.TypeSet     { return drowningDireTypes }
func (DrowningDireRed) GoAgain() bool           { return false }

// not implemented: on-hit "may put a non-attack action card from your graveyard on the bottom
// of your deck" rider
func (DrowningDireRed) NotImplemented() {}
func (DrowningDireRed) Play(s *sim.TurnState, self *sim.CardState) {
	drowningDirePlay(s, self)
}

type DrowningDireYellow struct{}

func (DrowningDireYellow) ID() ids.CardID          { return ids.DrowningDireYellow }
func (DrowningDireYellow) Name() string            { return "Drowning Dire" }
func (DrowningDireYellow) Cost(*sim.TurnState) int { return 2 }
func (DrowningDireYellow) Pitch() int              { return 2 }
func (DrowningDireYellow) Attack() int             { return 4 }
func (DrowningDireYellow) Defense() int            { return 3 }
func (DrowningDireYellow) Types() card.TypeSet     { return drowningDireTypes }
func (DrowningDireYellow) GoAgain() bool           { return false }

// not implemented: on-hit "may put a non-attack action card from your graveyard on the bottom
// of your deck" rider
func (DrowningDireYellow) NotImplemented() {}
func (DrowningDireYellow) Play(s *sim.TurnState, self *sim.CardState) {
	drowningDirePlay(s, self)
}

type DrowningDireBlue struct{}

func (DrowningDireBlue) ID() ids.CardID          { return ids.DrowningDireBlue }
func (DrowningDireBlue) Name() string            { return "Drowning Dire" }
func (DrowningDireBlue) Cost(*sim.TurnState) int { return 2 }
func (DrowningDireBlue) Pitch() int              { return 3 }
func (DrowningDireBlue) Attack() int             { return 3 }
func (DrowningDireBlue) Defense() int            { return 3 }
func (DrowningDireBlue) Types() card.TypeSet     { return drowningDireTypes }
func (DrowningDireBlue) GoAgain() bool           { return false }

// not implemented: on-hit "may put a non-attack action card from your graveyard on the bottom
// of your deck" rider
func (DrowningDireBlue) NotImplemented() {}
func (DrowningDireBlue) Play(s *sim.TurnState, self *sim.CardState) {
	drowningDirePlay(s, self)
}
