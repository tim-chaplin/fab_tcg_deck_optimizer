package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// cardAblation is one card's whole-deck contribution: how much mean hand value the deck loses
// when every copy of the card is removed and the deck is re-simulated. Positive means the deck
// is worse without it; near-zero or negative flags a card the deck doesn't need.
type cardAblation struct {
	name   string
	copies int
	delta  float64
}

// perCardAblation scores each unique card by coupled leave-one-out ablation (sim.PerCardAblation):
// delta = full-deck avg minus the avg with every copy of that card removed, measured on shared
// shuffles so the deltas are tight. Unlike a per-hand presence metric this is turn-agnostic and
// card-agnostic — a card whose payoff lands on a later turn still shows its worth, because pulling
// it lowers the whole deck's value.
func perCardAblation(d *deck.Deck, ev *sim.Evaluator, shuffles int, mp sim.Matchup, seed int64) []cardAblation {
	deltas := ev.PerCardAblation(d, shuffles, mp, rand.New(rand.NewSource(seed)))
	rows := make([]cardAblation, 0, len(deltas))
	for id, delta := range deltas {
		rows = append(rows, cardAblation{
			name:   registry.GetCard(id).DisplayName(),
			copies: d.Count(id),
			delta:  delta,
		})
	}
	return rows
}

// printCardAblation renders the per-card contribution table, sorted by contribution descending
// so load-bearing cards surface at the top and cut candidates (near-zero / negative) sink to the
// bottom. shuffles annotates the precision the deltas were measured at.
func printCardAblation(rows []cardAblation, shuffles int) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].delta != rows[j].delta {
			return rows[i].delta > rows[j].delta
		}
		return rows[i].name < rows[j].name
	})

	// Column widths take the larger of header label and longest data string so the dividers
	// line up. Labels fold in the copy count ("Name x2") so multi-copy contributions read clearly.
	labels := make([]string, len(rows))
	nameW := len("Card")
	for i, r := range rows {
		labels[i] = r.name
		if r.copies > 1 {
			labels[i] = fmt.Sprintf("%s x%d", r.name, r.copies)
		}
		if w := len(labels[i]); w > nameW {
			nameW = w
		}
	}
	const valHdr = "Contribution"
	valW := len(valHdr)
	for _, r := range rows {
		if w := len(fmt.Sprintf("%+.2f", r.delta)); w > valW {
			valW = w
		}
	}

	fmt.Println()
	fmt.Printf("Per-card contribution (mean hand value the deck loses without the card, %s shuffles):\n", commaInt(shuffles))
	fmt.Println()
	fmt.Printf("  %-*s | %*s\n", nameW, "Card", valW, valHdr)
	fmt.Printf("  %s-+-%s\n", strings.Repeat("-", nameW), strings.Repeat("-", valW))
	for i, r := range rows {
		fmt.Printf("  %-*s | %*s\n", nameW, labels[i], valW, fmt.Sprintf("%+.2f", r.delta))
	}
}
