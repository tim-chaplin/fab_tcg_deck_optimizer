package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a mid-turn-drawn card fills an empty arsenal slot at end of turn rather than
// being held into turn 2.
func TestEvalTwoTurns_MidTurnDrawArsenalsWhenSlotEmpty(t *testing.T) {
	beacon := testutils.RedAttack{}
	deckCards := []Card{
		cards.SnatchRed{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
		testutils.YellowAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, nil)

	wantHand := []Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.DealtHand, wantHand) {
		t.Errorf("turn 2 dealt hand = %v, want %v (full 4-card refill from positions 5..8; Yellow at slot 3 proves drawn card arsenaled rather than held)", got.Turn2.DealtHand, wantHand)
	}

	if got.Turn1.State.Arsenal != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (drawn card should take the empty arsenal slot)", got.Turn1.State.Arsenal, beacon)
	}

	// Remaining deck: one untouched Yellow from source position 9, then the pitched Blue
	// recycled to the bottom on turn 1.
	wantDeck := []Card{
		testutils.YellowAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.State.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", got.Turn2.State.Deck, wantDeck)
	}

	if got.Turn1.State.Runechants() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", got.Turn1.State.Runechants())
	}
}

// Tests that with 2 mid-turn draws and an empty arsenal, exactly one drawn card arsenals and
// the other stays Held into turn 2.
func TestEvalTwoTurns_TwoMidTurnDraws_OneArsenalsOneHeld(t *testing.T) {
	beacon := testutils.RedAttack{}
	deckCards := []Card{
		cards.FlyingHighRed{},
		cards.FlyingHighRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
		beacon,
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, nil)

	// One beacon arsenaled, the other held at slot 0; the remaining three slots are the fresh
	// refill from deck positions 6..8.
	wantHand := []Card{
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.DealtHand, wantHand) {
		t.Errorf("turn 2 dealt hand = %v, want %v (one beacon held + 3 fresh Blues; two beacons here would mean neither got arsenaled, a Yellow would mean the sim over-drew)", got.Turn2.DealtHand, wantHand)
	}

	if got.Turn1.State.Arsenal != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the two drawn beacons should fill the empty slot)", got.Turn1.State.Arsenal, beacon)
	}

	// Remaining deck: only the Yellow tripwire at source position 9. Turn 1 had no pitches
	// (all four cards played as attacks), so nothing recycled to the bottom.
	wantDeck := []Card{
		testutils.YellowAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.State.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", got.Turn2.State.Deck, wantDeck)
	}

	if got.Turn1.State.Runechants() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", got.Turn1.State.Runechants())
	}
}

// Tests 3 mid-turn draws filling a slot vacated by arsenal-in: one arsenals, two stay Held.
func TestEvalTwoTurns_ThreeMidTurnDraws_ArsenalFromDrawnPool(t *testing.T) {
	beacon := testutils.RedAttack{}
	arsenalIn := cards.SnatchRed{}
	deckCards := []Card{
		cards.FlyingHighRed{},
		cards.FlyingHighRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
		beacon,
		beacon,
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{Arsenal: arsenalIn}, nil)

	// Two held beacons plus two fresh Blues from deck positions 7..8.
	wantHand := []Card{
		beacon,
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.DealtHand, wantHand) {
		t.Errorf("turn 2 dealt hand = %v, want %v (two beacons held + 2 fresh Blues; a Yellow here would indicate the sim pulled more than 2 refill cards)", got.Turn2.DealtHand, wantHand)
	}

	if got.Turn1.State.Arsenal != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the three drawn beacons should fill the slot vacated by arsenal-in Snatch)", got.Turn1.State.Arsenal, beacon)
	}

	// Remaining deck: only the Yellow tripwire. Turn 1 had no pitches.
	wantDeck := []Card{
		testutils.YellowAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.State.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", got.Turn2.State.Deck, wantDeck)
	}

	if got.Turn1.State.Runechants() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", got.Turn1.State.Runechants())
	}
}

