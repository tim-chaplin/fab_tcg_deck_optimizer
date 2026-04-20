package deck

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
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
		generic.SnatchRed{},
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

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (nothing on turn 1 creates runechants)", state.RunechantCarryover)
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
		generic.FlyingHighRed{},
		generic.FlyingHighRed{},
		generic.SnatchRed{},
		generic.SnatchRed{},
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

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (nothing on turn 1 creates runechants)", state.RunechantCarryover)
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
	arsenalIn := generic.SnatchRed{}
	deckCards := []card.Card{
		generic.FlyingHighRed{},
		generic.FlyingHighRed{},
		generic.SnatchRed{},
		generic.SnatchRed{},
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

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (nothing on turn 1 creates runechants)", state.RunechantCarryover)
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
	arsenalIn := generic.ToughenUpBlue{} // DR, cost 2, defense 4 — stays in arsenal with incoming 0
	deckCards := []card.Card{
		generic.SnatchRed{},
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

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (nothing on turn 1 creates runechants)", state.RunechantCarryover)
	}
}

// TestEvalOneTurn_MidTurnDrawPitchesToFundHopefulAttacker pins the PITCH disposition: a
// partition that leaves an attacker under-funded from hand pitch is still legal when a
// mid-turn-drawn card's pitch plugs the gap. The winning line plays Flying High (go again
// grant to the next attack), Snatch (fires on-hit DrawOne consuming the Blue on top of deck),
// then Amplify the Arknight — Amplify's cost 3 is paid by the drawn Blue's pitch 3. Flying
// High Yellow avoids the +1{p} matching-colour bonus against Snatch Red (pitch 2 ≠ pitch 1),
// so the attack chain totals 0 + 4 + 6 damage, with Viserai's Runechant trigger adding +1
// when Amplify (a Runeblade attack) resolves after a non-attack action (Flying High) — turn
// Value lands at 11.
//
// Starts with an explicit 3-card hand so nothing extraneous competes for the arsenal slot:
// all three cards play out, the drawn Blue is pitched (not Held), and the arsenal stays
// empty at the start of turn 2.
//
// Deck layout (consumed in source order):
//   - position 0 = the Blue beacon consumed by Snatch's DrawOne and pitched for Amplify.
//   - positions 1..3 = Blues that make up turn 2's refill.
//   - position 4 = Yellow in turn 2's last hand slot.
//   - position 5 = Yellow tripwire that should stay in the deck.
func TestEvalOneTurn_MidTurnDrawPitchesToFundHopefulAttacker(t *testing.T) {
	initialHand := []card.Card{
		generic.FlyingHighYellow{},
		generic.SnatchRed{},
		runeblade.AmplifyTheArknightRed{},
	}
	deckCards := []card.Card{
		fake.BlueAttack{}, // beacon — drawn by Snatch, pitched for Amplify
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
		fake.YellowAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, initialHand)

	// Turn 2 hand: plain refill of three Blues plus a Yellow. All three starting cards played,
	// the drawn Blue was pitched mid-chain, so nothing carries over as Held.
	wantHand := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.YellowAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", state.Hand, wantHand)
	}

	// Arsenal stays empty: no Held candidates (all hand cards played, drawn Blue was pitched).
	if state.ArsenalCard != nil {
		t.Errorf("turn 2 arsenal = %v, want nil (no Held card to promote)", state.ArsenalCard)
	}

	// Remaining deck: the untouched Yellow tripwire, then the drawn Blue recycled to the
	// bottom after it was pitched to fund Amplify's cost.
	wantDeck := []card.Card{
		fake.YellowAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	// Amplify's attack step consumes the Runechant Viserai created on its resolution (Amplify
	// is Runeblade + Flying High is a non-attack action played earlier), so no tokens carry.
	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (Amplify consumed Viserai's token on resolve)", state.RunechantCarryover)
	}

	// Turn 1 damage: 0 (Flying High Yellow — no pitch-match with Snatch Red) + 4 (Snatch) + 6
	// (Amplify) + 1 (Runechant Viserai creates on Amplify's resolve, credited at creation).
	if state.PrevTurnValue != 11 {
		t.Errorf("turn 1 Value = %d, want 11 (Flying High 0 + Snatch 4 + Amplify 6 + Viserai Runechant 1)", state.PrevTurnValue)
	}
}

