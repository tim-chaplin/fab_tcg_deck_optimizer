package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// flipsAndReadsDominate is a stub Action - Attack whose Play increments *saw whenever
// self.GrantedDominate is already true at Play entry, then sets it true. Across the N!
// permutations bestSequence walks, *saw must stay 0: the chain runner has to reset the
// flag between permutations or a previous permutation's grant leaks into the next one.
//
// Mirrors the production shape of cards (Drowning Dire / Pound for Pound / Regurgitating
// Slog) that conditionally flip self.GrantedDominate — the bug would inflate their
// EffectiveDominate in later permutations even when the condition didn't fire.
type flipsAndReadsDominate struct{ saw *int }

func (flipsAndReadsDominate) ID() ids.CardID           { return ids.InvalidCard }
func (flipsAndReadsDominate) Name() string             { return "flipsAndReadsDominate" }
func (flipsAndReadsDominate) DisplayName() string      { return "flipsAndReadsDominate" }
func (flipsAndReadsDominate) Cost(card.GameEngine) int { return 0 }
func (flipsAndReadsDominate) Pitch() int               { return 0 }
func (flipsAndReadsDominate) Attack() int              { return 3 }
func (flipsAndReadsDominate) Defense() int             { return 0 }
func (flipsAndReadsDominate) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (flipsAndReadsDominate) GoAgain(card.GameEngine) bool { return true }
func (s flipsAndReadsDominate) Play(_ card.GameEngine, _ card.Logger, self *card.CardState) {
	if self.GrantedDominate {
		*s.saw++
	}
	self.GrantedDominate = true
}

// flipsAndReadsOverpower mirrors flipsAndReadsDominate for the GrantedOverpower field
// — Vantage Point's shape, which conditionally sets self.GrantedOverpower on an aura
// having been created this turn.
type flipsAndReadsOverpower struct{ saw *int }

func (flipsAndReadsOverpower) ID() ids.CardID           { return ids.InvalidCard }
func (flipsAndReadsOverpower) Name() string             { return "flipsAndReadsOverpower" }
func (flipsAndReadsOverpower) DisplayName() string      { return "flipsAndReadsOverpower" }
func (flipsAndReadsOverpower) Cost(card.GameEngine) int { return 0 }
func (flipsAndReadsOverpower) Pitch() int               { return 0 }
func (flipsAndReadsOverpower) Attack() int              { return 3 }
func (flipsAndReadsOverpower) Defense() int             { return 0 }
func (flipsAndReadsOverpower) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (flipsAndReadsOverpower) GoAgain(card.GameEngine) bool { return true }
func (s flipsAndReadsOverpower) Play(_ card.GameEngine, _ card.Logger, self *card.CardState) {
	if self.GrantedOverpower {
		*s.saw++
	}
	self.GrantedOverpower = true
}

// flipsAndReadsBonusDefense mirrors the pattern for BonusDefense — Sigil of Suffering's
// shape (self.BonusDefense++ inside Play). With the per-permutation reset missing, the
// counter would accumulate across the N! enumeration.
type flipsAndReadsBonusDefense struct{ maxSeen *int }

func (flipsAndReadsBonusDefense) ID() ids.CardID           { return ids.InvalidCard }
func (flipsAndReadsBonusDefense) Name() string             { return "flipsAndReadsBonusDefense" }
func (flipsAndReadsBonusDefense) DisplayName() string      { return "flipsAndReadsBonusDefense" }
func (flipsAndReadsBonusDefense) Cost(card.GameEngine) int { return 0 }
func (flipsAndReadsBonusDefense) Pitch() int               { return 0 }
func (flipsAndReadsBonusDefense) Attack() int              { return 3 }
func (flipsAndReadsBonusDefense) Defense() int             { return 0 }
func (flipsAndReadsBonusDefense) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (flipsAndReadsBonusDefense) GoAgain(card.GameEngine) bool { return true }
func (s flipsAndReadsBonusDefense) Play(_ card.GameEngine, _ card.Logger, self *card.CardState) {
	self.BonusDefense++
	if self.BonusDefense > *s.maxSeen {
		*s.maxSeen = self.BonusDefense
	}
}

// TestBestSequence_GrantedDominateResetPerPermutation: the partial reset inside
// playSequenceWithMeta (between permutations of the same bestSequence enumeration) must
// clear pcBuf[i].GrantedDominate so a card that conditionally sets the flag never
// observes a stale-true value from an earlier permutation.
func TestBestSequence_GrantedDominateResetPerPermutation(t *testing.T) {
	saw := 0
	attackers := []card.Card{flipsAndReadsDominate{saw: &saw}, testutils.RedAttack{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(attackers))
	ctx.BestSequence(attackers)
	if saw != 0 {
		t.Errorf("GrantedDominate leaked across permutations: Play observed stale-true %d times, want 0", saw)
	}
}

// TestBestSequence_GrantedOverpowerResetPerPermutation: ditto for GrantedOverpower.
func TestBestSequence_GrantedOverpowerResetPerPermutation(t *testing.T) {
	saw := 0
	attackers := []card.Card{flipsAndReadsOverpower{saw: &saw}, testutils.RedAttack{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(attackers))
	ctx.BestSequence(attackers)
	if saw != 0 {
		t.Errorf("GrantedOverpower leaked across permutations: Play observed stale-true %d times, want 0", saw)
	}
}

// TestBestSequence_BonusDefenseResetPerPermutation: BonusDefense must reset between
// permutations or a card incrementing it would accumulate across the N! enumeration.
func TestBestSequence_BonusDefenseResetPerPermutation(t *testing.T) {
	maxSeen := 0
	attackers := []card.Card{flipsAndReadsBonusDefense{maxSeen: &maxSeen}, testutils.RedAttack{}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(attackers))
	ctx.BestSequence(attackers)
	if maxSeen != 1 {
		t.Errorf("BonusDefense leaked across permutations: max observed = %d, want 1", maxSeen)
	}
}
