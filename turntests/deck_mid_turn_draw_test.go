package turntests

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a mid-turn-drawn card fills an empty arsenal slot at end of turn rather than
// being held into turn 2.
func TestEvalOneTurn_MidTurnDrawArsenalsWhenSlotEmpty(t *testing.T) {
	beacon := testutils.RedAttack{}
	hand := []card.Card{
		cards.SnatchRed{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	deckCards := []deck.Card{
		beacon, // Snatch draws this mid-turn.
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
		testutils.YellowAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	wantHand := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	if !reflect.DeepEqual(summary.State.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (full 4-card refill; Yellow at slot 3 proves drawn card arsenaled rather than held)", summary.State.Hand(), wantHand)
	}

	if summary.State.Arsenal() != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (drawn card should take the empty arsenal slot)", summary.State.Arsenal(), beacon)
	}

	if summary.State.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", summary.State.RunechantCount())
	}
}

// Tests that with 2 mid-turn draws and an empty arsenal, exactly one drawn card arsenals and
// the other stays Held into turn 2.
func TestEvalOneTurn_TwoMidTurnDraws_OneArsenalsOneHeld(t *testing.T) {
	beacon := testutils.RedAttack{}
	hand := []card.Card{
		cards.FlyingHighRed{},
		cards.FlyingHighRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
	}
	deckCards := []deck.Card{
		beacon, // Snatch #1 draws this.
		beacon, // Snatch #2 draws this.
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	// One beacon arsenaled, the other held at slot 0; the remaining three slots are the
	// fresh refill from the remaining Blues.
	wantHand := []card.Card{
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(summary.State.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (one beacon held + 3 fresh Blues; two beacons here would mean neither got arsenaled, a Yellow would mean the sim over-drew)", summary.State.Hand(), wantHand)
	}

	if summary.State.Arsenal() != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the two drawn beacons should fill the empty slot)", summary.State.Arsenal(), beacon)
	}

	if summary.State.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", summary.State.RunechantCount())
	}
}

// Tests 3 mid-turn draws filling a slot vacated by arsenal-in: one arsenals, two stay Held.
func TestEvalOneTurn_ThreeMidTurnDraws_ArsenalFromDrawnPool(t *testing.T) {
	beacon := testutils.RedAttack{}
	arsenalIn := cards.SnatchRed{}
	hand := []card.Card{
		cards.FlyingHighRed{},
		cards.FlyingHighRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
	}
	deckCards := []deck.Card{
		beacon, // Drawn by arsenal-in Snatch.
		beacon, // Drawn by hand Snatch #1.
		beacon, // Drawn by hand Snatch #2.
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetArsenal(arsenalIn).Build(), hand)

	// Two held beacons plus two fresh Blues from the remaining deck.
	wantHand := []card.Card{
		beacon,
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(summary.State.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (two beacons held + 2 fresh Blues; a Yellow here would indicate the sim pulled more than 2 refill cards)", summary.State.Hand(), wantHand)
	}

	if summary.State.Arsenal() != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the three drawn beacons should fill the slot vacated by arsenal-in Snatch)", summary.State.Arsenal(), beacon)
	}

	if summary.State.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", summary.State.RunechantCount())
	}
}

// Tests that with an occupied arsenal slot, a mid-turn-drawn card stays Held into turn 2.
func TestEvalOneTurn_MidTurnDrawHeldWhenArsenalFull(t *testing.T) {
	beacon := testutils.RedAttack{}
	arsenalIn := cards.ToughenUpBlue{} // DR, cost 2, defense 4 — stays in arsenal with incoming 0
	hand := []card.Card{
		cards.SnatchRed{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	deckCards := []deck.Card{
		beacon, // Snatch draws this mid-turn.
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
		testutils.YellowAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetArsenal(arsenalIn).Build(), hand)

	wantHand := []card.Card{
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(summary.State.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (beacon held + 3 fresh Blues; a Yellow here means the sim over-drew past the 3-card budget)", summary.State.Hand(), wantHand)
	}

	if summary.State.Arsenal() != arsenalIn {
		t.Errorf("turn 2 arsenal = %v, want %v (arsenal-in should remain untouched when no better candidate beats it)", summary.State.Arsenal(), arsenalIn)
	}

	if summary.State.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", summary.State.RunechantCount())
	}
}

// Tests that without go-again the drawn card and a Held card share the post-chain pool —
// exactly one arsenals, the other anchors turn 2's hand (either outcome is accepted).
func TestEvalOneTurn_MidTurnDrawSansGoAgainStaysHeld(t *testing.T) {
	initialHand := []card.Card{
		cards.SnatchRed{},
		cards.ToughenUpBlue{},
	}
	deckCards := []deck.Card{
		cards.AetherSlashRed{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, nil, initialHand)

	// Turn 1 damage: Snatch alone for 4 (no chain extension, no Viserai trigger — Snatch isn't
	// Runeblade and nothing else was played).
	if summary.Value != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch alone; chain couldn't extend)", summary.Value)
	}

	// One of {Toughen Up, Aether Slash} lands in arsenal; the other anchors turn 2's hand.
	if summary.State.Arsenal() == nil {
		t.Fatalf("turn 2 arsenal is nil; want one of {Toughen Up, Aether Slash}")
	}
	arsenalIsTU := summary.State.Arsenal().ID() == ids.ToughenUpBlue
	arsenalIsSlash := summary.State.Arsenal().ID() == ids.AetherSlashRed
	if !arsenalIsTU && !arsenalIsSlash {
		t.Errorf("turn 2 arsenal = %v, want Toughen Up Blue or Aether Slash Red", summary.State.Arsenal())
	}

	// Turn 2 hand: the non-promoted of the two anchors the held prefix, then three fresh Blues
	// from the deck (positions 1..3).
	var wantAnchor card.Card = cards.ToughenUpBlue{}
	if arsenalIsTU {
		wantAnchor = cards.AetherSlashRed{}
	}
	wantHand := []card.Card{
		wantAnchor,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(summary.State.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", summary.State.Hand(), wantHand)
	}

	// Deck is fully consumed: 4 deck cards minus 1 Slash drawn mid-turn = 3 Blues, all in the
	// turn 2 refill alongside the held anchor.
	if summary.State.Deck().Size() != 0 {
		t.Errorf("turn 2 deck = %v, want empty", summary.State.Deck())
	}

	if summary.State.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", summary.State.RunechantCount())
	}
}

// Tests that DrawOne against an empty deck is a no-op (no panic, no spurious draw).
func TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop(t *testing.T) {
	initialHand := []card.Card{cards.SnatchRed{}}
	d := deck.New(heroes.Viserai{}, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, nil, initialHand)

	if summary.Value != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch damage; DrawOne is a no-op on empty deck)", summary.Value)
	}
	if len(summary.State.Hand()) != 0 {
		t.Errorf("turn 2 hand = %v, want empty (deck was empty, can't refill)", summary.State.Hand())
	}
	if summary.State.Deck().Size() != 0 {
		t.Errorf("turn 2 deck = %v, want empty", summary.State.Deck())
	}
	if summary.State.Arsenal() != nil {
		t.Errorf("turn 2 arsenal = %v, want nil (nothing Held to promote)", summary.State.Arsenal())
	}
	if summary.State.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", summary.State.RunechantCount())
	}
}
