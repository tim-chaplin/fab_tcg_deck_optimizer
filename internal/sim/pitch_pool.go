package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// pitchPool tracks the state of the attack-phase pitch pool during a single chain run:
// the active pitch ordering (perm / vals), how many pitches have popped (idx, n), the
// partially-consumed front (front + remaining), and the flat backing slice (attr) that
// per-CardState PitchedToPlay slices index into.
//
// Lifecycle: one pitchPool per chain run (per attack-permutation × pitch-permutation
// pair). playSequenceWithMeta constructs the pool from ctx.attackPitchPerm / Vals and
// drains it step-by-step via pay. At end of chain a pool with idx < n means a pitched
// card was held back without funding any cost — illegal in FaB. Residual `remaining`
// is fine: it's the over-pitch surplus on the last popped pitch.
type pitchPool struct {
	perm []card.Card
	vals []int
	idx  int
	n    int
	// front + remaining track the partially-consumed pitched card carrying over from a
	// previous chain step. Between chain steps either front is empty (front==nil &&
	// remaining==0) or one pitched card sits at the front with leftover resources.
	// Tests bypass the real pool by seeding remaining with a synthetic budget and no
	// backing front — pay then drains the budget without contributing attribution.
	front     card.Card
	remaining int
	attr      []card.Card
}

// pay consumes `cost` resources from the front of the pool, popping new pitches as the
// front exhausts. Every pitched card whose resources contribute even partially to this
// payment lands in the returned slice — so pitching one 3-resource non-attack to fund
// three 1-cost plays attributes the non-attack to all three, not just the one whose
// payment popped it. Each newly popped card fires its triggertype.Pitch handlers, whose
// AddPitchBonus grants fold into that card's contribution. Returns ok=false if the pool
// ran out of pitches mid-payment.
func (p *pitchPool) pay(ge *gameengine.GameEngine, cost int) (contrib []card.Card, ok bool) {
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
