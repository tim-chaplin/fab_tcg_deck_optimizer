package deck

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestEvalOneTurn_MidTurnDrawArsenalsWhenSlotEmpty pins the ARSENAL disposition for mid-turn-
// drawn cards: when the arsenal slot is empty at end of turn 1, the card Snatch drew mid-turn
// becomes turn 2's arsenal-in. Turn 2's hand is then a full handSize refill from the top of
// the deck — including a Yellow tripwire at position 8 — rather than the beacon at slot 0
// plus three fresh Blues (which would indicate the drawn card was held instead of arsenaled).
//
// Deck layout (consumed in source order):
//   - positions 0..3 = turn 1's hand: Snatch Red (cost 0, attack 4, on-hit DrawOne) + three
//     Blues that chain for Value 6 (pitch 1 Blue, Blue + Blue + Snatch for 1 + 1 + 4 damage).
//   - position 4 = the beacon (fake.RedAttack) that Snatch draws mid-turn.
//   - positions 5..7 = Blues that make up turn 2's refill.
//   - positions 8..9 = Yellow tripwires — a Yellow only shows up in turn 2's hand when the
//     sim over-draws past the expected refill count.
func TestEvalOneTurn_MidTurnDrawArsenalsWhenSlotEmpty(t *testing.T) {
	beacon := fake.RedAttack{}
	deckCards := []card.Card{
		cards.SnatchRed{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
		fake.YellowAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, nil)

	wantHand := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (full 4-card refill from positions 5..8; Yellow at slot 3 proves drawn card arsenaled rather than held)", state.Hand, wantHand)
	}

	if state.ArsenalCard != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (drawn card should take the empty arsenal slot)", state.ArsenalCard, beacon)
	}

	// Remaining deck: one untouched Yellow from source position 9, then the pitched Blue
	// recycled to the bottom on turn 1.
	wantDeck := []card.Card{
		fake.YellowAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	if state.Runechants != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", state.Runechants)
	}
}

// TestEvalOneTurn_TwoMidTurnDraws_OneArsenalsOneHeld pins the arsenal pool's behaviour when
// more drawn cards exist than slots: with 2 mid-turn draws and an empty arsenal, exactly one
// drawn card takes the slot and the other stays HELD, carrying into turn 2's hand.
//
// Turn 1's chain is Flying High → Snatch → Flying High → Snatch: each Flying High grants go
// again to the next attack (and +1 power when the pitch matches — both Reds here), letting
// the two Snatches both fire and consume the top two cards of the deck via their on-hit
// DrawOne. The winning partition has no Held hand cards (all four played), so the arsenal
// pool is just the two drawn cards.
//
// Deck layout (consumed in source order):
//   - positions 0..3 = turn 1's hand: Flying High Red + Flying High Red + Snatch Red + Snatch Red.
//   - positions 4..5 = two identical Red beacons that Snatch's on-hit DrawOne consumes mid-turn.
//   - positions 6..8 = Blues that make up turn 2's refill behind the held beacon.
//   - position 9 = Yellow tripwire — showing up in turn 2's hand would indicate the sim
//     pulled more than handSize - 1 refill cards.
func TestEvalOneTurn_TwoMidTurnDraws_OneArsenalsOneHeld(t *testing.T) {
	beacon := fake.RedAttack{}
	deckCards := []card.Card{
		cards.FlyingHighRed{},
		cards.FlyingHighRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
		beacon,
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, nil)

	// One beacon arsenaled, the other held at slot 0; the remaining three slots are the fresh
	// refill from deck positions 6..8.
	wantHand := []card.Card{
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (one beacon held + 3 fresh Blues; two beacons here would mean neither got arsenaled, a Yellow would mean the sim over-drew)", state.Hand, wantHand)
	}

	if state.ArsenalCard != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the two drawn beacons should fill the empty slot)", state.ArsenalCard, beacon)
	}

	// Remaining deck: only the Yellow tripwire at source position 9. Turn 1 had no pitches
	// (all four cards played as attacks), so nothing recycled to the bottom.
	wantDeck := []card.Card{
		fake.YellowAttack{},
	}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	if state.Runechants != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", state.Runechants)
	}
}

