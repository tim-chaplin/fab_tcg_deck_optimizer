package card

// CardState wraps a Card with per-turn mutable flags that other cards' effects can
// toggle. Instances are created by the solver at the start of each attack chain and
// live only for that chain. Effects that grant keywords to "the next X" scan
// TurnState.CardsRemaining and flip flags on the matching entry; the card currently
// resolving receives its own CardState as the `self` parameter to Play.
type CardState struct {
	Card Card
	// GrantedGoAgain is set by a prior card's grant ("next X attack" riders) or by
	// the card's own Play flipping self.GrantedGoAgain = true. The solver's chain-
	// legality check ORs this with Card.GoAgain().
	GrantedGoAgain bool
	// GrantedDominate is the Dominate counterpart to GrantedGoAgain: set by a prior
	// card's grant or by this card's own Play flipping self.GrantedDominate = true
	// when a conditional "gains dominate" clause fires. LikelyToHit ORs this with the
	// card's Dominator marker (HasDominate) to decide whether to credit the "can't
	// over-block" bump.
	GrantedDominate bool
	// FromArsenal flags the single CardState whose Card came from the arsenal slot
	// at start of turn. The solver sets it before the chain runs; CardStates for
	// hand cards and mid-turn extensions stay false. Cards gate "if this is played
	// from arsenal" riders on self.FromArsenal.
	FromArsenal bool
	// Mode is the chosen mode for a ModalCard ("Choose 1" cards), set by the chain
	// runner before Play. Always 0 for non-modal cards. Sized int8 so it packs into
	// the bool block's padding without growing the CardState — every chain step
	// reads pcBuf[i], so a wider Mode would push more cache lines through the inner
	// loop.
	Mode int8
	// BonusAttack is the +{p} this card has accumulated from prior cards' "next
	// attack +N{p}" riders. Granters set pc.BonusAttack += N on the matching
	// CardState in CardsRemaining so the damage is attributed to the attack
	// receiving the buff, and EffectiveAttack folds it into hit-likelihood checks
	// (LikelyToHit) — a +N buff can bump a 4-power attack into the 5+ dominate
	// window or a 6 into the unblockable 7. The solver applies BonusAttack to every
	// CardState's contribution unconditionally; deciding which CardStates are legal
	// targets (attack actions, weapons, future card types) is the grantor's job.
	// Negative bonuses (defender-side -N{p} debuffs) clamp at 0 because FaB attack
	// power can't go below 0.
	BonusAttack int
	// BonusDefense is the +{d} this card has accumulated from "+N{d}" rider
	// clauses, the defender-side counterpart to BonusAttack. Cross-card grants from
	// other cards and self-riders ("if X, this gains +1{d}") both write into this
	// field; EffectiveDefense folds it into the chain step's (+N) so a buffed DR's
	// block reflects the grant. Negative grants clamp at 0.
	BonusDefense int
	// PitchedToPlay is the pitched cards the chain runner attributed to paying this
	// card's resource cost during the active permutation. Populated by the chain
	// runner before each Card.Play: as costs come up, pitched cards are popped from
	// the active pitch ordering (carrying over any excess to fund subsequent cards)
	// and the popped slice is exposed here. Cards whose printed text gates on "if X
	// was pitched to play this" iterate this slice instead of the unordered
	// s.Pitched bag — the same pitched bag still lives on TurnState for cards that
	// read it as a multiset. Empty for cards whose cost was fully paid by carry from
	// a prior pitch.
	PitchedToPlay []Card
	// OnHit holds "if this hits" handlers registered during Play. Stored as struct
	// values (function pointer + small data payload) rather than closures so the hot
	// anneal path doesn't allocate per registration. See docs/dev-standards.md
	// "OnHit registrations" for the wiring contract.
	OnHit []OnHitHandler
	// SkipGraveyard is set by Play helpers that route this card to a non-graveyard
	// zone — e.g. GameEngine.RecycleToDeckBottom for Relentless Pursuit's "put this on the
	// bottom of its owner's deck" clause. The chain dispatcher's "non-persistent →
	// graveyard" append checks this flag and skips the append when set, so the
	// helper that moved the card owns its destination zone fully.
	SkipGraveyard bool
}

// OnHitHandler is one registered on-hit rider on a CardState. The chain runner fires
// Fire at finalize-active-attack time when LikelyToHit(self) is true; self is the
// buffed attack's CardState. Source names the card that registered the handler so
// log attribution stays correct when the handler was added to a different card's
// OnHit (Mauvrion Skies, Runic Reaping). N and LogText are optional small payloads
// cards use to avoid closures.
type OnHitHandler struct {
	Fire    func(s GameEngine, l Logger, self *CardState, h *OnHitHandler)
	Source  Card
	LogText string
	N       int
}

