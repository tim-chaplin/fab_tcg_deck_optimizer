package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Back Alley Breakline's OnFaceUp hook grants +1 AP and that the extra AP
// enables an additional non-Go-Again attack step that otherwise wouldn't fit.
//
// Setup:
//   - Arsenal: Back Alley Breakline Red (cost 1, attack 5)
//   - Hand: a fake red attack whose Play scans CardsRemaining for BAB and calls
//     TurnFaceUp on that exact CardState pointer. 4 base power, Go Again so it refunds
//     its own AP. Stands in for an "effect that puts BAB face up in a zone."
//   - Hand: a second fake red attack — 3 power, cost 0, no Go Again. Needs its own AP
//     to play.
//   - Hand: BluePitch (pitch 3) to fund BAB's 1{r} cost when it plays from arsenal.
//
// Without OnFaceUp granting +1 AP, the attack turn has 1 AP total: the Go-Again caller
// refunds its AP, BAB consumes the only remaining AP, and the second attack is stuck
// in hand. Optimal there is caller (4) + BAB-from-arsenal (5) = 9.
//
// With OnFaceUp granting +1 AP, the caller's Play turns BAB face up, BAB's hook adds 1
// AP, and the second attack plays for 3 more. Optimal is caller (4) + BAB (5) +
// second attack (3) = 12.
func TestBackAlleyBreakline_OnFaceUpGrantsAPForExtraAttackStep(t *testing.T) {
	bab := cards.BackAlleyBreaklineRed{}
	caller := testutils.FakeRedAttack().
		WithName("FakeFaceUpCaller").
		WithPower(4).
		WithGoAgain().
		WithPlay(func(ge card.GameEngine, _ card.Logger, _ *card.CardState) {
			for _, pc := range ge.CardsRemaining() {
				if pc.Card == bab {
					ge.TurnFaceUp(pc)
					return
				}
			}
		})
	extraAttack := testutils.FakeRedAttack().
		WithName("FakeExtraAttack").
		WithPower(3)

	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{caller, extraAttack, testutils.FakeBlueResource()}
	initial := gameengine.GameStateBuilder().SetArsenal(bab).Build()

	summary := sim.EvalOneTurnForTesting(d, initial, hand)

	if summary.Value != 12 {
		t.Errorf("Value = %d, want 12 (caller 4 + BAB 5 + extra 3 — the +1 AP from OnFaceUp must enable the third attack step)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if !bestLineHasRole(summary.BestLine, bab, card.Attack) {
		t.Errorf("BestLine missing BAB as Attack (from arsenal): %s", formatBestLine(summary.BestLine))
	}
	if !bestLineHasRole(summary.BestLine, caller, card.Attack) {
		t.Errorf("BestLine missing caller as Attack: %s", formatBestLine(summary.BestLine))
	}
	if !bestLineHasRole(summary.BestLine, extraAttack, card.Attack) {
		t.Errorf("BestLine missing extra attack as Attack: %s", formatBestLine(summary.BestLine))
	}
}
