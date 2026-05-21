package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestBest_AllRedHand(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with the other 2 (cost 2, dealt 6). Value = 6.
	h := []card.Card{testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 6 {
		t.Fatalf("want value 6, got %d", summary.Value)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented
	// 3). Value = 5.
	h := []card.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 5 {
		t.Fatalf("want value 5, got %d", summary.Value)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented
	// 3). Value = 9.
	h := []card.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 9 {
		t.Fatalf("want value 9, got %d", summary.Value)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at
	// incoming=2). Value = 4.
	h := []card.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(2).Build(), h)
	if summary.Value != 4 {
		t.Fatalf("want value 4, got %d", summary.Value)
	}
}

func TestBest_DefenseReactionRequiresCostPaid(t *testing.T) {
	// Toughen Up [B]: Cost 2, Pitch 3, Defense 4. A hand of just this card can't pay its own
	// 2-resource cost to play as a Defense Reaction (there's nothing else to pitch). The only
	// legal lines are to pitch it (0 damage prevented) or do nothing — Value must be 0.
	h := []card.Card{cards.ToughenUpBlue{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 0 {
		t.Fatalf("want value 0 (cost unpaid), got %d", summary.Value)
	}
}

func TestBest_DefenseReactionAffordableResolves(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), pay Toughen Up [B]'s cost 2, prevent 4 damage (capped at
	// incoming=4). Value = 4.
	h := []card.Card{cards.MaleficIncantationBlue{}, cards.ToughenUpBlue{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 4 {
		t.Fatalf("want value 4 (cost paid, full block), got %d", summary.Value)
	}
}

func TestBest_PlainBlockStillFree(t *testing.T) {
	// Attack cards have no Defense-Reaction type, so using them as blockers costs nothing. One
	// Red attacker (Defense 1) alone, used as a blocker against 1 incoming, prevents 1. Value = 1.
	h := []card.Card{testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(1).Build(), h)
	if summary.Value != 1 {
		t.Fatalf("want value 1 (free plain block), got %d", summary.Value)
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must
	// cover costs.
	h := []card.Card{testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, nil, h)
	if summary.Value != 6 {
		t.Fatalf("want value 6, got %d", summary.Value)
	}
	var res, cost int
	for i, c := range h {
		switch summary.BestLine[i].Role {
		case card.Pitch:
			res += c.(card.Card).Pitch()
		case card.Attack:
			cost += c.(card.Card).Cost(gameengine.New())
		}
	}
	if res < cost {
		t.Fatalf("invalid play: resources %d < costs %d", res, cost)
	}
}

// Tests the "hand does nothing this turn" case: a Held card with an empty arsenal gets
// post-hoc promoted to Arsenal.
func TestBest_AllHeldWhenNoLegalPlay(t *testing.T) {
	h := []card.Card{cards.ToughenUpBlue{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 0 {
		t.Fatalf("Value = %d, want 0", summary.Value)
	}
	if summary.BestLine[0].Role != card.Arsenal {
		t.Errorf("role = %s, want ARSENAL (empty slot + Held card → promoted)", summary.BestLine[0].Role)
	}
	if summary.State.Arsenal() == nil || summary.State.Arsenal().ID() != ids.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue", summary.State.Arsenal())
	}
}

// Tests that attack-phase and defense-phase pitches draw from disjoint pools — a single
// pitched card can't fund both phases.
func TestBest_AttackPitchCantCoverDefense(t *testing.T) {
	h := []card.Card{cards.MaleficIncantationBlue{}, cards.ToughenUpBlue{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 5 {
		t.Fatalf("Value = %d, want 5 (attack and defense pitches are separate pools; Roles=[%s])",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests that two pitched cards unlock the attack/defense split — one pitch funds each phase.
func TestBest_DRPitchNeedsSecondPitchedCard(t *testing.T) {
	h := []card.Card{
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.ToughenUpBlue{},
		testutils.RedAttack{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(4).Build(), h)
	if summary.Value != 7 {
		t.Fatalf("Value = %d, want 7 (two pitched cards let attack + defense phases both pay; Roles=[%s])",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests attackBufs scratch sizing on a full 4-card hand of 0-cost attackers plus an
// arsenal-in attacker (5 attackers, no weapons) — guards against slice-bounds panics.
func TestBest_AllAttackHandPlusArsenalNoWeapons(t *testing.T) {
	h := []card.Card{
		cards.WoundingBlowRed{}, cards.WoundingBlowRed{},
		cards.WoundingBlowRed{}, cards.WoundingBlowRed{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetArsenal(cards.WoundingBlowRed{}).Build(), h)
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (one Wounding Blow Red lands; rest can't chain without GoAgain). Roles=[%s]",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
}
