// Package card defines the Card interface every Flesh and Blood card implements, the
// per-chain-step CardState wrapper that carries mutable flags between resolution phases,
// and the narrow GameEngine / Logger / Aura / EphemeralTrigger interfaces cards consume from
// the sim.
//
// The package owns the contract; it does NOT import the sim. gameengine.GameEngine and
// gameengine.NoopLogger satisfy these interfaces structurally.
package card

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
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
	// HandSize reports the total number of cards in hand, including drawn cards. Prefer
	// this over len(Hand()) for emptiness / size gates: it doesn't flip IsCacheable.
	HandSize() int
	// HandHasMatching reports whether any non-drawn hand entry satisfies pred. Drawn
	// entries are skipped — their identity is opaque mid-chain. Doesn't flip IsCacheable.
	HandHasMatching(pred func(Card) bool) bool
	// HeldHandSize reports the total Held-role entry count (Pitch / Attack excluded),
	// including drawn entries. Counting alone doesn't reveal drawn-card attributes, so
	// this accessor doesn't flip IsCacheable.
	HeldHandSize() int
	// HeldHand returns the Held subset of the hand (Pitch / Attack excluded, drawn
	// entries included). Flips IsCacheable since iterating the slice exposes attributes.
	HeldHand() []Card
	AppendHand(Card)
	PeekDeck() (Card, bool)
	PeekTopN(int) []Card
	PrependToDeck(Card)
	AppendToDeck(Card)
	AddToGraveyard(Card)
	// Discard pops the first Held card to the graveyard and logs under source.
	// Returns true on success; false when no Held card exists. Cache-safe — the
	// discarded identity never escapes the engine.
	Discard(source string) bool
	// DiscardToTopOfDeck pops the first Held card to the top of the deck and logs
	// under source. Returns true on success. Cache-safe.
	DiscardToTopOfDeck(source string) bool
	// DiscardToBottomOfDeck pops the first Held card to the bottom of the deck and
	// logs under source. Returns true on success. Cache-safe.
	DiscardToBottomOfDeck(source string) bool

	// CreateAura registers a multi-fire aura sourced from source (typically the playing
	// card's self.Card) that fires on every event in tt's bit set. oncePerTurn caps it to
	// one fire per turn; filter narrows the firing site to a card-type predicate (nil =
	// any) and is consulted only on events that carry a triggering card (so StartOfTurn,
	// EndOfTurn, and DamageTaken effectively ignore it). Handler signatures are inlined to
	// keep this package import-free of the concrete aura type.
	CreateAura(source Card, tt triggertype.Type, handler func(GameEngine, Logger, Aura), count int, oncePerTurn bool, filter func(TypeSet) bool)
	// DestroyAura removes the aura currently being fired. addToGraveyard sends the
	// originating card to the graveyard (token auras skip the append). Reached via the
	// per-fire ctx's Destroy method; exposed on GameEngine so the ctx can route the call
	// through its stored engine reference.
	DestroyAura(addToGraveyard bool)
	// DestroyItem removes the item currently being fired. The item counterpart of
	// DestroyAura — reached via the firing item's Destroy method.
	DestroyItem(addToGraveyard bool)
	// CreateItem puts a card-sourced item into play whose handler fires on every event in
	// tt's bit set. oncePerTurn caps it to one fire per turn; filter narrows the firing
	// site (nil = any).
	CreateItem(source Card, tt triggertype.Type, handler func(GameEngine, Logger, Item), oncePerTurn bool, filter func(TypeSet) bool)
	// AddResourcePoints adds n resources to the card currently being pitched — a Pitch handler
	// calls it to boost what that pitched card yields. No effect outside a pitch fire.
	AddResourcePoints(n int)
	// SacrificePayoffAura destroys one aura the player controls, reporting whether one was
	// destroyed. See GameEngine.SacrificePayoffAura for targeting rules.
	SacrificePayoffAura() bool

	// CreateTrigger registers a one-shot ephemeral trigger that fires once on the next
	// event in tt's bit set and is then dropped. filter narrows the firing event to a
	// card-type predicate (nil = any); it is consulted only when the triggering event has
	// a triggering card. A trigger registered from a card's own Play does not fire for
	// its own resolution — the CardOrAbility event has already resolved by then.
	CreateTrigger(source Card, tt triggertype.Type, handler func(GameEngine, Logger, EphemeralTrigger), filter func(TypeSet) bool)

	// Token economy
	CreateRunechants(int)
	CreatePonders(int)
	CreateGold(int)
	CreateSilver(int)
	CreateCopper(int)
	RunechantCount() int
	GoldCount() int
	SilverCount() int
	CopperCount() int

	// Status-token mints under the opponent's control. Opposing-side state isn't tracked
	// — these credit a flat heuristic value via AddValue instead of touching state.
	CreateFrailtyForOpponent()
	CreateInertiaForOpponent()
	CreateBloodrotPoxForOpponent()

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

	// PreventArcaneDamage caps incoming arcane damage by up to n, returning the amount
	// actually prevented (clamped at the remaining arcane). Mutates ArcaneIncomingDamage
	// so downstream readers see the reduced figure. Callers AddValue the returned amount
	// to credit the prevention.
	PreventArcaneDamage(n int) int

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
	// CrowdCheer / CrowdBoo land a crowd reaction on the active hero: flips the
	// HasCrowdCheered / HasCrowdBooed flag and fires the CrowdCheer / CrowdBoo trigger.
	// Callers gate on the source-side rule (e.g. "each Revered hero") themselves.
	CrowdCheer()
	CrowdBoo()
	// HasCrowdCheered / HasCrowdBooed reports whether a crowd reaction has landed on the
	// active hero this turn.
	HasCrowdCheered() bool
	HasCrowdBooed() bool
	// UntapHero untaps the owning player's hero — the printed "untap your hero" effect.
	UntapHero()
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
	// HeroHasType reports whether the active hero's type line contains t — used by cards
	// that gate on hero-only keywords like TypeRevered / TypeReviled.
	HeroWantsLowerHealth() bool
	CurrentHeroClass() CardType
	HeroHasType(t CardType) bool

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

// EphemeralTrigger is the minimal view cards' one-shot trigger handlers see of the firing
// trigger. Today only CardName is needed; expand as more triggers appear.
type EphemeralTrigger interface {
	CardName() string
}

// Item is the minimal view cards' item trigger handlers see of the firing item. The
// handler reads source-card identity and ends the item's life via Destroy.
type Item interface {
	// CardName is the originating card or token's display name — used for log attribution.
	CardName() string
	// CardID is the originating card's registry ID, or ids.InvalidCard for token items.
	CardID() ids.CardID
	// Destroy ends the item's life. addToGraveyard sends the originating card to the
	// graveyard (token items with no originating card skip the append).
	Destroy(addToGraveyard bool)
}

// Hero is the narrow view of the active hero internal/card needs (Class / Types only) so
// card-side helpers can answer "what's the current hero's class?" without importing sim.
type Hero interface {
	Class() CardType
	Types() TypeSet
}
