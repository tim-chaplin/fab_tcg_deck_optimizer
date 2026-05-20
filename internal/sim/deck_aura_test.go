package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// damageTrigger returns a StartOfTurn Aura crediting the given damage and destroying
// itself on the first fire.
func damageTrigger(self card.Card, damage int, calls *int) gameengine.Aura {
	return aura.NewFromCard(
		self,
		triggertype.StartOfTurn,
		func(ge card.GameEngine, _ card.Logger, a card.Aura) {
			*calls++
			ge.AddValue(damage)
			a.Destroy(true)
		},
		1,
		false,
	)
}

// TestProcessAurasAtStartOfTurn_FiresEachQueuedTriggerOnce verifies every queued
// start-of-turn trigger's handler is invoked exactly once per pass.
func TestProcessAurasAtStartOfTurn_FiresEachQueuedTriggerOnce(t *testing.T) {
	src := testutils.RedAttack{}
	var callsA, callsB int
	queue := []gameengine.Aura{damageTrigger(src, 2, &callsA), damageTrigger(src, 3, &callsB)}
	survivors, total, _, _ := ProcessAurasAtStartOfTurnForTest(queue, DeckOf())
	if total != 5 {
		t.Errorf("total = %d, want 5 (2+3)", total)
	}
	if len(survivors) != 0 {
		t.Errorf("survivors = %+v, want empty (both triggers had Count=1)", survivors)
	}
	if callsA != 1 || callsB != 1 {
		t.Errorf("handler call counts = (%d, %d), want (1, 1)", callsA, callsB)
	}
}

// TestProcessAurasAtStartOfTurn_EmptyQueue short-circuits.
func TestProcessAurasAtStartOfTurn_EmptyQueue(t *testing.T) {
	survivors, total, revealed, _ := ProcessAurasAtStartOfTurnForTest(nil, DeckOf())
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if len(survivors) != 0 || len(revealed) != 0 {
		t.Errorf("non-empty outputs on empty input: survivors=%v revealed=%v",
			survivors, revealed)
	}
}

// TestProcessAurasAtStartOfTurn_GraveyardsExhaustedAura: a handler that calls Destroy
// lands Self in the graveyard before subsequent handlers run.
func TestProcessAurasAtStartOfTurn_GraveyardsExhaustedAura(t *testing.T) {
	src := testutils.RedAttack{}
	var seen []card.Card
	watcher := aura.NewFromCard(
		testutils.YellowAttack{},
		triggertype.StartOfTurn,
		func(ge card.GameEngine, _ card.Logger, _ card.Aura) {
			eng := ge.(*gameengine.GameEngine)
			seen = append([]card.Card(nil), eng.Graveyard()...)
		},
		1,
		false,
	)
	first := aura.NewFromCard(
		src,
		triggertype.StartOfTurn,
		func(_ card.GameEngine, _ card.Logger, a card.Aura) {
			a.Destroy(true)
		},
		1,
		false,
	)
	_, _, _, _ = ProcessAurasAtStartOfTurnForTest([]gameengine.Aura{first, watcher}, DeckOf())
	if len(seen) != 1 || seen[0] != src {
		t.Errorf("second handler saw Graveyard = %v, want [%v]", seen, src)
	}
}

// TestProcessAurasAtStartOfTurn_IgnoresPonder: Ponder fires at end-of-turn (not start).
func TestProcessAurasAtStartOfTurn_IgnoresPonder(t *testing.T) {
	survivors, _, _, _ := ProcessAurasAtStartOfTurnForTest([]gameengine.Aura{token.NewPonder(1)}, DeckOf())
	if len(survivors) != 1 || survivors[0].CardName() != "Ponder" {
		t.Errorf("survivors = %+v, want one Ponder aura intact", survivors)
	}
}