// TestEvalOneTurn_MidTurnDrawExtendsChainUntilItCantContinue pins the PLAY disposition: a
// drawn card can attach itself to the end of the chain when the previous card granted Go
// again, the drawn card is free, and the card itself isn't a Defense Reaction. The solver
// keeps extending greedily — each extension fires its own Play (drawing more cards), which
// can unlock another extension — until the supply of playable drawn cards runs out.
//
// Hand is a single copy of a made-up free-cycling attack: cost 0, Go again, and fires
// DrawOne on play. Five more copies sit on top of the deck; an Aether Slash is at the
// bottom. Play the hand copy → draws copy #2 → extend → draws #3 → extend → … until all
// five deck copies have been pulled and played. The sixth draw is the Aether Slash, which
// costs 1 (no more resources, no hand pitches) and can't extend the chain. It lands in the
// arsenal via the empty-slot promotion: with every earlier drawn card already assigned
// Role=Attack, the only Held candidate is the Slash itself, so the promotion picks it.
//
// Turn 2's deck is empty (ten total cards consumed — one hand, five deck-top cantrips, the
// five cards each of those drew, covering the final Aether Slash) so the sim reports an
// empty hand with just the arsenal and prior-turn value set.
func TestEvalOneTurn_MidTurnDrawExtendsChainUntilItCantContinue(t *testing.T) {
	cantrip := fake.DrawCantrip{}
	slash := runeblade.AetherSlashRed{}
	initialHand := []card.Card{cantrip}
	deckCards := []card.Card{
		cantrip,
		cantrip,
		cantrip,
		cantrip,
		cantrip,
		slash,
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, initialHand)

	if state.ArsenalCard != slash {
		t.Errorf("turn 2 arsenal = %v, want %v (the one card that couldn't extend the chain should fill the empty slot)", state.ArsenalCard, slash)
	}

	// Deck consumed entirely: the hand cantrip plus five deck cantrips each drew a card,
	// chaining through the five refills and the final Aether Slash. Nothing was pitched, so
	// nothing recycled to the bottom either.
	if len(state.Hand) != 0 {
		t.Errorf("turn 2 hand = %v, want empty (deck should be exhausted by the extension chain)", state.Hand)
	}
	if len(state.Deck) != 0 {
		t.Errorf("turn 2 deck = %v, want empty (every card was consumed)", state.Deck)
	}

	// Six cantrip plays * 1 damage each = 6. Aether Slash never resolved. No Runechant triggers
	// (cantrip is Generic, not Runeblade, so Viserai doesn't fire).
	if state.PrevTurnValue != 6 {
		t.Errorf("turn 1 Value = %d, want 6 (six cantrip attacks)", state.PrevTurnValue)
	}

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0", state.RunechantCarryover)
	}
}

