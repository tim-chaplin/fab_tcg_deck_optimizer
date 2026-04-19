package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestHitTheHighNotes_NoAuraReturnsBase: without the aura clause satisfied the +2 buff can't
// fire regardless of hit likelihood.
func TestHitTheHighNotes_NoAuraReturnsBase(t *testing.T) {
	cases := []struct {
		c    card.Card
		base int
	}{
		{HitTheHighNotesRed{}, 4},
		{HitTheHighNotesYellow{}, 3},
		{HitTheHighNotesBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
	}
}

// TestHitTheHighNotes_LikelyBuffedTotalCreditsBonus: Blue (2+2=4) is the only variant whose
// buffed total lands in the likely set — the opponent can't comfortably block 4, so the buff
// delivers and we credit it.
func TestHitTheHighNotes_LikelyBuffedTotalCreditsBonus(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	if got := (HitTheHighNotesBlue{}).Play(&s); got != 2+2 {
		t.Errorf("Blue with aura played: Play() = %d, want 4 (buffed total 4 likely to hit)", got)
	}
}

// TestHitTheHighNotes_BlockableBuffedTotalSuppressesBonus: Red (4+2=6) and Yellow (3+2=5)
// produce buffed totals the opponent comfortably blocks; the buff delivers nothing and is not
// credited.
func TestHitTheHighNotes_BlockableBuffedTotalSuppressesBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{HitTheHighNotesRed{}, 4},
		{HitTheHighNotesYellow{}, 3},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s with aura played: Play() = %d, want %d (buffed total blockable)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestHitTheHighNotes_AuraCreatedTriggersBonus: the AuraCreated flag satisfies the aura clause
// the same as a played Aura card, subject to the same hit-likelihood gate.
func TestHitTheHighNotes_AuraCreatedTriggersBonus(t *testing.T) {
	s := card.TurnState{AuraCreated: true}
	if got := (HitTheHighNotesBlue{}).Play(&s); got != 2+2 {
		t.Errorf("Blue with AuraCreated: Play() = %d, want 4", got)
	}
}

// TestHitTheHighNotes_RunechantsDontRescueBuff: the +2{p} is a physical buff; runechants firing
// alongside are a separate arcane stream and don't change whether the physical attack hits.
func TestHitTheHighNotes_RunechantsDontRescueBuff(t *testing.T) {
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}, Runechants: 1}
	if got := (HitTheHighNotesRed{}).Play(&s); got != 4 {
		t.Errorf("Red with aura + 1 Runechant: Play() = %d, want 4 (buff physical, arcane separate)", got)
	}
}
