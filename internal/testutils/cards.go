// Package testutils provides fake Card, Hero, and Weapon implementations shared by tests
// in multiple packages (card, cards, deck, hand, sim, turntests). The main entry point is
// the colour-and-shape-named builder constructors in fakes.go (FakeRedAttack,
// FakeBlueResource, FakeRedDR, …) plus follow-up `With...` methods so every attribute the
// test cares about is visible at the call site.
package testutils

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// GrantAll is a Runeblade attack-action card that flips GrantedGoAgain=true on every
// remaining CardState in CardsRemaining when it resolves. Used with GrantSpy to detect
// cross-permutation CardState wrapper leakage in bestSequence: the fresh-wrapper
// invariant must keep grants from bleeding across permutations.
type GrantAll struct{ Fake }

// NewGrantAll returns a configured GrantAll with the default Runeblade attack-action
// type line and Go again.
func NewGrantAll() GrantAll {
	return GrantAll{Fake: FakeRedAttack().
		WithGoAgain().
		WithTypes(card.TypeRuneblade).
		WithName("GrantAll"),
	}
}

func (GrantAll) Play(ge card.GameEngine, _ card.Logger, _ *card.CardState) {
	for _, pc := range ge.CardsRemaining() {
		pc.GrantedGoAgain = true
	}
}

// GrantSpy is a Runeblade attack-action card. When it plays first in a permutation it
// records (via *Saw) whether any CardState in CardsRemaining already has
// GrantedGoAgain=true. With per-permutation fresh wrappers that should never happen —
// no prior card in this permutation has run yet. If wrappers leak across permutations,
// a prior permutation's GrantAll Play would still be visible and the spy trips.
type GrantSpy struct {
	Fake
	Saw *bool
}

// NewGrantSpy returns a configured GrantSpy bound to saw.
func NewGrantSpy(saw *bool) GrantSpy {
	return GrantSpy{
		Fake: FakeRedAttack().
			WithGoAgain().
			WithTypes(card.TypeRuneblade).
			WithName("GrantSpy"),
		Saw: saw,
	}
}

func (g GrantSpy) Play(ge card.GameEngine, _ card.Logger, _ *card.CardState) {
	if len(ge.CardsPlayed()) != 0 {
		return
	}
	for _, pc := range ge.CardsRemaining() {
		if pc.GrantedGoAgain {
			*g.Saw = true
		}
	}
}

// CardNamesSim is the []card.Card overload of CardNames so tests holding sim slices
// don't have to widen at the call site.
func CardNamesSim(cs []card.Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}

// FireOnHitIfLikely fires every OnHit handler on self when LikelyToHit(self) — what the
// attack-turn runner does at finalize-active-attack time. Unit tests that call Card.Play
// directly use this to exercise on-hit riders without driving the full attack-turn runner.
func FireOnHitIfLikely(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.LikelyToHit(self) {
		return
	}
	for i := range self.OnHit {
		h := &self.OnHit[i]
		h.Fire(ge, l, self, h)
	}
}
