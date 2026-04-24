package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
)

func TestBest_AllRedHand(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with the other 2 (cost 2, dealt 6). Value = 6.
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 6 {
		t.Fatalf("want value 6, got %d", got.Value)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented
	// 3). Value = 5.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 5 {
		t.Fatalf("want value 5, got %d", got.Value)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented
	// 3). Value = 9.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 9 {
		t.Fatalf("want value 9, got %d", got.Value)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at
	// incoming=2). Value = 4.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}}
	got := Best(stubHero, nil, h, 2, nil, 0, nil)
	if got.Value != 4 {
		t.Fatalf("want value 4, got %d", got.Value)
	}
}

func TestBest_DefenseReactionRequiresCostPaid(t *testing.T) {
	// Toughen Up (Blue): Cost 2, Pitch 3, Defense 4. A hand of just this card can't pay its own
	// 2-resource cost to play as a Defense Reaction (there's nothing else to pitch). The only
	// legal lines are to pitch it (0 damage prevented) or do nothing — Value must be 0.
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 0 {
		t.Fatalf("want value 0 (cost unpaid), got %d", got.Value)
	}
}

func TestBest_DefenseReactionAffordableResolves(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), pay Toughen Up (Blue)'s cost 2, prevent 4 damage (capped at
	// incoming=4). Value = 4.
	h := []card.Card{runeblade.MaleficIncantationBlue{}, generic.ToughenUpBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 4 {
		t.Fatalf("want value 4 (cost paid, full block), got %d", got.Value)
	}
}

func TestBest_PlainBlockStillFree(t *testing.T) {
	// Attack cards have no Defense-Reaction type, so using them as blockers costs nothing. One
	// Red attacker (Defense 1) alone, used as a blocker against 1 incoming, prevents 1. Value = 1.
	h := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 1, nil, 0, nil)
	if got.Value != 1 {
		t.Fatalf("want value 1 (free plain block), got %d", got.Value)
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must
	// cover costs.
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 0, nil, 0, nil)
	if got.Value != 6 {
		t.Fatalf("want value 6, got %d", got.Value)
	}
	var res, cost int
	for i, c := range h {
		switch got.BestLine[i].Role {
		case Pitch:
			res += c.Pitch()
		case Attack:
			cost += c.Cost(&card.TurnState{})
		}
	}
	if res < cost {
		t.Fatalf("invalid play: resources %d < costs %d", res, cost)
	}
}

// TestBest_AllHeldWhenNoLegalPlay covers the "hand does nothing this turn" case. A single
// Toughen Up Blue DR (cost 2) with no pitched cards to pay it has Value = 0. The partition
// leaves the card Held; post-hoc the empty arsenal slot claims it, so Role becomes Arsenal
// and Play.ArsenalCard records the card for next turn's carryover.
func TestBest_AllHeldWhenNoLegalPlay(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 0 {
		t.Fatalf("Value = %d, want 0", got.Value)
	}
	if got.BestLine[0].Role != Arsenal {
		t.Errorf("role = %s, want ARSENAL (empty slot + Held card → promoted)", got.BestLine[0].Role)
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue", got.ArsenalCard)
	}
	if got.BestLine[0].Contribution != 0 {
		t.Errorf("Contribution = %.1f, want 0 (card sits in arsenal, real value accrues on a later turn)", got.BestLine[0].Contribution)
	}
}

// TestBest_AttackPitchCantCoverDefense enforces that attack-phase and defense-phase pitches
// draw from disjoint pools (resources don't cross turns). Hand: Malefic Blue (cost 0, pitch 3,
// defense 2), Toughen Up Blue (DR, cost 2, pitch 3, defense 4), Red Attack (cost 1, pitch 1,
// attack 3). Against incoming 4: only one pitched card (pitch 3) can be paired with Toughen Up
// as DR, and that single pitch can cover either the 1-cost Red OR the 2-cost Toughen Up — not
// both. The solver takes the better single-phase line: pitch Toughen Up to pay Red's cost,
// plain-block with Malefic. Value = 3 (Red attack) + 2 (Malefic block) = 5. A single-pool
// fallback would score 7 by funding both phases from one pitch — illegal, locked out here.
func TestBest_AttackPitchCantCoverDefense(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}, generic.ToughenUpBlue{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 5 {
		t.Fatalf("Value = %d, want 5 (attack and defense pitches are separate pools; Roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_DRPitchNeedsSecondPitchedCard confirms that adding a second pitched card unlocks the
// split: Malefic Blue (pitch 3) + second Malefic Blue (pitch 3) is enough to fund Red's 1-cost
// attack from one Malefic and Toughen Up's 2-cost defense from the other. Value = 3 (attack) +
// 4 (full prevent) = 7.
func TestBest_DRPitchNeedsSecondPitchedCard(t *testing.T) {
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		generic.ToughenUpBlue{},
		fake.RedAttack{},
	}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Value != 7 {
		t.Fatalf("Value = %d, want 7 (two pitched cards let attack + defense phases both pay; Roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_AllAttackHandPlusArsenalNoWeapons guards attackBufs scratch sizing against slice-
// bounds panics when a full 4-card hand of 0-cost attackers plus an arsenal-in attacker (5
// entries) goes through bestAttackWithWeapons with no weapons. Wounding Blow Red is 0-cost
// attack 4, so the all-Attack partition is phase-feasible with zero pitches and enumerates the
// 5-attacker path. The winning line only chains one (no GoAgain), but the enumerator still
// evaluates the 5-attacker partition — the buffer must survive that.
func TestBest_AllAttackHandPlusArsenalNoWeapons(t *testing.T) {
	h := []card.Card{
		generic.WoundingBlowRed{}, generic.WoundingBlowRed{},
		generic.WoundingBlowRed{}, generic.WoundingBlowRed{},
	}
	got := Best(stubHero, nil, h, 0, nil, 0, generic.WoundingBlowRed{})
	if got.Value != 4 {
		t.Fatalf("Value = %d, want 4 (one Wounding Blow Red lands; rest can't chain without GoAgain). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}
