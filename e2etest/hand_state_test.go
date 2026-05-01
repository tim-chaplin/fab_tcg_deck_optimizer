package e2etest

// These tests pin the semantics of TurnState.Hand at chain-step Play time: Hand reflects
// the cards in hand AT THIS MOMENT — committed-to-the-turn cards (pitched, used to block,
// already played, the playing card itself) are out, but cards going to be Held, played later
// in the chain, or drawn earlier in the chain stay in. Spring Load Red's "+3{p} if you have
// no cards in hand" rider is the canary — the only implemented card whose firing depends on
// a precise-at-this-moment snapshot.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that pitching a single blue and playing Spring Load fires the +3{p} rider.
func TestHandState_SpringLoadAlonePitchEmptiesHand(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{testutils.BluePitch{}, cards.SpringLoadRed{}}
	if got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue; got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Spring Load 2 + rider 3)", got)
	}
}

// Tests that a card committed to blocking counts as out-of-hand for Spring Load's rider.
func TestHandState_BlockerEmptiesHandForSpringLoad(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{testutils.BluePitch{}, cards.DodgeBlue{}, cards.SpringLoadRed{}}
	// Incoming = 3 → BluePitch pitched (3 res), Dodge played as DR for 2 prevented,
	// Spring Load resolves with empty hand. Value = 5 (Spring Load + rider) + 2 (Dodge).
	if got := d.EvalOneTurnForTesting(3, nil, hand).PrevTurnValue; got != 7 {
		t.Fatalf("PrevTurnValue = %d, want 7 (Spring Load 2 + rider 3 + Dodge 2)", got)
	}
}

// Tests that an upcoming chain step keeps Hand non-empty: only ONE Spring Load fires the rider.
func TestHandState_UpcomingChainStepBlocksFirstSpringLoadRider(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		testutils.BluePitch{},
		cards.FlyingHighBlue{},
		cards.SpringLoadRed{}, cards.SpringLoadRed{},
	}
	// Pitch BluePitch (3 res) → fund Flying High (0) + Spring Load × 2 (1 + 1).
	// Chain order [FH, SL1, SL2]: at SL1's Play, SL2 is upcoming → Hand non-empty,
	// rider blocked. At SL2's Play, hand is empty → rider fires. Value = 0 + 2 + 5 = 7.
	if got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue; got != 7 {
		t.Fatalf("PrevTurnValue = %d, want 7 (FH 0 + SL no rider 2 + SL with rider 5)", got)
	}
}

// Tests that a mid-chain draw lands in Hand: Spring Load can never fire its rider
// alongside Snatch.
func TestHandState_MidChainDrawBlocksSpringLoadRider(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		testutils.BluePitch{},
		cards.FlyingHighBlue{},
		cards.SnatchRed{},
		cards.SpringLoadRed{},
	}
	// Pitch BluePitch (3 res) → fund Flying High (0) + Snatch (0) + Spring Load (1).
	// Either chain order keeps Spring Load's rider blocked:
	//   [FH, SL, Snatch] — SL plays before Snatch, Snatch in upcoming → Hand non-empty.
	//   [FH, Snatch, SL] — Snatch hits, draws into Hand → SL sees the drawn card.
	// Damage in both: FH 0 + Snatch 4 + SL 2 = 6. (Drawn card may itself extend the
	// chain, but its presence at the moment SL resolves keeps the rider off.)
	if got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue; got != 6 {
		t.Fatalf("PrevTurnValue = %d, want 6 (FH 0 + Snatch 4 + SL no rider 2)", got)
	}
}

// Tests that a card stuck in hand (no profitable role) keeps Spring Load's rider off.
func TestHandState_HeldCardBlocksSpringLoadRider(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{testutils.BluePitch{}, cards.DodgeBlue{}, cards.SpringLoadRed{}}
	// Same hand as TestHandState_BlockerEmptiesHandForSpringLoad but with incoming = 0
	// so there's no damage for Dodge to defend against. Per the test's stated intent,
	// Dodge sits Held → hand non-empty at Spring Load's Play → rider blocked. Value = 2.
	//
	// Caveat: under the strict "used to Block → out of hand" reading, the optimizer can
	// still commit Dodge to the Defend role even with 0 incoming (Dodge's DR cost is 0,
	// so it plays for free and prevents 0). That empties the hand and fires the rider
	// for Value 5. Whether a no-op Defend assignment should still count as "stuck in
	// hand" is a sim-semantics call beyond the chain-step Hand snapshot — flagged for
	// review.
	if got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue; got != 2 {
		t.Fatalf("PrevTurnValue = %d, want 2 (Spring Load base 2; rider blocked by held Dodge)", got)
	}
}

// Tests that a to-be-pitched card stays in hand long enough for Demolition Crew's reveal.
func TestHandState_DemolitionCrewSeesUncommittedPitchInHand(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	// Hand has 4 cards + Flying High in arsenal so the chain has 5 cards available without
	// exceeding Viserai's intel=4 hand size. Optimal line: Flying High plays from arsenal
	// (granting go-again + matching-pitch +1{p} to the next attack action — Demolition Crew),
	// pitch Drag Down [Y] to play Demolition Crew (which reveals Toughen Up — the only
	// non-self cost-2 card still in hand at that moment), then pitch Toughen Up to fund
	// Brandish. Toughen Up has to be in hand at Demolition Crew's Play even though it's
	// queued to pitch later in the chain — that's the pitch-tracking semantic under test.
	// Value: 0 (Flying High) + 7 (Demolition Crew base 6 + Flying High +1{p}) + 3 (Brandish) = 10.
	hand := []sim.Card{
		cards.DemolitionCrewRed{},
		cards.ToughenUpBlue{},
		notimpl.DragDownYellow{},
		notimpl.BrandishRed{},
	}
	if got := d.EvalOneTurnForTesting(0, cards.FlyingHighRed{}, hand).PrevTurnValue; got != 10 {
		t.Fatalf("PrevTurnValue = %d, want 10 (FH 0 + DC 7 + Brandish 3 — DC reveal sees pitched Toughen Up)", got)
	}
}
