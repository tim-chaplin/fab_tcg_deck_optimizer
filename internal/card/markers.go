package card

// Optional marker interfaces cards opt into to layer behaviour onto the base Card
// contract. The attack-turn runner type-asserts on these at the matching point in the
// pipeline (cost, partition, defense, …) and skips the branch when the card doesn't
// implement the marker.

// VariableCost is optionally implemented by cards whose Cost(g) varies with the engine
// state (e.g. discount-per-token effects). MinCost and MaxCost are static bounds; the
// solver uses them for cheap pre-screens before enumerating attack-turn permutations.
// Non-implementers must return the same value for Cost(g) regardless of g.
type VariableCost interface {
	MinCost() int
	MaxCost() int
}

// Modal is the marker for "Choose 1" cards. Modes returns the number of exclusive modes
// (typically 2); cards dispatch on self.Mode inside Play. See docs/dev-standards.md
// `Modal "Choose 1" cards` for the wiring contract.
type Modal interface {
	Modes() int
}

// ModalCost is an optional add-on to Modal for cards whose resource cost varies by mode.
// Implementers return the cost paid when self.Mode equals the given mode index.
type ModalCost interface {
	ModalCost(mode int8) int
}

// ModalTypes is an optional add-on to Modal for cards whose type-line varies by mode (e.g.
// Tip-Off: Generic Action - Attack in the printed-attack mode, Generic Instant in the
// Instant-discard mode). Implementers return the TypeSet for the given mode index. Cards
// reading another card's type via CardState should use CardState.EffectiveTypes so the
// mode-aware dispatch lands; the attack-turn runner's cardmeta cache stores a per-mode slice
// and dispatches the AP / mark / attack-clear gates on the live mode.
type ModalTypes interface {
	TypesForMode(g GameEngine, mode int8) TypeSet
}

// PlayPrecondition is a marker for cards whose printed text imposes a non-resource
// additional cost beyond Cost(). Implementers return false when THIS play can't legally
// happen (e.g. "reveal a card in your hand with cost 2 or greater" with no eligible
// target); the attack-turn runner rejects the permutation and Play is not called. The check runs
// after the attack-turn runner has removed the playing card and popped this card's funding
// pitches, so scans see only cards that genuinely remain in hand.
type PlayPrecondition interface {
	PlayPrecondition(g GameEngine, self *CardState) bool
}

// Blocker is an optional interface for plain-block cards with block-time logic. Block runs
// once for the card during defense resolution, before its block is counted. Implementations
// scan g.Defenders() and flip self.BonusDefense for conditional buffs ("+N{d} when defending
// alone", "+N{d} together with another card"), and may apply game-state side effects to pay
// a block-time cost (e.g. discarding a card for +N{d}). Cards without block-time logic don't
// need to implement Blocker; their plain-block contribution stays at the printed Defense().
type Blocker interface {
	Block(g GameEngine, l Logger, self *CardState)
}

// BlockCost is an optional add-on for Modal Blockers whose block-time bonus comes with a
// per-mode resource cost. Mode-0 cost is conventionally 0 — the printed "default" branch
// with no extra resources spent.
type BlockCost interface {
	BlockCost(mode int8) int
}

// DefensiveInstant marks a TypeInstant card whose printed effect prevents damage during the
// defense phase. Opting in routes the card through the Defense-Reaction partition slot;
// ResolveAttackStep treats it like a DR. Cards whose prevention is gated by hidden state
// (arcane-only, multi-source rationing) must not opt in — the marker promises a full
// single-bucket prevention. See docs/dev-standards.md: DefensiveInstant markers.
type DefensiveInstant interface {
	DefensiveInstant()
}

// Dominator is the marker the Dominate keyword maps to. EffectiveDominate ORs it with
// GrantedDominate.
type Dominator interface {
	Dominate()
}

// ArsenalDefenseBonus is the marker EffectiveDefense consults when this copy came from the
// arsenal slot. Returns the additional defense added to Defense(). Defense() itself stays
// the printed value so the hand-played path is unaffected.
type ArsenalDefenseBonus interface {
	ArsenalDefenseBonus() int
}

// ResourceSource is the marker for a card that adds resource points beyond printed pitch.
// MaxResourcePoints is a static upper bound on what it adds this turn; the attack-budget
// prune relaxes by it so such a card is never wrongly pruned. A lint flags card files that
// call GameEngine.AddResourcePoints without implementing this.
type ResourceSource interface {
	MaxResourcePoints() int
}

// Universal is the marker for cards whose printed type-line is "extended" by the active
// hero's class. Universal cards ask the engine for the active hero's class inside their
// own Types(g) body.
type Universal interface {
	Universal()
}

// AttackReaction is implemented by every Attack Reaction card. ARTargetAllowed reports
// whether c is a legal target for this AR's chosen mode. Non-modal ARs ignore the mode
// parameter (it's always 0). Modal ARs (card.Modal) dispatch on it: each mode's printed
// target text becomes its own predicate leg, and the attack-turn runner rejects the
// permutation when the chosen mode doesn't accept the active attack.
//
// The engine handle is threaded through so predicates that read variable cost (c.Cost(g))
// or class-aware types (c.Types(g)) don't have to fabricate a zero TurnState. Most ARs
// look only at printed type-line predicates and ignore g.
//
// CardState.GrantAttackReactionBuff (the method most ARs call from Play) lives alongside
// this interface in internal/card — it's pure GameEngine / Logger / CardState plumbing.
type AttackReaction interface {
	ARTargetAllowed(ge GameEngine, c Card, mode int8) bool
}

// LeavesArenaAura is the optional marker for an Aura card with a printed "when this leaves
// the arena" clause. The engine calls OnLeavesArena when the card's aura is destroyed —
// whether by the aura's own start-of-action-phase self-destroy or by another effect that
// destroys it. The aura's trigger handler need only destroy the aura; the leave payoff
// lives in OnLeavesArena so it fires on every path out of the arena.
type LeavesArenaAura interface {
	OnLeavesArena(g GameEngine, l Logger)
}

// FaceUpHook is the optional marker for cards that react to being turned face up by an
// effect (e.g. "If an activated ability or action card effect puts <self> face up into a
// zone from your deck, gain 1 action point"). The engine calls OnFaceUp from
// GameEngine.TurnFaceUp after setting the target CardState's FaceUp flag. Implementers
// shouldn't gate on FaceUp themselves — TurnFaceUp only fires the hook when the flip
// happens, not on every Play.
type FaceUpHook interface {
	OnFaceUp(g GameEngine, l Logger)
}
