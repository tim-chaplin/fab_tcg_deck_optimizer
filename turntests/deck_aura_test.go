package turntests

// End-to-end aura-trigger coverage driven through EvalOneTurnForTesting: pin that
// start-of-turn auras queued during turn 1 fire correctly at the turn-2 boundary —
// damage credit, graveyard exhaustion, hand reveals, and OncePerTurn re-arm.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Sigil of Fyendal plays turn 1 and its start-of-turn trigger fires on turn 2 —
// crediting 1 damage-equivalent and landing the sigil in graveyard.
func TestEvalOneTurn_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := cards.SigilOfFyendalBlue{}
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{sigil})

	sigilPlayed := false
	for _, a := range extras.BestLine {
		if a.Card.ID() == ids.SigilOfFyendalBlue && a.Role == deck.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", extras.BestLine)
	}
	if extras.TriggerDamage != 1 {
		t.Errorf("TriggerDamage = %d, want 1 (Fyendal's 1{h} gain fires at start of turn 2)",
			extras.TriggerDamage)
	}
	if len(extras.TriggerGraveyard) != 1 || extras.TriggerGraveyard[0].ID() != ids.SigilOfFyendalBlue {
		t.Errorf("TriggerGraveyard = %v, want [Sigil of Fyendal] (Count hit zero after firing)",
			extras.TriggerGraveyard)
	}
}

// Tests end-to-end Sigil of the Arknight: plays turn 1, on turn 2 reveals an attack action
// off the deck top into hand alongside the normal refill.
func TestEvalOneTurn_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := cards.SigilOfTheArknightBlue{}
	reveal := cards.AetherSlashRed{}
	// Deck layout: positions 0..3 are turn 2's normal refill (Blues), position 4 is the reveal
	// target at the post-draw top, positions 5+ are unused filler.
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{sigil})

	sigilPlayed := false
	for _, a := range extras.BestLine {
		if a.Card.ID() == ids.SigilOfTheArknightBlue && a.Role == deck.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", extras.BestLine)
	}

	// Turn 2: 4 normal draws + 1 revealed = 5 cards. deckCards[0..3] refill turn 2's hand;
	// deckCards[4] is the reveal target appended at the tail.
	wantHand := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
	}
	if len(gs.Hand()) != len(wantHand) {
		t.Fatalf("turn 2 hand size = %d, want %d (4 normal draws + 1 revealed)", len(gs.Hand()), len(wantHand))
	}
	for i, want := range wantHand {
		if gs.Hand()[i] != want {
			t.Errorf("turn 2 hand[%d] = %v, want %v", i, gs.Hand()[i], want)
		}
	}
	if len(extras.TriggerGraveyard) != 1 || extras.TriggerGraveyard[0].ID() != ids.SigilOfTheArknightBlue {
		t.Errorf("TriggerGraveyard = %v, want [Sigil of the Arknight]", extras.TriggerGraveyard)
	}
}

// Tests Blessing of Occult queueing a 3-rune start-of-turn trigger that fires on turn 2:
// turn-1 play credits 0 (deferred), turn 2's carryover holds 3 Runechants and the Blessing
// exhausts to graveyard.
func TestEvalOneTurn_BlessingOfOccultCreatesRunesAtStartOfNextTurn(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{blessing, pitch})

	if extras.Value != 0 {
		t.Errorf("Value = %d, want 0 (Blessing's rune credit is deferred)", extras.Value)
	}
	blessingPlayed := false
	for _, a := range extras.BestLine {
		if a.Card.ID() == ids.BlessingOfOccultRed && a.Role == deck.Attack {
			blessingPlayed = true
			break
		}
	}
	if !blessingPlayed {
		t.Errorf("turn 1 BestLine didn't play Blessing as Role=Attack: %+v", extras.BestLine)
	}
	if gs.RunechantCount() != 3 {
		t.Errorf("Runechants = %d, want 3 (Blessing's start-of-turn trigger creates 3 tokens)",
			gs.RunechantCount())
	}
	if extras.TriggerDamage != 3 {
		t.Errorf("TriggerDamage = %d, want 3", extras.TriggerDamage)
	}
	if len(extras.TriggerGraveyard) != 1 || extras.TriggerGraveyard[0].ID() != ids.BlessingOfOccultRed {
		t.Errorf("TriggerGraveyard = %v, want [Blessing [R]]", extras.TriggerGraveyard)
	}
}

// Tests that Malefic's OncePerTurn AttackAction trigger fires exactly once on the lone
// attack action and survives into next turn with Count decremented.
func TestEvalOneTurn_MaleficIncantationOncePerTurnLimitsToOneRune(t *testing.T) {
	malefic := cards.MaleficIncantationRed{}
	hocus := cards.HocusPocusRed{}
	// Filler deck so turn 2 can be dealt — content doesn't matter for what we assert.
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{malefic, hocus})

	maleficPlayed, hocusPlayed := false, false
	for _, a := range extras.BestLine {
		if a.Card.ID() == ids.MaleficIncantationRed && a.Role == deck.Attack {
			maleficPlayed = true
		}
		if a.Card.ID() == ids.HocusPocusRed && a.Role == deck.Attack {
			hocusPlayed = true
		}
	}
	if !maleficPlayed {
		t.Errorf("turn 1 BestLine didn't play Malefic as Role=Attack: %+v", extras.BestLine)
	}
	if !hocusPlayed {
		t.Errorf("turn 1 BestLine didn't play Hocus Pocus as Role=Attack: %+v", extras.BestLine)
	}
	if extras.Value != 6 {
		t.Errorf("Value = %d, want 6 (3 Hocus + 1 Hocus rune + 1 Viserai trigger + 1 Malefic trigger)",
			extras.Value)
	}
	// Malefic's AttackAction trigger doesn't fire at start of turn — it only ticks on
	// attack actions during the chain. Carry-only at the turn boundary.
	if extras.TriggerDamage != 0 {
		t.Errorf("TriggerDamage = %d, want 0 (Malefic only fires on attack actions)",
			extras.TriggerDamage)
	}
	if len(extras.TriggerGraveyard) != 0 {
		t.Errorf("TriggerGraveyard = %v, want empty (Malefic still has Count>0)",
			extras.TriggerGraveyard)
	}
}

// Tests that Runeblood Incantation's start-of-turn trigger fires once per turn and survives
// while Count > 0.
func TestEvalOneTurn_RunebloodIncantationTicksAcrossTurns(t *testing.T) {
	runeblood := cards.RunebloodIncantationRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{runeblood, pitch})

	runebloodPlayed := false
	for _, a := range extras.BestLine {
		if a.Card.ID() == ids.RunebloodIncantationRed && a.Role == deck.Attack {
			runebloodPlayed = true
			break
		}
	}
	if !runebloodPlayed {
		t.Errorf("turn 1 BestLine didn't play Runeblood as Role=Attack: %+v", extras.BestLine)
	}
	if extras.Value != 0 {
		t.Errorf("Value = %d, want 0 (every Runeblood rune is deferred to a future fire)",
			extras.Value)
	}
	if extras.TriggerDamage != 1 {
		t.Errorf("TriggerDamage = %d, want 1 (one tick per turn)", extras.TriggerDamage)
	}
	if gs.RunechantCount() != 1 {
		t.Errorf("Runechants = %d, want 1 (one rune per fire)", gs.RunechantCount())
	}
	if len(extras.TriggerGraveyard) != 0 {
		t.Errorf("TriggerGraveyard = %v, want empty (Red has Count=3, only one tick fired)",
			extras.TriggerGraveyard)
	}
}
