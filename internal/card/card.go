// Package card defines the Card interface used by the simulator and basic/test implementations.
//
// The per-card CardState wrapper, the Card interface itself, and the optional markers cards
// opt into (NoMemo, VariableCost, Dominator, AddsFutureValue, ArsenalDefenseBonus, …) live in
// this file. Cohesive concern groups are split across sibling files in this package:
// types.go (CardType + TypeSet bitfield), turn_state.go (TurnState and its mutation helpers),
// triggers.go (AuraTrigger + EphemeralAttackTrigger).
package card

// CardState wraps a Card with per-turn mutable flags that other cards' effects can toggle.
// Instances are created by the solver at the start of each attack chain and live only for that
// chain. Effects that grant keywords to "the next X" scan TurnState.CardsRemaining and flip
// flags on the matching entry; the card currently resolving receives its own CardState as
// the `self` parameter to Play.
type CardState struct {
	Card Card
	// GrantedGoAgain is set by a prior card's grant ("next X attack" riders) or by the card's
	// own Play flipping self.GrantedGoAgain = true. The solver's chain-legality check ORs
	// this with Card.GoAgain().
	GrantedGoAgain bool
	// GrantedDominate is the Dominate counterpart to GrantedGoAgain: set by a prior card's
	// grant or by this card's own Play flipping self.GrantedDominate = true when a
	// conditional "gains dominate" clause fires. LikelyToHit ORs this with the card's
	// Dominator marker (HasDominate) to decide whether to credit the "can't over-block" bump.
	GrantedDominate bool
	// FromArsenal flags the single CardState whose Card came from the arsenal slot at start of
	// turn. The solver sets it before the chain runs; CardStates for hand cards and mid-turn
	// extensions stay false. Cards gate "if this is played from arsenal" riders on
	// self.FromArsenal.
	FromArsenal bool
	// BonusAttack is the +{p} this card has accumulated from prior cards' "next attack +N{p}"
	// riders. Granters set pc.BonusAttack += N on the matching CardState in CardsRemaining
	// instead of returning the bonus from their own Play return — that way the damage is
	// attributed to the attack receiving the buff, and EffectiveAttack folds it into
	// hit-likelihood checks (LikelyToHit) so a +N buff bumps a 4-power attack into the 5+
	// dominate window or a 6 into the unblockable 7. The solver applies BonusAttack to every
	// CardState's contribution unconditionally; deciding which CardStates are legal targets
	// (attack actions, weapons, future card types) is the grantor's job, not the solver's.
	// Negative bonuses (defender-side -N{p} debuffs) clamp at 0 because FaB attack power
	// can't go below 0.
	BonusAttack int
}

// EffectiveGoAgain reports whether this card has Go again this turn — from printed text or a
// grant by a prior card's effect.
func (p *CardState) EffectiveGoAgain() bool {
	return p.Card.GoAgain() || p.GrantedGoAgain
}

// EffectiveDominate reports whether this card attacks with Dominate this turn — from its
// printed Dominator marker or a grant flipping GrantedDominate (either by a prior card or by
// this card's own Play when a conditional "gains dominate" clause fires).
func (p *CardState) EffectiveDominate() bool {
	return p.GrantedDominate || HasDominate(p.Card)
}

// EffectiveAttack returns the card's printed Attack() plus any granted BonusAttack from prior
// "next attack action card gains +N{p}" riders, clamped at 0. An attack's power can't be
// reduced below 0 in FaB, so a -2 grant on a 1-power attack resolves as a 0-power attack
// (not -1). Cards with "if this hits" clauses should pass this into LikelyToHit so the rider
// fires off the post-clamp value — a +1 grant bumps a base-3 attack to 4 (the 1/4/7 likely-to-
// hit window), and a -3 grant on a 3-power attack drops it to 0 (no rider fires).
func (p *CardState) EffectiveAttack() int {
	n := p.Card.Attack() + p.BonusAttack
	if n < 0 {
		return 0
	}
	return n
}

// Hero is the minimal hero profile card effects need. Narrower than hero.Hero to avoid an
// import cycle; package simstate holds the active hero for the run.
type Hero interface {
	Name() string
	Intelligence() int
}

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static
// profile plus a Play hook for on-play logic.
type Card interface {
	// ID returns the card's canonical registry identifier. Stable within a build. Lets callers
	// key maps / slices on cards without string-hashing Name().
	ID() ID
	// Name returns the card's printed name without any pitch-color suffix — all three
	// printings of "Aether Slash" return the same string. Cards comparing by name
	// (synergies, "if you have played a card named X this turn" effects) use this directly.
	// For display, callers route through DisplayName which appends the pitch tag.
	Name() string
	// Cost returns the card's current resource cost given the turn state. Cards with a static
	// printed cost ignore s and return a constant; cards that read s (e.g. discount-per-token
	// effects) additionally implement VariableCost so the solver can pre-screen with cheap
	// MinCost / MaxCost bounds before enumerating chain permutations.
	Cost(s *TurnState) int
	Pitch() int
	// Attack is the printed attack value. Conditional bonuses belong in Play, not here.
	Attack() int
	Defense() int
	// Types returns the card's type-line descriptors as a TypeSet bitfield, e.g.
	// NewTypeSet(TypeRuneblade, TypeAction, TypeAttack).
	Types() TypeSet
	// GoAgain reports whether playing this card grants an additional action point. Cards
	// printed with "Go again" return true.
	GoAgain() bool
	// Play is called when the card resolves — as an attack or as a defense reaction. Returns
	// damage dealt to the opposing hero (may differ from Attack() after conditional bonuses) and
	// may read state to decide effects. self is the CardState wrapper for this resolution:
	// cards read self.FromArsenal for arsenal-gated riders and write self.GrantedGoAgain = true
	// to grant themselves Go again. The dispatcher folds the return into s.Value via
	// s.RecordValue; cards don't call RecordValue themselves.
	Play(s *TurnState, self *CardState) int
}

