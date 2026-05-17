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
	h := []deck.Card{testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 6 {
		t.Fatalf("want value 6, got %d", extras.Value)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented
	// 3). Value = 5.
	h := []deck.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 5 {
		t.Fatalf("want value 5, got %d", extras.Value)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented
	// 3). Value = 9.
	h := []deck.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 9 {
		t.Fatalf("want value 9, got %d", extras.Value)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at
	// incoming=2). Value = 4.
	h := []deck.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.BlueAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 2}, nil, h)
	if extras.Value != 4 {
		t.Fatalf("want value 4, got %d", extras.Value)
	}
}

func TestBest_DefenseReactionRequiresCostPaid(t *testing.T) {
	// Toughen Up [B]: Cost 2, Pitch 3, Defense 4. A hand of just this card can't pay its own
	// 2-resource cost to play as a Defense Reaction (there's nothing else to pitch). The only
	// legal lines are to pitch it (0 damage prevented) or do nothing — Value must be 0.
	h := []deck.Card{cards.ToughenUpBlue{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 0 {
		t.Fatalf("want value 0 (cost unpaid), got %d", extras.Value)
	}
}

func TestBest_DefenseReactionAffordableResolves(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), pay Toughen Up [B]'s cost 2, prevent 4 damage (capped at
	// incoming=4). Value = 4.
	h := []deck.Card{cards.MaleficIncantationBlue{}, cards.ToughenUpBlue{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 4 {
		t.Fatalf("want value 4 (cost paid, full block), got %d", extras.Value)
	}
}

func TestBest_PlainBlockStillFree(t *testing.T) {
	// Attack cards have no Defense-Reaction type, so using them as blockers costs nothing. One
	// Red attacker (Defense 1) alone, used as a blocker against 1 incoming, prevents 1. Value = 1.
	h := []deck.Card{testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 1}, nil, h)
	if extras.Value != 1 {
		t.Fatalf("want value 1 (free plain block), got %d", extras.Value)
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must
	// cover costs.
	h := []deck.Card{testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, h)
	if extras.Value != 6 {
		t.Fatalf("want value 6, got %d", extras.Value)
	}
	var res, cost int
	for i, c := range h {
		switch extras.BestLine[i].Role {
		case deck.Pitch:
			res += c.(card.Card).Pitch()
		case deck.Attack:
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
	h := []deck.Card{cards.ToughenUpBlue{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 0 {
		t.Fatalf("Value = %d, want 0", extras.Value)
	}
	if extras.BestLine[0].Role != deck.Arsenal {
		t.Errorf("role = %s, want ARSENAL (empty slot + Held card → promoted)", extras.BestLine[0].Role)
	}
	if gs.Arsenal() == nil || gs.Arsenal().ID() != ids.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue", gs.Arsenal())
	}
}

// Tests that attack-phase and defense-phase pitches draw from disjoint pools — a single
// pitched card can't fund both phases.
func TestBest_AttackPitchCantCoverDefense(t *testing.T) {
	h := []deck.Card{cards.MaleficIncantationBlue{}, cards.ToughenUpBlue{}, testutils.RedAttack{}}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 5 {
		t.Fatalf("Value = %d, want 5 (attack and defense pitches are separate pools; Roles=[%s])",
			extras.Value, sim.FormatBestLine(extras.BestLine))
	}
}

// Tests that two pitched cards unlock the attack/defense split — one pitch funds each phase.
func TestBest_DRPitchNeedsSecondPitchedCard(t *testing.T) {
	h := []deck.Card{
		cards.MaleficIncantationBlue{},
		cards.MaleficIncantationBlue{},
		cards.ToughenUpBlue{},
		testutils.RedAttack{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 4}, nil, h)
	if extras.Value != 7 {
		t.Fatalf("Value = %d, want 7 (two pitched cards let attack + defense phases both pay; Roles=[%s])",
			extras.Value, sim.FormatBestLine(extras.BestLine))
	}
}

// Tests attackBufs scratch sizing on a full 4-card hand of 0-cost attackers plus an
// arsenal-in attacker (5 attackers, no weapons) — guards against slice-bounds panics.
func TestBest_AllAttackHandPlusArsenalNoWeapons(t *testing.T) {
	h := []deck.Card{
		cards.WoundingBlowRed{}, cards.WoundingBlowRed{},
		cards.WoundingBlowRed{}, cards.WoundingBlowRed{},
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.GameStateBuilder().SetArsenal(cards.WoundingBlowRed{}).Build(), h)
	if extras.Value != 4 {
		t.Fatalf("Value = %d, want 4 (one Wounding Blow Red lands; rest can't chain without GoAgain). Roles=[%s]",
			extras.Value, sim.FormatBestLine(extras.BestLine))
	}
}
