// Package card defines the Card interface every Flesh and Blood card implements, the
// per-chain-step CardState wrapper that carries mutable flags between resolution phases,
// and the narrow GameEngine / Logger / Aura / Trigger interfaces cards consume from the sim.
//
// The package owns the contract; it does NOT import the sim. gameengine.GameEngine and
// gameengine.NoopLogger satisfy these interfaces structurally.
package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// GameEngine is the cards-facing rules-engine handle the sim threads through every Card
// hook. The surface combines state queries (Hand, Graveyard, Runechants, …) with active-
// effect operations (DrawOne, CreateRunechants, AddValue, Opt, Clash, …). *sim.TurnState
// satisfies it structurally.
//
// The interface is method-only: it exposes no fields. Card code that needs raw field access
// type-asserts back to *sim.TurnState and lives with the sim import — the exception, not
// the rule.
type GameEngine interface {
	// Zones
	Hand() []Card
	AppendHand(Card)
	PopHandAt(int) Card
	PeekDeck() (Card, bool)
	PeekTopN(int) []Card
	PrependToDeck(Card)
	AppendToDeck(Card)
	AddToGraveyard(Card)
	// Discard pops the first hand card and appends it to the graveyard. Returns the
	// discarded card and true; returns (nil, false) when the hand is empty.
	Discard() (Card, bool)

	// Auras: per-trigger-type registration. Cards supply the handler and initial count;
	// the engine builds the underlying aura. Source is derived from pc.Card. Handler
	// signatures are inlined to keep this package import-free of the concrete aura type.
	CreateStartOfTurnAura(pc *CardState, handler func(GameEngine, Logger, Aura), count int)
	CreateOncePerTurnAttackActionAura(pc *CardState, handler func(GameEngine, Logger, Aura), count int)
	// CreateHitOrDamageTakenAura registers an aura that fires when an attack hits or when
	// the defense phase ends with damage unblocked. filter narrows the hit side (nil = any
	// hit); it never gates the damage-taken side.
	CreateHitOrDamageTakenAura(pc *CardState, handler func(GameEngine, Logger, Aura), count int, filter func(TypeSet) bool)
	// DestroyAura removes the aura currently being fired. addToGraveyard sends the
	// originating card to the graveyard (token auras skip the append). Reached via the
	// per-fire ctx's Destroy method; exposed on GameEngine so the ctx can route the call
	// through its stored engine reference.
	DestroyAura(addToGraveyard bool)

	// Triggers: one-shot, per-trigger-type. AddHitTrigger's filter narrows the firing event
	// to a card-type predicate (typically TypeSet.IsAttack or TypeSet.IsAttackAction); nil
	// = no filter. Handler signatures are inlined for the same reason as the aura methods.
	AddHitTrigger(pc *CardState, handler func(GameEngine, Logger, Trigger), filter func(TypeSet) bool)
	AddEndOfTurnTrigger(pc *CardState, handler func(GameEngine, Logger, Trigger))
	// AddCardOrAbilityTrigger fires a one-shot listener when the next card is played in
	// the chain; filter narrows the firing card (nil = any). The registering card never
	// fires it — the event resolves before a card's own effect, so the trigger isn't
	// queued yet when its own card resolves.
	AddCardOrAbilityTrigger(pc *CardState, handler func(GameEngine, Logger, Trigger), filter func(TypeSet) bool)

	// Token economy
	CreateRunechants(int)
	CreatePonders(int)
	CreateGold(int)
	CreateCopper(int)
	RunechantCount() int
	GoldCount() int
	SilverCount() int
	CopperCount() int

	// Value crediting and arcane damage. AddValue accepts negatives — Test of Strength's
	// clash-loss concedes value to the opponent.
	AddValue(int)
	// DealArcaneDamage credits n arcane damage from the named source: adds to Value, flips
	// ArcaneDamageDealt when the hit lands, and logs the rider under source. Pass
	// self.Card.DisplayName() from Play, or a.CardName() from an aura handler.
	DealArcaneDamage(l Logger, source string, n int)

	// Hit / damage heuristics. The card says "I attack for N"; the engine decides whether
	// the attack lands.
	LikelyToHit(pc *CardState) bool
	LikelyDamageHits(n int, dominate bool) bool

	// Tempo verbs. Cards say "this card does X" (force a discard); the engine decides how
	// much that's worth in damage-equivalent Value and credits it. The verb returns the
	// value it credited so cards can attribute the rider line.
	OpponentDiscard(n int) int

	// AP (chain-step controls cards grant).
	AddActionPoints(int)

	// Sticky flags cards read to gate their effects. Setters are implicit side effects of
	// engine verbs: DealArcaneDamage flips ArcaneDamageDealt; Create*Aura flips AuraCreated;
	// MarkOpponent flips OpponentMarked. NonAttackActionPlayed and CardBanished are
	// sim-side bookkeeping.
	ArcaneDamageDealt() bool
	AuraCreated() bool
	CardBanished() bool
	NonAttackActionPlayed() bool
	OpponentMarked() bool
	MarkOpponent()
	// LastAttackHit reports whether the most recent finalised attack on this combat chain
	// hit. False until the first attack finalises; each subsequent attack overwrites it.
	LastAttackHit() bool
	// IsMyTurn reports whether the active phase is the owning player's action phase (true)
	// or the defense phase (false) — backs "during your turn" riders.
	IsMyTurn() bool

	// Partition / matchup state.
	RemainingUnblockedDamage() int
	ArcaneIncomingDamage() int
	BlockTotal() int
	Defenders() []Card
	Pitched() []Card

	// Hero info. HeroWantsLowerHealth reports whether the current hero opts into the
	// LowerHealthWanter marker (proxy for "less {h} than the opponent" riders).
	// CurrentHeroClass returns the hero's primary class; Universal cards fold this into
	// their own Types(g) so class-gated triggers see the right type-line.
	HeroWantsLowerHealth() bool
	CurrentHeroClass() CardType

	// Chain queries
	HasPlayedType(CardType) bool
	CardsPlayed() []Card
	SetCardsPlayed([]Card)
	CardsRemaining() []*CardState
	TriggeringCard() Card
	// AuraCount is the count of live auras — used by "while you control an aura" gates.
	// Cards don't get a typed slice view; the engine owns the live aura set.
	AuraCount() int

	// Mid-chain draw / tutoring / recycling — cards moving to different zones.
	DrawOne() bool
	TutorFromDeck(func(Card) int) (Card, bool)
	RecycleToDeckBottom(*CardState)
	RecycleFromGraveyardToTop(func(Card) bool) (Card, bool)
	RecycleFromGraveyardToBottom(func(Card) bool) (Card, bool)
	BanishFromGraveyard(func(Card) bool) (Card, bool)

	// Opt (hero-driven top-of-deck reshape)
	Opt(Logger, int)

	// Clash (top-of-deck power compare)
	Clash(win, lose func())

	// PlayCard runs Card.Play on pc and emits the chain step. Used by cards that resolve
	// another card mid-handler (Moon Wish tutoring Sun Kiss into play on go-again).
	PlayCard(Logger, *CardState)

	// Attack reaction target accessor.
	// TODO: pass the target directly into the AttackReaction's Play instead of leaking the
	// "engine remembers the active attack" detail through this accessor.
	AttackReactionTarget() *CardState
}

