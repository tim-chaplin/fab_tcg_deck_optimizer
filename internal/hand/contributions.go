package hand

// Per-card contribution attribution: once the partition enumerator picks the winning line,
// fillContributions runs a tracked replay so every BestLine entry carries its own damage /
// block / pitch share, AttackChain surfaces the weapons that never appear on BestLine, and
// summary.State carries the winning permutation's end-of-chain CarryState.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// fillDefenseContributions writes Contribution on each Defend-role entry. The block-prevention
// share is proportional to the card's effective defense out of sumDef, capped by incomingDamage
// so over-blocking doesn't inflate attribution past what actually stopped. Effective defense is
// Defense() plus the arsenal bonus when FromArsenal is set on a card.ArsenalDefenseBonus
// implementer. Defense Reactions add their own Play return on top, evaluated against a fresh
// TurnState seeded with the turn's pitched pool and remaining deck so card effects see the same
// context the solver scored them in.
func fillDefenseContributions(line []CardAssignment, pitched []card.Card, deck []card.Card, bufs *attackBufs, sumDef, incomingDamage int) {
	prevented := sumDef
	if prevented > incomingDamage {
		prevented = incomingDamage
	}
	// Collect defenders first so each DR's Play sees the full set in state.Graveyard — mirroring
	// the seeding defenseReactionDamage does during partition enumeration.
	defenders := bufs.defendersBuf[:0]
	for i := range line {
		if line[i].Role == Defend {
			defenders = append(defenders, line[i].Card)
		}
	}
	for i := range line {
		if line[i].Role != Defend {
			continue
		}
		c := line[i].Card
		def := c.Defense()
		if line[i].FromArsenal {
			if ab, ok := c.(card.ArsenalDefenseBonus); ok {
				def += ab.ArsenalDefenseBonus()
			}
		}
		if sumDef > 0 {
			line[i].Contribution = float64(def) * float64(prevented) / float64(sumDef)
		}
		if c.Types().IsDefenseReaction() {
			bufs.defenseGravScratch = append(bufs.defenseGravScratch[:0], defenders...)
			*bufs.state = card.TurnState{Pitched: pitched, Deck: deck, Graveyard: bufs.defenseGravScratch}
			bufs.drCardStateScratch = card.CardState{Card: c, FromArsenal: line[i].FromArsenal}
			line[i].Contribution += float64(c.Play(bufs.state, &bufs.drCardStateScratch))
		}
	}
}

// fillContributions populates each BestLine entry's Contribution from the winning line:
//   - Pitch:  Card.Pitch() as resource value.
//   - Defend: proportional share of Prevented plus own Play return if a Defense Reaction.
//   - Attack: per-card damage from the winning attack-chain replay.
//   - Held / Arsenal: zero (contributed nothing this turn).
//
// Also captures the winning permutation's CarryState (end-of-chain Hand, Deck, Arsenal,
// Graveyard, Banish, Runechants, AuraTriggers) into summary.State so the deck loop can adopt
// it wholesale into next-turn state.
//
// arsenalAtChainStart is the arsenal-in card when the partition kept it in the slot, nil
// otherwise — same value bestAttackWithWeapons received during enumeration. Threading it
// through ensures the replay starts with the same s.Arsenal the winning permutation saw.
func fillContributions(summary *TurnSummary, hero hero.Hero, weapons []weapon.Weapon, swungNames []string, budget chainBudget, deck []card.Card, arsenalAtChainStart card.Card, bufs *attackBufs, incomingDamage, runechantCarryover int, priorAuraTriggers []card.AuraTrigger) {
	line := summary.BestLine

	// Reconstruct pitched, attackers, and held from the winning line. The arsenal-in entry
	// (FromArsenal=true, last slot) participates in attackers / defenders identically to hand
	// entries when its role is Attack / Defend; it can never be Held (roleAllowed bars it).
	pitched := bufs.pitchedBuf[:0]
	attackers := bufs.attackersBuf[:0]
	held := bufs.heldBuf[:0]
	arsenalInIdx := -1
	var sumDef int
	for _, a := range line {
		switch a.Role {
		case Pitch:
			pitched = append(pitched, a.Card)
		case Attack:
			if a.FromArsenal {
				arsenalInIdx = len(attackers)
			}
			attackers = append(attackers, a.Card)
		case Defend:
			def := a.Card.Defense()
			if a.FromArsenal {
				if ab, ok := a.Card.(card.ArsenalDefenseBonus); ok {
					def += ab.ArsenalDefenseBonus()
				}
			}
			sumDef += def
		case Held:
			held = append(held, a.Card)
		}
	}

	// Pitch contributions.
	for i := range line {
		if line[i].Role == Pitch {
			line[i].Contribution = float64(line[i].Card.Pitch())
		}
	}

	fillDefenseContributions(line, pitched, deck, bufs, sumDef, incomingDamage)

	chain := buildAttackChain(bufs.attackerBuf[:0], attackers, weapons, swungNames)
	if len(chain) > 0 {
		// Re-seed ctx with the winning phase split's chain-resource state so bestSequence
		// reproduces the exact permutation that won during enumeration; per-card damage
		// depends on order.
		ctx := &sequenceContext{
			hero:                hero,
			pitched:             pitched,
			deck:                deck,
			handStart:           held,
			arsenalAtChainStart: arsenalAtChainStart,
			bufs:                bufs,
			resourceBudget:      budget.resource,
			runechantCarryover:  runechantCarryover,
			incomingDamage:      incomingDamage,
			blockTotal:          sumDef,
			hasAttackPitches:    budget.hasAttackPitches,
			maxAttackPitch:      budget.maxPitch,
			arsenalInIdx:        arsenalInIdx,
			priorAuraTriggers:   priorAuraTriggers,
		}
		fillAttackChainContributions(summary, chain, ctx)
		summary.State = ctx.carryWinner
	} else {
		// No chain: synthesise a CarryState from the partition + carryover so the deck loop
		// still gets a usable next-turn seed (e.g. a hand with no attacks but a Held arsenal-in
		// card and the unchanged deck).
		summary.State = CarryState{
			Hand:         append([]card.Card(nil), held...),
			Deck:         append([]card.Card(nil), deck...),
			Arsenal:      arsenalAtChainStart,
			Runechants:   runechantCarryover,
			AuraTriggers: append([]card.AuraTrigger(nil), priorAuraTriggers...),
		}
	}
}

