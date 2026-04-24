package hand

// Human-readable rendering of TurnSummary: the compact one-liner FormatBestLine plus the
// sectioned play-order printout FormatBestTurn. Pure presentation layer — no solver state
// leaks in, and nothing in this file is called from the partition / sequence hot loops.

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
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

// assignmentName returns the card name, suffixed with " (from arsenal)" when the assignment came
// from the arsenal slot — that tag tells readers why the card isn't in the dealt-hand list the
// optimiser reports alongside. Used by FormatBestLine (debug one-liner) and the held/arsenal
// footer; FormatBestTurn's numbered play order attaches the arsenal tag to the role label
// instead ("PLAY from arsenal").
func assignmentName(a CardAssignment) string {
	if a.FromArsenal {
		return a.Card.Name() + " (from arsenal)"
	}
	return a.Card.Name()
}

// FormatBestLine pairs each card in BestLine with its assigned role for debug output, e.g.
// "Hocus Pocus (Blue): PITCH, Runic Reaping (Red): ATTACK". Compact one-line form; use
// FormatBestTurn for chronological play order.
func FormatBestLine(line []CardAssignment) string {
	parts := make([]string, len(line))
	for i, a := range line {
		parts[i] = assignmentName(a) + ": " + a.Role.String()
	}
	return strings.Join(parts, ", ")
}

// roleLabelWithArsenal attaches " from arsenal" to the role label when a is arsenal-in, so the
// numbered play order reads "Mauvrion Skies (Red): PLAY from arsenal" rather than tagging the
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
		parts = append(parts, fmt.Sprintf("drew %s into hand", d.Revealed.Name()))
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

// formatTriggerDamageTag joins non-zero hero and aura trigger damage into a single
// parenthesised group like "(+1 hero trigger, +1 aura trigger)" so the chain line doesn't
// stack two separate parens for what the reader sees as one composite effect block.
// Returns "" when both are zero.
func formatTriggerDamageTag(e AttackChainEntry) string {
	var parts []string
	if e.TriggerDamage > 0 {
		parts = append(parts, fmt.Sprintf("+%s hero trigger", formatContribution(e.TriggerDamage)))
	}
	if e.AuraTriggerDamage > 0 {
		parts = append(parts, fmt.Sprintf("+%s aura trigger", formatContribution(e.AuraTriggerDamage)))
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, ", ") + ")"
}

// appendAttackChainLines renders the Attack phase into playLines at 4-space indent (one level
// deeper than section headers): one numbered entry per AttackChain step in solver-chosen play
// order. nextStep advances the shared step counter so the chain's entries interleave with the
// other my-turn entries built around it. Non-weapon entries cross-reference BestLine by ID so
// arsenal-played cards get "PLAY from arsenal" / "ATTACK from arsenal" on the role label;
// weapons skip the match since they have no BestLine entry. Cards that aren't attacks
// (e.g. non-attack actions like Mauvrion Skies) use "PLAY" so the label matches what the card
// actually does on the chain. Non-zero TriggerDamage / AuraTriggerDamage merge into a single
// "(+M hero trigger, +M aura trigger)" tag so each source is attributed without fragmenting
// the line into two parenthesised groups.
func appendAttackChainLines(playLines []string, t TurnSummary, nextStep func() int) []string {
	used := make([]bool, len(t.BestLine))
	appendAttack := func(label, cardName string, e AttackChainEntry) {
		line := fmt.Sprintf("    %d. %s: %s (+%s)", nextStep(), cardName, label, formatContribution(e.Damage))
		line += formatTriggerDamageTag(e)
		playLines = append(playLines, line)
	}
	for _, e := range t.AttackChain {
		if _, isWeapon := e.Card.(weapon.Weapon); isWeapon {
			appendAttack("WEAPON ATTACK", e.Card.Name(), e)
			continue
		}
		// Match the first unused Attack-role BestLine entry by ID so we can detect FromArsenal
		// and pick the right role label.
		var matched CardAssignment
		for i := range t.BestLine {
			if used[i] || t.BestLine[i].Role != Attack || t.BestLine[i].Card.ID() != e.Card.ID() {
				continue
			}
			matched = t.BestLine[i]
			used[i] = true
			break
		}
		label := "ATTACK"
		if !e.Card.Types().Has(card.TypeAttack) {
			label = "PLAY"
		}
		appendAttack(roleLabelWithArsenal(matched, label), e.Card.Name(), e)
	}
	return playLines
}