// TestEvalOneTurn_ThreeMidTurnDraws_ArsenalFromDrawnPool pins the arsenal pool's behaviour
// when a starting arsenal-in card plays out and the drawn cards fill the vacated slot: with
// an arsenal-in Snatch plus two Flying Highs and two Snatches in hand, all three Snatches
// fire their on-hit DrawOne (the two hand Snatches inherit go again from the Flying Highs;
// the arsenal-in Snatch plays last, no chain constraint past it). That's three mid-turn
// draws. One drawn card takes the arsenal slot (vacated when arsenal-in played), the other
// two carry HELD into turn 2's hand — so turn 2 refills only handSize - 2 = 2 cards.
//
// Deck layout (consumed in source order):
//   - positions 0..3 = turn 1's hand: Flying High Red + Flying High Red + Snatch Red + Snatch Red.
//   - positions 4..6 = three identical Red beacons consumed by the three Snatch on-hit draws.
//   - positions 7..8 = Blues that make up turn 2's refill behind the two held beacons.
//   - position 9 = Yellow tripwire — appearing in turn 2's hand would mean the sim over-drew
//     past the 2-card refill budget.
func TestEvalOneTurn_ThreeMidTurnDraws_ArsenalFromDrawnPool(t *testing.T) {
	beacon := fake.RedAttack{}
	arsenalIn := cards.SnatchRed{}
	deckCards := []card.Card{
		cards.FlyingHighRed{},
		cards.FlyingHighRed{},
		cards.SnatchRed{},
		cards.SnatchRed{},
		beacon,
		beacon,
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, arsenalIn, nil)

	// Two held beacons plus two fresh Blues from deck positions 7..8.
	wantHand := []card.Card{
		beacon,
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (two beacons held + 2 fresh Blues; a Yellow here would indicate the sim pulled more than 2 refill cards)", state.Hand, wantHand)
	}

	if state.ArsenalCard != beacon {
		t.Errorf("turn 2 arsenal = %v, want %v (one of the three drawn beacons should fill the slot vacated by arsenal-in Snatch)", state.ArsenalCard, beacon)
	}

	// Remaining deck: only the Yellow tripwire. Turn 1 had no pitches.
	wantDeck := []card.Card{
		fake.YellowAttack{},
	}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	if state.Runechants != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", state.Runechants)
	}
}

