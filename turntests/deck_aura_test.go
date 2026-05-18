package turntests

// Pins that start-of-turn auras queued during turn 1 fire correctly at the turn-2
// boundary: damage credit, graveyard exhaustion, hand reveals, and OncePerTurn re-arm.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Sigil of Fyendal's start-of-turn trigger credits 1 and exhausts the sigil to graveyard.
func TestEvalTwoTurns_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := cards.SigilOfFyendalBlue{}
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras1, _ := sim.EvalTwoTurnsForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{sigil})

	if !bestLineHasRole(extras1.BestLine, ids.SigilOfFyendalBlue, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", extras1.BestLine)
	}
	if extras1.TriggerDamage != 1 {
		t.Errorf("TriggerDamage = %d, want 1 (Fyendal's 1{h} gain fires at start of turn 2)",
			extras1.TriggerDamage)
	}
	if len(extras1.TriggerGraveyard) != 1 || extras1.TriggerGraveyard[0].ID() != ids.SigilOfFyendalBlue {
		t.Errorf("TriggerGraveyard = %v, want [Sigil of Fyendal] (Count hit zero after firing)",
			extras1.TriggerGraveyard)
	}
}

// Sigil of the Arknight reveals an attack off the deck top into the turn-2 hand and the
// turn-2 chain plays it.
func TestEvalTwoTurns_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := cards.SigilOfTheArknightBlue{}
	reveal := cards.AetherSlashRed{}
	// Positions 0..3 are turn 2's normal refill (Blues); position 4 sits at the post-draw
	// top so it's the reveal target; position 5 is unused filler.
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras1, extras2 := sim.EvalTwoTurnsForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{sigil})

	if !bestLineHasRole(extras1.BestLine, ids.SigilOfTheArknightBlue, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", extras1.BestLine)
	}
	if len(extras1.TriggerGraveyard) != 1 || extras1.TriggerGraveyard[0].ID() != ids.SigilOfTheArknightBlue {
		t.Errorf("TriggerGraveyard = %v, want [Sigil of the Arknight]", extras1.TriggerGraveyard)
	}
	// Blues are pitch-only filler, so the reveal is turn 2's only attack option.
	if !bestLineHasRole(extras2.BestLine, ids.AetherSlashRed, deck.Attack) {
		t.Errorf("turn 2 BestLine didn't play the revealed Aether Slash: %+v", extras2.BestLine)
	}
}

// Blessing of Occult defers its 3-rune credit to the start-of-turn-2 tick, then exhausts.
func TestEvalTwoTurns_BlessingOfOccultCreatesRunesAtStartOfNextTurn(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras1, _ := sim.EvalTwoTurnsForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{blessing, pitch})

	if extras1.Value != 0 {
		t.Errorf("Value = %d, want 0 (Blessing's rune credit is deferred)", extras1.Value)
	}
	if !bestLineHasRole(extras1.BestLine, ids.BlessingOfOccultRed, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play Blessing as Role=Attack: %+v", extras1.BestLine)
	}
	if extras1.TriggerDamage != 3 {
		t.Errorf("TriggerDamage = %d, want 3 (Blessing's start-of-turn trigger creates 3 tokens)",
			extras1.TriggerDamage)
	}
	if len(extras1.TriggerGraveyard) != 1 || extras1.TriggerGraveyard[0].ID() != ids.BlessingOfOccultRed {
		t.Errorf("TriggerGraveyard = %v, want [Blessing [R]]", extras1.TriggerGraveyard)
	}
}

// Malefic's OncePerTurn AttackAction trigger ticks on chain attacks (one fire from Hocus
// Pocus) and stays silent at the start-of-turn boundary, surviving with Count>0.
func TestEvalTwoTurns_MaleficIncantationOncePerTurnLimitsToOneRune(t *testing.T) {
	malefic := cards.MaleficIncantationRed{}
	hocus := cards.HocusPocusRed{}
	// Filler so turn 2 can be dealt; contents don't affect the asserted invariants.
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras1, _ := sim.EvalTwoTurnsForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{malefic, hocus})

	if !bestLineHasRole(extras1.BestLine, ids.MaleficIncantationRed, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play Malefic as Role=Attack: %+v", extras1.BestLine)
	}
	if !bestLineHasRole(extras1.BestLine, ids.HocusPocusRed, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play Hocus Pocus as Role=Attack: %+v", extras1.BestLine)
	}
	if extras1.Value != 6 {
		t.Errorf("Value = %d, want 6 (3 Hocus + 1 Hocus rune + 1 Viserai trigger + 1 Malefic trigger)",
			extras1.Value)
	}
	if extras1.TriggerDamage != 0 {
		t.Errorf("TriggerDamage = %d, want 0 (Malefic only fires on attack actions)",
			extras1.TriggerDamage)
	}
	if len(extras1.TriggerGraveyard) != 0 {
		t.Errorf("TriggerGraveyard = %v, want empty (Malefic still has Count>0)",
			extras1.TriggerGraveyard)
	}
}

// Runeblood Incantation's start-of-turn trigger credits 1 per turn and survives until Count hits 0.
func TestEvalTwoTurns_RunebloodIncantationTicksAcrossTurns(t *testing.T) {
	runeblood := cards.RunebloodIncantationRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []deck.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	_, extras1, _ := sim.EvalTwoTurnsForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, []deck.Card{runeblood, pitch})

	if !bestLineHasRole(extras1.BestLine, ids.RunebloodIncantationRed, deck.Attack) {
		t.Errorf("turn 1 BestLine didn't play Runeblood as Role=Attack: %+v", extras1.BestLine)
	}
	if extras1.Value != 0 {
		t.Errorf("Value = %d, want 0 (every Runeblood rune is deferred to a future fire)",
			extras1.Value)
	}
	if extras1.TriggerDamage != 1 {
		t.Errorf("TriggerDamage = %d, want 1 (one tick per turn)", extras1.TriggerDamage)
	}
	if len(extras1.TriggerGraveyard) != 0 {
		t.Errorf("TriggerGraveyard = %v, want empty (Red has Count=3, only one tick fired)",
			extras1.TriggerGraveyard)
	}
}

// bestLineHasRole reports whether bestLine contains an entry for cardID with the given role.
func bestLineHasRole(bestLine []deck.CardAssignment, cardID ids.CardID, role deck.Role) bool {
	for _, a := range bestLine {
		if a.Card.ID() == cardID && a.Role == role {
			return true
		}
	}
	return false
}
