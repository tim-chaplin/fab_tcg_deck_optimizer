package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests post-hoc Arsenal promotion of a Held card when the slot is empty.
func TestBest_EmptyArsenalClaimsHeldCard(t *testing.T) {
	h := []card.Card{cards.ToughenUpBlue{}}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetIncomingDamage(4).
		Build())
	if got.BestLine[0].Role != card.Arsenal {
		t.Errorf("Roles[0] = %s, want ARSENAL", got.BestLine[0].Role)
	}
	if got.State.Arsenal() == nil || got.State.Arsenal().ID() != ids.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue", got.State.Arsenal())
	}
}

// Tests an arsenal-in DR being played (funded by a hand pitch) and vacating the slot.
func TestBest_ArsenalInPlayDR(t *testing.T) {
	h := []card.Card{cards.MaleficIncantationBlue{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(cards.ToughenUpBlue{}).
		SetIncomingDamage(4).
		Build()
	got := Best(nil, h, nil, state)
	if got.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Malefic pitches to pay arsenal DR, prevents 4). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.State.Arsenal() != nil {
		t.Errorf("ArsenalCard = %v, want nil (slot was vacated, no Held card to promote)", got.State.Arsenal())
	}
	// ArsenalIn surfaces the arsenal-in assignment so callers (the best-hand printout) can flag
	// that this card wasn't in hand this turn.
	ai, hasArsenal := got.ArsenalIn()
	if !hasArsenal || ai.Card.ID() != ids.ToughenUpBlue {
		t.Errorf("ArsenalIn = %v, want Toughen Up Blue", ai)
	}
	if ai.Role != card.Defend {
		t.Errorf("ArsenalIn role = %s, want DEFEND", ai.Role)
	}
}

// Tests that an occupied arsenal slot blocks Held → Arsenal post-hoc promotion.
func TestBest_ArsenalInStayBlocksNewArsenal(t *testing.T) {
	h := []card.Card{cards.ToughenUpBlue{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(cards.ToughenUpBlue{}).
		Build()
	got := Best(nil, h, nil, state)
	if got.BestLine[0].Role != card.Held {
		t.Errorf("Roles[0] = %s, want HELD (slot occupied by arsenal-in, can't promote)", got.BestLine[0].Role)
	}
	if got.State.Arsenal() == nil || got.State.Arsenal().ID() != ids.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue (the staying arsenal-in card)", got.State.Arsenal())
	}
}

// Tests the arsenal-card-played-as-attack branch (arsenal Red funded by pitching hand Red).
func TestBest_ArsenalInPlayAttack(t *testing.T) {
	h := []card.Card{testutils.FakeRedAttack().
		WithCost(1).
		WithPower(3)}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(testutils.FakeRedAttack().
			WithCost(1).
			WithPower(3)).
		Build()
	got := Best(nil, h, nil, state)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (arsenal Red played, hand Red pitched to fund it). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.State.Arsenal() != nil {
		t.Errorf("ArsenalCard = %v, want nil (slot vacated, no Held to promote)", got.State.Arsenal())
	}
}

// Tests that a non-Attack action card in arsenal is playable.
func TestBest_ArsenalInNonAttackActionPlays(t *testing.T) {
	h := []card.Card{cards.MaleficIncantationBlue{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(cards.HealingBalmRed{}).
		Build()
	got := Best(nil, h, nil, state)
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Malefic pitched, arsenal Healing Balm played for 3{h}). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.State.Arsenal() != nil {
		t.Errorf("ArsenalCard = %v, want nil (Healing Balm played out of arsenal)", got.State.Arsenal())
	}
}

// Tests that Unmovable's +1{d} arsenal-defense rider fires when played from arsenal.
func TestBest_ArsenalInUnmovableGrantsDefenseBonus(t *testing.T) {
	h := []card.Card{cards.MaleficIncantationBlue{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(cards.UnmovableRed{}).
		SetIncomingDamage(8).
		Build()
	got := Best(nil, h, nil, state)
	if got.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Unmovable from arsenal blocks 7+1). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that Unmovable's arsenal-defense rider does NOT fire when played from hand.
func TestBest_HandUnmovableNoDefenseBonus(t *testing.T) {
	h := []card.Card{cards.MaleficIncantationBlue{}, cards.UnmovableRed{}}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetIncomingDamage(8).
		Build())
	if got.Value != 7 {
		t.Fatalf("Value = %d, want 7 (hand-played Unmovable: no rider). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that Smashing Good Time's +3 rider only fires for the arsenal copy, not the hand copy.
func TestBest_ArsenalInSmashingGoodTimeGatesOnlyArsenalCopy(t *testing.T) {
	h := []card.Card{
		notimpl.SmashingGoodTimeRed{},
		cards.HocusPocusRed{},
	}
	state := gameengine.GameStateBuilder().
		SetHero(heroes.Viserai).
		SetArsenal(notimpl.SmashingGoodTimeRed{}).
		Build()
	got := Best(nil, h, nil, state)
	if got.Value != 8 {
		t.Fatalf("Value = %d, want 8 (only arsenal SGT grants +3). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// Tests that PromoteRandomHandCardToArsenal's hash-based selection lands on different slots
// across distinct hands (no bias toward slot 0).
func TestPromoteRandomHandCardToArsenal_SpreadsAcrossHands(t *testing.T) {
	// 20 different 4-card hands; varying which card sits in which slot exercises the hash.
	wbR := cards.WoundingBlowRed{}
	wbY := cards.WoundingBlowYellow{}
	wbB := cards.WoundingBlowBlue{}
	hands := [][]card.Card{
		{wbR, wbR, wbR, wbY}, {wbR, wbR, wbY, wbY}, {wbR, wbR, wbY, wbB}, {wbR, wbY, wbY, wbB},
		{wbR, wbY, wbB, wbB}, {wbR, wbR, wbR, wbB}, {wbR, wbR, wbB, wbB}, {wbY, wbY, wbY, wbB},
		{wbY, wbY, wbB, wbB}, {wbY, wbB, wbB, wbB}, {wbR, wbR, wbR, wbR}, {wbY, wbY, wbY, wbY},
		{wbB, wbB, wbB, wbB}, {wbR, wbY, wbY, wbY}, {wbR, wbR, wbY, wbR}, {wbB, wbR, wbY, wbR},
		{wbB, wbY, wbR, wbR}, {wbR, wbB, wbB, wbY}, {wbR, wbR, wbB, wbY}, {wbY, wbR, wbB, wbB},
	}
	picks := map[ids.CardID]int{}
	for _, h := range hands {
		handCopy := append([]card.Card(nil), h...)
		line := make([]card.CardAssignment, len(handCopy))
		for i, c := range handCopy {
			line[i] = card.CardAssignment{Card: c, Role: card.Held}
		}
		best := TurnSummary{
			BestLine: line,
			State:    EngineWithHand(append([]card.Card(nil), handCopy...)),
		}
		PromoteRandomHandCardToArsenal(&best, handCopy, nil)
		if best.State.Arsenal() == nil {
			t.Fatalf("hand %v: State.Arsenal nil after promotion", h)
		}
		picks[best.State.Arsenal().ID()]++
	}
	if len(picks) < 2 {
		t.Errorf("arsenal promotion only ever landed on card %v across %d hands; expected spread", picks, len(hands))
	}
}

// Tests that the same hand produces the same picked card across runs (reproducibility).
func TestPromoteRandomHandCardToArsenal_DeterministicPerHand(t *testing.T) {
	hand := []card.Card{
		cards.WoundingBlowRed{}, cards.WoundingBlowYellow{},
		cards.WoundingBlowBlue{}, cards.WoundingBlowBlue{},
	}
	var firstID ids.CardID
	for run := 0; run < 5; run++ {
		line := []card.CardAssignment{
			{Card: hand[0], Role: card.Held},
			{Card: hand[1], Role: card.Held},
			{Card: hand[2], Role: card.Held},
			{Card: hand[3], Role: card.Held},
		}
		best := TurnSummary{
			BestLine: line,
			State:    EngineWithHand(append([]card.Card(nil), hand...)),
		}
		PromoteRandomHandCardToArsenal(&best, hand, nil)
		if best.State.Arsenal() == nil {
			t.Fatalf("run %d: State.Arsenal nil", run)
		}
		got := best.State.Arsenal().ID()
		if run == 0 {
			firstID = got
			continue
		}
		if got != firstID {
			t.Errorf("run %d: Arsenal = %v, want %v (deterministic per-hand)", run, got, firstID)
		}
	}
}

// Tests the n=1 edge: with one State.Hand entry the only candidate gets promoted.
func TestPromoteRandomHandCardToArsenal_SingleCandidateAlwaysPicked(t *testing.T) {
	hand := []card.Card{cards.WoundingBlowRed{}, cards.WoundingBlowBlue{}}
	line := []card.CardAssignment{
		{Card: hand[0], Role: card.Attack},
		{Card: hand[1], Role: card.Held},
	}
	best := TurnSummary{
		BestLine: line,
		State:    EngineWithHand([]card.Card{hand[1]}),
	}
	PromoteRandomHandCardToArsenal(&best, hand, nil)
	if best.BestLine[1].Role != card.Arsenal {
		t.Errorf("Role[1] = %s, want Arsenal (only candidate)", best.BestLine[1].Role)
	}
	if best.State.Arsenal() == nil || best.State.Arsenal().ID() != hand[1].ID() {
		t.Errorf("State.Arsenal = %v, want %s", best.State.Arsenal(), hand[1].Name())
	}
}

// Tests that an empty State.Hand makes the promotion a no-op.
func TestPromoteRandomHandCardToArsenal_EmptyHandIsNoop(t *testing.T) {
	hand := []card.Card{cards.WoundingBlowRed{}, cards.WoundingBlowBlue{}}
	line := []card.CardAssignment{
		{Card: hand[0], Role: card.Attack},
		{Card: hand[1], Role: card.Pitch},
	}
	best := TurnSummary{BestLine: line, State: gameengine.GameStateBuilder().Build()}
	PromoteRandomHandCardToArsenal(&best, hand, nil)
	for i, a := range best.BestLine {
		if a.Role == card.Arsenal {
			t.Errorf("BestLine[%d].Role = Arsenal, want unchanged (no candidates)", i)
		}
	}
	if best.State.Arsenal() != nil {
		t.Errorf("State.Arsenal = %v, want nil (no promotion possible)", best.State.Arsenal())
	}
}

// Tests attackTurnScore.cmp's lexicographic ranking: value, then cardsPlayed, then totalCards,
// then totalCounters.
func TestAttackTurnScore_Cmp(t *testing.T) {
	cases := []struct {
		name string
		a, b attackTurnScore
		want int
	}{
		{"value outranks everything", attackTurnScore{value: 10}, attackTurnScore{value: 11, cardsPlayed: 9, totalCards: 9, totalCounters: 9}, -1},
		{"cardsPlayed breaks a value tie", attackTurnScore{value: 5, cardsPlayed: 2}, attackTurnScore{value: 5, cardsPlayed: 1}, 1},
		{"cardsPlayed outranks totalCards", attackTurnScore{value: 5, cardsPlayed: 2}, attackTurnScore{value: 5, cardsPlayed: 1, totalCards: 9}, 1},
		{"totalCards breaks a cardsPlayed tie", attackTurnScore{value: 5, cardsPlayed: 1, totalCards: 3}, attackTurnScore{value: 5, cardsPlayed: 1, totalCards: 2}, 1},
		{"totalCards outranks totalCounters", attackTurnScore{value: 5, totalCards: 1}, attackTurnScore{value: 5, totalCounters: 9}, 1},
		{"totalCounters breaks a totalCards tie", attackTurnScore{value: 5, totalCards: 2, totalCounters: 1}, attackTurnScore{value: 5, totalCards: 2}, 1},
		{"fully equal", attackTurnScore{value: 5, cardsPlayed: 1, totalCards: 2}, attackTurnScore{value: 5, cardsPlayed: 1, totalCards: 2}, 0},
	}
	for _, tc := range cases {
		if got := tc.a.cmp(tc.b); got != tc.want {
			t.Errorf("%s: a.cmp(b) = %d, want %d", tc.name, got, tc.want)
		}
		if got := tc.b.cmp(tc.a); got != -tc.want {
			t.Errorf("%s: b.cmp(a) = %d, want %d (antisymmetric)", tc.name, got, -tc.want)
		}
	}
}
