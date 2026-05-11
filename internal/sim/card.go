// Package card defines the Card interface used by the simulator and basic/test implementations.
//
// The per-card CardState wrapper, the Card interface itself, and the optional markers cards
// opt into (VariableCost, Dominator, ArsenalDefenseBonus, …) live in this file. Cohesive
// concern groups are split across sibling files in this package: types.go (card.CardType +
// card.TypeSet bitfield), turn_state.go (TurnState and its mutation helpers), triggers.go
// (Aura).
package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

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
	// Mode is the chosen mode for a ModalCard ("Choose 1" cards), set by the chain runner
	// before Play. Always 0 for non-modal cards. Sized int8 so it packs into the bool block's
	// padding without growing the CardState — every chain step reads pcBuf[i], so a wider
	// Mode would push more cache lines through the inner loop.
	Mode int8
	// BonusAttack is the +{p} this card has accumulated from prior cards' "next attack +N{p}"
	// riders. Granters set pc.BonusAttack += N on the matching CardState in CardsRemaining so
	// the damage is attributed to the attack receiving the buff, and EffectiveAttack folds it
	// into hit-likelihood checks (LikelyToHit) — a +N buff can bump a 4-power attack into the
	// 5+ dominate window or a 6 into the unblockable 7. The solver applies BonusAttack to
	// every CardState's contribution unconditionally; deciding which CardStates are legal
	// targets (attack actions, weapons, future card types) is the grantor's job. Negative
	// bonuses (defender-side -N{p} debuffs) clamp at 0 because FaB attack power can't go
	// below 0.
	BonusAttack int
	// BonusDefense is the +{d} this card has accumulated from "+N{d}" rider clauses, the
	// defender-side counterpart to BonusAttack. Cross-card grants from other cards and self-
	// riders ("if X, this gains +1{d}") both write into this field; EffectiveDefense folds it
	// into the chain step's (+N) so a buffed DR's block reflects the grant. Negative grants
	// clamp at 0.
	BonusDefense int
	// PitchedToPlay is the pitched cards the chain runner attributed to paying this card's
	// resource cost during the active permutation. Populated by the chain runner before each
	// Card.Play: as costs come up, pitched cards are popped from the active pitch ordering
	// (carrying over any excess to fund subsequent cards) and the popped slice is exposed
	// here. Cards whose printed text gates on "if X was pitched to play this" iterate this
	// slice instead of the unordered s.Pitched bag — the same pitched bag still lives on
	// TurnState for cards that read it as a multiset. Empty for cards whose cost was fully
	// paid by carry from a prior pitch.
	PitchedToPlay []Card
	// OnHit holds "if this hits" handlers registered during Play. Stored as struct values
	// (function pointer + small data payload) rather than closures so the hot anneal path
	// doesn't allocate per registration. See docs/dev-standards.md "OnHit registrations"
	// for the wiring contract.
	OnHit []OnHitHandler
	// SkipGraveyard is set by Play helpers that route this card to a non-graveyard zone —
	// e.g. TurnState.RecycleToDeckBottom for Relentless Pursuit's "put this on the bottom
	// of its owner's deck" clause. The chain dispatcher's "non-persistent → graveyard"
	// append checks this flag and skips the append when set, so the helper that moved
	// the card owns its destination zone fully.
	SkipGraveyard bool
}

// OnHitHandler is one registered on-hit rider on a CardState. The chain runner fires Fire
// at finalize-active-attack time when LikelyToHit(self) is true; self is the buffed
// attack's CardState. Source names the card that registered the handler so log attribution
// stays correct when the handler was added to a different card's OnHit (Mauvrion Skies,
// Runic Reaping). N and LogText are optional small payloads cards use to avoid closures.
type OnHitHandler struct {
	Fire    func(s *TurnState, l Logger, self *CardState, h *OnHitHandler)
	Source  Card
	LogText string
	N       int
}

// RegisterOnHit appends a fire-only on-hit handler — the common case for "if this hits, do
// X" riders. Cards needing Source / N / LogText payloads on the handler append an
// OnHitHandler literal directly.
func (p *CardState) RegisterOnHit(fire func(s *TurnState, l Logger, self *CardState, h *OnHitHandler)) {
	p.OnHit = append(p.OnHit, OnHitHandler{Fire: fire})
}

// EffectiveGoAgain reports whether this card has Go again this turn — from printed text or a
// grant by a prior card's effect.
func (p *CardState) EffectiveGoAgain() bool {
	return p.Card.GoAgain() || p.GrantedGoAgain
}

