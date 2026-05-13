package sim_test

import (
	"math/rand"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// damageTrigger returns a StartOfTurn Aura crediting the given damage and destroying
// itself on the first fire.
func damageTrigger(self card.Card, damage int, calls *int) gameengine.Aura {
	return NewCardAura(
		&card.CardState{Card: self},
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
	aura := testutils.RedAttack{}
	var callsA, callsB int
	queue := []gameengine.Aura{damageTrigger(aura, 2, &callsA), damageTrigger(aura, 3, &callsB)}
	survivors, contribs, total, _, _ := ProcessAurasAtStartOfTurn(queue, DeckOf())
	if total != 5 {
		t.Errorf("total = %d, want 5 (2+3)", total)
	}
	if len(contribs) != 2 || contribs[0].Damage != 2 || contribs[1].Damage != 3 {
		t.Errorf("contribs = %+v, want [{_, 2}, {_, 3}]", contribs)
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
	survivors, contribs, total, _, _ := ProcessAurasAtStartOfTurn(nil, DeckOf())
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if contribs != nil || len(survivors) != 0 {
		t.Errorf("non-empty outputs on empty input: contribs=%v survivors=%v",
			contribs, survivors)
	}
}

// TestProcessAurasAtStartOfTurn_GraveyardsExhaustedAura: a handler that calls Destroy
// lands Self in the graveyard before subsequent handlers run.
func TestProcessAurasAtStartOfTurn_GraveyardsExhaustedAura(t *testing.T) {
	aura := testutils.RedAttack{}
	var seen []card.Card
	watcher := NewCardAura(
		&card.CardState{Card: testutils.YellowAttack{}},
		triggertype.StartOfTurn,
		func(ge card.GameEngine, _ card.Logger, _ card.Aura) {
			eng := ge.(*gameengine.GameEngine)
			seen = append([]card.Card(nil), eng.Graveyard()...)
		},
		1,
		false,
	)
	first := NewCardAura(
		&card.CardState{Card: aura},
		triggertype.StartOfTurn,
		func(_ card.GameEngine, _ card.Logger, a card.Aura) {
			a.Destroy(true)
		},
		1,
		false,
	)
	_, _, _, _, _ = ProcessAurasAtStartOfTurn([]gameengine.Aura{first, watcher}, DeckOf())
	if len(seen) != 1 || seen[0] != aura {
		t.Errorf("second handler saw Graveyard = %v, want [%v]", seen, aura)
	}
}

// TestProcessAurasAtStartOfTurn_IgnoresPonder: Ponder fires at end-of-turn (not start).
func TestProcessAurasAtStartOfTurn_IgnoresPonder(t *testing.T) {
	survivors, _, _, _, _ := ProcessAurasAtStartOfTurn([]gameengine.Aura{NewPonderAura(1)}, DeckOf())
	if len(survivors) != 1 || survivors[0].CardName() != "Ponder" {
		t.Errorf("survivors = %+v, want one Ponder aura intact", survivors)
	}
}

// TestFireEndOfTurn_PonderPopsDeckTopIntoHand.
func TestFireEndOfTurn_PonderPopsDeckTopIntoHand(t *testing.T) {
	a, b, c := testutils.NewStubCard("a"), testutils.NewStubCard("b"), testutils.NewStubCard("c")
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{a, b, c}).Build()}
	ge.CreateAura(NewPonderAura(2))
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
	ge.CreateAura(NewPonderAura(1))
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
	_, contribs, total, revealed, _ := ProcessAurasAtStartOfTurn(play.Auras(), DeckOf(slash))
	if total != 0 {
		t.Errorf("total = %d, want 0 (reveal contributes via hand, not damage)", total)
	}
	if len(revealed) != 1 || revealed[0] != slash {
		t.Errorf("revealed = %v, want [%v]", revealed, slash)
	}
	if len(contribs) != 1 || contribs[0].Card.ID() != ids.SigilOfTheArknightBlue {
		t.Errorf("contribs = %+v, want one entry for the sigil", contribs)
	}
}

// Tests that the TriggerContribution carries the revealed card.
func TestProcessAurasAtStartOfTurn_AttributesRevealedToContribution(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	slash := cards.AetherSlashRed{}
	_, contribs, _, _, _ := ProcessAurasAtStartOfTurn(play.Auras(), DeckOf(slash))
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	if contribs[0].Revealed == nil || contribs[0].Revealed.ID() != ids.AetherSlashRed {
		t.Errorf("contribs[0].Revealed = %v, want Aether Slash [R]", contribs[0].Revealed)
	}
}

// TestProcessAurasAtStartOfTurn_CascadingReveals.
func TestProcessAurasAtStartOfTurn_CascadingReveals(t *testing.T) {
	play := gameengine.New()
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	first := cards.AetherSlashRed{}
	second := cards.ConsumingVolitionRed{}
	_, _, _, revealed, _ := ProcessAurasAtStartOfTurn(play.Auras(), DeckOf(first, second))
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
	_, _, total, revealed, _ := ProcessAurasAtStartOfTurn(play.Auras(), DeckOf(sigil))
	if total != 0 {
		t.Errorf("total = %d, want 0 (non-attack top, no credit)", total)
	}
	if revealed != nil {
		t.Errorf("revealed = %v, want nil (non-attack tops aren't moved)", revealed)
	}
}

// TestProcessAurasAtStartOfTurn_SigilHitAuthorsLogText.
func TestProcessAurasAtStartOfTurn_SigilHitAuthorsLogText(t *testing.T) {
	play := gameengine.New()
	sigil := cards.SigilOfTheArknightBlue{}
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: sigil})
	_, contribs, _, _, _ := ProcessAurasAtStartOfTurn(play.Auras(), DeckOf(cards.AetherSlashRed{}))
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	want := "Sigil of the Arknight [B] drew Aether Slash [R] into hand"
	if contribs[0].Text != want {
		t.Errorf("contribs[0].Text = %q, want %q", contribs[0].Text, want)
	}
}

