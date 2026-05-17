package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// modalBlocker is a test stub mirroring Brothers in Arms: mode 0 costs 0 and contributes
// the printed Defense (2); mode 1 costs 1 and contributes Defense + 2.
type modalBlocker struct{}

func (modalBlocker) ID() ids.CardID           { return ids.InvalidCard }
func (modalBlocker) Name() string             { return "modalBlocker" }
func (modalBlocker) DisplayName() string      { return "modalBlocker" }
func (modalBlocker) Cost(card.GameEngine) int { return 0 }
func (modalBlocker) Pitch() int               { return 0 }
func (modalBlocker) Attack() int              { return 0 }
func (modalBlocker) Defense() int             { return 2 }
func (modalBlocker) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (modalBlocker) GoAgain(card.GameEngine) bool                       { return false }
func (modalBlocker) Modes() int                                         { return 2 }
func (modalBlocker) BlockCost(mode int8) int                            { return int(mode) }
func (modalBlocker) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (modalBlocker) Block(_ card.GameEngine, _ card.Logger, self *card.CardState) {
	if self.Mode == 1 {
		self.BonusDefense += 2
	}
}

// Tests that defendersDamage picks mode 1 when the spare budget covers the per-mode cost.
func TestModalBlocker_FiresWhenBudgetCovers(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamageWithBudget(
		[]card.Card{modalBlocker{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		10, 1, -1,
	)
	if got != 4 {
		t.Errorf("budget=1 block = %d, want 4 (mode 1 fires: 2+2)", got)
	}
}

// Tests that defendersDamage falls back to mode 0 when the spare budget can't cover the
// modal cost.
func TestModalBlocker_FallsBackToMode0WhenBudgetTight(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamageWithBudget(
		[]card.Card{modalBlocker{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		10, 0, -1,
	)
	if got != 2 {
		t.Errorf("budget=0 block = %d, want 2 (mode 0: printed only)", got)
	}
}