// DisplayName returns the card name with a pitch-color suffix — "Mauvrion Skies [Y]" for a
// pitch-2 yellow printing. Use anywhere a human-readable identifier needs to disambiguate
// pitch variants (log lines, deck listings, debug printouts). Pitch values outside 1-3
// (printed-pitch-zero cards, weapons, items, hero cards) fall through to the bare Name()
// with no suffix — there's no color to disambiguate.
func DisplayName(c Card) string {
	switch c.Pitch() {
	case 1:
		return c.Name() + " [R]"
	case 2:
		return c.Name() + " [Y]"
	case 3:
		return c.Name() + " [B]"
	}
	return c.Name()
}

// NoMemo is an optional marker. Cards that implement it opt out of the hand-evaluation memo —
// typically because the card's Play output depends on context (e.g. remaining deck composition)
// that the memo key doesn't capture.
type NoMemo interface {
	NoMemo()
}

// VariableCost is optionally implemented by cards whose Cost(s) varies with TurnState (e.g.
// discount-per-token effects). MinCost and MaxCost are static bounds on the Cost output across
// any state; the solver uses them for cheap O(1) pre-screens before enumerating chain
// permutations. Non-implementers must return the same value for Cost(s) regardless of s.
type VariableCost interface {
	MinCost() int
	MaxCost() int
}

// NotSilverAgeLegal is an optional marker. Cards that implement it signal they're banned in the
// Silver Age format and must be excluded from format-restricted deck pools. Source of truth is
// data_sources/silver_age_banlist.txt — keep the two in sync.
type NotSilverAgeLegal interface {
	NotSilverAgeLegal()
}

// NotImplemented is an optional marker. Cards whose printed effect references mechanics the
// simulator doesn't model (e.g. Gold / Silver / Copper token economies, Landmarks) opt in so
// random deck generation and mutation pools skip them. A NotImplemented card is still a valid
// Card — hands that already contain one still evaluate using its static Attack / Pitch /
// Defense — but the optimizer won't introduce it into a new deck or swap one in via mutation.
// Orthogonal to NotSilverAgeLegal: a card can be format-legal yet unimplemented, or banned yet
// fully implemented, or both.
type NotImplemented interface {
	NotImplemented()
}

// Dominator is an optional marker. Attack action cards printed with the Dominate keyword
// implement it; the defender is capped at one blocking card, so LikelyToHit credits the
// "slips past one block" bump at 5+ power. Conditional grants ("if X, it gains dominate")
// stay off this marker and flow through CardState.GrantedDominate instead.
type Dominator interface {
	Dominate()
}

// HasDominate reports whether c is printed with the Dominate keyword — a type assertion to
// the Dominator marker. Used by CardState.EffectiveDominate and any future scanner that
// needs the static printed-keyword check without going through a CardState.
func HasDominate(c Card) bool {
	_, ok := c.(Dominator)
	return ok
}

// LowerHealthWanter is an optional Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than their opponent (deck building, sandbagging, self-damage) opt in. Cards with a
// "less {h} than an opposing hero" rider assume the clause always fires for these heroes and never
// fires for anyone else — a coarse proxy that skips per-turn life tracking.
type LowerHealthWanter interface {
	WantsLowerHealth()
}

// AddsFutureValue is an optional marker for cards whose printed effect delivers value on a
// LATER turn rather than the one they're played — next-turn triggers, cross-turn counters,
// and the like. The solver uses it as a beatsBest tiebreaker: at equal current-turn Value
// and equal leftover-runechants, a partition that plays more AddsFutureValue cards wins,
// because their hidden future payoff isn't reflected in this turn's score. Without the
// bias, a lone future-value aura loses to Held → arsenal promotion on the arsenal-occupancy
// tiebreak.
//
// The marker is intentionally decoupled from AuraTrigger so future hidden-value mechanisms
// can opt in without piggybacking on the trigger system.
type AddsFutureValue interface {
	AddsFutureValue()
}

// ArsenalDefenseBonus is an optional marker for Defense Reactions whose printed text grants
// extra defense only when the card is played from arsenal. Implementers return the
// additional defense added to Defense() when this copy came from the arsenal slot at start
// of turn. Defense() itself stays the printed value so the hand-played path is unaffected.
type ArsenalDefenseBonus interface {
	ArsenalDefenseBonus() int
}