// TestFireEndOfTurn_PonderPopsDeckTopIntoHand.
func TestFireEndOfTurn_PonderPopsDeckTopIntoHand(t *testing.T) {
	a, b, c := testutils.NewStubCard("a"), testutils.NewStubCard("b"), testutils.NewStubCard("c")
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b, c}).Build()}
	ge.CreateAura(token.NewPonder(2))
	FireEndOfTurn(ge)

	h := ge.Hand()
	if len(h) != 2 || h[0] != a || h[1] != b {
		t.Errorf("Hand = %v, want [a, b]", h)
	}
	if got := ge.Deck().Size(); got != 1 {
		t.Errorf("Deck size = %d, want 1 (two cards popped)", got)
	}
	if len(ge.Auras()) != 0 {
		t.Errorf("Auras = %+v, want empty (Ponder destroyed itself)", ge.Auras())
	}
}

// TestFireEndOfTurn_PonderEmptyDeckIsNoOp.
func TestFireEndOfTurn_PonderEmptyDeckIsNoOp(t *testing.T) {
	ge := gameengine.New()
	ge.CreateAura(token.NewPonder(1))
	FireEndOfTurn(ge)
	if h := ge.Hand(); len(h) != 0 {
		t.Errorf("Hand = %v, want empty (no deck to draw from)", h)
	}
	if len(ge.Auras()) != 0 {
		t.Errorf("Auras = %+v, want empty (Ponder still destroys with empty deck)", ge.Auras())
	}
}

// Tests that Sigil of the Arknight's handler reveals an attack action.
func TestProcessAurasAtStartOfTurn_RevealsAttackActionIntoHand(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	slash := cards.AetherSlashRed{}
	_, total, revealed, _ := ProcessAurasAtStartOfTurnForTest(play.Auras(), DeckOf(slash))
	if total != 0 {
		t.Errorf("total = %d, want 0 (reveal contributes via hand, not damage)", total)
	}
	if len(revealed) != 1 || revealed[0] != slash {
		t.Errorf("revealed = %v, want [%v]", revealed, slash)
	}
}

// TestProcessAurasAtStartOfTurn_CascadingReveals.
func TestProcessAurasAtStartOfTurn_CascadingReveals(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	first := cards.AetherSlashRed{}
	second := cards.ConsumingVolitionRed{}
	_, _, revealed, _ := ProcessAurasAtStartOfTurnForTest(play.Auras(), DeckOf(first, second))
	if len(revealed) != 2 {
		t.Fatalf("len(revealed) = %d, want 2 (two cascading reveals)", len(revealed))
	}
	if revealed[0] != first || revealed[1] != second {
		t.Errorf("revealed = %v, want [%v, %v]", revealed, first, second)
	}
}

// TestProcessAurasAtStartOfTurn_NonAttackActionTopSkipsReveal.
func TestProcessAurasAtStartOfTurn_NonAttackActionTopSkipsReveal(t *testing.T) {
	play := gameengine.New()
	sigil := cards.SigilOfTheArknightBlue{}
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: sigil})
	_, total, revealed, _ := ProcessAurasAtStartOfTurnForTest(play.Auras(), DeckOf(sigil))
	if total != 0 {
		t.Errorf("total = %d, want 0 (non-attack top, no credit)", total)
	}
	if revealed != nil {
		t.Errorf("revealed = %v, want nil (non-attack tops aren't moved)", revealed)
	}
}

// Tests that a deck whose Blessing of Occult fires on later turns scores at least 1.
func TestEvaluate_TriggersFromLastTurnSurfacesInMean(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	slash := cards.AetherSlashRed{}
	deckCards := make([]deck.Card, 0, 20)
	for i := 0; i < 8; i++ {
		deckCards = append(deckCards, blessing)
	}
	for i := 0; i < 6; i++ {
		deckCards = append(deckCards, slash)
	}
	for i := 0; i < 6; i++ {
		deckCards = append(deckCards, testutils.BlueAttack{})
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	rng := rand.New(rand.NewSource(42))
	stats := NewEvaluator().Evaluate(d, 20, Matchup{}, rng)
	if stats.Best.Value < 1 {
		t.Errorf("Best.Value = %d, want at least 1 (Blessing carryover should land)", stats.Best.Value)
	}
}