// Logger is the cards-facing log sink the chain runner threads through every Card hook.
//
// Cards use AppendPostTrigger / AppendPostTriggerf for self-riders ("Created a runechant"),
// AppendPreTrigger for hero or aura attack-action triggers. The AppendChainStep /
// AppendChainStepf / AmendLastChainStepN trio is sim-internal — used by ResolveChainStep,
// the Opt helper, and GrantAttackReactionBuff — but lives on Logger so both call sites
// share the same value.
type Logger interface {
	AppendChainStep(text string, n int)
	AppendChainStepf(n int, format string, args ...any)
	AppendPostTrigger(source, text string, n int)
	AppendPostTriggerf(source string, n int, format string, args ...any)
	AppendPreTrigger(source, text string, n int)
	AppendPreTriggerf(source string, n int, format string, args ...any)
	AmendLastChainStepN(n int)
}

// Aura is the minimal view cards' aura handlers see of the firing aura. The handler reads
// counts / source-card identity through this interface and ends the aura's life via Destroy.
type Aura interface {
	// Count returns the aura's current count. Fire-then-destroy handlers read it as a
	// payload; multi-fire auras read it as fires-remaining.
	Count() int
	// DecrementCount decrements the count by 1 and returns the new value. Multi-fire aura
	// handlers call this and Destroy when the result reaches 0.
	DecrementCount() int
	// CardName is the originating card or token's display name — used for log attribution.
	CardName() string
	// CardID is the originating card's registry ID, or ids.InvalidCard when the aura
	// belongs to a token. Handlers that match against a specific source card gate on this.
	CardID() ids.CardID
	// Destroy ends the aura's life. addToGraveyard sends the originating card to the
	// graveyard (token auras with no originating card skip the append).
	Destroy(addToGraveyard bool)
}

// Trigger is the minimal view cards' one-shot trigger handlers see of the firing
// trigger. Today only CardName is needed; expand as more triggers appear.
type Trigger interface {
	CardName() string
}

// Hero is the narrow view of the active hero internal/card needs (Class / Types only) so
// card-side helpers can answer "what's the current hero's class?" without importing sim.
type Hero interface {
	Class() CardType
	Types() TypeSet
}