// Tests that with an occupied arsenal slot, a mid-turn-drawn card stays Held into turn 2.
func TestEvalTwoTurns_MidTurnDrawHeldWhenArsenalFull(t *testing.T) {
	beacon := testutils.RedAttack{}
	arsenalIn := cards.ToughenUpBlue{} // DR, cost 2, defense 4 — stays in arsenal with incoming 0
	deckCards := []Card{
		cards.SnatchRed{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
		testutils.YellowAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{Arsenal: arsenalIn}, nil)

	wantHand := []Card{
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.DealtHand, wantHand) {
		t.Errorf("turn 2 dealt hand = %v, want %v (beacon held + 3 fresh Blues; a Yellow here means the sim over-drew past the 3-card budget)", got.Turn2.DealtHand, wantHand)
	}

	if got.Turn1.State.Arsenal != arsenalIn {
		t.Errorf("turn 2 arsenal = %v, want %v (arsenal-in should remain untouched when no better candidate beats it)", got.Turn1.State.Arsenal, arsenalIn)
	}

	// Remaining deck: two untouched Yellows from positions 8..9, then the pitched Blue
	// recycled to the bottom on turn 1.
	wantDeck := []Card{
		testutils.YellowAttack{},
		testutils.YellowAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.State.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", got.Turn2.State.Deck, wantDeck)
	}

	if got.Turn1.State.Runechants() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", got.Turn1.State.Runechants())
	}
}

// Tests that without go-again the drawn card and a Held card share the post-chain pool —
// exactly one arsenals, the other anchors turn 2's hand (either outcome is accepted).
func TestEvalTwoTurns_MidTurnDrawSansGoAgainStaysHeld(t *testing.T) {
	initialHand := []Card{
		cards.SnatchRed{},
		cards.ToughenUpBlue{},
	}
	deckCards := []Card{
		cards.AetherSlashRed{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, initialHand)

	// Turn 1 damage: Snatch alone for 4 (no chain extension, no Viserai trigger — Snatch isn't
	// Runeblade and nothing else was played).
	if got.Turn1.Value != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch alone; chain couldn't extend)", got.Turn1.Value)
	}

	// One of {Toughen Up, Aether Slash} lands in arsenal; the other anchors turn 2's hand.
	if got.Turn1.State.Arsenal == nil {
		t.Fatalf("turn 2 arsenal is nil; want one of {Toughen Up, Aether Slash}")
	}
	arsenalIsTU := got.Turn1.State.Arsenal.ID() == ids.ToughenUpBlue
	arsenalIsSlash := got.Turn1.State.Arsenal.ID() == ids.AetherSlashRed
	if !arsenalIsTU && !arsenalIsSlash {
		t.Errorf("turn 2 arsenal = %v, want Toughen Up Blue or Aether Slash Red", got.Turn1.State.Arsenal)
	}

	// Turn 2 hand: the non-promoted of the two anchors the held prefix, then three fresh Blues
	// from the deck (positions 1..3).
	var wantAnchor Card = cards.ToughenUpBlue{}
	if arsenalIsTU {
		wantAnchor = cards.AetherSlashRed{}
	}
	wantHand := []Card{
		wantAnchor,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(got.Turn2.DealtHand, wantHand) {
		t.Errorf("turn 2 dealt hand = %v, want %v", got.Turn2.DealtHand, wantHand)
	}

	// Deck is fully consumed: 4 deck cards minus 1 Slash drawn mid-turn = 3 Blues, all in the
	// turn 2 refill alongside the held anchor.
	if len(got.Turn2.State.Deck) != 0 {
		t.Errorf("turn 2 deck = %v, want empty", got.Turn2.State.Deck)
	}

	if got.Turn1.State.Runechants() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", got.Turn1.State.Runechants())
	}
}

// Tests that DrawOne against an empty deck is a no-op (no panic, no spurious draw).
func TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop(t *testing.T) {
	initialHand := []Card{cards.SnatchRed{}}
	d := New(heroes.Viserai{}, nil, nil)
	summary := d.EvalOneTurnForTesting(Matchup{IncomingDamage: 0}, TurnState{}, initialHand)
	state := summary.State

	if summary.Value != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch damage; DrawOne is a no-op on empty deck)", summary.Value)
	}
	if len(state.Hand) != 0 {
		t.Errorf("end-of-chain hand = %v, want empty (deck was empty, can't refill)", state.Hand)
	}
	if len(state.Deck) != 0 {
		t.Errorf("end-of-chain deck = %v, want empty", state.Deck)
	}
	if state.Arsenal != nil {
		t.Errorf("end-of-chain arsenal = %v, want nil (nothing Held to promote)", state.Arsenal)
	}
	if state.Runechants() != 0 {
		t.Errorf("runechants = %d, want 0", state.Runechants())
	}
}
