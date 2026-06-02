package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

var _ card.Legendary = cards.ImperialSealOfCommandRed{}

// Tests that a Royal hero activates Imperial Seal — the hit destroys the opponent's arsenal and
// the "Destroy this" cost leaves no item in play.
func TestImperialSeal_RoyalAbilityDestroysArsenalOnHit(t *testing.T) {
	// Power-4 lands the LikelyDamageDealt=1 leak path, firing OnHit.
	attacker := testutils.FakeRedAttack().
		WithPower(4).
		WithName("Attacker")
	hero := testutils.Hero{Intel: 4, TypeSet: card.NewTypeSet(card.TypeRoyal)}
	state := gameengine.GameStateBuilder().
		SetHero(hero).
		CreateItemFromCard(cards.ImperialSealOfCommandRed{}).
		Build()
	d := deck.New(hero, nil, nil)
	hand := []card.Card{attacker}
	summary := sim.EvalOneTurnForTesting(d, state, hand)

	if !summary.State.DestroyedOpponentArsenal() {
		t.Errorf("DestroyedOpponentArsenal still false\nValue=%d BestLine: %s",
			summary.Value, sim.FormatBestLine(summary.BestLine))
	}
	// 4 power attack hit + 3 OpponentDiscard credit for the destroyed arsenal.
	if got, want := summary.Value, 7; got != want {
		t.Errorf("Value = %d, want %d\nBestLine: %s",
			got, want, sim.FormatBestLine(summary.BestLine))
	}
	// "Destroy this" is the activation cost: the item must be gone after the ability resolves.
	if got := len(summary.State.Items()); got != 0 {
		t.Errorf("Items() = %d after the ability resolved, want 0 — the item should self-destruct\n"+
			"BestLine: %s", got, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests that a non-Royal hero skips the OnHit registration so the attack's hit doesn't
// destroy the arsenal, even though the ability is still enumerated and free to activate.
func TestImperialSeal_NonRoyalHeroNoArsenalDestroy(t *testing.T) {
	attacker := testutils.FakeRedAttack().
		WithPower(4).
		WithName("Attacker")
	hero := testutils.Hero{Intel: 4} // no Royal type
	state := gameengine.GameStateBuilder().
		SetHero(hero).
		CreateItemFromCard(cards.ImperialSealOfCommandRed{}).
		Build()
	d := deck.New(hero, nil, nil)
	hand := []card.Card{attacker}
	summary := sim.EvalOneTurnForTesting(d, state, hand)

	if summary.State.DestroyedOpponentArsenal() {
		t.Errorf("DestroyedOpponentArsenal flipped without a Royal hero\nBestLine: %s",
			sim.FormatBestLine(summary.BestLine))
	}
	// 4 power attack hit, no arsenal credit (Royal gate blocks the trigger registration).
	if got, want := summary.Value, 4; got != want {
		t.Errorf("Value = %d, want %d\nBestLine: %s",
			got, want, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests the full two-turn line: turn 1 plays Imperial Seal from hand to create the in-play
// item, turn 2 draws an attacker, the wmask enumerates the ability, and the attack hit
// destroys the opponent's arsenal — the one-turn-delay contract end-to-end.
func TestImperialSeal_PlayThenAbilityDestroysArsenalNextTurn(t *testing.T) {
	hero := testutils.Hero{Intel: 4, TypeSet: card.NewTypeSet(card.TypeRoyal)}
	attacker := testutils.FakeRedAttack().
		WithPower(4).
		WithName("Attacker")
	state := gameengine.GameStateBuilder().
		SetHero(hero).
		Build()
	d := deck.New(hero, nil, []deck.Card{attacker})
	hand1 := []card.Card{cards.ImperialSealOfCommandRed{}}

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, state, hand1)

	foundItem := false
	for _, it := range turn1.State.Items() {
		if it.CardID() == (cards.ImperialSealOfCommandRed{}).ID() {
			foundItem = true
			break
		}
	}
	if !foundItem {
		t.Fatalf("Imperial Seal not in turn1 Items() — parent Play should have created it\n"+
			"BestLine: %s", sim.FormatBestLine(turn1.BestLine))
	}
	if turn1.State.DestroyedOpponentArsenal() {
		t.Errorf("DestroyedOpponentArsenal flipped on turn 1; same-turn activation isn't "+
			"supported\nBestLine: %s", sim.FormatBestLine(turn1.BestLine))
	}
	if !turn2.State.DestroyedOpponentArsenal() {
		t.Errorf("DestroyedOpponentArsenal still false on turn 2\nBestLine: %s",
			sim.FormatBestLine(turn2.BestLine))
	}
	// The turn-2 activation pays "Destroy this", so the item is gone by end of turn 2.
	if got := len(turn2.State.Items()); got != 0 {
		t.Errorf("turn2 Items() = %d after the ability resolved, want 0 — the item should "+
			"self-destruct\nBestLine: %s", got, sim.FormatBestLine(turn2.BestLine))
	}
	// Turn 1 plays Imperial Seal for nothing — item creation credits no immediate value.
	if got, want := turn1.Value, 0; got != want {
		t.Errorf("turn1.Value = %d, want %d\nBestLine: %s",
			got, want, sim.FormatBestLine(turn1.BestLine))
	}
	// Turn 2: 4 power attack hit + 3 OpponentDiscard credit for the destroyed arsenal.
	if got, want := turn2.Value, 7; got != want {
		t.Errorf("turn2.Value = %d, want %d\nBestLine: %s",
			got, want, sim.FormatBestLine(turn2.BestLine))
	}
}
