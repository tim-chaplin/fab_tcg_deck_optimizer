package turntests

// End-to-end aura-trigger coverage driven through EvalOneTurnForTesting: pin that
// start-of-turn auras queued during turn 1 fire correctly at the turn-2 boundary —
// damage credit, graveyard exhaustion, hand reveals, and OncePerTurn re-arm.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Sigil of Fyendal plays turn 1 and its start-of-turn trigger fires on turn 2 —
// crediting 1 damage-equivalent and landing the sigil in graveyard.
func TestEvalOneTurn_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := cards.SigilOfFyendalBlue{}
	deckCards := []sim.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, []sim.Card{sigil})

	sigilPlayed := false
	for _, a := range state.BestLine {
		if a.Card.ID() == ids.SigilOfFyendalBlue && a.Role == sim.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", state.BestLine)
	}
	if state.StartOfNextTurnTriggerDamage != 1 {
		t.Errorf("StartOfNextTurnTriggerDamage = %d, want 1 (Fyendal's 1{h} gain fires at start of turn 2)",
			state.StartOfNextTurnTriggerDamage)
	}
	if len(state.StartOfNextTurnGraveyard) != 1 || state.StartOfNextTurnGraveyard[0].ID() != ids.SigilOfFyendalBlue {
		t.Errorf("StartOfNextTurnGraveyard = %v, want [Sigil of Fyendal] (Count hit zero after firing)",
			state.StartOfNextTurnGraveyard)
	}
}

// Tests end-to-end Sigil of the Arknight: plays turn 1, on turn 2 reveals an attack action
// off the deck top into hand alongside the normal refill.
func TestEvalOneTurn_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := cards.SigilOfTheArknightBlue{}
	reveal := cards.AetherSlashRed{}
	// Deck layout: positions 0..3 are turn 2's normal refill (Blues), position 4 is the reveal
	// target at the post-draw top, positions 5+ are unused filler.
	deckCards := []sim.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
		testutils.BlueAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, []sim.Card{sigil})

	sigilPlayed := false
	for _, a := range state.BestLine {
		if a.Card.ID() == ids.SigilOfTheArknightBlue && a.Role == sim.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", state.BestLine)
	}

	// Turn 2: 4 normal draws + 1 revealed = 5 cards. deckCards[0..3] refill turn 2's hand;
	// deckCards[4] is the reveal target appended at the tail.
	wantHand := []sim.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
	}
	if len(state.StartOfNextTurnHand) != len(wantHand) {
		t.Fatalf("turn 2 hand size = %d, want %d (4 normal draws + 1 revealed)", len(state.StartOfNextTurnHand), len(wantHand))
	}
	for i, want := range wantHand {
		if state.StartOfNextTurnHand[i] != want {
			t.Errorf("turn 2 hand[%d] = %v, want %v", i, state.StartOfNextTurnHand[i], want)
		}
	}
	if len(state.StartOfNextTurnGraveyard) != 1 || state.StartOfNextTurnGraveyard[0].ID() != ids.SigilOfTheArknightBlue {
		t.Errorf("StartOfNextTurnGraveyard = %v, want [Sigil of the Arknight]", state.StartOfNextTurnGraveyard)
	}
}

// Tests Blessing of Occult queueing a 3-rune start-of-turn trigger that fires on turn 2:
// turn-1 play credits 0 (deferred), turn 2's carryover holds 3 Runechants and the Blessing
// exhausts to graveyard.
func TestEvalOneTurn_BlessingOfOccultCreatesRunesAtStartOfNextTurn(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []sim.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, []sim.Card{blessing, pitch})

	if state.Value != 0 {
		t.Errorf("Value = %d, want 0 (Blessing's rune credit is deferred)", state.Value)
	}
	blessingPlayed := false
	for _, a := range state.BestLine {
		if a.Card.ID() == ids.BlessingOfOccultRed && a.Role == sim.Attack {
			blessingPlayed = true
			break
		}
	}
	if !blessingPlayed {
		t.Errorf("turn 1 BestLine didn't play Blessing as Role=Attack: %+v", state.BestLine)
	}
	if state.Runechants() != 3 {
		t.Errorf("Runechants = %d, want 3 (Blessing's start-of-turn trigger creates 3 tokens)",
			state.Runechants())
	}
	if state.StartOfNextTurnTriggerDamage != 3 {
		t.Errorf("StartOfNextTurnTriggerDamage = %d, want 3", state.StartOfNextTurnTriggerDamage)
	}
	if len(state.StartOfNextTurnGraveyard) != 1 || state.StartOfNextTurnGraveyard[0].ID() != ids.BlessingOfOccultRed {
		t.Errorf("StartOfNextTurnGraveyard = %v, want [Blessing [R]]", state.StartOfNextTurnGraveyard)
	}
}

