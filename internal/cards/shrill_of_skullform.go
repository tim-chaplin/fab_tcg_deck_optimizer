// Shrill of Skullform — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Shrill of Skullform gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var shrillTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type ShrillOfSkullformRed struct{}

func (ShrillOfSkullformRed) ID() ids.CardID          { return ids.ShrillOfSkullformRed }
func (ShrillOfSkullformRed) Name() string            { return "Shrill of Skullform" }
func (ShrillOfSkullformRed) Cost(*sim.TurnState) int { return 2 }
func (ShrillOfSkullformRed) Pitch() int              { return 1 }
func (ShrillOfSkullformRed) Attack() int             { return 4 }
func (ShrillOfSkullformRed) Defense() int            { return 3 }
func (ShrillOfSkullformRed) Types() card.TypeSet     { return shrillTypes }
func (ShrillOfSkullformRed) GoAgain() bool           { return false }
func (ShrillOfSkullformRed) Play(s *sim.TurnState, self *sim.CardState) {
	shrillPlay(s, self)
}

type ShrillOfSkullformYellow struct{}

func (ShrillOfSkullformYellow) ID() ids.CardID          { return ids.ShrillOfSkullformYellow }
func (ShrillOfSkullformYellow) Name() string            { return "Shrill of Skullform" }
func (ShrillOfSkullformYellow) Cost(*sim.TurnState) int { return 2 }
func (ShrillOfSkullformYellow) Pitch() int              { return 2 }
func (ShrillOfSkullformYellow) Attack() int             { return 3 }
func (ShrillOfSkullformYellow) Defense() int            { return 3 }
func (ShrillOfSkullformYellow) Types() card.TypeSet     { return shrillTypes }
func (ShrillOfSkullformYellow) GoAgain() bool           { return false }
func (ShrillOfSkullformYellow) Play(s *sim.TurnState, self *sim.CardState) {
	shrillPlay(s, self)
}

type ShrillOfSkullformBlue struct{}

func (ShrillOfSkullformBlue) ID() ids.CardID          { return ids.ShrillOfSkullformBlue }
func (ShrillOfSkullformBlue) Name() string            { return "Shrill of Skullform" }
func (ShrillOfSkullformBlue) Cost(*sim.TurnState) int { return 2 }
func (ShrillOfSkullformBlue) Pitch() int              { return 3 }
func (ShrillOfSkullformBlue) Attack() int             { return 2 }
func (ShrillOfSkullformBlue) Defense() int            { return 3 }
func (ShrillOfSkullformBlue) Types() card.TypeSet     { return shrillTypes }
func (ShrillOfSkullformBlue) GoAgain() bool           { return false }
func (ShrillOfSkullformBlue) Play(s *sim.TurnState, self *sim.CardState) {
	shrillPlay(s, self)
}

// shrillPlay routes the +3{p} aura-in-play buff through self.BonusAttack so EffectiveAttack
// and LikelyToHit see the buffed power, then emits the chain step at the buffed value. No
// rider sub-line — this is a power buff, not a separable effect.
func shrillPlay(s *sim.TurnState, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.BonusAttack += 3
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
