package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestFerventForerunner_BaseGoAgainFalse pins the simplification: the only go-again trigger is
// "played from arsenal", which isn't modelled, so GoAgain() must return false. Returning true
// would let Fervent Forerunner always chain, over-crediting every sequence where it wasn't
// actually played from arsenal (which is the vast majority).
func TestFerventForerunner_BaseGoAgainFalse(t *testing.T) {
	for _, c := range []sim.Card{FerventForerunnerRed{}, FerventForerunnerYellow{}, FerventForerunnerBlue{}} {
		if c.GoAgain() {
			t.Errorf("%s: GoAgain() = true, want false (arsenal-only go-again not modelled)", c.Name())
		}
	}
}
