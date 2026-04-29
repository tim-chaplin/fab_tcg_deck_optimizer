// Aether Slash — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Aether Slash attacks, if a 'non-attack' action card was pitched to play it, deal 1
// arcane damage to any target."
//
// Reads self.PitchedToPlay (via the PitchAttributionReader marker) to gate the +1 arcane
// rider on a non-attack action funding THIS copy specifically.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var aetherSlashTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type AetherSlashRed struct{}

func (AetherSlashRed) ID() ids.CardID          { return ids.AetherSlashRed }
func (AetherSlashRed) Name() string            { return "Aether Slash" }
func (AetherSlashRed) Cost(*sim.TurnState) int { return 1 }
func (AetherSlashRed) Pitch() int              { return 1 }
func (AetherSlashRed) Attack() int             { return 4 }
func (AetherSlashRed) Defense() int            { return 3 }
func (AetherSlashRed) Types() card.TypeSet     { return aetherSlashTypes }
func (AetherSlashRed) GoAgain() bool           { return false }
func (AetherSlashRed) ReadsPitchAttribution()  {}
func (AetherSlashRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	aetherSlashApplyRider(s, self)
}

type AetherSlashYellow struct{}

func (AetherSlashYellow) ID() ids.CardID          { return ids.AetherSlashYellow }
func (AetherSlashYellow) Name() string            { return "Aether Slash" }
func (AetherSlashYellow) Cost(*sim.TurnState) int { return 1 }
func (AetherSlashYellow) Pitch() int              { return 2 }
func (AetherSlashYellow) Attack() int             { return 3 }
func (AetherSlashYellow) Defense() int            { return 3 }
func (AetherSlashYellow) Types() card.TypeSet     { return aetherSlashTypes }
func (AetherSlashYellow) GoAgain() bool           { return false }
func (AetherSlashYellow) ReadsPitchAttribution()  {}
func (AetherSlashYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	aetherSlashApplyRider(s, self)
}

type AetherSlashBlue struct{}

func (AetherSlashBlue) ID() ids.CardID          { return ids.AetherSlashBlue }
func (AetherSlashBlue) Name() string            { return "Aether Slash" }
func (AetherSlashBlue) Cost(*sim.TurnState) int { return 1 }
func (AetherSlashBlue) Pitch() int              { return 3 }
func (AetherSlashBlue) Attack() int             { return 2 }
func (AetherSlashBlue) Defense() int            { return 3 }
func (AetherSlashBlue) Types() card.TypeSet     { return aetherSlashTypes }
func (AetherSlashBlue) GoAgain() bool           { return false }
func (AetherSlashBlue) ReadsPitchAttribution()  {}
func (AetherSlashBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	aetherSlashApplyRider(s, self)
}

// aetherSlashApplyRider deals 1 arcane and emits the rider sub-line when a non-attack action
// is among the pitched cards the runner attributed to paying for this Aether Slash.
func aetherSlashApplyRider(s *sim.TurnState, self *sim.CardState) {
	for _, p := range self.PitchedToPlay {
		if p.Types().IsNonAttackAction() {
			s.DealAndLogArcaneDamage(self, 1)
			return
		}
	}
}
