package turntests

// Multi-turn tests that pin state which must survive the turn boundary: equipped weapons,
// graveyard contents, and destroyed-aura source cards.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// TestEvalTwoTurns_EquippedWeaponPersistsAcrossTurns pins that an equipped weapon stays
// swingable on turn 2. Each turn pitches a Blue (cost 3) to swing Nebula Blade for
// 1 base + 1 on-hit Runechant = 2. Expected turn1.Value == turn2.Value == 2; if the
// weapon is lost across the turn boundary, turn2.Value drops to 0.
func TestEvalTwoTurns_EquippedWeaponPersistsAcrossTurns(t *testing.T) {
	d := deck.New(heroes.Viserai, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{
		testutils.FakeBlueResource(),
	})
	hand := []card.Card{testutils.FakeBlueResource()}

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, nil, hand)

	if turn1.Value != 2 {
		t.Errorf("turn 1 Value = %d, want 2\nBestLine: %s", turn1.Value, formatBestLine(turn1.BestLine))
	}
	if !bestLineHasRole(turn1.BestLine, testutils.FakeBlueResource(), card.Pitch) {
		t.Errorf("turn 1 BestLine missing BluePitch as Pitch: %s", formatBestLine(turn1.BestLine))
	}
	if turn2.Value != 2 {
		t.Errorf("turn 2 Value = %d, want 2 (weapon should still be equipped)\nBestLine: %s",
			turn2.Value, formatBestLine(turn2.BestLine))
	}
	if !bestLineHasRole(turn2.BestLine, testutils.FakeBlueResource(), card.Pitch) {
		t.Errorf("turn 2 BestLine missing BluePitch as Pitch: %s", formatBestLine(turn2.BestLine))
	}
}

// TestEvalTwoTurns_GraveyardPersistsAcrossTurns pins that a card in the graveyard at the
// start of turn 1 is still there for turn 2. Initial graveyard = [Sigil of Deadwood];
// turn 1 has an empty hand; turn 2 plays Sigil of Silphidae, whose on-enter banishes
// another aura from the graveyard for 1 arcane. Expected turn1.Value == 0, turn2.Value == 1;
// if the graveyard is wiped between turns, the banish finds nothing and turn2.Value == 0.
func TestEvalTwoTurns_GraveyardPersistsAcrossTurns(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, []deck.Card{
		cards.SigilOfSilphidaeBlue{},
	})
	initial := gameengine.GameStateBuilder().
		SetGraveyard([]card.Card{cards.SigilOfDeadwoodBlue{}}).
		Build()

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, initial, nil)

	if turn1.Value != 0 {
		t.Errorf("turn 1 Value = %d, want 0 (empty hand)\nBestLine: %s",
			turn1.Value, formatBestLine(turn1.BestLine))
	}
	if len(turn1.BestLine) != 0 {
		t.Errorf("turn 1 BestLine = %s, want empty", formatBestLine(turn1.BestLine))
	}
	if turn2.Value != 1 {
		t.Errorf("turn 2 Value = %d, want 1\nBestLine: %s", turn2.Value, formatBestLine(turn2.BestLine))
	}
	if !bestLineHasRole(turn2.BestLine, cards.SigilOfSilphidaeBlue{}, card.Attack) {
		t.Errorf("turn 2 BestLine missing Sigil of Silphidae as Attack: %s",
			formatBestLine(turn2.BestLine))
	}
}

// TestEvalTwoTurns_DestroyedAuraSourceReachesGraveyard pins the FaB rule that a destroyed
// aura's source card lands in the graveyard. Initial state has a Sigil of Fyendal aura in
// play; turn 1's hand is empty, so the only event is the start-of-turn aura tick (+1) which
// destroys the aura and deposits Fyendal in the graveyard. Turn 2 plays Sigil of Silphidae,
// whose on-enter banishes Fyendal from the graveyard for 1 arcane. Expected turn1.Value ==
// turn2.Value == 1; if destroyed aura sources don't reach the graveyard, turn2.Value == 0.
func TestEvalTwoTurns_DestroyedAuraSourceReachesGraveyard(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, []deck.Card{
		cards.SigilOfSilphidaeBlue{},
	})
	initial := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.SigilOfFyendalBlue{}).
		Build()

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, initial, nil)

	if turn1.Value != 1 {
		t.Errorf("turn 1 Value = %d, want 1 (Fyendal aura tick)\nBestLine: %s",
			turn1.Value, formatBestLine(turn1.BestLine))
	}
	if len(turn1.BestLine) != 0 {
		t.Errorf("turn 1 BestLine = %s, want empty (value is from the aura tick, not a chain play)",
			formatBestLine(turn1.BestLine))
	}
	if turn2.Value != 1 {
		t.Errorf("turn 2 Value = %d, want 1 (destroyed Fyendal must persist in graveyard for Silphidae to banish)\nBestLine: %s",
			turn2.Value, formatBestLine(turn2.BestLine))
	}
	if !bestLineHasRole(turn2.BestLine, cards.SigilOfSilphidaeBlue{}, card.Attack) {
		t.Errorf("turn 2 BestLine missing Sigil of Silphidae as Attack: %s",
			formatBestLine(turn2.BestLine))
	}
}

// Tests that an OpponentMarked flag carried over from turn 1 is still set when turn 2's
// marked-defender riders run.
func TestEvalTwoTurns_OpponentMarkedPersistsAcrossTurns(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, []deck.Card{
		cards.OutedRed{},
	})
	initial := gameengine.GameStateBuilder().SetOpponentMarked(true).Build()

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, initial, nil)

	if turn1.Value != 0 {
		t.Errorf("turn 1 Value = %d, want 0 (empty hand)\nBestLine: %s",
			turn1.Value, formatBestLine(turn1.BestLine))
	}
	if turn2.Value != 4 {
		t.Errorf("turn 2 Value = %d, want 4 (Outed 3{p} + 1 marked-defender bonus must see the carried mark)\nBestLine: %s",
			turn2.Value, formatBestLine(turn2.BestLine))
	}
	if !bestLineHasRole(turn2.BestLine, cards.OutedRed{}, card.Attack) {
		t.Errorf("turn 2 BestLine missing Outed as Attack: %s", formatBestLine(turn2.BestLine))
	}
}

// Tests that a token item (Gold) in play at the start of turn 1 is still on the items list
// at the start of turn 2 when neither turn can fund its ability cost.
func TestEvalTwoTurns_ItemTokenPersistsAcrossTurns(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	initial := gameengine.GameStateBuilder().AddItem(token.NewGold(1)).Build()

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, initial, nil)

	if turn1.State.GoldCount() != 1 {
		t.Errorf("turn 1 GoldCount = %d, want 1 (empty hand can't fund the {2} ability cost)",
			turn1.State.GoldCount())
	}
	if turn2.State.GoldCount() != 1 {
		t.Errorf("turn 2 GoldCount = %d, want 1 (items must persist across the turn boundary)",
			turn2.State.GoldCount())
	}
}