// buildAttackChain appends attackers first, then the weapons named in swungNames in that order,
// so the sequence search sees the same chain composition the partition loop priced. Uses the
// passed-in slice's backing array (typically bufs.attackerBuf) to stay allocation-free.
func buildAttackChain(dst []card.Card, attackers []card.Card, weapons []weapon.Weapon, swungNames []string) []card.Card {
	dst = append(dst, attackers...)
	for _, name := range swungNames {
		for _, w := range weapons {
			if w.Name() == name {
				dst = append(dst, w)
				break
			}
		}
	}
	return dst
}

// fillAttackChainContributions re-runs the sequence search with tracking enabled to recover
// the winning permutation, snapshots it into summary.AttackChain (fresh slice to avoid
// aliasing the buf-backed winnerOrder), and maps each position's damage back to BestLine's
// Attack-role entries. Weapons have no BestLine entry; their damage is already in
// summary.Value. Duplicate printings disambiguate by scan order. Contribution bundles Play
// return + hero-trigger + aura-trigger so per-card stats reflect total this-turn impact.
func fillAttackChainContributions(summary *TurnSummary, chain []card.Card, ctx *sequenceContext) {
	line := summary.BestLine
	total := len(line)
	bufs := ctx.bufs
	winnerOrder := bufs.fillContribWinnerOrder[:len(chain)]
	perCardDmg := bufs.fillContribPerCard[:len(chain)]
	if cap(bufs.fillContribTriggerDmg) < len(chain) {
		bufs.fillContribTriggerDmg = make([]float64, len(chain))
	}
	perCardTrigger := bufs.fillContribTriggerDmg[:len(chain)]
	if cap(bufs.fillContribAuraTriggerDmg) < len(chain) {
		bufs.fillContribAuraTriggerDmg = make([]float64, len(chain))
	}
	perCardAuraTrigger := bufs.fillContribAuraTriggerDmg[:len(chain)]
	ctx.bestSequence(chain, winnerOrder, perCardDmg, perCardTrigger, perCardAuraTrigger)
	summary.AttackChain = make([]AttackChainEntry, len(winnerOrder))
	for i := range winnerOrder {
		summary.AttackChain[i] = AttackChainEntry{
			Card:              winnerOrder[i],
			Damage:            perCardDmg[i],
			TriggerDamage:     perCardTrigger[i],
			AuraTriggerDamage: perCardAuraTrigger[i],
		}
	}
	if cap(bufs.fillContribUsed) < total {
		bufs.fillContribUsed = make([]bool, total)
	}
	used := bufs.fillContribUsed[:total]
	for i := range used {
		used[i] = false
	}
	for k, c := range winnerOrder {
		if _, isWeapon := c.(weapon.Weapon); isWeapon {
			continue
		}
		for i := range line {
			if used[i] || line[i].Role != Attack || line[i].Card.ID() != c.ID() {
				continue
			}
			line[i].Contribution = perCardDmg[k] + perCardTrigger[k] + perCardAuraTrigger[k]
			used[i] = true
			break
		}
	}
}