// TestEvalOneTurn_MidTurnDrawHeldWhenArsenalFull pins the fallback disposition: with the
// arsenal slot already occupied (and not displaced), the drawn card can't arsenal, so it stays
// HELD and carries into turn 2's hand. Turn 2 then refills handSize - 1 = 3 cards from the top
// of the deck behind the held beacon. The Yellow tripwires at positions 8..9 should NOT show
// up in turn 2's hand — if they do, the sim pulled too many refill cards.
//
// Deck layout (consumed in source order) — same shape as the arsenal-empty variant so the
// difference between the two tests is purely the starting arsenal slot:
//   - positions 0..3 = turn 1's hand (Snatch + three Blues; chains for Value 6).
//   - position 4 = the beacon Snatch draws mid-turn.
//   - positions 5..7 = Blues that make up turn 2's refill behind the held beacon.
//   - positions 8..9 = Yellow tripwires that should stay in the deck.
func TestEvalOneTurn_MidTurnDrawHeldWhenArsenalFull(t *testing.T) {
	beacon := fake.RedAttack{}
	arsenalIn := cards.ToughenUpBlue{} // DR, cost 2, defense 4 — stays in arsenal with incoming 0
	deckCards := []card.Card{
		cards.SnatchRed{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
		fake.YellowAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, arsenalIn, nil)

	wantHand := []card.Card{
		beacon,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v (beacon held + 3 fresh Blues; a Yellow here means the sim over-drew past the 3-card budget)", state.Hand, wantHand)
	}

	if state.ArsenalCard != arsenalIn {
		t.Errorf("turn 2 arsenal = %v, want %v (arsenal-in should remain untouched when no better candidate beats it)", state.ArsenalCard, arsenalIn)
	}

	// Remaining deck: two untouched Yellows from positions 8..9, then the pitched Blue
	// recycled to the bottom on turn 1.
	wantDeck := []card.Card{
		fake.YellowAttack{},
		fake.YellowAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	if state.Runechants != 0 {
		t.Errorf("turn 2 runechants = %d, want 0 (nothing on turn 1 creates runechants)", state.Runechants)
	}
}

// TestEvalOneTurn_MidTurnDrawSansGoAgainStaysHeld: hand is Snatch + Toughen Up Blue, top of
// deck is Aether Slash. Snatch has no baseline Go again and nothing grants it, so the chain
// ends right after Snatch — the Slash is drawn but not played. Turn 1 Value = 4 (Snatch
// alone). Toughen Up and the drawn Slash share the post-chain hand pool that feeds arsenal
// promotion; the deterministic hash picks one for the arsenal slot, leaving the other to
// carry into turn 2's hand alongside three fresh Blues. The test accepts either outcome of
// the hash modulo so it stays stable across hash-input changes; the contract is "exactly one
// of the two lands in arsenal, the other anchors turn 2's hand."
func TestEvalOneTurn_MidTurnDrawSansGoAgainStaysHeld(t *testing.T) {
	initialHand := []card.Card{
		cards.SnatchRed{},
		cards.ToughenUpBlue{},
	}
	deckCards := []card.Card{
		cards.AetherSlashRed{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, initialHand)

	// Turn 1 damage: Snatch alone for 4 (no chain extension, no Viserai trigger — Snatch isn't
	// Runeblade and nothing else was played).
	if state.PrevTurnValue != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch alone; chain couldn't extend)", state.PrevTurnValue)
	}

	// One of {Toughen Up, Aether Slash} lands in arsenal; the other anchors turn 2's hand.
	if state.ArsenalCard == nil {
		t.Fatalf("turn 2 arsenal is nil; want one of {Toughen Up, Aether Slash}")
	}
	arsenalIsTU := state.ArsenalCard.ID() == card.ToughenUpBlue
	arsenalIsSlash := state.ArsenalCard.ID() == card.AetherSlashRed
	if !arsenalIsTU && !arsenalIsSlash {
		t.Errorf("turn 2 arsenal = %v, want Toughen Up Blue or Aether Slash Red", state.ArsenalCard)
	}

	// Turn 2 hand: the non-promoted of the two anchors the held prefix, then three fresh Blues
	// from the deck (positions 1..3).
	var wantAnchor card.Card = cards.ToughenUpBlue{}
	if arsenalIsTU {
		wantAnchor = cards.AetherSlashRed{}
	}
	wantHand := []card.Card{
		wantAnchor,
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", state.Hand, wantHand)
	}

	// Deck is fully consumed: 4 deck cards minus 1 Slash drawn mid-turn = 3 Blues, all in the
	// turn 2 refill alongside the held anchor.
	if len(state.Deck) != 0 {
		t.Errorf("turn 2 deck = %v, want empty", state.Deck)
	}

	if state.Runechants != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", state.Runechants)
	}
}

// TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop pins the degenerate case: Snatch fires DrawOne
// against an empty deck and nothing happens — no panic, no spurious entry in state.Drawn, no
// index math goes sideways downstream. Hand is just Snatch; d.Cards is empty. Turn 2 can't
// deal (deck stays empty through the turn), so the sim returns a TurnStartState with just
// the previous-turn value and the arsenal/live runechants.
func TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop(t *testing.T) {
	initialHand := []card.Card{cards.SnatchRed{}}
	d := New(hero.Viserai{}, nil, nil)
	state := d.EvalOneTurnForTesting(0, nil, initialHand)

	if state.PrevTurnValue != 4 {
		t.Errorf("turn 1 Value = %d, want 4 (Snatch damage; DrawOne is a no-op on empty deck)", state.PrevTurnValue)
	}
	if len(state.Hand) != 0 {
		t.Errorf("turn 2 hand = %v, want empty (deck was empty, can't refill)", state.Hand)
	}
	if len(state.Deck) != 0 {
		t.Errorf("turn 2 deck = %v, want empty", state.Deck)
	}
	if state.ArsenalCard != nil {
		t.Errorf("turn 2 arsenal = %v, want nil (nothing Held to promote)", state.ArsenalCard)
	}
	if state.Runechants != 0 {
		t.Errorf("turn 2 runechants = %d, want 0", state.Runechants)
	}
}
