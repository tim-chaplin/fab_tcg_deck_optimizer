package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that mode 0 accepts Club and Hammer weapon attacks (probed via the activated-
// ability Cards the attack-turn runner enqueues for the swing).
func TestPummel_Mode0AcceptsClubAndHammer(t *testing.T) {
	if !(PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeWeaponSwing().WithTypes(card.TypeClub, card.TypeOneHand)}, 0) {
		t.Error("mode 0 should accept a Club weapon ability")
	}
	if !(PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeWeaponSwing().WithTypes(card.TypeHammer, card.TypeOneHand)}, 0) {
		t.Error("mode 0 should accept a Hammer weapon ability")
	}
}

// Tests that mode 0 rejects a non-club/hammer target.
func TestPummel_Mode0RejectsOtherTargets(t *testing.T) {
	if (PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack().WithCost(2)}, 0) {
		t.Error("mode 0 should reject a non-club/hammer attack")
	}
}

// Tests that mode 0 rejects a club-typed attack action card — the printed text says
// "weapon attack", so an action card sharing the Club subtype shouldn't qualify.
func TestPummel_Mode0RejectsClubAttackActionCard(t *testing.T) {
	clubAction := testutils.FakeRedAttack().WithTypes(card.TypeClub)
	if (PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: clubAction}, 0) {
		t.Error("mode 0 should reject a club attack action card (only weapon attacks qualify)")
	}
}

// Tests that mode 1 accepts a cost-≥2 attack action.
func TestPummel_Mode1AcceptsCostTwoAttackAction(t *testing.T) {
	if !(PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack().WithCost(2)}, 1) {
		t.Error("mode 1 should accept a cost-2 attack action")
	}
}

// Tests that mode 1 rejects cost-1 attack actions and non-attacks.
func TestPummel_Mode1RejectsCostOneAndNonAttacks(t *testing.T) {
	if (PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack().WithCost(1)}, 1) {
		t.Error("mode 1 should reject a cost-1 attack action")
	}
	if (PummelRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAction()}, 1) {
		t.Error("mode 1 should reject a non-attack action")
	}
}
