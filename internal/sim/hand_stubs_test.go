package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// heroStub is a minimal sim.Hero used by the package-private test fixtures in this file.
// Inlined (rather than reusing testutils.Hero) because testutils itself imports sim, so a
// sim test importing testutils would form a cycle in the test binary.
type heroStub struct{ intel int }

func (heroStub) ID() ids.HeroID                    { return ids.InvalidHero }
func (heroStub) Name() string                      { return "heroStub" }
func (heroStub) Health() int                       { return 20 }
func (h heroStub) Intelligence() int               { return h.intel }
func (heroStub) Types() card.TypeSet               { return 0 }
func (heroStub) OnCardPlayed(Card, *TurnState) int { return 0 }

// stubHero is the package-wide no-op hero for tests measuring raw hand value with no
// hero-ability contribution. Intel=4 matches real adult hero hand size so solver tests see
// the same draw-up depth as production.
var stubHero = heroStub{intel: 4}

// moonWishHero is the no-op hero used by the Moon Wish e2e tests so the assertions on
// Value / Hand / Deck / Arsenal reflect Moon Wish's plumbing alone (no Viserai-style
// OnCardPlayed runechant credit perturbing the numbers). Lives here (rather than next to
// the Moon Wish tests, which are in package sim_test) so exports_test.go can re-export it
// as MoonWishHero for the dot-importing test files.
var moonWishHero = heroStub{intel: 4}

// cardNames renders a slice of Card names for test failure messages. Lives here rather than
// next to the one ordering test that uses it so other tests can reuse it without a forced
// peer import.
func cardNames(cs []Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}

// instantStub is a 0-cost, 0-power Generic Action - Instant card with no Go again. Tests
// chain-runner behaviour around the Action Point debit: an Instant after a non-Go-again
// card should still resolve because Instants cost 0 AP. ID is a distinct fake-card slot
// so cardMetaCache doesn't share its meta with another stub at slot 0.
type instantStub struct{}

func (instantStub) ID() ids.CardID      { return ids.FakeInstant }
func (instantStub) Name() string        { return "instantStub" }
func (instantStub) Cost(*TurnState) int { return 0 }
func (instantStub) Pitch() int          { return 0 }
func (instantStub) Attack() int         { return 0 }
func (instantStub) Defense() int        { return 0 }
func (instantStub) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeInstant)
}
func (instantStub) GoAgain() bool                      { return false }
func (instantStub) Play(s *TurnState, self *CardState) { s.LogPlay(self) }

// noGoAgainAttackStub is a 0-cost, 1-power attack action card with no Go again. Tests
// chain-runner behaviour after the AP pool runs out: a non-Instant follow-up should be
// rejected.
type noGoAgainAttackStub struct{}

func (noGoAgainAttackStub) ID() ids.CardID      { return ids.FakeNoGoAgainAttack }
func (noGoAgainAttackStub) Name() string        { return "noGoAgainAttack" }
func (noGoAgainAttackStub) Cost(*TurnState) int { return 0 }
func (noGoAgainAttackStub) Pitch() int          { return 0 }
func (noGoAgainAttackStub) Attack() int         { return 1 }
func (noGoAgainAttackStub) Defense() int        { return 0 }
func (noGoAgainAttackStub) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (noGoAgainAttackStub) GoAgain() bool { return false }
func (noGoAgainAttackStub) Play(s *TurnState, self *CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

// grantAll is a test-only attacker that sets GrantedGoAgain=true on every CardState remaining in
// CardsRemaining. Used with grantSpy to detect cross-permutation CardState wrapper leakage in
// bestSequence: if grantAll runs first in one permutation, the fresh-wrapper invariant must
// keep its grants from bleeding into a later permutation where grantSpy runs first.
type grantAll struct{}

func (grantAll) ID() ids.CardID      { return ids.InvalidCard }
func (grantAll) Name() string        { return "grantAll" }
func (grantAll) Cost(*TurnState) int { return 0 }
func (grantAll) Pitch() int          { return 0 }
func (grantAll) Attack() int         { return 0 }
func (grantAll) Defense() int        { return 0 }
func (grantAll) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (grantAll) GoAgain() bool { return true }
func (grantAll) Play(s *TurnState, self *CardState) {
	for _, pc := range s.CardsRemaining {
		pc.GrantedGoAgain = true
	}
	s.LogPlay(self)
}

// grantSpy is a test-only attacker that, when it plays FIRST in a permutation, records whether
// any CardState in CardsRemaining already has GrantedGoAgain=true. With per-permutation fresh
// wrappers, that should never happen (no prior card in this permutation has run yet). If wrappers
// leak across permutations, a grant applied by a previous permutation's grantAll will still be
// visible here — tripping the spy.
type grantSpy struct{ Saw *bool }

func (grantSpy) ID() ids.CardID      { return ids.InvalidCard }
func (grantSpy) Name() string        { return "grantSpy" }
func (grantSpy) Cost(*TurnState) int { return 0 }
func (grantSpy) Pitch() int          { return 0 }
func (grantSpy) Attack() int         { return 0 }
func (grantSpy) Defense() int        { return 0 }
func (grantSpy) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (grantSpy) GoAgain() bool { return true }
func (g grantSpy) Play(s *TurnState, self *CardState) {
	defer s.LogPlay(self)
	if len(s.CardsPlayed) != 0 {
		return
	}
	for _, pc := range s.CardsRemaining {
		if pc.GrantedGoAgain {
			*g.Saw = true
		}
	}
}
