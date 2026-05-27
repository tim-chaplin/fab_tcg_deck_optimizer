package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// pitchOnlyRed is a 0-cost, pitch-1, 0-attack, 0-defense red action: pure resource fodder for
// integration tests that need to fund a 1-cost card without leaving any defensive or
// offensive option in the partition. Non-attack action so Nimblism / Come to Fight / etc.
// don't see it as a candidate target for their "next attack action" grants.
type pitchOnlyRed struct{}

func (pitchOnlyRed) ID() ids.CardID      { return ids.InvalidCard }
func (pitchOnlyRed) Name() string        { return "pitchOnlyRed" }
func (pitchOnlyRed) DisplayName() string { return "pitchOnlyRed [R]" }
func (pitchOnlyRed) Cost() int           { return 0 }
func (pitchOnlyRed) Pitch() int          { return 1 }
func (pitchOnlyRed) Attack() int         { return 0 }
func (pitchOnlyRed) Defense() int        { return 0 }
func (pitchOnlyRed) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (pitchOnlyRed) GoAgain(card.GameEngine) bool                       { return false }
func (pitchOnlyRed) Play(card.GameEngine, card.Logger, *card.CardState) {}

// grantBonusAttack is a test-only non-attack action card that scans CardsRemaining and adds n
// to BonusAttack on the first attack action card it finds. Mirrors production "next attack
// +N{p}" grants: the buff lands on the target's CardState so it's attributed to the attack
// being buffed and feeds EffectiveAttack for any "if this hits" rider on that target.
type grantBonusAttack struct{ n int }

func (grantBonusAttack) ID() ids.CardID      { return ids.InvalidCard }
func (grantBonusAttack) Name() string        { return "grantBonusAttack" }
func (grantBonusAttack) DisplayName() string { return "grantBonusAttack" }
func (grantBonusAttack) Cost() int           { return 0 }
func (grantBonusAttack) Pitch() int          { return 0 }
func (grantBonusAttack) Attack() int         { return 0 }
func (grantBonusAttack) Defense() int        { return 0 }
func (grantBonusAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (grantBonusAttack) GoAgain(card.GameEngine) bool { return true }
func (c grantBonusAttack) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	for _, pc := range ge.CardsRemaining() {
		if pc.Card.Types(nil).IsAttackAction() {
			pc.BonusAttack += c.n
			break
		}
	}
}

// grantBonusAttackWeapon scans CardsRemaining for the first weapon swing (a weapon
// activated ability, identified by IsWeaponAttack — TypeWeapon + TypeAttack) and adds n
// to its BonusAttack. Mirrors the production shape of Brandish's "next weapon attack
// +1{p}" rider — the target is a weapon swing, not an attack action.
type grantBonusAttackWeapon struct{ n int }

func (grantBonusAttackWeapon) ID() ids.CardID      { return ids.InvalidCard }
func (grantBonusAttackWeapon) Name() string        { return "grantBonusAttackWeapon" }
func (grantBonusAttackWeapon) DisplayName() string { return "grantBonusAttackWeapon" }
func (grantBonusAttackWeapon) Cost() int           { return 0 }
func (grantBonusAttackWeapon) Pitch() int          { return 0 }
func (grantBonusAttackWeapon) Attack() int         { return 0 }
func (grantBonusAttackWeapon) Defense() int        { return 0 }
func (grantBonusAttackWeapon) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (grantBonusAttackWeapon) GoAgain(card.GameEngine) bool { return true }
func (c grantBonusAttackWeapon) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	for _, pc := range ge.CardsRemaining() {
		if pc.Card.Types(nil).IsWeaponAttack() {
			pc.BonusAttack += c.n
			break
		}
	}
}

// Tests that a granter writes BonusAttack on the target's CardState and the attack turn total
// reflects printed-attack + bonus.
func TestPlaySequence_BonusAttackAppliedToTargetDamage(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, testutils.FakeRedAttack().WithPower(3)}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter→RedAttack to extend via go-again")
	}
	// Granter (cost 0, attack 0, go again) → RedAttack (cost 1, printed power 3, bonus +3).
	// Total: 0 + (3 + 3) = 6.
	if dmg != 6 {
		t.Fatalf("dmg = %d, want 6 (RedAttack 3 + granted bonus 3)", dmg)
	}
}

