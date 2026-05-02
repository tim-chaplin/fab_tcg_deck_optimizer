package sim

// evaluatePartition is the shared "given a fixed role assignment, score it" body used by
// both findBest's recurse leaf (one of many partitions explored during the search) and
// replayBest (the cached partition replayed without searching). It groups hand cards into
// pitched/attackers/defenders/held, folds the arsenal-in card into the right bucket per
// rolesBuf[n], computes the arsenal indices, and forwards everything to
// bestAttackWithWeapons. The output tuple is bestAttackWithWeapons's tuple plus the
// computed arsenalAtChainStart so callers wiring up TurnSummary don't recompute it.
//
// rolesBuf must be in sync with hand (rolesBuf[i] is the role assigned to hand[i]) and
// must use only Pitch/Attack/Defend/Held in hand-card slots — hand cards never have
// Arsenal role during the chain run. The cache-replay caller, whose stored BestLine may
// contain a post-hoc-promoted hand entry tagged Arsenal, flips that entry back to Held
// before calling and restores the Arsenal tag on the returned BestLine afterward.
//
// Mutates the bufs scratch slices (pitchedBuf, attackersBuf, defendersBuf, heldBuf) in
// place; both callers feed pooled scratch through bufs and tolerate the rewrite.
func (e *Evaluator) evaluatePartition(
	hero Hero, weapons []Weapon, hand []Card,
	deck []Card, arsenalCardIn Card,
	rolesBuf []Role, n int, bufs *attackBufs,
	runechantCarryover int, mp Matchup, defenseSum int,
	priorAuraTriggers []AuraTrigger, skipLog bool,
) (
	attackDealt, defenseDealt, leftoverRunechants int,
	swung []string, carry CarryState,
	ok, cacheable bool,
	arsenalAtChainStart Card,
) {
	// Group hand cards into played / pitched / defending buckets, then fold in the
	// arsenal-in card based on its slot's role.
	p, a, d := groupByRoleInto(
		hand, rolesBuf[:n],
		bufs.pitchedBuf[:0], bufs.attackersBuf[:0], bufs.defendersBuf[:0],
	)
	if arsenalCardIn != nil {
		switch rolesBuf[n] {
		case Attack:
			a = append(a, arsenalCardIn)
		case Defend:
			d = append(d, arsenalCardIn)
		}
	}
	// Arsenal-in is appended last to a / d above, so its index is len(slice)-1 when
	// present. -1 means no arsenal-in card in that bucket (either no arsenal-in card at
	// all, or it took a different role).
	arsenalInIdx := -1
	if arsenalCardIn != nil && rolesBuf[n] == Attack {
		arsenalInIdx = len(a) - 1
	}
	arsenalDefenderIdx := -1
	if arsenalCardIn != nil && rolesBuf[n] == Defend {
		arsenalDefenderIdx = len(d) - 1
	}
	// Held cards (hand cards left without a Pitch / Attack / Defend role) thread through
	// to TurnState.Hand so alt-cost effects (e.g. Moon Wish's "use a hand card") can read
	// len > 0. Arsenal-in can never be Held (roleAllowed bars it).
	h := gatherHeldCards(hand, rolesBuf[:n], bufs.heldBuf[:0])
	arsenalAtChainStart = findArsenalCard(rolesBuf, arsenalCardIn, n)

	// Hand off to the chain dispatcher — same call shape both callers used inline before
	// the extraction.
	attackDealt, defenseDealt, leftoverRunechants, _, swung, carry, ok, cacheable = bestAttackWithWeapons(
		hero, weapons, a, d, p, h, deck, bufs,
		runechantCarryover, mp, defenseSum,
		arsenalInIdx, arsenalDefenderIdx, arsenalAtChainStart,
		priorAuraTriggers, skipLog,
	)
	return
}
