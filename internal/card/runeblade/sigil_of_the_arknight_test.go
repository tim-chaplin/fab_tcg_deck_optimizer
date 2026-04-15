package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

// stubHero4 is a minimal card.Hero with Intelligence 4 for tests.
type stubHero4 struct{}

func (stubHero4) Name() string       { return "stubHero4" }
func (stubHero4) Intelligence() int  { return 4 }

func TestSigilOfTheArknight_EmptyDeckReturnsZero(t *testing.T) {
	// No deck → can't reach the reveal index, so 0. AuraCreated still flips.
	simstate.CurrentHero = stubHero4{}
	defer func() { simstate.CurrentHero = nil }()
	var s card.TurnState
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() with empty deck = %d, want 0", got)
	}
	if !s.AuraCreated {
		t.Errorf("Play() did not set AuraCreated")
	}
}

func TestSigilOfTheArknight_RevealsAttackActionAtIntelligenceIndex(t *testing.T) {
	// Intelligence=4: first 4 cards go to next hand; the card at index 4 is revealed. Here that's
	// an attack action → +3.
	simstate.CurrentHero = stubHero4{}
	defer func() { simstate.CurrentHero = nil }()
	deck := []card.Card{
		stubNonAttack{}, stubNonAttack{}, stubNonAttack{}, stubNonAttack{},
		stubRunebladeAttack{},
	}
	s := card.TurnState{Deck: deck}
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3 (attack action at deck[Intelligence])", got)
	}
}

func TestSigilOfTheArknight_RevealsNonAttackAtIntelligenceIndex(t *testing.T) {
	// Deck[Intelligence] is a non-attack card → 0, even though attack actions sit elsewhere in
	// the deck.
	simstate.CurrentHero = stubHero4{}
	defer func() { simstate.CurrentHero = nil }()
	deck := []card.Card{
		stubRunebladeAttack{}, stubRunebladeAttack{}, stubRunebladeAttack{}, stubRunebladeAttack{},
		stubAura{},
	}
	s := card.TurnState{Deck: deck}
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (revealed card at index 4 is non-attack)", got)
	}
}

func TestSigilOfTheArknight_DeckTooShortReturnsZero(t *testing.T) {
	// Intelligence=4 but only 3 cards remain — can't reach the reveal index → 0.
	simstate.CurrentHero = stubHero4{}
	defer func() { simstate.CurrentHero = nil }()
	deck := []card.Card{stubRunebladeAttack{}, stubRunebladeAttack{}, stubRunebladeAttack{}}
	s := card.TurnState{Deck: deck}
	if got := (SigilOfTheArknightBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (deck too short)", got)
	}
}

func TestSigilOfTheArknight_NoMemoMarker(t *testing.T) {
	// Sigil's Play depends on deck composition, so it must opt out of the hand-evaluation memo.
	var c card.Card = SigilOfTheArknightBlue{}
	if _, ok := c.(card.NoMemo); !ok {
		t.Errorf("SigilOfTheArknightBlue should implement card.NoMemo")
	}
}