// TestProcessAurasAtStartOfTurn_SigilWhiffStillLogs.
func TestProcessAurasAtStartOfTurn_SigilWhiffStillLogs(t *testing.T) {
	play := gameengine.New()
	sigil := cards.SigilOfTheArknightBlue{}
	play.ResolveChainStep(play.Logger(), &card.CardState{Card: sigil})
	_, contribs, _, _, _ := ProcessAurasAtStartOfTurn(play.Auras(), DeckOf(sigil))
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	want := "Sigil of the Arknight [B] revealed Sigil of the Arknight [B] but didn't draw it"
	if contribs[0].Text != want {
		t.Errorf("contribs[0].Text = %q, want %q", contribs[0].Text, want)
	}
}

// Tests that Evaluate surfaces a start-of-turn trigger as a TriggersFromLastTurn entry on
// the best-scoring hand.
func TestEvaluate_TriggersFromLastTurnSurfacesInBest(t *testing.T) {
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

	if len(stats.Best.Summary.TriggersFromLastTurn) == 0 {
		t.Errorf("Stats.Best.Summary.TriggersFromLastTurn is empty; Best.Value=%d",
			stats.Best.Summary.Value)
	}
	foundBlessing := false
	for _, a := range stats.Best.Summary.StartOfTurnAuras {
		if a.ID() == ids.BlessingOfOccultRed {
			foundBlessing = true
			break
		}
	}
	if !foundBlessing {
		t.Errorf("Stats.Best.Summary.StartOfTurnAuras missing Blessing; got %+v",
			stats.Best.Summary.StartOfTurnAuras)
	}
}

// Tests that the OncePerTurn FiredThisTurn flag is cleared at every turn boundary.
func TestProcessAurasAtStartOfTurn_ReArmsOncePerTurnGate(t *testing.T) {
	aura := testutils.RedAttack{}
	exhausted := NewCardAura(
		&card.CardState{Card: aura},
		triggertype.AttackAction,
		func(card.GameEngine, card.Logger, card.Aura) {},
		2,
		true,
	)
	exhausted.SetFiredThisTurn(true)
	survivors, _, _, _, _ := ProcessAurasAtStartOfTurn([]gameengine.Aura{exhausted}, DeckOf())
	if len(survivors) != 1 {
		t.Fatalf("survivors len = %d, want 1 (AttackAction trigger passes through)", len(survivors))
	}
	if survivors[0].FiredThisTurn() {
		t.Errorf("FiredThisTurn = true, want false (turn-boundary reset)")
	}
	if survivors[0].Count() != 2 {
		t.Errorf("Count = %d, want 2 (only re-arm; don't tick)", survivors[0].Count())
	}
}