// FormatBestTurn renders a TurnSummary as a numbered play-order list grouped into two
// sections, matching when the actions take place around the FaB turn boundary:
//
//	"My turn:"       — previous-turn AuraTrigger fires (start of action phase), then this
//	                   turn's attack-phase pitches, then the attack chain (plays, attacks,
//	                   weapon swings).
//	"Opponent's turn:" — defense-phase pitches (paying for Defense Reactions), plain blocks,
//	                   Defense Reactions. These resolve on the opponent's following turn.
//
// Section headers render at 2-space indent; numbered entries below them at 4-space indent.
// The step counter is continuous across sections so the reader can reference "step 5"
// without ambiguity. An empty section (no entries) is elided, header and all.
//
// A non-numbered `  Auras in play at start of turn: …` header precedes the sections when
// StartOfTurnAuras is non-empty or startingRunechants > 0. Runechants append to the header
// as "N Runechant(s)" since they're another piece of start-of-turn carryover state the
// reader wants to see at a glance.
//
// Held / Arsenal cards are summarized on trailing lines so the reader sees what's carrying over.
//
// Pitch-phase assignment uses a greedy split for display: smallest pitches first fund the defense
// pool until drCost is covered, the rest fund attack. The solver already validated some legal
// split exists; this picks one deterministically.
func FormatBestTurn(t TurnSummary, startingRunechants int) string {
	parts := partitionBestLineForDisplay(t.BestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	var lines []string
	if hdr := formatStartOfTurnHeader(t.StartOfTurnAuras, startingRunechants); hdr != "" {
		lines = append(lines, hdr)
	}

	step := 0
	nextStep := func() int { step++; return step }
	myTurn := buildMyTurnLines(t, attackPitches, nextStep)
	opponent := buildOpponentTurnLines(defensePitches, parts.plainBlocks, parts.defenseReactions, nextStep)

	if len(myTurn) > 0 {
		lines = append(lines, "  My turn:")
		lines = append(lines, myTurn...)
	}
	if len(opponent) > 0 {
		lines = append(lines, "  Opponent's turn:")
		lines = append(lines, opponent...)
	}
	lines = appendHeldArsenalFooter(lines, parts.held, parts.arsenal, t.Drawn)
	return strings.Join(lines, "\n")
}

// formatStartOfTurnHeader builds the "  Auras in play at start of turn: ..." line. Auras
// render by sorted name (duplicates preserved). Runechants append as "N Runechant" (singular
// for 1, plural otherwise) when startingRunechants > 0 so the reader sees both carryover
// inputs in one line. Returns "" when there's nothing to report so the caller skips the line
// entirely instead of emitting a dangling label.
func formatStartOfTurnHeader(auras []card.Card, startingRunechants int) string {
	if len(auras) == 0 && startingRunechants == 0 {
		return ""
	}
	var items []string
	if len(auras) > 0 {
		names := make([]string, len(auras))
		for i, a := range auras {
			names[i] = a.Name()
		}
		sort.Strings(names)
		items = append(items, names...)
	}
	if startingRunechants > 0 {
		noun := "Runechants"
		if startingRunechants == 1 {
			noun = "Runechant"
		}
		items = append(items, fmt.Sprintf("%d %s", startingRunechants, noun))
	}
	return "  Auras in play at start of turn: " + strings.Join(items, ", ")
}

// buildMyTurnLines returns the numbered entries for the "My turn:" section — previous-turn
// AuraTrigger fires, attack-phase pitches, and the attack chain — at 4-space indent. Shares
// the caller's step counter via nextStep so the returned lines carry globally unique step
// numbers. Zero-effect triggers (no damage, no reveal) are dropped so the section isn't
// padded with bare "(+0)" rows. Role labels route through roleLabelWithArsenal so an
// arsenal-in entry's provenance surfaces on the role label ("ATTACK from arsenal", "PLAY
// from arsenal", "PITCH from arsenal") via one rendering contract shared with the chain
// and opponent's-turn helpers.
func buildMyTurnLines(t TurnSummary, attackPitches []CardAssignment, nextStep func() int) []string {
	var lines []string
	for _, d := range t.TriggersFromLastTurn {
		suffix := formatTriggerEffect(d)
		if suffix == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("    %d. %s: %s", nextStep(), d.Card.Name(), suffix))
	}
	for _, a := range attackPitches {
		lines = append(lines, fmt.Sprintf("    %d. %s: %s",
			nextStep(), a.Card.Name(), roleLabelWithArsenal(a, "PITCH")))
	}
	return appendAttackChainLines(lines, t, nextStep)
}