// TestEvalOneTurn_MidTurnDrawExtensionPaysFromLeftoverPitch pins the PLAY disposition for a
// paid extension: a drawn card with non-zero cost can still attach to the chain tail when the
// initial partition's leftover pitch covers it. The hand commits a blue DR as the pitch source
// up front — a partition that would otherwise fail the pitch-timing rule (leftover == max
// attack-phase pitch) is salvaged when the extension consumes that leftover.
//
// Hand: Flying High Yellow, Snatch Red, Toughen Up Blue. Deck top: Aether Slash Red.
// Optimal line: Flying High (grants Go again to Snatch) → Snatch (attack 4, fires DrawOne to
// consume the Aether Slash) → Aether Slash as extension, funded by pitching Toughen Up Blue
// (pitch 3) from the initial partition. Aether Slash deals 4 base damage; Viserai's Runechant
// trigger adds +1 (Aether Slash is a Runeblade, Flying High is a non-attack action played
// earlier). No matching-colour bonus from Flying High Yellow (pitch 2 ≠ Snatch's pitch 1).
// Total Value: 0 + 4 + 4 + 1 = 9.
//
// Deck layout:
//   - position 0 = Aether Slash Red (consumed mid-turn 1, played as extension).
//   - positions 1..3 = three Blue attacks that refill turn 2's hand behind the recycled
//     pitched Toughen Up. Not strictly needed for the behaviour under test; included so the
//     function reports a valid turn-2 hand snapshot.
func TestEvalOneTurn_MidTurnDrawExtensionPaysFromLeftoverPitch(t *testing.T) {
	initialHand := []card.Card{
		generic.FlyingHighYellow{},
		generic.SnatchRed{},
		generic.ToughenUpBlue{},
	}
	deckCards := []card.Card{
		runeblade.AetherSlashRed{}, // drawn mid-turn, paid for by the pitched Toughen Up
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, initialHand)

	// Turn 1 damage: 0 (Flying High, no pitch-match) + 4 (Snatch) + 4 (Aether Slash base —
	// Toughen Up isn't a non-attack action so the Arcane bonus doesn't fire) + 1 (Viserai's
	// Runechant trigger on Aether Slash's resolve).
	if state.PrevTurnValue != 9 {
		t.Errorf("turn 1 Value = %d, want 9 (FH 0 + Snatch 4 + Aether Slash 4 + Viserai Runechant 1)", state.PrevTurnValue)
	}

	// All three starting cards left the hand (Flying High and Snatch played, Toughen Up
	// pitched); Aether Slash was drawn and played as an extension — nothing Held for promotion.
	if state.ArsenalCard != nil {
		t.Errorf("turn 2 arsenal = %v, want nil (no Held card to promote)", state.ArsenalCard)
	}

	// Turn 2 refill: the three Blues at source positions 1..3, then the pitched Toughen Up
	// recycled to the bottom of the deck slides into slot 3.
	wantHand := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		generic.ToughenUpBlue{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", state.Hand, wantHand)
	}

	if len(state.Deck) != 0 {
		t.Errorf("turn 2 deck = %v, want empty (deck of 4 minus the drawn Slash, plus the recycled Toughen Up, exactly fills turn 2's hand)", state.Deck)
	}

	// Aether Slash's attack consumes the Runechant Viserai created on its resolve.
	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (Aether Slash consumed the token on resolve)", state.RunechantCarryover)
	}
}

