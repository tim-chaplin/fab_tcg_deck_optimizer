package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRavenousRabble_EmptyDeckReturnsBasePower: with no deck, no card is revealed → no penalty.
func TestRavenousRabble_EmptyDeckReturnsBasePower(t *testing.T) {
	s := &card.TurnState{}
	cases := []struct {
		c    card.Card
		want int
	}{
		{RavenousRabbleRed{}, 5},
		{RavenousRabbleYellow{}, 4},
		{RavenousRabbleBlue{}, 3},
	}
	for _, tc := range cases {
		if got := tc.c.Play(s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (empty deck → base power)", tc.c.Name(), got, tc.want)
		}
	}
}

// TestRavenousRabble_TopPitchSubtracted: the top card's pitch is subtracted from base power.
func TestRavenousRabble_TopPitchSubtracted(t *testing.T) {
	cases := []struct {
		name     string
		topPitch int
		red      int // want for Red (base 5)
		yellow   int // want for Yellow (base 4)
		blue     int // want for Blue (base 3)
	}{
		{"pitch 1", 1, 4, 3, 2},
		{"pitch 2", 2, 3, 2, 1},
		{"pitch 3", 3, 2, 1, 0},
	}
	for _, tc := range cases {
		s := &card.TurnState{Deck: []card.Card{stubGenericAttackPitch(0, 0, tc.topPitch)}}
		if got := (RavenousRabbleRed{}).Play(s); got != tc.red {
			t.Errorf("%s Red: Play() = %d, want %d", tc.name, got, tc.red)
		}
		if got := (RavenousRabbleYellow{}).Play(s); got != tc.yellow {
			t.Errorf("%s Yellow: Play() = %d, want %d", tc.name, got, tc.yellow)
		}
		if got := (RavenousRabbleBlue{}).Play(s); got != tc.blue {
			t.Errorf("%s Blue: Play() = %d, want %d", tc.name, got, tc.blue)
		}
	}
}

// TestRavenousRabble_FloorsAtZero: a pitch-3 card vs Blue (base 3) would give 0, not negative.
// Verify the floor explicitly by reducing well past zero: Blue vs a (hypothetical) pitch-5 card
// should still return 0, not a negative number that'd turn into negative damage downstream.
func TestRavenousRabble_FloorsAtZero(t *testing.T) {
	s := &card.TurnState{Deck: []card.Card{stubGenericAttackPitch(0, 0, 5)}}
	if got := (RavenousRabbleBlue{}).Play(s); got != 0 {
		t.Errorf("Blue vs pitch-5 top: Play() = %d, want 0 (floor)", got)
	}
}

// TestRavenousRabble_OnlyFirstDeckCardMatters: the reveal is the top card; cards below it don't
// affect the result.
func TestRavenousRabble_OnlyFirstDeckCardMatters(t *testing.T) {
	s := &card.TurnState{Deck: []card.Card{
		stubGenericAttackPitch(0, 0, 1),
		stubGenericAttackPitch(0, 0, 3),
		stubGenericAttackPitch(0, 0, 3),
	}}
	if got := (RavenousRabbleRed{}).Play(s); got != 4 {
		t.Errorf("Play() = %d, want 4 (5 − top pitch 1, ignoring deeper cards)", got)
	}
}
