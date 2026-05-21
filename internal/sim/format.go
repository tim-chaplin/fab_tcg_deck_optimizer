package sim

// Compact one-liner FormatBestLine plus the partition helpers that slot the winning
// BestLine into display sections. Pure presentation, off the hot loops.

import (
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// assignmentName returns the card's display name, suffixed with " (from arsenal)" when
// the assignment came from the arsenal slot.
func assignmentName(a card.CardAssignment) string {
	if a.FromArsenal {
		return a.Card.DisplayName() + " (from arsenal)"
	}
	return a.Card.DisplayName()
}

// FormatBestLine pairs each card in BestLine with its assigned role for debug output, e.g.
// "Hocus Pocus [B]: PITCH, Runic Reaping [R]: ATTACK". Compact one-line form.
func FormatBestLine(line []card.CardAssignment) string {
	parts := make([]string, len(line))
	for i, a := range line {
		parts[i] = assignmentName(a) + ": " + a.Role.String()
	}
	return strings.Join(parts, ", ")
}

// splitPitchesByPhase assigns each pitch card to the defense or attack phase, simulating the
// order FaB prompts them in. Smallest pitches fund the defense bucket until drCost is covered;
// the rest pay for this turn's attacks. Stable on ties so display order is deterministic.
func splitPitchesByPhase(pitched []card.CardAssignment, drCost int) (defensePitches, attackPitches []card.CardAssignment) {
	sorted := append([]card.CardAssignment(nil), pitched...)
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

// bestLineDisplayParts groups BestLine entries by display section. Pitches pool before being
// split into defense / attack phases; blocks split again by whether the card is a Defense
// Reaction (its own tag). drCost sums DR costs so splitPitchesByPhase can decide how much
// of the pitch pool funds the opponent's turn.
type bestLineDisplayParts struct {
	pitched          []card.CardAssignment
	plainBlocks      []card.CardAssignment
	defenseReactions []card.CardAssignment
	held             []card.CardAssignment
	arsenal          []card.CardAssignment
	drCost           int
}

// partitionBestLineForDisplay sorts the winning line into the buckets the printout renders
// section by section. Defenders split on DR membership so DR-only lines get the right label
// and their cost contributes to the defense-phase pitch target.
func partitionBestLineForDisplay(line []card.CardAssignment) bestLineDisplayParts {
	var parts bestLineDisplayParts
	zeroState := gameengine.New()
	for _, a := range line {
		switch a.Role {
		case card.Pitch:
			parts.pitched = append(parts.pitched, a)
		case card.Attack:
			_ = a.Card.Cost(zeroState)
		case card.Defend:
			if a.Card.Types(nil).IsDefenseReaction() {
				parts.drCost += a.Card.Cost(zeroState)
				parts.defenseReactions = append(parts.defenseReactions, a)
			} else {
				parts.plainBlocks = append(parts.plainBlocks, a)
			}
		case card.Held:
			parts.held = append(parts.held, a)
		case card.Arsenal:
			parts.arsenal = append(parts.arsenal, a)
		}
	}
	return parts
}