// TestEvalOneTurn_MidTurnDrawSansGoAgainStaysHeld is the mirror of
// MidTurnDrawExtensionPaysFromLeftoverPitch: same shape of hand (Snatch + Toughen Up Blue) and
// deck (Aether Slash on top + refill) but without the Flying High that granted Go again to
// Snatch. Snatch's baseline Go again is false, no other card grants it one, so the chain ends
// after Snatch resolves — the Slash is drawn but sits out the rest of the turn.
//
// The partition enumerator can't pitch Toughen Up Blue here either: without an extension to
// consume the residual, the committed 3 pitch against a 0-cost chain trips the pitch-timing
// rule (leftover == maxAttackPitch). So the winning partition has Toughen Up Blue Held. Turn 1
// Value = 4 (Snatch alone). Toughen Up and the drawn Slash both land in the unified Held pool
// that feeds post-enumeration arsenal promotion; the deterministic hash picks Toughen Up for
// the arsenal slot this hand, leaving Aether Slash to carry into turn 2's hand as the Held
// prefix of the refill.
func TestEvalOneTurn_MidTurnDrawSansGoAgainStaysHeld(t *testing.T) {
	initialHand := []card.Card{
		generic.SnatchRed{},
		generic.ToughenUpBlue{},
	}
	deckCards := []card.Card{
		runeblade.AetherSlashRed{},
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

	// Toughen Up lands in the arsenal (hash pick over the {Toughen Up, Aether Slash} Held pool);
	// the Slash falls through to turn 2's hand instead.
	if _, ok := state.ArsenalCard.(generic.ToughenUpBlue); !ok {
		t.Errorf("turn 2 arsenal = %v, want Toughen Up Blue (hash picks TU from the {TU, Aether Slash} Held pool)", state.ArsenalCard)
	}

	// Turn 2 hand: the Held Aether Slash plus three fresh Blues from the deck (positions 1..3).
	wantHand := []card.Card{
		runeblade.AetherSlashRed{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", state.Hand, wantHand)
	}

	// Deck is fully consumed: 4 deck cards minus 1 Slash drawn mid-turn = 3 Blues, all in the
	// turn 2 refill alongside the Held Slash.
	if len(state.Deck) != 0 {
		t.Errorf("turn 2 deck = %v, want empty", state.Deck)
	}

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0", state.RunechantCarryover)
	}
}

// TestEvalOneTurn_TwoDrawRidersInOneChain pins a chain with two DrawOne calls — Snatch's
// on-hit plus Drawn to the Dark Dimension's on-play — so `state.Drawn` accumulates twice and
// both ends of the pipeline (the solver's draw-tracking, the sim's per-role routing, the
// fillContributions snapshot) exercise multi-entry behaviour. Flying High grants Go again to
// Snatch; Drawn to the Dark Dimension follows as the last attacker (doesn't need Go again).
//
// The winning partition leaves Toughen Up Blue HELD rather than pitching it: Drawn's cost 2
// is covered by the drawn Blue Snatch pulled (Phase 2 pitch-from-drawn, which adds 3 to
// resources), and leaving Toughen Up in hand lets the arsenal tiebreak land on an occupied
// slot. One drawn Blue is consumed as pitch; the second drawn Blue (pulled by Drawn's
// DrawOne) stays Held. Post-enumeration promotion then picks Toughen Up out of the combined
// {Toughen Up, Blue} Held pool for the arsenal.
//
// Damage: 0 (Flying High Yellow, no pitch-match) + 4 (Snatch) + 3 (Drawn to the Dark
// Dimension) + 1 (Viserai's Runechant trigger on Drawn's resolve — Flying High is a
// non-attack action played earlier) = 8.
func TestEvalOneTurn_TwoDrawRidersInOneChain(t *testing.T) {
	initialHand := []card.Card{
		generic.FlyingHighYellow{},
		generic.SnatchRed{},
		runeblade.DrawnToTheDarkDimensionRed{},
		generic.ToughenUpBlue{},
	}
	deckCards := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, initialHand)

	if state.PrevTurnValue != 8 {
		t.Errorf("turn 1 Value = %d, want 8 (FH 0 + Snatch 4 + DrawnToDark 3 + Viserai Runechant 1)", state.PrevTurnValue)
	}

	// Toughen Up Blue stays Held in hand and is promoted to arsenal; the pitched drawn Blue
	// recycles to the deck bottom, the second drawn Blue carries as Held into turn 2's hand.
	if _, ok := state.ArsenalCard.(generic.ToughenUpBlue); !ok {
		t.Errorf("turn 2 arsenal = %v, want Toughen Up Blue (Held hand card promoted over the held drawn Blue)", state.ArsenalCard)
	}

	// Turn 2 hand: the second drawn Blue (Held) plus three fresh Blues from the deck.
	wantHand := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	if !reflect.DeepEqual(state.Hand, wantHand) {
		t.Errorf("turn 2 hand = %v, want %v", state.Hand, wantHand)
	}

	// Deck bottom: the drawn Blue that was pitched mid-chain to fund Drawn to the Dark
	// Dimension, recycled to the tail of the deck like any other pitched card.
	wantDeck := []card.Card{fake.BlueAttack{}}
	if !reflect.DeepEqual(state.Deck, wantDeck) {
		t.Errorf("turn 2 deck = %v, want %v", state.Deck, wantDeck)
	}

	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0 (Drawn to the Dark Dimension consumed its own Runechant on resolve)", state.RunechantCarryover)
	}
}

// TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop pins the degenerate case: Snatch fires DrawOne
// against an empty deck and nothing happens — no panic, no spurious entry in state.Drawn, no
// index math goes sideways downstream. Hand is just Snatch; d.Cards is empty. Turn 2 can't
// deal (deck stays empty through the turn), so the sim returns a TurnStartState with just
// the previous-turn value and the arsenal/runechant carryover.
func TestEvalOneTurn_DrawOneOnEmptyDeckIsNoop(t *testing.T) {
	initialHand := []card.Card{generic.SnatchRed{}}
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
	if state.RunechantCarryover != 0 {
		t.Errorf("turn 2 runechant carryover = %d, want 0", state.RunechantCarryover)
	}
}