// TestPlaySequence_BonusAttackNoTargetFizzles pins the no-target case: a granter alone
// scans CardsRemaining, finds no attack action, and contributes nothing — the BonusAttack
// state simply stays 0.
func TestPlaySequence_BonusAttackNoTargetFizzles(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter alone to be a legal 1-card attack turn")
	}
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (granter has no damage and no target to buff)", dmg)
	}
}

// TestPlaySequence_BonusAttackStacksAcrossGranters pins that two granters in front of the
// same target both write to BonusAttack; the field accumulates rather than overwriting.
func TestPlaySequence_BonusAttackStacksAcrossGranters(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, grantBonusAttack{n: 2}, testutils.FakeRedAttack().WithPower(3)}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected two granters→RedAttack to extend via go-again")
	}
	// Granter +3 → granter +2 → RedAttack (cost 1, printed power 3, bonus 3+2=5). Total 0+0+8 = 8.
	if dmg != 8 {
		t.Fatalf("dmg = %d, want 8 (RedAttack 3 + stacked grants 5)", dmg)
	}
}

// Tests that BonusAttack applies to weapon swings (TypeWeapon + TypeAttack), not just
// attack action cards.
func TestPlaySequence_BonusAttackAppliesToWeapon(t *testing.T) {
	order := []card.Card{grantBonusAttackWeapon{n: 2}, weapons.ReapingBlade{}.Ability()}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter→weapon swing to extend via go-again")
	}
	// Granter (cost 0, returns 0) → Reaping Blade (cost 1, printed power 3, bonus +2 = 5).
	// Total: 0 + 5 = 5.
	if dmg != 5 {
		t.Fatalf("dmg = %d, want 5 (Reaping Blade 3 + granted bonus 2)", dmg)
	}
}

// Tests that a negative BonusAttack clamps the target's contribution at 0 (FaB attack-power
// floor) — a 1-power attack with a -3 grant deals 0, not -2.
func TestPlaySequence_BonusAttackNegativeClampsAtZero(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: -3}, testutils.FakeBlueAttack().WithPower(1)}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	// Granter (returns 0, no own attack) → BlueAttack (printed power 1, bonus -3 →
	// pre-clamp -2, post-clamp 0). Total 0+0 = 0.
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (1-power attack with -3 bonus floors at 0)", dmg)
	}
}

// TestPlaySequence_BonusAttackNegativePartialReduction pins the in-range case: a negative
// grant that doesn't drive the target below 0 reduces the contribution by the full bonus,
// no clamp.
func TestPlaySequence_BonusAttackNegativePartialReduction(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: -2}, testutils.FakeRedAttack().WithPower(3)}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	// Granter → RedAttack (printed power 3, bonus -2 → 1). Total 1.
	if dmg != 1 {
		t.Fatalf("dmg = %d, want 1 (RedAttack 3 - debuff 2)", dmg)
	}
}

// TestPlaySequence_BonusAttackNoAttackTargetFizzles pins the granter-side scan: if no attack
// action follows the granter, the rider has nowhere to land and total damage stays 0.
func TestPlaySequence_BonusAttackNoAttackTargetFizzles(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, grantBonusAttack{n: 2}}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.PlaySequence(order)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (no attack actions present; both granters fizzle)", dmg)
	}
}

// Tests the per-permutation BonusAttack reset: two back-to-back playSequence calls produce
// the same total, never a leaked-bonus regression.
func TestPlaySequence_BonusAttackPerPermutationReset(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, testutils.FakeRedAttack().WithPower(3)}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 10, 0, len(order))
	first, _, _, _ := ctx.PlaySequence(order)
	second, _, _, _ := ctx.PlaySequence(order)
	if first != 6 || second != 6 {
		t.Fatalf("non-deterministic damage across reuses: first=%d, second=%d, want both=6", first, second)
	}
}

// Tests end-to-end that a Nimblism BonusAttack grant pushes Consuming Volition past the
// likely-to-hit threshold so its arcane-damage discard rider fires.
func TestBest_NimblismGrantsConsumingVolitionDiscardRider(t *testing.T) {
	h := []card.Card{
		cards.ConsumingVolitionYellow{},
		cards.NimblismBlue{},
		pitchOnlyRed{},
	}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		AddAura(token.NewRunechant(1)).
		Build()
	got := Best(nil, h, nil, state)
	if got.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Volition 3 base + Nimblism +1 BonusAttack + discard rider 3 from runechant-driven ArcaneDamageDealt × s.LikelyToHit on buffed 4-power attack); line=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}
