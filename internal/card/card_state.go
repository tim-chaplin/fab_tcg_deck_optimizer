package card

// CardState wraps a Card with per-turn mutable flags that other cards' effects can toggle.
// Instances are created by the solver at the start of each attack turn and live only for
// that attack turn. Effects that grant keywords to "the next X" scan TurnState.CardsRemaining and
// flip flags on the matching entry; the card currently resolving receives its own CardState
// as the `self` parameter to Play.
//
// Fields split into two groups: identity (Card / Role / FromArsenal / Mode) that's stable
// across the card's play, and Ephemeral (anonymously embedded) — the granted keywords,
// accumulated bonuses, registered handlers, and pitch attribution that build up or get
// assigned around the card's play. Ephemeral.Reset zeroes the whole group in one statement
// so a new mutable field can't silently carry over when the slot is reused.
type CardState struct {
	Card Card
	// Role is the card's partition-assigned role for this turn — Pitch / Attack / Defend /
	// Held / Arsenal. Discard removes a Held-role entry from the hand.
	Role Role
	// FromArsenal flags the single CardState whose Card came from the arsenal slot at start
	// of turn. Cards gate "if this is played from arsenal" riders on self.FromArsenal.
	FromArsenal bool
	// FromDraw flags entries added to the hand by a mid-attack-turn DrawOne. The optimizer
	// treats drawn cards as a black box: they count toward HandSize but their identity
	// is hidden from in-attack-turn attribute reads (HandHasMatching / HeldHandSize skip them).
	FromDraw bool
	// Mode is the chosen mode for a Modal ("Choose 1") card. The attack-turn runner's modal-tuple
	// enumeration writes Mode before each permutation; Play reads it. Sized int8 so it
	// packs into the bool block's padding.
	Mode int8
	// Ephemeral groups every field whose value can change during this card's attack-step
	// resolution: granted keywords, bonus accumulators, the pitch-attribution slice, and
	// the on-hit handler queue. Anonymous embedding so callers keep writing pc.GrantedGoAgain
	// / pc.BonusAttack / etc. directly; struct literals must address the group explicitly
	// (e.g. CardState{Card: c, Ephemeral: Ephemeral{GrantedGoAgain: true}}).
	Ephemeral
}

// Ephemeral is the short-lived state attached to a CardState — granted keywords, bonus
// accumulators, registered on-hit handlers, and pitch attribution that build up around a
// card's play and don't survive the slot being reused. Reset zeroes the whole struct in a
// single statement; every new mutable field added here gets cleared for free.
type Ephemeral struct {
	// GrantedGoAgain is set by a prior card's grant ("next X attack" riders) or by the
	// card's own Play flipping self.GrantedGoAgain = true. Card.EffectiveGoAgain ORs this
	// with Card.GoAgain().
	GrantedGoAgain bool
	// GrantedDominate is the Dominate counterpart to GrantedGoAgain. EffectiveDominate ORs
	// this with the card's Dominator marker.
	GrantedDominate bool
	// GrantedOverpower flags an attack that has gained the Overpower keyword this turn.
	// The engine doesn't currently consume Overpower (the partition's incoming-damage /
	// block model accounts for it), so the flag is the rules-text record for cards that
	// read "this has Overpower".
	GrantedOverpower bool
	// GrantedInstant marks a card a prior effect lets you play "as though it were an
	// instant" — it pays no action point this turn, exactly like an innate Instant.
	GrantedInstant bool
	// BonusAttack is the +{p} this card has accumulated from "next attack +N{p}" riders or
	// self-riders. EffectiveAttack folds it into hit-likelihood checks — a +N buff can bump
	// a 4-power attack into the 5+ dominate window. Negative grants (defender-side debuffs)
	// clamp at 0 because FaB attack power can't go below 0.
	BonusAttack int
	// BonusDefense is the defender-side counterpart to BonusAttack. EffectiveDefense folds
	// it into the attack step's (+N). Negative grants clamp at 0.
	BonusDefense int
	// PitchedToPlay is the pitched cards the attack-turn runner attributed to paying this card's
	// resource cost during the active permutation. Cards whose printed text gates on "if X
	// was pitched to play this" iterate this slice instead of the unordered g.Pitched bag.
	// Each entry is the physical pitched-zone copy (carries in-attack-turn state). Empty
	// for cards whose cost was fully paid by carry from a prior pitch.
	PitchedToPlay []*CardState
	// FaceUp is set to true when GameEngine.TurnFaceUp flips this scheduled card's CardState
	// face-up earlier in the attack turn. Lives in Ephemeral because we don't yet model FaceUp
	// across the turn boundary (no card consumes the persistent flag); cards reading
	// self.FaceUp see the in-attack-turn TurnFaceUp result.
	FaceUp bool
	// PaidAlternativeCost is set by the attack-turn runner when this card was paid via its
	// AlternativeCost branch (Moon Wish's "put a card on top of your deck rather than pay
	// {r}"). Card.Play reads it to decide whether the alt-cost side effect (MoveFromHandToTopOfDeck)
	// should run. Cards without AlternativeCost never see this set.
	PaidAlternativeCost bool
	// OnHit holds "if this hits" handlers registered during Play. Stored as struct values
	// (function pointer + small data payload) rather than closures so registration is
	// alloc-free.
	OnHit []OnHitHandler
}

