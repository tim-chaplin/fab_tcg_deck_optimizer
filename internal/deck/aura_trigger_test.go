package deck

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// damageTrigger returns a StartOfTurn AuraTrigger crediting the given damage and exhausting
// itself on the first fire. calls is bumped each time the handler runs so tests can assert
// firing counts.
func damageTrigger(self card.Card, damage int, calls *int) card.AuraTrigger {
	return card.AuraTrigger{
		Self:  self,
		Type:  card.TriggerStartOfTurn,
		Count: 1,
		Handler: func(*card.TurnState) int {
			*calls++
			return damage
		},
	}
}

// TestFireStartOfTurnTriggers_FiresEachQueuedTriggerOnce verifies every queued start-of-turn
// trigger's handler is invoked exactly once per pass, contributions are reported, and a
// trigger whose Count hits zero drops out of survivors.
func TestFireStartOfTurnTriggers_FiresEachQueuedTriggerOnce(t *testing.T) {
	aura := fake.RedAttack{}
	var callsA, callsB int
	queue := []card.AuraTrigger{damageTrigger(aura, 2, &callsA), damageTrigger(aura, 3, &callsB)}
	survivors, contribs, total, _, _ := fireStartOfTurnTriggers(queue)
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

// TestFireStartOfTurnTriggers_EmptyQueue short-circuits: no contribs, no allocation, zero total.
func TestFireStartOfTurnTriggers_EmptyQueue(t *testing.T) {
	survivors, contribs, total, runes, _ := fireStartOfTurnTriggers(nil)
	if total != 0 || runes != 0 {
		t.Errorf("total/runes = %d/%d, want 0/0", total, runes)
	}
	if contribs != nil || len(survivors) != 0 {
		t.Errorf("non-empty outputs on empty input: contribs=%v survivors=%v",
			contribs, survivors)
	}
}

// TestFireStartOfTurnTriggers_GraveyardsExhaustedAura: when a trigger's Count hits zero after
// firing, the sim moves Self into the turn-state graveyard so subsequent handlers (e.g. an
// aura with a graveyard-banish rider) see it. Asserts the contract without relying on any
// specific card to model the "look at graveyard" side.
func TestFireStartOfTurnTriggers_GraveyardsExhaustedAura(t *testing.T) {
	aura := fake.RedAttack{}
	var seen []card.Card
	// Second trigger's handler records what's currently in the graveyard so we can check the
	// first trigger's destroy happened BEFORE the second fires.
	watcher := card.AuraTrigger{
		Self:  fake.YellowAttack{},
		Type:  card.TriggerStartOfTurn,
		Count: 1,
		Handler: func(s *card.TurnState) int {
			seen = append([]card.Card(nil), s.Graveyard...)
			return 0
		},
	}
	_, _, _, _, _ = fireStartOfTurnTriggers([]card.AuraTrigger{
		{Self: aura, Type: card.TriggerStartOfTurn, Count: 1, Handler: func(*card.TurnState) int { return 0 }},
		watcher,
	})
	if len(seen) != 1 || seen[0] != aura {
		t.Errorf("second handler saw Graveyard = %v, want [%v] (first trigger's Self graveyarded first)",
			seen, aura)
	}
}

// TestEvalOneTurn_SigilOfFyendalQueuesTrigger: turn 1 starts with Sigil of Fyendal alone in
// hand. The solver plays it (the beatsBest tiebreaker prefers playing trigger-creating
// auras at equal Value over Held → arsenal promotion). Turn 2's start-of-turn pass fires
// the registered trigger: credit 1 damage-equivalent (the 1{h} gain) and graveyard the
// sigil (Count hits zero).
func TestEvalOneTurn_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := generic.SigilOfFyendalBlue{}
	deckCards := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{sigil})

	sigilPlayed := false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == card.SigilOfFyendalBlue && a.Role == hand.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", state.PrevTurnBestLine)
	}
	if state.StartOfTurnTriggerDamage != 1 {
		t.Errorf("StartOfTurnTriggerDamage = %d, want 1 (Fyendal's 1{h} gain fires at start of turn 2)",
			state.StartOfTurnTriggerDamage)
	}
	if len(state.StartOfTurnGraveyard) != 1 || state.StartOfTurnGraveyard[0].ID() != card.SigilOfFyendalBlue {
		t.Errorf("StartOfTurnGraveyard = %v, want [Sigil of Fyendal] (Count hit zero after firing)",
			state.StartOfTurnGraveyard)
	}
}
