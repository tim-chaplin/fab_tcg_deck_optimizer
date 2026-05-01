// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var hitTheHighNotesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type HitTheHighNotesRed struct{}

func (HitTheHighNotesRed) ID() ids.CardID          { return ids.HitTheHighNotesRed }
func (HitTheHighNotesRed) Name() string            { return "Hit the High Notes" }
func (HitTheHighNotesRed) Cost(*sim.TurnState) int { return 1 }
func (HitTheHighNotesRed) Pitch() int              { return 1 }
func (HitTheHighNotesRed) Attack() int             { return 4 }
func (HitTheHighNotesRed) Defense() int            { return 3 }
func (HitTheHighNotesRed) Types() card.TypeSet     { return hitTheHighNotesTypes }
func (HitTheHighNotesRed) GoAgain() bool           { return false }
func (HitTheHighNotesRed) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type HitTheHighNotesYellow struct{}

func (HitTheHighNotesYellow) ID() ids.CardID          { return ids.HitTheHighNotesYellow }
func (HitTheHighNotesYellow) Name() string            { return "Hit the High Notes" }
func (HitTheHighNotesYellow) Cost(*sim.TurnState) int { return 1 }
func (HitTheHighNotesYellow) Pitch() int              { return 2 }
func (HitTheHighNotesYellow) Attack() int             { return 3 }
func (HitTheHighNotesYellow) Defense() int            { return 3 }
func (HitTheHighNotesYellow) Types() card.TypeSet     { return hitTheHighNotesTypes }
func (HitTheHighNotesYellow) GoAgain() bool           { return false }
func (HitTheHighNotesYellow) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type HitTheHighNotesBlue struct{}

func (HitTheHighNotesBlue) ID() ids.CardID          { return ids.HitTheHighNotesBlue }
func (HitTheHighNotesBlue) Name() string            { return "Hit the High Notes" }
func (HitTheHighNotesBlue) Cost(*sim.TurnState) int { return 1 }
func (HitTheHighNotesBlue) Pitch() int              { return 3 }
func (HitTheHighNotesBlue) Attack() int             { return 2 }
func (HitTheHighNotesBlue) Defense() int            { return 3 }
func (HitTheHighNotesBlue) Types() card.TypeSet     { return hitTheHighNotesTypes }
func (HitTheHighNotesBlue) GoAgain() bool           { return false }
func (HitTheHighNotesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
func hitTheHighNotesBonus(s *sim.TurnState) int {
	if s.HasPlayedOrCreatedAura() {
		return 2
	}
	return 0
}
