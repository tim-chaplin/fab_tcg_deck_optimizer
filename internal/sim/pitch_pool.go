package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// pitchPool tracks the state of the attack-phase pitch pool during a single attack turn run:
// the active pitch ordering (perm / vals), how many pitches have popped (idx, n), the
// partially-consumed front (front + remaining), and the flat backing slice (attr) that
// per-CardState PitchedToPlay slices index into.
//
// Each pitched entry is a *CardState — the physical pitched-zone copy carrying its own
// ephemeral state — so triggertype.Pitch handlers reading TriggeringCard see the actual
// scheduled instance rather than a stateless platonic Card.
//
// Lifecycle: one pitchPool per attack turn run (per attack-permutation × pitch-permutation
// pair). playSequenceWithMeta constructs the pool from ctx.attackPitchPerm / Vals and
// drains it step-by-step via pay. At end of attack turn a pool with idx < n means a pitched
// card was held back without funding any cost — illegal in FaB. Residual `remaining`
// is fine: it's the over-pitch surplus on the last popped pitch.
type pitchPool struct {
	perm []*card.CardState
	vals []int
	idx  int
	n    int
	// front + remaining track the partially-consumed pitched card carrying over from a
	// previous attack step. Between attack steps either front is empty (front==nil &&
	// remaining==0) or one pitched card sits at the front with leftover resources.
	// Tests bypass the real pool by seeding remaining with a synthetic budget and no
	// backing front — pay then drains the budget without contributing attribution.
	front     *card.CardState
	remaining int
	attr      []*card.CardState
}

// wrapPitchedCards allocates a *CardState (Role=Pitch) for each c in cs and returns the
// pointer slice. Used by replay / print paths where bufs.pitchPcBuf isn't in scope; the
// hot partition path uses groupByRole's pitchPcBuf-backed wrapping instead.
func wrapPitchedCards(cs []card.Card) []*card.CardState {
	if len(cs) == 0 {
		return nil
	}
	states := make([]card.CardState, len(cs))
	out := make([]*card.CardState, len(cs))
	for i, c := range cs {
		states[i] = card.CardState{Card: c, Role: card.Pitch}
		out[i] = &states[i]
	}
	return out
}

// pay consumes `cost` resources from the front of the pool, popping new pitches as the
// front exhausts. Every pitched card whose resources contribute even partially to this
// payment lands in the returned slice — so pitching one 3-resource non-attack to fund
// three 1-cost plays attributes the non-attack to all three, not just the one whose
// payment popped it. Each newly popped card fires its triggertype.Pitch handlers; any
// AddResourcePoints grant folds into that card's contribution. Returns ok=false if the pool
// ran out of pitches mid-payment.
func (p *pitchPool) pay(ge *gameengine.GameEngine, cost int) (contrib []*card.CardState, ok bool) {
	attrStart := len(p.attr)
	remaining := cost
	for remaining > 0 {
		if p.front == nil && p.remaining == 0 {
			if p.idx >= p.n {
				return nil, false
			}
			p.front = p.perm[p.idx]
			p.remaining = p.vals[p.idx]
			p.idx++
			p.remaining += ge.FirePitchTriggers(p.front)
		}
		if p.front != nil {
			p.attr = append(p.attr, p.front)
		}
		if p.remaining > remaining {
			p.remaining -= remaining
			remaining = 0
		} else {
			remaining -= p.remaining
			p.remaining = 0
			p.front = nil
		}
	}
	return p.attr[attrStart:len(p.attr)], true
}
