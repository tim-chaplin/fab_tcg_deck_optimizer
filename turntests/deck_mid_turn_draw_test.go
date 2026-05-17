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
	deckCards := []deck.Card{
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
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, _ := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, nil)

	wantHand := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.YellowAttack{},
	}
	if !reflect.DeepEqual(gs.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (full 4-card refill from positions 5..8; Yellow at slot 3 proves drawn card arsenaled rather than held)", gs.Hand(), wantHand)
	}

	if gs.Arsenal() != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (drawn card should take the empty arsenal slot)", gs.Arsenal(), beacon)
	}

	if gs.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", gs.RunechantCount())
	}
}

// Tests that with 2 mid-turn draws and an empty arsenal, exactly one drawn card arsenals and
// the other stays Held into turn 2.
func TestEvalOneTurn_TwoMidTurnDraws_OneArsenalsOneHeld(t *testing.T) {
	beacon := testutils.RedAttack{}
	deckCards := []deck.Card{
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
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, _ := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, nil)

	// One beacon arsenaled, the other held at slot 0; the remaining three slots are the fresh
	// refill from deck positions 6..8.
	wantHand := []card.Card{
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(gs.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (one beacon held + 3 fresh Blues; two beacons here would mean neither got arsenaled, a Yellow would mean the sim over-drew)", gs.Hand(), wantHand)
	}

	if gs.Arsenal() != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the two drawn beacons should fill the empty slot)", gs.Arsenal(), beacon)
	}

	if gs.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", gs.RunechantCount())
	}
}

// Tests 3 mid-turn draws filling a slot vacated by arsenal-in: one arsenals, two stay Held.
func TestEvalOneTurn_ThreeMidTurnDraws_ArsenalFromDrawnPool(t *testing.T) {
	beacon := testutils.RedAttack{}
	arsenalIn := cards.SnatchRed{}
	deckCards := []deck.Card{
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
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, _ := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.GameStateBuilder().SetArsenal(arsenalIn).Build(), nil)

	// Two held beacons plus two fresh Blues from deck positions 7..8.
	wantHand := []card.Card{
		beacon,
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(gs.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (two beacons held + 2 fresh Blues; a Yellow here would indicate the sim pulled more than 2 refill cards)", gs.Hand(), wantHand)
	}

	if gs.Arsenal() != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the three drawn beacons should fill the slot vacated by arsenal-in Snatch)", gs.Arsenal(), beacon)
	}

	if gs.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", gs.RunechantCount())
	}
}

// Tests that with an occupied arsenal slot, a mid-turn-drawn card stays Held into turn 2.
func TestEvalOneTurn_MidTurnDrawHeldWhenArsenalFull(t *testing.T) {
	beacon := testutils.RedAttack{}
	arsenalIn := cards.ToughenUpBlue{} // DR, cost 2, defense 4 — stays in arsenal with incoming 0
	deckCards := []deck.Card{
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
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	gs, _ := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.GameStateBuilder().SetArsenal(arsenalIn).Build(), nil)

	wantHand := []card.Card{
		beacon,
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	if !reflect.DeepEqual(gs.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (beacon held + 3 fresh Blues; a Yellow here means the sim over-drew past the 3-card budget)", gs.Hand(), wantHand)
	}

	if gs.Arsenal() != arsenalIn {
		t.Errorf("turn 2 arsenal = %v, want %v (arsenal-in should remain untouched when no better candidate beats it)", gs.Arsenal(), arsenalIn)
	}

	if gs.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", gs.RunechantCount())
	}
}

// Tests that without go-again the drawn card and a Held card share the post-chain pool —
// exactly one arsenals, the other anchors turn 2's hand (either outcome is accepted).
func TestEvalOneTurn_MidTurnDrawSansGoAgainStaysHeld(t *testing.T) {
	initialHand := []deck.Card{
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
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, initialHand)

	// Turn 1 damage: Snatch alone for 4 (no chain extension, no Viserai trigger — Snatch isn't
	// Runeblade and nothing else was played).
	if extras.Value != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch alone; chain couldn't extend)", extras.Value)
	}

	// One of {Toughen Up, Aether Slash} lands in arsenal; the other anchors turn 2's hand.
	if gs.Arsenal() == nil {
		t.Fatalf("turn 2 arsenal is nil; want one of {Toughen Up, Aether Slash}")
	}
	arsenalIsTU := gs.Arsenal().ID() == ids.ToughenUpBlue
	arsenalIsSlash := gs.Arsenal().ID() == ids.AetherSlashRed
	if !arsenalIsTU && !arsenalIsSlash {
		t.Errorf("turn 2 arsenal = %v, want Toughen Up Blue or Aether Slash Red", gs.Arsenal())
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
	if !reflect.DeepEqual(gs.Hand(), wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", gs.Hand(), wantHand)
	}

	// Deck is fully consumed: 4 deck cards minus 1 Slash drawn mid-turn = 3 Blues, all in the
	// turn 2 refill alongside the held anchor.
	if gs.Deck().Size() != 0 {
		t.Errorf("turn 2 deck = %v, want empty", gs.Deck())
	}

	if gs.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", gs.RunechantCount())
	}
}

// Tests that DrawOne against an empty deck is a no-op (no panic, no spurious draw).
func TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop(t *testing.T) {
	initialHand := []deck.Card{cards.SnatchRed{}}
	d := deck.New(heroes.Viserai{}, nil, nil)
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, initialHand)

	if extras.Value != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch damage; DrawOne is a no-op on empty deck)", extras.Value)
	}
	if len(gs.Hand()) != 0 {
		t.Errorf("turn 2 hand = %v, want empty (deck was empty, can't refill)", gs.Hand())
	}
	if gs.Deck().Size() != 0 {
		t.Errorf("turn 2 deck = %v, want empty", gs.Deck())
	}
	if gs.Arsenal() != nil {
		t.Errorf("turn 2 arsenal = %v, want nil (nothing Held to promote)", gs.Arsenal())
	}
	if gs.RunechantCount() != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", gs.RunechantCount())
	}
}
