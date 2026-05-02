package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// stubAR is a minimal AttackReaction stub for partition-validator tests. It implements both
// sim.Card and sim.AttackReaction; allow controls which targets the predicate accepts.
type stubAR struct {
	id    ids.CardID
	name  string
	allow func(Card) bool
}

func (s stubAR) ID() ids.CardID    { return s.id }
func (s stubAR) Name() string      { return s.name }
func (stubAR) Cost(*TurnState) int { return 0 }
func (stubAR) Pitch() int          { return 3 }
func (stubAR) Attack() int         { return 0 }
func (stubAR) Defense() int        { return 0 }
func (stubAR) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)
}
func (stubAR) GoAgain() bool                 { return false }
func (s stubAR) ARTargetAllowed(c Card) bool { return s.allow(c) }
func (stubAR) Play(*TurnState, *CardState)   {}

// stubAttack is a Generic Action - Attack stub used as a target candidate.
type stubAttack struct{}

func (stubAttack) ID() ids.CardID      { return ids.InvalidCard }
func (stubAttack) Name() string        { return "stubAttack" }
func (stubAttack) Cost(*TurnState) int { return 0 }
func (stubAttack) Pitch() int          { return 1 }
func (stubAttack) Attack() int         { return 1 }
func (stubAttack) Defense() int        { return 0 }
func (stubAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubAttack) GoAgain() bool               { return true }
func (stubAttack) Play(*TurnState, *CardState) {}

// Tests that a chain with no Attack Reactions is vacuously valid — the validator only
// gates ARs.
func TestPartitionHasValidARTargets_NoARsAlwaysValid(t *testing.T) {
	chain := []Card{stubAttack{}, stubAttack{}}
	if !partitionHasValidARTargets(chain) {
		t.Error("chain with no ARs should validate")
	}
}

// Tests that an AR with at least one matching non-self target validates.
func TestPartitionHasValidARTargets_ARWithTargetValidates(t *testing.T) {
	ar := stubAR{name: "AR", allow: func(c Card) bool { return c.Types().IsAttackAction() }}
	chain := []Card{ar, stubAttack{}}
	if !partitionHasValidARTargets(chain) {
		t.Error("AR with a matching attack-action target should validate")
	}
}

// Tests that an AR with no matching target invalidates the chain.
func TestPartitionHasValidARTargets_ARWithoutTargetInvalid(t *testing.T) {
	ar := stubAR{name: "AR", allow: func(c Card) bool { return c.Types().IsAttackAction() }}
	chain := []Card{ar} // AR alone, nothing else to react to
	if partitionHasValidARTargets(chain) {
		t.Error("AR alone should fail validation (no target)")
	}
}

// Tests that an AR can't satisfy its own target requirement — self-targeting is excluded.
// Even if the predicate accepts the AR's own type set (Attack Reaction is not an attack
// action card, so this is mostly defensive), the validator skips self.
func TestPartitionHasValidARTargets_SelfTargetExcluded(t *testing.T) {
	// Predicate matches everything — including the AR itself.
	ar := stubAR{name: "AR", allow: func(c Card) bool { return true }}
	chain := []Card{ar} // Only the AR; if self-targeting were allowed it'd validate.
	if partitionHasValidARTargets(chain) {
		t.Error("AR matching only itself should fail validation (no distinct target)")
	}
}

// Tests that two ARs in a chain both need their own valid target — the validator scans
// every AR independently. Here AR1 wants attack actions (matches stubAttack), AR2 wants
// nothing in the chain.
func TestPartitionHasValidARTargets_AllARsMustHaveTarget(t *testing.T) {
	ar1 := stubAR{name: "AR1", allow: func(c Card) bool { return c.Types().IsAttackAction() }}
	ar2 := stubAR{name: "AR2", allow: func(c Card) bool { return false }}
	chain := []Card{ar1, ar2, stubAttack{}}
	if partitionHasValidARTargets(chain) {
		t.Error("chain should fail when one of two ARs has no target")
	}
}

// Tests that GrantAttackReactionBuff lands +n on the first matching CardState in
// CardsRemaining and leaves non-matching entries untouched.
func TestGrantAttackReactionBuff_LandsOnFirstMatch(t *testing.T) {
	target := &CardState{Card: stubAttack{}}
	skipped := &CardState{Card: stubAR{name: "skipMe", allow: func(Card) bool { return false }}}
	s := TurnState{CardsRemaining: []*CardState{skipped, target}}
	GrantAttackReactionBuff(&s, func(c Card) bool { return c.Types().IsAttackAction() }, 3)
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
	if skipped.BonusAttack != 0 {
		t.Errorf("skipped CardState mutated: BonusAttack = %d, want 0", skipped.BonusAttack)
	}
}

// Tests that GrantAttackReactionBuff fizzles silently when no CardState matches — orderings
// where the AR plays after every target.
func TestGrantAttackReactionBuff_NoMatchFizzles(t *testing.T) {
	skipped := &CardState{Card: stubAR{name: "skipMe", allow: func(Card) bool { return false }}}
	s := TurnState{CardsRemaining: []*CardState{skipped}}
	GrantAttackReactionBuff(&s, func(c Card) bool { return false }, 5)
	if skipped.BonusAttack != 0 {
		t.Errorf("skipped CardState mutated: BonusAttack = %d, want 0", skipped.BonusAttack)
	}
}
