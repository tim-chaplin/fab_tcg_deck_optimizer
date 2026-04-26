package hand

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/stubs"
)

// stubHero is the package-wide no-op hero for tests measuring raw hand value with no
// hero-ability contribution. Intel=4 matches real adult hero hand size so solver tests see
// the same draw-up depth as production.
var stubHero = stubs.Hero{Intel: 4}

// cardNames renders a slice of Card names for test failure messages. Lives here rather than
// next to the one ordering test that uses it so other tests can reuse it without a forced
// peer import.
func cardNames(cs []card.Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}

// grantAll is a test-only attacker that sets GrantedGoAgain=true on every CardState remaining in
// CardsRemaining. Used with grantSpy to detect cross-permutation CardState wrapper leakage in
// bestSequence: if grantAll runs first in one permutation, the fresh-wrapper invariant must
// keep its grants from bleeding into a later permutation where grantSpy runs first.
type grantAll struct{}

func (grantAll) ID() card.ID              { return card.Invalid }
func (grantAll) Name() string             { return "grantAll" }
func (grantAll) Cost(*card.TurnState) int { return 0 }
func (grantAll) Pitch() int               { return 0 }
func (grantAll) Attack() int              { return 0 }
func (grantAll) Defense() int             { return 0 }
func (grantAll) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (grantAll) GoAgain() bool { return true }
func (grantAll) Play(s *card.TurnState, self *card.CardState) {
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
type grantSpy struct{ saw *bool }

func (grantSpy) ID() card.ID              { return card.Invalid }
func (grantSpy) Name() string             { return "grantSpy" }
func (grantSpy) Cost(*card.TurnState) int { return 0 }
func (grantSpy) Pitch() int               { return 0 }
func (grantSpy) Attack() int              { return 0 }
func (grantSpy) Defense() int             { return 0 }
func (grantSpy) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (grantSpy) GoAgain() bool { return true }
func (g grantSpy) Play(s *card.TurnState, self *card.CardState) {
	defer s.LogPlay(self)
	if len(s.CardsPlayed) != 0 {
		return
	}
	for _, pc := range s.CardsRemaining {
		if pc.GrantedGoAgain {
			*g.saw = true
		}
	}
}
