package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestConditionalGoAgain_NotGrantedByDefault pins the baseline GoAgain() return for cards whose
// printed go-again is gated on a condition we don't (or can't) model. Returning `true` from
// GoAgain() would make the keyword unconditional in the solver, overstating these cards' tempo.
// Cards that model the condition via TurnState.Self.GrantedGoAgain (Zealous Belting) also belong
// here: their GoAgain() must be false so EffectiveGoAgain doesn't OR-dominate the gated grant.
func TestConditionalGoAgain_NotGrantedByDefault(t *testing.T) {
	cards := []card.Card{
		SunKissRed{}, SunKissYellow{}, SunKissBlue{},
		OutMuscleRed{}, OutMuscleYellow{}, OutMuscleBlue{},
		ZealousBeltingRed{}, ZealousBeltingYellow{}, ZealousBeltingBlue{},
	}
	for _, c := range cards {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (printed keyword is conditional)", c.Name())
		}
	}
}