// Reset zeroes r, preserving OnHit's backing array so per-Best reuse stays allocation-free.
// PitchedToPlay's backing is owned by the attack-turn runner's pitch pool; nilling it here drops
// the alias so the next resolution receives a fresh contrib slice.
func (r *Ephemeral) Reset() {
	onHit := r.OnHit[:0]
	*r = Ephemeral{OnHit: onHit}
}

// OnHitHandler is one registered on-hit rider on a CardState. The attack-turn runner fires Fire
// at finalize-active-attack time when LikelyToHit(self) is true; self is the buffed attack's
// CardState. Source names the card that registered the handler — needed when the handler
// was attached to a different card's OnHit. N and LogText are optional small payloads cards
// use to avoid closures.
type OnHitHandler struct {
	Fire    func(g GameEngine, l Logger, self *CardState, h *OnHitHandler)
	Source  Card
	LogText string
	N       int
}

// RegisterOnHit appends a fire-only on-hit handler. Cards needing Source / N / LogText
// payloads append an OnHitHandler literal directly.
func (p *CardState) RegisterOnHit(fire func(g GameEngine, l Logger, self *CardState, h *OnHitHandler)) {
	p.OnHit = append(p.OnHit, OnHitHandler{Fire: fire})
}

// GrantAttackReactionBuff buffs the active attack target by n: adds to BonusAttack, credits
// g's value, amends the target's attack-step delta, and logs the rider under the target's
// entry. p is the Attack Reaction card granting the buff.
func (p *CardState) GrantAttackReactionBuff(g GameEngine, l Logger, n int) {
	target := g.AttackReactionTarget()
	if target == nil {
		return
	}
	target.BonusAttack += n
	g.AddValue(n)
	l.AmendLastAttackStepN(n)
	l.AppendPostTriggerf(target.Card.DisplayName(), 0, "%s buffed +%d{p}", p.Card.DisplayName(), n)
}

// EffectiveGoAgain reports whether this card has Go again this turn — from printed text or
// a grant by a prior card's effect. g is forwarded to Card.GoAgain so hero-conditional
// cards can read g.HeroWantsLowerHealth.
func (p *CardState) EffectiveGoAgain(g GameEngine) bool {
	return p.Card.GoAgain(g) || p.GrantedGoAgain
}

// EffectiveCost returns the card's live paid resource cost: VariableCost cards return
// their EffectiveCost(g); everyone else falls through to Card.Cost() (the printed cost).
// This is the value the sim's resource-pool drains against — never the printed cost for
// a VariableCost card. Predicates that read "a card with cost N" (Flock's reveal, etc.)
// use Card.Cost() directly so the read stays cache-key-safe; EffectiveCost is for the
// attack-turn runner's payment site.
func (p *CardState) EffectiveCost(g GameEngine) int {
	if vc, ok := p.Card.(VariableCost); ok {
		return vc.EffectiveCost(g)
	}
	return p.Card.Cost()
}

// GrantGoAgainIfFromArsenal flips p.GrantedGoAgain when this copy came from the arsenal
// slot. Safe to call unconditionally at the top of any Play whose printed text reads
// "If <Self> is played from arsenal, it gains go again."
func (p *CardState) GrantGoAgainIfFromArsenal() {
	if p.FromArsenal {
		p.GrantedGoAgain = true
	}
}

// EffectiveTypes returns the type-line this CardState resolves with for its chosen mode.
// Cards that implement ModalTypes (per-mode type-line) get the mode-dispatched TypeSet;
// every other card falls back to the static Card.Types(g). Use this whenever a predicate
// reads types from a *CardState — buffers scanning CardsRemaining, predicates inspecting
// the actively-resolving card, etc. — so the mode-aware result lands.
func (p *CardState) EffectiveTypes(g GameEngine) TypeSet {
	if mt, ok := p.Card.(ModalTypes); ok {
		return mt.TypesForMode(g, p.Mode)
	}
	return p.Card.Types(g)
}

// EffectiveAttack returns Attack() + BonusAttack, clamped at 0. An attack's power can't
// be reduced below 0 in FaB. Cards with "if this hits" clauses pass this into LikelyToHit
// so the rider fires off the post-clamp value.
func (p *CardState) EffectiveAttack() int {
	n := p.Card.Attack() + p.BonusAttack
	if n < 0 {
		return 0
	}
	return n
}

// EffectiveDefense returns Defense() + BonusDefense + ArsenalDefenseBonus (when this copy
// came from arsenal), clamped at 0. Read by ResolveAttackStep to credit the DR's attack-step
// (+N) and bank the block against IncomingDamage.
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

// EffectiveDominate reports whether this card attacks with Dominate this turn — from its
// printed Dominator marker or a GrantedDominate flip.
func (p *CardState) EffectiveDominate() bool {
	return p.GrantedDominate || HasDominate(p.Card)
}

// HasDominate reports whether c is printed with the Dominate keyword. Used by
// EffectiveDominate and by scanners that need the static check without a CardState.
func HasDominate(c Card) bool {
	_, ok := c.(Dominator)
	return ok
}

// arsenalDefenseBonusOf returns c's ArsenalDefenseBonus contribution, or 0 when c doesn't
// implement the marker. Callers gate on their own from-arsenal predicate before invoking —
// the helper does NOT decide whether the bonus applies, only how to extract it.
func arsenalDefenseBonusOf(c Card) int { return ArsenalDefenseBonusOf(c) }

// ArsenalDefenseBonusOf is the exported variant for callers outside this package.
func ArsenalDefenseBonusOf(c Card) int {
	if ab, ok := c.(ArsenalDefenseBonus); ok {
		return ab.ArsenalDefenseBonus()
	}
	return 0
}
