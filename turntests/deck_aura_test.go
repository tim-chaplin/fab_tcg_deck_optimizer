package turntests

// Pins start-of-turn aura behaviour at the turn-2 boundary: damage credit, graveyard
// exhaustion, hand reveals, OncePerTurn re-arm. Turn 2 usually gets an empty deck/hand so
// the attack turn produces 0 value and turn2.Value isolates the trigger's contribution.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Sigil of Fyendal's start-of-turn trigger credits 1{h} and exhausts to graveyard.
func TestEvalTwoTurns_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := cards.SigilOfFyendalBlue{}
	d := deck.New(heroes.Viserai, nil, nil)
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{sigil})

	if !bestLineHasRole(turn1.BestLine, cards.SigilOfFyendalBlue{}, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", turn1.BestLine)
	}
	if turn2.Value != 1 {
		t.Errorf("turn 2 Value = %d, want 1 (Fyendal's 1{h} gain fires at start of turn 2; attack turn has nothing to play)",
			turn2.Value)
	}
	if !graveyardContains(turn2.State.Graveyard(), cards.SigilOfFyendalBlue{}) {
		t.Errorf("turn 2 graveyard = %v, want it to contain Sigil of Fyendal (Count hit zero after firing)",
			turn2.State.Graveyard())
	}
}

// Sigil of the Arknight reveals an attack off the deck top into the turn-2 hand and the
// turn-2 attack turn plays it. Deck stacks 4 pitches (turn 2's refill) above the reveal so
// Sigil's PopDeckTop pulls the attack.
func TestEvalTwoTurns_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := cards.SigilOfTheArknightBlue{}
	reveal := cards.AetherSlashRed{}
	d := deck.New(heroes.Viserai, nil, []deck.Card{
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
		reveal,
	})
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{sigil})

	if !bestLineHasRole(turn1.BestLine, cards.SigilOfTheArknightBlue{}, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", turn1.BestLine)
	}
	if !bestLineHasRole(turn2.BestLine, cards.AetherSlashRed{}, card.Attack) {
		t.Errorf("turn 2 BestLine didn't play the revealed Aether Slash: %+v", turn2.BestLine)
	}
	if !graveyardContains(turn2.State.Graveyard(), cards.SigilOfTheArknightBlue{}) {
		t.Errorf("turn 2 graveyard = %v, want it to contain Sigil of the Arknight",
			turn2.State.Graveyard())
	}
}

// Blessing of Occult defers its 3-rune credit to the start-of-turn-2 tick, then exhausts.
// Turn 1 needs a pitch card to fund Blessing's cost.
func TestEvalTwoTurns_BlessingOfOccultCreatesRunesAtStartOfNextTurn(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	pitch := testutils.FakeBlueResource()
	d := deck.New(heroes.Viserai, nil, nil)
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{blessing, pitch})

	if turn1.Value != 0 {
		t.Errorf("turn 1 Value = %d, want 0 (Blessing's rune credit is deferred)", turn1.Value)
	}
	if !bestLineHasRole(turn1.BestLine, cards.BlessingOfOccultRed{}, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play Blessing as Role=Attack: %+v", turn1.BestLine)
	}
	if turn2.Value != 3 {
		t.Errorf("turn 2 Value = %d, want 3 (Blessing's trigger creates 3 runechants; attack turn has nothing to play)",
			turn2.Value)
	}
	if !graveyardContains(turn2.State.Graveyard(), cards.BlessingOfOccultRed{}) {
		t.Errorf("turn 2 graveyard = %v, want it to contain Blessing", turn2.State.Graveyard())
	}
}

// Malefic's OncePerTurn AttackAction trigger ticks on attack-turn attacks (one fire from Hocus
// Pocus on turn 1) and stays silent at the start-of-turn boundary, surviving with Count>0.
func TestEvalTwoTurns_MaleficIncantationOncePerTurnLimitsToOneRune(t *testing.T) {
	malefic := cards.MaleficIncantationRed{}
	hocus := cards.HocusPocusRed{}
	d := deck.New(heroes.Viserai, nil, nil)
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), []card.Card{malefic, hocus})

	if !bestLineHasRole(turn1.BestLine, cards.MaleficIncantationRed{}, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play Malefic as Role=Attack: %+v", turn1.BestLine)
	}
	if !bestLineHasRole(turn1.BestLine, cards.HocusPocusRed{}, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play Hocus Pocus as Role=Attack: %+v", turn1.BestLine)
	}
	if turn1.Value != 6 {
		t.Errorf("turn 1 Value = %d, want 6 (3 Hocus + 1 Hocus rune + 1 Viserai trigger + 1 Malefic trigger)",
			turn1.Value)
	}
	if turn2.Value != 0 {
		t.Errorf("turn 2 Value = %d, want 0 (Malefic only fires on attack actions; attack turn has nothing to play)",
			turn2.Value)
	}
	if graveyardContains(turn2.State.Graveyard(), cards.MaleficIncantationRed{}) {
		t.Errorf("turn 2 graveyard = %v, want it to NOT contain Malefic (still has Count>0)",
			turn2.State.Graveyard())
	}
}

// Runeblood Incantation's start-of-turn trigger credits 1 per turn and survives until Count hits 0.
// Turn 1 needs a pitch card to fund Runeblood's cost.
func TestEvalTwoTurns_RunebloodIncantationTicksAcrossTurns(t *testing.T) {
	runeblood := cards.RunebloodIncantationRed{}
	pitch := testutils.FakeBlueResource()
	d := deck.New(heroes.Viserai, nil, nil)
	turn1, turn2 := sim.EvalTwoTurnsForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), []card.Card{runeblood, pitch})

	if !bestLineHasRole(turn1.BestLine, cards.RunebloodIncantationRed{}, card.Attack) {
		t.Errorf("turn 1 BestLine didn't play Runeblood as Role=Attack: %+v", turn1.BestLine)
	}
	if turn1.Value != 0 {
		t.Errorf("turn 1 Value = %d, want 0 (every Runeblood rune is deferred to a future fire)",
			turn1.Value)
	}
	if turn2.Value != 1 {
		t.Errorf("turn 2 Value = %d, want 1 (one tick per turn; attack turn has nothing to play)", turn2.Value)
	}
	if graveyardContains(turn2.State.Graveyard(), cards.RunebloodIncantationRed{}) {
		t.Errorf("turn 2 graveyard = %v, want it to NOT contain Runeblood (Red has Count=3, only one tick fired)",
			turn2.State.Graveyard())
	}
}

// bestLineHasRole reports whether bestLine contains an entry for the given card (matched
// by Name) with the given role.
func bestLineHasRole(bestLine []card.CardAssignment, want card.Card, role card.Role) bool {
	for _, a := range bestLine {
		if a.Card.Name() == want.Name() && a.Role == role {
			return true
		}
	}
	return false
}

// graveyardContains reports whether the graveyard slice contains a card matching want by
// Name.
func graveyardContains(grav []card.Card, want card.Card) bool {
	for _, c := range grav {
		if c.Name() == want.Name() {
			return true
		}
	}
	return false
}

// containsSameCards reports whether got and want hold the same multiset of cards (matched
// by Name, order-insensitive). Use this for assertions on hands / decks / graveyards where
// the test cares which cards ended up in the zone, not what order they were sorted into.
func containsSameCards(got, want []card.Card) bool {
	if len(got) != len(want) {
		return false
	}
	counts := make(map[string]int, len(want))
	for _, c := range want {
		counts[c.Name()]++
	}
	for _, c := range got {
		counts[c.Name()]--
		if counts[c.Name()] < 0 {
			return false
		}
	}
	return true
}
