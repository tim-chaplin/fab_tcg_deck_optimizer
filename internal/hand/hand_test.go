package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// stubHero is a no-op Hero used by tests that want to measure raw hand
// value without any hero-ability contribution.
type stubHero struct{}

func (stubHero) Name() string                          { return "stubHero" }
func (stubHero) Health() int                           { return 20 }
func (stubHero) Intelligence() int                     { return 4 }
func (stubHero) Types() card.TypeSet                   { return 0 }
func (stubHero) OnCardPlayed(card.Card, *card.TurnState) int { return 0 }

func TestBest_AllRedHand(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with the other 2 (cost 2, dealt 6). Value = 6.
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.Value != 6 {
		t.Fatalf("want value 6, got %d", got.Value)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented
	// 3). Value = 5.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.Value != 5 {
		t.Fatalf("want value 5, got %d", got.Value)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented
	// 3). Value = 9.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.Value != 9 {
		t.Fatalf("want value 9, got %d", got.Value)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at
	// incoming=2). Value = 4.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}}
	got := Best(stubHero{}, nil, h, 2, nil, 0, nil)
	if got.Value != 4 {
		t.Fatalf("want value 4, got %d", got.Value)
	}
}

func TestBest_DefenseReactionRequiresCostPaid(t *testing.T) {
	// Toughen Up (Blue): Cost 2, Pitch 3, Defense 4. A hand of just this card can't pay its own
	// 2-resource cost to play as a Defense Reaction (there's nothing else to pitch). The only
	// legal lines are to pitch it (0 damage prevented) or do nothing — Value must be 0.
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.Value != 0 {
		t.Fatalf("want value 0 (cost unpaid), got %d", got.Value)
	}
}

func TestBest_DefenseReactionAffordableResolves(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), pay Toughen Up (Blue)'s cost 2, prevent 4 damage (capped at
	// incoming=4). Value = 4.
	h := []card.Card{runeblade.MaleficIncantationBlue{}, generic.ToughenUpBlue{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.Value != 4 {
		t.Fatalf("want value 4 (cost paid, full block), got %d", got.Value)
	}
}

func TestBest_PlainBlockStillFree(t *testing.T) {
	// Attack cards have no Defense-Reaction type, so using them as blockers costs nothing. One
	// Red attacker (Defense 1) alone, used as a blocker against 1 incoming, prevents 1. Value = 1.
	h := []card.Card{fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 1, nil, 0, nil)
	if got.Value != 1 {
		t.Fatalf("want value 1 (free plain block), got %d", got.Value)
	}
}

func TestBest_ViseraiMaleficShrillCombo(t *testing.T) {
	// Hero = Viserai. Best line: pitch the Blue Malefic, then play both Red Maleficas and the Red
	// Shrill. Value = 15.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationRed{},
		runeblade.MaleficIncantationRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 4, nil, 0, nil)
	if got.Value != 15 {
		t.Fatalf("want value 15, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiReapingBladeBlueMalefics(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play the other 3 Blue Malefics (Runechants from Viserai on #2
	// and #3), then swing Reaping Blade (cost 1, 3 dmg). Value = 3 + 2 + 3 = 8.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 8 {
		t.Fatalf("want value 8, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiReapingBladeMaleficsPlusShrill(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play 2 Blue Malefics (2 dmg + 1 Runechant), then Red Shrill
	// (cost 2, 4+3 aura bonus + 1 Runechant = 8). Reaping Blade stays holstered — Shrill has no
	// Go again, so nothing can follow it. Value = 2 + 1 + 8 = 11.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiOathBlueHocusRedMalefic(t *testing.T) {
	// Pitch Blue Hocus Pocus (3 res). Play Red Malefic (3 dmg, go again). Play Red Oath (+1
	// Runechant, peeks ahead and sees the Blade swing = +3 bonus, +1 Viserai Runechant from prior
	// non-attack action = 5). Swing Reaping Blade (cost 1, 3 dmg). Value = 3 + 5 + 3 = 11.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.OathOfTheArknightRed{},
		runeblade.MaleficIncantationRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_RunicReapingPrefersAttackPitch(t *testing.T) {
	// Pitching the Blue Hocus Pocus (attack-typed, pitch 3) pays for Runic Reaping + Shrill AND
	// satisfies Runic Reaping's pitched-attack rider. Pitching the Blue Malefic Aura instead would
	// lose the rider. Blue Malefic (1 arcane + 1 Viserai runechant = 2) → Runic Reaping (3 + 1
	// rider + 1 Viserai runechant = 5) → Shrill (4 base + 3 aura-created bonus = 7). Value = 2 + 5
	// + 7 = 14.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.RunicReapingRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	if got.Value != 14 {
		t.Fatalf("want value 14, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestBest_ViseraiMauvrionGrantsGoAgainToShrill(t *testing.T) {
	// Pitch Blue Hocus Pocus (3 res). Play Blue Malefic (1 arcane, go again). Play Red Mauvrion
	// Skies (0 cost, go again; grants go-again to the next Runeblade attack action card = Shrill,
	// and emits 3 runechants). Play Red Shrill (cost 2, 4 base + 3 aura-created bonus = 7; chains
	// thanks to Mauvrion's grant). Swing Reaping Blade (cost 1, 3 dmg). Viserai fires +1 on
	// Mauvrion (prior Malefic is a non-attack action) and +1 on Shrill (priors include non-attack
	// actions). Value = 1 + 3 + 7 + 3 + 2 = 16.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MauvrionSkiesRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 16 {
		t.Fatalf("want value 16, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

func TestIsLegalOrder_MauvrionCantSaveShrillWhenRuneragerIsAhead(t *testing.T) {
	// Mauvrion's grant lands on the first matching Runeblade attack action card in CardsRemaining.
	// In the ordering Mauvrion → Runerager → Shrill → weapon, Runerager is that first match, so
	// Shrill never gets the grant. Shrill has no printed go-again, so the Shrill → weapon chain
	// must break — isLegalOrder rejects the ordering.
	order := []card.Card{
		runeblade.MauvrionSkiesRed{},
		runeblade.RuneragerSwarmRed{},
		runeblade.ShrillOfSkullformRed{},
		weapon.ReapingBlade{},
	}
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.playSequence(order, nil, nil); legal {
		t.Fatalf("ordering %v should be illegal (Shrill has no go-again and Mauvrion granted Runerager instead)",
			cardNames(order))
	}
}

func cardNames(cs []card.Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}

func TestBest_ViseraiMauvrionChainsShrillIntoRuneragerIntoWeapon(t *testing.T) {
	// Pitch Blue Hocus → Mauvrion → Shrill → Runerager → Reaping Blade. Value = 3 + 7 + 3 + 3 + 2
	// Viserai runechants = 18.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MauvrionSkiesRed{},
		runeblade.RuneragerSwarmRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 0, nil)
	if got.Value != 18 {
		t.Fatalf("want value 18, got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// grantAll is a test-only attacker that sets GrantedGoAgain=true on every PlayedCard remaining in
// CardsRemaining. Used with grantSpy to detect cross-permutation PlayedCard wrapper leakage.
type grantAll struct{}

func (grantAll) ID() card.ID            { return card.Invalid }
func (grantAll) Name() string           { return "grantAll" }
func (grantAll) Cost() int              { return 0 }
func (grantAll) Pitch() int              { return 0 }
func (grantAll) Attack() int            { return 0 }
func (grantAll) Defense() int           { return 0 }
func (grantAll) Types() card.TypeSet    { return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack) }
func (grantAll) GoAgain() bool          { return true }
func (grantAll) Play(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		pc.GrantedGoAgain = true
	}
	return 0
}

// grantSpy is a test-only attacker that, when it plays FIRST in a permutation, records whether
// any PlayedCard in CardsRemaining already has GrantedGoAgain=true. With per-permutation fresh
// wrappers, that should never happen (no prior card in this permutation has run yet). If wrappers
// leak across permutations, a grant applied by a previous permutation's grantAll will still be
// visible here — tripping the spy.
type grantSpy struct{ saw *bool }

func (grantSpy) ID() card.ID              { return card.Invalid }
func (grantSpy) Name() string             { return "grantSpy" }
func (grantSpy) Cost() int                { return 0 }
func (grantSpy) Pitch() int               { return 0 }
func (grantSpy) Attack() int              { return 0 }
func (grantSpy) Defense() int             { return 0 }
func (grantSpy) Types() card.TypeSet      { return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack) }
func (grantSpy) GoAgain() bool            { return true }
func (g grantSpy) Play(s *card.TurnState) int {
	if len(s.CardsPlayed) != 0 {
		return 0
	}
	for _, pc := range s.CardsRemaining {
		if pc.GrantedGoAgain {
			*g.saw = true
		}
	}
	return 0
}

func TestBestSequence_PlayedCardGrantsDontLeakAcrossPermutations(t *testing.T) {
	// The permutation loop in bestSequence must allocate fresh *PlayedCard wrappers per
	// permutation so a grant applied by one permutation's Play() can't bleed into a later
	// permutation's legality/effect checks.
	//
	// Setup: attackers = [grantAll, grantSpy, grantAll]. The permute order emits grantAll-first
	// permutations before the first grantSpy-first permutation. Each grantAll-first permutation
	// mutates GrantedGoAgain on the wrappers for the cards behind it. When grantSpy later plays
	// FIRST in its own permutation, its CardsRemaining contains the other two cards' wrappers —
	// which must be fresh (GrantedGoAgain=false), since no card has played yet in this permutation.
	// If the wrappers were reused across permutations the spy would see leaked grants and trip.
	var sawLeak bool
	attackers := []card.Card{grantAll{}, grantSpy{saw: &sawLeak}, grantAll{}}
	ctx := newSequenceContextForTest(stubHero{}, nil, nil, 1_000_000, 0, len(attackers))
	_, _, _ = ctx.bestSequence(attackers, nil, nil, nil)
	if sawLeak {
		t.Fatalf("PlayedCard wrapper state leaked across permutations: grantSpy saw a pre-existing GrantedGoAgain when playing first")
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must
	// cover costs.
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, nil)
	if got.Value != 6 {
		t.Fatalf("want value 6, got %d", got.Value)
	}
	var res, cost int
	for i, c := range h {
		switch got.BestLine[i].Role {
		case Pitch:
			res += c.Pitch()
		case Attack:
			cost += c.Cost()
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
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
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
// plain-block with Malefic. Value = 3 (Red attack) + 2 (Malefic block) = 5. The OLD single-pool
// model would have let one pitch fund both phases and scored 7 (Red attack + 4 Toughen Up
// prevention) — an illegal split this test locks out.
func TestBest_AttackPitchCantCoverDefense(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}, generic.ToughenUpBlue{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
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
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.Value != 7 {
		t.Fatalf("Value = %d, want 7 (two pitched cards let attack + defense phases both pay; Roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_EmptyArsenalClaimsHeldCard confirms the post-hoc Arsenal promotion fires when the
// slot is empty and the winning partition has Held cards. A hand that can't play Toughen Up as
// DR (no other card to pitch for the 2-cost) leaves the DR Held; with arsenalCardIn=nil the
// slot is empty so the DR becomes Arsenal and rides into next turn as Play.ArsenalCard.
func TestBest_EmptyArsenalClaimsHeldCard(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	if got.BestLine[0].Role != Arsenal {
		t.Errorf("Roles[0] = %s, want ARSENAL", got.BestLine[0].Role)
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue", got.ArsenalCard)
	}
}

// TestBest_ArsenalInPlayDR covers the "arsenal card played as DR" branch. Previous turn left a
// Toughen Up Blue in arsenal; this turn we draw a Blue Malefic (pitch 3, cost 0). The pitched
// Malefic funds Toughen Up's 2-cost defense out of the arsenal, preventing 4 damage. Value = 4.
// Play.ArsenalCard is nil because the slot was vacated and no hand card ends up Held.
func TestBest_ArsenalInPlayDR(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, generic.ToughenUpBlue{})
	if got.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Malefic pitches to pay arsenal DR, prevents 4). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (slot was vacated, no Held card to promote)", got.ArsenalCard)
	}
	// ArsenalIn surfaces the arsenal-in assignment so callers (the best-hand printout) can flag
	// that this card wasn't in hand this turn.
	ai, hasArsenal := got.ArsenalIn()
	if !hasArsenal || ai.Card.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalIn = %v, want Toughen Up Blue", ai)
	}
	if ai.Role != Defend {
		t.Errorf("ArsenalIn role = %s, want DEFEND", ai.Role)
	}
}

// TestBest_ArsenalInStayBlocksNewArsenal locks in that while the arsenal slot is occupied, a
// hand card that would otherwise be promoted to Arsenal (because it's Held) stays Held instead —
// one arsenal slot, no replacement until the old card is played. A lone DR in hand is Held;
// the arsenal-in Toughen Up Blue stays (incoming=0 makes defending pointless, and the hand
// can't fund a DR anyway); post-hoc the slot is occupied so no promotion happens.
func TestBest_ArsenalInStayBlocksNewArsenal(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, generic.ToughenUpBlue{})
	if got.BestLine[0].Role != Held {
		t.Errorf("Roles[0] = %s, want HELD (slot occupied by arsenal-in, can't promote)", got.BestLine[0].Role)
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue (the staying arsenal-in card)", got.ArsenalCard)
	}
}

// TestBest_ArsenalInPlayAttack covers the "arsenal card played as attack" branch. A Red attack
// sits in arsenal from a previous turn; this turn we draw a single Red Attack which pitches
// (pitch 1) to fund both the hand Red's 1-cost and the arsenal Red's 1-cost... wait, one pitch
// can't pay two costs. Instead, the winning line plays the arsenal Red (funded by pitching the
// hand Red) and leaves the hand slot consumed. Value = 3 (arsenal Red's attack). With the
// arsenal slot now empty and no Held cards, ArsenalCard is nil.
func TestBest_ArsenalInPlayAttack(t *testing.T) {
	h := []card.Card{fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, fake.RedAttack{})
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (arsenal Red played, hand Red pitched to fund it). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (slot vacated, no Held to promote)", got.ArsenalCard)
	}
}

// TestBest_ArsenalInNonAttackActionPlays covers the "arsenal card isn't tagged Attack but can
// still be played on your turn" rule — non-attack actions (auras, item cards, etc.) are playable
// from arsenal. Hand: Malefic Incantation Red (cost 0, pitch 1). Arsenal: Blessing of Occult Red
// (cost 1, pitch 1, attack 0, Play returns 3 via DelayRunechants). The winning line pitches the
// Malefic to fund Blessing's 1-cost, plays Blessing from arsenal, and accrues 3 delayed
// runechants for next turn. Value = 3 (Blessing's Play return is counted as damage credit for
// the chain); LeftoverRunechants reflects the 3 tokens carrying over.
func TestBest_ArsenalInNonAttackActionPlays(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationRed{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, runeblade.BlessingOfOccultRed{})
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Malefic pitched, arsenal Blessing played for 3 runechants). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.LeftoverRunechants != 3 {
		t.Errorf("LeftoverRunechants = %d, want 3 (Blessing's Play delayed 3 tokens to next turn)",
			got.LeftoverRunechants)
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (Blessing played out of arsenal)", got.ArsenalCard)
	}
}
