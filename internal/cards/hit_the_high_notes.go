// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var hitTheHighNotesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type HitTheHighNotesRed struct{}

func (HitTheHighNotesRed) ID() card.ID              { return card.HitTheHighNotesRed }
func (HitTheHighNotesRed) Name() string             { return "Hit the High Notes" }
func (HitTheHighNotesRed) Cost(*card.TurnState) int { return 1 }
func (HitTheHighNotesRed) Pitch() int               { return 1 }
func (HitTheHighNotesRed) Attack() int              { return 4 }
func (HitTheHighNotesRed) Defense() int             { return 3 }
func (HitTheHighNotesRed) Types() card.TypeSet      { return hitTheHighNotesTypes }
func (HitTheHighNotesRed) GoAgain() bool            { return false }
func (HitTheHighNotesRed) Play(s *card.TurnState, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	s.ApplyAndLogEffectiveAttack(self)
}

type HitTheHighNotesYellow struct{}

func (HitTheHighNotesYellow) ID() card.ID              { return card.HitTheHighNotesYellow }
func (HitTheHighNotesYellow) Name() string             { return "Hit the High Notes" }
func (HitTheHighNotesYellow) Cost(*card.TurnState) int { return 1 }
func (HitTheHighNotesYellow) Pitch() int               { return 2 }
func (HitTheHighNotesYellow) Attack() int              { return 3 }
func (HitTheHighNotesYellow) Defense() int             { return 3 }
func (HitTheHighNotesYellow) Types() card.TypeSet      { return hitTheHighNotesTypes }
func (HitTheHighNotesYellow) GoAgain() bool            { return false }
func (HitTheHighNotesYellow) Play(s *card.TurnState, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	s.ApplyAndLogEffectiveAttack(self)
}

type HitTheHighNotesBlue struct{}

func (HitTheHighNotesBlue) ID() card.ID              { return card.HitTheHighNotesBlue }
func (HitTheHighNotesBlue) Name() string             { return "Hit the High Notes" }
func (HitTheHighNotesBlue) Cost(*card.TurnState) int { return 1 }
func (HitTheHighNotesBlue) Pitch() int               { return 3 }
func (HitTheHighNotesBlue) Attack() int              { return 2 }
func (HitTheHighNotesBlue) Defense() int             { return 3 }
func (HitTheHighNotesBlue) Types() card.TypeSet      { return hitTheHighNotesTypes }
func (HitTheHighNotesBlue) GoAgain() bool            { return false }
func (HitTheHighNotesBlue) Play(s *card.TurnState, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
	s.ApplyAndLogEffectiveAttack(self)
}
func hitTheHighNotesBonus(s *card.TurnState) int {
	if s.HasPlayedOrCreatedAura() {
		return 2
	}
	return 0
}