// Tests that Malefic's OncePerTurn AttackAction trigger fires exactly once on the lone
// attack action and survives into next turn with Count decremented.
func TestEvalOneTurn_MaleficIncantationOncePerTurnLimitsToOneRune(t *testing.T) {
	malefic := cards.MaleficIncantationRed{}
	hocus := cards.HocusPocusRed{}
	// Filler deck so turn 2 can be dealt — content doesn't matter for what we assert.
	deckCards := []sim.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, []sim.Card{malefic, hocus})

	maleficPlayed, hocusPlayed := false, false
	for _, a := range state.BestLine {
		if a.Card.ID() == ids.MaleficIncantationRed && a.Role == sim.Attack {
			maleficPlayed = true
		}
		if a.Card.ID() == ids.HocusPocusRed && a.Role == sim.Attack {
			hocusPlayed = true
		}
	}
	if !maleficPlayed {
		t.Errorf("turn 1 BestLine didn't play Malefic as Role=Attack: %+v", state.BestLine)
	}
	if !hocusPlayed {
		t.Errorf("turn 1 BestLine didn't play Hocus Pocus as Role=Attack: %+v", state.BestLine)
	}
	if state.Value != 6 {
		t.Errorf("Value = %d, want 6 (3 Hocus + 1 Hocus rune + 1 Viserai trigger + 1 Malefic trigger)",
			state.Value)
	}
	// Malefic's AttackAction trigger doesn't fire at start of turn — it only ticks on
	// attack actions during the chain. Carry-only at the turn boundary.
	if state.StartOfNextTurnTriggerDamage != 0 {
		t.Errorf("StartOfNextTurnTriggerDamage = %d, want 0 (Malefic only fires on attack actions)",
			state.StartOfNextTurnTriggerDamage)
	}
	if len(state.StartOfNextTurnGraveyard) != 0 {
		t.Errorf("StartOfNextTurnGraveyard = %v, want empty (Malefic still has Count>0)",
			state.StartOfNextTurnGraveyard)
	}
}

// Tests that Runeblood Incantation's start-of-turn trigger fires once per turn and survives
// while Count > 0.
func TestEvalOneTurn_RunebloodIncantationTicksAcrossTurns(t *testing.T) {
	runeblood := cards.RunebloodIncantationRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []sim.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, []sim.Card{runeblood, pitch})

	runebloodPlayed := false
	for _, a := range state.BestLine {
		if a.Card.ID() == ids.RunebloodIncantationRed && a.Role == sim.Attack {
			runebloodPlayed = true
			break
		}
	}
	if !runebloodPlayed {
		t.Errorf("turn 1 BestLine didn't play Runeblood as Role=Attack: %+v", state.BestLine)
	}
	if state.Value != 0 {
		t.Errorf("Value = %d, want 0 (every Runeblood rune is deferred to a future fire)",
			state.Value)
	}
	if state.StartOfNextTurnTriggerDamage != 1 {
		t.Errorf("StartOfNextTurnTriggerDamage = %d, want 1 (one tick per turn)", state.StartOfNextTurnTriggerDamage)
	}
	if state.Runechants() != 1 {
		t.Errorf("Runechants = %d, want 1 (one rune per fire)", state.Runechants())
	}
	if len(state.StartOfNextTurnGraveyard) != 0 {
		t.Errorf("StartOfNextTurnGraveyard = %v, want empty (Red has Count=3, only one tick fired)",
			state.StartOfNextTurnGraveyard)
	}
}