// RegisterOnHit appends a fire-only on-hit handler — the common case for "if this
// hits, do X" riders. Cards needing Source / N / LogText payloads on the handler
// append an OnHitHandler literal directly.
func (p *CardState) RegisterOnHit(fire func(s GameEngine, l Logger, self *CardState, h *OnHitHandler)) {
	p.OnHit = append(p.OnHit, OnHitHandler{Fire: fire})
}

// EffectiveGoAgain reports whether this card has Go again this turn — from printed
// text or a grant by a prior card's effect.
func (p *CardState) EffectiveGoAgain() bool {
	return p.Card.GoAgain() || p.GrantedGoAgain
}

// GrantGoAgainIfFromArsenal flips p.GrantedGoAgain when this copy came from the
// arsenal slot (p.FromArsenal). Names the standard "played-from-arsenal go again"
// rider — see the docs/dev-standards.md "Played-from-arsenal go-again" entry — so
// card Play bodies don't need to spell out the three-line if. No-op when FromArsenal
// is false; safe to call unconditionally at the top of any Play whose printed text
// reads "If <Self> is played from arsenal, it gains go again."
func (p *CardState) GrantGoAgainIfFromArsenal() {
	if p.FromArsenal {
		p.GrantedGoAgain = true
	}
}

// EffectiveAttack returns the card's printed Attack() plus any granted BonusAttack
// from prior "next attack action card gains +N{p}" riders, clamped at 0. An attack's
// power can't be reduced below 0 in FaB, so a -2 grant on a 1-power attack resolves
// as a 0-power attack (not -1). Cards with "if this hits" clauses should pass this
// into LikelyToHit so the rider fires off the post-clamp value — a +1 grant bumps a
// base-3 attack to 4 (the 1/4/7 likely-to-hit window), and a -3 grant on a 3-power
// attack drops it to 0 (no rider fires).
func (p *CardState) EffectiveAttack() int {
	n := p.Card.Attack() + p.BonusAttack
	if n < 0 {
		return 0
	}
	return n
}

// EffectiveDefense returns the card's printed Defense() plus any granted
// BonusDefense plus the ArsenalDefenseBonus when this copy came from the arsenal
// slot, clamped at 0. The sim reads this in ResolveChainStep to credit the DR's
// chain-step "(+N)" delta and to decrement IncomingDamage.
func (p *CardState) EffectiveDefense() int {
	n := p.Card.Defense() + p.BonusDefense
	if p.FromArsenal {
		n += arsenalDefenseBonusOf(p.Card)
	}
	if n < 0 {
		return 0
	}
	return n
}

// EffectiveDominate reports whether this card attacks with Dominate this turn —
// from its printed Dominator marker or a grant flipping GrantedDominate (either by
// a prior card or by this card's own Play when a conditional "gains dominate"
// clause fires).
func (p *CardState) EffectiveDominate() bool {
	return p.GrantedDominate || HasDominate(p.Card)
}

// HasDominate reports whether c is printed with the Dominate keyword — a type
// assertion to the Dominator marker. Used by CardState.EffectiveDominate and any
// future scanner that needs the static printed-keyword check without going through
// a CardState.
func HasDominate(c Card) bool {
	_, ok := c.(Dominator)
	return ok
}

// Dominator is the marker the Dominate keyword maps to. Defined here (alongside the
// HasDominate helper that consumes it) rather than in card.go because every other
// optional marker — VariableCost, ArsenalDefenseBonus, Blocker, etc. — lives in
// sim alongside the chain-runner code that consults it; Dominator is the exception
// that v2/card needs to interpret directly via the EffectiveDominate helper, so it
// stays close to its only reader.
type Dominator interface {
	Dominate()
}

// ArsenalDefenseBonus is the marker EffectiveDefense consults when this copy came
// from the arsenal slot. Returns the additional defense added to Defense() in that
// case. Defense() itself stays the printed value so the hand-played path is
// unaffected.
type ArsenalDefenseBonus interface {
	ArsenalDefenseBonus() int
}

// arsenalDefenseBonusOf returns c's ArsenalDefenseBonus contribution, or 0 when c
// doesn't implement the marker. Centralises the type assertion behind a single
// named call so every "if this came from arsenal, fold in the rider" site reads as
// one arithmetic line. Callers gate on their own from-arsenal predicate
// (CardState.FromArsenal, partition arsenal-slot index, BestLine.FromArsenal)
// before invoking — the helper does NOT decide whether the bonus applies, only how
// to extract it.
func arsenalDefenseBonusOf(c Card) int {
	if ab, ok := c.(ArsenalDefenseBonus); ok {
		return ab.ArsenalDefenseBonus()
	}
	return 0
}
