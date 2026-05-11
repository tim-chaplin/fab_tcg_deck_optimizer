// Package sim's card-side surface. The Card interface, CardState wrapper, Logger,
// OnHitHandler, GameEngine handle, and small markers (Dominator,
// ArsenalDefenseBonus) live in v2/card — this file re-exports them as type
// aliases so sibling files in sim and external callers keep referring to
// sim.Card / sim.CardState / sim.Logger etc. without an import switch. Markers
// whose signatures touch *TurnState (PlayPrecondition, Blocker) and sim-internal
// optional shapes (VariableCost, ModalCard, …) live below.
package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Aliases for the cards-facing types. The concrete implementations live in v2/card
// and the chain runner consumes them through these names.
type (
	Card                = card.Card
	CardState           = card.CardState
	Logger              = card.Logger
	OnHitHandler        = card.OnHitHandler
	GameEngine          = card.GameEngine
	Dominator           = card.Dominator
	ArsenalDefenseBonus = card.ArsenalDefenseBonus
)

// HasDominate reports whether c is printed with the Dominate keyword — a type assertion to
// the Dominator marker. Used by CardState.EffectiveDominate and any scanner that needs the
// static printed-keyword check without going through a CardState.
func HasDominate(c Card) bool { return card.HasDominate(c) }

// arsenalDefenseBonusOf returns c's ArsenalDefenseBonus contribution, or 0 when c doesn't
// implement the marker. Used by the partition pre-screen and the log-format pass to fold
// the bonus into a defender's reported (+N) when the copy came from the arsenal slot.
func arsenalDefenseBonusOf(c Card) int { return card.ArsenalDefenseBonusOf(c) }

// VariableCost is optionally implemented by cards whose Cost(g) varies with TurnState (e.g.
// discount-per-token effects). MinCost and MaxCost are static bounds on the Cost output across
// any state; the solver uses them for cheap O(1) pre-screens before enumerating chain
// permutations. Non-implementers must return the same value for Cost(g) regardless of g.
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
// the given mode index. The chain runner reads ModalCost in place of Card.Cost(g) when
// the card declares both ModalCard and ModalCost, and folds min/max of the per-mode
// costs into the partition pre-screen via VariableCost.
type ModalCost interface {
	ModalCost(mode int8) int
}

// PlayPrecondition is an optional Card marker for cards whose printed text imposes a
// non-resource additional cost beyond Cost(). Implementers return false when THIS play
// can't legally happen (e.g. Demolition Crew's "reveal a card in your hand with cost 2 or
// greater" with no eligible target); the chain runner rejects the permutation and the
// card's Play is not called. The check runs after the chain runner has removed the
// playing card and popped this card's funding pitches from s.hand, so scans see only
// cards that genuinely remain in hand — a pitch source can't double as a reveal target.
type PlayPrecondition interface {
	PlayPrecondition(g GameEngine, self *CardState) bool
}

// LowerHealthWanter is an optional Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than their opponent (deck building, sandbagging, self-damage) opt in. Cards with a
// "less {h} than an opposing hero" rider assume the clause always fires for these heroes and never
// fires for anyone else — a coarse proxy that skips per-turn life tracking.
type LowerHealthWanter interface {
	WantsLowerHealth()
}

// Blocker is an optional interface for plain-block cards that need to react to other
// defenders before contributing their block. The chain runner calls Block on every plain
// blocker that implements it, with TurnState.Defenders populated with the partition's
// full defender slice (DRs + plain blocks). Implementations typically scan Defenders and
// flip self.BonusDefense for "+N{d} when defending alone" / "+N{d} when defending with
// another card" / similar conditional buffs. Cards without block-time logic don't need
// to implement Blocker; their plain-block contribution stays at the printed Defense().
type Blocker interface {
	Block(g GameEngine, l Logger, self *CardState)
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
