package hand

// Human-readable rendering of TurnSummary: the compact one-liner FormatBestLine plus the
// sectioned play-order printout FormatBestTurn. Pure presentation layer — no solver state
// leaks in, and nothing in this file is called from the partition / sequence hot loops.

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// formatContribution renders a contribution/damage value for the best-turn printout. Integers
// render bare; fractional values (proportional defense share when multiple blockers split an
// incoming attack) show one decimal place.
func formatContribution(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%.1f", v)
}

// assignmentName returns the card's display name, suffixed with " (from arsenal)" when the
// assignment came from the arsenal slot — that tag tells readers why the card isn't in the
// dealt-hand list the optimiser reports alongside. Used by FormatBestLine (debug one-liner)
// and the held/arsenal footer; FormatBestTurn's numbered play order attaches the arsenal
// tag to the role label instead ("PLAY from arsenal").
func assignmentName(a CardAssignment) string {
	if a.FromArsenal {
		return card.DisplayName(a.Card) + " (from arsenal)"
	}
	return card.DisplayName(a.Card)
}

// FormatBestLine pairs each card in BestLine with its assigned role for debug output, e.g.
// "Hocus Pocus [B]: PITCH, Runic Reaping [R]: ATTACK". Compact one-line form; use
// FormatBestTurn for chronological play order.
func FormatBestLine(line []CardAssignment) string {
	parts := make([]string, len(line))
	for i, a := range line {
		parts[i] = assignmentName(a) + ": " + a.Role.String()
	}
	return strings.Join(parts, ", ")
}

// roleLabelWithArsenal attaches " from arsenal" to the role label when a is arsenal-in, so the
// numbered play order reads "Mauvrion Skies [R]: PLAY from arsenal" rather than tagging the
// card name. Bare label otherwise.
func roleLabelWithArsenal(a CardAssignment, label string) string {
	if a.FromArsenal {
		return label + " from arsenal"
	}
	return label
}

// formatTriggerEffect renders the effect suffix for a cross-turn AuraTrigger line — the
// portion after the aura name. Damage > 0 surfaces as "START OF ACTION PHASE (+N)"; a
// non-nil Revealed card surfaces as "drew X into hand". No current card both damages and
// reveals; the comma-join handles it generically in case one is added. Returns "" when the
// trigger had no visible effect — caller drops the line entirely so a zero-impact
// reveal-capable aura (e.g. Sigil of the Arknight when the top card wasn't an attack
// action) doesn't clutter the output with a bare "(+0)".
func formatTriggerEffect(d TriggerContribution) string {
	var parts []string
	if d.Damage > 0 {
		parts = append(parts, fmt.Sprintf("START OF ACTION PHASE (+%d)", d.Damage))
	}
	if d.Revealed != nil {
		parts = append(parts, fmt.Sprintf("drew %s into hand", card.DisplayName(d.Revealed)))
	}
	return strings.Join(parts, ", ")
}

// splitPitchesByPhase assigns each pitch card to the defense or attack phase, simulating the
// order FaB prompts them in. Smallest pitches fund the defense bucket until drCost is covered;
// the rest pay for this turn's attacks. Stable on ties so display order is deterministic.
func splitPitchesByPhase(pitched []CardAssignment, drCost int) (defensePitches, attackPitches []CardAssignment) {
	sorted := append([]CardAssignment(nil), pitched...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Card.Pitch() < sorted[j].Card.Pitch() })
	covered := 0
	for _, a := range sorted {
		if covered < drCost {
			defensePitches = append(defensePitches, a)
			covered += a.Card.Pitch()
		} else {
			attackPitches = append(attackPitches, a)
		}
	}
	return defensePitches, attackPitches
}

// FormatBestTurn renders a TurnSummary's best-turn printout in one call, equivalent to
// FormatTurnLog(BuildTurnLog(t, startingRunechants)). Convenient for one-shot callers
// (tests, ad-hoc tools) that don't need to retain the TurnLog separately.
func FormatBestTurn(t TurnSummary, startingRunechants int) string {
	return FormatTurnLog(BuildTurnLog(t, startingRunechants))
}

// bestLineDisplayParts groups BestLine entries by the display section each belongs to. Pitches
// pool before being split into defense / attack phases; blocks split again by whether the card
// is a Defense Reaction (which has its own "DEFENSE REACTION" tag). drCost sums Defense-Reaction
// costs so splitPitchesByPhase can decide how much of the pitch pool funds the opponent's turn.
type bestLineDisplayParts struct {
	pitched          []CardAssignment
	plainBlocks      []CardAssignment
	defenseReactions []CardAssignment
	held             []CardAssignment
	arsenal          []CardAssignment
	drCost           int
}

// partitionBestLineForDisplay sorts the winning line into the buckets FormatBestTurn renders
// section-by-section. Defenders split on DR membership so DR-only lines get the right label
// and their cost contributes to the defense-phase pitch target.
func partitionBestLineForDisplay(line []CardAssignment) bestLineDisplayParts {
	var parts bestLineDisplayParts
	zeroState := &card.TurnState{}
	for _, a := range line {
		switch a.Role {
		case Pitch:
			parts.pitched = append(parts.pitched, a)
		case Attack:
			// Attack-phase cost sum is computed here to match the turn's modeling, but is not
			// surfaced in the rendered output — the attack chain's per-card lines already show
			// damage credit rather than cost.
			_ = a.Card.Cost(zeroState)
		case Defend:
			if a.Card.Types().IsDefenseReaction() {
				parts.drCost += a.Card.Cost(zeroState)
				parts.defenseReactions = append(parts.defenseReactions, a)
			} else {
				parts.plainBlocks = append(parts.plainBlocks, a)
			}
		case Held:
			parts.held = append(parts.held, a)
		case Arsenal:
			parts.arsenal = append(parts.arsenal, a)
		}
	}
	return parts
}