// GrantGoAgainIfFromArsenal flips p.GrantedGoAgain when this copy came from the arsenal
// slot (p.FromArsenal). Names the standard "played-from-arsenal go again" rider — see the
// docs/dev-standards.md "Played-from-arsenal go-again" entry — so card Play bodies don't
// need to spell out the three-line if. No-op when FromArsenal is false; safe to call
// unconditionally at the top of any Play whose printed text reads "If <Self> is played
// from arsenal, it gains go again."
func (p *CardState) GrantGoAgainIfFromArsenal() {
	if p.FromArsenal {
		p.GrantedGoAgain = true
	}
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

// EffectiveDefense returns the card's printed Defense() plus any granted BonusDefense plus
// the ArsenalDefenseBonus when this copy came from the arsenal slot, clamped at 0. The
// sim reads this in ResolveChainStep to credit the DR's chain-step "(+N)" delta and to
// decrement IncomingDamage.
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

// Logger is the cards-facing log sink the chain runner threads through Play, Block,
// OnHitHandler.Fire, TriggerHandler, and Hero.OnCardPlayed. The method set matches
// *turnlogger.TurnLogger exactly so sim hands one of those directly to cards — no
// adapter sits in the middle. The interface lives in card.go so it travels with Card
// when the package later moves to v2/card; the concrete implementation
// (turnlogger.TurnLogger) stays in v2/turnlogger and the move requires no import
// rewiring on the card side. Cards convert CardState into log-line text at the call
// site via ChainStepText / Card.DisplayName.
type Logger interface {
	// AppendChainStep appends a main-line chain-step entry with the given text and
	// damage-equivalent display delta n. Cards typically wrap the call as
	// l.AppendChainStep(ChainStepText(self), n).
	AppendChainStep(text string, n int)
	// AppendChainStepf is the format variant of AppendChainStep — fmt.Sprintf runs only
	// on the recording branch.
	AppendChainStepf(n int, format string, args ...any)
	// AppendPostTrigger appends an indented post-trigger sub-line attributed to source.
	// Self-riders pass self.Card.DisplayName() as source; cross-card riders (OnHit
	// handlers attached to a target card) pass the target's name.
	AppendPostTrigger(source, text string, n int)
	// AppendPostTriggerf is the format variant of AppendPostTrigger.
	AppendPostTriggerf(source string, n int, format string, args ...any)
	// AppendPreTrigger appends an indented pre-trigger sub-line attributed to source —
	// a hero or aura-attack-action trigger that fires before its parent chain entry.
	AppendPreTrigger(source, text string, n int)
	// AppendPreTriggerf is the format variant of AppendPreTrigger.
	AppendPreTriggerf(source string, n int, format string, args ...any)
	// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. ARs use
	// this to fold their +{p} buff into the buffed attack's display delta.
	AmendLastChainStepN(n int)
}

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static
// profile plus a Play hook for on-play logic.
type Card interface {
	// ID returns the card's canonical registry identifier. Stable within a build. Lets callers
	// key maps / slices on cards without string-hashing Name().
	ID() ids.CardID
	// Name returns the card's printed name without any pitch-color suffix — all three
	// printings of "Aether Slash" return the same string. Cards comparing by name
	// (synergies, "if you have played a card named X this turn" effects) use this directly.
	// For display, callers route through DisplayName which appends the pitch tag.
	Name() string
	// DisplayName returns the human-readable identifier including the pitch suffix —
	// "Aether Slash [R]", "Aether Slash [Y]", "Aether Slash [B]". Use this anywhere a
	// printout needs to disambiguate pitch printings (log lines, deck listings, debug).
	DisplayName() string
	// Cost returns the card's current resource cost given the turn state. Cards with a static
	// printed cost ignore s and return a constant; cards that read s (e.g. discount-per-token
	// effects) additionally implement VariableCost so the solver can pre-screen with cheap
	// MinCost / MaxCost bounds before enumerating chain permutations.
	Cost(s *TurnState) int
	Pitch() int
	// Attack is the printed attack value. Conditional bonuses belong in Play, not here.
	Attack() int
	Defense() int
	// Types returns the card's type-line descriptors as a card.TypeSet bitfield, e.g.
	// card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack).
	Types() card.TypeSet
	// GoAgain reports whether playing this card grants an additional action point. Cards
	// printed with "Go again" return true.
	GoAgain() bool
	// Play is called when the card resolves — as an attack or as a defense reaction.
	// Cards own card-specific behaviour only: conditional self-buffs (flip
	// self.BonusAttack / self.BonusDefense), riders (l.AppendPostTrigger sub-lines),
	// OnHit registration, mid-chain effects. The standard "credit EffectiveAttack /
	// capped EffectiveDefense to s.Value and emit the <Card>: <VERB> (+N) chain step"
	// mechanic happens in sim.ResolveChainStep after Play returns — vanilla attack /
	// DR cards have an empty Play body.
	Play(s *TurnState, l Logger, self *CardState)
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

// ModalCard is an optional marker for "Choose 1" cards. Modes returns the number of
// exclusive modes (typically 2); the chain runner enumerates 0..Modes()-1 per ordering and
// cards dispatch on self.Mode inside Play. Modes that are no-ops for the current state
// should resolve as zero-Value no-ops so the runner picks a sibling mode that contributes
// more. See docs/dev-standards.md `Modal "Choose 1" cards` for the wiring contract.
type ModalCard interface {
	Modes() int
}

// ModalCost is an optional add-on to ModalCard for cards whose resource cost varies by
// mode (Bluster Buff: "this gets -1{p} unless you pay {r}" — mode 0 is the printed cost,
// mode 1 spends one more {r}). Implementers return the cost paid when self.Mode equals
// the given mode index. The chain runner reads ModalCost in place of Card.Cost(s) when
// the card declares both ModalCard and ModalCost, and folds min/max of the per-mode
// costs into the partition pre-screen via VariableCost.
type ModalCost interface {
	ModalCost(mode int8) int
}

// Dominator is an optional marker. Attack action cards printed with the Dominate keyword
// implement it; the defender is capped at one blocking card, so LikelyToHit credits the
// "slips past one block" bump at 5+ power. Conditional grants ("if X, it gains dominate")
// stay off this marker and flow through CardState.GrantedDominate instead.
type Dominator interface {
	Dominate()
}

// PlayPrecondition is an optional Card marker for cards whose printed text imposes a
// non-resource additional cost beyond Cost(). Implementers return false when THIS play
// can't legally happen (e.g. Demolition Crew's "reveal a card in your hand with cost 2 or
// greater" with no eligible target); the chain runner rejects the permutation and the
// card's Play is not called. The check runs after the chain runner has removed the
// playing card and popped this card's funding pitches from s.hand, so scans see only
// cards that genuinely remain in hand — a pitch source can't double as a reveal target.
type PlayPrecondition interface {
	PlayPrecondition(s *TurnState, self *CardState) bool
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

// ArsenalDefenseBonus is an optional marker for Defense Reactions whose printed text grants
// extra defense only when the card is played from arsenal. Implementers return the
// additional defense added to Defense() when this copy came from the arsenal slot at start
// of turn. Defense() itself stays the printed value so the hand-played path is unaffected.
type ArsenalDefenseBonus interface {
	ArsenalDefenseBonus() int
}

// arsenalDefenseBonusOf returns c's ArsenalDefenseBonus contribution, or 0 when c doesn't
// implement the marker. Centralises the type assertion behind a single named call so every
// "if this came from arsenal, fold in the rider" site reads as one arithmetic line. Callers
// gate on their own from-arsenal predicate (CardState.FromArsenal, partition arsenal-slot
// index, BestLine.FromArsenal) before invoking — the helper does NOT decide whether the bonus
// applies, only how to extract it.
func arsenalDefenseBonusOf(c Card) int {
	if ab, ok := c.(ArsenalDefenseBonus); ok {
		return ab.ArsenalDefenseBonus()
	}
	return 0
}

// Blocker is an optional interface for plain-block cards that need to react to other
// defenders before contributing their block. The chain runner calls Block on every plain
// blocker that implements it, with TurnState.Defenders populated with the partition's
// full defender slice (DRs + plain blocks). Implementations typically scan Defenders and
// flip self.BonusDefense for "+N{d} when defending alone" / "+N{d} when defending with
// another card" / similar conditional buffs. Cards without block-time logic don't need
// to implement Blocker; their plain-block contribution stays at the printed Defense().
type Blocker interface {
	Block(s *TurnState, l Logger, self *CardState)
}

// BlockCost is an optional add-on for ModalCard Blockers whose block-time bonus comes
// with a per-mode resource cost (Brothers in Arms: "may pay {r} for +2{d}"). The chain
// runner enumerates each modal blocker's modes whose cost fits the partition's spare
// defense budget (defendBudget − drCost) and picks the one that yields the highest
// effective defense. Mode-0 cost is conventionally 0 — the printed "default" branch with
// no extra resources spent.
type BlockCost interface {
	BlockCost(mode int8) int
}

// DefensiveInstant marks a TypeInstant card whose printed effect prevents damage during
// the defense phase. Opting in routes the card through the Defense-Reaction partition
// slot, cost summing, and chain-step Play; the type stays TypeInstant. ResolveChainStep
// treats DefensiveInstant cards like DRs — caps EffectiveDefense at IncomingDamage and
// decrements — so an empty Play body covers the standard prevention case. Cards whose
// prevention is gated by hidden state (arcane-only, multi-source rationing) must not
// opt in — the marker promises a full single-bucket prevention. See
// docs/dev-standards.md: DefensiveInstant markers.
type DefensiveInstant interface {
	DefensiveInstant()
}