// buildOpponentTurnLines returns the numbered entries for the "Opponent's turn:" section —
// defense-phase pitches (paying for Defense Reactions), plain blocks, then Defense Reactions
// — at 4-space indent. Every role label routes through roleLabelWithArsenal so arsenal
// provenance surfaces consistently on whichever role the card takes. The solver currently
// only lets arsenal-in cards reach Defense Reactions in this section (plain blocking and
// pitching from arsenal aren't legal), but the helper keeps one rendering contract for any
// future role-permission widening. Prevented damage renders as "(+N prevented)" on defense
// lines.
func buildOpponentTurnLines(defensePitches, plainBlocks, defenseReactions []CardAssignment, nextStep func() int) []string {
	var lines []string
	for _, a := range defensePitches {
		lines = append(lines, fmt.Sprintf("    %d. %s: %s",
			nextStep(), a.Card.Name(), roleLabelWithArsenal(a, "PITCH")))
	}
	for _, a := range plainBlocks {
		lines = append(lines, fmt.Sprintf("    %d. %s: %s (+%s prevented)",
			nextStep(), a.Card.Name(), roleLabelWithArsenal(a, "BLOCK"),
			formatContribution(a.Contribution)))
	}
	for _, a := range defenseReactions {
		lines = append(lines, fmt.Sprintf("    %d. %s: %s (+%s prevented)",
			nextStep(), a.Card.Name(), roleLabelWithArsenal(a, "DEFENSE REACTION"),
			formatContribution(a.Contribution)))
	}
	return lines
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

// appendHeldArsenalFooter appends the trailing "(held: ...)" / "(arsenal: ...)" lines that show
// unplayed cards outside the numbered sequence. Mid-turn-drawn Held / Arsenal cards render with
// a "(drawn)" suffix so the reader can tell them from starting-hand entries; the arsenal card
// itself shows "(stayed)" vs "(new)" so staying-in-place is distinguishable from being newly
// placed this turn.
func appendHeldArsenalFooter(lines []string, held, arsenal, drawn []CardAssignment) []string {
	var footers []string
	for _, a := range held {
		footers = append(footers, fmt.Sprintf("  (held: %s)", a.Card.Name()))
	}
	for _, d := range drawn {
		if d.Role == Held {
			footers = append(footers, fmt.Sprintf("  (held: %s (drawn))", d.Card.Name()))
		}
	}
	for _, a := range arsenal {
		label := a.Card.Name()
		if a.FromArsenal {
			label += " (stayed)"
		} else {
			label += " (new)"
		}
		footers = append(footers, fmt.Sprintf("  (arsenal: %s)", label))
	}
	for _, d := range drawn {
		if d.Role == Arsenal {
			footers = append(footers, fmt.Sprintf("  (arsenal: %s (drawn))", d.Card.Name()))
		}
	}
	return append(lines, footers...)
}
