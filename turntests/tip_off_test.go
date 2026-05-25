package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Tip-Off's Instant-discard mode wins the partition when its mark sets up a
// marked-defender rider on a future turn. Turn 1 hand has Tip-Off Red alone (no resources
// to fund mode 0 at cost 1); mode 1 costs 0, plays as a 0-damage attack, queues an
// end-of-turn mark, and refunds the AP via Go Again. Turn 2 draws Outed Red from the
// one-card deck and sees the carried mark for +1{p} (3 printed + 1 = 4).
func TestTipOff_InstantModeMarksOpponentForNextTurnOuted(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, []deck.Card{cards.OutedRed{}})
	hand1 := []card.Card{cards.TipOffRed{}}

	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, nil, hand1)

	if turn1.Value != 0 {
		t.Errorf("turn 1 Value = %d, want 0 (mode 1 deals 0 damage)\nBestLine: %s",
			turn1.Value, formatBestLine(turn1.BestLine))
	}
	if !bestLineHasRole(turn1.BestLine, cards.TipOffRed{}, card.Attack) {
		t.Errorf("turn 1 BestLine missing Tip-Off as Attack: %s", formatBestLine(turn1.BestLine))
	}
	if !turn1.State.OpponentMarked() {
		t.Errorf("turn 1 end state OpponentMarked = false, want true (Tip-Off mode 1's end-of-turn mark must land and carry to next turn)")
	}
	if turn2.Value != 4 {
		t.Errorf("turn 2 Value = %d, want 4 (Outed 3{p} + 1 marked-defender bonus)\nBestLine: %s",
			turn2.Value, formatBestLine(turn2.BestLine))
	}
	if !bestLineHasRole(turn2.BestLine, cards.OutedRed{}, card.Attack) {
		t.Errorf("turn 2 BestLine missing Outed as Attack: %s", formatBestLine(turn2.BestLine))
	}
}
